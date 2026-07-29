package gate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// reuseMarkerRepo is a throwaway repo whose gate appends one line per run to
// .git/runs, so "the oracle did not run again" is a counted fact rather than an
// inference from the returned exit code.
func reuseMarkerRepo(t *testing.T, exit int, manifest string) string {
	t.Helper()
	return gateTestRepo(t, fmt.Sprintf("#!/usr/bin/env bash\necho run >> .git/runs\nexit %d\n", exit), manifest)
}

func gateRunCount(t *testing.T, root string) int {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, ".git", "runs"))
	if errors.Is(err, os.ErrNotExist) {
		return 0
	}
	if err != nil {
		t.Fatal(err)
	}
	return len(strings.Fields(string(data)))
}

// TestReusableGreenIsReusedWithoutRunningOrWriting grades the three halves of a reuse
// that no single assertion covers: the oracle is not re-run, the durable record is not
// rewritten (a rewritten RecordedAt would slide the freshness window forward on every
// read), and the returned tuple projects the green rather than a zero Inspection that
// would read as `absent` on a green tree.
func TestReusableGreenIsReusedWithoutRunningOrWriting(t *testing.T) {
	root := reuseMarkerRepo(t, 0, `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	first := Execute(context.Background(), root, io.Discard, io.Discard)
	if first.ActionExit != 0 || !first.Inspection.ReusableGreen {
		t.Fatalf("first execution = %+v, want a reusable green", first)
	}
	before := mustRead(t, cachePath(t, root))

	var stdout bytes.Buffer
	second := Execute(context.Background(), root, &stdout, io.Discard)
	if got := gateRunCount(t, root); got != 1 {
		t.Fatalf("gate runs = %d, want 1 — the reusable green paid a second run", got)
	}
	if got, want := stdout.String(), "gate: green (fresh verdict reused for this tree)\n"; got != want {
		t.Fatalf("reuse stdout = %q, want %q", got, want)
	}
	if after := mustRead(t, cachePath(t, root)); !bytes.Equal(before, after) {
		t.Fatalf("reuse rewrote the verdict record:\nbefore %q\nafter  %q", before, after)
	}
	if second.GateExit != 0 || second.ActionExit != 0 || second.Inspection.State != Ready || second.Inspection.Status != "green" || !second.Inspection.ReusableGreen {
		t.Fatalf("reused result = %+v, want 0/0 with an inspection projecting the reusable green", second)
	}
}

// TestNonReusableSubjectsPayARealRun pins the short-circuit to the ReusableGreen
// predicate. The cheapest wrong implementation short-circuits on any cached green; each
// case here leaves exactly one run behind and requires the next execution to reach the
// oracle anyway.
func TestNonReusableSubjectsPayARealRun(t *testing.T) {
	const closed = `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`
	seedGreen := func(t *testing.T, root string) {
		t.Helper()
		if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
			t.Fatalf("seed execution = %+v, want green", got)
		}
	}
	cases := []struct {
		name    string
		arrange func(t *testing.T) (string, gateEngine)
	}{
		{"recorded red", func(t *testing.T) (string, gateEngine) {
			root := reuseMarkerRepo(t, 1, closed)
			if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit == 0 {
				t.Fatalf("red gate execution = %+v, want a red verdict", got)
			}
			return root, productionGateEngine{}
		}},
		{"expired verdict", func(t *testing.T) (string, gateEngine) {
			root := reuseMarkerRepo(t, 0, closed)
			seedGreen(t, root)
			return root, &faultEngine{now: time.Now().UTC().Truncate(time.Second).Add(freshness + time.Minute)}
		}},
		{"pending record", func(t *testing.T) (string, gateEngine) {
			root := reuseMarkerRepo(t, 0, closed)
			seedGreen(t, root)
			plan := mustSubject(t, root)
			pending := verdictRecord{Schema: 1, State: Pending, Tree: plan.Tree, Oracle: plan.Oracle, StartedAt: time.Now().UTC().Add(-time.Hour).Truncate(time.Second).Format(time.RFC3339), OwnerPID: 99999999}
			if err := durableReplace(filepath.Dir(cachePath(t, root)), pending); err != nil {
				t.Fatal(err)
			}
			return root, productionGateEngine{}
		}},
		{"open subject", func(t *testing.T) (string, gateEngine) {
			root := reuseMarkerRepo(t, 0, "")
			seedGreen(t, root)
			return root, productionGateEngine{}
		}},
		{"changed tree", func(t *testing.T) (string, gateEngine) {
			root := reuseMarkerRepo(t, 0, closed)
			seedGreen(t, root)
			writeGateTestFile(t, root, "changed.txt", "changed\n", 0o644)
			return root, productionGateEngine{}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, engine := tc.arrange(t)
			if got := gateRunCount(t, root); got != 1 {
				t.Fatalf("arranged gate runs = %d, want 1", got)
			}
			if got := inspectAt(root, engine.Now()); got.ReusableGreen {
				t.Fatalf("arranged subject is reusable: %+v", got)
			}
			executeWithEngine(context.Background(), root, io.Discard, io.Discard, engine)
			if got := gateRunCount(t, root); got != 2 {
				t.Fatalf("gate runs = %d, want 2 — a non-reusable subject skipped the oracle", got)
			}
		})
	}
}
