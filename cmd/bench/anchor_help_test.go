package main

import "testing"

func TestAnchorsAppendsHonestEmptyHelpWithoutChangingPrimaryResponse(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantPrefix string
		wantCode   int
	}{
		{name: "populated", args: []string{"AGENTS.md"}, wantPrefix: "anchors[", wantCode: 0},
		{name: "empty", args: []string{"specs/axi-query-disclosure/tickets/render-honest-anchor-help.md"}, wantPrefix: "anchors[0]{kind,section,needle}:\n", wantCode: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, code := anchorsCommand(tt.args)
			if code != tt.wantCode {
				t.Fatalf("anchorsCommand(%v) exit = %d, want %d", tt.args, code, tt.wantCode)
			}
			if tt.name == "populated" {
				const want = "anchors[4]{kind,section,needle}:\n" +
					"  require,\"\",never by polling a self-matching pattern\n" +
					"  require,\"\",\"runs plan-before-apply: print the exact target list, sample it, then apply\"\n" +
					"  require,\"\",repository-wide sweep uses `rg --hidden`\n" +
					"  require,\"\",Discover Bench verbs non-interactively\n" +
					"help[0]{cmd,why}:\n"
				if got != want {
					t.Fatalf("anchorsCommand(%v) = %q, want captured response %q", tt.args, got, want)
				}
			} else if got != tt.wantPrefix+"help[0]{cmd,why}:\n" {
				t.Fatalf("anchorsCommand(%v) = %q, want exact empty response", tt.args, got)
			}
		})
	}

	got, code := anchorsCommand(nil)
	if code != 2 || got != "usage: bench anchors (missing argument: argument)\n" {
		t.Fatalf("anchorsCommand(nil) = (%q, %d), want unchanged usage refusal", got, code)
	}
}
