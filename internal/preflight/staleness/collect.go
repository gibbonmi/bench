package staleness

import (
	"fmt"
	"path"
	"slices"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/preflight/ownership"
	specref "github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/tickets"
)

// unsupported names each mechanical check that no collector grades yet. Each one stays a
// diagnostic, so the audit stays incomplete and manual review stays required.
var unsupported = []string{CheckCitedSymbols, CheckCitedTests}

// Collect reads the audit scope of the committed pair that source names at root. It reads
// the spec and its planning-facts inventory at both identities from committed objects, never
// from the checkout. sourceFacts is the closure snapshot of the same pair, whose diagnostics
// stay open in the audit. An absent or special spec at either identity, an unreadable tree,
// or a member of the scope that is no regular file refuses the collection with an *Error and
// no facts.
func Collect(root string, source ownership.Source, sourceFacts ownership.Snapshot) (AuditFacts, error) {
	baseline, err := readSpec(root, source.Base, source.SpecPath, StageBaselineSpec)
	if err != nil {
		return AuditFacts{}, err
	}
	current, err := readSpec(root, source.Tip, source.SpecPath, StageTipSpec)
	if err != nil {
		return AuditFacts{}, err
	}
	entries, unread := scope(source.SpecPath, [][]byte{baseline, current}, []ownership.Manifest{ownership.ParseManifest(baseline), ownership.ParseManifest(current)})
	before, err := members(root, source.Base, entries)
	if err != nil {
		return AuditFacts{}, err
	}
	after, err := members(root, source.Tip, entries)
	if err != nil {
		return AuditFacts{}, err
	}
	facts := AuditFacts{Baseline: source.Base, Tip: source.Tip, SpecPath: source.SpecPath, Paths: union(before, after)}
	for _, diagnostic := range sourceFacts.Diagnostics {
		facts.Diagnostics = append(facts.Diagnostics, Diagnostic{Subject: diagnostic.ID, Reason: diagnostic.Reason})
	}
	facts.Diagnostics = append(facts.Diagnostics, unread...)
	for _, check := range unsupported {
		facts.Diagnostics = append(facts.Diagnostics, Diagnostic{Subject: check, Reason: "no " + check + " collector grades this audit"})
	}
	return facts, nil
}

// readSpec reads the committed spec bytes at commit. A spec that is absent there is a
// refusal, never an empty baseline.
func readSpec(root, commit, specPath, stage string) ([]byte, error) {
	content, err := git.ReadTreeFile(root, commit, specPath)
	if err != nil {
		return nil, &Error{Stage: stage, Subject: specPath, Cause: fmt.Errorf("at %s: %w", commit, err)}
	}
	return content, nil
}

// scope is every entry the audit enumerates: the active spec and its tickets, the ownership
// fence entries of both spec identities, and every path that either manifest names. Each
// entry takes the one Writes: split, so a marker or a trailing slash never becomes a path
// segment. unread holds one diagnostic for each premise whose subject paths do not read.
func scope(specPath string, specs [][]byte, manifests []ownership.Manifest) (split []string, unread []Diagnostic) {
	entries := []string{specPath, path.Join(path.Dir(specPath), "tickets")}
	for _, content := range specs {
		fences, _ := specref.FenceTokens(content)
		entries = append(entries, fences...)
	}
	for _, manifest := range manifests {
		paths, diagnostics := claimSources(manifest)
		entries = append(entries, paths...)
		for _, diagnostic := range diagnostics {
			if !slices.Contains(unread, diagnostic) {
				unread = append(unread, diagnostic)
			}
		}
	}
	for _, entry := range entries {
		if p, _ := tickets.WritesPath(entry); p != "" && !slices.Contains(split, p) {
			split = append(split, p)
		}
	}
	return split, unread
}

// claimSources is every path one manifest names: its scope paths, each item location, each
// declared owner, each test seam, each decision source, each premise subject path, and each
// native witness phase owner. A premise whose subject paths do not read is a diagnostic,
// because its claim sources then stay outside the audit scope.
func claimSources(m ownership.Manifest) ([]string, []Diagnostic) {
	var paths []string
	var unread []Diagnostic
	for _, item := range m.Scope {
		paths = append(append(paths, item.Location.Path), item.Paths...)
	}
	for _, change := range m.Changes {
		paths = append(paths, change.Location.Path, change.Owner.Path)
	}
	for _, premise := range m.Premises {
		subject, err := ownership.SubjectPaths(premise)
		if err != nil {
			unread = append(unread, Diagnostic{Subject: premise.ID, Reason: "the premise subject paths stay outside the audit scope: " + err.Error()})
		}
		paths = append(append(paths, premise.Location.Path), subject...)
	}
	for _, seam := range m.Seams {
		paths = append(paths, seam.Location.Path, seam.Path)
	}
	for _, decision := range m.Decisions {
		paths = append(paths, decision.Location.Path, decision.Source.Path)
	}
	for _, claim := range m.ManualClaims {
		paths = append(paths, claim.Location.Path)
		if claim.NativeWitness != nil {
			paths = append(paths, claim.NativeWitness.PhaseOwner.Path)
		}
	}
	return paths, unread
}

// members is every tracked member of the scope entries at commit, by path. An entry covers
// a path that it names exactly or contains at a segment boundary, and the recursive tree
// listing holds every tracked path, a dot path included. A member that is no regular file
// refuses the collection, so the audit never reads through a link or a gitlink.
func members(root, commit string, entries []string) (map[string]Object, error) {
	listing, err := git.ListTree(root, "-r", "--full-tree", commit)
	if err != nil {
		return nil, &Error{Stage: StageMembers, Subject: commit, Cause: err}
	}
	found := map[string]Object{}
	for _, member := range listing {
		if !slices.ContainsFunc(entries, func(entry string) bool { return tickets.Covers(entry, member.Path) }) {
			continue
		}
		if !member.IsRegularFile() {
			return nil, &Error{Stage: StageSpecial, Subject: member.Path, Cause: fmt.Errorf("mode %s at %s is no regular file", member.Mode, commit)}
		}
		found[member.Path] = Object{Mode: member.Mode, ID: member.Object}
	}
	return found, nil
}
