package benchguard

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// blockedPrefixWant is the refusal prefix every side keeps, written here independently
// of the production constant.
const blockedPrefixWant = "BLOCKED: Bench response is bounded, complete, and self-contained."

func TestClassifyFollowOns(t *testing.T) {
	resolver := Resolver{Getwd: func() (string, error) { return "/work", nil }, EvalSymlinks: func(path string) (string, error) {
		if path == "/work/kit-command" {
			return "/work/bin/bench.sh", nil
		}
		return "", errors.New("not bench")
	}}
	for _, command := range []string{"bench help | touch marker", "bench gate 2>&1", "X=1 bench help && touch marker", "bash -lc 'bench help || touch marker'", "bench worktree exec a -- bench gate > marker"} {
		if !Classify(command, resolver).Blocked {
			t.Errorf("Classify(%q) = false, want true", command)
		}
	}
	for _, command := range []string{"bench gate --fresh", "rg bench AGENTS.md", "bench worktree exec a -- bash -lc 'go test && go vet'", "cat <<'EOF'\nbench gate | tail\nEOF"} {
		if Classify(command, resolver).Blocked {
			t.Errorf("Classify(%q) = true, want false", command)
		}
	}
}

func TestClassifyResolvesBareAliasFromProcessPath(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "bench.sh")
	if err := os.WriteFile(target, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(dir, "kit-command")); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	if !Classify("kit-command help | touch marker", DefaultResolver()).Blocked {
		t.Fatal("bare PATH alias did not classify as Bench")
	}
}

func TestCommandFromEnvelope(t *testing.T) {
	if _, err := CommandFromEnvelope([]byte(`{"tool_input":{"command":"bench help\\u0026\\u0026 touch marker"}}`)); err != nil {
		t.Fatal(err)
	}
	for _, input := range []string{`{`, `{}`, `{"tool_input":{"command":1}}`, `{"tool_input":{"command":""}}`, `{"tool_input":{"command":"\u0000"}}`} {
		if _, err := CommandFromEnvelope([]byte(input)); err == nil {
			t.Errorf("CommandFromEnvelope(%q) succeeded", input)
		}
	}
}

// TestInvokesBenchWalksOneWrapperLevel proves the exported Bench test reads a
// wrapped string, so a `bash -c` call that runs Bench is a Bench call.
func TestInvokesBenchWalksOneWrapperLevel(t *testing.T) {
	resolver := Resolver{Getwd: func() (string, error) { return "/work", nil }, EvalSymlinks: func(string) (string, error) { return "", errors.New("not bench") }}
	for _, command := range []string{"bench status", "bash -c 'bench status'", "bench worktree exec a -- sed -i x /pool/id/b", "sh -c 'cd /tmp && bench gate'"} {
		if !InvokesBench(command, resolver) {
			t.Errorf("InvokesBench(%q) = false, want true", command)
		}
	}
	for _, command := range []string{"ls", "rg bench AGENTS.md", "bash -c 'ls /pool/id'"} {
		if InvokesBench(command, resolver) {
			t.Errorf("InvokesBench(%q) = true, want false", command)
		}
	}
}

// TestClassifySpanScopedFollowOns drives the span-scoped rule. A `bench worktree
// exec` head allows a heredoc in its span and a `;` or `&&` after it; every other
// shape refuses. A refuse row that carries a want also grades the refusal line.
// (Coverage rows G1, G2, G3, G4, G5, G7, G14, KG14, KG15.)
func TestClassifySpanScopedFollowOns(t *testing.T) {
	resolver := Resolver{Getwd: func() (string, error) { return "/work", nil }, EvalSymlinks: func(string) (string, error) { return "", errors.New("not bench") }}
	for _, tc := range []struct{ row, command string }{
		{"G1", "bench worktree exec L -- cp a b; cp b a"},
		{"G1", "bench worktree exec L -- cp a b && rg -n x b"},
		{"G5", "bench worktree exec L -- cat <<'EOF'\nbench gate\nEOF"},
		{"G7", "cat <<'EOF'\nfirst line\nsecond line\nEOF"},
		{"KG14", "bench worktree exec L -- cat <<'EOF'\nEOF"},
	} {
		if Classify(tc.command, resolver).Blocked {
			t.Errorf("%s: Classify(%q) = true, want false", tc.row, tc.command)
		}
	}
	for _, tc := range []struct{ row, command, want string }{
		{row: "G2", command: "cp a b && bench worktree exec L -- true"},
		{row: "G3", command: "bench worktree exec L -- true | cat"},
		{row: "G3", command: "bench worktree exec L -- true || echo x"},
		{row: "G3", command: "bench worktree exec L -- true &"},
		{row: "G3", command: "bench worktree exec L -- true > out"},
		{row: "G3", command: "bench worktree exec L -- true <<< x"},
		{row: "G4", command: "bench worktree exec L -- true; bench maps"},
		{row: "G14", command: "bench gate; cp a b"},
		{row: "KG15", command: "bench worktree exec L -- cat < /dev/null", want: "Run the Bench command without a redirection. segment=bench worktree exec L -- cat operator=</dev/null"},
	} {
		verdict := Classify(tc.command, resolver)
		if !verdict.Blocked {
			t.Errorf("%s: Classify(%q) = false, want true", tc.row, tc.command)
			continue
		}
		if tc.want != "" && !strings.HasSuffix(verdict.Message(), tc.want) {
			t.Errorf("%s: message %q does not end with %q", tc.row, verdict.Message(), tc.want)
		}
	}
}

// TestClassifyNamesTheSegmentAndTheOperator proves the refusal line keeps the fixed
// sentence and then names the Bench segment and the operator that caused it.
// (Coverage rows G8, G9, G10.)
func TestClassifyNamesTheSegmentAndTheOperator(t *testing.T) {
	resolver := Resolver{Getwd: func() (string, error) { return "/work", nil }, EvalSymlinks: func(string) (string, error) { return "", errors.New("not bench") }}
	for _, tc := range []struct{ row, command, want string }{
		{"G8", "cat a && echo x; bench maps", "segment=bench maps operator=;"},
		{"G9", "bench gate 2>&1", "segment=bench gate operator=2>&1"},
		{"G10", "cp a b && bench worktree exec L -- true | cat", "segment=bench worktree exec L -- true operator=&&"},
		{"G9b", "bench worktree exec L -- cat <<'EOF' 2>&1\nx\nEOF", "segment=bench worktree exec L -- cat operator=2>&1"},
	} {
		verdict := Classify(tc.command, resolver)
		if !verdict.Blocked {
			t.Errorf("%s: Classify(%q) allowed the call", tc.row, tc.command)
			continue
		}
		message := verdict.Message()
		if !strings.HasPrefix(message, blockedPrefixWant) {
			t.Errorf("%s: message %q does not start with the fixed prefix", tc.row, message)
		}
		if !strings.HasSuffix(message, tc.want) {
			t.Errorf("%s: message %q does not end with %q", tc.row, message, tc.want)
		}
	}
}

// TestClassifyNamesTheSideOfTheOperator proves the refusal names the side the operator
// sits on. A leading operator is not a follow-on, so the sentence that tells the reader
// to remove a follow-on points at nothing.
func TestClassifyNamesTheSideOfTheOperator(t *testing.T) {
	resolver := Resolver{Getwd: func() (string, error) { return "/work", nil }, EvalSymlinks: func(string) (string, error) { return "", errors.New("not bench") }}
	const prefix = blockedPrefixWant
	for _, tc := range []struct{ name, command, want string }{
		{
			name:    "leading operator",
			command: "cd /tmp && bench gate",
			want:    prefix + " Run the Bench command from the current directory; it resolves the worktree itself. segment=bench gate operator=&&",
		},
		{
			name:    "trailing operator",
			command: "bench gate && echo done",
			want:    prefix + " Run the Bench command without a shell follow-on. segment=bench gate operator=&&",
		},
		{
			name:    "redirection inside the span",
			command: "bench gate 2>&1",
			want:    prefix + " Run the Bench command without a redirection. segment=bench gate operator=2>&1",
		},
	} {
		verdict := Classify(tc.command, resolver)
		if !verdict.Blocked {
			t.Errorf("%s: Classify(%q) allowed the call", tc.name, tc.command)
			continue
		}
		if got := verdict.Message(); got != tc.want {
			t.Errorf("%s: message = %q, want %q", tc.name, got, tc.want)
		}
		if !strings.HasPrefix(verdict.Message(), prefix) {
			t.Errorf("%s: message %q does not start with the fixed prefix", tc.name, verdict.Message())
		}
	}
}

// TestPoolReferenceRefusesAPoolPathOutsideTheExecVerb drives the pool denial. A worktree
// runs through `bench worktree exec`, so a `cd`, an assignment, or a git repository
// option that reaches the pool path is the mistake the guard names. A read of a file
// under the pool, a pool path the `bench worktree` grammar takes, and a relative target
// inside a wrapper string stay allowed, because the guard reads text only and never
// resolves a path. A shell loop stays allowed, whether it iterates the pool path or cds
// into it inside its body: the loop reads as one command whose first word is `for`, so
// neither the `cd` arm nor the assignment arm reaches it. `craft-delegate` therefore
// carries the loop as a guidance rule, and these two rows record how far Go reaches.
func TestPoolReferenceRefusesAPoolPathOutsideTheExecVerb(t *testing.T) {
	const pools = "/home/agent/.bench/worktrees"
	const assignment = pools + "/bench-123/abc-def"
	for _, tc := range []struct{ name, command, want string }{
		{"bare cd into the pool", "cd " + assignment, assignment},
		{"cd into the pool then a follow-on", "cd " + assignment + " && go test ./...", assignment},
		{"assignment then git -C the variable", "W=" + assignment + "; git -C $W log", assignment},
		{"git -C into the pool", "git -C " + assignment + " log", assignment},
		{"git --git-dir into the pool", "git --git-dir=" + assignment + "/.git log", assignment + "/.git"},
		{"git --work-tree into the pool", "git --work-tree " + assignment + " status", assignment},
		{"export of the pool path", "export W=" + assignment, assignment},
		{"env assignment of the pool path", "env W=" + assignment + " git log", assignment},
		{"relative cd inside a wrapper string", "bench worktree exec \"x\" -- sh -c 'cd sub && go test ./...'", ""},
		{"cd outside the pool", "cd /tmp", ""},
		{"cd through an unexpanded variable", "cd \"$W\"", ""},
		{"cat a file in the pool", "cat " + assignment + "/specs/x.md", ""},
		{"sed a file in the pool", "sed -n 1,5p " + assignment + "/f", ""},
		{"worktree release takes the pool path", "bench worktree release --request r " + assignment, ""},
		{"worktree land takes the pool path", "bench worktree land --request r --base a --source-tip b -m m " + assignment, ""},
		{"assignment and git -C outside the pool", "W=/tmp/x; git -C /tmp/x log", ""},
		{"shell loop over the pool path", "for f in " + assignment + "/*.md; do cat \"$f\"; done", ""},
		{"shell loop whose body cds into the pool", "for f in a b; do cd " + assignment + "; done", ""},
	} {
		if got := PoolReference(tc.command, pools); got != tc.want {
			t.Errorf("%s: PoolReference(%q) = %q, want %q", tc.name, tc.command, got, tc.want)
		}
	}
}

// TestPoolReferenceMessageNamesTheTarget proves the refusal states the one command form
// and the path it read.
func TestPoolReferenceMessageNamesTheTarget(t *testing.T) {
	const want = `BLOCKED: a Bench worktree runs through bench worktree exec "<label>" -- <command>; never cd, assign, or git -C into the pool path. target=/home/agent/.bench/worktrees/bench-123/abc-def`
	if got := PoolReferenceMessage("/home/agent/.bench/worktrees/bench-123/abc-def"); got != want {
		t.Errorf("PoolReferenceMessage = %q, want %q", got, want)
	}
}
