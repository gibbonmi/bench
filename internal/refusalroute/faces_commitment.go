package refusalroute

// The commitment faces. The commitment policy and the commitment verb raise them from
// several packages, so each name has this one spelling.
const (
	CommitmentNeedsAssignment  = "commitment-needs-assignment"
	CommitmentVerifyEvidence   = "commitment-verify-evidence"
	CommitmentDecision         = "commitment-decision"
	CommitmentUnbound          = "commitment-unbound"
	CommitmentLightPathOutside = "commitment-light-path-outside"
	CommitmentRunUnknown       = "commitment-run-unknown"
	CommitmentPlanInput        = "commitment-plan-input"
	CommitmentHandback         = "commitment-handback"
)

// commitmentFaces are the commitment refusal faces, in registry order. The policy or the
// verb supplies each sentence: the refusal names the cause it observed.
var commitmentFaces = []Face{
	{
		// The policy decides for the one active assignment that owns the checkout, and the
		// agent makes its own.
		Verb:      Commitment,
		Name:      CommitmentNeedsAssignment,
		Authority: Agent,
		Route:     []Step{Command(Text("bench worktree create --request"), Operator("request"), Text("--label"), Operator("label"))},
	},
	{
		// A milestone completion needs the receipt that a verification of its evidence records.
		Verb:      Commitment,
		Name:      CommitmentVerifyEvidence,
		Authority: Agent,
		Route:     []Step{Command(Text("bench commitment verify --milestone"), Fact(FactMilestone), Text("--evidence"), Operator("file"))},
	},
	{
		// The clear changes the active commitment: its adoption, an exact approval, the
		// protected commitment, a legacy scope, or the obligations of a deliverable. A
		// commitment change is the reviewer's.
		Verb:      Commitment,
		Name:      CommitmentDecision,
		Authority: Reviewer,
		Route:     []Step{Command(Text("bench commitment plan --input"), Operator("file"))},
	},
	{
		// The assignment holds no current delivery binding, so the agent binds its own
		// assignment to the eligible outcome.
		Verb:      Commitment,
		Name:      CommitmentUnbound,
		Authority: Agent,
		Route:     []Step{Command(Text("bench commitment start --outcome"), Operator("id"), Text("--request"), Operator("request"), Text("--deliverable"), Operator("path"))},
	},
	{
		// A light-path change carries one ticket, so a path outside its Writes line is one
		// edit of that line. The policy reads the ticket from the candidate, so the printing
		// verb's own re-run names the ticket after its paths.
		Verb:      Commitment,
		Name:      CommitmentLightPathOutside,
		Authority: Agent,
		Route:     []Step{Instruction(Text("add the path to the Writes: line of"), Fact(FactTicket)), Command(Composed(FactRerun), Fact(FactTicket))},
	},
	{
		// A plan lists a run that the ledger does not hold, and the inventory lists the runs
		// it holds.
		Verb:      Commitment,
		Name:      CommitmentRunUnknown,
		Authority: Agent,
		Route:     []Step{Command(Text("bench commitment inventory"))},
	},
	{
		// The plan input does not read or does not validate, and the caller corrects it.
		Verb:      Commitment,
		Name:      CommitmentPlanInput,
		Authority: Agent,
		Route:     []Step{Instruction(Text("correct the input file")), rerun},
	},
	{
		// A cause with no face of its own, such as a board that the default branch cannot
		// read. Its clear is outside the agent's authority.
		Verb:      Commitment,
		Name:      CommitmentHandback,
		Authority: Reviewer,
		Route:     handback,
	},
}

// CommitRerun renders the caller's own commit re-run over values: the label and the
// arguments that the commit composed. A commitment face whose route ends with the printing
// verb's re-run reads it, so the commit supplies its own in the form its own faces print.
func CommitRerun(values map[string]string) string { return commitRerun.render(factsFill(values)) }
