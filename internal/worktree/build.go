// The worktree build verb: one executable at <worktree>/dist/bench through the tree's own
// build script, and its install form, which publishes that executable at the primary checkout.
package worktree

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// BuildCommand builds one active Bench-owned worktree's tree into that worktree's
// `dist/bench`, so the assignment's own grammar has an executable. The build runs the
// tree's sanctioned build script, which writes the executable and its seal, so the verb
// authors no artifact of its own.
func BuildCommand(root, home string, args []string, stdout, stderr io.Writer) int {
	return buildWith(defaultJoins(), newAmbient(home, stderr), root, args, stdout, stderr)
}

// buildWith is BuildCommand with the seam set and the ambient value resolved at the
// caller's boundary. It is also the verb's record boundary. The record opens after the
// grammar answers, and the verb resolves its target inside the span, so a target refusal
// records the verb with no subject.
func buildWith(j joins, a ambient, root string, args []string, stdout, stderr io.Writer) int {
	parsed, line, code := usage.Parse(buildGrammar, args)
	if line != "" {
		if code == 0 {
			fmt.Fprintln(stdout, line)
			return 0
		}
		fmt.Fprintln(stderr, line)
		return code
	}
	var assignment string
	finishSpan := beginVerbSpan(a.home, root, otelBuildSeam)
	exit := buildAttributed(&assignment, j, root, parsed, stdout, stderr)
	finishSpan(exit, assignment)
	return exit
}

// buildAttributed is the build verb's own work, with the assignment the build wrote for
// written to attributed. The target authority is read once, so the record names the same
// assignment the executable landed in.
func buildAttributed(attributed *string, j joins, root string, parsed usage.Result, stdout, stderr io.Writer) int {
	if _, install := parsed.Flags[installFlag]; install {
		return installAttributed(attributed, j, root, parsed.Positionals[0], stdout, stderr)
	}
	path, assignment, err := resolveBuildTarget(root, parsed.Positionals[0])
	if err != nil {
		return printTargetRefusal(stderr, buildVerb, err)
	}
	*attributed = assignment.ID
	output := freshness.PublishedExecutable(path)
	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	if buildErr := j.buildSubject(ctx, path, output); buildErr != nil {
		fmt.Fprintf(stderr, "%s: %v\n", buildVerb, buildErr)
		return nameWorktree(stderr, path, buildExitCode(buildErr))
	}
	table, err := toon.Table(buildTable, []string{"worktree", "executable"}, [][]string{{assignment.ID, output}})
	if err != nil {
		fmt.Fprintln(stderr, toon.RenderError(err))
		return 1
	}
	fmt.Fprint(stdout, table)
	fmt.Fprintf(stdout, "next[1]:\n%s\n", buildExecNext(assignment))
	return 0
}

// buildVerb is the name every refusal this verb prints carries.
var buildVerb = usage.CommandName(usage.WorktreeBuild)

// installFlag selects the install form of the build verb.
const installFlag = "--install"

var buildGrammar = usage.Grammar{
	Cmd:     buildVerb,
	Help:    "usage: " + usage.WorktreeBuild + "\n       " + usage.WorktreeBuildInstall,
	Flags:   []usage.Flag{{Name: installFlag}},
	MinArgs: 1,
	MaxArgs: 1,
}

// resolveBuildTarget returns the authorized worktree path and the record the output
// names. resolveAssignment is the whole authority, and the record it answers carries both.
func resolveBuildTarget(root, target string) (string, intent.Assignment, error) {
	assignment, err := resolveAssignment(root, target)
	if err != nil {
		return "", intent.Assignment{}, err
	}
	return assignment.Worktree, assignment, nil
}

// buildExitCode maps an interrupted build apart from a failed one, so a reader can tell
// a signal from a broken tree.
func buildExitCode(err error) int {
	if errors.Is(err, context.Canceled) {
		return 130
	}
	return 1
}

// buildExecNext names the command that runs what the build produced. A label that is not
// line-safe gives way to the assignment id, by the conflictFacts precedent, because
// no quoting makes a control byte pasteable.
func buildExecNext(assignment intent.Assignment) string {
	address := assignment.ID
	if lineSafe(assignment.Label) {
		address = sanitize.ShellQuote(assignment.Label)
	}
	return "  bench worktree exec " + address + " -- ./dist/bench <verb>"
}

// installTarget is the authority of one install: the primary checkout, the bound
// assignment whose tree the install builds, and that assignment's bound deliverable.
type installTarget struct {
	primary     string
	assignment  intent.Assignment
	deliverable string
}

// installAttributed is the install form of the build verb. It builds the target's tree
// into the target's own `dist/bench`, publishes a copy as the primary executable, and runs
// the primary checkout's doctor repair, so the broker manifest binds the candidate. Each
// authority refusal comes before the build, so such a refusal writes no file. A failed
// repair leaves the published candidate in place, and a rerun repeats both steps.
func installAttributed(attributed *string, j joins, root, target string, stdout, stderr io.Writer) int {
	authority, err := installAuthority(j, root, target)
	if err != nil {
		return printTargetRefusal(stderr, buildVerb, err)
	}
	*attributed = authority.assignment.ID
	built := freshness.PublishedExecutable(authority.assignment.Worktree)
	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	if buildErr := j.buildSubject(ctx, authority.assignment.Worktree, built); buildErr != nil {
		fmt.Fprintf(stderr, "%s: %v\n", buildVerb, buildErr)
		return nameWorktree(stderr, authority.assignment.Worktree, buildExitCode(buildErr))
	}
	executable, err := publishInstall(j, authority, built)
	if err != nil {
		return printTargetRefusal(stderr, buildVerb, refusalError{refusal{detail: "publication refused: " + err.Error(), next: usage.WorktreeBuildInstallFor(targetAddress(target))}})
	}
	argv := []string{filepath.Join(wrapperDir(authority.primary), "bench.sh"), "doctor", "--fix"}
	if code := j.doctorChild(authority.primary, argv, installChildEnv(os.Environ(), executable), stderr); code != 0 {
		detail := fmt.Sprintf("the doctor repair of the primary checkout exited %d; the candidate stays published", code)
		return printTargetRefusal(stderr, buildVerb, refusalError{refusal{detail: detail, next: usage.WorktreeBuildInstallFor(targetAddress(target))}})
	}
	return printInstalled(stdout, stderr, authority, executable)
}

// installAuthority applies the install's refusals in their decided order: the build
// inputs, the running executable, the bound assignment, and the default-branch
// declaration. The running executable is checked before any ledger or spec read, so a
// candidate that the wrapper selected never reads its own authority.
func installAuthority(j joins, root, target string) (installTarget, error) {
	primary, err := git.PrimaryCheckout(root)
	if err != nil {
		return installTarget{}, refusalError{refusal{detail: "the primary checkout does not resolve: " + err.Error()}}
	}
	if !freshness.DeclaresBuildInputs(primary) {
		return installTarget{}, refusalError{refusal{detail: "the repository declares no Bench build inputs, so it has no candidate executable to install", next: "bench status"}}
	}
	executable := freshness.PublishedExecutable(primary)
	if running, primaryRuns := runsExecutable(j, executable); !primaryRuns {
		detail := "the running executable " + running + " is not the primary executable " + executable
		return installTarget{}, refusalError{refusal{detail: detail, next: "cd " + sanitize.ShellQuote(primary) + refusalroute.StepJoiner + usage.WorktreeBuildInstallFor(targetAddress(target))}}
	}
	ledger, err := intent.ReadTolerant(root)
	if err != nil {
		return installTarget{}, refusalError{refusal{detail: "the intent ledger does not read: " + err.Error(), next: "bench status"}}
	}
	assignment, err := installAssignment(root, ledger.Assignments, target)
	if err != nil {
		return installTarget{}, err
	}
	binding, bound := boundDelivery(ledger, assignment.ID)
	if !bound {
		return installTarget{}, refusalError{refusal{detail: "assignment " + assignment.ID + " has no bound deliverable, so no reviewed declaration authorizes an install", next: "bench commitment show"}}
	}
	plainBuild := usage.WorktreeBuildFor(targetAddress(target))
	fields, err := commitrepo.Store{Root: root}.DefaultDeclaration(binding.Deliverable)
	if err != nil {
		return installTarget{}, refusalError{refusal{detail: err.Error(), next: plainBuild}}
	}
	if len(fields) == 0 {
		return installTarget{}, refusalError{refusal{detail: "the default-branch copy of " + binding.Deliverable + " declares no persisted field, so the build needs no installed candidate", next: plainBuild}}
	}
	return installTarget{primary: primary, assignment: assignment, deliverable: binding.Deliverable}, nil
}

// runsExecutable reports whether this process runs executable, by canonical physical path,
// and answers the path that it runs. The wrapper stays the one source of its resolution
// order, so the check reads the wrapper's answer and does not repeat that order.
func runsExecutable(j joins, executable string) (string, bool) {
	running, err := j.resolveRunningBinary()
	if err != nil {
		return "that does not resolve (" + err.Error() + ")", false
	}
	physical, err := canonicalPath(running)
	if err != nil {
		return running, false
	}
	want, err := canonicalPath(executable)
	return physical, err == nil && physical == want
}

// installAssignment resolves target among the assignments of the tolerant ledger read and
// proves its identity bundle, as the strict resolution proves the bundle of a plain build.
func installAssignment(root string, assignments []intent.Assignment, target string) (intent.Assignment, error) {
	if !lineSafe(target) {
		return intent.Assignment{}, errTargetControls
	}
	selected, err := selectAssignment(assignments, target)
	if err != nil {
		return intent.Assignment{}, err
	}
	if err := identityBundleRefusal(root, selected.Worktree, selected, landingActiveState); err != nil {
		return intent.Assignment{}, err
	}
	return selected, nil
}

// boundDelivery is the delivery binding of the assignment id in ledger, if one exists.
func boundDelivery(ledger intent.Ledger, id string) (intent.DeliveryBinding, bool) {
	if ledger.Commitment == nil {
		return intent.DeliveryBinding{}, false
	}
	for _, binding := range ledger.Commitment.Bindings {
		if binding.Assignment == id {
			return binding, true
		}
	}
	return intent.DeliveryBinding{}, false
}

// publishInstall publishes a copy of built as the primary executable, through the
// publication owner. The seal root is the target tree, whose sources produced built, and
// the broker manifest lands in the primary wrapper directory. The version is the one that
// the target's own build bound in its manifest. The owner restores the earlier set on a
// failure, so the staged copy is the only file that this step removes itself.
func publishInstall(j joins, target installTarget, built string) (string, error) {
	manifest, err := brokermanifest.Read(filepath.Join(wrapperDir(target.assignment.Worktree), brokermanifest.Name))
	if err != nil {
		return "", err
	}
	executable := freshness.PublishedExecutable(target.primary)
	staged, err := stageCopy(built, executable)
	if err != nil {
		return "", err
	}
	defer os.Remove(staged)
	if err := j.publishCandidate(target.assignment.Worktree, staged, executable, wrapperDir(target.primary), manifest["version"]); err != nil {
		return "", err
	}
	return executable, nil
}

// stageCopy writes a copy of source beside executable and answers its path, so the
// publication's rename stays inside one directory.
func stageCopy(source, executable string) (string, error) {
	data, err := os.ReadFile(source)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
		return "", err
	}
	staged, err := os.CreateTemp(filepath.Dir(executable), ".bench-install-*")
	if err != nil {
		return "", err
	}
	_, err = staged.Write(data)
	if err == nil {
		err = staged.Chmod(0o755)
	}
	if closeErr := staged.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(staged.Name())
		return "", err
	}
	return staged.Name(), nil
}

// wrapperDir is the directory of a checkout's wrapper, which holds its broker manifest.
func wrapperDir(root string) string { return filepath.Join(root, "bin") }

// installChildEnv is the environment of the doctor repair that follows an install: base
// without the wrapper routing, with executable selected. The child then resolves the
// primary checkout's own kit, and the repair binds the executable that it runs.
func installChildEnv(base []string, executable string) []string {
	return runbinary.WithEnv(env.WithoutWrapperRouting(base), executable)
}

// runDoctorChild is the default doctorChild. The child writes to stderr, so stdout keeps
// the install's one table.
func runDoctorChild(dir string, argv, environment []string, stderr io.Writer) int {
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir, cmd.Env, cmd.Stdout, cmd.Stderr = dir, environment, stderr, stderr
	err := cmd.Run()
	if err != nil && cmd.ProcessState == nil {
		fmt.Fprintf(stderr, "%s: %v\n", buildVerb, err)
	}
	return childExitCode(cmd, err)
}

// printInstalled prints the install's table and the first ticket dispatch of the bound spec.
func printInstalled(stdout, stderr io.Writer, target installTarget, executable string) int {
	digest, err := brokermanifest.Digest(executable)
	if err != nil {
		fmt.Fprintf(stderr, "%s: digest the primary executable: %v\n", buildVerb, err)
		return 1
	}
	manifest := filepath.Join(wrapperDir(target.primary), brokermanifest.Name)
	table, err := toon.Table(installTable, []string{"worktree", "sha256", "manifest"}, [][]string{{target.assignment.ID, digest, manifest}})
	if err != nil {
		fmt.Fprintln(stderr, toon.RenderError(err))
		return 1
	}
	slug := spec.LiveSpecSlug(target.deliverable)
	if lineSafe(slug) {
		slug = sanitize.ShellQuote(slug)
	} else {
		slug = "<slug>"
	}
	fmt.Fprint(stdout, table)
	fmt.Fprintf(stdout, "next[1]:\n  bench record assignment %s --help\n", slug)
	return 0
}

// targetAddress is the operator's target as a shell word for a route, or the grammar slot
// when no line can carry it.
func targetAddress(target string) string {
	if lineSafe(target) {
		return sanitize.ShellQuote(target)
	}
	return "<target>"
}
