package models

import (
	"reflect"
	"testing"
)

func TestParseIDs(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    []string
		wantErr bool
	}{
		{
			name: "two ids",
			body: `{"data":[{"id":"claude-a"},{"id":"claude-b"}]}`,
			want: []string{"claude-a", "claude-b"},
		},
		{
			name: "empty data",
			body: `{"data":[]}`,
			want: []string{},
		},
		{
			name:    "malformed body",
			body:    `{"data":[{"id":`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIDs([]byte(tt.body))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseIDs(%q) = %v, want error", tt.body, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseIDs(%q) unexpected error: %v", tt.body, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("parseIDs(%q) = %v, want %v", tt.body, got, tt.want)
			}
		})
	}
}

func TestNoKeyText(t *testing.T) {
	want := "No ANTHROPIC_API_KEY set, so I can't query the model list directly. Discover from\n" +
		"your harness instead, then bind the tiers (cheap / mid / top) in projects/<name>.md:\n" +
		"  - Claude Code: `claude --help`, or the in-app /model picker\n" +
		"  - Codex:       `codex --help`, or its model config\n" +
		"  - or export ANTHROPIC_API_KEY and re-run `bench models`\n"
	if got := noKeyText(); got != want {
		t.Fatalf("noKeyText() = %q, want %q", got, want)
	}
}
