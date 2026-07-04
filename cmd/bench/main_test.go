package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// The idiom-setting table test for the module: pure logic, table-driven, no process
// boundary. Acceptance for the version line lives at the shell seam (the gate's
// version-routing contract); this pins the format one layer down.
func TestVersionLine(t *testing.T) {
	cases := []struct {
		v, goos, goarch, want string
	}{
		{"0.2.0", "linux", "amd64", "benchkit 0.2.0 (linux/amd64)"},
		{"dev", "darwin", "arm64", "benchkit dev (darwin/arm64)"},
		{"1.0.0", "linux", "arm64", "benchkit 1.0.0 (linux/arm64)"},
	}
	for _, c := range cases {
		if got := versionLine(c.v, c.goos, c.goarch); got != c.want {
			t.Errorf("versionLine(%q,%q,%q) = %q, want %q", c.v, c.goos, c.goarch, got, c.want)
		}
	}
}

func TestRunVersionExits0(t *testing.T) {
	if rc := run([]string{"version"}, nil, nil); rc != 0 {
		t.Errorf("run version exit = %d, want 0", rc)
	}
}

func TestRunUnknownExits2(t *testing.T) {
	if rc := run([]string{"nope"}, nil, nil); rc != 2 {
		t.Errorf("run nope exit = %d, want 2", rc)
	}
}

func TestGuardGitDescribeClasses(t *testing.T) {
	var out bytes.Buffer
	if code := guardGit([]string{"--describe-classes"}, strings.NewReader(""), &out, io.Discard); code != 0 {
		t.Fatalf("--describe-classes exit = %d, want 0", code)
	}
	if !strings.HasPrefix(out.String(), "git push, ") {
		t.Errorf("--describe-classes did not print the deny surface to stdout: %q", out.String())
	}
}

func TestGuardGitBlockAllow(t *testing.T) {
	var errb bytes.Buffer
	block := `{"tool_input":{"command":"git push"}}`
	if code := guardGit(nil, strings.NewReader(block), io.Discard, &errb); code != 2 {
		t.Errorf("block exit = %d, want 2", code)
	}
	if !strings.Contains(errb.String(), "BLOCKED:") {
		t.Errorf("block did not emit BLOCKED on stderr: %q", errb.String())
	}
	for _, in := range []string{`{"tool_input":{"command":"git status"}}`, "not json", `{"tool_input":{"command":""}}`} {
		if code := guardGit(nil, strings.NewReader(in), io.Discard, io.Discard); code != 0 {
			t.Errorf("allow exit for %q = %d, want 0", in, code)
		}
	}
}

// panicReader forces guardGit's stdin read to panic, exercising the recover→exit-3
// rim so a crash can never masquerade as an exit-2 block.
type panicReader struct{}

func (panicReader) Read([]byte) (int, error) { panic("boom") }

func TestGuardGitRecoversToExit3(t *testing.T) {
	if code := guardGit(nil, panicReader{}, io.Discard, io.Discard); code != 3 {
		t.Errorf("panic mapped to exit %d, want 3", code)
	}
}
