package handoff

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/handoffdoc"
	"github.com/gibbonmi/bench/internal/status"
)

// LC21, story 22. The flag's whole point: the file's bytes become the owned section's
// State, and the words the document held go. The file ends in the newline an editor adds,
// which the document's own writer trims, so the written body carries neither.
func TestHandoffStateFileWritesTheState(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	seedState(t, document, "The earlier words.")
	file := stateFile(t, "The build is live.\n\nTicket 3 is diff-ready.\n")

	runIn(t, root, []string{"--state-file", file})

	if body := stateBody(t, document); body != "The build is live.\n\nTicket 3 is diff-ready." {
		t.Fatalf("State body = %q, want the file's bytes under one trimmed newline", body)
	}
}

// LC22, stories 23 and 24. A drafted State gets the same scan a hand-written one gets. An
// off-ancestry pin is a stale resume target whichever way it arrived, and the refusal comes
// before the render, so the document keeps the bytes the run read.
func TestHandoffStateFileRefusesAnOffAncestryPin(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	off := offAncestryCommit(t, root)
	before := seedState(t, document, "The reviewer's own words.")
	file := stateFile(t, "The build resumes from `"+off+"`.\n")

	out, code := runAt(t, root, []string{"--state-file", file})
	if code != 1 {
		t.Fatalf("handoff over a drafted off-ancestry pin = (%q, %d), want exit 1", out, code)
	}
	if !strings.Contains(out, faultOffAncestry) {
		t.Errorf("the refusal does not give the ancestry reason %q\n%s", faultOffAncestry, out)
	}
	if !strings.Contains(out, "The build resumes from `"+off+"`.") {
		t.Errorf("the refusal does not print the offending line\n%s", out)
	}
	if after := read(t, document); after != before {
		t.Fatalf("a refused run rewrote the document\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// LC23, story 25. A tree hash resolves in the object store and is still no resume target.
// The reason is asserted rather than the exit code, because the repair the writer owes for
// a lost commit and for an object that was never one are different.
func TestHandoffStateFileRefusesANonCommitPin(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	seedState(t, document, "The reviewer's own words.")
	file := stateFile(t, "The tree is `"+gitOut(t, root, "rev-parse", "HEAD^{tree}")+"`.\n")

	out, code := runAt(t, root, []string{"--state-file", file})
	if code != 1 {
		t.Fatalf("handoff over a drafted tree hash = (%q, %d), want exit 1", out, code)
	}
	if !strings.Contains(out, faultNotACommit) {
		t.Errorf("the refusal does not give the not-a-commit reason %q\n%s", faultNotACommit, out)
	}
	if strings.Contains(out, faultOffAncestry) {
		t.Errorf("a tree hash drew the ancestry reason, so the commit peel never ran\n%s", out)
	}
}

// LC24, story 26. An abbreviation that expands to two objects resolves to neither, so both
// object probes fail and the token reads as prose under the exit code alone.
func TestHandoffStateFileRefusesAnAmbiguousPin(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	seedState(t, document, "The reviewer's own words.")
	prefix := ambiguousPrefix(t, root)
	file := stateFile(t, "Resume from `"+prefix+"`.\n")

	out, code := runAt(t, root, []string{"--state-file", file})
	if code != 1 {
		t.Fatalf("handoff over a drafted %q = (%q, %d), want exit 1", prefix, out, code)
	}
	if !strings.Contains(out, faultAmbiguous) {
		t.Errorf("the refusal does not give the ambiguity reason %q\n%s", faultAmbiguous, out)
	}
}

// LC25, story 27. The four file kinds a plain read would get wrong. An absent draft read as
// an empty State would erase the section's State on a typo, and a link or a special file
// would send the read outside the file the caller named. Each refuses, and each leaves the
// document alone.
func TestHandoffStateFileRefusesAnUnreadableFile(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	before := seedState(t, document, "The reviewer's own words.")
	dir := t.TempDir()
	target := filepath.Join(dir, "draft.md")
	write(t, target, "The build is live.\n")

	for _, tc := range []struct {
		name string
		path func(*testing.T) string
	}{
		{"absent", func(*testing.T) string {
			return filepath.Join(dir, "missing.md")
		}},
		{"symbolic link", func(t *testing.T) string {
			link := filepath.Join(dir, "link.md")
			if err := os.Symlink(target, link); err != nil {
				t.Fatalf("Symlink: %v", err)
			}
			return link
		}},
		{"named pipe", func(t *testing.T) string {
			pipe := filepath.Join(dir, "pipe.md")
			if err := syscall.Mkfifo(pipe, 0o600); err != nil {
				t.Fatalf("Mkfifo: %v", err)
			}
			return pipe
		}},
		{"unreadable mode", func(t *testing.T) string {
			if os.Geteuid() == 0 {
				capability.Capability(t, capability.Privilege, "root reads 0000-mode files; the permission case is unobservable")
			}
			locked := filepath.Join(dir, "locked.md")
			write(t, locked, "The build is live.\n")
			if err := os.Chmod(locked, 0o000); err != nil {
				t.Fatalf("Chmod: %v", err)
			}
			return locked
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, code := runAt(t, root, []string{"--state-file", tc.path(t)})
			if code != 1 {
				t.Fatalf("handoff over a %s state file = (%q, %d), want exit 1", tc.name, out, code)
			}
			if !strings.Contains(out, faultStateFileUnreadable) {
				t.Errorf("the refusal does not give the read reason %q\n%s", faultStateFileUnreadable, out)
			}
			if after := read(t, document); after != before {
				t.Fatalf("a refused run rewrote the document\nbefore:\n%s\nafter:\n%s", before, after)
			}
		})
	}
}

// LC26, story 28. A control byte in the draft would ride into every downstream reader of
// the artifact, so it refuses at the same predicate the derived fields answer to.
func TestHandoffStateFileRefusesAControlByte(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	before := seedState(t, document, "The reviewer's own words.")
	file := stateFile(t, "The gate is \x1b[31mred\x1b[0m.\n")

	out, code := runAt(t, root, []string{"--state-file", file})
	if code != 1 {
		t.Fatalf("handoff over a drafted escape byte = (%q, %d), want exit 1", out, code)
	}
	if !strings.Contains(out, faultStateFileByte) {
		t.Errorf("the refusal does not give the control-byte reason %q\n%s", faultStateFileByte, out)
	}
	if after := read(t, document); after != before {
		t.Fatalf("a refused run rewrote the document\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// LC27, story 29. A level-two heading in the draft opens a block below the State it was
// written into, and the document the next run reads is then a shape this grammar cannot
// name. The refusal happens while the writer still holds the text.
func TestHandoffStateFileRefusesASectionHeading(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	for _, heading := range []string{handoffdoc.MainHeading, "## Open questions"} {
		t.Run(heading, func(t *testing.T) {
			before := seedState(t, document, "The reviewer's own words.")
			file := stateFile(t, "The build is live.\n\n"+heading+"\n\nNone.\n")

			out, code := runAt(t, root, []string{"--state-file", file})
			if code != 1 {
				t.Fatalf("handoff over a drafted %q = (%q, %d), want exit 1", heading, out, code)
			}
			if !strings.Contains(out, heading) {
				t.Errorf("the refusal does not print the offending line\n%s", out)
			}
			if after := read(t, document); after != before {
				t.Fatalf("a refused run rewrote the document\nbefore:\n%s\nafter:\n%s", before, after)
			}
			// The refusal exists so the next run still parses the file, and only a second
			// run states that. A written heading makes that run refuse the document.
			runIn(t, root, nil)
		})
	}

	// A heading inside a fence is an example the writer pasted, not a block the parser
	// opens, so the same rule that skips a fenced pin skips this line.
	fenced := stateFile(t, "The earlier close wrote:\n\n```\n"+handoffdoc.MainHeading+"\n```\n")
	runIn(t, root, []string{"--state-file", fenced})
}

// LC28, story 30. An empty draft is a deliberate reset rather than a failed read, so it
// writes an empty State and the header carries the first-session guidance again.
func TestHandoffStateFileAcceptsAnEmptyFile(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	seedState(t, document, "The earlier words.")
	file := stateFile(t, "")

	runIn(t, root, []string{"--state-file", file})

	if body := stateBody(t, document); body != "" {
		t.Fatalf("State body = %q, want an empty State", body)
	}
	if after := read(t, document); !strings.Contains(after, scaffoldGuidance) {
		t.Errorf("an emptied State does not bring the scaffold guidance back\n%s", after)
	}
}

// LC29, story 31. Without the flag the State is the reviewer's, byte for byte. The fenced
// block is the sharp case: a form that recomposed State rather than re-emitting it would
// reflow the pasted lines.
func TestCommandKeepsTheOwnedStateWithoutTheFlag(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	const state = "The reviewer's own words.\n\n```console\n$ bench gate\ngate: green\n```"
	seedState(t, document, state)

	runIn(t, root, nil)

	if body := stateBody(t, document); body != state {
		t.Fatalf("State body = %q, want the document's own bytes %q", body, state)
	}
}

// stateFile writes one drafted State body outside the checkout and returns its path, which
// is where a coordinator drafts one.
func stateFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "state.md")
	write(t, path, content)
	return path
}

// stateBody returns the owned section's State as the document holds it. It reads through
// the section split the other fixtures use, so a change to the rendered shape moves the
// expectation with it.
func stateBody(t *testing.T, document string) string {
	t.Helper()
	_, after, found := strings.Cut(sectionBytes(t, document, handoffdoc.MainKey), handoffdoc.StateHeading+"\n")
	if !found {
		t.Fatalf("the document carries no %q heading\n%s", handoffdoc.StateHeading, read(t, document))
	}
	return strings.Trim(after, "\n")
}

// seedSection plants main's reviewer-owned fields and returns the document's bytes, so a
// refusal test can compare against exactly what the refused run read. It goes in through
// the leaf package's own writer, so the planted bytes are the ones a real run would parse.
func seedSection(t *testing.T, document, next, state string) string {
	t.Helper()
	if err := handoffdoc.WriteSection(document, handoffdoc.Section{Key: handoffdoc.MainKey, Next: next, State: state}); err != nil {
		t.Fatalf("seed section: %v", err)
	}
	return read(t, document)
}

func seedState(t *testing.T, document, state string) string {
	t.Helper()
	return seedSection(t, document, "", state)
}
