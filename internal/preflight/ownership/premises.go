package ownership

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/tickets"
)

// seamCite is one supported `_test.go` (`Test...`) citation in a coverage row's seam cell.
type seamCite struct {
	path, function string
}

// rowSeamCitations extracts each supported Go test seam citation from a coverage row's
// seam cell. Multiple backticked `Test...` names inside one parenthesized group produce
// independent citations so omitting any one of them leaves the seam inventory incomplete.
func rowSeamCitations(cell string) []seamCite {
	var out []seamCite
	rest := cell
	for {
		open := strings.IndexByte(rest, '`')
		if open < 0 {
			break
		}
		rest = rest[open+1:]
		close := strings.IndexByte(rest, '`')
		if close < 0 {
			break
		}
		path := rest[:close]
		rest = rest[close+1:]
		if !strings.HasSuffix(path, "_test.go") {
			continue
		}
		trimmed := strings.TrimLeft(rest, " \t")
		if !strings.HasPrefix(trimmed, "(") {
			continue
		}
		end := strings.IndexByte(trimmed, ')')
		if end < 0 {
			continue
		}
		group := trimmed[1:end]
		rest = trimmed[end+1:]
		var funcs []string
		for {
			fnOpen := strings.IndexByte(group, '`')
			if fnOpen < 0 {
				break
			}
			group = group[fnOpen+1:]
			fnClose := strings.IndexByte(group, '`')
			if fnClose < 0 {
				break
			}
			fn := group[:fnClose]
			group = group[fnClose+1:]
			if strings.HasPrefix(fn, "Test") {
				funcs = append(funcs, fn)
			}
		}
		if len(funcs) > 0 {
			for _, fn := range funcs {
				out = append(out, seamCite{path: path, function: fn})
			}
		}
	}
	return out
}

// checkSeams grades the manifest's seam inventory against the coverage map and the
// committed source tip. Every supported citation of every coverage row of a
// declared-complete manifest must appear in seams, and every existing seam must resolve its
// test function at the source tip.
func (c *census) checkSeams(m Manifest) []Diagnostic {
	if !m.Present || len(m.Diagnostics) > 0 || m.Completeness != completenessDeclared {
		return nil
	}
	var diagnostics []Diagnostic
	for _, row := range m.rows {
		for _, cite := range rowSeamCitations(row.Seam) {
			ok := slices.ContainsFunc(m.Seams, func(s Seam) bool {
				return s.Path == cite.path && s.Function == cite.function && slices.Contains(s.Rows, row.ID)
			})
			if !ok {
				diagnostics = append(diagnostics, Diagnostic{
					ID:     metadataID,
					Reason: fmt.Sprintf("row %s seam omits %s (%s) from seams", row.ID, cite.path, cite.function),
				})
			}
		}
	}
	if c.source.Tip != "" {
		for _, s := range m.Seams {
			if err := c.locateSpec(m, s.Location); err != nil {
				diagnostics = append(diagnostics, Diagnostic{ID: metadataID, Reason: fmt.Sprintf("seam %s: %v", s.ID, err)})
			} else if s.State == "existing" {
				src, err := git.ReadTreeFile(c.root, c.source.Tip, s.Path)
				if err != nil || !declaresTestFunction(s.Path, src, s.Function) {
					diagnostics = append(diagnostics, Diagnostic{ID: metadataID, Reason: fmt.Sprintf("seam %s: %s declares no test function %s", s.ID, s.Path, s.Function)})
				}
			}
		}
	}
	return diagnostics
}

// decisions grades each decision item's exact predicate-to-row mapping, with no duplicate
// row or multi-row predicate mapping.
func (c *census) decisions(m Manifest) []Premise {
	if !m.Present || len(m.Diagnostics) > 0 || len(m.Decisions) == 0 {
		return nil
	}
	rowByID := make(map[string]coverage.MapRow, len(m.rows))
	for _, row := range m.rows {
		rowByID[row.ID] = row
	}
	rowCount := make(map[string]int, len(m.Decisions))
	predRows := make(map[string][]string, len(m.Decisions))
	for _, d := range m.Decisions {
		rowCount[d.Row]++
		if !slices.Contains(predRows[d.Predicate], d.Row) {
			predRows[d.Predicate] = append(predRows[d.Predicate], d.Row)
		}
	}
	rows := make([]Premise, 0, len(m.Decisions))
	for _, d := range m.Decisions {
		reason, err := c.checkDecision(m, d, rowByID, rowCount[d.Row], predRows[d.Predicate])
		row := Premise{ID: d.ID, State: StateGreen, Source: c.source.Tip, Reason: reason}
		if err != nil {
			row.State, row.Reason = StateUnknown, "decision premise: "+err.Error()
		}
		rows = append(rows, row)
	}
	return rows
}

// checkDecision verifies one decision item against its spec location, its pinned source
// line, its coverage row's exact behavior text, and uniqueness across rows and predicates.
func (c *census) checkDecision(m Manifest, d Decision, rows map[string]coverage.MapRow, sameRow int, samePred []string) (string, error) {
	row, ok := rows[d.Row]
	if !ok {
		return "", fmt.Errorf("coverage row %s is not in the coverage map", d.Row)
	}
	var rev string
	if c.source.Tip != "" {
		if err := c.locateSpec(m, d.Location); err != nil {
			return "", err
		}
		var err error
		if rev, err = c.locate(d.Source); err != nil {
			return "", err
		}
	}
	if row.Behavior != d.Predicate {
		return "", fmt.Errorf("coverage row %s behavior is %q, not %q", d.Row, row.Behavior, d.Predicate)
	}
	if rev != "" {
		data, err := git.ReadTreeFile(c.root, rev, d.Source.Path)
		if err != nil {
			return "", err
		}
		lines := strings.Split(string(data), "\n")
		if d.Source.Line > len(lines) || !strings.Contains(strings.TrimSuffix(lines[d.Source.Line-1], "\r"), d.Predicate) {
			return "", fmt.Errorf("%s line %d does not hold predicate %q", d.Source.Path, d.Source.Line, d.Predicate)
		}
	}
	if sameRow > 1 {
		return "", fmt.Errorf("multiple decisions map coverage row %s", d.Row)
	}
	if len(samePred) > 1 {
		return "", fmt.Errorf("predicate %q maps to multiple coverage rows (%s)", d.Predicate, strings.Join(samePred, ", "))
	}
	return "decision " + d.ID + " maps predicate to coverage row " + d.Row, nil
}

// manualClaims projects each manual claim onto a manual premise row and validates any
// native-future witness against the coverage map, tickets, completion plan, and phase owner.
func (c *census) manualClaims(m Manifest, parsed []tickets.Ticket) ([]Premise, []ManualClaim, bool) {
	var rows []Premise
	var pending []ManualClaim
	valid := true
	for _, claim := range m.ManualClaims {
		reason := string(claim.Applicability) + " manual claim: " + claim.Reason
		if claim.Applicability == NativeFuture {
			if err := c.validateNativeFuture(m, claim, parsed); err != nil {
				valid = false
				reason += ": " + err.Error()
			} else {
				pending = append(pending, claim)
			}
		}
		rows = append(rows, Premise{ID: claim.ID, State: StateManual, Source: c.source.Tip, Reason: reason})
	}
	return rows, pending, valid
}

// validateNativeFuture checks that a native-future manual claim sits at its own coverage
// row's spec line, as the coverage parser reports it, and names an exclusively owned manual row, a valid completion-plan chunk
// and verification requirement, and the native implementation procedure as its phase owner,
// located at its pinned source.
func (c *census) validateNativeFuture(m Manifest, claim ManualClaim, parsed []tickets.Ticket) error {
	if c.source.Tip == "" || claim.NativeWitness == nil {
		return fmt.Errorf("native-future witness has no committed source tip")
	}
	w := claim.NativeWitness
	idx := slices.IndexFunc(m.rows, func(r coverage.MapRow) bool { return r.ID == w.Row })
	if idx < 0 {
		return fmt.Errorf("row %s is not in the coverage map", w.Row)
	}
	if err := c.locateSpec(m, claim.Location); err != nil {
		return err
	}
	if claim.Location.Line != m.rows[idx].Line {
		return fmt.Errorf("%s line %d is not the location of row %s", claim.Location.Path, claim.Location.Line, w.Row)
	}
	if w.PhaseOwner.Path != chargesource.BuildPhase {
		return fmt.Errorf("phase owner %s is not the native implementation procedure %s", w.PhaseOwner.Path, chargesource.BuildPhase)
	}
	if _, err := c.locate(w.PhaseOwner); err != nil {
		return err
	}
	if len(rowSeamCitations(m.rows[idx].Seam)) > 0 || slices.ContainsFunc(m.Seams, func(s Seam) bool { return slices.Contains(s.Rows, w.Row) }) {
		return fmt.Errorf("row %s declares a supported test seam", w.Row)
	}
	var owners []string
	for _, t := range parsed {
		if slices.Contains(t.Covers, w.Row) {
			owners = append(owners, t.Name)
		}
	}
	if len(owners) != 1 || owners[0] != w.Ticket {
		return fmt.Errorf("row %s must be owned by ticket %s alone", w.Row, w.Ticket)
	}
	plan, err := reviewrecord.ReadPlan(c.root, c.source.Tip, c.source.SpecPath)
	if err != nil {
		return err
	}
	cIdx := slices.IndexFunc(plan.Chunks, func(ch reviewrecord.PlannedChunk) bool {
		return ch.ID == w.Chunk && slices.Contains(ch.Tickets, w.Ticket)
	})
	if cIdx < 0 {
		return fmt.Errorf("chunk %s does not own ticket %s", w.Chunk, w.Ticket)
	}
	vIdx := slices.IndexFunc(plan.Chunks[cIdx].Verification, func(v reviewrecord.Requirement) bool {
		return v.ID == w.Verification && (plan.Version < 2 || v.Ticket == w.Ticket)
	})
	if vIdx < 0 {
		return fmt.Errorf("chunk %s declares no verification %s for %s", w.Chunk, w.Verification, w.Ticket)
	}
	return nil
}

// locateSpec verifies that a spec location's line holds its needle exactly once in the
// parsed spec bytes.
func (c *census) locateSpec(m Manifest, loc Location) error {
	if c.source.SpecPath != "" && loc.Path != c.source.SpecPath {
		return fmt.Errorf("spec location path %s is not %s", loc.Path, c.source.SpecPath)
	}
	lines := strings.Split(string(m.spec), "\n")
	if loc.Line < 1 || loc.Line > len(lines) || strings.Count(strings.TrimSuffix(lines[loc.Line-1], "\r"), loc.Needle) != 1 {
		return fmt.Errorf("%s line %d does not hold its needle %q once", loc.Path, loc.Line, loc.Needle)
	}
	return nil
}

// chargeReadiness computes the fail-closed charge projection from the graded premises,
// collection diagnostics, and validated manual claims. Its blocker names the first unmet
// condition: an inventory fault first, then an unproved required current premise.
func chargeReadiness(m Manifest, premises []Premise, diagnostics []Diagnostic, pending []ManualClaim, nativeValid bool) ChargeReadiness {
	ready := ChargeReadiness{
		InventoryComplete: m.Present && m.Completeness == completenessDeclared && len(diagnostics) == 0 && nativeValid,
		PendingNative:     pending,
	}
	byID := make(map[string]Premise, len(premises))
	for _, p := range premises {
		byID[p.ID] = p
	}
	block := func(id, reason string) {
		if ready.Blocker == "" {
			ready.Blocker = id + ": " + reason
		}
	}
	switch {
	case len(diagnostics) > 0:
		block(diagnostics[0].ID, diagnostics[0].Reason)
	case m.Completeness != completenessDeclared:
		block(metadataID, "the premise inventory is "+m.Completeness+", not "+completenessDeclared)
	}
	for _, claim := range m.ManualClaims {
		if claim.Applicability == NativeFuture && !slices.ContainsFunc(pending, func(p ManualClaim) bool { return p.ID == claim.ID }) {
			block(claim.ID, byID[claim.ID].Reason)
		}
	}
	var requiredCurrent []string
	for _, item := range m.Premises {
		if item.State == "current" {
			requiredCurrent = append(requiredCurrent, item.ID)
		}
	}
	for _, d := range m.Decisions {
		requiredCurrent = append(requiredCurrent, d.ID)
	}
	for _, claim := range m.ManualClaims {
		if claim.Applicability == CurrentPremise {
			requiredCurrent = append(requiredCurrent, claim.ID)
		}
	}
	for _, id := range requiredCurrent {
		if byID[id].State != StateGreen {
			block(id, byID[id].Reason)
		}
	}
	return ready
}
