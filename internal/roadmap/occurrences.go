package roadmap

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/retros"
)

var occurrenceTokenRe = regexp.MustCompile(`^\[occurrence:([^/]*)/([^\]]*)\]$`)

type captureUnit struct {
	Source, CaptureUnit, Body string
}

func projectCaptureOccurrences(doc *Document, units []captureUnit) ([]CaptureOccurrence, []OccurrencePair) {
	owners := make(map[string]map[string]bool, len(doc.Rows))
	for _, row := range doc.Rows {
		keys := map[string]bool{}
		for _, key := range strings.Split(row.OccurrenceKeys, ", ") {
			if key != "" {
				keys[key] = true
			}
		}
		owners[row.ID] = keys
	}
	var occurrences []CaptureOccurrence
	pending := map[OccurrencePair]struct{}{}
	for _, unit := range units {
		owner, incident, kind, cited := finalOccurrence(unit.Body)
		if !cited {
			continue
		}
		if kind != "" {
			doc.OccurrenceDiscrepancies = append(doc.OccurrenceDiscrepancies, OccurrenceDiscrepancy{Source: unit.Source, CaptureUnit: unit.CaptureUnit, Kind: kind, Owner: owner, Incident: incident, Structural: true})
			continue
		}
		keys, current := owners[owner]
		if !current {
			doc.OccurrenceDiscrepancies = append(doc.OccurrenceDiscrepancies, OccurrenceDiscrepancy{Source: unit.Source, CaptureUnit: unit.CaptureUnit, Kind: "unknown-owner", Owner: owner, Incident: incident, Structural: true})
			continue
		}
		state := "pending"
		if keys[incident] {
			state = "already-recorded"
			doc.OccurrenceDiscrepancies = append(doc.OccurrenceDiscrepancies, OccurrenceDiscrepancy{Source: unit.Source, CaptureUnit: unit.CaptureUnit, Kind: "already-recorded", Owner: owner, Incident: incident})
		}
		occurrences = append(occurrences, CaptureOccurrence{Owner: owner, Incident: incident, Source: unit.Source, CaptureUnit: unit.CaptureUnit, State: state})
		if state == "pending" {
			pending[OccurrencePair{Owner: owner, Incident: incident}] = struct{}{}
		}
	}
	sort.Slice(occurrences, func(i, j int) bool {
		left, right := occurrences[i], occurrences[j]
		return occurrenceOrder(left.Owner, left.Incident, left.Source, left.CaptureUnit, left.State) < occurrenceOrder(right.Owner, right.Incident, right.Source, right.CaptureUnit, right.State)
	})
	sort.Slice(doc.OccurrenceDiscrepancies, func(i, j int) bool {
		left, right := doc.OccurrenceDiscrepancies[i], doc.OccurrenceDiscrepancies[j]
		return occurrenceOrder(left.Source, left.CaptureUnit, left.Kind, left.Owner, left.Incident, strconv.FormatBool(left.Structural)) < occurrenceOrder(right.Source, right.CaptureUnit, right.Kind, right.Owner, right.Incident, strconv.FormatBool(right.Structural))
	})
	pairs := make([]OccurrencePair, 0, len(pending))
	for pair := range pending {
		pairs = append(pairs, pair)
	}
	sort.Slice(pairs, func(i, j int) bool {
		return occurrenceOrder(pairs[i].Owner, pairs[i].Incident) < occurrenceOrder(pairs[j].Owner, pairs[j].Incident)
	})
	return occurrences, pairs
}

// occurrenceSequenceTrusted reports whether the recommended sequence rests on
// evidence the reader may act on. Any structural discrepancy, any integrity
// diagnostic over the split board, or any capture source that did not read cleanly
// withdraws that trust: a sequence derived from an unread row file or a broken tree
// is a guess, and the index may never look clean over one.
func occurrenceSequenceTrusted(discrepancies []OccurrenceDiscrepancy, diagnostics []string, sources []SourceFact) bool {
	if len(diagnostics) != 0 {
		return false
	}
	for _, discrepancy := range discrepancies {
		if discrepancy.Structural {
			return false
		}
	}
	states := make(map[string]string, len(sources))
	for _, source := range sources {
		states[source.Source] = source.State
	}
	for _, source := range []string{RoadmapFile, IdeasFile, learnings.JournalPath, retros.Directory + "/"} {
		switch bounds.FileState(states[source]) {
		case bounds.StateAbsent, bounds.StateEmpty, bounds.StateParsed:
			continue
		default:
			return false
		}
	}
	return true
}

func occurrenceOrder(fields ...string) string { return strings.Join(fields, "\x00") }

func finalOccurrence(body string) (owner, incident, kind string, cited bool) {
	final := strings.TrimRightFunc(body, unicode.IsSpace)
	start := strings.LastIndex(final, "[occurrence:")
	if start < 0 {
		return "", "", "", false
	}
	token := final[start:]
	if !strings.HasSuffix(token, "]") {
		if !strings.Contains(token, "]") {
			owner, incident = malformedOccurrenceParts(token)
			return owner, incident, "malformed-token", true
		}
		return "", "", "", false
	}
	matches := occurrenceTokenRe.FindStringSubmatch(token)
	if matches == nil {
		return "", "", "malformed-token", true
	}
	owner, incident = matches[1], matches[2]
	if hasCompleteOccurrenceToken(final[:start]) {
		return owner, incident, "multiple-tokens", true
	}
	if !ValidOccurrenceOwner(owner) || !ValidOccurrenceIncident(incident) {
		return owner, incident, "malformed-token", true
	}
	return owner, incident, "", true
}

func hasCompleteOccurrenceToken(body string) bool {
	for from := 0; from < len(body); {
		start := strings.Index(body[from:], "[occurrence:")
		if start < 0 {
			return false
		}
		start += from
		end := strings.IndexByte(body[start:], ']')
		if end < 0 {
			return false
		}
		end += start + 1
		if match := occurrenceTokenRe.FindStringSubmatch(body[start:end]); match != nil && ValidOccurrenceOwner(match[1]) && ValidOccurrenceIncident(match[2]) {
			return true
		}
		from = end
	}
	return false
}

func malformedOccurrenceParts(token string) (string, string) {
	parts := strings.SplitN(strings.TrimSuffix(strings.TrimPrefix(token, "[occurrence:"), "]"), "/", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}
