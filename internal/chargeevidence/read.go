package chargeevidence

import (
	"os"
	"strconv"
	"strings"
	"syscall"
)

// Cursor parts. Cursor.String and ParseCursor own the field order; the manifest stream uses
// source ordinal zero. A cursor names a position only; it never names a path and the store
// never records one.
const (
	CursorVersion  = "v1"
	cursorManifest = "m"
	cursorSource   = "s"
	cursorFields   = 5
)

// RefuseCursor names a malformed, foreign, or out-of-range cursor. RefuseSource names a
// well-formed source identifier the manifest does not declare.
const (
	RefuseCursor = "invalid-cursor"
	RefuseSource = "unknown-source"
)

// Cursor is one stateless position in an artifact's default evidence stream.
type Cursor struct {
	Identity string
	Ordinal  int
	Index    int
}

// String renders the canonical cursor.
func (c Cursor) String() string {
	kind := cursorSource
	if c.Ordinal == 0 {
		kind = cursorManifest
	}
	return strings.Join([]string{CursorVersion, strings.TrimPrefix(c.Identity, IdentityPrefix), kind,
		strconv.Itoa(c.Ordinal), strconv.Itoa(c.Index)}, ".")
}

// ParseCursor strictly parses a cursor for the artifact named by identity.
func ParseCursor(text, identity string) (Cursor, error) {
	parts := strings.Split(text, ".")
	if len(parts) != cursorFields || parts[0] != CursorVersion {
		return Cursor{}, refuse(RefuseCursor, "cursor does not use the %s grammar", CursorVersion)
	}
	if !validDigest(parts[1]) {
		return Cursor{}, refuse(RefuseCursor, "cursor identity is not 64 lowercase hexadecimal digits")
	}
	if IdentityPrefix+parts[1] != identity {
		return Cursor{}, refuse(RefuseCursor, "cursor belongs to another artifact")
	}
	ordinal, ok := ParseDecimal(parts[3])
	index, indexOK := ParseDecimal(parts[4])
	if !ok || !indexOK || ordinal > maxInteger || index > maxInteger {
		return Cursor{}, refuse(RefuseCursor, "cursor numbers are not canonical unsigned decimals")
	}
	switch {
	case parts[2] == cursorManifest && ordinal != 0:
		return Cursor{}, refuse(RefuseCursor, "a manifest cursor uses source ordinal zero")
	case parts[2] == cursorSource && ordinal == 0:
		return Cursor{}, refuse(RefuseCursor, "a source cursor names a source ordinal")
	case parts[2] != cursorManifest && parts[2] != cursorSource:
		return Cursor{}, refuse(RefuseCursor, "cursor stream is neither %s nor %s", cursorManifest, cursorSource)
	}
	return Cursor{Identity: identity, Ordinal: int(ordinal), Index: int(index)}, nil
}

// ParseDecimal parses a canonical unsigned decimal: digits only, no sign, and no leading
// zero except the single digit zero. It refuses a value outside the unsigned 64-bit range.
func ParseDecimal(text string) (uint64, bool) {
	if text == "" || len(text) > 1 && text[0] == '0' {
		return 0, false
	}
	for i := 0; i < len(text); i++ {
		if text[i] < '0' || text[i] > '9' {
			return 0, false
		}
	}
	value, err := strconv.ParseUint(text, 10, 64)
	return value, err == nil
}

// Artifact is one opened published pack. It holds the shared operation lock until Close.
type Artifact struct {
	identity  string
	dir       *os.Root
	operation *lock
	file      *os.File
	info      os.FileInfo
	manifest  []byte
	m         Manifest
	starts    []int64
}

// Open validates the identifier before any path use, then opens and strictly validates the
// artifact's complete manifest. An absent store refuses without being created.
func (s *Store) Open(identity string) (*Artifact, error) {
	if !ValidIdentity(identity) {
		return nil, refuse(RefuseIdentifier, "artifact identifier is not sha256 followed by 64 lowercase hexadecimal digits")
	}
	dir, err := s.openDir(false)
	if err != nil {
		return nil, err
	}
	a := &Artifact{identity: identity, dir: dir}
	if a.operation, err = acquire(dir, OperationLockName, syscall.LOCK_SH, false); err != nil {
		a.Close()
		return nil, err
	}
	if a.file, a.info, err = openRegular(dir, packName(identity)); err != nil {
		a.Close()
		return nil, err
	}
	if err := a.readManifest(); err != nil {
		a.Close()
		return nil, err
	}
	return a, nil
}

func (a *Artifact) readManifest() error {
	size := uint64(a.info.Size())
	if size < HeaderBytes {
		return refuse(RefuseTruncated, "pack holds %d bytes, fewer than the %d-byte header", size, HeaderBytes)
	}
	header, err := a.readAt(0, HeaderBytes)
	if err != nil {
		return err
	}
	end, err := headerEnd(header, size)
	if err != nil {
		return err
	}
	if a.manifest, err = a.readAt(HeaderBytes, int(end-HeaderBytes)); err != nil {
		return err
	}
	if a.m, err = manifestAt(a.manifest, a.identity, end, size); err != nil {
		return err
	}
	offset := int64(end)
	for _, s := range a.m.Sources {
		a.starts = append(a.starts, offset)
		offset += int64(s.Bytes)
	}
	return unchanged(a.dir, packName(a.identity), a.file, a.info)
}

func (a *Artifact) readAt(offset int64, length int) ([]byte, error) {
	buffer := make([]byte, length)
	if _, err := a.file.ReadAt(buffer, offset); err != nil {
		return nil, refuse(RefuseTruncated, "pack is shorter than its validated layout: %v", err)
	}
	return buffer, nil
}

// Close releases the file and the shared operation lock.
func (a *Artifact) Close() {
	if a.file != nil {
		a.file.Close()
	}
	a.operation.release()
	if a.dir != nil {
		a.dir.Close()
	}
}

// Manifest returns the validated manifest and its exact bytes.
func (a *Artifact) Manifest() (Manifest, []byte) { return a.m, a.manifest }

// Fragment is one bounded piece of the default evidence stream.
type Fragment struct {
	Stream, Source string
	Index, Offset  int
	Bytes, Total   int
	SHA256         string
	Content        []byte
	Next           *Cursor
}

// First returns the cursor that starts the default stream: the first manifest fragment.
func (a *Artifact) First() Cursor { return Cursor{Identity: a.identity} }

// SourceOrdinal returns the one-based ordinal of the source the manifest declares under id.
func (a *Artifact) SourceOrdinal(id string) (int, error) {
	for i, source := range a.m.Sources {
		if source.ID == id {
			return i + 1, nil
		}
	}
	return 0, refuse(RefuseSource, "the manifest declares no source %s", id)
}

// Read returns the fragment at cursor and its successor in the default stream.
func (a *Artifact) Read(c Cursor) (Fragment, error) { return a.read(c, false) }

// ReadWithin returns the fragment at cursor with a successor that stays inside the cursor's
// own source, so an explicit source read ends after that source instead of continuing into
// the default stream.
func (a *Artifact) ReadWithin(c Cursor) (Fragment, error) { return a.read(c, true) }

// read returns the fragment at cursor. A source page is read from the opened file, checked
// against its manifest digest, and returned only when the directory still names the same
// unchanged file. within keeps the successor inside the cursor's own source.
func (a *Artifact) read(c Cursor, within bool) (Fragment, error) {
	if c.Identity != a.identity {
		return Fragment{}, refuse(RefuseCursor, "cursor belongs to another artifact")
	}
	if within && c.Ordinal == 0 {
		return Fragment{}, refuse(RefuseCursor, "a source read takes a source cursor")
	}
	if c.Ordinal == 0 {
		return a.manifestFragment(c.Index)
	}
	if c.Ordinal > len(a.m.Sources) {
		return Fragment{}, refuse(RefuseCursor, "cursor names source ordinal %d beyond %d sources", c.Ordinal, len(a.m.Sources))
	}
	source := a.m.Sources[c.Ordinal-1]
	pages := a.pagesOf(source.ID)
	if c.Index >= len(pages) {
		return Fragment{}, refuse(RefuseCursor, "cursor names page %d beyond %d pages of %s", c.Index, len(pages), source.ID)
	}
	page := pages[c.Index]
	content, err := a.readAt(a.starts[c.Ordinal-1]+int64(page.Offset), page.Bytes)
	if err != nil {
		return Fragment{}, err
	}
	if err := checkPage(source.ID, page, content); err != nil {
		return Fragment{}, err
	}
	if err := unchanged(a.dir, packName(a.identity), a.file, a.info); err != nil {
		return Fragment{}, err
	}
	next := a.after(c.Ordinal, c.Index)
	if within && next != nil && next.Ordinal != c.Ordinal {
		next = nil
	}
	return Fragment{Stream: StreamSource, Source: source.ID, Index: page.Index, Offset: page.Offset, Bytes: page.Bytes,
		Total: source.Bytes, SHA256: page.SHA256, Content: content, Next: next}, nil
}

func (a *Artifact) manifestFragment(index int) (Fragment, error) {
	sizes := pageSizes(a.manifest)
	if index >= len(sizes) {
		return Fragment{}, refuse(RefuseCursor, "cursor names manifest fragment %d beyond %d fragments", index, len(sizes))
	}
	offset := 0
	for _, size := range sizes[:index] {
		offset += size
	}
	content := a.manifest[offset : offset+sizes[index]]
	next := &Cursor{Identity: a.identity, Index: index + 1}
	if index+1 == len(sizes) {
		next = a.after(0, 0)
	}
	return Fragment{Stream: StreamManifest, Index: index, Offset: offset, Bytes: len(content), Total: len(a.manifest),
		SHA256: Digest(content), Content: content, Next: next}, nil
}

// after returns the default-stream position that follows page index of source ordinal,
// or nil at the end. Ordinal zero with any index means the end of the manifest stream.
func (a *Artifact) after(ordinal, index int) *Cursor {
	if ordinal > 0 && index+1 < len(a.pagesOf(a.m.Sources[ordinal-1].ID)) {
		return &Cursor{Identity: a.identity, Ordinal: ordinal, Index: index + 1}
	}
	for next := ordinal + 1; next <= len(a.m.Sources); next++ {
		if len(a.pagesOf(a.m.Sources[next-1].ID)) > 0 {
			return &Cursor{Identity: a.identity, Ordinal: next}
		}
	}
	return nil
}

func (a *Artifact) pagesOf(id string) []Page {
	var pages []Page
	for _, page := range a.m.Pages {
		if page.Source == id {
			pages = append(pages, page)
		}
	}
	return pages
}

// Verified counts one artifact's complete verification.
type Verified struct {
	Evidence       string
	Pages, Sources int
}

// Verify reads every page of every source, checks each page digest, and checks each
// reconstructed source against its declared length and digest. Opening the artifact already
// validated the complete manifest against the expected identity, so a successful return
// covers every stored byte.
func (a *Artifact) Verify() (Verified, error) {
	verified := Verified{Evidence: a.identity}
	for ordinal, source := range a.m.Sources {
		body := make([]byte, 0, source.Bytes)
		for _, page := range a.pagesOf(source.ID) {
			content, err := a.readAt(a.starts[ordinal]+int64(page.Offset), page.Bytes)
			if err != nil {
				return Verified{}, err
			}
			if err := checkPage(source.ID, page, content); err != nil {
				return Verified{}, err
			}
			body = append(body, content...)
			verified.Pages++
		}
		if err := checkSource(source, body); err != nil {
			return Verified{}, err
		}
		verified.Sources++
	}
	if err := unchanged(a.dir, packName(a.identity), a.file, a.info); err != nil {
		return Verified{}, err
	}
	return verified, nil
}

// Encode renders the verified response through the registered schema.
func (v Verified) Encode() (string, error) {
	return encodeResponse(blockVerified, []any{v.Evidence, true, v.Pages, v.Sources, DeliveryUnverified})
}

// Current is the registered current-action binding response row. Assignment names the
// current assignment, never the assignment that prepared the artifact.
type Current struct {
	Evidence, Assignment, Base, SourceTip string
}

// Encode renders the current response through the registered schema.
func (c Current) Encode() (string, error) {
	return encodeResponse(blockCurrent, []any{c.Evidence, c.Assignment, c.Base, c.SourceTip, true, DeliveryUnverified})
}

// Prepared is the registered prepared response row.
type Prepared struct {
	Evidence, Mode, Base, SourceTip, Assignment string
	Sources, Pages, ManifestBytes               int
	Next                                        string
}

// Encode renders the prepared response through the registered schema.
func (p Prepared) Encode() (string, error) {
	return encodeResponse(blockPrepared, []any{p.Evidence, p.Mode, p.Base, p.SourceTip, p.Assignment,
		SelectionReference, SourceID(1), p.Sources, p.Pages, p.ManifestBytes, true, DeliveryUnverified, p.Next})
}

// Encode renders one fragment response through the registered schema. next is the exact
// successor command, or empty at the end of the stream.
func (f Fragment) Encode(identity, next string) (string, error) {
	return encodeResponse(blockPage, []any{identity, f.Stream, f.Source, f.Index, f.Offset, f.Bytes, f.Total,
		f.SHA256, string(f.Content), true, f.Next == nil, next})
}

func encodeResponse(name string, row []any) (string, error) {
	for _, block := range ResponseBlocks {
		if block.Name == name {
			out, err := encodeBlocks([]Block{block}, map[string][][]any{name: {row}})
			return string(out), err
		}
	}
	return "", refuse(RefuseSchema, "unregistered response block %q", name)
}
