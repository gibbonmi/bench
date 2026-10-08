package refusalroute

// destinationClean is the repair a destination that carries uncommitted work demands. The
// first run and the resume refuse on the same destination state, so both faces name it.
var destinationClean = Instruction(Text("commit the destination's uncommitted work, or discard it"))

// review sends a repaired source back through review before the landing re-runs. The
// review is a phase the agent runs, not a shell command, so the step is an instruction.
var review = Instruction(Text("/bench-review-implementation"))

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
			Instruction(Text("take the refusal_paths entries out of the reviewed range, or declare them under the spec's ## Ownership fences")),
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
			Instruction(Text("repair each failure that the gate reports in"), Fact(FactLabel)),
			commitAt(FactLabel),
			review,
			rerun,
		},
	},
	{
		// The gate stopped on an infrastructure outcome, so no failure is the diff's. The
		// authorization policy owns the sentence.
		Verb:      Land,
		Name:      "land-infrastructure",
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
