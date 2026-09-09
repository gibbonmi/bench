package preflight

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/diff"
)

func assertProposalRefusal(t *testing.T, out string, code int, want string) {
	t.Helper()
	if code == 0 || !strings.Contains(out, want) || strings.Contains(out, "writes_proposal[") {
		t.Fatalf("proposal refusal = (%d), want %q and no proposal:\n%s", code, want, out)
	}
}

func advanceProposalFromCurrent(t *testing.T) {
	t.Helper()
	runGit(t, "update-ref", "refs/heads/main", "HEAD")
	mustWriteFile(t, "internal/example/after.go", "package example\n")
	runGit(t, "add", "internal/example/after.go")
	runGit(t, "commit", "-q", "-m", "advance proposal source")
}

func TestWritesProposalGrammar(t *testing.T) {
	base := []string{"build", "example", "--propose-writes", "--ticket", "one.md", "--base", "base", "--source-tip", "tip"}
	cases := []struct {
		name string
		args []string
		want string
	}{
		{"duplicate", append(append([]string{}, base...), "--propose-writes"), "--propose-writes"},
		{"unknown", append(append([]string{}, base...), "--proposal"), "--proposal"},
		{"missing ticket value", []string{"build", "example", "--propose-writes", "--ticket"}, "--ticket"},
		{"missing base value", []string{"build", "example", "--propose-writes", "--ticket", "one.md", "--base"}, "--base"},
		{"missing source tip value", base[:len(base)-1], "--source-tip"},
		{"review mode", append([]string{"review"}, base[1:]...), "--propose-writes requires build"},
		{"combined modes", append(append([]string{}, base...), "--charge"), "cannot be combined"},
		{"full unsupported", append(append([]string{}, base...), "--full"), "--full requires --charge"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			out, code := Command(test.args)
			if code != 2 || !strings.Contains(out, test.want) || strings.Contains(out, "writes_proposal[") {
				t.Fatalf("grammar = (%d), want %q:\n%s", code, test.want, out)
			}
		})
	}

	t.Run("ticket traversal", func(t *testing.T) {
		root, slug := seedProposal(t)
		args := proposalArgs(t, root, slug)
		args[4] = "../one.md"
		out, code := Command(args)
		assertProposalRefusal(t, out, code, "selected ticket")
	})
}

func TestWritesProposalSnapshotRefusals(t *testing.T) {
	t.Run("pin mismatch", func(t *testing.T) {
		root, slug := seedProposal(t)
		args := proposalArgs(t, root, slug)
		args[8] = runGit(t, "rev-parse", "main")
		out, code := Command(args)
		assertProposalRefusal(t, out, code, "tip-current")
	})

	t.Run("dirty checkout", func(t *testing.T) {
		root, slug := seedProposal(t)
		args := proposalArgs(t, root, slug)
		mustWriteFile(t, "dirty.txt", "dirty\n")
		out, code := Command(args)
		assertProposalRefusal(t, out, code, "checkout required")
	})

	for _, test := range []struct {
		name   string
		mutate func(*testing.T, int)
	}{
		{"head movement", func(t *testing.T, n int) {
			mustWriteFile(t, "head-move.txt", string(rune('0'+n))+"\n")
			runGit(t, "add", "head-move.txt")
			runGit(t, "commit", "-q", "-m", "move head")
		}},
		{"index movement", func(t *testing.T, n int) {
			mustWriteFile(t, "index-move.txt", string(rune('0'+n))+"\n")
			runGit(t, "add", "index-move.txt")
		}},
		{"required byte movement", func(t *testing.T, n int) {
			mustWriteFile(t, buildPhase, "# moved "+string(rune('0'+n))+"\n")
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, slug := seedProposal(t)
			args := proposalArgs(t, root, slug)
			calls := 0
			restore := diff.SetSnapshotAfterReadForTest(func() {
				calls++
				test.mutate(t, calls)
			})
			out, code := Command(args)
			restore()
			if calls != 2 {
				t.Fatalf("movement callbacks = %d, want 2", calls)
			}
			assertProposalRefusal(t, out, code, "snapshot drift")
		})
	}

	t.Run("retry ends in early refusal", func(t *testing.T) {
		root, slug := seedProposal(t)
		args := proposalArgs(t, root, slug)
		calls := 0
		restore := diff.SetSnapshotAfterReadForTest(func() {
			calls++
			if calls == 1 {
				path := "specs/" + slug + "/spec.md"
				mustWriteFile(t, path, strings.Replace(specBody(slug), "Status: staged", "Status: draft", 1))
			}
		})
		out, code := Command(args)
		restore()
		if calls != 2 {
			t.Fatalf("retry callbacks = %d, want 2", calls)
		}
		assertProposalRefusal(t, out, code, "spec not staged")
	})

	t.Run("final recapture error discards output", func(t *testing.T) {
		root, slug := seedProposal(t)
		args := proposalArgs(t, root, slug)
		index := strings.TrimSpace(runGit(t, "rev-parse", "--git-path", "index"))
		restore := diff.SetSnapshotAfterReadForTest(func() {
			if err := os.Remove(index); err != nil {
				t.Fatal(err)
			}
		})
		out, code := Command(args)
		restore()
		assertProposalRefusal(t, out, code, "snapshot identity failed")
	})
}

func TestWritesProposalRequiredInputs(t *testing.T) {
	t.Run("ignored required source absent at tip", func(t *testing.T) {
		root, slug := seedProposal(t)
		if err := os.Remove(buildPhase); err != nil {
			t.Fatal(err)
		}
		mustWriteFile(t, ".gitignore", buildPhase+"\n")
		runGit(t, "add", "-A")
		runGit(t, "commit", "-q", "-m", "remove required source")
		mustWriteFile(t, buildPhase, "# ignored live source\n")
		advanceProposalFromCurrent(t)
		out, code := Command(proposalArgs(t, root, slug))
		assertProposalRefusal(t, out, code, buildPhase+" is absent or unreadable at source tip")
	})

	t.Run("required source differs from tip", func(t *testing.T) {
		root, slug := seedProposal(t)
		runGit(t, "update-index", "--assume-unchanged", buildPhase)
		mustWriteFile(t, buildPhase, "# live drift\n")
		out, code := Command(proposalArgs(t, root, slug))
		assertProposalRefusal(t, out, code, buildPhase+" does not match source tip")
	})

	t.Run("unreadable fixture inventory", func(t *testing.T) {
		root, slug := seedProposal(t)
		path := "tests/canary/example-family/pinning-fixture/BASE"
		if err := os.Chmod(path, 0); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chmod(path, 0o644) })
		out, code := Command(proposalArgs(t, root, slug))
		assertProposalRefusal(t, out, code, "fixture")
	})

	for _, test := range []struct {
		name    string
		prepare func(*testing.T)
	}{
		{"live link", func(t *testing.T) {
			if err := os.Remove(buildPhase); err != nil {
				t.Fatal(err)
			}
			mustWriteFile(t, "phase-target.md", "# phase\n")
			if err := os.Symlink("../../phase-target.md", buildPhase); err != nil {
				t.Fatal(err)
			}
		}},
		{"dangling link", func(t *testing.T) {
			if err := os.Remove(buildPhase); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink("missing", buildPhase); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, slug := seedProposal(t)
			test.prepare(t)
			runGit(t, "add", "-A")
			runGit(t, "commit", "-q", "-m", "special source")
			advanceProposalFromCurrent(t)
			out, code := Command(proposalArgs(t, root, slug))
			assertProposalRefusal(t, out, code, buildPhase)
		})
	}

	t.Run("fifo", func(t *testing.T) {
		root, slug := seedProposal(t)
		if err := os.Remove(buildPhase); err != nil {
			t.Fatal(err)
		}
		if err := syscall.Mkfifo(filepath.Clean(buildPhase), 0o600); err != nil {
			t.Fatal(err)
		}
		out, code := Command(proposalArgs(t, root, slug))
		assertProposalRefusal(t, out, code, "checkout required")
	})
}
