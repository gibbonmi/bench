package otelrecord

import (
	"context"
	"os"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// ScopeName is the instrumentation scope every Bench span carries.
const ScopeName = "github.com/gibbonmi/bench"

// Provider owns one repository's tracer and the processor behind it. A verb boundary
// builds one, puts its tracer on the context, and shuts it down when the verb exits.
type Provider struct {
	provider *sdktrace.TracerProvider
}

// NewProvider returns the provider that records root's spans below home. An empty home
// resolves through internal/benchhome, the one BENCH_HOME read in the tree, so a caller
// that has already resolved the home passes it in rather than reading it twice.
//
// While a test binary runs, a resolved home that names the user's own Bench home — the
// fallback home, whether BENCH_HOME named that exact path or the fallback supplied it —
// returns a provider with no span processor. The refusal is a silent no-op: every span
// still starts and ends, but nothing reaches disk. A test that names its own BENCH_HOME
// elsewhere still records, so a test that wants its own record still gets one.
func NewProvider(home, root string) *Provider {
	if home == "" {
		home = benchhome.Dir()
	}
	if testing.Testing() && benchhome.IsFallback(home) {
		return &Provider{provider: sdktrace.NewTracerProvider()}
	}
	return &Provider{provider: sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(newProcessor(home, root)),
	)}
}

// Tracer returns the tracer that starts Bench seam spans.
func (p *Provider) Tracer() trace.Tracer {
	return p.provider.Tracer(ScopeName)
}

// Shutdown releases the provider. The processor buffers nothing, so a missed shutdown
// loses no line.
func (p *Provider) Shutdown(ctx context.Context) error {
	return p.provider.Shutdown(ctx)
}

// tracerKey addresses the tracer on a context. The gate threads its run log the same
// way, and the tracer follows that pattern rather than a package-level global.
type tracerKey struct{}

// WithTracer puts the tracer on the context for the verb's whole call tree.
func WithTracer(ctx context.Context, tracer trace.Tracer) context.Context {
	return context.WithValue(ctx, tracerKey{}, tracer)
}

// TracerFrom returns the context's tracer, and a no-op tracer when no boundary put one
// there. A seam therefore starts its span without asking whether recording is on.
func TracerFrom(ctx context.Context) trace.Tracer {
	if tracer, ok := ctx.Value(tracerKey{}).(trace.Tracer); ok && tracer != nil {
		return tracer
	}
	return noop.NewTracerProvider().Tracer(ScopeName)
}

// Begin opens one seam's span below home and returns the span with the closer that ends
// it. The closer ends the span and shuts the provider down, in that order, so the end
// line is written before the provider goes away. The returned context carries the tracer
// and the span, so a seam that runs work below it parents that work here.
//
// Every instrumented seam opens through this call. The protocol — build the provider,
// start the span with its seam attribute, end, shut down — has one source, so a new seam
// cannot record a span the consumer cannot read.
func Begin(home, root, seam string) (context.Context, trace.Span, func()) {
	return BeginIn(context.Background(), home, root, seam, seam)
}

// BeginIn is Begin with the parent context given, and with the span name separate from
// the seam. A gate run and a lane run start below a context that is already threaded,
// and both name the span for the mode or the lane while the seam attribute stays the
// seam. An empty name takes the seam. The attrs join the seam attribute at start, so the
// start line carries them too.
func BeginIn(ctx context.Context, home, root, seam, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span, func()) {
	if name == "" {
		name = seam
	}
	provider := NewProvider(home, root)
	tracer := provider.Tracer()
	ctx = WithTracer(ctx, tracer)
	attrs = append([]attribute.KeyValue{attribute.String(AttrSeam, seam)}, attrs...)
	ctx, span := tracer.Start(ctx, name, trace.WithAttributes(attrs...))
	return ctx, span, func() {
		span.End()
		_ = provider.Shutdown(context.WithoutCancel(ctx))
	}
}

// The trace handoff to a child process. A child that records under its parent's trace
// reads the repository and the parent span from these two variables. This package owns
// both names, so every composer and every child reads one spelling.
const (
	handoffRootEnv        = "BENCH_OTEL_ROOT"
	handoffTraceparentEnv = "BENCH_OTEL_TRACEPARENT"
)

// HandoffVariables returns every handoff variable name. A composer strips each inherited
// value before it hands a child its own handoff.
func HandoffVariables() []string {
	return []string{handoffRootEnv, handoffTraceparentEnv}
}

// WithHandoff appends to env the handoff for the span on ctx, recorded under root. A
// context with no span adds the root only, so the child records in a trace of its own.
func WithHandoff(ctx context.Context, root string, env []string) []string {
	env = append(env, handoffRootEnv+"="+root)
	carrier := propagation.MapCarrier{}
	propagation.TraceContext{}.Inject(ctx, carrier)
	if parent := carrier.Get("traceparent"); parent != "" {
		env = append(env, handoffTraceparentEnv+"="+parent)
	}
	return env
}

// AttachHandoff attaches this process to the record and the trace that its parent handed
// off. The returned context carries the tracer and the parent span, and the closer shuts
// the provider down. A process with no handoff gets its context back unchanged and a
// closer that does nothing, so TracerFrom answers a no-op tracer.
func AttachHandoff(ctx context.Context) (context.Context, func()) {
	root := os.Getenv(handoffRootEnv)
	if root == "" {
		return ctx, func() {}
	}
	provider := NewProvider("", root)
	ctx = WithTracer(ctx, provider.Tracer())
	if parent := os.Getenv(handoffTraceparentEnv); parent != "" {
		ctx = propagation.TraceContext{}.Extract(ctx, propagation.MapCarrier{"traceparent": parent})
	}
	return ctx, func() { _ = provider.Shutdown(context.WithoutCancel(ctx)) }
}
