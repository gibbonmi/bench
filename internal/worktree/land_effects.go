// Landing effects: the work a landing owns after publication and release, each with its
// own recorded result and its own resume path.
package worktree

import (
	"context"
	"fmt"
	"io"

	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// An effect's result word. A landing reports one of these for every effect it owns, so a
// reader tells a skipped effect from a finished one without a second read of the tree.
const (
	effectComplete      = "complete"
	effectFailed        = "failed"
	effectNotApplicable = "not-applicable"
)

// refreshEffect names the broker refresh in the effects row and in the incomplete
// record's step cell, so `incomplete:refresh` and the row cell cannot disagree.
const refreshEffect = "refresh"

// effectRow is one effect's line in the landing's effects table.
type effectRow struct {
	effect string
	result string
}

// landedAfterEffects runs the landing's post-release effects at root, prints their table,
// and renders the landing's terminal record. It is the one place the effects and the
// record meet, so the first run and the resume state the same thing about the same tree.
// A failed effect reaches the existing incomplete render under its own step name, which
// carries the resume command every other incomplete step carries.
func landedAfterEffects(j joins, root string, result landing.ReviewedResult, specArg, path, assignment string, active bool, records int, stdout, stderr io.Writer) int {
	rows := []effectRow{{effect: refreshEffect, result: refreshBroker(j, root, stderr)}}
	fmt.Fprintf(stdout, "effects[%d]{effect,result}:\n", len(rows))
	for _, row := range rows {
		fmt.Fprintf(stdout, "  %s,%s\n", row.effect, row.result)
	}
	for _, row := range rows {
		if row.result == effectFailed {
			return landedIncomplete(stdout, result, specArg, path, assignment, row.effect, records)
		}
	}
	return landedComplete(stdout, result, active, records)
}

// refreshBroker republishes root's own Bench executable after the landing published the
// sources it is built from. It applies only where the destination declares Bench build
// inputs, which is the predicate the landing's broker notice reads; a linked installation
// declares none and stays untouched.
//
// Completion is a predicate over the tree, never a recorded step: the published
// executable verifies against the destination's own sources. The publication owner
// promotes the seal last and installs the manifest with it, so a verified seal also
// proves the manifest of that same transaction landed. A resume therefore needs no
// journal, and an interrupted refresh finishes on the next call.
func refreshBroker(j joins, root string, stderr io.Writer) string {
	if !freshness.DeclaresBuildInputs(root) {
		return effectNotApplicable
	}
	executable := freshness.PublishedExecutable(root)
	if freshness.Verify(root, executable) == nil {
		return effectComplete
	}
	// The build runs the destination's own entry point, which is candidate code. It runs
	// here and nowhere else, after the release, so it can reach neither the gate, nor the
	// composition, nor the destination compare-and-swap that published these bytes. The
	// refresh never executes the published executable itself.
	if err := j.buildSubject(context.Background(), root, executable); err != nil {
		fmt.Fprintln(stderr, "landing refresh build failed: "+sanitize.Controls(err.Error()))
	}
	// The build's own verdict does not settle the effect, because a build can return
	// zero and still leave an executable no seal authenticates. So the predicate is read
	// a second time, and that read is what reports the effect.
	if err := freshness.Verify(root, executable); err != nil {
		fmt.Fprintln(stderr, "landing refresh failed: "+sanitize.Controls(err.Error()))
		return effectFailed
	}
	return effectComplete
}
