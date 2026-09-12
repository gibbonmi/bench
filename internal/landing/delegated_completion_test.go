package landing

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/spec"
)

// newDelegatedFixture mirrors newCompletionFixture in the version 2 delegated
// form. Two chunks give the run two distinct ticket authors, and the
// orchestrator owns the final obligations.
func newDelegatedFixture(t *testing.T) completionFixture {
	t.Helper()
	f := recordtest.AttachDelegated(t, fixture(t), 2)
	f.Write(".gitignore", ".logs/\n")
	f.Write(".bench/gate.sh", "#!/bin/sh\ngrep -q '^Status: implemented$' specs/example/spec.md\n")
	if err := os.Chmod(filepath.Join(f.Root, ".bench/gate.sh"), 0755); err != nil {
		t.Fatal(err)
	}
	f.Write(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`+"\n")
	f.Commit("landing oracle")
	root, base := f.Root, f.Tip()
	source := filepath.Join(t.TempDir(), "source")
	f.Git("worktree", "add", "-qb", "evidence-source", source, base)
	f.Root = source
	f.AddChunk()
	f.Save()
	f.Commit("retain first chunk")
	f.AddChunk()
	f.Complete()
	f.Save()
	f.Commit("retain completion")
	return completionFixture{f, root, base}
}

// DI8, DI9, DI10: delegated landing needs the orchestrator's final verification
// on the final source, a complete reconciliation, and a reviewed destination.
func TestDelegatedFinalVerifier(t *testing.T) {
	cases := []struct {
		name, reason string
		change       func(*testing.T, completionFixture)
	}{
		{"complete", "", func(*testing.T, completionFixture) {}},
		{"ticket author as final verifier", "verification integration", func(t *testing.T, f completionFixture) {
			f.evidence.Record.Completion.Verification[1].Performer = recordtest.Author("1.md")
			f.evidence.Save()
		}},
		{"author role on the final obligation", "verification integration", func(t *testing.T, f completionFixture) {
			f.evidence.Record.Completion.Verification[1].Role = "author-verification"
			f.evidence.Save()
		}},
		{"older final source", "verification integration", func(t *testing.T, f completionFixture) {
			f.evidence.Record.Completion.Verification[1].SourceDigest = f.evidence.Record.Chunks[0].SourceDigest
			f.evidence.Save()
		}},
		{"ticket author reconciles completion", "completion is incomplete or stale", func(t *testing.T, f completionFixture) {
			f.evidence.Record.Completion.Performer = recordtest.Author("2.md")
			f.evidence.Save()
		}},
		// DI9
		{"missing reconciliation", "reconciliation E1", func(t *testing.T, f completionFixture) {
			delete(f.evidence.Record.Completion.Reconciliation, "E1")
			f.evidence.Save()
		}},
		// DI10
		{"destination addition", "completion composition changes added.txt", func(t *testing.T, f completionFixture) {
			write(t, f.root, "added.txt", "unreviewed destination content\n")
			git(t, f.root, "add", "added.txt")
			git(t, f.root, "commit", "-qm", "destination changed")
		}},
		{"destination deletion", "completion composition changes foreign", func(t *testing.T, f completionFixture) {
			git(t, f.root, "rm", "foreign")
			git(t, f.root, "commit", "-qm", "destination removed content")
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newDelegatedFixture(t)
			tc.change(t, f)
			if f.evidence.Git("status", "--porcelain") != "" {
				f.evidence.Commit("evidence case")
			}
			request := f.request(t)
			result, err := New().LandReviewed(t.Context(), request)
			if tc.reason == "" {
				if err != nil {
					t.Fatalf("complete delegated evidence refused: %v", err)
				}
				want, err := spec.Implemented(request.SpecBytes)
				if err != nil {
					t.Fatal(err)
				}
				if got := gitBytes(t, f.root, "show", result.Commit+":"+request.SpecPath); !bytes.Equal(got, want) {
					t.Fatalf("published status transform = %q", got)
				}
				if got := git(t, f.root, "rev-parse", "main"); got != result.Commit {
					t.Fatalf("publication missing: %s", got)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.reason) {
				t.Fatalf("incomplete delegated evidence published or lost reason %q: %v", tc.reason, err)
			}
			if got := git(t, f.root, "rev-parse", "main"); got != request.DestinationBase {
				t.Fatalf("ref moved before proof: %s", got)
			}
		})
	}
}
