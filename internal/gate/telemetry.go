package gate

import (
	"context"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

// The gate's seam record. The verb boundary resolves the Bench home once, builds the
// provider, and puts its tracer on the context for the whole run, the way the run log
// threads its own file. The phase table runs in a second process, so the record's
// repository and this run's trace reach that process through the environment.
const (
	otelGateSeam      = "gate"
	otelGatePhaseSeam = "gate.phase"

	otelRootEnv        = "BENCH_OTEL_ROOT"
	otelTraceparentEnv = "BENCH_OTEL_TRACEPARENT"
)

// otelGateEnv is the whole set the phases process inherits, so the stripper and the
// composer read one list rather than two that can disagree.
var otelGateEnv = []string{otelRootEnv, otelTraceparentEnv}

// otelRecordRootKey addresses the repository the run records under. The child that runs
// the phases records under the same repository, so the run's spans land in one file.
type otelRecordRootKey struct{}

// beginGateSpan starts the run's root span and returns the closer that ends it. The
// mode rides in the span name: story 19's declared attribute set names the seam, the
// subject, the outcome, and the measures, and none of those carries a run mode.
func beginGateSpan(ctx context.Context, root, mode string) (context.Context, func(Result)) {
	ctx, span, finish := otelrecord.BeginIn(ctx, "", root, otelGateSeam, otelGateSeam+"."+mode)
	ctx = context.WithValue(ctx, otelRecordRootKey{}, root)
	return ctx, func(result Result) {
		// A subject the run never resolved has no digest to group its iterations by, and
		// an empty attribute would read as one.
		if subject := result.Inspection.CurrentTree; subject != "" {
			span.SetAttributes(attribute.String(otelrecord.AttrSubjectID, subject))
		}
		// The action exit is the one the operator sees, so the record and the shell
		// agree about the same run.
		span.SetAttributes(attribute.String(otelrecord.AttrOutcome, otelrecord.ExitOutcome(result.ActionExit)))
		finish()
	}
}

// withGateSpanEnv hands the phases child the repository it records under and this run's
// trace, so its phase spans join the root span rather than start a trace of their own.
// A context outside a recorded run adds nothing.
func withGateSpanEnv(ctx context.Context, base []string) []string {
	root, _ := ctx.Value(otelRecordRootKey{}).(string)
	if root == "" {
		return base
	}
	env := append(base, otelRootEnv+"="+root)
	carrier := propagation.MapCarrier{}
	propagation.TraceContext{}.Inject(ctx, carrier)
	if parent := carrier.Get("traceparent"); parent != "" {
		env = append(env, otelTraceparentEnv+"="+parent)
	}
	return env
}
