package worktree

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const recoverySchema = "bench-recovery/v1"

type recoveryManifest struct {
	Schema string            `json:"schema"`
	Base   string            `json:"base"`
	Layers map[string]string `json:"layers"`
}
type indexEntry struct {
	mode, oid, path string
	stage           int
}

// recoverAssignmentWithFault writes every Git-visible layer through temporary indexes. It
// never points HEAD elsewhere and never opens the real index for writing.
func recoverAssignmentWithFault(root string, assignment intent.Assignment, fault Fault) (intent.Assignment, error) {
	if len(assignment.Recovery) > 0 {
		if len(assignment.Recovery) != 1 {
			return assignment, errors.New("existing recovery metadata is ambiguous")
		}
		if err := ensureRecoveryRef(root, assignment, assignment.Recovery[0]); err != nil {
			return assignment, err
		}
		return assignment, nil
	}
	head, err := git.Output("-C", assignment.Worktree, "rev-parse", "HEAD")
	if err != nil {
		return assignment, fmt.Errorf("read recovery HEAD: %w", err)
	}
	headTree, err := git.Output("-C", assignment.Worktree, "rev-parse", "HEAD^{tree}")
	if err != nil {
		return assignment, fmt.Errorf("read recovery base tree: %w", err)
	}
	admin, err := git.Output("-C", assignment.Worktree, "rev-parse", "--path-format=absolute", "--git-dir")
	if err != nil {
		return assignment, err
	}
	layerTrees := map[string]string{}
	workingTree, err := worktreeTree(assignment.Worktree, admin)
	if err != nil {
		return assignment, fmt.Errorf("capture working layer: %w", err)
	}
	if workingTree != headTree {
		layerTrees["working"] = workingTree
	}
	entries, conflicted, err := readIndexEntries(assignment.Worktree)
	if err != nil {
		return assignment, err
	}
	if conflicted {
		for stage, name := range map[int]string{1: "base", 2: "ours", 3: "theirs"} {
			tree, err := conflictTree(assignment.Worktree, admin, entries, stage)
			if err != nil {
				return assignment, fmt.Errorf("capture conflict %s layer: %w", name, err)
			}
			layerTrees[name] = tree
		}
	} else {
		stagedTree, err := realIndexTree(assignment.Worktree, admin)
		if err != nil {
			return assignment, fmt.Errorf("capture staged layer: %w", err)
		}
		if stagedTree != headTree {
			layerTrees["staged"] = stagedTree
		}
	}
	if len(layerTrees) == 0 {
		return assignment, errors.New("recovery requested for a clean assignment")
	}
	treePayload := map[string]string{}
	layers := map[string]string{}
	names := make([]string, 0, len(layerTrees))
	for name := range layerTrees {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		tree := layerTrees[name]
		payload := treePayload[tree]
		if payload == "" {
			payload, err = commitTree(root, tree, []string{head}, "bench recovery payload: "+name+"\n")
			if err != nil {
				return assignment, err
			}
			treePayload[tree] = payload
		}
		layers[name] = payload
	}
	payloads := make([]string, 0, len(treePayload))
	for _, payload := range treePayload {
		payloads = append(payloads, payload)
	}
	sort.Strings(payloads)
	manifest := recoveryManifest{Schema: recoverySchema, Base: head, Layers: layers}
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		return assignment, err
	}
	manifestBytes = append(manifestBytes, '\n')
	blob, err := gitInput(root, nil, manifestBytes, "hash-object", "-w", "--stdin")
	if err != nil {
		return assignment, err
	}
	rootTree, err := gitInput(root, nil, []byte("100644 blob "+blob+"\tmanifest.json\n"), "mktree")
	if err != nil {
		return assignment, err
	}
	rootOID, err := commitTree(root, rootTree, payloads, "bench recovery root\n")
	if err != nil {
		return assignment, err
	}
	ref, err := nextRecoveryRef(root, assignment)
	if err != nil {
		return assignment, err
	}
	recovery := intent.Recovery{Ref: ref, Root: rootOID, Payloads: payloads}
	assignment.Recovery = append(assignment.Recovery, recovery)
	if err := intent.PutAssignment(root, assignment); err != nil {
		return assignment, fmt.Errorf("persist recovery metadata: %w", err)
	}
	if err := hit(fault, StepRecoveryMetadata); err != nil {
		return assignment, err
	}
	if err := ensureRecoveryRef(root, assignment, recovery); err != nil {
		return assignment, err
	}
	return assignment, nil
}

// readRecoveryManifest is the one reader of a recovery envelope. It answers with a
// manifest only when the commit carries a well-formed one of this schema. No caller has
// to decide for itself what a usable envelope is.
func readRecoveryManifest(root, commitish string) (recoveryManifest, bool) {
	body, err := git.Output("-C", root, "show", commitish+":manifest.json")
	if err != nil {
		return recoveryManifest{}, false
	}
	var manifest recoveryManifest
	if json.Unmarshal([]byte(body), &manifest) != nil || manifest.Schema != recoverySchema || manifest.Base == "" || len(manifest.Layers) == 0 {
		return recoveryManifest{}, false
	}
	return manifest, true
}
func recoveryEnvelopeValid(root string, recovery intent.Recovery) bool {
	manifest, ok := readRecoveryManifest(root, recovery.Root)
	if !ok {
		return false
	}
	payloads := map[string]bool{}
	for _, payload := range recovery.Payloads {
		payloads[payload] = true
	}
	seen := map[string]bool{}
	for _, payload := range manifest.Layers {
		if !payloads[payload] {
			return false
		}
		seen[payload] = true
	}
	return len(seen) == len(payloads) && git.OK("-C", root, "cat-file", "-e", manifest.Base+"^{commit}")
}
func temporaryIndex(admin string) (string, func(), error) {
	file, err := os.CreateTemp(admin, "bench-recovery-index-")
	if err != nil {
		return "", nil, err
	}
	name := file.Name()
	if err := file.Close(); err != nil {
		return "", nil, err
	}
	if err := os.Remove(name); err != nil {
		return "", nil, err
	}
	return name, func() { _ = os.Remove(name) }, nil
}
func worktreeTree(path, admin string) (string, error) {
	index, cleanup, err := temporaryIndex(admin)
	if err != nil {
		return "", err
	}
	defer cleanup()
	if _, err := gitInput(path, []string{"GIT_INDEX_FILE=" + index}, nil, "read-tree", "HEAD"); err != nil {
		return "", err
	}
	if _, err := gitInput(path, []string{"GIT_INDEX_FILE=" + index}, nil, "add", "-A"); err != nil {
		return "", err
	}
	return gitInput(path, []string{"GIT_INDEX_FILE=" + index}, nil, "write-tree")
}
func realIndexTree(path, admin string) (string, error) {
	realIndex, err := git.Output("-C", path, "rev-parse", "--path-format=absolute", "--git-path", "index")
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(realIndex)
	if err != nil {
		return "", err
	}
	index, cleanup, err := temporaryIndex(admin)
	if err != nil {
		return "", err
	}
	defer cleanup()
	if err := os.WriteFile(index, data, 0o600); err != nil {
		return "", err
	}
	return gitInput(path, []string{"GIT_INDEX_FILE=" + index}, nil, "write-tree")
}
func readIndexEntries(path string) ([]indexEntry, bool, error) {
	raw, err := git.Raw("-C", path, "ls-files", "--stage", "-z")
	if err != nil {
		return nil, false, err
	}
	var entries []indexEntry
	conflicted := false
	for record := range bytes.SplitSeq(raw, []byte{0}) {
		if len(record) == 0 {
			continue
		}
		tab := bytes.IndexByte(record, '\t')
		if tab < 0 {
			return nil, false, errors.New("malformed staged index entry")
		}
		fields := strings.Fields(string(record[:tab]))
		if len(fields) != 3 {
			return nil, false, errors.New("malformed staged index metadata")
		}
		stage, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, false, err
		}
		conflicted = conflicted || stage != 0
		entries = append(entries, indexEntry{mode: fields[0], oid: fields[1], stage: stage, path: string(record[tab+1:])})
	}
	return entries, conflicted, nil
}
func conflictTree(path, admin string, entries []indexEntry, wanted int) (string, error) {
	index, cleanup, err := temporaryIndex(admin)
	if err != nil {
		return "", err
	}
	defer cleanup()
	if _, err := gitInput(path, []string{"GIT_INDEX_FILE=" + index}, nil, "read-tree", "--empty"); err != nil {
		return "", err
	}
	var input bytes.Buffer
	for _, entry := range entries {
		if entry.stage != 0 && entry.stage != wanted {
			continue
		}
		fmt.Fprintf(&input, "%s %s\t%s%c", entry.mode, entry.oid, entry.path, byte(0))
	}
	if _, err := gitInput(path, []string{"GIT_INDEX_FILE=" + index}, input.Bytes(), "update-index", "-z", "--index-info"); err != nil {
		return "", err
	}
	return gitInput(path, []string{"GIT_INDEX_FILE=" + index}, nil, "write-tree")
}
func commitTree(root, tree string, parents []string, message string) (string, error) {
	args := []string{"commit-tree", tree}
	for _, parent := range parents {
		args = append(args, "-p", parent)
	}
	env := []string{
		"GIT_AUTHOR_NAME=bench", "GIT_AUTHOR_EMAIL=bench@local",
		"GIT_COMMITTER_NAME=bench", "GIT_COMMITTER_EMAIL=bench@local",
	}
	return gitInput(root, env, []byte(message), args...)
}
func gitInput(root string, extraEnv []string, input []byte, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	cmd.Env = append(os.Environ(), extraEnv...)
	cmd.Stdin = bytes.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %s", args[0], strings.TrimSpace(stderr.String()))
	}
	return strings.TrimRight(stdout.String(), "\n"), nil
}
func nextRecoveryRef(root string, assignment intent.Assignment) (string, error) {
	prefix := intent.RecoveryRefPrefix(assignment.OwnerID, assignment.ID)
	out, err := git.Output("-C", root, "for-each-ref", "--format=%(refname)", prefix)
	if err != nil {
		return "", err
	}
	used := map[int]bool{}
	for _, ref := range strings.Split(out, "\n") {
		ordinal, _ := strconv.Atoi(strings.TrimPrefix(ref, prefix))
		used[ordinal] = true
	}
	for ordinal := 1; ; ordinal++ {
		if !used[ordinal] {
			return prefix + strconv.Itoa(ordinal), nil
		}
	}
}
func verifyRecovery(root string, assignment intent.Assignment, recovery intent.Recovery) error {
	resolved, err := git.Output("-C", root, "rev-parse", "--verify", recovery.Ref+"^{commit}")
	if err != nil || resolved != recovery.Root {
		return errors.New("recovery ref does not resolve to the recorded root")
	}
	for _, payload := range recovery.Payloads {
		if !git.OK("-C", root, "cat-file", "-e", payload+"^{commit}") || !git.OK("-C", root, "merge-base", "--is-ancestor", payload, recovery.Root) {
			return errors.New("recovery payload is missing or unreachable")
		}
	}
	for _, recorded := range assignment.Recovery {
		if recorded.Ref == recovery.Ref && recorded.Root == recovery.Root && strings.Join(recorded.Payloads, "\x00") == strings.Join(recovery.Payloads, "\x00") {
			return nil
		}
	}
	return errors.New("assignment does not name the verified recovery envelope")
}
func anchorDetached(root string, plan CleanupPlan) error {
	head, err := git.Output("-C", plan.Target, "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	zero := strings.Repeat("0", len(head))
	if out, err := exec.Command("git", "-C", root, "update-ref", plan.Recovery, head, zero).CombinedOutput(); err != nil {
		if existing, readErr := git.Output("-C", root, "rev-parse", "--verify", plan.Recovery+"^{commit}"); readErr != nil || existing != head {
			return fmt.Errorf("anchor detached HEAD: %s", strings.TrimSpace(string(out)))
		}
	}
	resolved, err := git.Output("-C", root, "rev-parse", "--verify", plan.Recovery+"^{commit}")
	if err != nil || resolved != head {
		return errors.New("detached recovery ref failed verification")
	}
	return nil
}
func discardIgnored(plan CleanupPlan) error {
	current, _, err := inventoryIgnored(plan.Target, false)
	if err != nil || current.Digest != plan.Ignored.Digest || current.Count != plan.Ignored.Count || current.Bytes != plan.Ignored.Bytes {
		return errStaleFingerprint
	}
	for _, name := range current.Paths {
		full := filepath.Join(plan.Target, filepath.Clean(filepath.FromSlash(name)))
		if _, err := ignoredLstat(full); err != nil {
			return errStaleFingerprint
		}
		warnBeforeRemovingLiveBinary(plan.Target, full)
		if err := os.Remove(full); err != nil {
			return fmt.Errorf("discard ignored path: %w", err)
		}
	}
	return nil
}
