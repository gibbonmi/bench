package gitguard

import "testing"

// TestGitProjectionAuthority grades the Git policy around the shared exec projection: a child
// in another checkout and a redirected repository leave push unresolved, and a non-Bench head
// projects no child.
func TestGitProjectionAuthority(t *testing.T) {
	cases := []struct {
		name string
		cmd  string
		chk  Checker
		want string
	}{
		{"child push", "bench worktree exec X -- git push origin topic", pushMain, "git push with an unresolved destination"},
		{"redirected push", "git -C /other push origin topic", pushMain, "git push with an unresolved destination"},
		{"non-Bench exec head", "mybench worktree exec X -- git reset --hard", refYes, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Classify(tc.cmd, tc.chk); got != tc.want {
				t.Errorf("Classify(%q) = %q, want %q", tc.cmd, got, tc.want)
			}
		})
	}
}

// TestGuardWrapperSelection grades that Git reads only the first matching shell child.
func TestGuardWrapperSelection(t *testing.T) {
	cases := []struct {
		name string
		cmd  string
		want string
	}{
		{"later reset child", "bash -c 'echo ok' -c 'git reset --hard'", ""},
		{"first reset child", "bash -c 'git reset --hard' -c 'echo ok'", "git reset"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Classify(tc.cmd, refYes); got != tc.want {
				t.Errorf("Classify(%q) = %q, want %q", tc.cmd, got, tc.want)
			}
		})
	}
}
