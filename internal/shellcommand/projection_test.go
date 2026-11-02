package shellcommand

import (
	"reflect"
	"testing"
)

// child is one projected wrapper child: its operand index and its command words.
type child struct {
	index int
	words []string
}

// childWords reduces each wrapper child to its operand index and its command words.
func childWords(children []WrapperChild) []child {
	var got []child
	for _, wrapper := range children {
		var words []string
		for _, span := range wrapper.Stream.Commands {
			words = append(words, ProjectCommandWords(wrapper.Stream.Tokens[span.Start:span.End])...)
		}
		got = append(got, child{wrapper.ArgumentIndex, words})
	}
	return got
}

func TestShellChildren(t *testing.T) {
	cases := []struct {
		name       string
		words      []string
		executable int
		want       []child
	}{
		{"sh", []string{"sh", "-c", "git status"}, 0, []child{{2, []string{"git", "status"}}}},
		{"bash", []string{"bash", "-c", "git status"}, 0, []child{{2, []string{"git", "status"}}}},
		{"zsh", []string{"zsh", "-c", "git status"}, 0, []child{{2, []string{"git", "status"}}}},
		{"absolute path and -lc", []string{"/bin/bash", "-lc", "git push"}, 0, []child{{2, []string{"git", "push"}}}},
		{"after a prefix", []string{"env", "A=1", "sh", "-c", "git log"}, 2, []child{{4, []string{"git", "log"}}}},
		{"flag-shaped operand", []string{"bash", "-c", "-c", "git status"}, 0,
			[]child{{2, []string{"-c"}}, {3, []string{"git", "status"}}}},
		{"two flags in argv order", []string{"bash", "-c", "echo ok", "-c", "git reset --hard"}, 0,
			[]child{{2, []string{"echo", "ok"}}, {4, []string{"git", "reset", "--hard"}}}},
		{"nested wrapper stays one child", []string{"bash", "-c", "sh -c 'git push'"}, 0,
			[]child{{2, []string{"sh", "-c", "git push"}}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := childWords(ShellChildren(tc.words, tc.executable)); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ShellChildren(%q, %d) = %v, want %v", tc.words, tc.executable, got, tc.want)
			}
		})
	}
}

func TestShellChildrenBounds(t *testing.T) {
	cases := []struct {
		name       string
		words      []string
		executable int
		want       []int
	}{
		{"negative index", []string{"bash", "-c", "git status"}, -1, nil},
		{"index past the words", []string{"bash", "-c", "git status"}, 3, nil},
		{"missing operand", []string{"bash", "-c"}, 0, nil},
		{"missing last operand", []string{"bash", "-c", "git status", "-lc"}, 0, []int{2}},
		{"unsupported shell", []string{"dash", "-c", "git status"}, 0, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got []int
			for _, child := range ShellChildren(tc.words, tc.executable) {
				got = append(got, child.ArgumentIndex)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ShellChildren(%q, %d) operands = %v, want %v", tc.words, tc.executable, got, tc.want)
			}
		})
	}
}

func TestIsBenchName(t *testing.T) {
	for word, want := range map[string]bool{
		"bench": true, "bench.sh": true, "./dist/bench": true, "/opt/kit/bin/bench.sh": true,
		"benchmark": false, "bench.py": false, "git": false, "": false,
	} {
		if got := IsBenchName(word); got != want {
			t.Errorf("IsBenchName(%q) = %t, want %t", word, got, want)
		}
	}
}

func TestWorktreeExec(t *testing.T) {
	cases := []struct {
		name       string
		words      []string
		executable int
		want       []string
	}{
		{"quoted argument", []string{"bench", "worktree", "exec", "X", "--env", "A=1", "--", "git", "commit", "-m", "two words"}, 0,
			[]string{"git", "commit", "-m", "two words"}},
		{"built binary after a prefix", []string{"env", "A=1", "./dist/bench", "worktree", "exec", "X", "--", "git", "status"}, 2,
			[]string{"git", "status"}},
		{"later delimiter stays in the child", []string{"bench", "worktree", "exec", "X", "--", "git", "checkout", "--", "a b"}, 0,
			[]string{"git", "checkout", "--", "a b"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := WorktreeExec(tc.words, tc.executable)
			if !got.Head || !got.HasSeparator || !reflect.DeepEqual(got.Child, tc.want) {
				t.Errorf("WorktreeExec(%q, %d) = %+v, want child %q", tc.words, tc.executable, got, tc.want)
			}
		})
	}
}

func TestWorktreeExecMissingDelimiter(t *testing.T) {
	got := WorktreeExec([]string{"bench", "worktree", "exec", "X", "git", "status"}, 0)
	if want := (ExecProjection{Head: true}); !reflect.DeepEqual(got, want) {
		t.Errorf("WorktreeExec without -- = %+v, want %+v", got, want)
	}
}

func TestWorktreeExecEmptyChild(t *testing.T) {
	got := WorktreeExec([]string{"bench", "worktree", "exec", "X", "--"}, 0)
	if !got.Head || !got.HasSeparator || len(got.Child) != 0 || got.DifferentCheckout {
		t.Errorf("WorktreeExec with an empty child = %+v, want a head, a separator, and no child context", got)
	}
}

func TestWorktreeExecContext(t *testing.T) {
	if got := WorktreeExec([]string{"bench", "worktree", "exec", "X", "--", "git", "push"}, 0); !got.DifferentCheckout {
		t.Errorf("WorktreeExec with a child = %+v, want DifferentCheckout", got)
	}
	for _, words := range [][]string{
		{"bench", "worktree", "list"},
		{"bench", "worktree"},
		{"echo", "exec", "worktree", "--", "git"},
		{"bench", "worktree", "list", "exec", "--", "git", "reset", "--hard"},
	} {
		if got := WorktreeExec(words, 0); !reflect.DeepEqual(got, ExecProjection{}) {
			t.Errorf("WorktreeExec(%q) = %+v, want no projection", words, got)
		}
	}
	if got := WorktreeExec([]string{"bench", "worktree", "exec", "X", "--", "git"}, -1); !reflect.DeepEqual(got, ExecProjection{}) {
		t.Errorf("WorktreeExec at index -1 = %+v, want no projection", got)
	}
}
