package anchors

import (
	"strings"

	"github.com/gibbonmi/bench/internal/retros"
)

// retainedWorkflowAnchors is the whole retained-workflow family: the retired
// default and delegated workflow forms, the
// lane-and-landing split, the retro capture route, the declared author line, the
// fact owners that the guides point to, the drain and craft skill references, and
// the stale implementation instructions. The registry reads this one name, so its
// groups evaluate together.
var retainedWorkflowAnchors = append(append(append(append(append(append(append(append([]Anchor{}, defaultWorkflowAnchors...), delegatedWorkflowAnchors...), laneAndLandingAnchors...), retroCaptureAnchors...), declaredLineAnchors...), factOwnerAnchors...), referenceRouteAnchors...), staleImplementationAnchors...)

// staleImplementationAnchors keep each retired implementation instruction out of the
// live workflow reader that once carried it. The matcher collapses whitespace, so a
// needle also matches the instruction when it wraps across lines.
var staleImplementationAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Forbid, Needle: "Every spec-backed run assigns genuine write work to a write subagent", Diagnostic: "retained workflow: implementation phase restored the write subagent for the write work of every spec-backed run"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-write-spec.md", Kind: Forbid, Needle: "A fresh session builds after ticket approval", Diagnostic: "retained workflow: spec phase restored the fresh build session after ticket approval"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-write-spec.md", Kind: Forbid, Needle: "fresh mid-tier build session", Diagnostic: "retained workflow: spec phase restored the fresh mid-tier build session"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-drain.md", Kind: Forbid, Needle: "one later write delegate authors the complete tracked batch", Diagnostic: "retained workflow: drain restored the later write delegate that authors the tracked batch"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-drain.md", Kind: Forbid, Needle: "one later write delegate to draft the complete tracked pass", Diagnostic: "retained workflow: drain restored the later write delegate that drafts the tracked pass"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Forbid, Needle: "Initial review blocks on the full", Diagnostic: "retained workflow: review phase restored the initial review block on the full review"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Forbid, Needle: "repair-scoped re-review", Diagnostic: "retained workflow: review phase restored the repair-scoped re-review"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-tickets/SKILL.md", Kind: Forbid, Needle: "one fresh write-delegate charge", Diagnostic: "retained workflow: craft-tickets restored the one fresh write-delegate charge"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/SKILL.md", Kind: Forbid, Needle: "All other code authorship runs as a write-delegation", Diagnostic: "retained workflow: craft-delegate restored the write delegation for all other code authorship"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/SKILL.md", Kind: Forbid, Needle: "Otherwise a fresh charge in an isolated worktree carries the finding", Diagnostic: "retained workflow: craft-delegate restored the fresh charge that carries a finding"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-spec/SKILL.md", Kind: Forbid, Needle: "A build may not edit its own spec's acceptance rows", Diagnostic: "retained workflow: craft-spec restored the ban on a build that edits its own acceptance rows"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/references/delegation-discipline.md", Kind: Forbid, Needle: "takes a fence extension in a continuation", Diagnostic: "retained workflow: delegation discipline restored the fence extension in a continuation"},
	{Group: AfterImplementSpec, File: "docs/field-guide.html", Kind: Forbid, Needle: "their build starts in a fresh session", Diagnostic: "retained workflow: field guide restored the fresh build session"},
	{Group: AfterImplementSpec, File: "docs/field-guide.html", Kind: Forbid, Needle: "per fresh write-delegate context", Diagnostic: "retained workflow: field guide restored the fresh write-delegate context"},
	{Group: AfterImplementSpec, File: "projects/benchkit.md", Kind: Forbid, Needle: "After ticket approval, a fresh mid-tier session starts the build.", Diagnostic: "retained workflow: project profile restored the fresh mid-tier build session after ticket approval"},
}

// referenceRouteAnchors pin the reachable routes of the drain and assess commands and
// each craft skill's pointer to the owner of a fact that another file holds. A Require
// row pins the worktree list target slot rule. Forbid rows keep out the wrong ledger
// file, the wrong rule owner, the hand-written handoff and its whole-file dating, the
// hand-append park route, the blanket adapter rule, the terminal-only disclosure, the
// literal story count, and each retired copy of an owner's sentence.
var referenceRouteAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".agents/commands/bench-drain.md", Kind: Forbid, Needle: "`Occurrences:` line in `ROADMAP.md`", Diagnostic: "reference route: drain restored ROADMAP.md as the file of the Occurrences: ledger line"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-drain.md", Kind: Forbid, Needle: "(the AGENTS.md rule)", Diagnostic: "reference route: drain restored AGENTS.md as the owner of the batch approval rule"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-drain.md", Kind: Forbid, Needle: "A batch approval covers per-spec sign-offs when I'm unreachable.", Diagnostic: "reference route: drain restored a copy of the operating guide's batch approval rule"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-drain.md", Kind: Forbid, Needle: "a standing batch approval counts, with contestable calls flagged for post-hoc veto", Diagnostic: "reference route: drain restored its paraphrase of the batch approval rule that the operating guide owns"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-drain.md", Kind: Forbid, Needle: "When the handoff and the tree disagree, the tree wins.", Diagnostic: "reference route: drain restored its copy of the working agreement's tree-wins rule"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-drain.md", Kind: Forbid, Needle: "The coordinator writes ignored `capture/session-handoff.md` last.", Diagnostic: "reference route: drain restored the hand-written handoff as its last write"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-drain.md", Kind: Forbid, Needle: "`bench status` dates the ignored handoff by its write time.", Diagnostic: "reference route: drain restored the write-time dating of the whole handoff"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-drain.md", Kind: Forbid, Needle: "Its pin block names the pre-commit HEAD", Diagnostic: "reference route: drain restored the pre-commit HEAD pin of a handoff that it writes after the landing"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-assess.md", Kind: Forbid, Needle: "or into `capture/IDEAS.md`", Diagnostic: "reference route: assess restored the hand-append park route into capture/IDEAS.md"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-skills/SKILL.md", Kind: Forbid, Needle: "and the phase adapters are not", Diagnostic: "reference route: craft-skills restored the blanket rule that no phase adapter is model-invoked"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-cli/SKILL.md", Kind: Forbid, Needle: "Only an over-cap default discloses", Diagnostic: "reference route: craft-cli bench consumers row restored the claim that only an over-cap default discloses"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-cli/SKILL.md", Kind: Require, Class: Command, Needle: "The active rows with a present tree share one `bench worktree path <target>` action and one `bench worktree exec <target> -- <command>` action.", Diagnostic: "reference route: craft-cli bench worktree list row dropped the target slot rule for the active-row actions"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-spec/SKILL.md", Kind: Forbid, Needle: "more than four stories", Diagnostic: "reference route: craft-spec restored the literal story count that bench coverage --check owns"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-spec/SKILL.md", Kind: Forbid, Needle: "Each planned chunk has a stable ID and names its tickets", Diagnostic: "reference route: craft-spec restored its copy of the implementation-chunk contract"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-spec/SKILL.md", Kind: Forbid, Needle: "Each planned chunk names its stable ID, outcome, acceptance rows, and tests.", Diagnostic: "reference route: craft-spec restored a copy of the operating guide's implementation-chunk contract"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-tdd/references/tests.md", Kind: Forbid, Needle: "An independently authored test expectation counts as duplicated", Diagnostic: "reference route: craft-tdd tests reference restored a copy of the project's test-expectation standard"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-adr/SKILL.md", Kind: Forbid, Needle: "No file paths, no code snippets", Diagnostic: "reference route: craft-adr restored its copy of invariant 3"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-adr/SKILL.md", Kind: Forbid, Needle: "they go stale the next session and then mislead", Diagnostic: "reference route: craft-adr restored its copy of the reason for invariant 3"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-adr/SKILL.md", Kind: Forbid, Needle: "Put no file paths and no code snippets in an ADR", Diagnostic: "reference route: craft-adr restored a copy of the operating guide's invariant 3"},
}

// Each constant below is an owner sentence. Every row that keeps a copy out of a former
// copy site, and any row that pins the sentence at its owner, reads the same constant.
const (
	// handoffSectionRule is the reference's owner sentence for the handoff verb.
	handoffSectionRule = "`bench handoff` rewrites only the calling worktree's assignment section."
	// handoffMainOwner is the reference's owner sentence for the handoff's main section.
	handoffMainOwner = "The primary checkout owns the `main` section."
	// handoffNextRule is the reference's owner sentence for the handoff's Next command.
	handoffNextRule = "The verb keeps a non-empty Next command."
	// handoffAncestryRefusal is the reference's owner clause for the handoff State refusal.
	handoffAncestryRefusal = "refuses a State that pins a commit outside the tip's ancestry"
	// censusSignalAccount is the reference's owner sentence for the census signal.
	censusSignalAccount = "The `census` signal counts raw calls per assignment from `$BENCH_HOME/census/<repo-key>/`."
	// greenRunGateOutput opens the reference's green-run gate output sentence.
	greenRunGateOutput = "A green run prints one `phases[N]{phase,verdict,elapsed_ms}` table"
	// redRunGateOutput opens the reference's red-run gate output sentence.
	redRunGateOutput = "A red run prints one `failures[N]{phase,line}` table"
	// skillLinkRule is the Claude README's owner sentence for the skill links.
	skillLinkRule = "`.claude/skills/` links every `.agents/skills/` skill that has no same-named command."
)

// factOwnerAnchors keep each guide's facts at the file that holds them. A Require row
// pins the command registry as the owner of the plumbing subcommands. Forbid rows keep
// out the retired copies, the reference's verb grammar, the false skill-link claims,
// the any-branch commit claim, the gate-priced capture commit, and the gate-then-commit
// wording of the light path and the drain.
var factOwnerAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".bench/BENCH.md", Kind: Forbid, Needle: "live in `.bench/BENCH-reference.md`", Diagnostic: "fact owner: operating guide restored the reference as the home of the plumbing subcommands"},
	{Group: AfterImplementSpec, File: ".bench/BENCH.md", Kind: Require, Class: Invariant, Needle: "The command registry owns the plumbing subcommands that hooks and adapters drive.", Diagnostic: "fact owner: operating guide dropped the command registry as the owner of the plumbing subcommands"},
	{Group: AfterImplementSpec, File: ".bench/BENCH.md", Kind: Forbid, Needle: "gate and commit on green", Diagnostic: "fact owner: operating guide light path restored the gate-then-commit wording"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-drain.md", Kind: Forbid, Needle: "commit on green", Diagnostic: "fact owner: drain restored the gate-then-commit wording for a worktree commit"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-drain.md", Kind: Forbid, Needle: "The gate is what a commit costs", Diagnostic: "fact owner: drain restored the gate as the price of a worktree commit"},
	{Group: AfterImplementSpec, File: "AGENTS.md", Kind: Forbid, Needle: "the four invariants, how the pieces fit", Diagnostic: "fact owner: working agreement restored how the pieces fit as a shared platform rule"},
	{Group: AfterImplementSpec, File: "AGENTS.md", Kind: Forbid, Needle: "the communication rules, and the skills index", Diagnostic: "fact owner: working agreement restored the skills index as a shared platform rule"},
	{Group: AfterImplementSpec, File: "AGENTS.md", Kind: Forbid, Needle: handoffMainOwner, Diagnostic: "fact owner: working agreement restored its copy of the handoff main-section owner"},
	{Group: AfterImplementSpec, File: "AGENTS.md", Kind: Forbid, Needle: handoffNextRule, Diagnostic: "fact owner: working agreement restored its copy of the handoff Next-command rule"},
	{Group: AfterImplementSpec, File: "AGENTS.md", Kind: Forbid, Needle: handoffAncestryRefusal, Diagnostic: "fact owner: working agreement restored its copy of the handoff State ancestry refusal"},
	{Group: AfterImplementSpec, File: "AGENTS.md", Kind: Forbid, Needle: "one gate-priced commit", Diagnostic: "fact owner: working agreement restored the gate-priced phase-close capture commit"},
	{Group: AfterImplementSpec, File: "AGENTS.md", Kind: Forbid, Needle: "Give the drafted State to `bench handoff --state-file <path>`.", Diagnostic: "fact owner: working agreement restored its copy of the handoff state-file route"},
	{Group: AfterImplementSpec, File: "projects/benchkit.md", Kind: Forbid, Needle: redRunGateOutput, Diagnostic: "fact owner: project profile restored its copy of the red-run gate output that the reference owns"},
	{Group: AfterImplementSpec, File: "projects/benchkit.md", Kind: Forbid, Needle: "one `capability-skips` line", Diagnostic: "fact owner: project profile restored its copy of the green-run capability-skips line that the reference owns"},
	{Group: AfterImplementSpec, File: "projects/benchkit.md", Kind: Forbid, Needle: "`phases: N/N green`", Diagnostic: "fact owner: project profile restored its copy of the collapsed green-run row that the reference owns"},
	{Group: AfterImplementSpec, File: "projects/benchkit.md", Kind: Forbid, Needle: "one more row names the file that holds the complete stream", Diagnostic: "fact owner: project profile restored its copy of the red-run stream row that the reference owns"},
	{Group: AfterImplementSpec, File: "projects/benchkit.md", Kind: Forbid, Needle: "The complete phase stream goes to `.logs/gate-<run>.out`", Diagnostic: "fact owner: project profile restored its copy of the gate stream log that the reference owns"},
	{Group: AfterImplementSpec, File: "projects/benchkit.md", Kind: Forbid, Needle: "`.logs/gate-<run>.jsonl` progress log", Diagnostic: "fact owner: project profile restored its copy of the gate progress log that the reference owns"},
	{Group: AfterImplementSpec, File: "projects/benchkit.md", Kind: Forbid, Needle: "`bench commit` works on any branch", Diagnostic: "fact owner: project profile restored the false claim that bench commit works on any branch"},
	{Group: AfterImplementSpec, File: ".bench/BENCH-reference.md", Kind: Forbid, Needle: "`bench handoff [--harness <name>]", Diagnostic: "fact owner: reference restored the bench handoff grammar that executable help owns"},
	{Group: AfterImplementSpec, File: ".bench/BENCH-reference.md", Kind: Forbid, Needle: "`bench worktree reset --to <commit> <target>` plans", Diagnostic: "fact owner: reference restored the bench worktree reset grammar that executable help owns"},
	{Group: AfterImplementSpec, File: ".bench/BENCH-reference.md", Kind: Forbid, Needle: "`bench retro <slug> (--body <markdown>", Diagnostic: "fact owner: reference restored the bench retro grammar that executable help owns"},
	{Group: AfterImplementSpec, File: ".bench/BENCH-reference.md", Kind: Forbid, Needle: "`.claude/skills/` carries only the", Diagnostic: "fact owner: reference restored the false claim that .claude/skills/ links only the craft skills"},
	{Group: AfterImplementSpec, File: ".bench/BENCH-reference.md", Kind: Forbid, Needle: "because Claude already has each phase as a command", Diagnostic: "fact owner: reference restored its copy of the Codex-only phase adapter rule that the Claude README owns"},
	{Group: AfterImplementSpec, File: ".bench/BENCH-reference.md", Kind: Forbid, Needle: skillLinkRule, Diagnostic: "fact owner: reference restored a copy of the Claude README's skill-link rule"},
	{Group: AfterImplementSpec, File: ".claude/README.md", Kind: Forbid, Needle: "links only the `bench-craft-*` skills", Diagnostic: "fact owner: Claude README restored the false claim that .claude/skills/ links only the craft skills"},
	{Group: AfterImplementSpec, File: ".claude/README.md", Kind: Forbid, Needle: "`bench-writer` runs a user-directed write delegation", Diagnostic: "fact owner: Claude README restored bench-writer as the type of a user-directed write delegation only"},
}

// declaredLineAnchors keep the retired author-line routes out of craft-line and
// craft-spec. Forbid rows keep out the per-story ceiling, the per-ticket
// re-run of the decision table, the per-story collapse, the retired step 3 and
// step 5 wording, the unscoped top-tier rule, and craft-spec's per-story lines in
// its approval table.
var declaredLineAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-line/SKILL.md", Kind: Forbid, Needle: "ceiling, not a binding", Diagnostic: "declared line: craft-line restored the per-story line as a ceiling, not a binding"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-line/SKILL.md", Kind: Forbid, Needle: "Re-run the decision table per ticket at charge time.", Diagnostic: "declared line: craft-line restored the per-ticket re-run of the decision table"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-line/SKILL.md", Kind: Forbid, Needle: "use the highest tier any story needs", Diagnostic: "declared line: craft-line restored the per-story collapse to the highest tier"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-line/SKILL.md", Kind: Forbid, Needle: "Report each collapsed line.", Diagnostic: "declared line: craft-line restored the report of each collapsed per-story line"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-line/SKILL.md", Kind: Forbid, Needle: "escalate immediately; no retry burned.", Diagnostic: "declared line: craft-line restored the step 3 escalation without the step 2 tier-move rule"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-line/SKILL.md", Kind: Forbid, Needle: "The top tier implements nothing unless the reviewer names it.", Diagnostic: "declared line: craft-line restored the top-tier rule for every run, `--delegate` included"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-line/SKILL.md", Kind: Forbid, Needle: "Any bump to the top tier pauses and asks the reviewer", Diagnostic: "declared line: craft-line restored the top-tier pause for every run, `--delegate` included"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-spec/SKILL.md", Kind: Forbid, Needle: "stories and their lines", Diagnostic: "declared line: craft-spec restored per-story lines in the approval table"},
}

// retroCaptureAnchors pin how the final check writes and commits the retro. A Require
// row pins the scaffold-then-body route through `bench retro`. Forbid rows keep out the
// report's copy of the
// retired no-commit rule, the drain-only exit, the reference's claim that the drain
// owns the retro's capture commit, the reference's copy of the commit rule, and a
// paste of the scaffold's heading list, which the retros package owns.
var retroCaptureAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Require, Class: Command, Needle: "Read `bench retro <slug> --scaffold`, then write the retro once with `bench retro <slug> --body <markdown>`.", Diagnostic: "retro capture: final check dropped the scaffold-then-body route through bench retro"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: "Capture the retro without another", Diagnostic: "retro capture: final check report restored the no-commit retro rule"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: "The retro leaves through the next reviewer-approved capture drain.", Diagnostic: "retro capture: final check restored the drain-only retro exit"},
	{Group: AfterImplementSpec, File: ".bench/BENCH-reference.md", Kind: Forbid, Needle: "owns their reviewed drain and its capture commit", Diagnostic: "retro capture: reference restored the drain as the owner of the retro capture commit"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: strings.Join(retros.RequiredHeadings(), "\n"), Diagnostic: "retro capture: final check restored a copy of the scaffold's retro heading list"},
	{Group: AfterImplementSpec, File: ".bench/BENCH-reference.md", Kind: Forbid, Needle: "A tracked retro commits with the phase close", Diagnostic: "retro capture: reference restored its copy of the final check's tracked-or-ignored retro rule"},
}

// laneAndLandingAnchors guard the commit and landing split. Forbid rows keep out the
// retired claims that the commit is the gate run or the source of the green verdict,
// the second gate before the landing, and the reference's claims that the landing
// boundary is unenforced and that a stale executable reruns the landing. A Forbid row
// reads the operating guide's enforcement sentence through its named constant, so no
// copy returns there. Require rows pin the synthesis gate source, the final check's
// broker rehearsal, and its spec retirement.
var laneAndLandingAnchors = []Anchor{
	{Group: AfterImplementSpec, File: synthesisSkill, Kind: Forbid, Needle: "`bench commit` gates the tree it lands", Diagnostic: "lane and landing: craft-synthesis restored the claim that bench commit gates the tree it lands"},
	{Group: AfterImplementSpec, File: synthesisSkill, Kind: Forbid, Needle: "Take that verdict from the commit itself.", Diagnostic: "lane and landing: craft-synthesis restored the commit as the source of the prose-only green verdict"},
	{Group: AfterImplementSpec, File: synthesisSkill, Kind: Forbid, Needle: "A `bench gate` run before it grades the same tree twice and pays the oracle twice.", Diagnostic: "lane and landing: craft-synthesis restored the claim that a gate run before the commit pays the oracle twice"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: "Run the external gate and commit work on green", Diagnostic: "lane and landing: final check description restored the claim that it runs the gate and commits on green"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: "After a lane-only repair commit and before the landing, run the whole-tree gate on the source.", Diagnostic: "lane and landing: final check restored a second whole-project gate before the landing"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: "Do not run `bench gate` first.", Diagnostic: "lane and landing: final check restored the ban on a gate run before the commit"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: "the gate reuses a fresh green verdict for the identical tree and never re-pays it", Diagnostic: "lane and landing: final check restored the claim that the commit reuses a fresh green gate verdict"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: "This command gates and commits them atomically.", Diagnostic: "lane and landing: final check restored the claim that bench commit gates and commits atomically"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: "The commit already is the gate run", Diagnostic: "lane and landing: final check restored the claim that the commit is the gate run"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: "the oracle run and landing are one command", Diagnostic: "lane and landing: final check restored the claim that the oracle run and the landing are one command"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: "gate-then-commit path", Diagnostic: "lane and landing: final check restored the gate-then-commit path"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: "then runs the gate and commits only on green", Diagnostic: "lane and landing: final check restored the claim that bench commit runs the gate and commits only on green"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Require, Class: Command, Needle: "Before the first landing of such a spec, run `bench worktree build <target>` and then `bench worktree exec <target> -- ./dist/bench doctor`.", Diagnostic: "lane and landing: final check dropped the broker rehearsal before the first landing of a broker-changing spec"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Require, Class: Command, Needle: "A merged spec awaiting retirement gets `bench spec retire <slug>` in a Bench worktree, and that worktree lands its `spec-retire: <slug>` commit through `bench worktree land`.", Diagnostic: "lane and landing: final check dropped the spec retirement in a Bench worktree and its landing"},
	{Group: AfterImplementSpec, File: ".bench/BENCH-reference.md", Kind: Forbid, Needle: "the rule is guidance, not a hook", Diagnostic: "lane and landing: reference restored the claim that the landing rule is guidance, not a hook"},
	{Group: AfterImplementSpec, File: ".bench/BENCH-reference.md", Kind: Forbid, Needle: WorktreeEnforcementMarker, Diagnostic: "lane and landing: reference restored a copy of the operating guide's worktree enforcement sentence"},
	{Group: AfterImplementSpec, File: ".bench/BENCH-reference.md", Kind: Forbid, Needle: "A stale Bench executable is rebuilt, and the landing re-runs under it.", Diagnostic: "lane and landing: reference restored the claim that the landing rebuilds a stale Bench executable and re-runs"},
}

// defaultWorkflowAnchors keep out the plan-expansion learning without a dispatch
// deadline and the Writes expectation scoped to a user-directed write delegate.
var defaultWorkflowAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".bench/BENCH.md", Kind: Forbid, Needle: "Record every plan or gate expansion with `bench learning`, including what changed, why, and verification; `/bench-drain` owns its later disposition.", Diagnostic: "retained workflow: operating guide restored expansion learning without a dispatch deadline"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/references/delegation-discipline.md", Kind: Forbid, Needle: "A user-directed write delegate treats `Writes:` as an expectation.", Diagnostic: "retained workflow: delegation discipline scoped the Writes expectation to a user-directed write delegate"},
}

// delegatedWorkflowAnchors keep out the retired rules of who authors, repairs, probes,
// commits, and lands a ticket: the retained-session repair, the user-directed-only
// bench-writer role, the coordinator commit, the diff-ready stop, the copy-aside
// probes, the main-checkout build, the resume-clean step, and phase-command copies of
// owner rules. A Forbid row that keeps a copy of an owner sentence out reads that
// sentence's named constant. No anchor claims to prove native dispatch.
var delegatedWorkflowAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/SKILL.md", Kind: Forbid, Needle: "A spec-backed ticket goes to a fresh author session on its integration source", Diagnostic: "fresh ticket author: craft-delegate restored its copy of the fresh author session rule that the operating guide owns"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/SKILL.md", Kind: Forbid, Needle: "Repairs return to the retained implementation session.", Diagnostic: "fresh ticket author: craft-delegate restored the repair return to the retained session"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/SKILL.md", Kind: Forbid, Needle: "delegation runs as `bench-writer`", Diagnostic: "fresh ticket author: craft-delegate restored bench-writer as the type of a user-directed write delegation only"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/SKILL.md", Kind: Forbid, Needle: "A user-directed write-delegation from a spec carries", Diagnostic: "fresh ticket author: craft-delegate scoped the coverage-row charge to a user-directed write-delegation"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/SKILL.md", Kind: Forbid, Needle: "the coordinator runs `bench commit` per worktree", Diagnostic: "fresh ticket author: craft-delegate restored the coordinator commit per worktree"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/SKILL.md", Kind: Forbid, Needle: "stops at diff-ready", Diagnostic: "fresh ticket author: craft-delegate restored the diff-ready stop for every write delegate"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/SKILL.md", Kind: Forbid, Needle: "Stop at diff ready;", Diagnostic: "fresh ticket author: craft-delegate charge example restored the uncommitted diff-ready stop"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/SKILL.md", Kind: Forbid, Needle: "Share a worktree only when a delegate's work depends on another's output.", Diagnostic: "fresh ticket author: craft-delegate restored the dependency-only worktree share rule"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/SKILL.md", Kind: Forbid, Needle: "A ticket delegate returns focused evidence", Diagnostic: "fresh ticket author: craft-delegate restored the returned-tree coordinator probe for every ticket delegate"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/SKILL.md", Kind: Forbid, Needle: "the destructive-git guard refuses it", Diagnostic: "fresh ticket author: craft-delegate restored the false claim that the guard refuses every git stash"},
	{Group: AfterImplementSpec, File: ".claude/agents/bench-writer.md", Kind: Forbid, Needle: "The user directed this delegation", Diagnostic: "fresh ticket author: bench-writer restored the user-directed premise for every write delegate"},
	{Group: AfterImplementSpec, File: ".claude/agents/bench-writer.md", Kind: Forbid, Needle: "Stop at a diff that is ready, with the focused checks green.", Diagnostic: "fresh ticket author: bench-writer restored the uncommitted-diff stop for every write delegate"},
	{Group: AfterImplementSpec, File: ".bench/BENCH.md", Kind: Forbid, Needle: "Production repairs stay with the recorded ticket author.", Diagnostic: "fresh ticket author: operating guide restored the ticket-author repair ownership"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/references/delegation-discipline.md", Kind: Forbid, Needle: "Probe a tracked file that has pending changes with a copy aside.", Diagnostic: "fresh ticket author: delegation discipline restored the copy-aside probe that bench probe replaces"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/references/delegation-discipline.md", Kind: Forbid, Needle: "mutated bytes against the copy aside", Diagnostic: "fresh ticket author: delegation discipline restored the copy-aside confirmation of the mutated bytes"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/references/delegation-discipline.md", Kind: Forbid, Needle: "may run in the main checkout", Diagnostic: "fresh ticket author: delegation discipline restored the main-checkout build exception"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/references/delegation-discipline.md", Kind: Forbid, Needle: "`bench resume-clean`", Diagnostic: "fresh ticket author: delegation discipline restored the internal resume-clean verb as a recovery step"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/references/delegation-discipline.md", Kind: Forbid, Needle: "on the integration source after every ticket commit and before the next charge", Diagnostic: "fresh ticket author: delegation discipline restored its copy of the build preflight rule that the build phase owns"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/references/delegation-discipline.md", Kind: Forbid, Needle: "a slow tool call is not a failed attempt", Diagnostic: "fresh ticket author: delegation discipline restored its copy of the attempt definition that craft-line owns"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Forbid, Needle: "A delegated chunk review starts after every ticket of the chunk reaches the integrated chunk tip.", Diagnostic: "fresh ticket author: review phase restored its copy of the chunk-review start rule that the operating guide owns"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Forbid, Needle: "Repeat delegated review only for a later semantic delta", Diagnostic: "fresh ticket author: review phase restored its copy of the review-repeat rule that the operating guide owns"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Forbid, Needle: "After the last chunk, the orchestrator reconciles overall acceptance and integration before landing.", Diagnostic: "fresh ticket author: review phase restored its copy of the final reconciliation rule that the operating guide owns"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Forbid, Needle: "Repeat delegated review only when a later delta or cross-chunk concern invalidates prior evidence.", Diagnostic: "fresh ticket author: implementation phase restored its copy of the review-repeat rule that the operating guide owns"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Forbid, Needle: reviewRepeatRule, Diagnostic: "fresh ticket author: review phase restored a copy of the operating guide's review-repeat rule"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Forbid, Needle: finalReconciliationRule, Diagnostic: "fresh ticket author: review phase restored a copy of the operating guide's final reconciliation rule"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Forbid, Needle: chunkTipReviewFence, Diagnostic: "fresh ticket author: review phase restored a copy of the operating guide's chunk-tip review fence"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Forbid, Needle: reviewRepeatRule, Diagnostic: "fresh ticket author: implementation phase restored a copy of the operating guide's review-repeat rule"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Forbid, Needle: "For a ticket author, raise the effort and resume; a tier move asks the reviewer first.", Diagnostic: "fresh ticket author: implementation phase restored its copy of the tier ladder without the `--delegate` exception that craft-line owns"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-debug.md", Kind: Forbid, Needle: "The coordinator validates the report and reslices", Diagnostic: "fresh ticket author: debug phase restored the coordinator reslice of an out-of-fence cause"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: "the author's final acceptance and integration command results", Diagnostic: "fresh ticket author: final check restored the author as the performer of the final acceptance and integration results"},
}

// chunkTipReviewFence is the `.bench/BENCH.md` owner sentence for the chunk-review
// start. The Forbid row that keeps a review-phase copy out reads this one constant.
const chunkTipReviewFence = "Every ticket contribution reaches the integrated chunk tip before that chunk's review begins."

// implementationContinuationAnchors keep out the craft-delegate copy of the diagnostic
// helper boundary and the bounded repair policy's retired mode list.
var implementationContinuationAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/SKILL.md", Kind: Forbid, Needle: "A diagnostic helper can inspect evidence, but it receives no implementation or repair assignment.", Diagnostic: "implementation continuation: craft-delegate restored its copy of the diagnostic helper boundary that the operating guide owns"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", Kind: Forbid, Needle: "This allowance applies to retained, full, delegated, unattended, and light-path implementation runs.", Diagnostic: "implementation continuation: bounded repair restored the retired mode list"},
}
