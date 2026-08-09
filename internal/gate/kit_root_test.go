package gate

import (
	"bytes"
	"context"
	"io"
	"reflect"
	"testing"
	"time"
)

func TestInjectedKitSurvivesHostileAmbientThroughEvaluationAndComposition(t *testing.T) {
	fixture := newKitShapedFixture(t)

	evaluation := newWorkingTreeEvaluationAtKit(fixture.root, fixture.root)
	plan, err := evaluation.acceptPre()
	if err != nil {
		t.Fatalf("accept fixture subject: %v", err)
	}
	wantIdentities, err := resolveComponentIdentitiesAtKit(fixture.root, fixture.root, evaluation.pre)
	if err != nil {
		t.Fatalf("resolve injected fixture identities: %v", err)
	}
	t.Setenv("BENCH_KIT", t.TempDir())
	scoping := evaluation.scope(plan.Resolution, forceRun, time.Now())
	if !scoping.eligible || !sameDirectory(scoping.runnerRoot, fixture.root) {
		t.Fatalf("evaluation scoping = %+v, want fixture-root eligibility", scoping)
	}
	if len(scoping.identities) == 0 {
		t.Fatal("evaluation resolved no fixture component identities")
	}
	table, err := phaseTable(fixture.root, fixture.root)
	if err != nil {
		t.Fatalf("resolve injected fixture phase table: %v", err)
	}
	if !reflect.DeepEqual(scoping.identities, wantIdentities) {
		t.Fatalf("evaluation identities = %v, want fixture identities %v", scoping.identities, wantIdentities)
	}
	if got, want := phaseNames(table), fixture.phaseNames(); !reflect.DeepEqual(got, want) {
		t.Fatalf("injected phase table = %v, want fixture table %v", got, want)
	}
	if !partialComposedGreen(fixture.root, fixture.root, plan, verdictRecord{}, time.Now()) {
		t.Fatal("composed-green scoping rejected the injected fixture kit")
	}
	if result := executeWithEngineAtKit(context.Background(), fixture.root, fixture.root, io.Discard, io.Discard, productionGateEngine{}); result.ActionExit != 0 {
		t.Fatalf("injected execution = %+v, want green", result)
	}
}

func TestPhasesCommandEmptyKitFallsBackToRoot(t *testing.T) {
	root := t.TempDir()
	t.Setenv("BENCH_KIT", "")
	original := benchkitPhasesForCommand
	t.Cleanup(func() { benchkitPhasesForCommand = original })
	var gotKit string
	benchkitPhasesForCommand = func(gotRoot, kit string) []Phase {
		if gotRoot != root {
			t.Fatalf("phase root = %q, want %q", gotRoot, root)
		}
		gotKit = kit
		return []Phase{{Name: "fallback", Argv: []string{"true"}}}
	}

	var stdout, stderr bytes.Buffer
	if code := phasesCommandAtKitForTest(root, kitRoot(root), &stdout, &stderr); code != 0 {
		t.Fatalf("PhasesCommand = %d, want 0; stderr=%q", code, stderr.String())
	}
	if gotKit != root {
		t.Fatalf("empty BENCH_KIT resolved to %q, want root %q", gotKit, root)
	}
}

func TestPhasesCommandNonEmptyKitSelectsBuiltinTableAtKit(t *testing.T) {
	root := t.TempDir()
	kit := t.TempDir()
	t.Setenv("BENCH_KIT", kit)
	original := benchkitPhasesForCommand
	t.Cleanup(func() { benchkitPhasesForCommand = original })
	var gotKit string
	benchkitPhasesForCommand = func(gotRoot, selectedKit string) []Phase {
		if gotRoot != root {
			t.Fatalf("phase root = %q, want %q", gotRoot, root)
		}
		gotKit = selectedKit
		return []Phase{{Name: "builtin", Argv: []string{"true"}}}
	}

	var stdout, stderr bytes.Buffer
	if code := phasesCommandAtKitForTest(root, kitRoot(root), &stdout, &stderr); code != 0 {
		t.Fatalf("PhasesCommand = %d, want 0; stderr=%q", code, stderr.String())
	}
	if gotKit != kit {
		t.Fatalf("non-empty BENCH_KIT resolved to %q, want %q", gotKit, kit)
	}
}
