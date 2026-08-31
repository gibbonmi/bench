package otelrecord

import (
	"context"

	"github.com/gibbonmi/bench/internal/benchhome"
	"go.opentelemetry.io/otel/attribute"
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
func NewProvider(home, root string) *Provider {
	if home == "" {
		home = benchhome.Dir()
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
// seam. An empty name takes the seam.
func BeginIn(ctx context.Context, home, root, seam, name string) (context.Context, trace.Span, func()) {
	if name == "" {
		name = seam
	}
	provider := NewProvider(home, root)
	tracer := provider.Tracer()
	ctx = WithTracer(ctx, tracer)
	ctx, span := tracer.Start(ctx, name, trace.WithAttributes(attribute.String(AttrSeam, seam)))
	return ctx, span, func() {
		span.End()
		_ = provider.Shutdown(context.WithoutCancel(ctx))
	}
}
