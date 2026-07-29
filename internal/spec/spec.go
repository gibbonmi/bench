// Package spec owns spec-file addressing and the two spec-lifecycle operations: the
// `bench spec implemented <slug>` status flip (Flip) and the `bench spec retire <slug>`
// deletion. Resolve is the one source of the spec-argument convention (path-first, then a
// specs/<slug>/spec.md and specs/<slug>.md fallbacks) that `bench coverage`, `bench commit --spec`, and both
// operations take their argument through. Flip is the single source of the status-line
// flip: it turns exactly one line-start `Status: staged` into the retirement-detector form
// `Status: implemented`, preserving every other byte, and is composed by `bench commit`.
// AwaitsRetirement is the single source of the merged-implemented predicate — the
// `implemented` twin of the staged form — shared by retire's validation and the
// `bench status` specs-awaiting-retirement counter.
package spec

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

// stagedRe matches the flip's input: a `Status:` line whose sole value is `staged`,
// tab/space separated, only whitespace trailing — the `staged` twin of the retirement
// detector's `^Status:[ \t]+implemented[ \t]*$`. Swapping `staged` for `implemented` in
// place therefore yields exactly the detector's accepted form by construction.
var stagedRe = regexp.MustCompile(`^Status:[ \t]+staged[ \t]*$`)

// Fact is one typed live spec record. Path is its repository-relative path.
type Fact struct {
	Slug, Path, Status, RoadmapID string
}

func metadata(content []byte) (status, roadmapID string) {
	inFence := false
	for _, raw := range strings.Split(string(content), "\n") {
		line := strings.TrimSuffix(raw, "\r")
		if strings.HasPrefix(line, "```") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		switch {
		case strings.HasPrefix(line, "Status:"):
			status = strings.TrimSpace(strings.TrimPrefix(line, "Status:"))
		case strings.HasPrefix(line, "Roadmap:"):
			roadmapID = strings.TrimSpace(strings.TrimPrefix(line, "Roadmap:"))
		}
	}
	return status, roadmapID
}

// Facts reads folder and flat specs in path order. Malformed files are retained with an
// empty status so callers can report the evidence rather than silently omit it.
func Facts(root string) ([]Fact, error) {
	flat, err := filepath.Glob(filepath.Join(root, "specs", "*.md"))
	if err != nil {
		return nil, err
	}
	folder, err := filepath.Glob(filepath.Join(root, "specs", "*", "spec.md"))
	if err != nil {
		return nil, err
	}
	paths := append(flat, folder...)
	sort.Strings(paths)
	out := make([]Fact, 0, len(paths))
	for _, path := range paths {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil, err
		}
		f := Fact{
			Slug: strings.TrimSuffix(filepath.Base(path), ".md"),
			Path: filepath.ToSlash(rel),
		}
		if filepath.Base(path) == "spec.md" {
			f.Slug = filepath.Base(filepath.Dir(path))
		}
		f.Status, f.RoadmapID = metadata(b)
		out = append(out, f)
	}
	return out, nil
}

// AwaitsRetirement reports whether spec content carries an unfenced retirement marker: the
// one definition of "a merged spec awaiting retirement", shared by the `bench status`
// specs row (which counts it across specs/) and `bench spec retire` (which validates it
// before deleting). Each line is CRLF-stripped; a line whose first three bytes are a code
// fence toggles fence state and is skipped; lines inside a fence are skipped; the first
// line matching the retirement regex marks the content.
func AwaitsRetirement(content []byte) bool {
	status, _ := metadata(content)
	return status == StatusImplemented
}

// StatusImplemented is the Status value a finished spec carries. It is exported as the one
// source of the token, so a caller classifying a Fact by status compares against the same
// spelling the retirement predicate does.
const StatusImplemented = "implemented"

// Resolve finds the readable file backing a spec argument: the argument as given
// (path-first, so a same-named readable file shadows the fallback), then — for a
// separator-free argument only — a flat fallback followed by a folder fallback. base
// anchors those fallbacks: pass the repo root to
// resolve it repo-root-relative (so a cwd deeper than the root still finds it), or "" to
// resolve it relative to the process cwd. ok is false when no form resolves; tried holds
// every form attempted, for the not-found error. A non-nil err is a read failure on an
// existing file (e.g. permissions), reported instead of not-found.
func Resolve(base, arg string) (content []byte, resolved string, tried []string, ok bool, err error) {
	tried = []string{arg}
	if b, err := readCandidate(arg); err != nil || b != nil {
		return b, arg, tried, err == nil, err
	}
	if !strings.ContainsRune(arg, '/') {
		slug := strings.TrimSuffix(arg, ".md")
		flat := filepath.Join(base, "specs", slug+".md")
		tried = append(tried, flat)
		if b, err := readCandidate(flat); err != nil || b != nil {
			return b, flat, tried, err == nil, err
		}
		folder := filepath.Join(base, "specs", slug, "spec.md")
		tried = append(tried, folder)
		if b, err := readCandidate(folder); err != nil || b != nil {
			return b, folder, tried, err == nil, err
		}
	}
	return nil, "", tried, false, nil
}

// readCandidate reads path as a candidate spec, through the classifier so a FIFO or
// device parked where a spec belongs is rejected before the open rather than blocking
// forever. An absent path or a directory is not a candidate (nil, nil — try the next
// form); any other failure is a real error to surface, never masked as not-found.
func readCandidate(path string) ([]byte, error) {
	c := bounds.Classify(path, bounds.ControlRecordLimit)
	switch {
	case c.State == bounds.StateAbsent:
		return nil, nil
	case c.State == bounds.StateWrongType && isDir(path):
		return nil, nil
	case c.State.Failed():
		return nil, errors.New(c.Reason)
	}
	return c.Data, nil
}

// isDir separates the one non-regular path that means "keep resolving" from the ones that
// mean "this read failed": a directory named like the candidate simply is not the spec
// file, while a special file where a spec belongs is a problem the caller must hear about.
func isDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// locateStaged resolves arg (base anchors both spec fallbacks) and requires exactly
// one line-start `Status: staged`, returning the resolved path, the file split into lines, and
// the index of that one line. It never writes — it is the shared core of CheckStaged
// (validate only) and Flip (validate then rewrite), so the resolution + single-staged-line
// rule has one source. The error names the file and the reason on not-found, not-readable, no
// `Status: staged` line (missing or already implemented), or more than one.
func locateStaged(base, arg string) (resolved string, lines [][]byte, idx int, err error) {
	content, resolved, tried, ok, readErr := Resolve(base, arg)
	if readErr != nil {
		return "", nil, -1, fmt.Errorf("spec not readable: %s: %v", resolved, readErr)
	}
	if !ok {
		return "", nil, -1, fmt.Errorf("spec not found: %s", strings.Join(tried, ", "))
	}
	lines = bytes.Split(content, []byte("\n"))
	matches := 0
	idx = -1
	for i, line := range lines {
		if stagedRe.Match(line) {
			matches++
			idx = i
		}
	}
	if matches == 0 {
		return "", nil, -1, fmt.Errorf("no `Status: staged` line in %s (already implemented, or missing)", resolved)
	}
	if matches > 1 {
		return "", nil, -1, fmt.Errorf("%d `Status: staged` lines in %s (expected exactly one)", matches, resolved)
	}
	return resolved, lines, idx, nil
}

// CheckStaged resolves arg and confirms it carries exactly one line-start `Status: staged`,
// returning the resolved path. It never writes: it is the fail-fast validation `bench commit
// --spec` runs before the gate, so a bad or already-implemented spec is rejected before the
// gate burns rather than after Flip runs on a green tree. Flip re-checks the (possibly changed)
// file before it rewrites, so the two share locateStaged's one validation rule.
func CheckStaged(base, arg string) (resolved string, err error) {
	resolved, _, _, err = locateStaged(base, arg)
	return resolved, err
}

// Flip resolves arg (base anchors the specs/<slug>.md fallback), requires exactly one
// line-start `Status: staged`, rewrites that one line to `Status: implemented` in place —
// every other byte, including a missing final newline, preserved — writes the file, and
// returns the resolved path. It edits the file only; it never stages. The error names the
// file and the reason on not-found, not-readable, no `Status: staged` line (missing or
// already implemented), or more than one — so a typo or a re-run is non-destructive.
func Flip(base, arg string) (resolved string, err error) {
	resolved, lines, idx, err := locateStaged(base, arg)
	if err != nil {
		return "", err
	}
	lines[idx] = bytes.Replace(lines[idx], []byte("staged"), []byte("implemented"), 1)
	out := bytes.Join(lines, []byte("\n"))
	mode := os.FileMode(0o644)
	if fi, statErr := os.Stat(resolved); statErr == nil {
		mode = fi.Mode().Perm()
	}
	if err := os.WriteFile(resolved, out, mode); err != nil {
		return "", fmt.Errorf("write %s: %v", resolved, err)
	}
	return resolved, nil
}

// Command implements `bench spec <subcommand> <slug>`: `implemented` flips the status line
// (Flip), `retire` deletes a merged spec and its review pickup, `history` reports the
// commits that retired or deleted it (a read-only AXI query — see history.go). The
// slug's specs/<slug>.md fallback is anchored at the repo root, so it resolves from any
// cwd inside the repo; a path argument stays cwd-relative. Usage errors (missing/unknown
// subcommand, missing or extra argument, unknown flag) exit 2; a resolve/validate/delete
// failure exits 1 naming the file and reason.
func Command(args []string) (string, int) {
	if len(args) == 0 {
		return toon.Usage("bench spec", "expected a subcommand: implemented, retire, history") + "\n", 2
	}
	switch args[0] {
	case "implemented":
		return implementedCommand(args[1:])
	case "retire":
		return retireCommand(args[1:])
	case "history":
		return historyCommand(args[1:])
	default:
		return toon.Usage("bench spec", args[0]) + "\n", 2
	}
}

// specArg extracts the single spec argument shared by every `bench spec` subcommand:
// exactly one positional (a path or slug). `-h/--help` prints help at exit 0; an unknown
// flag, a second positional, or a missing argument is a usage error at exit 2. ok is false
// on any terminal outcome, with out/code carrying it; otherwise arg holds the one
// positional. cmd names the subcommand for the usage line, help is its `-h` text.
func specArg(cmd, help string, rest []string) (arg, out string, code int, ok bool) {
	for _, a := range rest {
		switch {
		case a == "-h" || a == "--help":
			return "", help, 0, false
		case strings.HasPrefix(a, "-"):
			return "", toon.Usage(cmd, a) + "\n", 2, false
		default:
			if arg != "" {
				return "", toon.Usage(cmd, a) + "\n", 2, false
			}
			arg = a
		}
	}
	if arg == "" {
		return "", toon.MissingArg(cmd, "<spec.md | slug> is required") + "\n", 2, false
	}
	return arg, "", 0, true
}

// RepoBase returns the repo root that anchors the specs/<slug>.md fallback, or "" when the
// cwd is outside a repo (then the fallback resolves relative to the process cwd). It is the
// one source of that anchoring for every subcommand that resolves a bare slug, so where the
// caller happens to be standing never changes which spec a slug names.
func RepoBase() string {
	if root, err := git.Root(); err == nil {
		return root
	}
	return ""
}

// implementedCommand runs `bench spec implemented <slug>`: it flips the one `Status: staged`
// line to `Status: implemented`. A resolve/validate/write failure exits 1 naming the file.
func implementedCommand(rest []string) (string, int) {
	arg, out, code, ok := specArg("bench spec implemented", "usage: bench spec implemented <spec.md | slug>\n", rest)
	if !ok {
		return out, code
	}
	resolved, err := Flip(RepoBase(), arg)
	if err != nil {
		return toon.Errorf(err.Error(), "pass a spec with a single `Status: staged` line") + "\n", 1
	}
	return fmt.Sprintf("spec implemented: %s\n", resolved), 0
}

// retireCommand runs `bench spec retire <slug>`: on a merged-implemented spec it deletes
// the review pickup (when present) then the complete folder spec, and prints what it removed plus the
// judgment duties that remain. It never commits and never runs the gate — `bench commit`
// owns commit discipline. Every unsafe input refuses at exit 1 without deleting anything: a
// spec that is not merged-implemented (staged, or implemented only in the working tree and
// not yet at HEAD), an unknown slug, or an orphaned review pickup with no spec.
func retireCommand(rest []string) (string, int) {
	arg, out, code, ok := specArg("bench spec retire", "usage: bench spec retire <spec.md | slug>\n", rest)
	if !ok {
		return out, code
	}
	base := RepoBase()
	content, resolved, tried, found, err := Resolve(base, arg)
	if err != nil {
		return toon.Errorf(fmt.Sprintf("spec not readable: %s: %v", resolved, err), "check file permissions") + "\n", 1
	}
	if !found {
		if residue, ok := folderResidue(base, arg); ok {
			return toon.Errorf("incomplete retired spec folder: "+residue,
				"remove "+residue+" by hand after reviewing its residue; retire will not auto-clean it") + "\n", 1
		}
		// A review pickup with no spec is a reviewer judgment, never an auto-clean.
		if orphan, ok := orphanPickup(base, arg); ok {
			return toon.Errorf("orphaned review pickup: "+orphan+" has no spec",
				"delete "+orphan+" by hand if its residual risk is accepted; retire never auto-cleans a pickup") + "\n", 1
		}
		return toon.Errorf("spec not found: "+strings.Join(tried, ", "), "pass a merged spec path or slug") + "\n", 1
	}
	// Merged-implemented reuses AwaitsRetirement (the status detector), and additionally
	// requires the marker at HEAD: an implemented-but-uncommitted spec means the finishing
	// commit has not landed, so retiring it would delete work git cannot recover.
	if !AwaitsRetirement(content) {
		return toon.Errorf("spec not merged-implemented: "+RelTo(base, resolved),
			"retire only a spec whose Status line reads implemented") + "\n", 1
	}
	if !implementedAtHEAD(base, resolved) {
		return toon.Errorf("spec implemented in the working tree but not at HEAD: "+RelTo(base, resolved),
			"commit the finishing flip before retiring") + "\n", 1
	}
	// Deletion order: review pickup first, then the spec, so an interrupt between the two
	// leaves a valid spec that a re-run retires cleanly — never an orphaned review file.
	var b strings.Builder
	slug := slugOf(resolved)
	if pickup := filepath.Join(base, "reviews", slug+".md"); fileExists(pickup) {
		if err := os.Remove(pickup); err != nil {
			return toon.Errorf(fmt.Sprintf("remove %s: %v", RelTo(base, pickup), err), "check file permissions") + "\n", 1
		}
		fmt.Fprintf(&b, "retired: %s\n", RelTo(base, pickup))
	}
	if filepath.Base(resolved) == "spec.md" {
		folder := filepath.Dir(resolved)
		if err := os.RemoveAll(folder); err != nil {
			return b.String() + toon.Errorf(fmt.Sprintf("remove %s: %v", RelTo(base, folder), err), "check file permissions") + "\n", 1
		}
		fmt.Fprintf(&b, "retired: %s\n", RelTo(base, folder))
	} else if err := os.Remove(resolved); err != nil {
		return b.String() + toon.Errorf(fmt.Sprintf("remove %s: %v", RelTo(base, resolved), err), "check file permissions") + "\n", 1
	} else {
		fmt.Fprintf(&b, "retired: %s\n", RelTo(base, resolved))
	}
	fmt.Fprintf(&b, "next: promote durable content, remove the ROADMAP row, commit as `spec-retire: %s`\n", slug)
	return b.String(), 0
}

// folderResidue identifies the terminal interrupted-retire state: a folder remains but
// its spec file is gone. It is not safe to infer which residual files a reviewer accepts.
func folderResidue(base, arg string) (string, bool) {
	if strings.ContainsRune(arg, '/') {
		return "", false
	}
	folder := filepath.Join(base, "specs", strings.TrimSuffix(arg, ".md"))
	fi, err := os.Stat(folder)
	if err != nil || !fi.IsDir() {
		return "", false
	}
	if _, err := os.Stat(filepath.Join(folder, "spec.md")); !os.IsNotExist(err) {
		return "", false
	}
	return RelTo(base, folder), true
}

// implementedAtHEAD reports whether the spec's content at HEAD carries the retirement
// marker — the "finishing commit has landed" guard. It reads the blob through git (stderr
// is discarded), so a spec absent from HEAD or still staged there reads as false.
func implementedAtHEAD(base, resolved string) bool {
	rel := filepath.ToSlash(RelTo(base, resolved))
	args := []string{"show", "HEAD:" + rel}
	if base != "" {
		args = append([]string{"-C", base}, args...)
	}
	content, err := git.Raw(args...)
	if err != nil {
		return false
	}
	return AwaitsRetirement(content)
}

// orphanPickup reports the repo-relative reviews/<slug>.md when a review pickup exists for a
// slug whose spec did not resolve — the orphaned-pickup refusal names it.
func orphanPickup(base, arg string) (string, bool) {
	pickup := filepath.Join(base, "reviews", slugOf(arg)+".md")
	if fileExists(pickup) {
		return RelTo(base, pickup), true
	}
	return "", false
}

// slugOf is the spec slug of a path or bare argument. A folder spec's filename is always
// spec.md, so its parent folder supplies the slug rather than the generic filename.
func slugOf(arg string) string {
	if filepath.Base(arg) == "spec.md" {
		return filepath.Base(filepath.Dir(arg))
	}
	return strings.TrimSuffix(filepath.Base(arg), ".md")
}

// fileExists reports whether path is an existing regular file (not a directory).
func fileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

// RelTo renders path repo-relative to base for a stable `retired: specs/<slug>.md` line,
// falling back to the path verbatim when base is empty or path lies outside base. It is
// the one source of that rendering, so a resolved spec reads the same from any cwd and
// never leaks the machine's directory layout into agent-facing output.
func RelTo(base, path string) string {
	if base == "" {
		return path
	}
	abs := path
	if !filepath.IsAbs(abs) {
		if a, err := filepath.Abs(path); err == nil {
			abs = a
		}
	}
	if rel, err := filepath.Rel(base, abs); err == nil && !strings.HasPrefix(rel, "..") {
		return rel
	}
	return path
}
