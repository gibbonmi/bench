package refusalroute

import (
	"strings"

	"github.com/gibbonmi/bench/internal/usage"
)

// gateRerun re-runs the caller's own gate. The raising site composes the arguments from the
// run it refuses, and a face whose recovery forces a fresh run gets them with --fresh.
var gateRerun = callerRerun("bench gate")

// GateRerun renders the caller's own gate re-run over values.
func GateRerun(values map[string]string) string { return renderRerun(gateRerun, values) }

// FactForm is the record form that writes the entry a checkpoint refusal lacks.
const FactForm = "form"

// CandidateNotInstalled is the face of a declared build's guarded write whose running
// executable is not the installed candidate. The candidate check raises it, and the gate
// checkpoint and the record verb print it. The build verb's restore prints its route when the
// running executable cannot read the ledger.
const CandidateNotInstalled = "candidate-not-installed"

// DeclaredBuildsOverlap is the face of a guarded write while more than one declared build is
// open. The candidate check raises it, and the gate checkpoint and the record verb print it.
// The build verb's restore prints its route while the ledger holds an open build of any
// outcome.
const DeclaredBuildsOverlap = "declared-builds-overlap"

// installForm is the install form of the build verb up to its target slot, which the target
// fact fills.
var installForm = strings.TrimSuffix(usage.WorktreeBuildInstall, " "+placeholder(FactTarget))

// gateFaces are the gate checkpoint's refusal faces, in registry order. The raising site
// supplies each sentence: the checkpoint names the cause it observed.
var gateFaces = []Face{
	{
		// The checkpoint cannot show what the completion evidence lacks, and one read
		// reports the plan row and the record state. The slug comes from the spec-path
		// grammar's one owner.
		Verb:      Gate,
		Name:      "checkpoint-completion-evidence",
		Authority: Agent,
		Route:     []Step{Command(Text("bench preflight review"), Fact(FactSlug))},
	},
	{
		// The completion record lacks an entry that one record form writes, and the
		// refused cause names that form. The form's help prints its grammar.
		Verb:      Gate,
		Name:      "checkpoint-record-repair",
		Authority: Agent,
		Route:     []Step{Command(Text("bench record"), Fact(FactForm), Text("--help"))},
	},
	{
		// The landing publishes only committed bytes, and the chunk checkpoint grades only
		// committed bytes. So the caller commits its own uncommitted work and grades the
		// committed source again.
		Verb:      Gate,
		Name:      "checkpoint-dirty-checkout",
		Authority: Agent,
		Route:     []Step{commitAt(FactLabel), gateRerun},
	},
	{
		// The landing's closure transform refuses the spec: its status or its delivery
		// closure is the reviewer's to change.
		Verb:      Gate,
		Name:      "checkpoint-composition",
		Authority: Reviewer,
		Route:     []Step{Instruction(Text("the delivery closure of the spec does not compose; hand back"))},
	},
	{
		// The source tip moved between the accepted subject and the run, so the rerun
		// grades the moved tip.
		Verb:      Gate,
		Name:      "checkpoint-tip-moved",
		Authority: Agent,
		Route:     []Step{gateRerun},
	},
	{
		// A fault captured no subject. The diagnosis comes first, and the rerun does not
		// reuse a verdict.
		Verb:      Gate,
		Name:      "checkpoint-subject-unavailable",
		Authority: Agent,
		Route:     []Step{doctor, gateRerun},
	},
	{
		// The graded completion is not the exact transform of its reviewed source, which
		// no repair in the caller's worktree reaches.
		Verb:      Gate,
		Name:      "gate-handback",
		Authority: Reviewer,
		Route:     []Step{clearCause, gateRerun},
	},
	{
		// The install runs from the primary checkout, which no step addressed at a worktree
		// reaches, so it is an instruction. The printing verb supplies its own re-run.
		Verb:      Gate,
		Name:      CandidateNotInstalled,
		Authority: Agent,
		Route:     []Step{Instruction(Text("from the primary checkout, run"), Text(installForm), Fact(FactTarget)), rerun},
	},
	{
		// One primary executable serves one installed candidate, so one of the open declared
		// builds must end. The choice of the build to end is the reviewer's. The printing verb
		// supplies its own re-run.
		Verb:      Gate,
		Name:      DeclaredBuildsOverlap,
		Authority: Reviewer,
		Route:     []Step{Command(Text("bench commitment block --outcome"), Operator("id"), Text("--reason"), Operator("text")), rerun},
	},
}
