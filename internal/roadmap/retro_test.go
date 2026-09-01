package roadmap

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/retros"
	retrotestdata "github.com/gibbonmi/bench/internal/retros/testdata"
)

func eligibleRetro(t *testing.T) string {
	t.Helper()
	return retrotestdata.Eligible()
}

func TestRetroRejectsMalformedBodyWithoutWriting(t *testing.T) {
	root := newRepo(t)
	_, code := RetroCommand([]string{"broken", "--body", "## Outcome"})
	if code == 0 {
		t.Fatal("malformed retrospective exit = 0, want a refusal")
	}
	if _, err := os.Stat(filepath.Join(root, retros.Path("broken"))); !os.IsNotExist(err) {
		t.Fatalf("malformed retrospective created a file: %v", err)
	}
}

func TestRetroRoundTripsOneEligibleArtifact(t *testing.T) {
	root := newRepo(t)
	out, code := RetroCommand([]string{"first", "--body", eligibleRetro(t)})
	if code != 0 || out != "captured: first\n" {
		t.Fatalf("retro = %q/%d, want captured on exit 0", out, code)
	}
	facts := retros.Facts(root)
	if facts.State.Failed() || len(facts.Entries) != 1 {
		t.Fatalf("facts = %#v, want one eligible retrospective", facts)
	}
	if err := retros.Parse(facts.Entries[0].Body); err != nil {
		t.Fatalf("%s does not parse after write: %v", facts.Entries[0].Path, err)
	}

}

func TestRetroRepeatedWritesPreserveEarlierCapture(t *testing.T) {
	root := newRepo(t)
	body := eligibleRetro(t)
	for _, slug := range []string{"first", "second"} {
		if _, code := RetroCommand([]string{slug, "--body", body}); code != 0 {
			t.Fatalf("retro %q exit = %d, want 0", slug, code)
		}
	}
	facts := retros.Facts(root)
	if facts.State.Failed() || len(facts.Entries) != 2 {
		t.Fatalf("facts = %#v, want two eligible retrospectives", facts)
	}
	if !strings.Contains(string(facts.Entries[0].Body), "Land the ticket.") {
		t.Fatalf("first retrospective body = %q", facts.Entries[0].Body)
	}
}

func TestRetroRefusesSameSlugWithoutChangingEarlierCapture(t *testing.T) {
	root := newRepo(t)
	body := eligibleRetro(t)
	if _, code := RetroCommand([]string{"repeat", "--body", body}); code != 0 {
		t.Fatalf("first retro exit = %d, want 0", code)
	}
	path := filepath.Join(root, retros.Path("repeat"))
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	secondBody := strings.Replace(body, "Land the ticket.", "Preserve the first capture.", 1)
	if _, code := RetroCommand([]string{"repeat", "--body", secondBody}); code == 0 {
		t.Fatal("second retro exit = 0, want a refusal")
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Fatalf("retrospective changed after refusal:\n got %q\nwant %q", after, before)
	}
}

func TestRetroIgnoredDirectoryWritesToPrimaryCheckout(t *testing.T) {
	primary := resolvedToplevel(t, newPrimaryRepo(t))
	if err := os.WriteFile(filepath.Join(primary, ".gitignore"), []byte(retros.Directory+"/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command("git", "-C", primary, "add", ".gitignore").CombinedOutput(); err != nil {
		t.Fatalf("git add: %v: %s", err, out)
	}
	if out, err := exec.Command("git", "-C", primary, "commit", "-q", "-m", "ignore retros").CombinedOutput(); err != nil {
		t.Fatalf("git commit: %v: %s", err, out)
	}
	linked := newLinkedWorktree(t, primary)
	t.Chdir(linked)
	if out, code := RetroCommand([]string{"primary-local", "--body", eligibleRetro(t)}); code != 0 || out != "captured: primary-local\n" {
		t.Fatalf("ignored retro = %q/%d, want captured on exit 0", out, code)
	}
	if _, err := os.Stat(filepath.Join(linked, retros.Path("primary-local"))); !os.IsNotExist(err) {
		t.Fatalf("retro landed in linked worktree: %v", err)
	}
	if _, err := os.Stat(filepath.Join(primary, retros.Path("primary-local"))); err != nil {
		t.Fatalf("retro did not land in primary checkout: %v", err)
	}
}

func TestRetroRefusesSymlinkedDestinationComponents(t *testing.T) {
	for _, target := range []string{"outside", "missing"} {
		t.Run(target, func(t *testing.T) {
			root := newRepo(t)
			outside := filepath.Join(t.TempDir(), target)
			if target == "outside" {
				if err := os.Mkdir(outside, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			if err := os.Symlink(outside, filepath.Join(root, retros.Directory)); err != nil {
				t.Fatal(err)
			}

			out, code := RetroCommand([]string{"linked", "--body", eligibleRetro(t)})
			if code == 0 || !strings.Contains(out, "symbolic link") {
				t.Fatalf("retro through %s symlink = %q/%d, want a symbolic-link refusal", target, out, code)
			}
			if _, err := os.Stat(filepath.Join(outside, "linked.md")); !os.IsNotExist(err) {
				t.Fatalf("retro wrote through %s symlink: %v", target, err)
			}
		})
	}
}

func TestRetroContainsDestinationComponentReplacement(t *testing.T) {
	root := newRepo(t)
	outside := t.TempDir()
	originalOpenRoot := openRetroRoot
	t.Cleanup(func() { openRetroRoot = originalOpenRoot })
	openRetroRoot = func(name string) (*os.Root, error) {
		opened, err := os.OpenRoot(name)
		if err != nil {
			return nil, err
		}
		capture := filepath.Join(root, "capture")
		if err := os.Remove(capture); err != nil {
			_ = opened.Close()
			t.Fatalf("remove capture before replacement: %v", err)
		}
		if err := os.Symlink(outside, capture); err != nil {
			_ = opened.Close()
			t.Fatalf("replace capture with symlink: %v", err)
		}
		target, err := os.Readlink(capture)
		if err != nil {
			_ = opened.Close()
			t.Fatalf("read capture replacement: %v", err)
		}
		if target != outside {
			_ = opened.Close()
			t.Fatalf("capture replacement target = %q, want %q", target, outside)
		}
		return opened, nil
	}

	out, code := RetroCommand([]string{"raced", "--body", eligibleRetro(t)})
	if code == 0 {
		t.Fatalf("retro after destination replacement = %q/%d, want a refusal", out, code)
	}
	if _, err := os.Stat(filepath.Join(outside, "retrospectives", "raced.md")); !os.IsNotExist(err) {
		t.Fatalf("retro escaped through replaced component: %v", err)
	}
}
