package reviewrecord_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

func TestReviewRecordSource(t *testing.T) {
	f := recordtest.New(t, 2)
	f.AddChunk()
	f.Save()
	f.Commit("retain first results")
	f.AddChunk()
	f.Complete()
	f.Save()
	f.Commit("retain final results")
	loaded, err := rr.Read(f.Root, recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	if err := rr.CheckSource(f.Root, f.Tree(), f.Tip(), loaded, "", true); err != nil {
		t.Fatalf("continuous chain: %v", err)
	}
	omitted := loaded
	omitted.Chunks = omitted.Chunks[1:]
	if err := rr.CheckSource(f.Root, f.Tree(), f.Tip(), omitted, "2", false); err == nil || !strings.Contains(err.Error(), "missing planned chunk 1") {
		t.Fatalf("checkpoint accepted omitted predecessor: %v", err)
	}
	first := loaded.Chunks[0]
	loaded.Chunks[1].Base = first.Base
	if err := rr.CheckSource(f.Root, f.Tree(), f.Tip(), loaded, "", true); err == nil || !strings.Contains(err.Error(), "chain gap") {
		t.Fatalf("gap accepted: %v", err)
	}
	digest, err := rr.SourceDigest(f.Root, f.Tree(), recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	f.Record.Chunks[1].Reviews[0].NativeRef = recordtest.Native("completed and inspected without a local log")
	f.Save()
	f.Commit("append native excerpt")
	after, err := rr.SourceDigest(f.Root, f.Tree(), recordtest.Spec)
	if err != nil || after != digest {
		t.Fatalf("record-only change staled source: %s %s %v", digest, after, err)
	}
	loaded, err = rr.Read(f.Root, recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(loaded.Chunks[1].Reviews[0].NativeRef.Excerpt, "without a local log") {
		t.Fatal("native result not retained")
	}
	ticket, err := os.ReadFile(filepath.Join(f.Root, "specs/example/tickets/2.md"))
	if err != nil {
		t.Fatal(err)
	}
	f.Write("specs/example/tickets/2.md", strings.Replace(string(ticket), "Covers: E2", "Covers: E20", 1))
	f.Commit("unreviewed plan")
	if err := rr.CheckSource(f.Root, f.Tree(), f.Tip(), loaded, "", true); err == nil || !strings.Contains(err.Error(), "stale plan digest") {
		t.Fatalf("unreviewed valid plan accepted: %v", err)
	}
}

func TestReviewRecordTerminal(t *testing.T) {
	f := recordtest.New(t, 1)
	f.AddChunk()
	f.Complete()
	data, _ := json.Marshal(f.Record)
	if _, err := rr.Parse(data); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name, want string
		mutate     func(*rr.Record)
	}{
		{"missing completion source", "terminal completion", func(r *rr.Record) { r.Completion.SourceDigest = "" }},
		{"missing completion performer", "terminal completion", func(r *rr.Record) { r.Completion.Performer = "" }},
		{"author as reviewer", "independent review", func(r *rr.Record) { r.Chunks[0].Reviews[0].Performer = r.ImplementationSession }},
		{"verification as review", "independent review", func(r *rr.Record) { r.Chunks[0].Reviews[0].Role = "author-verification" }},
		{"missing source", "terminal source", func(r *rr.Record) { r.Chunks[0].Reviews[0].SourceDigest = "" }},
		{"missing performer", "terminal source", func(r *rr.Record) { r.Chunks[0].Reviews[0].Performer = "" }},
		{"duplicate id", "duplicate", func(r *rr.Record) { r.Chunks[0].Reviews[1].ID = r.Chunks[0].Reviews[0].ID }},
		{"computed state stored", "occurrence state", func(r *rr.Record) { r.Chunks[0].Reviews[0].State = "current" }},
		{"no embedded result", "native result", func(r *rr.Record) { r.Chunks[0].Reviews[0].NativeRef.Excerpt = "" }},
		{"bad result digest", "native result", func(r *rr.Record) { r.Chunks[0].Reviews[0].NativeRef.Digest = "bad" }},
		{"version before hostile reference", "unsupported version", func(r *rr.Record) { r.Version = 99; r.Chunks[0].Reviews[0].NativeRef.Ref = "../../unsafe" }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var r rr.Record
			_ = json.Unmarshal(data, &r)
			tc.mutate(&r)
			bad, _ := json.Marshal(r)
			if _, err := rr.Parse(bad); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("%s: %v", tc.name, err)
			}
		})
	}
	prior := f.Record.Chunks[0].Reviews[0]
	f.Record.Chunks[0].Reviews[0].Outcome = "findings"
	f.Record.Chunks[0].Reviews[0].FindingIDs = []string{"S1"}
	prior.ID = "standards-repair"
	prior.Supersedes = []string{f.Record.Chunks[0].Reviews[0].ID}
	f.Record.Chunks[0].Reviews = append(f.Record.Chunks[0].Reviews, prior)
	data, _ = json.Marshal(f.Record)
	loaded, err := rr.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Chunks[0].Reviews) != 4 || loaded.Chunks[0].Reviews[0].FindingIDs[0] != "S1" {
		t.Fatal("supersession erased earlier findings")
	}
	if err := rr.CheckReviews(loaded.Chunks[0], []string{loaded.ImplementationSession}, false); err != nil {
		t.Fatal(err)
	}
}

func TestReviewRecordPaths(t *testing.T) {
	for _, kind := range []string{"absent", "directory", "fifo", "symlink", "parent symlink", "oversized", "empty", "duplicate fence", "duplicate key"} {
		t.Run(kind, func(t *testing.T) {
			f := recordtest.New(t, 1)
			f.AddChunk()
			f.Save()
			file := filepath.Join(f.Root, "reviews/example.md")
			switch kind {
			case "absent":
				if err := os.RemoveAll(filepath.Dir(file)); err != nil {
					t.Fatal(err)
				}
			case "directory":
				_ = os.Remove(file)
				if err := os.Mkdir(file, 0o755); err != nil {
					t.Fatal(err)
				}
			case "fifo":
				_ = os.Remove(file)
				if err := syscall.Mkfifo(file, 0o600); err != nil {
					t.Fatal(err)
				}
			case "symlink":
				_ = os.Remove(file)
				if err := os.Symlink("missing", file); err != nil {
					t.Fatal(err)
				}
			case "parent symlink":
				_ = os.RemoveAll(filepath.Dir(file))
				if err := os.Symlink(t.TempDir(), filepath.Dir(file)); err != nil {
					t.Fatal(err)
				}
			case "oversized":
				f.Write("reviews/example.md", strings.Repeat("x", 3<<20))
			case "empty":
				f.Write("reviews/example.md", "")
			case "duplicate fence":
				data, _ := os.ReadFile(file)
				f.Write("reviews/example.md", string(data)+string(data))
			case "duplicate key":
				f.Write("reviews/example.md", "```bench-review-record\n{\"version\":1,\"version\":1}\n```\n")
			}
			if _, err := rr.Read(f.Root, recordtest.Spec); err == nil {
				t.Fatalf("%s record accepted", kind)
			}
		})
	}
	for _, spec := range []string{"specs/../spec.md", "specs/x\n/spec.md", "reviews", "/specs/x/spec.md"} {
		if _, err := rr.RecordPath(spec); err == nil {
			t.Fatalf("unsafe path %q accepted", spec)
		}
	}
}

// TestReviewRecordSlug grades the checkpoint spec grammar's one reader. The
// record path and the checkpoint's refusal route both read it, so a valid path
// yields the slug, and a shape RecordPath refuses is refused here too.
func TestReviewRecordSlug(t *testing.T) {
	slug, err := rr.Slug(recordtest.Spec)
	if err != nil || slug != "example" {
		t.Fatalf("Slug(%q) = (%q, %v), want example", recordtest.Spec, slug, err)
	}
	for _, spec := range []string{"specs/../spec.md", "reviews"} {
		if _, err := rr.Slug(spec); err == nil {
			t.Errorf("unsafe path %q accepted", spec)
		}
	}
}

// TestReviewRecordMissingFenceNamesTheFence grades the absent-fence wording. One
// sentinel answers every fence, so the message must name the fence that is
// absent. A missing completion plan that reads as a missing review record sends
// the operator to the wrong file. The sentinel stays testable with errors.Is.
func TestReviewRecordMissingFenceNamesTheFence(t *testing.T) {
	f := recordtest.New(t, 1)
	f.Write(recordtest.Spec, "# Example\n\nStatus: staged\n")
	f.Commit("drop the completion plan")
	_, err := rr.ReadPlan(f.Root, f.Tree(), recordtest.Spec)
	if !errors.Is(err, rr.ErrMissing) || !strings.Contains(err.Error(), "missing bench-completion-plan fence") {
		t.Fatalf("plan error = %v, want the completion-plan fence named", err)
	}
	f.Write("reviews/example.md", "# Review outcomes\n\nNo fence here.\n")
	_, err = rr.Read(f.Root, recordtest.Spec)
	if !errors.Is(err, rr.ErrMissing) || !strings.Contains(err.Error(), "missing bench-review-record fence") {
		t.Fatalf("record error = %v, want the review-record fence named", err)
	}
}

func TestReviewRecordControlPath(t *testing.T) {
	for _, control := range []string{"\n", "\r", "\t", "\x00", "\x1f", "\x7f"} {
		spec := "specs/x" + control + "/spec.md"
		if _, err := rr.RecordPath(spec); err == nil {
			t.Errorf("control path accepted: %q", spec)
		}
	}
	f := recordtest.New(t, 1)
	f.AddChunk()
	f.Record.Spec = "specs/x\n/spec.md"
	data, err := json.Marshal(f.Record)
	if err != nil {
		t.Fatal(err)
	}
	f.Write("reviews/x\n.md", "```bench-review-record\n"+string(data)+"\n```\n")
	if _, err := rr.Read(f.Root, f.Record.Spec); err == nil || !strings.Contains(err.Error(), "invalid checkpoint spec path") {
		t.Fatalf("existing control path accepted: %v", err)
	}
}

func TestReviewRecordLiteralTickets(t *testing.T) {
	f := recordtest.New(t, 1)
	f.Plan.Chunks[0].Tickets = []string{"a*.md", "ab.md"}
	for i, name := range f.Plan.Chunks[0].Tickets {
		f.Write("specs/example/tickets/"+name, fmt.Sprintf("# Literal ticket\n\nBlocked by: none\nWrites: source.txt\nCovers: E%d\n\n## What to build\n\nA behavior.\n\n## Acceptance\n\n- [ ] Works.\n", i+1))
	}
	data, _ := json.Marshal(f.Plan)
	f.Write(recordtest.Spec, "# Example\n\nStatus: staged\n\n```bench-completion-plan\n"+string(data)+"\n```\n")
	f.Commit("literal ticket names")
	plan, err := rr.ReadPlan(f.Root, f.Tree(), recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Chunks[0].Rows) != 2 {
		t.Fatalf("lost literal ticket: %+v", plan.Chunks)
	}
}

func TestReviewRecordPlanAmendment(t *testing.T) {
	f := recordtest.New(t, 2)
	f.AddChunk()
	f.Save()
	f.Commit("first review")
	base := f.Tip()
	old := f.Plan.Digest
	ticket, err := os.ReadFile(filepath.Join(f.Root, "specs/example/tickets/2.md"))
	if err != nil {
		t.Fatal(err)
	}
	f.Write("specs/example/tickets/2.md", strings.Replace(string(ticket), "Covers: E2", "Covers: E20", 1))
	f.Commit("valid amended plan")
	f.Plan, err = rr.ReadPlan(f.Root, f.Tree(), recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	if err := rr.CheckSource(f.Root, f.Tree(), f.Tip(), f.Record, "1", false); err == nil || !strings.Contains(err.Error(), "stale plan digest") {
		t.Fatalf("unreviewed amendment: %v", err)
	}
	f.Record.PlanDigest = f.Plan.Digest
	f.AddChunk()
	f.Record.Chunks[1].Base = base
	for i := range f.Record.Chunks[1].Reviews {
		f.Record.Chunks[1].Reviews[i].Base = base
	}
	if err := rr.CheckSource(f.Root, f.Tree(), f.Tip(), f.Record, "2", false); err == nil || !strings.Contains(err.Error(), "explicit old-to-new") {
		t.Fatalf("unmapped amendment: %v", err)
	}
	f.Record.Amendments = []rr.Amendment{{From: old, To: f.Plan.Digest, ChunkIDs: map[string][]string{"1": {"1"}, "2": {"2"}}}}
	if err := rr.CheckSource(f.Root, f.Tree(), f.Tip(), f.Record, "2", false); err != nil {
		t.Fatalf("reviewed mapped amendment: %v", err)
	}
}
