package worktree

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/gibbonmi/bench/internal/git"
)

// checkoutStatus owns the tracked-state argv that the cleanup planner, the nested
// classifier, and the reset planner all read, so the four flags never drift apart.
func checkoutStatus(target string) ([]byte, error) {
	return git.Raw("--no-optional-locks", "-C", target, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--ignore-submodules=none")
}

// ignoredListing owns the ignored-inventory argv that the cleanup inventory and the
// reset's collision check both read.
func ignoredListing(target string) ([]byte, error) {
	return git.Raw("--no-optional-locks", "-C", target, "ls-files", "--others", "--ignored", "--exclude-standard", "-z", "--")
}

// explicitContentIdentity digests the tracked diff and every untracked path's bytes, so
// an edit inside an already-dirty file changes the identity.
func explicitContentIdentity(target string) (string, error) {
	diff, err := git.Raw("--no-optional-locks", "-C", target, "diff", "--no-ext-diff", "--binary", "HEAD", "--")
	if err != nil {
		return "", fmt.Errorf("read worktree content identity: %w", err)
	}
	untracked, err := git.Raw("--no-optional-locks", "-C", target, "ls-files", "--others", "--exclude-standard", "-z")
	if err != nil {
		return "", fmt.Errorf("read untracked content identity: %w", err)
	}
	parts := [][]byte{diff}
	for record := range bytes.SplitSeq(untracked, []byte{0}) {
		if len(record) == 0 {
			continue
		}
		name := string(record)
		full := filepath.Join(target, filepath.FromSlash(name))
		info, statErr := os.Lstat(full)
		if statErr != nil {
			return "", fmt.Errorf("stat untracked path: %w", statErr)
		}
		parts = append(parts, []byte(name), []byte(strconv.FormatUint(uint64(info.Mode()), 10)))
		switch {
		case info.Mode().IsRegular():
			body, readErr := os.ReadFile(full)
			if readErr != nil {
				return "", fmt.Errorf("read untracked path: %w", readErr)
			}
			parts = append(parts, body)
		case info.Mode()&os.ModeSymlink != 0:
			link, readErr := os.Readlink(full)
			if readErr != nil {
				return "", fmt.Errorf("read untracked symlink: %w", readErr)
			}
			parts = append(parts, []byte(link))
		}
	}
	return fingerprintParts(parts...), nil
}

// captureLayers writes one recovery envelope for the checkout at path and returns the
// root commit and its payload commits. The manifest carries the base head.
func captureLayers(root, path string, keepWorking bool, tip string) (string, []string, error) {
	head, err := git.Output("-C", path, "rev-parse", "HEAD")
	if err != nil {
		return "", nil, fmt.Errorf("read recovery HEAD: %w", err)
	}
	headTree, err := git.Output("-C", path, "rev-parse", "HEAD^{tree}")
	if err != nil {
		return "", nil, fmt.Errorf("read recovery base tree: %w", err)
	}
	admin, err := git.AdminDir(path)
	if err != nil {
		return "", nil, err
	}
	layerTrees := map[string]string{}
	workingTree, err := worktreeTree(path, admin)
	if err != nil {
		return "", nil, fmt.Errorf("capture working layer: %w", err)
	}
	if keepWorking || workingTree != headTree {
		layerTrees["working"] = workingTree
	}
	entries, conflicted, err := readIndexEntries(path)
	if err != nil {
		return "", nil, err
	}
	if conflicted {
		for stage, name := range map[int]string{1: "base", 2: "ours", 3: "theirs"} {
			tree, err := conflictTree(path, admin, entries, stage)
			if err != nil {
				return "", nil, fmt.Errorf("capture conflict %s layer: %w", name, err)
			}
			layerTrees[name] = tree
		}
	} else {
		stagedTree, err := realIndexTree(path, admin)
		if err != nil {
			return "", nil, fmt.Errorf("capture staged layer: %w", err)
		}
		if stagedTree != headTree {
			layerTrees["staged"] = stagedTree
		}
	}
	if len(layerTrees) == 0 {
		return "", nil, errors.New("recovery requested for a clean assignment")
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
				return "", nil, err
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
	if tip != "" {
		manifest.Tip = tip
	}
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		return "", nil, err
	}
	manifestBytes = append(manifestBytes, '\n')
	blob, err := gitInput(root, nil, manifestBytes, "hash-object", "-w", "--stdin")
	if err != nil {
		return "", nil, err
	}
	rootTree, err := gitInput(root, nil, []byte("100644 blob "+blob+"\tmanifest.json\n"), "mktree")
	if err != nil {
		return "", nil, err
	}
	parents := append([]string(nil), payloads...)
	if manifest.Tip != "" && manifest.Tip != head {
		parents = append(parents, manifest.Tip)
	}
	rootOID, err := commitTree(root, rootTree, parents, "bench recovery root\n")
	if err != nil {
		return "", nil, err
	}
	return rootOID, payloads, nil
}
