package conformance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/roadmap/roadmaptest"
)

func checkOccurrenceLedgerAndMaintenance(root string) []string {
	return append(checkOccurrenceLedgerMigration(root), checkRecurrenceMaintenanceContract(root)...)
}

func checkRecurrenceMaintenanceContract(root string) []string {
	path := filepath.Join(root, ".agents", "commands", "bench-drain.md")
	text := stripHTMLComments(readIfExists(path))
	if strings.TrimSpace(text) == "" {
		return []string{"bench-drain recurrence maintenance command is unavailable"}
	}

	entry, entryCount := markdownH2Sections(text, "Entry orientation")
	delegation, delegationCount := markdownH2Sections(text, "Delegate the evidence")
	occurrences, occurrenceCount := markdownH2Sections(text, "2. Drain occurrence evidence")
	sequence, sequenceCount := markdownH2Sections(text, "7. Refresh the sequence")
	batch, batchCount := markdownH2Sections(text, "8. Batch-propose, then commit once on green")
	entry = collapseSpace(entry)
	delegation = collapseSpace(delegation)
	occurrences = collapseSpace(occurrences)
	sequence = collapseSpace(sequence)
	batch = collapseSpace(batch)
	var diags []string
	if entryCount != 1 || !strings.Contains(entry, "`context.schema = 4`") ||
		!strings.Contains(entry, "every other schema stops the phase before any batch mutation") ||
		!strings.Contains(entry, "Do not guess recurrence facts from an older schema") {
		diags = append(diags, "bench-drain does not require schema 4 before recurrence maintenance")
	}
	if entryCount != 1 || !strings.Contains(entry, "schema-4 index is the complete local inventory") ||
		!strings.Contains(entry, "The index proves what exists; targeted fetches and named body reads prove content") ||
		!strings.Contains(entry, "Fetch complete roadmap detail only for rows the reconcile touches") ||
		!strings.Contains(entry, "Read idea, learning, and retro bodies from the paths the index names") {
		diags = append(diags, "bench-drain does not use index-first evidence for recurrence maintenance")
	}
	schemaAt := strings.Index(entry, "`context.schema = 4`")
	trustAt := strings.Index(entry, "Before beginning `## 1. Reconcile first`, read `context.sequence_trusted` from that snapshot.")
	trustStop := "snapshot. When it is false, stop before any batch mutation"
	if entryCount != 1 || trustAt <= schemaAt || !strings.Contains(entry, trustStop) {
		diags = append(diags, "bench-drain does not stop untrusted recurrence evidence before reconciliation")
	}
	if entryCount != 1 || trustAt <= schemaAt || !strings.Contains(entry, "`occurrence_discrepancies` row together with the complete index snapshot") {
		diags = append(diags, "bench-drain does not preserve complete index evidence before reconciliation")
	}
	parallelRead := "Use at most three read-only delegates after the snapshot passes its trust check."
	readScopes := "Assign at most one delegate to each scope: roadmap reconcile, idea and journal mapping, and retro analysis."
	skipEmpty := "Skip a scope when its indexed source set is empty."
	delegateRouting := "Charge every delegate under `craft-delegate` and `craft-line`."
	readFence := "Read delegates edit nothing, take no new inventory, and use only the snapshot and its named paths."
	returnShape := "Each read delegate returns these fields: proposed owner, classification, occurrence, evidence, and reviewer decision."
	if delegationCount != 1 || !strings.Contains(delegation, parallelRead) {
		diags = append(diags, "bench-drain does not retain the bounded parallel-read contract")
	}
	if delegationCount != 1 || !strings.Contains(delegation, readScopes) {
		diags = append(diags, "bench-drain does not assign the three independent read scopes")
	}
	if delegationCount != 1 || !strings.Contains(delegation, skipEmpty) {
		diags = append(diags, "bench-drain does not skip empty indexed read scopes")
	}
	if delegationCount != 1 || !strings.Contains(delegation, delegateRouting) {
		diags = append(diags, "bench-drain does not route every delegate through craft-delegate and craft-line")
	}
	if delegationCount != 1 || !strings.Contains(delegation, readFence) {
		diags = append(diags, "bench-drain does not keep read delegates inside the snapshot-only fence")
	}
	if delegationCount != 1 || !strings.Contains(delegation, returnShape) {
		diags = append(diags, "bench-drain does not require the fixed read-delegate return shape")
	}

	decisionsBeforeWriter := "Resolve duplicate incidents and reviewer decisions before any batch writer starts."
	treeVerification := "Verify that the tree stayed unchanged."
	coordinatorOwnership := "Keep ignored capture removal, the handoff, verification, and landing with the coordinator."
	implementOverlap := "An implement-now writer may run while other reads continue."
	implementRouting := "Route each implement-now writer through `craft-delegate` isolation and `craft-line` routing."
	withImplementTiming := "If an implement-now item exists, create the batch worktree only after every such item lands green on `main`."
	withoutImplementTiming := "If no implement-now item exists, create the batch worktree after all reads finish and the coordinator resolves duplicate incidents and reviewer decisions."
	singleWriter := "If tracked changes remain, one later write delegate authors the complete tracked batch."
	decisionAt := strings.Index(delegation, decisionsBeforeWriter)
	writerAt := strings.Index(delegation, singleWriter)
	if delegationCount != 1 || decisionAt < 0 || writerAt <= decisionAt {
		diags = append(diags, "bench-drain does not resolve cross-source decisions before batch writing")
	}
	if delegationCount != 1 || !strings.Contains(delegation, treeVerification) {
		diags = append(diags, "bench-drain does not verify the tree stayed unchanged after reading")
	}
	if delegationCount != 1 || !strings.Contains(delegation, coordinatorOwnership) {
		diags = append(diags, "bench-drain does not retain coordinator ownership of local and landing work")
	}
	if delegationCount != 1 || !strings.Contains(delegation, implementOverlap) {
		diags = append(diags, "bench-drain does not allow implement-now work to overlap remaining reads")
	}
	if delegationCount != 1 || !strings.Contains(delegation, implementRouting) {
		diags = append(diags, "bench-drain does not route implement-now writers through isolation and line routing")
	}
	if delegationCount != 1 || !strings.Contains(delegation, withImplementTiming) {
		diags = append(diags, "bench-drain does not wait for every implement-now landing before batch creation")
	}
	if delegationCount != 1 || !strings.Contains(delegation, withoutImplementTiming) {
		diags = append(diags, "bench-drain does not create the batch after reads when no implement-now item exists")
	}
	if delegationCount != 1 || !strings.Contains(delegation, singleWriter) {
		diags = append(diags, "bench-drain does not retain one conditional tracked batch writer")
	}

	noTrackedWriter := "If no tracked changes remain, start no batch writer."
	ignoredAfterApproval := "After approval, the coordinator empties ignored inbox and journal sources."
	ignoredHandoff := "It writes ignored `capture/session-handoff.md` last."
	reviewerBatch := "The reviewer batch contains the tracked diff, proposed ignored-source removals, and every journal verdict."
	trackedDiff := "The tracked diff contains roadmap dispositions, retro removals, earned `bench spec retire` work, and provider scorecards."
	ignoredNotTracked := "Ignored local changes do not enter that diff or commit."
	handoffWriteTime := "`bench status` dates the ignored handoff by its write time."
	all := collapseSpace(text)
	if batchCount != 1 || !strings.Contains(batch, noTrackedWriter) {
		diags = append(diags, "bench-drain does not suppress the batch writer when no tracked changes remain")
	}
	if !strings.Contains(all, ignoredAfterApproval) {
		diags = append(diags, "bench-drain does not delay ignored-source removal until approval")
	}
	if !strings.Contains(all, ignoredHandoff) {
		diags = append(diags, "bench-drain does not keep the ignored handoff last")
	}
	if batchCount != 1 || !strings.Contains(batch, reviewerBatch) {
		diags = append(diags, "bench-drain does not preserve the tracked-versus-ignored reviewer batch boundary")
	}
	if batchCount != 1 || !strings.Contains(batch, trackedDiff) {
		diags = append(diags, "bench-drain does not constrain the tracked diff contents")
	}
	if batchCount != 1 || !strings.Contains(batch, ignoredNotTracked) {
		diags = append(diags, "bench-drain does not exclude ignored local changes from the tracked diff and commit")
	}
	if !strings.Contains(all, handoffWriteTime) {
		diags = append(diags, "bench-drain does not date the ignored handoff by write time")
	}
	pending := "add its incident key to that owner's `Occurrences:` line in `ROADMAP.md` before removing any source unit"
	if occurrenceCount != 1 || !strings.Contains(occurrences, "For every `pending` owner/incident pair") || !strings.Contains(occurrences, pending) {
		diags = append(diags, "bench-drain does not ledger every pending occurrence before source removal")
	}
	if occurrenceCount != 1 || !strings.Contains(occurrences, "Every `already-recorded` source already has that key: remove its source unit without adding another key") {
		diags = append(diags, "bench-drain does not remove already-recorded occurrence sources without another key")
	}

	normalizeHeading := "### Normalize touched rows"
	normalizeBeforeBatch := "Normalize every touched row before batch proposal."
	if occurrenceCount != 1 || strings.Count(occurrences, normalizeHeading) != 1 || !strings.Contains(occurrences, normalizeBeforeBatch) {
		diags = append(diags, "bench-drain does not normalize every touched row before batch proposal")
	}
	physicalOccurrence := "exactly one physical `Occurrence: <when/source> — <short situation>.` line"
	if occurrenceCount != 1 || !strings.Contains(occurrences, physicalOccurrence) {
		diags = append(diags, "bench-drain does not require one physical occurrence line per drained event")
	}
	eventOnly := "Occurrence lines contain event evidence only."
	coreProse := "New feature faces and decisions remain concise core prose."
	if occurrenceCount != 1 || !strings.Contains(occurrences, eventOnly) || !strings.Contains(occurrences, coreProse) {
		diags = append(diags, "bench-drain does not keep occurrence evidence separate from core remedies")
	}
	goodOccurrence := "**Good — event-only evidence:** `Occurrence: 2026-08-15 gate build — primary-checkout preflight failed on a stale base.`"
	badOccurrence := "**Bad — remedy derivation:** `Occurrence: 2026-08-15 gate build — change preflight selection to fix the stale base.`"
	goodAt := strings.Index(occurrences, goodOccurrence)
	badAt := strings.Index(occurrences, badOccurrence)
	if occurrenceCount != 1 || strings.Count(occurrences, goodOccurrence) != 1 || strings.Count(occurrences, badOccurrence) != 1 || goodAt < 0 || badAt <= goodAt {
		diags = append(diags, "bench-drain does not retain the event-only occurrence contrast")
	}

	precedence := []string{
		"Rank rows by severity.",
		"choose actionable work over blocked work",
		"apply literal dependencies, then explicit reviewer pricing.",
		"rank by descending occurrence count.",
		"existing reproduced defect-over-feature rule, then cheapest-first cost rule.",
	}
	previous := -1
	if sequenceCount != 1 {
		diags = append(diags, "bench-drain recurrence sequence precedence is unavailable")
	} else {
		for _, anchor := range precedence {
			at := strings.Index(sequence, anchor)
			if at < 0 || at <= previous {
				diags = append(diags, "bench-drain recurrence sequence precedence is incomplete or out of order")
				break
			}
			previous = at
		}
	}
	return diags
}

func TestRecurrenceMaintenanceContractCheckBites(t *testing.T) {
	h := NewHarness(t)
	command := readIfExists(filepath.Join(h.KitRoot, ".agents", "commands", "bench-drain.md"))
	if command == "" {
		t.Fatal("bench-drain command unavailable")
	}

	mutations := []struct {
		name, old, new, want string
	}{
		{"schema", "`context.schema = 4`", "`context.schema = 3`", "bench-drain does not require schema 4 before recurrence maintenance"},
		{"schema guessing", "Do not guess recurrence facts from an older schema.", "Infer recurrence facts from an older schema.", "bench-drain does not require schema 4 before recurrence maintenance"},
		{"index inventory", "schema-4\nindex is the complete local inventory", "schema-4\nindex is a partial inventory", "bench-drain does not use index-first evidence for recurrence maintenance"},
		{"targeted detail", "Fetch complete roadmap detail only for rows the reconcile touches", "Fetch complete roadmap detail for every row", "bench-drain does not use index-first evidence for recurrence maintenance"},
		{"named capture bodies", "Read idea, learning, and retro bodies\nfrom the paths the index names", "Read capture bodies from a directory sweep", "bench-drain does not use index-first evidence for recurrence maintenance"},
		{"trust after reconcile", "Before beginning `## 1. Reconcile first`, read `context.sequence_trusted` from that\nsnapshot.", "After `## 1. Reconcile first`, read `context.sequence_trusted` from that\nsnapshot.", "bench-drain does not stop untrusted recurrence evidence before reconciliation"},
		{"trust polarity", "snapshot. When it is false, stop before any batch mutation", "snapshot. When it is true, stop before any batch mutation", "bench-drain does not stop untrusted recurrence evidence before reconciliation"},
		{"trust stop", "stop before any batch mutation", "continue into batch mutation", "bench-drain does not stop untrusted recurrence evidence before reconciliation"},
		{"context evidence", "`occurrence_discrepancies` row together with the complete index snapshot", "`occurrence_discrepancies` row", "bench-drain does not preserve complete index evidence before reconciliation"},
		{"bounded parallel reads", "Use at most three read-only delegates after the snapshot passes its trust check.", "Use unlimited read-only delegates after the snapshot passes its trust check.", "bench-drain does not retain the bounded parallel-read contract"},
		{"read scopes", "Assign at most one delegate to each scope: roadmap reconcile, idea and journal mapping, and retro analysis.", "Assign one delegate to roadmap reconcile only.", "bench-drain does not assign the three independent read scopes"},
		{"skip empty scope", "Skip a scope when its indexed source set is empty.", "Run every scope even when its indexed source set is empty.", "bench-drain does not skip empty indexed read scopes"},
		{"delegate routing", "Charge every delegate under\n`craft-delegate` and `craft-line`.", "Charge every delegate without routing.", "bench-drain does not route every delegate through craft-delegate and craft-line"},
		{"read fence", "Read delegates edit nothing, take no new inventory, and use only the snapshot and its named paths.", "Read delegates may edit notes and take a new inventory.", "bench-drain does not keep read delegates inside the snapshot-only fence"},
		{"return shape", "Each read delegate returns these fields: proposed owner, classification, occurrence, evidence, and reviewer decision.", "Each read delegate returns a summary.", "bench-drain does not require the fixed read-delegate return shape"},
		{"decisions before writer", "Resolve duplicate incidents and reviewer decisions before any batch writer starts.", "Resolve duplicate incidents after the batch writer starts.", "bench-drain does not resolve cross-source decisions before batch writing"},
		{"tree verification", "Verify that the tree stayed unchanged.", "Assume that the tree stayed unchanged.", "bench-drain does not verify the tree stayed unchanged after reading"},
		{"coordinator ownership", "Keep ignored capture removal, the handoff,\nverification, and landing with the coordinator.", "Delegate ignored capture removal and landing.", "bench-drain does not retain coordinator ownership of local and landing work"},
		{"implement-now overlap", "An implement-now writer may run\nwhile other reads continue.", "Start implement-now work after every read finishes.", "bench-drain does not allow implement-now work to overlap remaining reads"},
		{"implement-now routing", "Route each implement-now writer through\n`craft-delegate` isolation and `craft-line` routing.", "Route implement-now writers without isolation.", "bench-drain does not route implement-now writers through isolation and line routing"},
		{"implement-now landing timing", "If an implement-now item exists, create the batch worktree only after every such item lands green on `main`.", "Create the batch worktree before implement-now items land.", "bench-drain does not wait for every implement-now landing before batch creation"},
		{"no-implement-now timing", "If no implement-now item exists, create the batch worktree after all reads finish and the coordinator resolves duplicate incidents and reviewer decisions.", "If no implement-now item exists, never create the batch worktree.", "bench-drain does not create the batch after reads when no implement-now item exists"},
		{"single batch writer", "If tracked changes remain, one later write delegate authors the complete tracked batch.", "Several later write delegates author the complete tracked batch.", "bench-drain does not retain one conditional tracked batch writer"},
		{"no tracked writer", "If no tracked changes\nremain, start no batch writer.", "If no tracked changes remain, start a batch writer.", "bench-drain does not suppress the batch writer when no tracked changes remain"},
		{"ignored removal after approval", "After approval, the coordinator empties ignored inbox and journal sources.", "Before approval, remove ignored sources.", "bench-drain does not delay ignored-source removal until approval"},
		{"ignored handoff last", "It\nwrites ignored `capture/session-handoff.md` last.", "It writes the ignored handoff first.", "bench-drain does not keep the ignored handoff last"},
		{"reviewer batch boundary", "The reviewer batch contains the tracked diff, proposed ignored-source removals,\nand every journal verdict.", "The tracked diff contains every journal verdict.", "bench-drain does not preserve the tracked-versus-ignored reviewer batch boundary"},
		{"tracked diff contents", "The tracked diff contains roadmap dispositions, retro removals, earned `bench spec retire` work, and provider scorecards.", "The tracked diff contains journal verdicts.", "bench-drain does not constrain the tracked diff contents"},
		{"ignored changes excluded", "Ignored local\nchanges do not enter that diff or commit.", "Ignored local changes enter the tracked diff.", "bench-drain does not exclude ignored local changes from the tracked diff and commit"},
		{"ignored handoff write time", "`bench\nstatus` dates the ignored handoff by its write time.", "`bench status` dates the ignored handoff by its commit time.", "bench-drain does not date the ignored handoff by write time"},
		{"pending ledger", "add its incident\nkey to that owner's `Occurrences:` line in `ROADMAP.md` before removing any source\nunit", "record the pending pair for later", "bench-drain does not ledger every pending occurrence before source removal"},
		{"already recorded", "remove its\nsource unit without adding another key", "remove its source unit after adding another key", "bench-drain does not remove already-recorded occurrence sources without another key"},
		{"normalize heading", "### Normalize touched rows", "### Touched rows", "bench-drain does not normalize every touched row before batch proposal"},
		{"normalize before proposal", "Normalize every touched row before batch proposal.", "Normalize touched rows after batch proposal.", "bench-drain does not normalize every touched row before batch proposal"},
		{"physical occurrence shape", "exactly one physical\n`Occurrence: <when/source> — <short situation>.` line", "one or more `Occurrence:` paragraphs", "bench-drain does not require one physical occurrence line per drained event"},
		{"event-only boundary", "Occurrence lines contain\nevent evidence only.", "Occurrence lines contain event evidence and remedy derivation.", "bench-drain does not keep occurrence evidence separate from core remedies"},
		{"contrastive example", "**Good — event-only evidence:** `Occurrence: 2026-08-15 gate build — primary-checkout preflight failed on a stale base.`", "**Good — remedy derivation:** `Occurrence: 2026-08-15 gate build — change preflight selection to fix the stale base.`", "bench-drain does not retain the event-only occurrence contrast"},
		{"equal class", "apply literal dependencies, then explicit\nreviewer pricing. Only when all four stronger inputs tie, rank by descending\noccurrence count.", "rank by descending occurrence count before explicit reviewer pricing.", "bench-drain recurrence sequence precedence is incomplete or out of order"},
	}
	for _, mutation := range mutations {
		t.Run(mutation.name, func(t *testing.T) {
			if baseline := checkRecurrenceMaintenanceContract(h.KitRoot); len(baseline) != 0 {
				t.Fatalf("live command has baseline diagnostics before mutation: %s", strings.Join(baseline, "; "))
			}
			if strings.Count(command, mutation.old) != 1 {
				t.Fatalf("mutation anchor %q occurs %d times", mutation.old, strings.Count(command, mutation.old))
			}
			root := t.TempDir()
			path := filepath.Join(root, ".agents", "commands", "bench-drain.md")
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(strings.Replace(command, mutation.old, mutation.new, 1)), 0o644); err != nil {
				t.Fatal(err)
			}
			if !containsDiagnostic(checkRecurrenceMaintenanceContract(root), mutation.want) {
				t.Fatalf("mutation did not bite: want %q", mutation.want)
			}
			t.Logf("observed red: %s", mutation.want)
		})
	}
}

func TestOccurrenceLedgerMigrationCheckBitesOnFT158Count(t *testing.T) {
	root := t.TempDir()
	kit := NewHarness(t).KitRoot
	index, err := os.ReadFile(filepath.Join(kit, "ROADMAP.md"))
	if err != nil {
		t.Fatal(err)
	}
	// The mutation drops a key from the ledger, which lives in the row's own detail file.
	// The fixture copies the live board's index beside a mutated roadmap/FT158.md.
	rowFile := filepath.Join(kit, "roadmap", "FT158.md")
	row, err := os.ReadFile(rowFile)
	if err != nil {
		t.Fatal(err)
	}
	mutated := strings.Replace(string(row), "Occurrences: baseline-01, baseline-02, baseline-03", "Occurrences: baseline-01, baseline-02", 1)
	if mutated == string(row) {
		t.Fatal("FT158 migration-count mutation did not change its ledger")
	}
	roadmaptest.WriteSplitBoard(t, root, string(index), map[string]string{"FT158.md": mutated})
	found := false
	for _, diag := range checkOccurrenceLedgerMigration(root) {
		if diag == "ROADMAP.md occurrence-ledger migration count for FT158 is wrong" {
			found = true
		}
	}
	if !found {
		t.Fatal("FT158 count mutation passed occurrence-ledger migration check")
	}
}

func TestOccurrenceLedgerMigrationAllowsRetiredFT126(t *testing.T) {
	root := t.TempDir()
	index, files := "", map[string]string{}
	for _, row := range [][2]string{
		{"FT71", "baseline-01"},
		{"FT158", "baseline-01, baseline-02, baseline-03"},
		{"FT128", "baseline-01"},
		{"FT98", "baseline-01, baseline-02, baseline-03"},
		{"FT169", "baseline-01"},
		{"FT141", "baseline-01"},
		{"FT94", "baseline-01"},
		{"FT125", "baseline-01"},
	} {
		heading := "**" + row[0] + " — migration baseline.**"
		index += heading + "\n\n"
		files[row[0]+".md"] = heading + "\nOccurrences: " + row[1] + "\n"
	}
	roadmaptest.WriteSplitBoard(t, root, index, files)
	if containsDiagnostic(checkOccurrenceLedgerMigration(root), "ROADMAP.md occurrence-ledger migration count for FT126 is wrong") {
		t.Fatal("retired FT126 remained required by occurrence-ledger migration")
	}
}
