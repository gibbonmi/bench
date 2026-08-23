package axitest

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	toonlib "github.com/toon-format/toon-go"

	"github.com/gibbonmi/bench/internal/toon"
)

// helpFields is the help block's schema, in the order the AXI envelope emits it.
var helpFields = []string{"cmd", "why"}

// Document is a command's complete stdout, decoded as one TOON document: every byte
// took part in the decode, so anything before, between, or after the blocks fails
// rather than passing a substring check. Blocks names the document's top-level keys in
// emission order and Values holds the decode.
type Document struct {
	Output string
	Blocks []string
	Values map[string]any
}

// HelpAction is one recovered row of the terminal help table.
type HelpAction struct {
	Cmd, Why string
}

// DecodeDocument decodes a command's whole stdout with the official decoder. Block order
// is recovered by decoding growing line prefixes rather than by reading the block grammar
// here: a key is placed where the shortest decodable prefix that carries it ends, so this
// helper stays a consumer of the library's parse instead of a second derivation of TOON.
func DecodeDocument(output string) (Document, error) {
	values, err := decodeDocumentMap(output)
	if err != nil {
		return Document{}, err
	}
	blocks, err := documentBlockOrder(output)
	if err != nil {
		return Document{}, err
	}
	if !reflect.DeepEqual(sorted(blocks), sorted(keysOf(values))) {
		return Document{}, fmt.Errorf("recovered block order %q does not cover decoded keys %q", blocks, keysOf(values))
	}
	return Document{Output: output, Blocks: blocks, Values: values}, nil
}

// Rows returns the decoded rows of one table block.
func (d Document) Rows(block string) ([]any, error) {
	rows, ok := d.Values[block].([]any)
	if !ok {
		return nil, fmt.Errorf("block %q decoded as %T, want a table", block, d.Values[block])
	}
	return rows, nil
}

// HelpActions returns the envelope's help rows. It requires the help block to be
// terminal and schema-correct: the document has to end with exactly what the kit's
// encoder renders for the decoded rows under the `{cmd,why}` schema. A reordered,
// misnamed, or non-final help table fails, including the zero-row case, whose schema
// survives only in the header line.
func (d Document) HelpActions() ([]HelpAction, error) {
	if len(d.Blocks) == 0 || d.Blocks[len(d.Blocks)-1] != "help" {
		return nil, fmt.Errorf("document blocks = %q, want a terminal help block", d.Blocks)
	}
	decoded, err := d.Rows("help")
	if err != nil {
		return nil, err
	}
	actions := make([]HelpAction, len(decoded))
	cells := make([][]string, len(decoded))
	for i, row := range decoded {
		fields, ok := row.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("help row %d decoded as %T, want an object", i, row)
		}
		if !reflect.DeepEqual(sorted(keysOf(fields)), sorted(helpFields)) {
			return nil, fmt.Errorf("help row %d has fields %q, want %q", i, sorted(keysOf(fields)), sorted(helpFields))
		}
		action := HelpAction{}
		for _, field := range helpFields {
			value, ok := fields[field].(string)
			if !ok || value == "" {
				return nil, fmt.Errorf("help row %d field %q = %#v, want a non-empty string", i, field, fields[field])
			}
			switch field {
			case "cmd":
				action.Cmd = value
			case "why":
				action.Why = value
			}
		}
		actions[i] = action
		cells[i] = []string{action.Cmd, action.Why}
	}
	canonical, err := toon.Table("help", helpFields, cells)
	if err != nil {
		return nil, fmt.Errorf("re-render help block: %w", err)
	}
	if !strings.HasSuffix(d.Output, canonical) {
		return nil, fmt.Errorf("document does not end with the canonical help block %q", canonical)
	}
	return actions, nil
}

func decodeDocumentMap(output string) (map[string]any, error) {
	decoded, err := toonlib.DecodeString(output)
	if err != nil {
		return nil, fmt.Errorf("decode document: %w", err)
	}
	values, ok := decoded.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("document decoded as %T, want an object", decoded)
	}
	if len(values) == 0 {
		return nil, errors.New("document decoded as no blocks")
	}
	return values, nil
}

// documentBlockOrder walks the document one line at a time and records each key at the
// first prefix that decodes with it present. A prefix cut inside a block does not decode
// and is skipped, so keys land in emission order.
func documentBlockOrder(output string) ([]string, error) {
	var blocks []string
	seen := map[string]bool{}
	prefix := ""
	for _, line := range strings.SplitAfter(output, "\n") {
		if line == "" {
			continue
		}
		prefix += line
		values, err := decodeDocumentMap(prefix)
		if err != nil {
			continue
		}
		var revealed []string
		for _, key := range keysOf(values) {
			if !seen[key] {
				revealed = append(revealed, key)
			}
		}
		if len(revealed) > 1 {
			return nil, fmt.Errorf("prefix ending %q reveals %q at once, so block order is ambiguous", line, sorted(revealed))
		}
		for _, key := range revealed {
			seen[key] = true
			blocks = append(blocks, key)
		}
	}
	return blocks, nil
}

func keysOf(values map[string]any) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}

func sorted(values []string) []string {
	out := append([]string(nil), values...)
	sort.Strings(out)
	return out
}
