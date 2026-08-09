package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
)

func TestLoadValidityMetadataFixturesBite(t *testing.T) {
	h := NewHarness(t)
	kitRoot := h.KitRoot
	fixtures := []string{
		"invalid-json",
		"codex-hooks-broken",
		"codex-hooks-timeout",
		"codex-hooks-timeout-typed",
		"bad-frontmatter",
		"claude-skills-unmirrored",
		"extensionless-gate-ref",
		"shared-rule-drift",
		"readme-shared-rule-drift",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			runFixtureBite(t, kitRoot, fixture)
		})
	}
}

func TestSkillsIndexAndCommandAdapterFixturesBite(t *testing.T) {
	h := NewHarness(t)
	kitRoot := h.KitRoot
	fixtures := []string{
		"dangling-index",
		"missing-index-field",
		"stale-index-wording",
		"unindexed-skill",
		"roadmap-promotion-persistence",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			runFixtureBite(t, kitRoot, fixture)
		})
	}
}

func TestDocsCurrencyTokenDietAndWorkflowFixturesBite(t *testing.T) {
	h := NewHarness(t)
	kitRoot := h.KitRoot
	all, err := canary.Fixtures(filepath.Join(kitRoot, "tests", "canary"))
	requireFixtureNoError(t, err)
	var fixtures []string
	for name, fixture := range all {
		if fixture.Family == "docs-currency-token-diet" || fixture.Family == "workflow-guidance-anchors" {
			fixtures = append(fixtures, name)
		}
	}
	slices.Sort(fixtures)
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			runFixtureBite(t, kitRoot, fixture)
		})
	}
}

func TestSpecBuildCadenceAnchorsRejectDeletionSwapAndRawGitRouting(t *testing.T) {
	const (
		bootstrapDeletionDiag = ".agents/skills/bench-craft-spec/SKILL.md dropped the bootstrap-authority pre-execution trace"
		bootstrapAfterDiag    = ".agents/skills/bench-craft-spec/SKILL.md validates a bootstrap authority after launch"
		bootstrapPointerDiag  = "bench-write-spec.md does not apply craft-spec's named bootstrap-authority rule during edge walking and falsification"
		repairEnvelopeDiag    = ".agents/skills/bench-craft-tickets/SKILL.md dropped the debug receipt's maximum repair envelope"
		repairResultDiag      = ".agents/skills/bench-craft-tickets/SKILL.md dropped the one-ticket-or-reciprocal-chain repair result"
		repairChainOnlyDiag   = ".agents/skills/bench-craft-tickets/SKILL.md contains an additive chain-only repair mandate for validated debug receipts"
		repairUnionDiag       = ".agents/skills/bench-craft-tickets/SKILL.md dropped repair-ticket fence union containment"
		repairEscapeDiag      = ".agents/skills/bench-craft-tickets/SKILL.md permits a repair ticket to escape the debug receipt's required fence"
		repairOwnerDiag       = "bench-implement-spec dropped the craft-tickets repair-reslicing owner pointer"
		repairCommonDiag      = "bench-implement-spec dropped the one-repair-ticket common case"
		repairSingularDiag    = "bench-implement-spec restores the singular exactly-one repair-ticket mandate"
		repairAssignDiag      = "bench-implement-spec dropped ordinary assign for every repair-chain ticket"
		repairCheckDiag       = "bench-implement-spec dropped ordinary checkpoint for every repair-chain ticket"
		repairIntegrateDiag   = "bench-implement-spec dropped ordinary integrate for every repair-chain ticket"
		repairGitDiag         = "bench-implement-spec permits chain-local synthesized Git checkpoint plumbing"
		repairTerminalDiag    = "bench-implement-spec dropped the terminal repair-ticket refresh precondition"
		repairEarlyDiag       = "bench-implement-spec permits premature refresh after a non-terminal repair ticket"
		repairAssignmentDiag  = "bench-implement-spec replaced the original blocked assignment refresh identity"
		repairReceiptDiag     = "bench-implement-spec replaced the original validated debug receipt identity"
	)
	h := NewHarness(t)
	owner, ok := conformanceChecks["docs-currency-workflow"]
	if !ok {
		t.Fatal("docs-currency-workflow conformance owner is not bound")
	}
	root := h.KitRoot
	diags := owner.run(root, h.KitRoot, registry.Dev)
	for _, diag := range []string{
		bootstrapDeletionDiag, bootstrapAfterDiag, bootstrapPointerDiag,
		repairEnvelopeDiag, repairResultDiag, repairChainOnlyDiag,
		repairUnionDiag, repairEscapeDiag,
		repairOwnerDiag,
		repairCommonDiag, repairSingularDiag, repairAssignDiag, repairCheckDiag,
		repairIntegrateDiag, repairGitDiag,
		repairTerminalDiag, repairEarlyDiag, repairAssignmentDiag, repairReceiptDiag,
	} {
		if containsDiagnostic(diags, diag) {
			t.Fatalf("finished workflow guidance is not conformant with %q:\n%s", diag, strings.Join(diags, "\n"))
		}
	}
	tests := []struct{ name, rel, old, replacement, diag string }{
		{"lifecycle deletion", ".agents/commands/bench-implement-spec.md", "`start` → `assign` → `checkpoint` →\n`integrate` → `review` → `promote`; `status` inspects the run and `abandon`\nplans or applies cleanup.", "", "bench-implement-spec dropped or reordered the eight-operation spec-build lifecycle"},
		{"initial capacity deletion", ".agents/commands/bench-implement-spec.md", "Re-derive the complete ready frontier and the harness's live capacity before\ndispatch. Assign every ownership-safe ticket up to the smaller of frontier size\nand available capacity.", "", "bench-implement-spec dropped initial frontier capacity dispatch"},
		{"additive generic unused-slot reason", ".agents/commands/bench-implement-spec.md", "For every unused harness slot, record exactly one\nreason: dependency, overlapping ownership fence, unavailable harness capacity,\nor measured resource constraint.", "For every unused harness slot, record exactly one reason: dependency, overlapping ownership fence, unavailable harness capacity, or measured resource constraint. An unused slot may instead be `NOT\n  PARALLELIZABLE`.", "bench-implement-spec permits a generic unused-slot reason outside the closed set"},
		{"exact candidate input deletion", ".agents/commands/bench-review-implementation.md", "For an active spec build, read `bench spec build status\n   <slug> --full` and bind the review inputs to the exact candidate subject and\n   recorded run base it reports. Confirm that subject is unchanged immediately\n   before receipt submission; a changed candidate invalidates the review rather\n   than letting a delta review authorize a new composition. ", "", "bench-review-implementation dropped exact-candidate review input binding"},
		{"frontier swap", ".agents/commands/bench-implement-spec.md", "or measured resource constraint. Refill the ownership-safe frontier after every\nintegration or assignment release while another delegate remains active.", "or measured resource constraint. Wait for the ownership-safe frontier to drain before refill after every integration or assignment release.", "bench-implement-spec replaced continuous frontier refill with drain-then-refill cadence"},
		{"repair deletion", ".agents/commands/bench-implement-spec.md", "Accepted findings become new ownership-fenced repair tickets and re-enter\n  `assign`, `checkpoint`, and `integrate` before a fresh composed review.", "", "bench-implement-spec routes an accepted repair outside the provisional lifecycle"},
		{"recomposition purpose deletion", ".bench/BENCH-reference.md", "; a moved tip recomposes through `promote`, discarding the review |", " |", "BENCH-reference dropped promote's moved-tip recomposition from the lifecycle lookup"},
		{"recomposition discard deletion", ".agents/commands/bench-implement-spec.md", "When the branch tip moves, `promote` is the operation that recomposes the run\n  onto the new tip, and recomposition discards the review.\n  ", "", "bench-implement-spec dropped moved-tip recomposition through promote or its review discard"},
		{"repair round deletion", ".agents/commands/bench-implement-spec.md", "\n  The repair round is therefore\n  repair → `promote` → `review` → `assign` … `integrate` → `review` → `promote`.", "", "bench-implement-spec dropped the ordered moved-tip repair round"},
		{"probe kind", ".agents/skills/bench-craft-delegate/SKILL.md", "The\ncoordinator probe's mutation kind differs from the delegate author's mutation\nkind.", "The\ncoordinator probe's mutation kind matches the delegate author's mutation\nkind.", "craft-delegate allows the coordinator probe to repeat the author's mutation kind"},
		{"ordinary commit route", ".bench/BENCH.md", "Provisional\ncadence is exclusive to reviewed spec-backed builds; light-path work, `bench\nshift`, and ordinary `bench commit` remain commit-on-green.", "Provisional cadence covers reviewed spec-backed builds, light-path work, `bench shift`, and ordinary `bench commit`.", ".bench/BENCH.md broadened provisional cadence beyond reviewed spec-backed builds"},
		{"purpose swap", ".bench/BENCH-reference.md", "| `assign` | lease one ownership-fenced ticket worktree; `--refresh <receipt>` re-bases a blocked assignment onto the repaired candidate on a validated debug receipt |", "| `assign` | validate focused evidence and bind a provisional commit |", "BENCH-reference misroutes spec build assign"},
		{"flag positional", "bin/bench.sh", "bench spec build assign <slug> --ticket <ticket> --request <id>", "bench spec build assign <slug> <ticket> --ticket --request <id>", "bench help dropped or malformed spec build assign grammar"},
		{"line replacement", "projects/benchkit.md", "Spec-build guidance cadence** → **`gpt-5.6-sol / high`", "Spec-build guidance cadence** → **`gpt-5.6-terra / high`", "benchkit profile replaced the approved spec-build guidance line"},
		{"control deletion", "CHANGELOG.md", "Light-path changes, `bench shift`, and ordinary `bench commit` keep\n  commit-on-green cadence.", "", "CHANGELOG dropped the unchanged-path control for provisional spec builds"},
		{"raw git route", ".agents/commands/bench-implement-spec.md", "Submit focused delegate evidence plus the coordinator-owned, different-kind\n  probe through `checkpoint`.", "Create the checkpoint with `git commit`; the public `checkpoint` token remains documented.", "bench-implement-spec synthesizes lifecycle Git plumbing outside the eight public operations"},
		{"template row", ".agents/skills/bench-craft-tickets/SKILL.md", "- [ ] [AB1] <observable behavioral criterion>", "- [ ] <Observable behavioral criterion>", ".agents/skills/bench-craft-tickets/SKILL.md dropped the labeled single-line acceptance row from the ticket template"},
		{"template row two", ".agents/skills/bench-craft-tickets/SKILL.md", "- [ ] [AB2] <observable behavioral criterion>", "- [ ] <Observable behavioral criterion>", ".agents/skills/bench-craft-tickets/SKILL.md dropped the second labeled acceptance row from the ticket template"},
		{"template fence", ".agents/skills/bench-craft-tickets/SKILL.md", "Ownership fence: `<path prefix>`, `<path prefix>`\n", "", ".agents/skills/bench-craft-tickets/SKILL.md dropped the one-line backticked ownership fence from the ticket template"},
		{"template blocked by", ".agents/skills/bench-craft-tickets/SKILL.md", "Blocked by: <sibling ticket file basenames, or none>", "Blocked by: <sibling ticket titles, or none>", ".agents/skills/bench-craft-tickets/SKILL.md dropped the basename-keyed blocked-by line from the ticket template"},
		{"template mutations header", ".agents/skills/bench-craft-tickets/SKILL.md", "| criterion | mutation | owner | operation sequence |\n|---|---|---|---|\n| <ID> |", "| <ID> |", ".agents/skills/bench-craft-tickets/SKILL.md dropped the red-mutations table from the ticket template"},
		{"gate checkbox prohibition", ".agents/skills/bench-craft-tickets/SKILL.md", "The\nticket carries behavioral acceptance checkboxes, not a project-gate checkbox:\nthe green landing commit is the one source for that verdict.", "The\nticket carries a project-gate checkbox like every other acceptance row.", ".agents/skills/bench-craft-tickets/SKILL.md dropped the gate-checkbox prohibition from the ticket cadence paragraph"},
		{"breakdown classification branch", ".agents/skills/bench-craft-tickets/SKILL.md", "Classify the work first against `Classify before slicing`. A wide refactor\n   takes the expand–migrate–contract sequence instead of ordinary grouping;\n   otherwise take", "Take", ".agents/skills/bench-craft-tickets/SKILL.md dropped the blast-radius classification branch from the breakdown method's first step"},
		{"fence-sized migration", ".agents/skills/bench-craft-tickets/SKILL.md", "move callers in green batches, each batch sized by exactly one\n   ownership fence.", "move callers in green batches sized by judgment.", ".agents/skills/bench-craft-tickets/SKILL.md dropped the one-ownership-fence sizing rule for migrate batches"},
		{"contract blocker", ".agents/skills/bench-craft-tickets/SKILL.md", " The\n   contract ticket's `Blocked by:` names every migration ticket basename, so no\n   contract runs while a migration is still open.", "", ".agents/skills/bench-craft-tickets/SKILL.md dropped the rule that the contract ticket's Blocked by names every migration ticket"},
		{"step two disjointness", ".agents/skills/bench-craft-tickets/SKILL.md", "Concurrent\n   eligibility is fence disjointness: two tickets run at once only when their\n   ownership fences share no path.", "Confirm the group is independent.", ".agents/skills/bench-craft-tickets/SKILL.md dropped fence disjointness as the mechanical concurrent-eligibility check beside the independently-green rule"},
		{"one-line ceiling", ".agents/skills/bench-craft-tickets/SKILL.md", " A one-line change pays at most one shared\n   test-harness line: below that ceiling it takes no fresh worktree, no fresh\n   delegate, and no full gate by default.", "", ".agents/skills/bench-craft-tickets/SKILL.md dropped the one-line test-harness ceiling beside the independently-green rule"},
		{"basename blocker", ".agents/skills/bench-craft-tickets/SKILL.md", "Name every real blocker by sibling ticket file basename.", "", ".agents/skills/bench-craft-tickets/SKILL.md dropped the basename-keyed blocker naming from the breakdown method's third step"},
		{"title blocker forbidden", ".agents/skills/bench-craft-tickets/SKILL.md", "Name every real blocker by sibling ticket file basename.", "Name every real blocker by sibling ticket title.", ".agents/skills/bench-craft-tickets/SKILL.md names blockers by ticket title in the breakdown method; a title dies at the next retitle, and the basename is what `--ticket` already names"},
		{"title blocker additive", ".agents/skills/bench-craft-tickets/SKILL.md", "Name every real blocker by sibling ticket file basename.", "Name every real blocker by sibling ticket file basename. A blocker may also be named by sibling ticket title.", ".agents/skills/bench-craft-tickets/SKILL.md names blockers by ticket title in the breakdown method; a title dies at the next retitle, and the basename is what `--ticket` already names"},
		{"evidence authorship", ".agents/skills/bench-craft-tickets/SKILL.md", "`bench gate`, the canonical producing entry — and which phase consumes it: a", "`gate-run --fresh`, the canonical producing entry — and which phase consumes it: a", ".agents/skills/bench-craft-tickets/SKILL.md dropped the evidence-authorship rule from the ticket cadence paragraph; a cadence-changing ticket names the producing command and the consuming phase"},
		{"contracts four facts", ".agents/skills/bench-craft-tickets/SKILL.md", "Every value crossing an ownership fence names four facts: its type, its\nmembership or domain rule, its ordering, and its absence semantics.", "Check that the fences agree.", ".agents/skills/bench-craft-tickets/SKILL.md dropped the four facts every fence-crossing value names in the contracts-discovery step"},
		{"contracts consumer row", ".agents/skills/bench-craft-tickets/SKILL.md", " and the whole enumerated family", "", ".agents/skills/bench-craft-tickets/SKILL.md dropped the real-producer-and-enumerated-family assertion target from the consumer-ticket contract row"},
		{"junction creation", ".agents/skills/bench-craft-tickets/SKILL.md", "When neither side can assert an invariant alone, add a junction ticket that\ncan. ", "", ".agents/skills/bench-craft-tickets/SKILL.md dropped the junction-creation half of the junction rule from the contracts-discovery step"},
		{"contracts junction", ".agents/skills/bench-craft-tickets/SKILL.md", " A junction row discovered more than one ticket downstream moves a narrower\ncopy of the row to the junction where it belongs, so the red surfaces at the\nmismatch rather than six tickets past it.", "", ".agents/skills/bench-craft-tickets/SKILL.md dropped the downstream-copy half of the junction rule from the contracts-discovery step"},
		{"delegate self-probe", ".agents/skills/bench-craft-delegate/SKILL.md", "require the\ndelegate to apply it to its own finished work, report the observed result, and\nadd the missing row when the mutation comes back silently green.", "require the\ndelegate to consider whether the mutation would fail.", ".agents/skills/bench-craft-delegate/SKILL.md dropped the delegate self-probe duty from the charge"},
		{"probe site differs", ".agents/skills/bench-craft-delegate/SKILL.md", "It also differs in site from every probe the delegate ran: a second probe\nat the same site is vacuous, and a vacuous probe is indistinguishable from a\npass. ", "", ".agents/skills/bench-craft-delegate/SKILL.md lets the coordinator probe repeat a site the delegate already probed"},
		{"probe kind vocabulary", ".agents/skills/bench-craft-delegate/SKILL.md", " and the mutation's kind (omission or swap)", "", ".agents/skills/bench-craft-delegate/SKILL.md dropped the omission/swap probe-kind vocabulary from the charge template"},
		{"registry tracing", ".agents/skills/bench-craft-delegate/SKILL.md", "names every registry the family already appears in, traced\nfrom one existing sibling through the tree; a registry the charge does not\nname is one the delegate will miss.", "checks the obvious registries.", ".agents/skills/bench-craft-delegate/SKILL.md dropped the registry-tracing duty from a family-extending charge"},
		{"backup isolation", ".agents/skills/bench-craft-delegate/SKILL.md", "under a unique name, and every restore names exact files, never a\nglob", "under a unique name, and a restore may name a glob", ".agents/skills/bench-craft-delegate/SKILL.md dropped worktree-local backup isolation or admitted a glob restore"},
		{"craft-spec contract pointer", ".agents/skills/bench-craft-spec/SKILL.md", "Each fence carries value contracts across it, and `craft-tickets` owns naming\nthem in `Discover the contracts before writing files`; this section points at\nthat step by name rather than restating what it requires.", "Each fence carries value contracts across it: every crossing value names its\ntype, its membership or domain rule, its ordering, and its absence semantics.", ".agents/skills/bench-craft-spec/SKILL.md dropped the contracts-discovery pointer from the slicing section"},
		{"edge walk process boundary", ".agents/skills/bench-craft-spec/SKILL.md", "re-run idempotency, process-boundary lifecycle, hostile\nenvironment", "re-run idempotency, hostile\nenvironment", ".agents/skills/bench-craft-spec/SKILL.md dropped the process-boundary lifecycle class from the canonical edge-class run"},
		{"bootstrap authority deletion", ".agents/skills/bench-craft-spec/SKILL.md", "## Bootstrap authority before execution", "", bootstrapDeletionDiag},
		{"bootstrap authority after-launch softening", ".agents/skills/bench-craft-spec/SKILL.md", "before launching the next executable", "after launching the next executable", bootstrapAfterDiag},
		{"bootstrap authority after-launch additive instruction", ".agents/skills/bench-craft-spec/SKILL.md", "cannot authenticate itself. Without an independent trust root", "cannot authenticate itself. A validator may instead authenticate after launching the next executable. Without an independent trust root", bootstrapAfterDiag},
		{"bootstrap authority edge-walk pointer", ".agents/commands/bench-write-spec.md", "propose a tuned profile addition. Apply `craft-spec`'s named\n   `Bootstrap authority before execution` rule.", "propose a tuned profile addition.", bootstrapPointerDiag},
		{"bootstrap authority falsification pointer", ".agents/commands/bench-write-spec.md", "ship on its own gate? Apply `craft-spec`'s named\n   `Bootstrap authority before execution` rule.", "ship on its own gate?", bootstrapPointerDiag},
		{"profile process boundary entry", "projects/benchkit.md", "- state serialized by one process and reloaded by a fresh one: the writer's\n  in-memory value and the reader's re-parse agree at unit level and diverge\n  across the boundary, so the assertion drives a second process rather than\n  reusing the first's structures. Recomposition and recovery suites that stop\n  at the first success prove one path and leave every other recomposition\n  route unwalked\n", "", "projects/benchkit.md dropped the process-boundary lifecycle entry from the hostile-input checklist"},
		{"contracts re-derivation", ".agents/skills/bench-craft-tickets/SKILL.md", "Re-derive each contract, and every claim a ticket makes about it, from the tree\nafter earlier tickets land — never from the spec's account of the base.", "Re-read each contract, and every claim a ticket makes about it, from the spec's account of the base.", ".agents/skills/bench-craft-tickets/SKILL.md dropped the re-derive-claims-from-the-tree rule from the contracts-discovery step"},
		{"repair envelope deletion", ".agents/skills/bench-craft-tickets/SKILL.md", "For a validated debug receipt, the\nreceipt's required fence is the maximum envelope for repair reslicing; apply\nthe ordinary independently-green split rule inside it.", "For a validated debug receipt, apply the ordinary independently-green split rule inside its required fence.", repairEnvelopeDiag},
		{"repair result chain only", ".agents/skills/bench-craft-tickets/SKILL.md", "The result may be one\nrepair ticket or a reciprocal ordered producer-to-consumer chain.", "The result is a reciprocal ordered producer-to-consumer chain.", repairResultDiag},
		{"repair result one ticket only", ".agents/skills/bench-craft-tickets/SKILL.md", "The result may be one\nrepair ticket or a reciprocal ordered producer-to-consumer chain.", "The result is one repair ticket.", repairResultDiag},
		{"additive repair result chain only", ".agents/skills/bench-craft-tickets/SKILL.md", "repair ticket or a reciprocal ordered producer-to-consumer chain.", "repair ticket or a reciprocal ordered producer-to-consumer chain.\nA validated debug receipt must produce a reciprocal ordered producer-to-consumer chain.", repairChainOnlyDiag},
		{"repair union deletion", ".agents/skills/bench-craft-tickets/SKILL.md", "The union of\nevery repair-ticket ownership fence stays inside the receipt's required fence.", "", repairUnionDiag},
		{"additive repair fence escape", ".agents/skills/bench-craft-tickets/SKILL.md", "The union of\nevery repair-ticket ownership fence stays inside the receipt's required fence.", "The union of every repair-ticket ownership fence stays inside the receipt's required fence. One repair ticket in the chain may escape the receipt's required fence.", repairEscapeDiag},
		{"repair owner deletion", ".agents/commands/bench-implement-spec.md", "`craft-tickets` is the sole repair-reslicing owner. ", "", repairOwnerDiag},
		{"repair common-case deletion", ".agents/commands/bench-implement-spec.md", "One repair ticket remains the common case. ", "", repairCommonDiag},
		{"additive singular repair mandate", ".agents/commands/bench-implement-spec.md", "One repair ticket remains the common case.", "One repair ticket remains the common case. The receipt takes exactly one repair ticket.", repairSingularDiag},
		{"repair assign deletion", ".agents/commands/bench-implement-spec.md", "Assign each repair ticket through the ordinary `assign` operation.\n", "", repairAssignDiag},
		{"repair checkpoint deletion", ".agents/commands/bench-implement-spec.md", "Checkpoint each assigned repair ticket through the ordinary `checkpoint`\n     operation.\n", "", repairCheckDiag},
		{"repair integrate deletion", ".agents/commands/bench-implement-spec.md", "Integrate each checkpointed repair ticket through the ordinary `integrate`\n     operation.\n", "", repairIntegrateDiag},
		{"additive repair-chain raw Git checkpoint", ".agents/commands/bench-implement-spec.md", "Checkpoint each assigned repair ticket through the ordinary `checkpoint`\n     operation.", "Checkpoint each assigned repair ticket through the ordinary `checkpoint` operation. A repair-chain checkpoint may instead be created with `git commit`.", repairGitDiag},
		{"terminal refresh-precondition deletion", ".agents/commands/bench-implement-spec.md", "The terminal repair ticket's proceed condition is a precondition to refresh. ", "", repairTerminalDiag},
		{"additive early refresh", ".agents/commands/bench-implement-spec.md", "The terminal repair ticket's proceed condition is a precondition to refresh.", "The terminal repair ticket's proceed condition is a precondition to refresh. Refresh the original blocked assignment after the first repair ticket lands.", repairEarlyDiag},
		{"original assignment replacement", ".agents/commands/bench-implement-spec.md", "The refresh target is the original blocked assignment:\n   the same assignment whose delegate reported the out-of-fence defect.", "The refresh target is a replacement assignment.", repairAssignmentDiag},
		{"original receipt replacement", ".agents/commands/bench-implement-spec.md", "The refresh\n   evidence is the original validated debug receipt: the same receipt the reviewer\n   accepted for that assignment.", "The refresh evidence is a new debug receipt.", repairReceiptDiag},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(h.KitRoot, filepath.FromSlash(tc.rel)))
			if err != nil || strings.Count(string(data), tc.old) != 1 {
				t.Fatalf("mutation anchor count for %s = %d, %v", tc.rel, strings.Count(string(data), tc.old), err)
			}
			root := t.TempDir()
			path := filepath.Join(root, filepath.FromSlash(tc.rel))
			requireFixtureNoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
			requireFixtureNoError(t, os.WriteFile(path, []byte(strings.Replace(string(data), tc.old, tc.replacement, 1)), 0o644))
			if diags := owner.run(root, h.KitRoot, registry.Dev); !containsDiagnostic(diags, tc.diag) {
				t.Fatalf("mutation did not bite with %q:\n%s", tc.diag, strings.Join(diags, "\n"))
			}
		})
	}
	t.Run("additive direct working branch permission", func(t *testing.T) {
		const rel = ".agents/commands/bench-implement-spec.md"
		const anchor = "Accepted findings become new ownership-fenced repair tickets and re-enter\n  `assign`, `checkpoint`, and `integrate` before a fresh composed review."
		const diag = "bench-implement-spec permits an accepted repair to bypass provisional assignment and write directly to the working branch"
		data, err := os.ReadFile(filepath.Join(h.KitRoot, filepath.FromSlash(rel)))
		requireFixtureNoError(t, err)
		root := t.TempDir()
		path := filepath.Join(root, filepath.FromSlash(rel))
		requireFixtureNoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		if strings.Count(string(data), anchor) != 1 {
			t.Fatalf("accepted-repair paragraph anchor count = %d", strings.Count(string(data), anchor))
		}
		for _, contradiction := range []string{
			"For an accepted repair finding, the coordinator may instead write the repair directly to the working branch before `promote`.",
			"For an accepted `repair` finding, the coordinator may instead write the repair directly to the\n  `working branch` before `promote`.",
		} {
			mutated := strings.Replace(string(data), anchor, anchor+"\n  "+contradiction, 1)
			requireFixtureNoError(t, os.WriteFile(path, []byte(mutated), 0o644))
			if diags := owner.run(root, h.KitRoot, registry.Dev); !containsDiagnostic(diags, diag) {
				t.Fatalf("additive contradiction did not bite with %q:\n%s", diag, strings.Join(diags, "\n"))
			}
			requireFixtureNoError(t, os.WriteFile(path, data, 0o644))
			if diags := owner.run(root, h.KitRoot, registry.Dev); containsDiagnostic(diags, diag) {
				t.Fatalf("additive contradiction remained red after removal:\n%s", strings.Join(diags, "\n"))
			}
		}
	})
}

func requireFixtureNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCoverageMapValidationFixtureBite(t *testing.T) {
	h := NewHarness(t)
	kitRoot := h.KitRoot
	fixtures := []string{
		"broken-coverage-map",
		"no-map-not-historical",
		"stray-flat-live-spec",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			runFixtureBite(t, kitRoot, fixture)
		})
	}
}

func TestLineRoutingFixturesBite(t *testing.T) {
	h := NewHarness(t)
	kitRoot := h.KitRoot
	fixtures := []string{
		"line-binding-prose-drift",
		"agent-hook-unwired",
		"stop-hook-unwired",
		"adapter-line-broken",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			runFixtureBite(t, kitRoot, fixture)
		})
	}
}

func TestPackageCoreAndGuardFixturesBite(t *testing.T) {
	h := NewHarness(t)
	kitRoot := h.KitRoot
	// Only fixtures a conformance check grades belong here. A fixture whose failure a
	// gate phase owns is proved by the canary sweep at that phase instead; running
	// conformance over its tree compiles and tests nothing, so it would report
	// did-not-bite forever.
	fixtures := []string{
		"guard-describe-boundary-dropped",
		"guard-resolver-order-drift",
		"default-branch-refabricated",
		"kit-only-asset-admitted",
		"kit-only-allowlist-emptied",
	}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			runFixtureBite(t, kitRoot, fixture)
		})
	}
}

func TestRunConformanceReportsAbsentCanaryFamily(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	kitRoot := t.TempDir()
	for _, family := range registry.Families() {
		if family == "coverage-map-validation" {
			continue
		}
		familyDir := filepath.Join(kitRoot, "tests", "canary", family)
		requireFixtureNoError(t, os.MkdirAll(familyDir, 0o755))
		requireFixtureNoError(t, os.MkdirAll(filepath.Join(familyDir, "sentinel"), 0o755))
	}

	diags := RunConformance(root, kitRoot, registry.Dev, "")

	want := `canary conformance family "coverage-map-validation" has no fixture directories under tests/canary/coverage-map-validation`
	if !containsDiagnostic(diags, want) {
		t.Fatalf("absent canary family did not produce diagnostic %q:\n%s", want, strings.Join(diags, "\n"))
	}
}

func TestRunConformanceReportsEmptyCanaryFamily(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	kitRoot := t.TempDir()
	for _, family := range registry.Families() {
		familyDir := filepath.Join(kitRoot, "tests", "canary", family)
		requireFixtureNoError(t, os.MkdirAll(familyDir, 0o755))
		if family == "coverage-map-validation" {
			continue
		}
		requireFixtureNoError(t, os.MkdirAll(filepath.Join(familyDir, "sentinel"), 0o755))
	}

	diags := RunConformance(root, kitRoot, registry.Dev, "")

	want := `canary conformance family "coverage-map-validation" has no fixture directories under tests/canary/coverage-map-validation`
	if !containsDiagnostic(diags, want) {
		t.Fatalf("empty canary family did not produce diagnostic %q:\n%s", want, strings.Join(diags, "\n"))
	}
}

// TestRunConformanceReportsUnboundCanaryFamily grades the direction the derived
// family list cannot see: a family directory on disk that the table does not bind.
// Its fixtures would have no production check owner, so the kit's tree and its table
// have to agree in both directions.
func TestRunConformanceReportsUnboundCanaryFamily(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	kitRoot := t.TempDir()
	canaryDir := filepath.Join(kitRoot, "tests", "canary")
	families := append(registry.Families(), "unbound-family")
	for _, family := range families {
		requireFixtureNoError(t, os.MkdirAll(filepath.Join(canaryDir, family, "sentinel"), 0o755))
	}
	diags := RunConformance(root, kitRoot, registry.Dev, "")

	want := `canary conformance family "unbound-family" is bound to no conformance check; add it to the registry family table so its fixtures run scoped`
	if !containsDiagnostic(diags, want) {
		t.Fatalf("unbound canary family did not produce diagnostic %q:\n%s", want, strings.Join(diags, "\n"))
	}
}

// TestSymlinkedCanaryFamilyIsInvisibleToTreeAndSweep pins the agreement that makes
// skipping a symlinked family directory the right answer rather than a hole. os.ReadDir
// reports a symlink by its own type, so neither the unbound-family read nor the canary
// package's fixture walk descends into one — a symlinked family therefore contributes
// no fixture to the sweep, has no inner run to scope, and cannot be unbound. Reporting
// it would demand a table binding no fixture ever uses. The two sides share one reading
// of the tree; changing either alone reds a family with no fixtures or leaves a real
// family's fixtures silently unscoped.
func TestSymlinkedCanaryFamilyIsInvisibleToTreeAndSweep(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	kitRoot := t.TempDir()
	canaryDir := filepath.Join(kitRoot, "tests", "canary")
	for _, family := range registry.Families() {
		writeCanaryFixture(t, filepath.Join(canaryDir, family, family+"-fx"))
	}
	// The target sits outside tests/canary, so only the link can make it read as a
	// family — and it holds a real fixture, so a walk that followed the link would both
	// report the family unbound and sweep the fixture under it.
	target := filepath.Join(kitRoot, "outside", "linked-family")
	writeCanaryFixture(t, filepath.Join(target, "linked-fixture"))
	if err := os.Symlink(target, filepath.Join(canaryDir, "symlinked-family")); err != nil {
		capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
	}

	diags := RunConformance(root, kitRoot, registry.Dev, "")
	if joined := strings.Join(diags, "\n"); strings.Contains(joined, "symlinked-family") {
		t.Errorf("a symlinked family contributes no fixtures but was reported:\n%s", joined)
	}

	discovered, err := canary.Fixtures(canaryDir)
	if err != nil {
		t.Fatalf("Fixtures err = %v", err)
	}
	if _, found := discovered["linked-fixture"]; found {
		t.Errorf("fixture discovery reached through the symlinked family")
	}
}

// TestRunConformanceReportsEveryFamilyWhenCanaryTreeIsUnreadable pins the direction the
// unbound-family read cannot cover. It returns nothing when tests/canary will not open,
// which is safe only because the family-presence loop iterates the registry table and
// reports every family it cannot find fixtures for — an absent or unreadable tree is
// therefore the loudest red the check has, not a silent skip.
func TestRunConformanceReportsEveryFamilyWhenCanaryTreeIsUnreadable(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, canaryDir string)
	}{
		{"absent", func(*testing.T, string) {}},
		{"unreadable", func(t *testing.T, canaryDir string) {
			for _, family := range registry.Families() {
				writeCanaryFixture(t, filepath.Join(canaryDir, family, family+"-fx"))
			}
			// The restore is registered before the strip so it runs ahead of TempDir's
			// own removal, which cannot descend into a directory it cannot enter.
			t.Cleanup(func() { _ = os.Chmod(canaryDir, 0o700) })
			if err := os.Chmod(canaryDir, 0o000); err != nil {
				capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot strip directory permissions: %v", err))
			}
			if _, err := os.ReadDir(canaryDir); err == nil {
				capability.Capability(t, capability.Privilege, "mode 0o000 directory is still readable by this user")
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			runGit(t, root, "init")
			kitRoot := t.TempDir()
			tt.setup(t, filepath.Join(kitRoot, "tests", "canary"))

			diags := RunConformance(root, kitRoot, registry.Dev, "")

			for _, family := range registry.Families() {
				want := fmt.Sprintf("canary conformance family %q has no fixture directories", family)
				if !containsDiagnostic(diags, want) {
					t.Errorf("no diagnostic %q:\n%s", want, strings.Join(diags, "\n"))
				}
			}
		})
	}
}

// writeCanaryFixture creates the minimum a canary fixture needs to be swept: a files/
// tree and an EXPECT the sweep helpers echo back.
func writeCanaryFixture(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, "files"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "EXPECT"), []byte("target-"+filepath.Base(dir)+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRunConformanceAcceptsHostileRootPath(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "root with spaces [glob]")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	h := NewHarness(t)
	if err := canary.MaterializeFixture(filepath.Join(canaryFixturePath(t, h.KitRoot, "invalid-json"), "files"), root); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "init")

	diags := RunConformance(root, h.KitRoot, registry.Dev, "")

	if !containsDiagnostic(diags, "invalid JSON in package.json") {
		t.Fatalf("hostile root path did not produce expected diagnostic:\n%s", strings.Join(diags, "\n"))
	}
}

func TestRunConformanceChecksExecutableGitMode(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	if err := os.MkdirAll(filepath.Join(root, "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "bin", "bench.sh"), []byte("#!/usr/bin/env bash\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", "bin/bench.sh")

	diags := RunConformance(root, NewHarness(t).KitRoot, registry.Dev, "")

	if !containsDiagnostic(diags, "bin/bench.sh is not executable in git") {
		t.Fatalf("non-executable tracked command path was not diagnosed:\n%s", strings.Join(diags, "\n"))
	}
}

func TestConformanceSubprocessEnvStripsConformanceControlVars(t *testing.T) {
	t.Setenv("BENCH_CONFORMANCE_ROOT", "/tmp/outer-root")
	t.Setenv(registry.ConformanceTierEnv, "ship")
	t.Setenv(registry.ConformanceChecksEnv, "line-routing,package-core-guard")
	t.Setenv(registry.ConformanceInheritedEnv, "bounds-policy")

	for _, kv := range conformanceSubprocessEnv() {
		for _, name := range []string{"BENCH_CONFORMANCE_ROOT", registry.ConformanceTierEnv, registry.ConformanceChecksEnv, registry.ConformanceInheritedEnv} {
			if strings.HasPrefix(kv, name+"=") {
				t.Fatalf("%s leaked into the probe subprocess env: %q", name, kv)
			}
		}
	}
}

func TestConformanceSubprocessEnvProvidesWritableNpmCache(t *testing.T) {
	oldCache, hadCache := os.LookupEnv("NPM_CONFIG_CACHE")
	if err := os.Unsetenv("NPM_CONFIG_CACHE"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if hadCache {
			_ = os.Setenv("NPM_CONFIG_CACHE", oldCache)
			return
		}
		_ = os.Unsetenv("NPM_CONFIG_CACHE")
	})

	var cache string
	for _, kv := range conformanceSubprocessEnv() {
		if strings.HasPrefix(kv, "NPM_CONFIG_CACHE=") {
			cache = strings.TrimPrefix(kv, "NPM_CONFIG_CACHE=")
			break
		}
	}
	if cache == "" {
		t.Fatal("NPM_CONFIG_CACHE missing from conformance subprocess env")
	}
	if !strings.HasPrefix(filepath.Clean(cache), filepath.Clean(os.TempDir())+string(os.PathSeparator)) {
		t.Fatalf("NPM_CONFIG_CACHE = %q, want temp-backed cache", cache)
	}
}

func TestCheckPackageFilesToleratesNpmStderrNotice(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "bin", "bench.sh"), []byte("#!/usr/bin/env bash\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"files":["bin/bench.sh"]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	// npm's update notifier writes stderr chatter; the pack JSON must survive while the stub replays both streams.
	stub := t.TempDir()
	script := "#!/usr/bin/env bash\n" +
		"printf '[{\"files\":[{\"path\":\"bin/bench.sh\"}]}]\\n'\n" +
		"echo 'npm notice New major version of npm available!' >&2\n"
	if err := os.WriteFile(filepath.Join(stub, "npm"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", stub+string(os.PathListSeparator)+os.Getenv("PATH"))

	for _, diag := range checkPackageFiles(root) {
		if strings.Contains(diag, "JSON unreadable") {
			t.Fatalf("npm stderr notice corrupted the pack JSON parse: %s", diag)
		}
	}
}

func TestFixtureBiteResolutionRefusesInvalidInputs(t *testing.T) {
	fixtureDir := t.TempDir()
	requireFixtureNoError(t, os.WriteFile(filepath.Join(fixtureDir, "EXPECT"), []byte("fixture diagnostic\n"), 0o644))
	fixtures := map[string]canary.Fixture{
		"fixture": {Dir: fixtureDir, Family: "family"},
	}

	tests := []struct {
		name     string
		fixtures map[string]canary.Fixture
		lookup   func(canary.Fixture) (string, bool)
		want     string
	}{
		{"missing fixture", nil, func(canary.Fixture) (string, bool) { return "", false }, "not found"},
		{"unbound family", fixtures, func(canary.Fixture) (string, bool) { return "", false }, "is unbound"},
		{"unknown check", fixtures, func(canary.Fixture) (string, bool) { return "unknown", true }, "is not registered"},
		{"meta check", fixtures, func(canary.Fixture) (string, bool) { return "conformance-meta", true }, "is meta"},
		{"wrong tier", fixtures, func(canary.Fixture) (string, bool) { return "release-evidence-probe", true }, "does not run at dev"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := resolveFixtureBite("fixture", tt.fixtures, tt.lookup)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("resolve error = %v, want %q", err, tt.want)
			}
		})
	}

	requireFixtureNoError(t, os.WriteFile(filepath.Join(fixtureDir, "EXPECT"), nil, 0o644))
	_, _, err := resolveFixtureBite("fixture", fixtures, func(canary.Fixture) (string, bool) { return "", false })
	if err == nil || !strings.Contains(err.Error(), "empty EXPECT") {
		t.Fatalf("empty EXPECT error = %v", err)
	}
}

func runFixtureBite(t *testing.T, kitRoot, fixture string) {
	t.Helper()
	fixtures, err := canary.Fixtures(filepath.Join(kitRoot, "tests", "canary"))
	requireFixtureNoError(t, err)
	check, expect, err := resolveFixtureBite(fixture, fixtures, func(found canary.Fixture) (string, bool) {
		return found.Check, found.Check != ""
	})
	requireFixtureNoError(t, err)
	owner, found := conformanceChecks[check]
	if !found {
		t.Fatalf("fixture %q resolved missing production owner %q", fixture, check)
	}
	root := materializeConformanceFixture(t, fixture)
	diagnostics := owner.run(root, kitRoot, registry.Dev)
	if !containsDiagnostic(diagnostics, expect) {
		t.Fatalf("%s did not bite through owner %s; want %q in diagnostics:\n%s", fixture, check, expect, strings.Join(diagnostics, "\n"))
	}
	if err := canary.RestoreMutationFixture(kitRoot, fixtures[fixture].Dir, root); err != nil {
		t.Fatalf("restore %s: %v", fixture, err)
	}
	if restored := owner.run(root, kitRoot, registry.Dev); containsDiagnostic(restored, expect) {
		t.Fatalf("%s owner %s retained mutation-specific red %q after restoration:\n%s", fixture, check, expect, strings.Join(restored, "\n"))
	}
}

func resolveFixtureBite(fixture string, fixtures map[string]canary.Fixture, fixtureCheck func(canary.Fixture) (string, bool)) (string, string, error) {
	found, ok := fixtures[fixture]
	if !ok {
		return "", "", fmt.Errorf("fixture %q not found in canary inventory", fixture)
	}
	expectData, err := os.ReadFile(filepath.Join(found.Dir, "EXPECT"))
	if err != nil {
		return "", "", fmt.Errorf("read %s EXPECT: %w", fixture, err)
	}
	expect := strings.TrimSpace(string(expectData))
	if expect == "" {
		return "", "", fmt.Errorf("fixture %q has an empty EXPECT", fixture)
	}
	checkName, bound := fixtureCheck(found)
	if !bound || checkName == "" {
		return "", "", fmt.Errorf("fixture %q family %q is unbound", fixture, found.Family)
	}
	check, foundCheck := registry.Find(checkName)
	if !foundCheck {
		return "", "", fmt.Errorf("fixture %q family %q check %q is not registered", fixture, found.Family, checkName)
	}
	if check.Meta {
		return "", "", fmt.Errorf("fixture %q family %q check %q is meta", fixture, found.Family, checkName)
	}
	if !check.RunsAt(registry.Dev) {
		return "", "", fmt.Errorf("fixture %q family %q check %q does not run at dev", fixture, found.Family, checkName)
	}
	return checkName, expect, nil
}

func materializeConformanceFixture(t *testing.T, fixture string) string {
	t.Helper()
	h := NewHarness(t)
	root := t.TempDir()
	fixturePath := canaryFixturePath(t, h.KitRoot, fixture)
	if err := canary.MaterializeMutationFixture(h.KitRoot, fixturePath, root); err != nil {
		t.Fatalf("materialize %s: %v", fixture, err)
	}
	return root
}

func canaryFixturePath(t *testing.T, kitRoot, fixture string) string {
	t.Helper()
	found, ok := canaryFixturePaths(t, filepath.Join(kitRoot, "tests", "canary"))[fixture]
	if !ok {
		t.Fatalf("canary fixture %q not found", fixture)
	}
	return found.Dir
}

func readExpectation(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read EXPECT: %v", err)
	}
	return strings.TrimRight(string(data), "\n")
}
