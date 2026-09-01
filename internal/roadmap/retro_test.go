package roadmap

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/retros"
)

const retroBody = "## Outcome\n\nLand the ticket.\n\n## Gate-stage timings\n\n- test: 1s\n\n## Ticket-versus-spec-slice and delegate performance\n\nOne ticket.\n\n## Coordinator catches\n\nNone.\n\n## Repair attribution\n\n| ticket | rounds | causes |\n|---|---|---|\n| retro | 1 | none |\n\n## Agent-experience improvements\n\n### Bench CLI\n\n- Add the writer.\n  Feeds: none\n\n### Skills\n\n### Process\n"

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
	out, code := RetroCommand([]string{"first", "--body", retroBody})
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
	for _, slug := range []string{"first", "second"} {
		if _, code := RetroCommand([]string{slug, "--body", retroBody}); code != 0 {
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
	if out, code := RetroCommand([]string{"primary-local", "--body", retroBody}); code != 0 || out != "captured: primary-local\n" {
		t.Fatalf("ignored retro = %q/%d, want captured on exit 0", out, code)
	}
	if _, err := os.Stat(filepath.Join(linked, retros.Path("primary-local"))); !os.IsNotExist(err) {
		t.Fatalf("retro landed in linked worktree: %v", err)
	}
	if _, err := os.Stat(filepath.Join(primary, retros.Path("primary-local"))); err != nil {
		t.Fatalf("retro did not land in primary checkout: %v", err)
	}
}
