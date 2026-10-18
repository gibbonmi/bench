package staleness

import (
	"encoding/json"
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/preflight/ownership"
	specref "github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/tickets"
)

// unsupported names each mechanical check that no collector grades yet. Each one stays a
// diagnostic, so the audit stays incomplete and manual review stays required.
var unsupported = []string{CheckCitedSymbols, CheckCitedTests}

// Collect reads the audit scope of the committed pair that source names at root. It reads
// the spec at both identities from committed objects, never from the checkout. manifest is
// the planning-facts inventory of the tip spec, and sourceFacts is the closure snapshot of the
// same pair, whose diagnostics stay open in the audit. An absent or special spec at either
// identity, an unreadable tree, or a member of the scope that is no regular file refuses the
// collection with an *Error and no facts.
func Collect(root string, source ownership.Source, manifest ownership.Manifest, sourceFacts ownership.Snapshot) (AuditFacts, error) {
	baseline, err := readSpec(root, source.Base, source.SpecPath, StageBaselineSpec)
	if err != nil {
		return AuditFacts{}, err
	}
	current, err := readSpec(root, source.Tip, source.SpecPath, StageTipSpec)
	if err != nil {
		return AuditFacts{}, err
	}
	entries := scope(source.SpecPath, [][]byte{baseline, current}, []ownership.Manifest{ownership.ParseManifest(baseline), manifest})
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
// segment.
func scope(specPath string, specs [][]byte, manifests []ownership.Manifest) []string {
	entries := []string{specPath, path.Join(path.Dir(specPath), "tickets")}
	for _, content := range specs {
		fences, _ := specref.FenceTokens(content)
		entries = append(entries, fences...)
	}
	for _, manifest := range manifests {
		entries = append(entries, claimSources(manifest)...)
	}
	var split []string
	for _, entry := range entries {
		if p, _ := tickets.WritesPath(entry); p != "" && !slices.Contains(split, p) {
			split = append(split, p)
		}
	}
	return split
}

// claimSources is every path one manifest names: its scope paths, each item location, each
// declared owner, each test seam, each decision source, each premise subject path, and each
// native witness phase owner.
func claimSources(m ownership.Manifest) []string {
	var paths []string
	for _, item := range m.Scope {
		paths = append(append(paths, item.Location.Path), item.Paths...)
	}
	for _, change := range m.Changes {
		paths = append(paths, change.Location.Path, change.Owner.Path)
	}
	for _, premise := range m.Premises {
		paths = append(paths, premise.Location.Path, subjectPath(premise.Subject))
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
	return paths
}

// subjectPath is the path that a premise subject names, or empty when its tagged object
// names no path.
func subjectPath(subject json.RawMessage) string {
	var named struct {
		Path string `json:"path"`
	}
	if json.Unmarshal(subject, &named) != nil {
		return ""
	}
	return named.Path
}

// members is every tracked member of the scope entries at commit, by path. An entry covers
// a path that it names exactly or contains at a segment boundary, and the recursive tree
// listing holds every tracked path, a dot path included. A member that is no regular file
// refuses the collection, so the audit never reads through a link or a gitlink.
func members(root, commit string, entries []string) (map[string]Object, error) {
	raw, err := git.Raw("-C", root, "ls-tree", "-r", "-z", "--full-tree", commit)
	if err != nil {
		return nil, &Error{Stage: StageMembers, Subject: commit, Cause: err}
	}
	found := map[string]Object{}
	for _, record := range strings.Split(string(raw), "\x00") {
		if record == "" {
			continue
		}
		meta, member, ok := strings.Cut(record, "\t")
		fields := strings.Fields(meta)
		if !ok || len(fields) != 3 {
			return nil, &Error{Stage: StageMembers, Subject: commit, Cause: fmt.Errorf("unreadable tree entry %q", record)}
		}
		if !slices.ContainsFunc(entries, func(entry string) bool { return tickets.Covers(entry, member) }) {
			continue
		}
		if !(git.IndexEntry{Mode: fields[0]}).IsRegularFile() {
			return nil, &Error{Stage: StageSpecial, Subject: member, Cause: fmt.Errorf("mode %s at %s is no regular file", fields[0], commit)}
		}
		found[member] = Object{Mode: fields[0], ID: fields[2]}
	}
	return found, nil
}
