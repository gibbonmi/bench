package assessment

import (
	"encoding/json"
	"fmt"
	"github.com/gibbonmi/bench/internal/bounds"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestAssessmentCommand(t *testing.T) {
	s := Store{Home: t.TempDir(), Root: t.TempDir()}
	out, code := Command(s, []string{"list"})
	if code != 0 || !strings.Contains(out, "runs[0]{run_id,state,attempts,task}:") || !strings.HasSuffix(out, "help[0]{cmd,why}:\n") {
		t.Fatalf("A15 empty list: %d %s", code, out)
	}
}

func TestAssessmentCommandDetail(t *testing.T) {
	s := Store{Home: t.TempDir(), Root: t.TempDir()}
	r := fixtureRun(s.Root)
	r.Attempts = nil
	for i, role := range Roles() {
		a := fixtureRun(s.Root).Attempts[0]
		a.AttemptID = fmt.Sprint("a", i)
		a.Role = role
		a.ChunkID = "chunk-one"
		r.Attempts = append(r.Attempts, a)
	}
	input := filepath.Join(t.TempDir(), "input.json")
	data, _ := json.Marshal(r)
	os.WriteFile(input, data, 0600)
	if out, code := Command(s, []string{"record", "--input", input}); code != 0 {
		t.Fatalf("record: %d %s", code, out)
	}
	out, code := Command(s, []string{"show", r.RunID})
	if code != 0 {
		t.Fatal(out)
	}
	for _, a := range r.Attempts {
		if !strings.Contains(out, strings.Join([]string{a.AttemptID, a.ChunkID, a.Role, a.State}, ",")) {
			t.Fatalf("A3 missing role %s", a.Role)
		}
	}
	if !strings.Contains(out, "/attempts/0/session_id") || !strings.Contains(out, "summary[1]") {
		t.Fatal("A15 detail omits native metadata or summary")
	}
	for _, verb := range []string{"list", "show", "record"} {
		for _, h := range []string{"help", "--help", "-h"} {
			if out, code := Command(s, []string{verb, h}); code != 0 || !strings.HasPrefix(out, "usage:") {
				t.Fatalf("help: %d %s", code, out)
			}
		}
	}
	if _, code := Command(s, nil); code != 2 {
		t.Fatal("bare command must return usage")
	}
	if _, code := Command(s, []string{"show", "missing"}); code != 1 {
		t.Fatal("missing run must refuse")
	}
}
func TestAssessmentCommandUnsafe(t *testing.T) {
	for _, kind := range []string{"traversal", "fifo", "symlink", "parent-link", "oversized", "version", "duplicate", "unknown-field", "duplicate-key", "task-esc", "task-bel", "chunk-esc", "chunk-bel"} {
		t.Run(kind, func(t *testing.T) {
			s := Store{Home: t.TempDir(), Root: t.TempDir()}
			r := fixtureRun(s.Root)
			input := filepath.Join(t.TempDir(), "input.json")
			switch kind {
			case "task-esc":
				r.TaskID = "task\x1b"
			case "task-bel":
				r.TaskID = "task\a"
			case "chunk-esc":
				r.Attempts[0].ChunkID = "chunk\x1b"
			case "chunk-bel":
				r.Attempts[0].ChunkID = "chunk\a"
			case "traversal":
				r.RunID = "../escape"
			case "version":
				r.Version = 9
			case "duplicate":
				r.Attempts = append(r.Attempts, r.Attempts[0])
			}
			data, _ := json.Marshal(r)
			if kind == "oversized" {
				data = []byte(strings.Repeat(" ", int(bounds.ControlRecordLimit)+1))
			}
			if kind == "duplicate-key" {
				data = []byte(strings.Replace(string(data), "{", "{\"version\":1,", 1))
			}
			if kind == "unknown-field" {
				data = []byte(strings.Replace(string(data), "{", "{\"execute\":\"false\",", 1))
			}
			if err := os.WriteFile(input, data, 0600); err != nil {
				t.Fatal(err)
			}
			switch kind {
			case "fifo":
				os.Remove(input)
				if err := syscall.Mkfifo(input, 0600); err != nil {
					t.Fatal(err)
				}
			case "symlink":
				link := input + ".link"
				if err := os.Symlink(input, link); err != nil {
					t.Fatal(err)
				}
				input = link
			case "parent-link":
				link := filepath.Join(t.TempDir(), "link")
				if err := os.Symlink(filepath.Dir(input), link); err != nil {
					t.Fatal(err)
				}
				input = filepath.Join(link, "input.json")
			}
			if out, code := Command(s, []string{"record", "--input", input}); code != 1 {
				t.Fatalf("A16 unsafe import accepted: %d %s", code, out)
			}
			out, code := Command(s, []string{"list"})
			if code != 0 || !strings.Contains(out, "runs[0]") {
				t.Fatalf("unsafe input left a run: %d %s", code, out)
			}
		})
	}
}
