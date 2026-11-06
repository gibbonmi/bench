package releaseevidence

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// TestMain installs the release-evidence row of the composition root's table, so promotion
// consumes protection as the CLI does. The package owns the only vector, so an outer run's
// resources stay out of its generations.
func TestMain(m *testing.M) {
	if err := os.Unsetenv(protection.Env); err != nil {
		os.Exit(3)
	}
	protection.Install(protection.Protocol{Validators: protection.Validators{protection.KindReleaseEvidence: ValidateProtection}, ReadLimit: bounds.ControlRecordLimit})
	os.Exit(m.Run())
}

// TestProtectedStageCleanup is PL119. An abandoned stage carries a valid owner marker, but a
// live user of the release evidence still holds it, so promotion must retain the stage and
// replace nothing. With no user, the same marker lets promotion clean the stage.
func TestProtectedStageCleanup(t *testing.T) {
	t.Run("protected", func(t *testing.T) {
		root, stage := abandonedStage(t)
		generation, err := Protect(root)
		if err != nil {
			t.Fatal(err)
		}
		user, err := protection.Installed().Prepare(generation.Resources())
		if err != nil {
			t.Fatal(err)
		}
		// The test's own group is a live group, so it stands for the stage's writer.
		if err := user.Started(os.Getpid(), syscall.Getpgrp()); err != nil {
			t.Fatal(err)
		}

		err = PromoteEvidence(root, ModeVerify, nil, Manifest{SchemaVersion: 1})
		if !errors.Is(err, protection.ErrRetained) {
			t.Fatalf("promotion over a protected stage = %v, want the retention", err)
		}
		for _, path := range []string{stage, stage + ".owner"} {
			if _, err := os.Lstat(path); err != nil {
				t.Fatalf("protected %s = %v, want it retained", filepath.Base(path), err)
			}
		}
		if _, err := os.Lstat(filepath.Join(root, "dist", "preflight")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("preflight target after the refusal = %v, want none", err)
		}
	})
	t.Run("unused", func(t *testing.T) {
		root, stage := abandonedStage(t)
		if err := PromoteEvidence(root, ModeVerify, nil, Manifest{SchemaVersion: 1}); err != nil {
			t.Fatal(err)
		}
		for _, path := range []string{stage, stage + ".owner"} {
			if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("unused %s = %v, want it cleaned", filepath.Base(path), err)
			}
		}
	})
}

// TestExchangeFailureKeepsThePriorEvidence drives the exchange seam: a real filesystem
// cannot fail the atomic directory exchange on demand. A failed exchange refuses the
// promotion, keeps the prior preflight evidence, and removes the new stage.
func TestExchangeFailureKeepsThePriorEvidence(t *testing.T) {
	root, _ := abandonedStage(t)
	if err := PromoteEvidence(root, ModeVerify, nil, Manifest{SchemaVersion: 1}); err != nil {
		t.Fatal(err)
	}
	manifest := filepath.Join(root, "dist", "preflight", "manifest.json")
	prior, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	refused := errors.New("forced exchange failure")
	t.Cleanup(SetExchangeForTesting(func(string, string) error { return refused }))
	if err := PromoteEvidence(root, ModeVerify, nil, Manifest{SchemaVersion: 1, Mode: ModeVerify}); !errors.Is(err, refused) {
		t.Fatalf("promotion over a failed exchange = %v, want the exchange failure", err)
	}
	if after, err := os.ReadFile(manifest); err != nil || string(after) != string(prior) {
		t.Fatalf("prior manifest after the refusal = %q, %v; want %q", after, err, prior)
	}
	if stages, err := filepath.Glob(filepath.Join(root, "dist", ".preflight-stage-*")); err != nil || len(stages) != 0 {
		t.Fatalf("stages after the refusal = %v, %v; want none", stages, err)
	}
}

// abandonedStage plants a stage and its valid owner marker in a fresh checkout's dist.
func abandonedStage(t *testing.T) (root, stage string) {
	t.Helper()
	root = t.TempDir()
	dist := filepath.Join(root, "dist")
	if err := os.MkdirAll(dist, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { discardStore(t, root) })
	stage, err := os.MkdirTemp(dist, ".preflight-stage-")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stage+".owner", stageOwnerBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	return root, stage
}

// discardStore removes the protection store of root's evidence, which lives under the shared
// store root and outlives the temporary checkout.
func discardStore(t *testing.T, root string) {
	t.Helper()
	d, err := Descriptor(root)
	if err == nil {
		err = protection.Generation{Resource: d}.Discard()
	}
	if err != nil {
		t.Error(err)
	}
}
