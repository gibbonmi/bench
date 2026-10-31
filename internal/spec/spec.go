// Package spec owns spec-file addressing and the spec-lifecycle operations. Retire
// implements the `bench spec retire <slug>` deletion, history reports the commits that
// retired or deleted a spec, and widen adds paths to one ticket's Writes: line. Resolve
// is the one source of the spec-argument convention: path-first, then a
// specs/<slug>/spec.md fallback. `bench coverage`, retire, and history resolve their
// argument through it.
//
// Implemented is the single source of the status-line flip. It derives the bytes that
// turn exactly one line-start `Status: staged` into the retirement-detector form
// `Status: implemented`, preserving every other byte, and `bench worktree land --spec`
// — the one author of the flip — composes them. AwaitsRetirement is the source of the
// merged-implemented predicate — the `implemented` twin of the staged form — shared by
// retire's validation and the `bench status` specs-awaiting-retirement counter.
package spec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"unicode"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/markdown"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// stagedRe matches the flip's input: a `Status:` line whose sole value is `staged`,
// tab/space separated, with only whitespace trailing. This is the `staged` twin of the
// retirement detector's `^Status:[ \t]+implemented[ \t]*$`. Swapping `staged` for
// `implemented` in place therefore yields exactly the detector's accepted form.
var stagedRe = regexp.MustCompile(`^Status:[ \t]+staged[ \t]*$`)

var liveSpecPathTokenRe = regexp.MustCompile(`specs/([A-Za-z0-9_-]+)/spec\.md`)

// Fact is one typed live spec record. Path is its repository-relative path.
type Fact struct {
	Slug, Path, Status, RoadmapID string
}

// LiveSpecSlugs enumerates distinct live folder-spec slugs that the header lines name.
func LiveSpecSlugs(content []byte) []string {
	seen := map[string]bool{}
	var slugs []string
	for _, line := range headerLines(content) {
		for _, match := range liveSpecPathTokenRe.FindAllStringSubmatch(line.Raw, -1) {
			slug := match[1]
			if !seen[slug] {
				seen[slug] = true
				slugs = append(slugs, slug)
			}
		}
	}
	return slugs
}

// headerLines is the header scan's view of spec content: the body lines of the block reader.
// Every spec grammar matches the raw text of these lines, so a fenced line or a frontmatter
// line declares nothing.
func headerLines(content []byte) []markdown.Line {
	lines := markdown.Read(string(content)).Lines
	return slices.DeleteFunc(lines, func(line markdown.Line) bool { return line.Class != markdown.Body })
}

const (
	statusPrefix  = "Status:"
	roadmapPrefix = "Roadmap:"
)

func metadata(content []byte) (status, roadmapID string) {
	statusLine, roadmapLine := metadataLines(content)
	return fieldValue(statusLine, statusPrefix), fieldValue(roadmapLine, roadmapPrefix)
}

// metadataLines selects the last visible Status: line and the last visible Roadmap: line,
// or nil for a field that no visible line declares.
func metadataLines(content []byte) (status, roadmap *markdown.Line) {
	lines := headerLines(content)
	for i := range lines {
		switch raw := lines[i].Raw; {
		case strings.HasPrefix(raw, statusPrefix):
			status = &lines[i]
		case strings.HasPrefix(raw, roadmapPrefix):
			roadmap = &lines[i]
		}
	}
	return status, roadmap
}

// fieldValue answers the trimmed value of a metadata line, or "" for no line.
func fieldValue(line *markdown.Line, prefix string) string {
	if line == nil {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(line.Raw, prefix))
}

// ErrPersistedFieldsInvalid is the sentinel that each PersistedFields refusal wraps.
var ErrPersistedFieldsInvalid = errors.New("persisted-fields-invalid")

// RecordIntentLedger and RecordReviewRecord are the record names that a declaration item can
// name: the shared intent ledger and the review record.
const (
	RecordIntentLedger = "intent-ledger"
	RecordReviewRecord = "review-record"
)

// persistedRecords is the closed set of record names that a declaration item can name.
var persistedRecords = map[string]bool{RecordIntentLedger: true, RecordReviewRecord: true}

// PersistedField is one `<record>:<path>` item of the `Persisted fields:` declaration. Path
// holds the dot-separated JSON tag names from the record root.
type PersistedField struct {
	Record string
	Path   []string
}

// PersistedFields reads the `Persisted fields:` declaration through headerLines, the shared
// fence-aware line view, so a fenced line declares nothing. A spec with no line returns an
// empty list. The value is a comma-separated item list. A malformed item or a second
// unfenced line refuses with ErrPersistedFieldsInvalid; the second-line refusal keeps a body
// line from replacing the reviewed header declaration.
func PersistedFields(content []byte) ([]PersistedField, error) {
	value, declared := "", false
	for _, line := range headerLines(content) {
		v, ok := strings.CutPrefix(line.Raw, "Persisted fields:")
		if !ok {
			continue
		}
		if declared {
			return nil, fmt.Errorf("%w: a second Persisted fields: line", ErrPersistedFieldsInvalid)
		}
		value, declared = v, true
	}
	if !declared {
		return []PersistedField{}, nil
	}
	var fields []PersistedField
	seen := map[string]bool{}
	for _, raw := range strings.Split(value, ",") {
		item := strings.TrimSpace(raw)
		record, path, _ := strings.Cut(item, ":")
		segments := strings.Split(path, ".")
		switch {
		case item == "":
			return nil, fmt.Errorf("%w: empty item", ErrPersistedFieldsInvalid)
		case !persistedRecords[record]:
			return nil, fmt.Errorf("%w: unknown record %q in %q", ErrPersistedFieldsInvalid, record, item)
		case slices.Contains(segments, ""):
			return nil, fmt.Errorf("%w: empty path segment in %q", ErrPersistedFieldsInvalid, item)
		case strings.ContainsFunc(path, func(r rune) bool { return r == ':' || unicode.IsSpace(r) }):
			return nil, fmt.Errorf("%w: whitespace or a colon in the path of %q", ErrPersistedFieldsInvalid, item)
		case seen[item]:
			return nil, fmt.Errorf("%w: duplicate item %q", ErrPersistedFieldsInvalid, item)
		}
		seen[item] = true
		fields = append(fields, PersistedField{Record: record, Path: segments})
	}
	return fields, nil
}

// Facts reads folder specs in slug order. It classifies every present spec.md before
// opening it. Only bounds.StateParsed supplies metadata; every other bounds.FileState
// stays an empty-metadata evidence row instead of an omission or a returned read
// error. This fail-closed posture gives ambient consumers evidence without letting a
// special file block.
func Facts(root string) ([]Fact, error) {
	entries, err := os.ReadDir(filepath.Join(root, "specs"))
	if err != nil {
		return []Fact{}, nil
	}
	out := make([]Fact, 0, len(entries))
	for _, entry := range entries {
		candidate := filepath.Join(root, "specs", entry.Name())
		info, err := os.Stat(candidate)
		if err != nil || !info.IsDir() {
			continue
		}
		path := filepath.Join(candidate, "spec.md")
		if _, err := os.Lstat(path); err != nil {
			continue
		}
		f := Fact{
			Slug: entry.Name(),
			Path: specPath(entry.Name()),
		}
		c := bounds.Classify(path, bounds.ControlRecordLimit)
		if c.State == bounds.StateParsed {
			f.Status, f.RoadmapID = metadata(c.Data)
		}
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Slug < out[j].Slug })
	return out, nil
}

// AwaitsRetirement reports whether spec content carries an unfenced retirement marker.
// This is the one definition of "a merged spec awaiting retirement." The `bench
// status` specs row shares it and counts it across specs/. `bench spec retire` shares
// it too, and validates it before deleting.
//
// It reads the `Status:` value through headerLines, so a fenced line never marks the
// content.
func AwaitsRetirement(content []byte) bool {
	status, _ := metadata(content)
	return status == StatusImplemented
}

// StatusImplemented is the Status value a finished spec carries. It is exported as the
// one source of the token. A caller that classifies a Fact by status compares against
// the same spelling the retirement predicate does.
const StatusImplemented = "implemented"

// Implemented derives the flipped bytes without mutating a checkout. Prospective
// composition uses it, so the caller authorizes lifecycle bytes before publication.
func Implemented(content []byte) ([]byte, error) {
	implemented, matches := deriveImplemented(content)
	if matches > 1 {
		return nil, errors.New("spec has more than one Status: staged line")
	}
	if matches == 0 {
		return nil, errors.New("spec has no Status: staged line")
	}
	return implemented, nil
}

// deriveImplemented flips the one staged status line in the source bytes by its line offset, so
// each other byte, a carriage return included, stays.
func deriveImplemented(content []byte) ([]byte, int) {
	matches, at := 0, 0
	for _, line := range headerLines(content) {
		if stagedRe.MatchString(line.Raw) {
			matches++
			at = line.Offset + strings.Index(line.Raw, "staged")
		}
	}
	if matches != 1 {
		return content, matches
	}
	return slices.Concat(content[:at], []byte(StatusImplemented), content[at+len("staged"):]), matches
}

// Command implements `bench spec <subcommand> <slug>`. Its help lists `show`, `retire`,
// `history`, and `widen`. The root command adapter dispatches `show` with ShowGrammar, so
// it can compose the spec selections with other owners; see read.go. `retire` deletes a
// merged spec and its review pickup. `history` reports the commits that retired or deleted
// it; see history.go. `widen` adds paths to one ticket's Writes: line inside the
// signed-off scope; see widen.go.
//
// The slug's specs/<slug>/spec.md fallback is anchored at the repo root, so it
// resolves from any cwd inside the repo, while a path argument stays cwd-relative. A
// usage error — a missing or unknown subcommand, a missing or extra argument, or an
// unknown flag — exits 2. A resolve, validate, or delete failure exits 1, naming the
// file and the reason. retireOperand is the operand a retire parsed, or "" for other runs.
func Command(args []string) (out string, code int, retireOperand string) {
	if len(args) == 0 {
		return toon.Usage("bench spec", "expected a subcommand: "+ShowSubcommand+", retire, history, widen") + "\n", 2, ""
	}
	switch args[0] {
	case "-h", "--help":
		return "usage: bench spec <subcommand>\n\nsubcommands:\n  " + strings.TrimPrefix(ShowUsage, "bench spec ") + "\n  retire <slug>\n  history <slug>\n  " + strings.TrimPrefix(WidenUsage, "bench spec ") + "\n", 0, ""
	case "retire":
		arg, out, code, ok := specArg("bench spec retire", "usage: bench spec retire <spec.md | slug>\n", args[1:])
		if ok {
			out, code = retireCommand(arg)
		}
		return out, code, arg
	case "history":
		out, code = historyCommand(args[1:])
	case "widen":
		out, code = widenCommand(args[1:])
	default:
		out, code = toon.Usage("bench spec", args[0])+"\n", 2
	}
	return out, code, ""
}

// specArg extracts the single positional that a subcommand grammar has already
// validated. Parser-owned help and usage errors retain their standard rendering.
func specArg(cmd, help string, rest []string) (arg, out string, code int, ok bool) {
	grammar := usage.Grammar{Cmd: cmd, Help: strings.TrimSuffix(help, "\n"), MinArgs: 1, MaxArgs: 1}
	parsed, line, code := usage.Parse(grammar, rest)
	if line != "" {
		return "", line + "\n", code, false
	}
	return parsed.Positionals[0], "", 0, true
}

// RepoBase returns the repo root that anchors the specs/<slug>/spec.md fallback. It
// returns "" when the cwd is outside a repo, so the fallback then resolves relative to
// the process cwd. It is the one source of that anchoring for every subcommand that
// resolves a bare slug, so the caller's cwd never changes which spec a slug names.
func RepoBase() string {
	if root, err := git.Root(); err == nil {
		return root
	}
	return ""
}

// SlugOf is the spec slug of a path or bare argument, which retire, history, and the retire
// listing share. A folder spec's filename is always spec.md, so its parent folder supplies it.
func SlugOf(arg string) string {
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

// RelTo renders path repo-relative to base for stable command output. It falls back
// to the path verbatim when base is empty or path lies outside base. It is the one
// source of that rendering, so a resolved spec reads the same from any cwd. The
// rendering never leaks the machine's directory layout into agent-facing output.
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
