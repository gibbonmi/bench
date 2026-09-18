package chargeevidence

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	toonlib "github.com/toon-format/toon-go"

	"github.com/gibbonmi/bench/internal/toon"
)

// Refusal classes. A strict read reports exactly one class for the first failed
// predicate, so a caller and a test can name the refused property.
const (
	RefuseIdentifier       = "invalid-identifier"
	RefuseIdentity         = "identity-mismatch"
	RefuseMarker           = "invalid-marker"
	RefuseContainerVersion = "unsupported-container-version"
	RefuseReserved         = "nonzero-reserved"
	RefuseOverflow         = "length-overflow"
	RefuseTruncated        = "truncated"
	RefuseTrailing         = "trailing-data"
	RefuseProfile          = "unsupported-profile"
	RefuseNoncanonical     = "noncanonical"
	RefuseSchema           = "schema"
	RefuseType             = "cell-type"
	RefuseInvalid          = "invalid-manifest"
	RefusePageOverlap      = "page-overlap"
	RefusePageGap          = "page-gap"
	RefusePageRange        = "page-range"
	RefusePageSplit        = "page-split"
	RefusePageDigest       = "page-digest"
	RefuseSourceDigest     = "source-digest"
	RefuseSourceBytes      = "source-bytes"
	RefuseCandidate        = "invalid-candidate"
)

// Refusal is the typed error every strict operation returns.
type Refusal struct {
	Class  string
	Detail string
}

func (r *Refusal) Error() string { return r.Class + ": " + r.Detail }

func refuse(class, format string, args ...any) error {
	return &Refusal{Class: class, Detail: fmt.Sprintf(format, args...)}
}

// maxInteger is the largest integer the shared TOON encoder emits as a bare number.
const maxInteger = 1<<53 - 1

// encodeBlocks renders every registered block once, in registry order, through the shared
// TOON adapter. A block without rows still renders its schema header.
func encodeBlocks(blocks []Block, values map[string][][]any) ([]byte, error) {
	var b strings.Builder
	for _, block := range blocks {
		text, err := toon.TableTyped(block.Name, block.FieldNames(), values[block.Name])
		if err != nil {
			return nil, refuse(RefuseCandidate, "%s block: %v", block.Name, err)
		}
		b.WriteString(text)
	}
	return []byte(b.String()), nil
}

// decodeBlocks strictly decodes a canonical block document. It refuses unknown or missing
// blocks and columns and every cell whose decoded type differs from the registry. The
// canonical re-encode must then reproduce the exact input bytes, which refuses reordered
// blocks or columns, duplicates, and every alternative serialization.
func decodeBlocks(blocks []Block, data []byte) (map[string][][]any, error) {
	if !utf8.Valid(data) {
		return nil, refuse(RefuseNoncanonical, "document is not valid UTF-8")
	}
	decoded, err := toonlib.Decode(data)
	if err != nil {
		return nil, refuse(RefuseNoncanonical, "document does not decode: %v", err)
	}
	document, ok := decoded.(map[string]any)
	if !ok {
		return nil, refuse(RefuseSchema, "document decodes as %T, want blocks", decoded)
	}
	registered := map[string]bool{}
	for _, block := range blocks {
		registered[block.Name] = true
	}
	for name := range document {
		if !registered[name] {
			return nil, refuse(RefuseSchema, "unknown block %q", name)
		}
	}
	values := make(map[string][][]any, len(blocks))
	for _, block := range blocks {
		raw, present := document[block.Name]
		if !present {
			return nil, refuse(RefuseSchema, "missing block %q", block.Name)
		}
		list, ok := raw.([]any)
		if !ok {
			return nil, refuse(RefuseSchema, "block %q decodes as %T, want a table", block.Name, raw)
		}
		rows := make([][]any, len(list))
		for i, item := range list {
			row, err := decodeRow(block, i, item)
			if err != nil {
				return nil, err
			}
			rows[i] = row
		}
		values[block.Name] = rows
	}
	canonical, err := encodeBlocks(blocks, values)
	if err != nil {
		return nil, refuse(RefuseNoncanonical, "document does not re-encode: %v", err)
	}
	if !bytes.Equal(canonical, data) {
		return nil, refuse(RefuseNoncanonical, "document bytes differ from their canonical encoding")
	}
	return values, nil
}

func decodeRow(block Block, index int, item any) ([]any, error) {
	object, ok := item.(map[string]any)
	if !ok {
		return nil, refuse(RefuseSchema, "%s row %d decodes as %T, want a record", block.Name, index, item)
	}
	if len(object) != len(block.Fields) {
		return nil, refuse(RefuseSchema, "%s row %d has %d columns, want %d", block.Name, index, len(object), len(block.Fields))
	}
	row := make([]any, len(block.Fields))
	for j, field := range block.Fields {
		value, present := object[field.Name]
		if !present {
			return nil, refuse(RefuseSchema, "%s row %d has no %s column", block.Name, index, field.Name)
		}
		cell, err := typedCell(field, value)
		if err != nil {
			return nil, refuse(RefuseType, "%s row %d %s: %v", block.Name, index, field.Name, err)
		}
		row[j] = cell
	}
	return row, nil
}

func typedCell(field Field, value any) (any, error) {
	switch field.Type {
	case Integer:
		number, ok := value.(float64)
		if !ok {
			return nil, fmt.Errorf("decodes as %T, want an integer", value)
		}
		if number < 0 || number > maxInteger || number != math.Trunc(number) {
			return nil, fmt.Errorf("%v is not an unsigned integer in range", number)
		}
		return int(number), nil
	case Boolean:
		flag, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("decodes as %T, want a boolean", value)
		}
		return flag, nil
	default:
		text, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("decodes as %T, want a string", value)
		}
		return text, nil
	}
}
