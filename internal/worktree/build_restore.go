// The restore form of the build verb: after a declared build stops before its landing, it
// clears the build's declared ledger fields under the installed candidate and returns the
// primary checkout to its own build.
package worktree

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// restoreAttributed is the restore form of the build verb, for the stopped build of the spec
// at slug. It copies the ledger under the Bench home and removes each intent-ledger field
// that the default-branch copy of the spec declares. Then it rebuilds the primary executable
// from the primary sources and runs the doctor repair under it, so `main` code binds the
// broker manifest again. Each authority refusal comes before any write.
func restoreAttributed(j joins, a ambient, root, slug string, stdout, stderr io.Writer) int {
	rerun := usage.WorktreeBuildRestoreFor(routeWord(slug, slugSlot))
	primary, paths, err := restoreAuthority(root, spec.LiveSpecPath(slug), rerun)
	if err != nil {
		return printTargetRefusal(stderr, buildVerb, err)
	}
	preserved := filepath.Join(a.home, "restore", poolkey.Key(primary), fmt.Sprintf("%s.%d", intent.Filename, a.now.UnixNano()))
	if err := intent.RemoveFields(primary, preserved, paths); err != nil {
		return printTargetRefusal(stderr, buildVerb, refusalError{refusal{detail: "the ledger clear refused: " + err.Error(), next: rerun}})
	}
	executable := freshness.PublishedExecutable(primary)
	ctx, stop := subprocess.NotifyCancel(context.Background())
	defer stop()
	if buildErr := j.buildSubject(ctx, primary, executable); buildErr != nil {
		fmt.Fprintf(stderr, "%s: %v\n", buildVerb, buildErr)
		return nameWorktree(stderr, primary, buildExitCode(buildErr))
	}
	if code := repairPrimary(j, primary, executable, stderr); code != 0 {
		detail := fmt.Sprintf("the doctor repair of the primary checkout exited %d; the ledger copy is %s", code, preserved)
		return printTargetRefusal(stderr, buildVerb, refusalError{refusal{detail: detail, next: rerun}})
	}
	table, err := toon.Table(restoreTable, []string{"ledger_copy", "executable"}, [][]string{{preserved, executable}})
	if err != nil {
		fmt.Fprintln(stderr, toon.RenderError(err))
		return 1
	}
	fmt.Fprint(stdout, table)
	return 0
}

// restoreAuthority applies the restore's refusals in their decided order: the build inputs,
// the default-branch declaration of deliverable, the strict ledger read, and the open builds.
// The strict read comes before the open-build test, which reads the ledger too and would
// print a decode error with no route. It answers the primary checkout and the declared
// intent-ledger paths.
func restoreAuthority(root, deliverable, rerun string) (string, [][]string, error) {
	primary, err := primaryWithInputs(root)
	if err != nil {
		return "", nil, err
	}
	store := commitrepo.Store{Root: primary}
	fields, err := store.DefaultDeclaration(deliverable)
	if err != nil {
		return "", nil, refusalError{refusal{detail: err.Error(), next: "bench status"}}
	}
	if len(fields) == 0 {
		return "", nil, refusalError{refusal{detail: "the default-branch copy of " + deliverable + " declares no persisted field, so the restore has no field to clear", next: "bench status"}}
	}
	if _, err := intent.Read(primary); err != nil {
		next := usage.WorktreeBuildInstallFor(restoreTarget(primary, deliverable)) + refusalroute.StepJoiner + rerun
		return "", nil, refusalError{refusal{detail: "this executable cannot read the intent ledger: " + err.Error(), next: next}}
	}
	open, err := store.HoldsOpenBuild()
	if err != nil {
		return "", nil, err
	}
	if open {
		detail := "the ledger holds an open build, and the clear is ledger-wide; end the build, then run the restore before bench worktree release, because a recovery install needs the target worktree"
		return "", nil, refusalError{refusal{detail: detail, next: refusalroute.BlockRoute(rerun)}}
	}
	var paths [][]string
	for _, field := range fields {
		if field.Record == spec.RecordIntentLedger {
			paths = append(paths, field.Path)
		}
	}
	return primary, paths, nil
}

// restoreTarget is the install target of the stopped build of deliverable, as a route word:
// the label of the first bound assignment that the tolerant read finds with its tree, or the
// target slot when that lookup fails.
func restoreTarget(root, deliverable string) string {
	ledger, err := intent.ReadTolerant(root)
	if err != nil || ledger.Commitment == nil {
		return targetSlot
	}
	for _, binding := range ledger.Commitment.Bindings {
		if binding.Deliverable != deliverable {
			continue
		}
		if assignment, err := resolveAmong(root, ledger.Assignments, binding.Assignment, landingActiveState); err == nil {
			return targetAddress(assignment.Label)
		}
	}
	return targetSlot
}
