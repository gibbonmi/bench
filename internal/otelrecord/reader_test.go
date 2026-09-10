package otelrecord

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
	"go.opentelemetry.io/otel/attribute"
)

// beginRecorded opens one span through the real record path and returns its closer.
func beginRecorded(ctx context.Context, home, root, seam, name string) (context.Context, func()) {
	child, _, finish := BeginIn(ctx, home, root, seam, name)
	return child, finish
}

func TestReadSpansAnswersTheFinishedSpansOfTheRecord(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	landing, landingSpan, finishLanding := BeginIn(context.Background(), home, root, SeamLanding, SeamLanding)
	landingSpan.SetAttributes(attribute.String(AttrSubjectID, "abc123"))
	_, phase, finishPhase := BeginIn(landing, home, root, SeamGatePhase, "build")
	phase.SetAttributes(attribute.String(AttrOutcome, OutcomeGreen))
	time.Sleep(15 * time.Millisecond)
	finishPhase()
	finishLanding()

	spans, err := ReadSpans(home, root)
	if err != nil || len(spans) != 2 {
		t.Fatalf("ReadSpans = %d spans, %v, want the two finished spans", len(spans), err)
	}
	stage, landed := spans[0], spans[1]
	if stage.Name != "build" || stage.Seam != SeamGatePhase || stage.Attributes[AttrOutcome] != OutcomeGreen {
		t.Fatalf("stage span = %+v, want the named gate phase with its outcome", stage)
	}
	if landed.Seam != SeamLanding || landed.Attributes[AttrSubjectID] != "abc123" {
		t.Fatalf("landing span = %+v, want the landing seam with its subject", landed)
	}
	if stage.TraceID != landed.TraceID || stage.ParentSpanID != landed.SpanID {
		t.Fatalf("stage %s/%s does not sit under landing %s/%s", stage.TraceID, stage.ParentSpanID, landed.TraceID, landed.SpanID)
	}
	if stage.Elapsed() < 15*time.Millisecond {
		t.Fatalf("stage elapsed = %s, want at least the 15ms the phase ran", stage.Elapsed())
	}
}

func TestReadSpansRefusesAnAbsentRecord(t *testing.T) {
	spans, err := ReadSpans(t.TempDir(), t.TempDir())
	if err == nil || spans != nil {
		t.Fatalf("ReadSpans over an absent record = %v, %v, want a refusal", spans, err)
	}
}

func TestReadSpansSkipsAMalformedLine(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	if err := NewWriter(home, root).Append([]byte("not a record line")); err != nil {
		t.Fatal(err)
	}
	_, _, finish := BeginIn(context.Background(), home, root, SeamLanding, SeamLanding)
	finish()

	spans, err := ReadSpans(home, root)
	if err != nil || len(spans) != 1 || spans[0].Seam != SeamLanding {
		t.Fatalf("ReadSpans past a malformed line = %+v, %v, want the one finished landing", spans, err)
	}
}

func TestReadSpansRefusesAnUnreadableRecord(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	if err := NewWriter(home, root).Append([]byte("{}")); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(Path(home, root), 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(Path(home, root), 0o600) })
	if os.Geteuid() == 0 {
		capability.Capability(t, capability.Privilege, "root reads mode 0000 files; the unreadable record is unobservable")
	}
	if _, err := ReadSpans(home, root); err == nil {
		t.Fatal("ReadSpans over an unreadable record = nil error, want a refusal")
	}
}

func TestNewestLandingChoosesTheNewestCompletedLandingTrace(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	for _, subject := range []string{"older", "newer"} {
		landing, span, finish := BeginIn(context.Background(), home, root, SeamLanding, SeamLanding)
		span.SetAttributes(attribute.String(AttrSubjectID, subject))
		_, closePhase := beginRecorded(landing, home, root, SeamGatePhase, subject+"-build")
		closePhase()
		finish()
	}
	// A gate run of the operator's own, outside every landing trace, stays unlisted.
	foreign, closeForeign := beginRecorded(context.Background(), home, root, "gate", "gate.ordinary")
	_, closeForeignPhase := beginRecorded(foreign, home, root, SeamGatePhase, "foreign-build")
	closeForeignPhase()
	closeForeign()

	got, ok := NewestLanding(home, root)
	if !ok || got.Commit != "newer" {
		t.Fatalf("NewestLanding = %+v, %v, want the newer landing", got, ok)
	}
	if len(got.Stages) != 1 || got.Stages[0].Name != "newer-build" {
		t.Fatalf("stages = %+v, want only the newer landing's own phase", got.Stages)
	}
}

func TestNewestLandingAnswersFalseWithoutALanding(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	if _, ok := NewestLanding(home, root); ok {
		t.Fatal("NewestLanding over an absent record = ok, want false")
	}
	_, finish := beginRecorded(context.Background(), home, root, "gate", "gate.ordinary")
	finish()
	if got, ok := NewestLanding(home, root); ok {
		t.Fatalf("NewestLanding over a record with no landing = %+v, ok, want false", got)
	}
}

// The reader selects by two seam names the packages that open those spans declare for
// themselves. This test is the reconciliation that keeps the pair from drifting.
func TestReaderSeamsAreRegisteredSeams(t *testing.T) {
	registered := map[string]bool{}
	for _, entry := range Registry {
		registered[entry.Seam] = true
	}
	for _, seam := range []string{SeamLanding, SeamGatePhase} {
		if !registered[seam] {
			t.Fatalf("reader seam %q names no Registry row", seam)
		}
	}
}
