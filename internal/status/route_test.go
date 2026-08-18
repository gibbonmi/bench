package status

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsInvocable(t *testing.T) {
	cases := []struct {
		action string
		want   bool
	}{
		{"bench gate --fresh for a whole-tree verdict", false},
		{"/bench-final-check / push", false},
		{"benchmark the loop", false},
		{"git push", true},
		{"bench spec retire <slug>", true},
		{"/bench-drain", true},
	}
	for _, tc := range cases {
		if got := IsInvocable(tc.action); got != tc.want {
			t.Errorf("IsInvocable(%q) = %v, want %v", tc.action, got, tc.want)
		}
	}
}

func TestRouteTranslatesHarnessPrefix(t *testing.T) {
	signals := []Signal{{Severity: 4, Name: "drain", Detail: "1 idea", Action: "/bench-drain"}}
	for _, tc := range []struct {
		harness, want string
	}{
		{HarnessClaude, "/bench-drain"},
		{"codex", "$bench-drain"},
	} {
		if got := Route(signals, tc.harness).Lead.Action; got != tc.want {
			t.Errorf("Route(%q).Lead.Action = %q, want %q", tc.harness, got, tc.want)
		}
	}
}

func TestFirstInvocable(t *testing.T) {
	prose := []Signal{
		{Severity: 0, Name: "gate", Detail: "red", Action: "fix before commit"},
		{Severity: 1, Name: "git", Detail: "1 dirty path", Action: "commit on green"},
		{Severity: 5, Name: "structure", Detail: "2 issue(s)", Action: "split (craft-seams)"},
	}
	cases := []struct {
		what        string
		signals     []Signal
		action      string
		wantName    string
		wantPresent bool
	}{
		{what: "no board at all", signals: nil},
		{what: "empty board", signals: []Signal{}},
		{what: "every action is prose", signals: prose},
		{
			what: "prose ahead of a command",
			signals: append(append([]Signal{}, prose...),
				Signal{Severity: 4, Name: "drain", Detail: "1 idea(s)", Action: "/bench-what-next"}),
			action: "/bench-what-next", wantName: "drain", wantPresent: true,
		},
		{
			what: "a command leads",
			signals: []Signal{
				{Severity: 3, Name: "guards", Detail: "pre-push hook missing", Action: "bench link"},
				{Severity: 4, Name: "drain", Detail: "1 idea(s)", Action: "/bench-what-next"},
			},
			action: "bench link", wantName: "guards", wantPresent: true,
		},
		{
			what: "a bench subcommand carrying an argument placeholder",
			signals: []Signal{
				{Severity: 5, Name: "structure", Detail: "2 issue(s)", Action: "split (craft-seams)"},
				{Severity: 8, Name: "specs", Detail: "1 merged spec(s)", Action: "bench spec retire <slug>"},
			},
			action: "bench spec retire <slug>", wantName: "specs", wantPresent: true,
		},
		{
			what: "a compound action naming several steps is not one invocation",
			signals: []Signal{
				{Severity: 1, Name: "git", Detail: "7 unpushed commits", Action: "/bench-final-check / push"},
				{Severity: 4, Name: "drain", Detail: "1 idea(s)", Action: "/bench-what-next"},
			},
			action: "/bench-what-next", wantName: "drain", wantPresent: true,
		},
		{
			what:    "a word beginning with bench is not a command",
			signals: []Signal{{Severity: 11, Name: "handoff", Detail: "stale", Action: "benchmark the loop"}},
		},
	}
	for _, tc := range cases {
		_, signal, ok := firstInvocable(tc.signals)
		if ok != tc.wantPresent || signal.Action != tc.action || signal.Name != tc.wantName {
			t.Fatalf("%s: firstInvocable = (%q, %q, %v), want (%q, %q, %v)",
				tc.what, signal.Action, signal.Name, ok, tc.action, tc.wantName, tc.wantPresent)
		}
	}
}

func TestRouteForUsesCleanFallbackOnlyForAnEmptyBoard(t *testing.T) {
	root := t.TempDir()
	drain := RouteFor(root, nil, HarnessClaude)
	if drain.Lead != (Signal{Name: "clean", Detail: "nothing pending", Action: "/bench-drain"}) || drain.NoCommand {
		t.Fatalf("drain fallback = %#v", drain)
	}
	if got, code := renderRoute(drain); code != 0 || got != "next[1]{state,why,command}:\n  clean,nothing pending,/bench-drain\nalso: none\n" {
		t.Fatalf("drain fallback output = (%q, %d)", got, code)
	}
	if err := os.WriteFile(filepath.Join(root, "ROADMAP.md"), []byte("# Roadmap\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	roadmap := RouteFor(root, nil, HarnessClaude)
	if roadmap.Lead != (Signal{Name: "clean", Detail: "nothing pending", Action: "bench roadmap"}) || roadmap.NoCommand {
		t.Fatalf("roadmap fallback = %#v", roadmap)
	}
	if got, code := renderRoute(roadmap); code != 0 || got != "next[1]{state,why,command}:\n  clean,nothing pending,bench roadmap\nalso: none\n" {
		t.Fatalf("roadmap fallback output = (%q, %d)", got, code)
	}

	locked := RouteFor(root, []Signal{{Name: "gate", Detail: "locked-pending"}}, HarnessClaude)
	if locked.Lead != (Signal{Name: "gate", Detail: "locked-pending"}) || !locked.NoCommand || len(locked.RunnersUp) != 0 {
		t.Fatalf("locked board = %#v, want no-command lead without fallback", locked)
	}
	if got, code := renderRoute(locked); code != 0 || got != "next[1]{state,why,command}:\n  gate,locked-pending,\"\"\nalso: none\n" {
		t.Fatalf("locked route = (%q, %d)", got, code)
	}
}

func TestRouteForTranslatesLeadAndRunnersUpForCodex(t *testing.T) {
	route := RouteFor(t.TempDir(), []Signal{
		{Name: "drain", Detail: "1 idea", Action: "/bench-what-next"},
		{Name: "specs", Detail: "1 staged spec", Action: "/bench-implement-spec"},
	}, "codex")
	if got, code := renderRoute(route); code != 0 || got != "next[1]{state,why,command}:\n  drain,1 idea,$bench-what-next\nalso: specs (1 staged spec) → $bench-implement-spec\n" {
		t.Fatalf("codex route = (%q, %d)", got, code)
	}
}
