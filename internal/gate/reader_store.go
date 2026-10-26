package gate

// The reader-table store. A whole run of the kit's built-in table owns one test log
// directory. After its green verdict persists, the run derives the reader table from those
// logs and stores it in the common Git directory, keyed by the graded tree.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate/greenmarker"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/runbinary"
)

const (
	readerTableSchema = 1
	// readerTableLimit bounds the tables that the store keeps.
	readerTableLimit = 32
	// readerTableSuffix ends the file name of each stored table, after its tree.
	readerTableSuffix = ".json"
	// testPackageTemplate makes go list print the directory of each package that holds a
	// test file, and an empty line for each other package.
	testPackageTemplate = "{{if or .TestGoFiles .XTestGoFiles}}{{.Dir}}{{end}}"
)

// readerTableRecord is one stored reader table: the graded tree, the identity digest of the
// run that derived it, and the readers of each covered family.
type readerTableRecord struct {
	Schema   int         `json:"schema"`
	Tree     string      `json:"tree"`
	Identity string      `json:"identity"`
	Readers  readerTable `json:"readers"`
}

// readerTrace is the test log directory of one whole run of the kit's built-in table, with
// the graded root and the kit of its traced test phase.
type readerTrace struct {
	logs, root, kit string
}

// beginReaderTrace opens the test log directory of a whole run whose schedule is the kit's
// built-in table with its traced test phase. It answers nil for every other run, and when
// the directory cannot be made. A nil trace records no log and stores no table.
func beginReaderTrace(runtimeRoot, storageRoot string, plan subject, selection *runbinary.Selection) *readerTrace {
	if selection == nil || !narrowsKitTable(runtimeRoot, storageRoot, plan.Resolution) || !RootConformanceRuns(runtimeRoot, selection.SourceRoot) {
		return nil
	}
	logs, err := os.MkdirTemp("", "bench-gate-testlog-")
	if err != nil {
		return nil
	}
	return &readerTrace{logs: logs, root: runtimeRoot, kit: selection.SourceRoot}
}

// env names the log directory in env, for the gate-phases process that runs the table.
func (t *readerTrace) env(env []string) []string {
	if t == nil {
		return env
	}
	return mergeEnv(env, []string{testLogEnv + "=" + t.logs})
}

// close removes the log directory.
func (t *readerTrace) close() {
	if t != nil {
		_ = os.RemoveAll(t.logs)
	}
}

// persist derives the reader table of a run whose green verdict has persisted, and stores
// it under storageRoot through replacer. A failure stores no table and leaves the verdict as
// it is; the run log names the failure.
func (t *readerTrace) persist(ctx context.Context, replacer durablefile.Replacer, storageRoot string, plan subject) {
	if t == nil {
		return
	}
	path, err := t.store(replacer, storageRoot, plan)
	record := gateLogRecord{Event: "reader.table", Root: storageRoot, Path: path}
	if err != nil {
		record.Detail = err.Error()
	}
	logGateEvent(ctx, record)
}

func (t *readerTrace) store(replacer durablefile.Replacer, storageRoot string, plan subject) (string, error) {
	packages, err := testPackages(t.root, t.kit, plan.Env)
	if err != nil {
		return "", err
	}
	links, err := benchgit.TreeSymlinks(storageRoot, plan.Tree)
	if err != nil {
		return "", err
	}
	logs := map[string][]byte{}
	for _, pkg := range packages {
		file, err := testLogFile(t.logs, t.root, filepath.Join(t.root, filepath.FromSlash(pkg)))
		if err != nil {
			return "", err
		}
		data, err := os.ReadFile(file)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return "", err
		}
		logs[pkg] = data
	}
	table, err := deriveReaderTable(readerSources{Root: t.root, Packages: packages, Logs: logs, Links: links, Families: coveredFamilies, Children: childReaders})
	if err != nil {
		return "", err
	}
	return storeReaderTable(replacer, storageRoot, readerTableRecord{Schema: readerTableSchema, Tree: plan.Tree, Identity: plan.Identity, Readers: table})
}

// testPackages lists the test packages of the traced test phase of root: each package of
// the phase's pattern that holds a test file, as a slash-separated directory relative to
// root. The listing passes the phase's build flags and runs under env, the closed
// environment of the run, with the phase's own entries, so an ambient Go variable of the
// operator changes no listed package.
func testPackages(root, kit string, env []string) ([]string, error) {
	phase := wholeTestPhase(root, kit)
	argv := append(append([]string{"list", DisableBuildVCS}, testBuildFlags...), "-f", testPackageTemplate, wholeTestPackages)
	cmd := exec.Command("go", argv...)
	cmd.Dir, cmd.Env = phase.Dir, mergeEnv(env, phase.Env)
	listed, err := cmd.Output()
	if err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			return nil, fmt.Errorf("go list failed: %w: %s", err, exit.Stderr)
		}
		return nil, fmt.Errorf("go list failed: %w", err)
	}
	var packages []string
	for _, dir := range listedFiles(string(listed)) {
		rel, inside := rootRelative(root, dir)
		if !inside {
			return nil, fmt.Errorf("test package %s lies outside %s", dir, root)
		}
		packages = append(packages, rel)
	}
	return packages, nil
}

// readerTablePath is the store path of the table of tree under the common Git directory.
func readerTablePath(gitdir, tree string) string {
	return filepath.Join(gitdir, benchgit.ReaderTableDir, tree+readerTableSuffix)
}

// storeReaderTable installs record as the table of its tree in the store of root through
// replacer, and then prunes the store. It answers the path of the table.
func storeReaderTable(replacer durablefile.Replacer, root string, record readerTableRecord) (string, error) {
	gitdir, err := benchgit.CommonDir(root)
	if err != nil {
		return "", err
	}
	path := readerTablePath(gitdir, record.Tree)
	if err := ensureEvidenceDir(gitdir, filepath.Dir(path)); err != nil {
		return "", err
	}
	data, err := json.Marshal(record)
	if err != nil {
		return "", err
	}
	if err := durableReplaceRecordAt(replacer, filepath.Dir(path), filepath.Base(path), data); err != nil {
		return "", err
	}
	return path, pruneReaderTables(root, filepath.Dir(path))
}

// loadReaderTable answers the stored table of tree in the store of root.
func loadReaderTable(root, tree string) (readerTableRecord, error) {
	gitdir, err := benchgit.CommonDir(root)
	if err != nil {
		return readerTableRecord{}, err
	}
	data, err := os.ReadFile(readerTablePath(gitdir, tree))
	if err != nil {
		return readerTableRecord{}, err
	}
	var record readerTableRecord
	if err := strictJSON(data, &record); err != nil {
		return readerTableRecord{}, err
	}
	if record.Schema != readerTableSchema || record.Tree != tree || !isContentAddress(record.Identity) {
		return readerTableRecord{}, errors.New("invalid reader table")
	}
	return record, nil
}

// pruneReaderTables removes the oldest tables of dir until it holds readerTableLimit. It
// never removes the table of the green marker's tree.
func pruneReaderTables(root, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	type stored struct {
		name     string
		modified time.Time
	}
	var tables []stored
	for _, entry := range entries {
		if !entry.Type().IsRegular() || strings.HasPrefix(entry.Name(), ".") || !strings.HasSuffix(entry.Name(), readerTableSuffix) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		tables = append(tables, stored{entry.Name(), info.ModTime()})
	}
	slices.SortFunc(tables, func(a, b stored) int {
		if order := a.modified.Compare(b.modified); order != 0 {
			return order
		}
		return strings.Compare(a.name, b.name)
	})
	kept, _ := greenMarkerTree(root)
	excess := len(tables) - readerTableLimit
	for _, table := range tables {
		if excess <= 0 {
			break
		}
		if table.name == kept+readerTableSuffix {
			continue
		}
		if err := os.Remove(filepath.Join(dir, table.name)); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		excess--
	}
	return nil
}

// greenMarkerTree answers the tree of the green marker of root's default branch, and false
// when no marker resolves.
func greenMarkerTree(root string) (string, bool) {
	branch, ok := benchgit.ResolvedDefault(root)
	if !ok {
		return "", false
	}
	commit, present, err := greenmarker.Read(root, branch)
	if err != nil || !present {
		return "", false
	}
	tree, err := benchgit.Output("-C", root, "rev-parse", commit+"^{tree}")
	return tree, err == nil
}
