package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/rgguard"
)

// TestGuardRgVerb covers R44 and the verb's refusal. An unreadable envelope warns and
// allows, a refused call exits 2 with the classifier's message, and a call with a path
// passes in silence. The refusal text comes from its one owner.
func TestGuardRgVerb(t *testing.T) {
	for _, tc := range []struct {
		name, envelope string
		code           int
		stderr         string
	}{
		{"R44 unreadable envelope", `{`, 0, "WARNING: block-unsafe-rg: unreadable command field — allowing Bash.\n"},
		{"pathless refusal", bashEnvelope("rg foo"), 2, rgguard.Classify("rg foo").Message() + "\n"},
		{"call with a path", bashEnvelope("rg foo ."), 0, ""},
	} {
		var errb bytes.Buffer
		if code := guardRg(nil, strings.NewReader(tc.envelope), io.Discard, &errb); code != tc.code || errb.String() != tc.stderr {
			t.Errorf("%s = (%d, %q), want (%d, %q)", tc.name, code, errb.String(), tc.code, tc.stderr)
		}
	}
}
