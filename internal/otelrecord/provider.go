package otelrecord

import (
	"context"

	"github.com/gibbonmi/bench/internal/benchhome"
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
