package refusalroute

// A policy outside the landing raises these land faces. The landing lease names its owner,
// and a lease lock fault stops the landing as an infrastructure outcome. The build window
// names the open build that holds the destination.
const (
	LandLeaseHeld      = "land-lease-held"
	LandInfrastructure = "land-infrastructure"
	LandBuildWindow    = "land-build-window"
)

// destinationClean is the repair a destination that carries uncommitted work demands. The
// first run and the resume refuse on the same destination state, so both faces name it.
var destinationClean = Instruction(Text("commit the destination's uncommitted work, or discard it"))

// review sends a repaired source back through review before the landing re-runs. The
// review is a phase the agent runs, not a shell command, so the step is an instruction.
var review = Instruction(Text("/bench-review-implementation"))

// repairRed sends the repair of a red grade to the source that the label fact addresses. A
// single landing and a stack both name it.
var repairRed = Instruction(Text("repair each failure that the gate reports in"), Fact(FactLabel))

// landFaces are the landing's refusal faces, in registry order.
var landFaces = []Face{
	{
		Verb:      Land,
		Name:      "destination-not-clean",
		Sentence:  "landing destination is not clean",
		Authority: Reviewer,
		Route:     []Step{destinationClean, rerun},
	},
	{
		// The paths are the operator's own files where the landing writes. Git refuses to
		// overwrite an untracked one and overwrites an ignored one without a word, so the
		// operator moves them before the landing runs.
		Verb:      Land,
		Name:      "destination-collision",
		Sentence:  "landing destination has untracked or ignored files where the landing writes",
		Authority: Reviewer,
		Route:     []Step{Instruction(Text("move the refusal_paths entries out of the landing checkout")), rerun},
	},
	{
		// The raising site re-points the caller's re-run at the source tip the tree holds,
		// so a moved tip leaves exactly one command to run.
		Verb:      Land,
		Name:      "source-tip-mismatch",
		Sentence:  "worktree source tip mismatch",
		Authority: Agent,
		Route:     []Step{rerun},
	},
	{
		// The commit moves the reviewed source, so the source goes back through review, and
		// the raising site points the re-run at the repaired tip.
		Verb:      Land,
		Name:      "source-not-clean",
		Sentence:  "reviewed source is not clean",
		Authority: Agent,
		Route:     []Step{commitAt(FactLabel), review, rerun},
	},
	{
		// The fence already authorizes a path that a fold of the default branch brings in
		// unchanged, so every refused path is a write of the build, and the route keeps the
		// caller's --base.
		Verb:      Land,
		Name:      "source-not-fenced",
		Sentence:  "reviewed source range or ownership fence is invalid",
		Authority: Agent,
		Route: []Step{
			Instruction(Text("take the refusal_paths entries out of the reviewed range, or add them to a ticket's Writes: line with bench spec widen")),
			rerun,
		},
	},
	{
		// The lease policy owns the sentence, which names the owner's assignment. The owner's
		// landing releases the lease when it publishes or refuses. The id is a hex identity,
		// so the route prints it as the sentence does.
		Verb:      Land,
		Name:      LandLeaseHeld,
		Authority: Agent,
		Route:     []Step{Instruction(Text("wait until the landing of assignment"), Composed(FactLeaseOwner), Text("ends")), rerun},
	},
	{
		// The window policy owns the sentence, which names the build and the gap to wait for.
		// The build's gap, or a landing override for this assignment, lets the re-run through.
		Verb:      Land,
		Name:      LandBuildWindow,
		Authority: Agent,
		Route:     []Step{Instruction(Text("wait for the gap that the refusal names, or ask the reviewer for an override")), rerun},
	},
	{
		// The fold brings the moved destination into the source. The raising site points the
		// re-run at the repaired tip.
		Verb:      Land,
		Name:      "land-destination-moved",
		Sentence:  "landing destination moved past the reviewed source",
		Authority: Agent,
		Route: []Step{
			Command(Text("bench worktree merge --from"), Fact(FactFrom), Fact(FactLabel)),
			Instruction(Text("for a spec landing, record the fold, its review, the merged-test verification, and the completion at the new tip")),
			rerun,
		},
	},
	{
		// The composition names the conflicted paths in its own sentence, so the face
		// declares none.
		Verb:      Land,
		Name:      "composition-conflict",
		Authority: Reviewer,
		Route:     append(append([]Step{}, handMerge...), review, rerun),
	},
	{
		// The source already holds a merge in progress, so a second merge is not the
		// repair: finishing the merge records the resolution and needs no separate commit.
		Verb:      Land,
		Name:      "composition-conflict-pending",
		Authority: Reviewer,
		Route: []Step{
			Command(Composed(FactCheckoutGit), Text("merge --continue (resolve the conflicted paths of the merge in progress first)")),
			review,
			rerun,
		},
	},
	{
		// The residue policy owns the sentence this face prints.
		Verb:      Land,
		Name:      "resume-destination-residue",
		Authority: Reviewer,
		Route:     []Step{destinationClean, rerun},
	},
	{
		// The landing policy owns the sentence this face prints. Only a landing on the
		// primary checkout advances the green marker, so the repair is the reviewer's.
		Verb:      Land,
		Name:      "resume-marker",
		Authority: Reviewer,
		Route: []Step{
			Instruction(Text("land a green landing on main that covers the published commit, or restore main to it")),
			rerun,
		},
	},
	{
		// The gate graded the composed tree red. The authorization policy owns the sentence.
		// The repair commits in the source, so the source goes back through review, and the
		// raising site points the re-run at the repaired tip.
		Verb:      Land,
		Name:      "land-red",
		Authority: Agent,
		Route: []Step{
			repairRed,
			commitAt(FactLabel),
			review,
			rerun,
		},
	},
	{
		// A stack's bisect attributed the red to one member. The authorization policy owns
		// the sentence. The repair commits in that member, so it goes back through review.
		// A published prefix moves the destination, so the fold follows the review, and the
		// raising site points the re-run at the member's tip after the fold.
		Verb:      Land,
		Name:      "stack-red",
		Authority: Agent,
		Route: []Step{
			repairRed,
			commitAt(FactLabel),
			review,
			Instruction(Text("when the members table shows a landed member, remove each landed member from the stack file, fold the destination into each remaining member, and set the new source_tip of each folded member in the stack file")),
			rerun,
		},
	},
	{
		// The member's landing closes its deliverable or transitions its staged spec, which a
		// stack does not do, so the member leaves the stack and takes its own landing.
		Verb:      Land,
		Name:      "stack-member-alone",
		Sentence:  "stack member must land alone",
		Authority: Agent,
		Route: []Step{
			Instruction(Text("remove the member"), Fact(FactLabel), Text("from the stack file, and land it with its own landing")),
			rerun,
		},
	},
	{
		// Two members change one path outside capture/, so the later member of the pair leaves
		// the stack and lands later on its own route.
		Verb:      Land,
		Name:      "stack-members-overlap",
		Sentence:  "stack members change one path",
		Authority: Agent,
		Route: []Step{
			Instruction(Text("remove the member"), Fact(FactLabel), Text("from the stack file; it lands later on its own route")),
			rerun,
		},
	},
	{
		// The gate stopped on an infrastructure outcome, so no failure is the diff's. The
		// authorization policy owns the sentence.
		Verb:      Land,
		Name:      LandInfrastructure,
		Authority: Agent,
		Route:     []Step{doctor, rerun},
	},
	{
		// The landing published, and a later step did not finish. The re-run is the
		// caller's resume of that published landing.
		Verb:      Land,
		Name:      "land-incomplete",
		Authority: Agent,
		Route:     []Step{rerun},
	},
	{
		Verb:      Land,
		Name:      "land-handback",
		Authority: Reviewer,
		Route:     handback,
	},
}
