package responsebound

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
)

// The cursor flags of every exact reader, and the cursor grammar that each reader's usage
// line carries. A supplied offset requires its source token, even at offset zero.
const (
	OffsetFlag  = "--offset"
	SourceFlag  = "--source"
	CursorUsage = "[" + OffsetFlag + " <byte-offset> " + SourceFlag + " <identity>]"
)

// readSchema names the page projection that each source identity binds.
const readSchema = "bench-read/1"

// The page encodings. Valid UTF-8 travels as a JSON string cell with every control escaped,
// and any other bytes travel as standard base64. Both decode to the exact selected bytes.
const (
	encodingUTF8   = "utf-8"
	encodingBase64 = "base64"
)

var (
	pageFields    = []string{"source", "kind", "present", "encoding", "offset", "next", "selected_bytes", "source_bytes", "complete"}
	contentFields = []string{"data"}
)

// ReadRequest is one exact page read of one selection.
type ReadRequest struct {
	// Command is the reader's argv after `bench`, without the cursor flags. A continuation
	// appends the cursor to it, and a refusal offers it bare as the first page.
	Command []string
	// Kind names the selection kind, and Selection is the selector's identity within it.
	Kind, Selection string
	// Projection names the versioned serialization that the payload holds, or is empty for
	// a payload of source bytes. A cursor from another projection refuses.
	Projection string
	// Scope names the repository or worktree identity, and Path the resolved path.
	Scope, Path string
	// Digest is the SHA-256 of the complete source bytes, and SourceSize is their count.
	Digest     [sha256.Size]byte
	SourceSize int64
	// Payload holds the selected bytes, and Size is their count.
	Payload io.ReaderAt
	Size    int64
	// Absent reports a selection that the source does not hold. Its one page has present
	// false, no content, and no successor, and its payload is empty.
	Absent bool
	// Offset and Token are the cursor. Continued reports that the call supplied them.
	Offset    int64
	Token     string
	Continued bool
	// Frame is the dispatcher's tree frame, which prints before the page and counts in the
	// page budget.
	Frame string
}

// ReadPage is one rendered page and its successor offset.
type ReadPage struct {
	Output   string
	Next     int64
	Complete bool
}

// ReadRefusal refuses one read before any content. Its message is the structured AXI error
// line.
type ReadRefusal struct{ Kind, Hint string }

func (r *ReadRefusal) Error() string { return toon.Errorf(r.Kind, r.Hint) }

// ParseCursor reads the cursor flags from a parsed flag set. The offset is a nonnegative
// decimal byte offset, and offset and source come together.
func ParseCursor(flags map[string]string) (int64, string, bool, error) {
	offset, hasOffset := flags[OffsetFlag]
	token, hasSource := flags[SourceFlag]
	if hasOffset != hasSource {
		return 0, "", false, errors.New(OffsetFlag + " and " + SourceFlag + " go together")
	}
	if !hasOffset {
		return 0, "", false, nil
	}
	value, err := decimal(offset)
	if err != nil {
		return 0, "", false, errors.New(OffsetFlag + " takes a nonnegative decimal byte offset")
	}
	return value, token, true, nil
}

// decimal reads one cursor or selector number: ASCII digits alone, with no sign, within
// int64. Every exact-reader number parses through it.
func decimal(text string) (int64, error) {
	if text == "" || strings.Trim(text, "0123456789") != "" {
		return 0, fmt.Errorf("%q is not a decimal", text)
	}
	return strconv.ParseInt(text, 10, 64)
}

// ParsePositive reads one positive selector number, such as a story number, through the one
// exact-reader decimal grammar.
func ParsePositive(text string) (int, error) {
	value, err := decimal(text)
	if err != nil || value < 1 {
		return 0, fmt.Errorf("%q is not a positive decimal", text)
	}
	return int(value), nil
}

// PageRead answers one bounded page of the selection, or a refusal before any content. The
// page carries the source identity, the selection kind, presence, the encoding, the byte
// offset, the selected and source byte counts, the next offset, and completion. A
// nonterminal page carries one executable continuation; a terminal page has none. The
// rendered page, its frame, and its continuation fit the shared response bound after every
// JSON and TOON escape, and each nonterminal page advances at least one source byte.
func PageRead(r ReadRequest) (ReadPage, error) {
	token := r.identity()
	if r.Continued && r.Token != token {
		return ReadPage{}, &ReadRefusal{Kind: "source changed since the cursor", Hint: "read again from the first page: " + r.firstPage()}
	}
	if r.Offset < 0 || r.Offset > r.Size || r.Offset == r.Size && r.Size > 0 {
		return ReadPage{}, &ReadRefusal{Kind: "cursor out of range", Hint: fmt.Sprintf("offset %d is outside the %d selected bytes; read again from the first page: %s", r.Offset, r.Size, r.firstPage())}
	}
	if err := r.runeBoundary(); err != nil {
		return ReadPage{}, err
	}
	window := make([]byte, min(r.Size-r.Offset, bounds.ResponseBytes))
	if len(window) > 0 {
		if n, err := r.Payload.ReadAt(window, r.Offset); n < len(window) {
			return ReadPage{}, &ReadRefusal{Kind: "selection unreadable", Hint: fmt.Sprint(err)}
		}
	}
	// The smallest page takes one whole unit: the first rune when the window opens with a
	// valid one, and otherwise one byte. A page never splits a valid rune.
	_, unit := utf8.DecodeRune(window)
	best, fits, err := r.render(token, window[:unit])
	if err != nil {
		return ReadPage{}, err
	}
	if !fits {
		return ReadPage{}, &ReadRefusal{Kind: "page metadata exceeds the response bound", Hint: "the frame, the page row, and the continuation leave no room for the next content byte or whole character"}
	}
	for low, high := unit, len(window); low < high; {
		middle := (low + high + 1) / 2
		page, fits, err := r.render(token, window[:middle])
		if err != nil {
			return ReadPage{}, err
		}
		if fits {
			low, best = middle, page
		} else {
			high = middle - 1
		}
	}
	return best, nil
}

// render renders the page that takes window, which starts at the request offset. It reports
// whether the frame and the page fit the response bound.
func (r ReadRequest) render(token string, window []byte) (ReadPage, bool, error) {
	consumed, encoding, data := encodeWindow(window, r.Offset+int64(len(window)) == r.Size)
	next := r.Offset + int64(consumed)
	complete := next == r.Size
	page, err := toon.TableTyped("page", pageFields, [][]any{{token, r.Kind, !r.Absent, encoding, r.Offset, next, r.Size, r.SourceSize, complete}})
	if err != nil {
		return ReadPage{}, false, &ReadRefusal{Kind: "page metadata unrepresentable", Hint: err.Error()}
	}
	var rows [][]any
	if consumed > 0 {
		rows = [][]any{{data}}
	}
	content, err := toon.TableTyped("content", contentFields, rows)
	if err != nil {
		return ReadPage{}, false, &ReadRefusal{Kind: "page content unrepresentable", Hint: err.Error()}
	}
	var actions []axi.Action
	if !complete {
		actions = append(actions, axi.ExecutableInvocation("read the next page", r.arguments(OffsetFlag, strconv.FormatInt(next, 10), SourceFlag, token)...))
	}
	help, err := axi.RenderHelp(actions)
	if err != nil || len(actions) > 0 && !strings.HasPrefix(help, "help[1]") {
		return ReadPage{}, false, &ReadRefusal{Kind: "continuation cannot be rendered", Hint: "a selector argument holds a character that an executable command cannot carry"}
	}
	output := page + content + help
	return ReadPage{Output: output, Next: next, Complete: complete}, Fits(r.Frame + output), nil
}

// runeBoundary refuses a cursor offset inside a rune of a valid UTF-8 selection, because each
// page of such a selection ends at a rune boundary, so no page offered that offset. An
// invalid UTF-8 selection travels as base64 with byte cursors, so any offset in it holds. Only
// an offset at a continuation byte reads the whole selection to decide its validity.
func (r ReadRequest) runeBoundary() error {
	if r.Offset == 0 || r.Offset >= r.Size {
		return nil
	}
	var at [1]byte
	if n, _ := r.Payload.ReadAt(at[:], r.Offset); n == 1 && utf8.RuneStart(at[0]) {
		return nil
	}
	reader := bufio.NewReader(io.NewSectionReader(r.Payload, 0, r.Size))
	for {
		char, size, err := reader.ReadRune()
		switch {
		case err == io.EOF:
			return &ReadRefusal{Kind: "invalid cursor", Hint: fmt.Sprintf("offset %d falls inside a character of the UTF-8 selection; read again from the first page: %s", r.Offset, r.firstPage())}
		case err != nil:
			return &ReadRefusal{Kind: "selection unreadable", Hint: err.Error()}
		case char == utf8.RuneError && size == 1:
			return nil
		}
	}
}

// encodeWindow encodes the page window. A nonfinal window that ends inside a rune stops at
// that rune's first byte, so a valid UTF-8 selection splits only at rune boundaries. Any
// other invalid window travels whole as base64. It answers the consumed byte count.
func encodeWindow(window []byte, final bool) (int, string, string) {
	text := window
	if !final && !utf8.Valid(text) {
		start := len(text) - 1
		for start > 0 && len(text)-start < utf8.UTFMax && !utf8.RuneStart(text[start]) {
			start--
		}
		if start > 0 && !utf8.FullRune(text[start:]) && utf8.Valid(text[:start]) {
			text = text[:start]
		}
	}
	if !utf8.Valid(text) {
		return len(window), encodingBase64, base64.StdEncoding.EncodeToString(window)
	}
	var cell strings.Builder
	encoder := json.NewEncoder(&cell)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(string(text))
	return len(text), encodingUTF8, strings.TrimSuffix(cell.String(), "\n")
}

// DecodeContent answers the selected bytes of one content cell in its page encoding.
func DecodeContent(encoding, cell string) ([]byte, error) {
	switch encoding {
	case encodingUTF8:
		var text string
		err := json.Unmarshal([]byte(cell), &text)
		return []byte(text), err
	case encodingBase64:
		return base64.StdEncoding.DecodeString(cell)
	}
	return nil, fmt.Errorf("unknown page encoding %q", encoding)
}

// identity hashes the length-framed schema, scope, path, selection, projection, and
// complete-source digest. It is a reader identity, not an evidence digest.
func (r ReadRequest) identity() string {
	hash := sha256.New()
	for _, field := range []string{readSchema, r.Scope, r.Path, r.Kind, r.Selection, r.Projection, hex.EncodeToString(r.Digest[:])} {
		fmt.Fprintf(hash, "%d:%s,", len(field), field)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

// arguments declares the reader command and extra as known invocation arguments.
func (r ReadRequest) arguments(extra ...string) []axi.InvocationArgument {
	var arguments []axi.InvocationArgument
	for _, argument := range append(append([]string(nil), r.Command...), extra...) {
		arguments = append(arguments, axi.KnownArgument(argument))
	}
	return arguments
}

// firstPage renders the exact first-page command of the selection.
func (r ReadRequest) firstPage() string {
	quoted := make([]string, len(r.Command))
	for i, argument := range r.Command {
		quoted[i] = sanitize.ShellQuote(argument)
	}
	return "bench " + strings.Join(quoted, " ")
}
