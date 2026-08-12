package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnchorsAppendsHonestEmptyHelpWithoutChangingPrimaryResponse(t *testing.T) {
	tests := []struct {
		name, fixture string
		args          []string
		wantPrefix    string
		wantCode      int
	}{
		{name: "populated", fixture: "pre-disclosure-populated.stdout", args: []string{"AGENTS.md"}, wantPrefix: "anchors[", wantCode: 0},
		{name: "empty", fixture: "pre-disclosure-empty.stdout", args: []string{"specs/axi-query-disclosure/tickets/render-honest-anchor-help.md"}, wantPrefix: "anchors[0]{kind,section,needle}:\n", wantCode: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			primary, err := os.ReadFile(filepath.Join("testdata", "anchors", tt.fixture))
			if err != nil {
				t.Fatal(err)
			}
			got, code := anchorsCommand(tt.args)
			if code != tt.wantCode {
				t.Fatalf("anchorsCommand(%v) exit = %d, want %d", tt.args, code, tt.wantCode)
			}
			if tt.name == "populated" {
				want := string(primary) + "help[0]{cmd,why}:\n"
				if got != want {
					t.Fatalf("anchorsCommand(%v) = %q, want captured response %q", tt.args, got, want)
				}
			} else if got != string(primary)+"help[0]{cmd,why}:\n" {
				t.Fatalf("anchorsCommand(%v) = %q, want exact empty response", tt.args, got)
			}
		})
	}

	got, code := anchorsCommand(nil)
	if code != 2 || got != "usage: bench anchors (missing argument: argument)\n" {
		t.Fatalf("anchorsCommand(nil) = (%q, %d), want unchanged usage refusal", got, code)
	}
}
