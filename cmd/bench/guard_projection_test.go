package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchguard"
	"github.com/gibbonmi/bench/internal/census"
	"github.com/gibbonmi/bench/internal/commitment/commitcmd"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/poolkey"
)

// guardEnvelopeCase is one raw stdin envelope and the exit a guard entry keeps for it.
// warn marks an exit-0 case whose stderr carries the entry's unreadable-input warning.
// A nonempty message is the one refusal or advice line that stderr must hold.
type guardEnvelopeCase struct {
	name    string
	in      string
	want    int
	warn    bool
	message string
}

// checkGuardEnvelope runs each case through the actual guard entry and grades its exit and streams.
// A case with a message carries exactly that line, exit 2 otherwise carries BLOCKED, a warned exit 0 carries WARNING, and any other exit 0 writes nothing.
func checkGuardEnvelope(t *testing.T, entry func([]string, io.Reader, io.Writer, io.Writer) int, cases []guardEnvelopeCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := entry(nil, strings.NewReader(tc.in), &stdout, &stderr)
			if code != tc.want {
				t.Errorf("entry(%q) exit = %d, want %d", tc.in, code, tc.want)
			}
			if stdout.Len() != 0 {
				t.Errorf("entry(%q) stdout = %q, want empty", tc.in, stdout.String())
			}
			got := stderr.String()
			switch {
			case tc.message != "":
				if got != tc.message+"\n" {
					t.Errorf("entry(%q) stderr = %q, want %q", tc.in, got, tc.message+"\n")
				}
			case tc.want == 2 && !strings.HasPrefix(got, "BLOCKED:"):
				t.Errorf("entry(%q) stderr = %q, want BLOCKED", tc.in, got)
			case tc.want == 0 && tc.warn && !strings.HasPrefix(got, "WARNING: "):
				t.Errorf("entry(%q) stderr = %q, want WARNING", tc.in, got)
			case tc.want == 0 && !tc.warn && got != "":
				t.Errorf("entry(%q) stderr = %q, want empty", tc.in, got)
			}
		})
	}
}

// TestGuardGitEnvelopePolicy grades guard-git exits for envelope family E0 and the Git decoder rows.
// Each git reset case blocks only when the standard struct decode keeps the command.
func TestGuardGitEnvelopePolicy(t *testing.T) {
	checkGuardEnvelope(t, guardGit, []guardEnvelopeCase{
		{name: "empty stdin", in: ``},
		{name: "invalid JSON", in: `not json`},
		{name: "absent tool_input", in: `{}`},
		{name: "null tool_input", in: `{"tool_input":null}`},
		{name: "absent command", in: `{"tool_input":{}}`},
		{name: "null command", in: `{"tool_input":{"command":null}}`},
		{name: "numeric command", in: `{"tool_input":{"command":7}}`},
		{name: "array command", in: `{"tool_input":{"command":["git reset --hard"]}}`},
		{name: "empty command", in: `{"tool_input":{"command":""}}`},
		{name: "NUL command", in: `{"tool_input":{"command":"\u0000"}}`},
		{name: "nested key casing", in: `{"tool_input":{"COMMAND":"git reset --hard"}}`, want: 2},
		{name: "duplicate command null", in: `{"tool_input":{"command":"git reset --hard","command":null}}`, want: 2},
		{name: "duplicate input null", in: `{"tool_input":{"command":"git reset --hard","file_path":"README.md"},"tool_input":null,"cwd":"/repo"}`, want: 2},
		{name: "earlier command type error", in: `{"tool_input":{"command":7,"command":"git reset --hard"}}`},
		{name: "root key casing", in: `{"TOOL_INPUT":{"command":"git reset --hard"}}`, want: 2},
		{name: "unrelated cwd type", in: `{"cwd":7,"tool_input":{"command":"git reset --hard"}}`, want: 2},
	})
}

// TestGuardBenchEnvelopePolicy grades the Bench follow-on entry for envelope family E0.
// Every unreadable command keeps the warning allow, not Git's silent allow or a refusal.
func TestGuardBenchEnvelopePolicy(t *testing.T) {
	checkGuardEnvelope(t, guardBenchFollowOn, []guardEnvelopeCase{
		{name: "empty stdin", in: ``, warn: true},
		{name: "invalid JSON", in: `not json`, warn: true},
		{name: "absent tool_input", in: `{}`, warn: true},
		{name: "null tool_input", in: `{"tool_input":null}`, warn: true},
		{name: "absent command", in: `{"tool_input":{}}`, warn: true},
		{name: "null command", in: `{"tool_input":{"command":null}}`, warn: true},
		{name: "numeric command", in: `{"tool_input":{"command":7}}`, warn: true},
		{name: "array command", in: `{"tool_input":{"command":["bench help"]}}`, warn: true},
		{name: "empty command", in: `{"tool_input":{"command":""}}`, warn: true},
		{name: "NUL command", in: `{"tool_input":{"command":"\u0000"}}`, warn: true},
	})
}

// TestGuardWriteEnvelopePolicy grades the file-write entry against a real primary checkout with a tracked README.md.
// Unreadable paths keep the warning allow, an unreadable fact allows silently, and kept cwd facts block.
func TestGuardWriteEnvelopePolicy(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("tracked\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitInitRepo(t, root)
	quote := func(s string) string { b, _ := json.Marshal(s); return string(b) }
	// The matrix spells the fixture root "/repo"; each case substitutes the real checkout.
	at := func(in string) string { return strings.ReplaceAll(in, `"/repo"`, quote(root)) }
	checkGuardEnvelope(t, guardFileWrite, []guardEnvelopeCase{
		{name: "empty stdin", in: ``, warn: true},
		{name: "invalid JSON", in: `not json`, warn: true},
		{name: "absent tool_input", in: at(`{"cwd":"/repo"}`), warn: true},
		{name: "absent file_path", in: at(`{"cwd":"/repo","tool_input":{}}`), warn: true},
		{name: "numeric file_path", in: at(`{"cwd":"/repo","tool_input":{"file_path":7}}`), warn: true},
		{name: "empty file_path", in: at(`{"cwd":"/repo","tool_input":{"file_path":""}}`), warn: true},
		{name: "NUL file_path", in: at(`{"cwd":"/repo","tool_input":{"file_path":"README\u0000.md"}}`), warn: true},
		{name: "relative file_path with no cwd", in: `{"tool_input":{"file_path":"README.md"}}`, warn: true},
		{name: "earlier cwd type error", in: at(`{"cwd":7,"cwd":"/repo","tool_input":{"file_path":"README.md"}}`), warn: true},
		{name: "unreadable root fact", in: `{"tool_input":{"file_path":` + quote(filepath.Join(t.TempDir(), "notes.md")) + `}}`},
		{name: "duplicate cwd null", in: at(`{"cwd":"/repo","cwd":null,"tool_input":{"file_path":"README.md"}}`), want: 2},
		{name: "write root key casing", in: at(`{"CWD":"/repo","TOOL_INPUT":{"file_path":"README.md"}}`), want: 2},
	})
}

// poolTarget sets a fresh Bench home and returns a path under its worktree pool.
func poolTarget(t *testing.T) string {
	t.Helper()
	t.Chdir(t.TempDir())
	home := filepath.Join(t.TempDir(), "home")
	t.Setenv("BENCH_HOME", home)
	return filepath.Join(poolkey.Pools(home), "repository", "assignment")
}

// TestGuardProjectionPoolPrecedence is GP13: a pool reference with a follow-on emits the pool refusal.
func TestGuardProjectionPoolPrecedence(t *testing.T) {
	target := poolTarget(t)
	checkGuardEnvelope(t, guardBenchFollowOn, []guardEnvelopeCase{
		{name: "pool then follow-on", in: bashEnvelope("cd " + target + " && bench status | cat"), want: 2, message: benchguard.PoolReferenceMessage(target)},
	})
}

// TestGuardProjectionReviewerPrecedence is GP65 and GP66. The pool refusal wins over the
// reviewer verb, and the reviewer-verb refusal wins over a shell follow-on.
func TestGuardProjectionReviewerPrecedence(t *testing.T) {
	target := poolTarget(t)
	checkGuardEnvelope(t, guardBenchFollowOn, []guardEnvelopeCase{
		{name: "pool then reviewer verb", in: bashEnvelope("cd " + target + " && bench commitment override --assignment x --reason y"), want: 2, message: benchguard.PoolReferenceMessage(target)},
		{name: "reviewer verb then follow-on", in: bashEnvelope("bench commitment override --assignment x --reason y | cat"), want: 2, message: benchguard.ReviewerVerbMessage(commitcmd.ReviewerVerb)},
	})
}

// TestGuardProjectionCensusExclusion is GP44. A supported Bench invocation that names a pool
// path allows silently and writes no raw census record. The raw form of the same read writes one.
func TestGuardProjectionCensusExclusion(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	t.Chdir(root)
	home := filepath.Join(t.TempDir(), "home")
	t.Setenv("BENCH_HOME", home)
	raw, id := followOnEnvelope(home, root)
	record := filepath.Join(census.Dir(home, root), id)
	// Each form wraps the raw read, which starts `"sed ` and ends `"}}` in the envelope.
	wrap := func(head, tail string) string {
		return strings.Replace(strings.Replace(raw, `"sed `, `"`+head+`sed `, 1), `"}}`, tail+`"}}`, 1)
	}
	for _, tc := range []guardEnvelopeCase{
		{name: "exec child", in: wrap("bench worktree exec t -- ", "")},
		{name: "routine prefix", in: wrap("timeout 30 bench worktree exec t -- ", "")},
		{name: "wrapper child", in: wrap("bash -c 'bench worktree exec t -- ", "'")},
	} {
		checkGuardEnvelope(t, guardBenchFollowOn, []guardEnvelopeCase{tc})
		if _, err := os.Lstat(record); !os.IsNotExist(err) {
			t.Fatalf("%s wrote a raw census record (stat error %v)", tc.name, err)
		}
	}
	checkGuardEnvelope(t, guardBenchFollowOn, []guardEnvelopeCase{{name: "raw read", in: raw}})
	if _, err := os.Lstat(record); err != nil {
		t.Fatalf("raw read wrote no census record: %v", err)
	}
}
