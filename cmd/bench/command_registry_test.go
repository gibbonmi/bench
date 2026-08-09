package main

import (
	"bytes"
	"reflect"
	"runtime"
	"sort"
	"testing"
)

func TestCommandRunsVersionInProcess(t *testing.T) {
	var stdout, stderr, observation bytes.Buffer
	command := Command{Stdout: &stdout, Stderr: &stderr, Executable: "/selected/bench", Observe: &observation}
	if code := command.Run([]string{"version"}); code != 0 {
		t.Fatalf("version exit = %d, want 0; stderr=%q", code, stderr.String())
	}
	want := versionLine(version, runtime.GOOS, runtime.GOARCH) + "\n"
	if stdout.String() != want {
		t.Fatalf("version stdout = %q, want %q", stdout.String(), want)
	}
	if want := "command-registry:version\n"; observation.String() != want {
		t.Fatalf("version implementation observation = %q, want %q", observation.String(), want)
	}
}

func TestCommandDispositionsAreComplete(t *testing.T) {
	want := map[processAttachment][]string{
		attachmentDirect: {"check-agent-line", "commit", "gate-go", "guard-git", "resume-clean", "session-inspect", "shift", "spec", "version", "worktree"},
		attachmentSystem: {"canary", "doctor", "freshness-check", "freshness-publish", "gate", "gate-phases", "gate-pin", "gate-run", "init", "link", "setup", "stop-verdict", "unlink", "upgrade", "worktree-hook"},
		attachmentShip:   {"prep-release", "release", "release-preflight"},
	}
	got := map[processAttachment][]string{}
	seen := map[string]bool{}
	for _, disposition := range commandDispositions() {
		if seen[disposition.Name] {
			t.Fatalf("command disposition repeats %q", disposition.Name)
		}
		seen[disposition.Name] = true
		got[disposition.Attachment] = append(got[disposition.Attachment], disposition.Name)
	}
	for attachment := range got {
		sort.Strings(got[attachment])
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("command dispositions = %#v, want %#v", got, want)
	}
}
