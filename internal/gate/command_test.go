package gate

import (
	"bytes"
	"strings"
	"testing"
)

func TestCommandHandlesPublicGateUsageWithoutStartingTheOracle(t *testing.T) {
	for _, test := range []struct {
		name       string
		args       []string
		wantExit   int
		wantStdout string
		wantStderr string
	}{
		{name: "help", args: []string{"--help"}, wantExit: 0, wantStdout: commandUsage},
		{name: "unknown", args: []string{"unexpected"}, wantExit: 2, wantStderr: commandUsage},
		{name: "fresh extra argument", args: []string{"--fresh", "unexpected"}, wantExit: 2, wantStderr: commandUsage},
		{name: "pin extra argument", args: []string{"pin", "unexpected"}, wantExit: 2, wantStderr: "usage: bench gate pin"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if got := Command(test.args, strings.NewReader(""), &stdout, &stderr); got != test.wantExit {
				t.Fatalf("exit = %d, want %d", got, test.wantExit)
			}
			if !strings.Contains(stdout.String(), test.wantStdout) {
				t.Fatalf("stdout = %q, want %q", stdout.String(), test.wantStdout)
			}
			if !strings.Contains(stderr.String(), test.wantStderr) {
				t.Fatalf("stderr = %q, want %q", stderr.String(), test.wantStderr)
			}
		})
	}
}
