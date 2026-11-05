package responsebound

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/bounds"
)

// TestReadPageArguments is S42: the next-page action and the first-page command of a stale
// cursor quote every reader argument, a hostile path included. Each command is authored by
// hand from the path and from the page's own source and next cells, a POSIX shell reads it
// back, and a sentinel file reds a path that the shell runs.
func TestReadPageArguments(t *testing.T) {
	sentinel := axitest.NewShellSentinel(t)
	path := "it's; " + sentinel.Touch()
	quoted := `'it'\''s; ` + sentinel.Touch() + `'`
	payload := []byte(strings.Repeat("x", 3*bounds.ResponseBytes))
	request := ReadRequest{Command: []string{"read", path}, Kind: "file", Scope: "/repo", Path: path, Digest: sha256.Sum256(payload), SourceSize: int64(len(payload)), Payload: bytes.NewReader(payload), Size: int64(len(payload)), Frame: readFrame}

	t.Run("next page", func(t *testing.T) {
		page, err := PageRead(request)
		if err != nil || page.Complete {
			t.Fatalf("PageRead = (%+v, %v), want a nonterminal page", page, err)
		}
		document, err := axitest.DecodeDocument(page.Output)
		if err != nil {
			t.Fatal(err)
		}
		rows, err := document.Rows("page")
		if err != nil {
			t.Fatal(err)
		}
		row := rows[0].(map[string]any)
		next, source := cell(row, "next"), cell(row, "source")
		actions, err := document.HelpActions()
		want := []axitest.HelpAction{{Cmd: "bench 'read' " + quoted + " '--offset' '" + next + "' '--source' '" + source + "'", Why: "read the next page"}}
		if err != nil || !slices.Equal(actions, want) {
			t.Fatalf("help actions = (%q, %v), want %q", actions, err, want)
		}
		argv, err := actions[0].Argv()
		if wantArgv := []string{"bench", "read", path, "--offset", next, "--source", source}; err != nil || !slices.Equal(argv, wantArgv) {
			t.Fatalf("argv of %q = (%q, %v), want %q", actions[0].Cmd, argv, err, wantArgv)
		}
	})
	t.Run("first page", func(t *testing.T) {
		stale := request
		stale.Offset, stale.Token, stale.Continued = 1, "stale", true
		_, err := PageRead(stale)
		var refusal *ReadRefusal
		const lead = "read again from the first page: "
		if want := lead + "bench 'read' " + quoted; !errors.As(err, &refusal) || refusal.Hint != want {
			t.Fatalf("PageRead refusal = %v, want hint %q", err, want)
		}
		argv, err := axitest.HelpAction{Cmd: strings.TrimPrefix(refusal.Hint, lead)}.Argv()
		if wantArgv := []string{"bench", "read", path}; err != nil || !slices.Equal(argv, wantArgv) {
			t.Fatalf("first-page argv = (%q, %v), want %q", argv, err, wantArgv)
		}
	})
	sentinel.RequireUnexecuted(t)
}
