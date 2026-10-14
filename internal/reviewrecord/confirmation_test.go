package reviewrecord_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// TestConfirmationParse grades the parser on one chunk confirmation: it accepts a valid
// confirmation and refuses a wrong role, a bad excerpt digest, and an ID that a review holds.
func TestConfirmationParse(t *testing.T) {
	for _, tc := range []struct {
		name, reason string
		change       func(*rr.Chunk)
	}{
		{"valid", "", func(*rr.Chunk) {}},
		{"wrong role", "invalid confirmation role", func(c *rr.Chunk) { c.Confirmations[0].Role = "independent-review" }},
		{"bad excerpt digest", "excerpt digest", func(c *rr.Chunk) { c.Confirmations[0].NativeRef.Digest = rr.Digest([]byte("another excerpt")) }},
		{"review ID", "invalid duplicate or empty evidence id", func(c *rr.Chunk) { c.Confirmations[0].ID = c.Reviews[0].ID }},
		{"invalid pair", "invalid confirmation pair", func(c *rr.Chunk) { c.Confirmations[0].ReviewedTip = "main" }},
		{"unnamed predecessor", "invalid supersession", func(c *rr.Chunk) { c.Confirmations[0].Supersedes = []string{"confirm-0"} }},
		{"missing probe mutation", "missing probe mutation", func(c *rr.Chunk) { c.Confirmations[0].Probe.Mutation = "" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := recordtest.New(t, 1)
			f.AddChunk()
			f.Confirm("confirm-1", f.Record.Chunks[0].Base, "F-1")
			tc.change(&f.Record.Chunks[0])
			data, err := json.Marshal(f.Record)
			if err != nil {
				t.Fatal(err)
			}
			parsed, err := rr.Parse(data)
			if tc.reason == "" {
				if err != nil || !reflect.DeepEqual(parsed.Chunks[0].Confirmations, f.Record.Chunks[0].Confirmations) {
					t.Fatalf("Parse = %v, %+v; want the confirmation as written", err, parsed.Chunks)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.reason) {
				t.Fatalf("Parse error = %v, want %q", err, tc.reason)
			}
		})
	}
}

// TestConfirmationAbsentKey grades a record without confirmations: it renders no
// confirmations key, so every existing record keeps its bytes.
func TestConfirmationAbsentKey(t *testing.T) {
	f := recordtest.New(t, 1)
	f.AddChunk()
	data, err := rr.Render(nil, f.Record)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), `"confirmations"`) {
		t.Fatalf("record without confirmations renders the key:\n%s", data)
	}
}
