package repairpilot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/jsonfile"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

type Options struct {
	Home      string
	Root      string
	RepoKey   string
	KitSource bool
	Now       time.Time
	Files     FileOps
}

type AtomicFile interface {
	io.Writer
	Close() error
	Name() string
	Sync() error
	Chmod(os.FileMode) error
}

type FileOps struct {
	MkdirAll   func(string, os.FileMode) error
	Chmod      func(string, os.FileMode) error
	OpenFile   func(string, int, os.FileMode) (AtomicFile, error)
	CreateTemp func(string, string) (AtomicFile, error)
	Rename     func(string, string) error
	Remove     func(string) error
}

type Document struct {
	Version       int           `json:"version"`
	RepositoryKey string        `json:"repository_key"`
	ActivatedAt   time.Time     `json:"activated_at"`
	CutoffAt      *time.Time    `json:"cutoff_at,omitempty"`
	Observations  []observation `json:"observations"`
	Audits        []audit       `json:"audits"`
}

// Store owns atomic replacement for one pilot document.
type Store struct {
	Path  string
	Files FileOps
}

// Replace preserves the prior document unless the final rename succeeds.
func (store Store) Replace(document Document) error {
	return writeDocument(store.Files.withDefaults(), store.Path, document)
}

var activateGrammar = usage.Grammar{
	Cmd:     "bench repair-pilot activate",
	Help:    "usage: bench repair-pilot activate",
	MaxArgs: 0,
}

var reportGrammar = usage.Grammar{
	Cmd:     "bench repair-pilot report",
	Help:    "usage: bench repair-pilot report [--full]",
	Flags:   []usage.Flag{{Name: "--full"}},
	MaxArgs: 0,
}

var recordGrammar = usage.Grammar{
	Cmd:     "bench repair-pilot record",
	Help:    "usage: bench repair-pilot record --input <file>",
	Flags:   []usage.Flag{{Name: "--input", HasValue: true, NoEmptyValue: true, Required: true}},
	MaxArgs: 0,
}

var FamilyUsage = familyUsage(activateGrammar, recordGrammar, reportGrammar)

func Command(options Options, args []string) (string, int) {
	if len(args) == 0 {
		return FamilyUsage, 2
	}
	if len(args) == 1 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		return FamilyUsage, 0
	}
	switch args[0] {
	case "activate":
		if _, line, code := usage.Parse(activateGrammar, args[1:]); line != "" {
			return line + "\n", code
		}
		return activate(options)
	case "report":
		parsed, line, code := usage.Parse(reportGrammar, args[1:])
		if line != "" {
			return line + "\n", code
		}
		_, full := parsed.Flags["--full"]
		return report(options, full)
	case "record":
		parsed, line, code := usage.Parse(recordGrammar, args[1:])
		if line != "" {
			return line + "\n", code
		}
		return record(options, parsed.Flags["--input"])
	default:
		return toon.Usage("bench repair-pilot", args[0]) + "\n", 2
	}
}

func activate(options Options) (string, int) {
	if !options.KitSource {
		return refusal("activation is limited to the Bench kit repository")
	}
	var document Document
	err := withPilotLock(options, true, func(files FileOps) error {
		loaded, state, loadErr := load(options)
		if loadErr != nil {
			return loadErr
		}
		document = loaded
		if state != bounds.StateAbsent {
			return nil
		}
		document = Document{Version: 1, RepositoryKey: options.RepoKey, ActivatedAt: options.Now.UTC(), Observations: []observation{}, Audits: []audit{}}
		return (Store{Path: documentPath(options), Files: files}).Replace(document)
	})
	if err != nil {
		return refusal(err.Error())
	}
	return renderDocumentStatus(document, options.Now)
}

func withPilotLock(options Options, create bool, update func(FileOps) error) error {
	files := options.Files.withDefaults()
	directory := filepath.Dir(documentPath(options))
	parentState, err := classifyDocumentParents(options)
	if err != nil {
		return err
	}
	if parentState == bounds.StateAbsent && !create {
		return errors.New("repair pilot is inactive; activate it before recording evidence")
	}
	if create {
		if err := files.MkdirAll(directory, 0o700); err != nil {
			return err
		}
		if err := files.Chmod(directory, 0o700); err != nil {
			return err
		}
	}
	lockPath := filepath.Join(directory, "pilot.lock")
	lock, err := files.OpenFile(lockPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return errors.New("pilot update is locked; retry after the active writer completes")
	}
	if err := lock.Close(); err != nil {
		files.Remove(lockPath)
		return err
	}
	defer files.Remove(lockPath)
	return update(files)
}

func report(options Options, full bool) (string, int) {
	document, state, err := load(options)
	if err != nil {
		return refusal(err.Error())
	}
	if state == bounds.StateAbsent {
		return renderPilotReport(Document{}, options.Now, full)
	}
	return renderPilotReport(document, options.Now, full)
}

func documentPath(options Options) string {
	return filepath.Join(options.Home, "repair-pilot", options.RepoKey, "pilot.json")
}

func load(options Options) (Document, bounds.FileState, error) {
	parentState, err := classifyDocumentParents(options)
	if err != nil {
		return Document{}, parentState, err
	}
	if parentState == bounds.StateAbsent {
		return Document{}, parentState, nil
	}
	classified := bounds.ClassifyNoFollow(documentPath(options))
	if classified.State == bounds.StateAbsent {
		return Document{}, classified.State, nil
	}
	if classified.State != bounds.StateParsed {
		return Document{}, classified.State, fmt.Errorf("pilot document is %s", classified.State)
	}
	var document Document
	if err := jsonfile.Decode(classified.Data, &document); err != nil {
		return Document{}, bounds.StateMalformed, fmt.Errorf("pilot document is malformed: %w", err)
	}
	if document.Version != 1 {
		return Document{}, bounds.StateUnsupportedSchema, fmt.Errorf("pilot document has unsupported version %d", document.Version)
	}
	if document.RepositoryKey != options.RepoKey {
		return Document{}, bounds.StateUnsupportedSchema, errors.New("pilot document names a different repository")
	}
	if document.ActivatedAt.IsZero() {
		return Document{}, bounds.StateUnsupportedSchema, errors.New("pilot document has no activation time")
	}
	return document, bounds.StateParsed, nil
}

func classifyDocumentParents(options Options) (bounds.FileState, error) {
	for _, directory := range []string{
		filepath.Join(options.Home, "repair-pilot"),
		filepath.Dir(documentPath(options)),
	} {
		classified := bounds.ClassifyDirNoFollow(directory)
		switch classified.State {
		case bounds.StateAbsent:
			return bounds.StateAbsent, nil
		case bounds.StateEmpty, bounds.StateParsed:
			continue
		default:
			return classified.State, fmt.Errorf("pilot directory is %s", classified.State)
		}
	}
	return bounds.StateParsed, nil
}

func (ops FileOps) withDefaults() FileOps {
	if ops.MkdirAll == nil {
		ops.MkdirAll = os.MkdirAll
	}
	if ops.Chmod == nil {
		ops.Chmod = os.Chmod
	}
	if ops.OpenFile == nil {
		ops.OpenFile = func(name string, flag int, perm os.FileMode) (AtomicFile, error) {
			return os.OpenFile(name, flag, perm)
		}
	}
	if ops.CreateTemp == nil {
		ops.CreateTemp = func(directory, pattern string) (AtomicFile, error) { return os.CreateTemp(directory, pattern) }
	}
	if ops.Rename == nil {
		ops.Rename = os.Rename
	}
	if ops.Remove == nil {
		ops.Remove = os.Remove
	}
	return ops
}

func writeDocument(files FileOps, path string, document Document) error {
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	temporary, err := files.CreateTemp(filepath.Dir(path), ".pilot-*")
	if err != nil {
		return fmt.Errorf("create temporary pilot document: %w", err)
	}
	temporaryPath := temporary.Name()
	keep := false
	defer func() {
		if !keep {
			files.Remove(temporaryPath)
		}
	}()
	if err := temporary.Chmod(0o600); err != nil {
		temporary.Close()
		return err
	}
	if n, err := temporary.Write(data); err != nil {
		temporary.Close()
		return fmt.Errorf("write temporary pilot document: %w", err)
	} else if n != len(data) {
		temporary.Close()
		return io.ErrShortWrite
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := files.Rename(temporaryPath, path); err != nil {
		return fmt.Errorf("replace pilot document: %w", err)
	}
	keep = true
	return nil
}

func renderDocumentStatus(document Document, now time.Time) (string, int) {
	deadline := collectionDeadline(document)
	state := "active"
	if !now.Before(deadline) || document.CutoffAt != nil {
		state = "stopped"
	}
	return renderStatus(state, document.ActivatedAt.UTC().Format(time.RFC3339), deadline.Format(time.RFC3339), summarizeDocument(document))
}

func familyUsage(grammars ...usage.Grammar) string {
	lines := make([]string, 0, len(grammars))
	for i, grammar := range grammars {
		line := grammar.Help
		if i != 0 {
			line = "       " + strings.TrimPrefix(line, "usage: ")
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n") + "\n"
}

func refusal(detail string) (string, int) {
	return toon.Errorf("repair pilot refused", detail) + "\n", 1
}

type pilotSummary struct {
	ComparableRepeats int
	Progress          map[string]int
}

func summarizeDocument(document Document) pilotSummary {
	summary := pilotSummary{Progress: map[string]int{}}
	prior := map[sequence][]failure{}
	for _, item := range document.Observations {
		for _, current := range item.Failures {
			for _, earlier := range prior[item.Sequence] {
				if comparableFailure(earlier, current) {
					summary.ComparableRepeats++
					break
				}
			}
		}
		prior[item.Sequence] = append(prior[item.Sequence], item.Failures...)
		summary.Progress[effectiveProgressLabel(document, item.ID)]++
	}
	return summary
}

func renderStatus(state, activated, deadline string, summary pilotSummary) (string, int) {
	columns := []string{"state", "activated_at", "deadline", "comparable_repeats"}
	row := []string{state, activated, deadline, strconv.Itoa(summary.ComparableRepeats)}
	for _, label := range append(append([]string{}, progressLabels...), "unknown") {
		columns = append(columns, "progress_"+label)
		row = append(row, strconv.Itoa(summary.Progress[label]))
	}
	out, err := toon.Table("repair_pilot", columns, [][]string{row})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out, 0
}
