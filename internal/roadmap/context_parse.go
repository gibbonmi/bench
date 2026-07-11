package roadmap

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/structure"
)

var roadmapStartRe = regexp.MustCompile(`^\*\*([A-Za-z]+[0-9]+)(.*)$`)
var specPathRe = regexp.MustCompile(`specs/([A-Za-z0-9_-]+)\.md`)
var commandRe = regexp.MustCompile(`/bench-[A-Za-z0-9-]+`)
var ideaRe = regexp.MustCompile(`^- ([0-9]{4}-[0-9]{2}-[0-9]{2})  (.*)$`)
var promotionDateRe = regexp.MustCompile(`([0-9]{4}-[0-9]{2}-[0-9]{2})`)
var roadmapIDRe = regexp.MustCompile(`[A-Z]+[0-9]+`)

func ParseDocument(content []byte, statuses map[string]string, full bool) (Document, []ParseFailure) {
	lines := strings.Split(string(content), "\n")
	doc := Document{Text: string(content)}
	var failures []ParseFailure
	for i := 0; i < len(lines); {
		m := roadmapStartRe.FindStringSubmatch(lines[i])
		if m == nil {
			if strings.HasPrefix(lines[i], "**") {
				raw, n, tr := limited(lines[i], full)
				failures = append(failures, ParseFailure{"ROADMAP.md", "malformed roadmap row", raw, n, tr})
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
			failures = append(failures, ParseFailure{"ROADMAP.md", "unclosed roadmap heading", raw, n, tr})
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
		if slugs := specSlugsFromText([]byte(joined)); len(slugs) > 0 {
			r.Spec = slugs[0]
			r.SpecStatus = statuses[slugs[0]]
		}
		lower := strings.ToLower(joined)
		r.ExternalTrigger = strings.Contains(lower, "pending ") || strings.Contains(lower, "graduate on") || strings.Contains(lower, "scheduled")
		doc.Rows = append(doc.Rows, r)
	}
	inSequence, inFence, sequenceStart, sequenceEnd := false, false, -1, len(lines)
	for idx, line := range lines {
		trimmed := strings.TrimRight(line, " \t\r")
		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
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
	if len(content) > 0 && len(doc.Rows) == 0 && len(failures) == 0 {
		raw, n, tr := limited(string(content), full)
		failures = append(failures, ParseFailure{"ROADMAP.md", "no roadmap rows recognized", raw, n, tr})
	}
	return doc, failures
}

// SpecSlugs returns roadmap-linked spec slugs through the canonical document parser.
func SpecSlugs(content []byte) []string {
	return specSlugsFromText(content)
}

func specSlugsFromText(content []byte) []string {
	seen := map[string]bool{}
	var out []string
	inFence := false
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		for _, m := range specPathRe.FindAllStringSubmatch(line, -1) {
			slug := m[1]
			if !seen[slug] {
				seen[slug] = true
				out = append(out, slug)
			}
		}
	}
	return out
}

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
			failures = append(failures, ParseFailure{"IDEAS.md", "malformed idea row", raw, n, tr})
			continue
		}
		text, n, tr := limited(m[2], full)
		facts = append(facts, IdeaFact{m[1], text, n, tr})
	}
	return facts, failures, rawLines
}

func parsePromotions(content []byte, full bool) []PromotionFact {
	lines := strings.Split(string(content), "\n")
	var out []PromotionFact
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		lower := strings.ToLower(line)
		if !strings.HasPrefix(line, "- **") || !strings.Contains(lower, "promotion") {
			continue
		}
		start := i
		i++
		for i < len(lines) && !strings.HasPrefix(lines[i], "- **") {
			i++
		}
		i--
		bodyRaw := strings.TrimSpace(strings.Join(lines[start:i+1], "\n"))
		body, n, tr := limited(bodyRaw, full)
		date := ""
		if m := promotionDateRe.FindString(bodyRaw); m != "" {
			date = m
		}
		ids := roadmapIDRe.FindAllString(bodyRaw, -1)
		ids = unique(ids)
		kind := "promotion"
		if strings.Contains(lower, "build") {
			kind = "build"
		}
		out = append(out, PromotionFact{kind, date, strings.Trim(strings.TrimPrefix(line, "- **"), "* ."), strings.Join(ids, " "), body, n, tr})
	}
	return out
}

func unique(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

func BuildContext(root string, full bool, gate GateCacheFact) (ContextSnapshot, error) {
	s := ContextSnapshot{Full: full}
	labels := []string{"ROADMAP.md", "IDEAS.md", ".bench/learnings.md", ".bench/structure.budgets", ".bench/structure-accept", "specs/", "CHANGELOG.md"}
	data := map[string][]byte{}
	for _, label := range labels {
		if label == "specs/" {
			_, state, n, err := readDirSource(sourcePath(root, label))
			if err != nil {
				return s, fmt.Errorf("read %s: %w", label, err)
			}
			s.Sources = append(s.Sources, SourceFact{label, state, n})
			continue
		}
		b, state, err := readSource(sourcePath(root, label))
		if err != nil {
			return s, fmt.Errorf("read %s: %w", label, err)
		}
		data[label] = b
		s.Sources = append(s.Sources, SourceFact{label, state, len(b)})
	}
	gitDir, err := benchgit.Output("-C", root, "rev-parse", "--absolute-git-dir")
	if err != nil {
		return s, err
	}
	cache, cacheState, err := readSource(filepath.Join(gitDir, benchgit.GateCacheFile))
	if err != nil {
		return s, fmt.Errorf("read gate cache: %w", err)
	}
	s.Sources = append(s.Sources, SourceFact{".git/bench-last-gate", cacheState, len(cache)})

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
	s.Roadmap, roadFails = ParseDocument(data["ROADMAP.md"], statuses, full)
	s.Failures = append(s.Failures, roadFails...)
	s.Ideas, roadFails, _ = parseIdeas(data["IDEAS.md"], full)
	s.Failures = append(s.Failures, roadFails...)
	learningFacts, malformedLearnings := learnings.Parse(data[".bench/learnings.md"])
	for _, e := range learningFacts {
		body, n, tr := limited(e.Body, full)
		s.Learnings = append(s.Learnings, LearningFact{e.Date, e.Title, e.State, body, n, tr})
	}
	for _, m := range malformedLearnings {
		raw, n, tr := limited(m.Raw, full)
		s.Failures = append(s.Failures, ParseFailure{".bench/learnings.md", m.Reason, raw, n, tr})
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
	s.Git = [][]any{{gf.Branch, gf.DefaultBranch, gf.Dirty, gf.Ahead, gf.Behind}}
	sort.Slice(gf.Changes, func(i, j int) bool { return gf.Changes[i].Path < gf.Changes[j].Path })
	for _, c := range gf.Changes {
		s.GitChanges = append(s.GitChanges, []string{c.Status, c.Path})
	}
	s.GateCache = [][]any{{gate.Present, gate.Status, gate.CachedTree, gate.WorkTree, gate.Timestamp, gate.Stale}}
	s.Promotions = parsePromotions(data["CHANGELOG.md"], full)
	for i := range s.Sources {
		for _, f := range s.Failures {
			if f.Source == s.Sources[i].Source {
				s.Sources[i].State = "malformed"
			}
		}
	}
	return s, nil
}
