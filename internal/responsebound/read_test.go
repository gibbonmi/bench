package responsebound

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/bounds"
)

// readFrame is a dispatcher tree frame of realistic size; every page budgets it.
var readFrame = "tree[1]{target,head,dirty}:\n  ft125-build," + strings.Repeat("a", 40) + ",false\n"

func pageRequest(frame string, payload []byte) ReadRequest {
	return ReadRequest{Command: []string{"read", "p.txt"}, Kind: "file", Scope: "/repo", Path: "p.txt", Digest: sha256.Sum256(payload), SourceSize: int64(len(payload)), Payload: bytes.NewReader(payload), Size: int64(len(payload)), Frame: frame}
}

// decodedPage is one page read back through the TOON decoder and the declared encoding.
type decodedPage struct {
	row  map[string]any
	data []byte
	argv []string
}

func decodePage(t *testing.T, output string) decodedPage {
	t.Helper()
	document, err := axitest.DecodeDocument(output)
	if err != nil || !reflect.DeepEqual(document.Blocks, []string{"page", "content", "help"}) {
		t.Fatalf("page = %q, want page, content, and help blocks: %v", output, err)
	}
	rows, err := document.Rows("page")
	if err != nil || len(rows) != 1 {
		t.Fatalf("page block = %v, want one row: %v", rows, err)
	}
	page := decodedPage{row: rows[0].(map[string]any)}
	contents, err := document.Rows("content")
	if err != nil {
		t.Fatal(err)
	}
	for _, content := range contents {
		data, err := DecodeContent(cell(page.row, "encoding"), content.(map[string]any)["data"].(string))
		if err != nil {
			t.Fatal(err)
		}
		page.data = append(page.data, data...)
	}
	if actions, err := document.HelpActions(); err != nil {
		t.Fatal(err)
	} else if len(actions) > 0 {
		if page.argv, err = axitest.RecoverHelpCommandArgv(output); err != nil {
			t.Fatal(err)
		}
	}
	return page
}

func cell(row map[string]any, field string) string { return fmt.Sprint(row[field]) }

// pageAll reads every page of payload from offset zero. It grades each page's bound, source,
// offset, encoding, and continuation, and the exact reconstruction, and it answers the
// rendered pages.
func pageAll(t *testing.T, payload []byte, encoding string) []string {
	t.Helper()
	request := pageRequest(readFrame, payload)
	var outputs []string
	var rebuilt []byte
	token := ""
	for len(outputs) < 1000 {
		page, err := PageRead(request)
		if err != nil {
			t.Fatalf("offset %d: %v", request.Offset, err)
		}
		if framed := readFrame + page.Output; !Fits(framed) {
			t.Fatalf("framed page has %d bytes and %d lines, past the response bound", len(framed), strings.Count(framed, "\n"))
		}
		outputs = append(outputs, page.Output)
		decoded := decodePage(t, page.Output)
		if token == "" {
			token = cell(decoded.row, "source")
		}
		if cell(decoded.row, "source") != token || cell(decoded.row, "offset") != strconv.FormatInt(request.Offset, 10) {
			t.Fatalf("page row = %v, want source %s at offset %d", decoded.row, token, request.Offset)
		}
		if len(payload) > 0 && cell(decoded.row, "encoding") != encoding {
			t.Fatalf("page encoding = %v, want %s", decoded.row["encoding"], encoding)
		}
		rebuilt = append(rebuilt, decoded.data...)
		if cell(decoded.row, "complete") == "true" {
			if decoded.argv != nil {
				t.Fatalf("terminal page offers a successor %q", decoded.argv)
			}
			break
		}
		next, _ := strconv.ParseInt(cell(decoded.row, "next"), 10, 64)
		want := []string{"bench", "read", "p.txt", OffsetFlag, strconv.FormatInt(next, 10), SourceFlag, token}
		if next <= request.Offset || !reflect.DeepEqual(decoded.argv, want) {
			t.Fatalf("continuation = %q at next %d, want %q past offset %d", decoded.argv, next, want, request.Offset)
		}
		request.Offset, request.Token, request.Continued = next, token, true
	}
	if !bytes.Equal(rebuilt, payload) {
		t.Fatalf("reconstruction differs: %d bytes, want %d", len(rebuilt), len(payload))
	}
	return outputs
}

// TestExactPageProtocol grades the page owner: source and offset on every page, executable
// continuations, rendered escaped bytes within the bound, reversible controls and invalid
// UTF-8, deterministic repetition, and refused cursors.
func TestExactPageProtocol(t *testing.T) {
	for _, tc := range []struct{ name, payload, encoding string }{
		{"ascii", strings.Repeat("plain text line\n", 600), "utf-8"},
		{"controls", strings.Repeat("\x00\x01\x1b\t\"\\", 1500), "utf-8"},
		{"unicode", strings.Repeat("é😀", 1500), "utf-8"},
		{"invalid utf-8", strings.Repeat("\xff\x80ab", 2000), "base64"},
		{"empty", "", "utf-8"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			first := pageAll(t, []byte(tc.payload), tc.encoding)
			if again := pageAll(t, []byte(tc.payload), tc.encoding); !reflect.DeepEqual(again, first) {
				t.Fatal("a repeated read at an unchanged source answered other pages")
			}
		})
	}

	payload := []byte("abcdef")
	first, err := PageRead(pageRequest(readFrame, payload))
	if err != nil {
		t.Fatal(err)
	}
	token := cell(decodePage(t, first.Output).row, "source")
	for _, tc := range []struct {
		name   string
		offset int64
		token  string
		frame  string
		want   string
	}{
		{"changed source", 0, strings.Repeat("0", 64), readFrame, "source changed since the cursor"},
		{"changed projection", 0, token, readFrame, "source changed since the cursor"},
		{"offset at the end", 6, token, readFrame, "cursor out of range"},
		{"offset past the end", 7, token, readFrame, "cursor out of range"},
		{"negative offset", -1, token, readFrame, "cursor out of range"},
		{"oversized metadata", 0, token, strings.Repeat("x", bounds.ResponseBytes-100) + "\n", "page metadata exceeds the response bound"},
		{"frame of many lines", 0, token, strings.Repeat("frame\n", bounds.ResponseLines), "page metadata exceeds the response bound"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := pageRequest(tc.frame, payload)
			request.Offset, request.Token, request.Continued = tc.offset, tc.token, true
			if tc.name == "changed projection" {
				request.Projection = "bench-other-projection/1"
			}
			page, err := PageRead(request)
			var refusal *ReadRefusal
			if !errors.As(err, &refusal) || refusal.Kind != tc.want || page.Output != "" {
				t.Fatalf("PageRead = %q, %v; want refusal %q and no content", page.Output, err, tc.want)
			}
		})
	}

	// A cursor inside a rune of a valid UTF-8 selection refuses before content, since no page
	// offered it. An invalid UTF-8 selection keeps its byte cursor anywhere.
	t.Run("cursor inside a rune", func(t *testing.T) {
		for _, tc := range []struct {
			name, payload string
			want          string
		}{
			{"valid selection", strings.Repeat("é😀", 40), "invalid cursor"},
			{"invalid selection", strings.Repeat("é😀", 40) + "\xff", ""},
			{"valid selection with an authored replacement character", strings.Repeat("é😀", 40) + "\xef\xbf\xbd", "invalid cursor"},
		} {
			request := pageRequest(readFrame, []byte(tc.payload))
			request.Offset, request.Token, request.Continued = int64(len("é😀")+1), request.identity(), true
			page, err := PageRead(request)
			var refusal *ReadRefusal
			switch {
			case tc.want != "" && (!errors.As(err, &refusal) || refusal.Kind != tc.want || page.Output != ""):
				t.Fatalf("%s: PageRead = %q, %v; want refusal %q and no content", tc.name, page.Output, err, tc.want)
			case tc.want == "" && (err != nil || cell(decodePage(t, page.Output).row, "encoding") != "base64"):
				t.Fatalf("%s: PageRead = %q, %v; want a base64 page at the byte cursor", tc.name, page.Output, err)
			}
		}
	})

	// An absent selection is one terminal page with present false and no content.
	t.Run("absent selection", func(t *testing.T) {
		request := pageRequest(readFrame, nil)
		request.Absent = true
		page, err := PageRead(request)
		if err != nil {
			t.Fatal(err)
		}
		decoded := decodePage(t, page.Output)
		if cell(decoded.row, "present") != "false" || cell(decoded.row, "complete") != "true" || len(decoded.data) != 0 || decoded.argv != nil || !strings.Contains(page.Output, "content[0]") {
			t.Fatalf("absent page = %q, want present false, complete, no content, and no successor", page.Output)
		}
	})

	// A frame leaves room for one byte of the next rune but not for the whole rune, so the
	// page refuses rather than split it.
	t.Run("window smaller than the next rune", func(t *testing.T) {
		payload := []byte(strings.Repeat("é😀", 3))
		request := pageRequest("", payload)
		token := request.identity()
		request.Offset, request.Token, request.Continued = int64(len("é")), token, true
		whole, _, err := request.render(token, payload[request.Offset:request.Offset+int64(len("😀"))])
		if err != nil {
			t.Fatal(err)
		}
		split, _, err := request.render(token, payload[request.Offset:request.Offset+1])
		if err != nil || len(split.Output) >= len(whole.Output) {
			t.Fatalf("a split-rune page of %d bytes is no shorter than the whole-rune page of %d bytes: %v", len(split.Output), len(whole.Output), err)
		}
		request.Frame = strings.Repeat("x", bounds.ResponseBytes-len(whole.Output)) + "\n"
		page, err := PageRead(request)
		var refusal *ReadRefusal
		if !errors.As(err, &refusal) || refusal.Kind != "page metadata exceeds the response bound" || page.Output != "" {
			t.Fatalf("PageRead = %q, %v; want the bounded refusal and no split rune", page.Output, err)
		}
	})
}
