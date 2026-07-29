package roadmap

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/structure"
)

var roadmapStartRe = regexp.MustCompile(`^\*\*([A-Za-z]+[0-9]+)(.*)$`)
var commandRe = regexp.MustCompile(`/bench-[A-Za-z0-9-]+`)
var ideaRe = regexp.MustCompile(`^- ([0-9]{4}-[0-9]{2}-[0-9]{2})  (.*)$`)

func ParseDocument(content []byte, statuses map[string]string, full bool) (Document, []ParseFailure) {
	lines := strings.Split(string(content), "\n")
	doc := Document{Text: string(content)}
	var failures []ParseFailure
	for i := 0; i < len(lines); {
		m := roadmapStartRe.FindStringSubmatch(lines[i])
		if m == nil {
			if strings.HasPrefix(lines[i], "**") {
				raw, n, tr := limited(lines[i], full)
				failures = append(failures, ParseFailure{RoadmapFile, "malformed roadmap row", raw, n, tr})
			}
			i++
			continue
		}
		start := i
		headerParts := []string{m[2]}
		closed := strings.Contains(m[2], "**")
		for !closed && i+1 < len(lines) {
			i++
			headerParts = append(headerParts, lines[i])
			closed = strings.Contains(lines[i], "**")
		}
		if !closed {
			rawText := strings.Join(lines[start:i+1], "\n")
			raw, n, tr := limited(rawText, full)
			failures = append(failures, ParseFailure{RoadmapFile, "unclosed roadmap heading", raw, n, tr})
			continue
		}
		headerJoined := strings.Join(headerParts, "\n")
		closeAt := strings.Index(headerJoined, "**")
		title := strings.Trim(strings.ReplaceAll(headerJoined[:closeAt], "\n", " "), " —:-\t")
		inlineBody := strings.TrimSpace(headerJoined[closeAt+2:])
		i++
		bodyStart := i
		for i < len(lines) && !strings.HasPrefix(lines[i], "**") && !strings.HasPrefix(lines[i], "## ") {
			i++
		}
		bodyRaw := strings.TrimSpace(strings.Join(append([]string{inlineBody}, lines[bodyStart:i]...), "\n"))
		body, bodyBytes, truncated := limited(bodyRaw, full)
		r := RoadmapRow{ID: m[1], Title: title, Body: body, BodyBytes: bodyBytes, Truncated: truncated}
		joined := strings.Join(lines[start:i], "\n")
		if slugs := spec.LiveSpecSlugs([]byte(joined)); len(slugs) > 0 {
			r.Spec = slugs[0]
			r.SpecStatus = statuses[slugs[0]]
		}
		lower := strings.ToLower(joined)
		r.ExternalTrigger = strings.Contains(lower, "pending ") || strings.Contains(lower, "graduate on") || strings.Contains(lower, "scheduled")
		doc.Rows = append(doc.Rows, r)
	}
	inSequence, inFence, sequenceStart, sequenceEnd := false, false, -1, len(lines)
	hasSection := false
	for idx, line := range lines {
		trimmed := strings.TrimRight(line, " \t\r")
		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		// Any unfenced `## ` heading — not only Recommended sequence — is roadmap
		// structure: a document with sections is a working roadmap the reader
		// recognizes, possibly one with nothing due yet, never "not the document
		// you think".
		if strings.HasPrefix(trimmed, "## ") {
			hasSection = true
		}
		if !inSequence && trimmed == "## Recommended sequence" {
			inSequence = true
			sequenceStart = idx
			continue
		}
		if inSequence && strings.HasPrefix(trimmed, "## ") {
			sequenceEnd = idx
			break
		}
		if !inSequence {
			continue
		}
		parts := strings.SplitN(line, ". ", 2)
		if len(parts) != 2 {
			continue
		}
		rank, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		cmd := ""
		if m := commandRe.FindString(parts[1]); m != "" {
			cmd = m
		}
		doc.Sequence = append(doc.Sequence, SequenceRow{Rank: rank, Text: parts[1], Command: cmd})
	}
	if sequenceStart >= 0 {
		doc.SequenceText = strings.Join(lines[sequenceStart:sequenceEnd], "\n")
		if !strings.HasSuffix(doc.SequenceText, "\n") {
			doc.SequenceText += "\n"
		}
	}
	if len(content) > 0 && len(doc.Rows) == 0 && len(failures) == 0 && !hasSection {
		raw, n, tr := limited(string(content), full)
		failures = append(failures, ParseFailure{RoadmapFile, noRoadmapRowsReason, raw, n, tr})
	}
	return doc, failures
}

// noRoadmapRowsReason is roadmap's unsupported-schema predicate — bytes that read
// cleanly but in which ParseDocument recognized no roadmap structure at all: no
// `**ID**` row and no unfenced `## ` section, Recommended sequence or otherwise. A
// document with a section is a working roadmap, possibly an early-stage or fully
// drained one, and that is "nothing to report" — the empty state, not
// unsupported-schema. Named once so the AXI surface (RoadmapCommand) and the
// --context snapshot (BuildContext) agree on exactly which parser failure means
// "unsupported-schema" rather than "malformed".
const noRoadmapRowsReason = "no roadmap rows recognized"

func parseIdeas(content []byte, full bool) ([]IdeaFact, []ParseFailure, []string) {
	var facts []IdeaFact
	var failures []ParseFailure
	var rawLines []string
	for _, line := range strings.Split(string(content), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.HasPrefix(line, "- ") {
			rawLines = append(rawLines, line)
		}
		m := ideaRe.FindStringSubmatch(line)
		if m == nil {
			raw, n, tr := limited(line, full)
			failures = append(failures, ParseFailure{IdeasFile, "malformed idea row", raw, n, tr})
			continue
		}
		text, n, tr := limited(m[2], full)
		facts = append(facts, IdeaFact{m[1], text, n, tr})
	}
	return facts, failures, rawLines
}

func BuildContext(root string, full bool, gate GateCacheFact) (ContextSnapshot, error) {
	s := ContextSnapshot{Full: full}
	labels := []string{RoadmapFile, IdeasFile, learnings.JournalPath, ".bench/structure.budgets", ".bench/structure-accept", "specs/"}
	data := map[string][]byte{}
	for _, label := range labels {
		if label == "specs/" {
			cd := bounds.ClassifyDir(sourcePath(root, label))
			s.Sources = append(s.Sources, SourceFact{label, string(cd.State), dirBytes(cd.Entries)})
			if degradedState(cd.State) {
				s.Failures = append(s.Failures, ParseFailure{label, string(cd.State) + ": " + cd.Reason, "", 0, false})
			}
			continue
		}
		c := bounds.Classify(sourcePath(root, label), bounds.ControlRecordLimit)
		if c.State == bounds.StateParsed {
			data[label] = c.Data
		}
		s.Sources = append(s.Sources, SourceFact{label, string(c.State), len(c.Data)})
		if degradedState(c.State) {
			s.Failures = append(s.Failures, ParseFailure{label, string(c.State) + ": " + c.Reason, "", 0, false})
		}
	}
	cacheState := "absent"
	if gate.Present {
		cacheState = "present"
	}
	s.Sources = append(s.Sources, SourceFact{".git/bench-last-gate", cacheState, gate.CacheBytes})

	specFacts, err := spec.Facts(root)
	if err != nil {
		return s, err
	}
	statuses := map[string]string{}
	for _, f := range specFacts {
		statuses[f.Slug] = f.Status
		s.Specs = append(s.Specs, []string{f.Slug, f.Status, f.RoadmapID})
		hist, e := spec.History(f.Slug)
		if e != nil {
			return s, e
		}
		for _, h := range hist {
			s.SpecHistory = append(s.SpecHistory, []string{h.Slug, h.Hash, h.Date, h.Kind, h.Subject})
		}
	}
	var roadFails []ParseFailure
	s.Roadmap, roadFails = ParseDocument(data[RoadmapFile], statuses, full)
	s.Failures = append(s.Failures, roadFails...)
	s.Ideas, roadFails, _ = parseIdeas(data[IdeasFile], full)
	s.Failures = append(s.Failures, roadFails...)
	learningFacts, malformedLearnings := learnings.Parse(data[learnings.JournalPath])
	for _, e := range learningFacts {
		body, n, tr := limited(e.Body, full)
		s.Learnings = append(s.Learnings, LearningFact{e.Date, e.Title, e.State, body, n, tr})
	}
	for _, m := range malformedLearnings {
		raw, n, tr := limited(m.Raw, full)
		s.Failures = append(s.Failures, ParseFailure{learnings.JournalPath, m.Reason, raw, n, tr})
	}
	structFacts, err := structure.Facts(root)
	if err != nil {
		return s, err
	}
	for _, f := range structFacts {
		s.Structure = append(s.Structure, []any{f.Kind, f.Path, f.Actual, f.Limit, f.State, f.Detail})
	}
	gf, err := benchgit.Facts(root)
	if err != nil {
		return s, err
	}
	defaultBranch, ahead, behind := any(gf.DefaultBranch), any(gf.Ahead), any(gf.Behind)
	if !gf.DefaultResolved {
		// The three cells the unresolved default makes unknowable, named rather than
		// rendered as the zeros they would otherwise fabricate.
		defaultBranch, ahead, behind = "unknown", "unknown", "unknown"
	}
	s.Git = [][]any{{gf.Branch, defaultBranch, gf.Dirty, ahead, behind}}
	sort.Slice(gf.Changes, func(i, j int) bool { return gf.Changes[i].Path < gf.Changes[j].Path })
	for _, c := range gf.Changes {
		s.GitChanges = append(s.GitChanges, []string{c.Status, c.Path})
	}
	s.GateCache = [][]any{{gate.Present, gate.State, gate.PendingStatus, gate.Status, gate.CachedTree, gate.WorkTree, gate.Timestamp, gate.Stale}}
	for i := range s.Sources {
		// A source already carrying a classifier-level state (unreadable, wrong-type,
		// a byte-level malformed read) keeps it: the parser-level failures below only
		// ever fire against bytes the classifier already handed over as parsed.
		if s.Sources[i].State != string(bounds.StateParsed) {
			continue
		}
		for _, f := range s.Failures {
			if f.Source != s.Sources[i].Source {
				continue
			}
			if f.Reason == noRoadmapRowsReason {
				// unsupported-schema is the parser's own state: the bytes read fine
				// but roadmap.go recognized no row in them, distinct from a
				// byte-level malformed read the classifier would have caught.
				s.Sources[i].State = string(bounds.StateUnsupportedSchema)
			} else {
				s.Sources[i].State = string(bounds.StateMalformed)
			}
		}
	}
	return s, nil
}

// degradedState reports whether a classifier state needs its own parse_failures row
// naming it — every state but the two that carry no problem to explain.
func degradedState(state bounds.FileState) bool {
	return state != bounds.StateAbsent && state != bounds.StateEmpty && state != bounds.StateParsed
}
