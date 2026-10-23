package testreport

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	benchenv "github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

const (
	// selectionCloseHelperEnv runs TestSelectionCloseFailure as the focused-run owner.
	selectionCloseHelperEnv = "BENCH_TEST_SELECTION_CLOSE_HELPER"
	// nestedUserGroupEnv runs TestNestedSelectionUser and names the live group it publishes.
	nestedUserGroupEnv = "BENCH_TEST_NESTED_USER_GROUP"
	// passingPackages is the packages table of the canned passing set.
	passingPackages = "packages[1]{package,status,elapsed_ms,tests_run}:\n  canned,pass,250,1\n"
)

// TestFocusedPipeCancel is PL64. The parked descendant ignores INT and keeps the focused
// child's stdout open, so the stdout consumer cannot finish before the group does. The
// cancellation must still end that descendant and report the interruption.
func TestFocusedPipeCancel(t *testing.T) {
	if mode := os.Getenv(cancelGoModeEnv); mode != "" {
		runGoHopCancelHelper(t, mode)
		return
	}
	signalGoHopCancelHelper(t, "TestFocusedPipeCancel", "test")
}

// TestListCancel is PL80. A cancelled package list keeps its own go-list interruption line,
// and its parked descendant does not survive.
func TestListCancel(t *testing.T) {
	if mode := os.Getenv(cancelGoModeEnv); mode != "" {
		runGoHopCancelHelper(t, mode)
		return
	}
	signalGoHopCancelHelper(t, "TestListCancel", "list")
}

// TestMalformedJSON is PL65. The Go child exits zero after a truncated event, so only the
// decode failure keeps that exit from certifying unreadable report data.
func TestMalformedJSON(t *testing.T) {
	output, code := commandOverEvents(t, t.TempDir(), []string{"--package", "./..."}, []string{
		`{"Action":"run","Package":"canned","Test":"TestCanned"}`,
		`{"Action":"pass","Package":"canned"`,
	}, 0)
	if want := "error: go test output malformed — "; code != 1 || !strings.HasPrefix(output, want) {
		t.Fatalf("Command = (%d, %q), want (1, prefix %q)", code, output, want)
	}
}

// TestIncompletePackage is PL79. A report is complete only when its stream completes and every
// package it names has a terminal event.
func TestIncompletePackage(t *testing.T) {
	t.Run("missing terminal", func(t *testing.T) {
		output, code := commandOverEvents(t, t.TempDir(), []string{"--package", "./..."}, append(cannedSetNamed(t, "passing").events,
			`{"Action":"run","Package":"other","Test":"TestOther"}`,
		), 0)
		if want := "error: go test reported incomplete packages — other\n"; code != 1 || output != want {
			t.Fatalf("Command = (%d, %q), want (1, %q)", code, output, want)
		}
	})
	// The terminal event arrives after the child exits, from a descendant that holds stdout.
	// The run reads the report only after the stream completes, so the event counts.
	t.Run("late terminal", func(t *testing.T) {
		passing := cannedSetNamed(t, "passing").events
		dir := t.TempDir()
		writeGoStub(t, filepath.Join(dir, "go"), "#!/usr/bin/env bash\n"+
			"printf '%s\\n' "+sanitize.ShellQuote(passing[0])+" "+sanitize.ShellQuote(passing[1])+"\n"+
			"( sleep 0.3; printf '%s\\n' "+sanitize.ShellQuote(passing[2])+" ) &\n")
		t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
		installCannedSelection(t)
		output, code := Command(t.TempDir(), []string{"--package", "./..."})
		if code != 0 || !strings.HasPrefix(output, passingPackages) {
			t.Fatalf("Command = (%d, %q), want (0, prefix %q)", code, output, passingPackages)
		}
	})
	// A forced closure ends the stream with a partial report that decodes as complete, and a
	// child exit zero cannot certify it.
	t.Run("forced closure", func(t *testing.T) {
		data := []byte(strings.Join(cannedSetNamed(t, "passing").events, "\n") + "\n")
		none := subprocess.Stream{Status: subprocess.StreamNone}
		child := subprocess.Outcome{
			Process:    subprocess.Process{Started: true, Waited: true},
			Streams:    subprocess.StreamSet{Stdin: none, Stdout: subprocess.Stream{Status: subprocess.StreamClosed, Err: os.ErrClosed, Data: data}, Stderr: none},
			Protection: protection.Disposition{State: protection.Safe},
		}
		outcome, output, code := focusedResult(focusedRequest{}, child)
		want := "error: go test cleanup incomplete — process exit 0; process stream " + string(subprocess.StreamClosed)
		if code != 1 || outcome.Kind != OutcomeRefused || !strings.HasPrefix(output, want) {
			t.Fatalf("focusedResult = (%q, %d, %q), want (refused, 1, prefix %q)", outcome.Kind, code, output, want)
		}
	})
}

// TestSelectionCloseFailure is PL122. The focused Go child exits zero with a passing report,
// but a nested user that it started keeps the selected executable protected. The run must
// consume the selection's Close failure beside the child's report and keep the executable.
func TestSelectionCloseFailure(t *testing.T) {
	if os.Getenv(selectionCloseHelperEnv) != "" {
		runSelectionCloseHelper(t)
		return
	}
	events := cannedSetNamed(t, "passing").events
	quoted := make([]string, len(events))
	for i, event := range events {
		quoted[i] = sanitize.ShellQuote(event)
	}
	goDir := t.TempDir()
	// The test's own group is a live group, so it stands for the group the nested user started.
	writeGoStub(t, filepath.Join(goDir, "go"), "#!/usr/bin/env bash\n"+
		"env "+nestedUserGroupEnv+"="+strconv.Itoa(syscall.Getpgrp())+" "+sanitize.ShellQuote(os.Args[0])+" -test.run='^TestNestedSelectionUser$' >/dev/null 2>&1 || exit 9\n"+
		"printf '%s\\n' "+strings.Join(quoted, " ")+"\n")
	helper := exec.Command(os.Args[0], "-test.run=^TestSelectionCloseFailure$", "-test.timeout="+deadline().String())
	// The run under test owns the only vector, so an outer run's resources stay out of it.
	helper.Env = protection.Environ(append(os.Environ(),
		selectionCloseHelperEnv+"=1",
		"PATH="+goDir+string(os.PathListSeparator)+os.Getenv("PATH"),
		"HOME="+t.TempDir(),
	), nil)
	if output, err := helper.CombinedOutput(); err != nil {
		t.Fatalf("selection close helper: %v\n%s", err, output)
	}
}

// runSelectionCloseHelper is the focused-run owner under the protocol of the selection and
// cache kinds. It discards the stores that the retained generations leave under the store root.
func runSelectionCloseHelper(t *testing.T) {
	installSelectionProtocol()
	var selected string
	installTestSelectionFactory(t, runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, _, output string) error {
			selected = output
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	})
	t.Cleanup(func() { discardSelectionStores(t, selected) })
	root := t.TempDir()
	request, line, _ := Prepare(root, []string{"--package", "./..."})
	if line != "" {
		t.Fatalf("Prepare = %q", line)
	}
	outcome, output, code := Execute(root, request)
	if code != 1 || outcome.Kind != OutcomeRefused || !strings.HasPrefix(output, passingPackages) {
		t.Fatalf("Execute = (%q, %d, %q), want (refused, 1) with the child report %q", outcome.Kind, code, output, passingPackages)
	}
	for _, want := range []string{"error: Bench executable selection cleanup failed — ", protection.ErrRetained.Error()} {
		if !strings.Contains(output, want) {
			t.Fatalf("Execute output = %q, want %q", output, want)
		}
	}
	if info, err := os.Stat(selected); err != nil || !info.Mode().IsRegular() {
		t.Fatalf("retained executable = %v, %v; want the selected file", info, err)
	}
}

// TestNestedSelectionUser is the nested user of TestSelectionCloseFailure. It registers on the
// selection alone in its inherited vector, publishes the live group that nestedUserGroupEnv
// names, and exits zero without a terminal record, so that group keeps the selection protected.
func TestNestedSelectionUser(t *testing.T) {
	value := os.Getenv(nestedUserGroupEnv)
	if value == "" {
		return
	}
	installSelectionProtocol()
	group, err := strconv.Atoi(value)
	if err != nil {
		t.Fatal(err)
	}
	inherited, err := protection.Inherited()
	selection := slices.DeleteFunc(inherited, func(d protection.Descriptor) bool { return d.Kind != protection.KindBinary })
	if err != nil || len(selection) != 1 {
		t.Fatalf("inherited selection = %v, %v; want the focused child's one selection", selection, err)
	}
	user, err := protection.Installed().Prepare(selection)
	if err == nil {
		err = user.Started(group, group)
	}
	if err != nil {
		t.Fatal(err)
	}
}

// installSelectionProtocol admits the kinds that a focused run registers on, as the composition
// root's table does. Only a helper process installs it.
func installSelectionProtocol() {
	protection.Install(protection.Protocol{Validators: protection.Validators{
		protection.KindBinary:  runbinary.ValidateProtection,
		protection.KindCache:   gocache.ValidateProtection,
		protection.KindKitTest: benchenv.ValidateProtection,
	}, ReadLimit: bounds.ControlRecordLimit})
}

// discardSelectionStores removes the stores of a retained selection directory and of the
// helper's cache, which live under the shared store root and outlive the temporary directories.
// A selection directory that Close removed took its store with it.
func discardSelectionStores(t *testing.T, selected string) {
	t.Helper()
	if selected != "" {
		d, err := protection.Directory{Kind: protection.KindBinary}.Describe(filepath.Dir(selected))
		if err == nil {
			err = protection.Generation{Resource: d}.Discard()
		}
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			t.Error(err)
		}
	}
	dir, err := gocache.Dir(os.Environ())
	if err == nil {
		var store string
		if store, err = gocache.ProtectionStore(dir); err == nil {
			err = os.RemoveAll(store)
		}
	}
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		t.Error(err)
	}
}
