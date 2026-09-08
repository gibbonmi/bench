package preflight

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/diff"
)

func TestChargeHeadAndIndexMovementRefuseAfterOneRetry(t *testing.T) {
	for _, movement := range []string{"head", "index", "required source"} {
		t.Run(movement, func(t *testing.T) {
			root, slug := seedConformant(t)
			args := chargeArgs(t, root, slug, false)
			calls := 0
			restore := diff.SetSnapshotAfterReadForTest(func() {
				calls++
				path := "internal/example/movement.go"
				body := "package example\n// movement " + string(rune('0'+calls)) + "\n"
				if movement == "required source" {
					path = buildPhase
					body = "# Build phase\n\nmovement " + string(rune('0'+calls)) + "\n"
				}
				mustWriteFile(t, path, body)
				if movement != "required source" {
					runGit(t, "add", path)
				}
				if movement == "head" {
					runGit(t, "commit", "-q", "-m", "move head")
				}
			})
			out, code := Command(args)
			restore()
			if code != 1 || calls != 2 || !strings.Contains(out, "error: snapshot drift") ||
				!strings.Contains(out, "retry the exact invocation") ||
				strings.Contains(out, "complete") {
				t.Fatalf("persistent %s movement = (%d, %d):\n%s", movement, code, calls, out)
			}
		})
	}
}

func TestChargeDirtySourceRefusesWithoutCompleteOutput(t *testing.T) {
	root, slug := seedConformant(t)
	args := chargeArgs(t, root, slug, false)
	mustWriteFile(t, "internal/example/dirty.go", "package example\n")
	out, code := Command(args)
	if code != 1 || !strings.Contains(out, "checkout required") ||
		!strings.Contains(out, "source checkout is dirty") ||
		!strings.Contains(out, "commit or remove local changes") ||
		strings.Contains(out, "complete") {
		t.Fatalf("dirty charge = (%d):\n%s", code, out)
	}
}

func TestChargeRepeatedPinnedInputsAreIdentical(t *testing.T) {
	root, slug := seedConformant(t)
	args := chargeArgs(t, root, slug, true)
	first, firstCode := Command(args)
	second, secondCode := Command(args)
	if firstCode != 0 || secondCode != 0 || first != second {
		t.Fatalf("repeated charge = (%d, %d, equal=%t)\nfirst:\n%s\nsecond:\n%s",
			firstCode, secondCode, first == second, first, second)
	}
}

func TestChargeCompactNamesExactFullRetrieval(t *testing.T) {
	root, slug := seedConformant(t)
	args := chargeArgs(t, root, slug, false)
	out, code := Command(args)
	want := "bench preflight build specs/example/spec.md --charge --ticket one.md --base " +
		args[6] + " --source-tip " + args[8] + " --full"
	if code != 0 || !strings.Contains(out, "\"false\"") ||
		!strings.Contains(out, want) || !strings.Contains(out, "omitted[5]{source}") {
		t.Fatalf("compact retrieval = (%d), want %q:\n%s", code, want, out)
	}
}

func TestChargeFinalSnapshotFailureDiscardsPreparedOutput(t *testing.T) {
	root, slug := seedConformant(t)
	head := filepath.Join(root, ".git", "HEAD")
	moved := head + ".during-charge"
	restoreSeam := diff.SetSnapshotAfterReadForTest(func() {
		if err := os.Rename(head, moved); err != nil {
			t.Fatal(err)
		}
	})
	defer restoreSeam()

	out, code := Command(chargeArgs(t, root, slug, true))
	if err := os.Rename(moved, head); err != nil {
		t.Fatal(err)
	}
	if code != 1 || !strings.Contains(out, "snapshot identity failed") ||
		strings.Contains(out, "charge[") || strings.Contains(out, "complete") {
		t.Fatalf("final snapshot failure = (%d):\n%s", code, out)
	}
}
