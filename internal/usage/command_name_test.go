package usage

import "testing"

func TestCommandNameStopsAtTheFirstOperandFlagOrGroup(t *testing.T) {
	for _, tc := range []struct{ line, want string }{
		{WorktreeShow, "bench worktree show"},
		{WorktreeLand, "bench worktree land"},
		{WorktreeReset, "bench worktree reset"},
		{WorktreeMerge, "bench worktree merge"},
		{WorktreeList, "bench worktree list"},
	} {
		if got := CommandName(tc.line); got != tc.want {
			t.Errorf("CommandName(%q) = %q, want %q", tc.line, got, tc.want)
		}
	}
}
