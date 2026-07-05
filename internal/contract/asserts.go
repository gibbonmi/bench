package contract

import (
	"regexp"
	"strings"
	"testing"
)

func RequireFileMatches(t testing.TB, f Fixture, path, pattern, msg string) {
	t.Helper()
	if !regexp.MustCompile(pattern).MatchString(f.ReadFile(path)) {
		t.Fatalf("%s:\n%s", msg, f.ReadFile(path))
	}
}

func RequireContains(t testing.TB, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Fatalf("missing %q in:\n%s", needle, haystack)
	}
}

func RequireNotContains(t testing.TB, haystack, needle string) {
	t.Helper()
	if strings.Contains(haystack, needle) {
		t.Fatalf("unexpected %q in:\n%s", needle, haystack)
	}
}

func (p Probe) RequireExit(want int) {
	if p.ExitCode != want {
		p.t.Helper()
		p.t.Fatalf("exit code = %d, want %d\nstdout:\n%s\nstderr:\n%s", p.ExitCode, want, p.Stdout, p.Stderr)
	}
}

func (p Probe) RequireContains(haystack, needle string) {
	p.t.Helper()
	if !strings.Contains(haystack, needle) {
		p.t.Fatalf("missing %q\nstdout:\n%s\nstderr:\n%s", needle, p.Stdout, p.Stderr)
	}
}

func (p Probe) RequireNotContains(haystack, needle string) {
	p.t.Helper()
	if strings.Contains(haystack, needle) {
		p.t.Fatalf("unexpected %q\nstdout:\n%s\nstderr:\n%s", needle, p.Stdout, p.Stderr)
	}
}

func RequireIntEqual(t testing.TB, got, want int, msg string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %d, want %d", msg, got, want)
	}
}

func LineCount(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n")
}

func NonEmptyLines(s string) []string {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			out = append(out, line)
		}
	}
	return out
}
