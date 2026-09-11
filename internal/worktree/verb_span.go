package worktree

import (
	"github.com/gibbonmi/bench/internal/otelrecord"
	"go.opentelemetry.io/otel/attribute"
)

// The worktree mutation verbs share one span boundary and identify the assignment as their subject.
// The landing keeps its own boundary because it carries different measures.
const (
	otelCreateSeam      = "worktree.create"
	otelExecSeam        = "worktree.exec"
	otelMergeSeam       = "worktree.merge"
	otelResetSeam       = "worktree.reset"
	otelReleaseSeam     = "worktree.release"
	otelBuildSeam       = "worktree.build"
	otelReauthorizeSeam = "worktree.reauthorize"
)

// otelVerbSeams names the mutation seams in registry order.
// Plans and read-only verbs open no span.
var otelVerbSeams = []string{
	otelCreateSeam,
	otelExecSeam,
	otelMergeSeam,
	otelResetSeam,
	otelReleaseSeam,
	otelBuildSeam,
	otelReauthorizeSeam,
}

// beginVerbSpan starts one worktree verb's span and returns the closer that ends it. The
// closer takes the assignment id, because a verb resolves its assignment inside the span
// and a refusal before that resolution has none. The id passes through the encoder, which
// escapes every control rune, so a hostile id forges no second record line.
func beginVerbSpan(home, root, seam string) func(int, string) {
	_, span, finish := otelrecord.Begin(home, root, seam)
	return func(exit int, assignment string) {
		if assignment != "" {
			span.SetAttributes(attribute.String(otelrecord.AttrSubjectID, assignment))
		}
		span.SetAttributes(attribute.String(otelrecord.AttrOutcome, otelrecord.PublishedExitOutcome(exit)))
		finish()
	}
}
