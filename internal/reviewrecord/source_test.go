package reviewrecord_test

import (
	"encoding/json"
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
	f.Write("specs/example/tickets/2.md", "# Changed plan\n")
	f.Commit("unreviewed plan")
	if err := rr.CheckSource(f.Root, f.Tree(), f.Tip(), loaded, "", true); err == nil {
		t.Fatal("unreviewed plan accepted")
	}
}

func TestReviewRecordTerminal(t *testing.T) {
	f := recordtest.New(t, 1)
	f.AddChunk()
	data, _ := json.Marshal(f.Record)
	if _, err := rr.Parse(data); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name, want string
		mutate     func(*rr.Record)
	}{
		{"author as reviewer", "independent review", func(r *rr.Record) { r.Chunks[0].Reviews[0].Performer = r.ImplementationSession }},
		{"verification as review", "independent review", func(r *rr.Record) { r.Chunks[0].Reviews[0].Role = "author-verification" }},
		{"missing source", "terminal source", func(r *rr.Record) { r.Chunks[0].Reviews[0].SourceDigest = "" }},
		{"missing performer", "terminal source", func(r *rr.Record) { r.Chunks[0].Reviews[0].Performer = "" }},
		{"duplicate id", "duplicate", func(r *rr.Record) { r.Chunks[0].Reviews[1].ID = r.Chunks[0].Reviews[0].ID }},
		{"computed state stored", "occurrence state", func(r *rr.Record) { r.Chunks[0].Reviews[0].State = "current" }},
		{"no embedded result", "native result", func(r *rr.Record) { r.Chunks[0].Reviews[0].NativeRef.Excerpt = "" }},
		{"bad result digest", "native result", func(r *rr.Record) { r.Chunks[0].Reviews[0].NativeRef.Digest = "bad" }},
		{"version before hostile reference", "unsupported version", func(r *rr.Record) { r.Version = 2; r.Chunks[0].Reviews[0].NativeRef.Ref = "../../unsafe" }},
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
	if err := rr.CheckReviews(loaded.Chunks[0], loaded.ImplementationSession); err != nil {
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
		if _, err := rr.Read(t.TempDir(), spec); err == nil {
			t.Fatalf("unsafe path %q accepted", spec)
		}
	}
}
