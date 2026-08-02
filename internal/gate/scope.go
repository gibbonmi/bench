package gate

// The reduced-run declaration: the repository paths the excludable phases cannot
// observe, and the phases that claim they cannot. The two lists are halves of one
// claim, so they live together — split apart, one half moves without the other.

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/canary"
)

// Scope is the declaration. Its fields are unexported and reached through accessors so
// the single source cannot be edited in place by a consumer holding a copy.
type Scope struct {
	directories []string
	files       []string
	excludable  []string
	included    []string
}

// ReducedScope is the kit's declaration. A directory entry carries its trailing slash:
// membership under it is location, not enumeration, so a surface that lands inside it
// tomorrow is covered without an edit here. `.bench-notes.md` is per-worktree shift
// scratch rather than repository capture, and is declared because the same paths govern
// the ambient staleness signal that already carried it.
//
// The build phase is in neither phase list. It produces the binary the other phases
// exec, so it runs in both modes: skipping it would leave a reduced run nothing to exec,
// and calling it included would claim it grades the declared paths, which it does not.
func ReducedScope() Scope {
	return Scope{
		directories: []string{"capture/", "specs/"},
		files:       []string{".bench-notes.md", "ROADMAP.md"},
		excludable: []string{
			canary.PhaseGofmt, canary.PhaseVet, canary.PhaseTest, canary.PhaseRace,
			canary.PhaseContract, "shellcheck", "canary",
		},
		included: []string{conformancePhaseName, canary.PhaseConformanceSuite},
	}
}

// declaresDirectory reports whether a declared entry names a directory rather than a file.
// The trailing slash is the marker, and this is where that convention is read.
func declaresDirectory(entry string) bool { return strings.HasSuffix(entry, "/") }

// declaredEntryCovers reports whether a declared entry covers a repository-relative,
// slash-separated path. The two entry kinds are told apart by declaresDirectory: a
// directory entry covers every descendant, and a file entry is matched byte-for-byte, with
// no case folding and no Unicode normalization — so a homoglyph filename is not covered,
// the direction that costs a full gate rather than an ungraded change.
//
// This is the one place membership is decided. Both the reduced-run declaration and the
// per-component input declarations name their surfaces in this shape, and a second
// membership rule beside it would let one of them cover a file the other did not.
func declaredEntryCovers(entry, path string) bool {
	if declaresDirectory(entry) {
		return strings.HasPrefix(path, entry)
	}
	return path == entry
}

// Member reports whether a repository-relative, slash-separated path is declared.
func (s Scope) Member(path string) bool {
	if !repositoryRelative(path) {
		return false
	}
	for _, entry := range slices.Concat(s.files, s.directories) {
		if declaredEntryCovers(entry, path) {
			return true
		}
	}
	return false
}

// Confines reports whether a changeset is entirely declared. Every path must be a
// member: one unlisted path in the set is a change the excludable phases can observe,
// and an empty changeset inherits nothing, so both answer no.
func (s Scope) Confines(paths []string) bool {
	if len(paths) == 0 {
		return false
	}
	for _, path := range paths {
		if !s.Member(path) {
			return false
		}
	}
	return true
}

// Excludable reports whether a phase is declared unable to observe the paths, and so
// may inherit a full-green ancestor's evidence for them.
func (s Scope) Excludable(phase string) bool { return slices.Contains(s.excludable, phase) }

// ExcludablePhases are the phases a reduced run skips.
func (s Scope) ExcludablePhases() []string { return slices.Clone(s.excludable) }

// splitTable partitions a resolved phase table by the declaration: the phases a reduced
// run executes — the complement of the excludable set, which carries the build phase
// because it produces the binary the others exec — and the names it skips. Every
// consumer of "what does a reduced run keep" derives from this one partition.
func (s Scope) splitTable(table []Phase) (run []Phase, skipped []string) {
	for _, phase := range table {
		if s.Excludable(phase.Name) {
			skipped = append(skipped, phase.Name)
			continue
		}
		run = append(run, phase)
	}
	return run, skipped
}

// IncludedPhaseNames derives the included set from the phase table that actually runs:
// root's resolved table minus the excludable declaration, minus the build phase (in
// neither declared list). The profile's included-phases row and the advertisement
// accessor below are both graded against this derivation, so a non-excludable phase
// added to the table cannot join every reduced run while the prose keeps advertising
// the old set.
func IncludedPhaseNames(root, kit string) ([]string, error) {
	table, err := phaseTable(root, kit)
	if err != nil {
		return nil, err
	}
	run, _ := ReducedScope().splitTable(table)
	names := make([]string, 0, len(run))
	for _, phase := range run {
		if phase.Name != canary.PhaseBuild {
			names = append(names, phase.Name)
		}
	}
	return names, nil
}

// IncludedPhases advertises the phases a reduced run executes against the real tree
// (minus the build phase, declared in neither list). What actually runs is the table
// complement IncludedPhaseNames derives; this list is the declaration's readable summary
// of it, pinned equal to that derivation over the kit's own table by a unit test, so it
// cannot drift from what runs.
func (s Scope) IncludedPhases() []string { return slices.Clone(s.included) }

// Files are the declared file entries, matched byte-for-byte.
func (s Scope) Files() []string { return slices.Clone(s.files) }

// Directories are the declared directory entries, each with its trailing slash. Every
// descendant is declared; the directory itself names no file.
func (s Scope) Directories() []string { return slices.Clone(s.directories) }

// repositoryRelative applies the lexical containment the gate input manifest's loader
// requires of a declared path: slash-separated, already clean, and unable to name
// anything above the root. A path that escapes is refused rather than resolved, because
// resolution would let a traversal reach a declared directory from outside it.
func repositoryRelative(path string) bool {
	if path == "" || path == "." || path == ".." || strings.HasPrefix(path, "../") {
		return false
	}
	if filepath.IsAbs(path) || hasUnsafeText(path) {
		return false
	}
	return filepath.ToSlash(filepath.FromSlash(path)) == path &&
		filepath.Clean(filepath.FromSlash(path)) == filepath.FromSlash(path)
}
