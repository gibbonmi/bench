package recordcmd

import (
	"errors"
	"strconv"

	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/usage"
)

var assignmentForm = form{name: reviewrecord.FormAssignment, description: "append one author dispatch to a ticket's history in the version 2 completion plan at HEAD",
	flags: []flag{{"--source", "<commit>", false, once}, {"--ticket", "<basename>", true, once}, {"--performer", "<session>", true, once},
		{"--assignment", "<id>", true, once}, {"--model", "<model>", true, once}, {"--effort", "<effort>", true, once},
		{"--ref", "<native-dispatch-ref>", true, once}, {"--run", "<run-id>", true, once}, {"--orchestrator", "<session>", true, once},
		{"--author-limit", "<positive>", false, once}, {"--review-mode", "unified", false, once}, {"--predecessor", "<session>", true, once},
		{"--trigger", "<trigger>", true, once}, {"--stopped", "<ref>", true, once}, {"--preserved", "<ref>", true, once},
		{"--reassessment", "<ref>", true, once}},
	layout: []string{"--source", "--ticket", "--performer", "--assignment", "--model", "--effort", "--ref", "--run", "--orchestrator", "--author-limit",
		"[", "--review-mode", "]", "[", "--predecessor", "--trigger", "--stopped", "--preserved", "[", "--reassessment", "]", "]"},
	valid: assignmentValid, run: assignment}

// assignmentValid closes the --author-limit and --review-mode value sets.
func assignmentValid(f form, parsed usage.Result) bool {
	limit, ok := integer(parsed.Flags, "--author-limit")
	_, mode := parsed.Flags["--review-mode"]
	return ok && limit > 0 && (!mode || chosen(f, parsed.Flags, "--review-mode"))
}

// assignment records one dispatch. The first dispatch of a chunk in an open build closes the
// build's gap in the ledger before the plan write. So a refused dispatch writes no plan.
func assignment(f form, root, spec string, parsed usage.Result) (string, int) {
	values := parsed.Flags
	source, out, code := commit(f, root, values, "--source")
	if out != "" {
		return out, code
	}
	limit, _ := integer(values, "--author-limit")
	call := reviewrecord.AssignmentCall{Source: source, Ticket: values["--ticket"], RunID: values["--run"], OrchestratorSession: values["--orchestrator"],
		AuthorLimit: limit, ReviewMode: values["--review-mode"], Assignment: reviewrecord.Assignment{Session: values["--performer"],
			Assignment: values["--assignment"], Model: values["--model"], Effort: values["--effort"], NativeRef: values["--ref"],
			Predecessor: values["--predecessor"], Trigger: values["--trigger"], Stopped: values["--stopped"], Preserved: values["--preserved"],
			Reassessment: values["--reassessment"]}}
	call.BeforeWrite = func(plan reviewrecord.Plan) error { return commitrepo.Store{Root: root}.CloseGap(plan, call.Ticket) }
	history, err := reviewrecord.RecordAssignment(root, spec, call)
	var unfolded commitrepo.UnfoldedDefault
	if errors.As(err, &unfolded) {
		return f.refuse(foldFirst(spec, unfolded))
	}
	if err != nil {
		return f.refuse(err.Error())
	}
	last := history[len(history)-1]
	return row("assignment", []string{"ticket", "session", "occurrence", "predecessor", "source"},
		[]string{call.Ticket, last.Session, strconv.Itoa(len(history)), last.Predecessor, last.Source})
}

// foldFirst is the refusal of a first chunk dispatch over a source that lacks the
// default-branch tip. It names the fold and the fold record that let the dispatch through.
func foldFirst(spec string, unfolded commitrepo.UnfoldedDefault) string {
	slug, _ := reviewrecord.Slug(spec)
	return "dispatch refused: " + unfolded.Error() + "; fold it with bench worktree merge --from " + unfolded.Branch + " " + unfolded.Assignment +
		" and record it with " + command + " " + reviewrecord.FormFold + " " + slug + " --source <fold-commit> before the dispatch"
}
