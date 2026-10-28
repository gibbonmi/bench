package gocache

import (
	"os/exec"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	"github.com/gibbonmi/bench/internal/subprocess"
)

// TestCacheCleanStopsBackgroundGrades is BG82. The machine queue under the Bench home of the
// clean holds a running grade and a queued grade. The clean stops the running grade and waits
// for its exit, drops the queued grade, and then cleans.
func TestCacheCleanStopsBackgroundGrades(t *testing.T) {
	sleep, err := exec.LookPath("sleep")
	if err != nil {
		capability.Environment(t, "the stand-in grade is a sleep process")
	}
	grade := exec.Command(sleep, "300")
	if err := grade.Start(); err != nil {
		t.Fatal(err)
	}
	reaped := make(chan struct{})
	go func() {
		_ = grade.Wait()
		close(reaped)
	}()
	t.Cleanup(func() {
		_ = grade.Process.Kill()
		<-reaped
	})
	home := t.TempDir()
	tip := strings.Repeat("c", 40)
	slot, err := backgroundstate.SlotFor(backgroundstate.Entry{Key: strings.Repeat("a", 64), Worktree: "/work/running", Tip: tip}, grade.Process.Pid)
	if err != nil {
		t.Fatal(err)
	}
	queued := backgroundstate.Entry{Key: strings.Repeat("b", 64), Worktree: "/work/queued", Tip: tip}
	if err := backgroundstate.WriteQueue(durablefile.Native(), home, backgroundstate.Queue{Running: &slot, Queued: []backgroundstate.Entry{queued}}); err != nil {
		t.Fatal(err)
	}
	fixture := newCacheFixture(t)
	out, code := clean(append(cleanEnv(fixture.home), benchhome.Env+"="+home))
	if code != 0 {
		t.Fatalf("clean = exit %d, want 0; out=%q", code, out)
	}
	if alive, err := subprocess.Alive(grade.Process.Pid); err != nil || alive {
		t.Fatalf("grade alive after the clean = %t, %v; want it stopped", alive, err)
	}
	queue, err := backgroundstate.ReadQueue(home)
	if err != nil || len(queue.Queued) != 0 {
		t.Fatalf("queue after the clean = %+v, %v; want no queued grade", queue, err)
	}
	if got, want := entries(t, fixture.dir), []string{"README", LockFile, trimFile}; !slices.Equal(got, want) {
		t.Fatalf("directory after the clean = %q, want %q", got, want)
	}
}
