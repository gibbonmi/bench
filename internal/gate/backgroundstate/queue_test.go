package backgroundstate_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
)

// TestWithdrawTakesOutOnlyTheWorktreeAndItsChildren is BG-C8-C1 at the queue. The withdrawal of
// a worktree drops the queued grades and the latest keys of that worktree and of each worktree
// under it, and it answers the running grade of that worktree. A worktree whose path only starts
// with the same text keeps its grades.
func TestWithdrawTakesOutOnlyTheWorktreeAndItsChildren(t *testing.T) {
	entry := func(worktree, key string) backgroundstate.Entry {
		return backgroundstate.Entry{Key: key, Worktree: worktree, Tip: strings.Repeat("c", 40)}
	}
	child, prefix := entry("/work/wt/child", strings.Repeat("b", 64)), entry("/work/wt2", strings.Repeat("d", 64))
	running := backgroundstate.Slot{Entry: entry("/work/wt", key), Owner: 2, Boot: "boot"}
	queue := backgroundstate.Queue{
		Running: &running,
		Queued:  []backgroundstate.Entry{entry("/work/wt", strings.Repeat("a", 64)), child, prefix},
		Latest:  map[string]string{"/work/wt": key, child.Worktree: child.Key, prefix.Worktree: prefix.Key},
	}
	stopped := queue.Withdraw("/work/wt/")
	want := backgroundstate.Queue{Running: &running, Queued: []backgroundstate.Entry{prefix}, Latest: map[string]string{prefix.Worktree: prefix.Key}}
	if stopped == nil || *stopped != running || !reflect.DeepEqual(queue, want) {
		t.Fatalf("Withdraw = %+v, queue %+v; want the running grade and the queue %+v", stopped, queue, want)
	}
	other := backgroundstate.Slot{Entry: prefix, Owner: 2, Boot: "boot"}
	kept := backgroundstate.Queue{Running: &other}
	if stopped := kept.Withdraw("/work/wt"); stopped != nil || kept.Running != &other {
		t.Fatalf("Withdraw of /work/wt beside the grade of %s = %+v; want no grade to stop", other.Worktree, stopped)
	}
}
