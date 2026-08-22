package axitest

import (
	"reflect"
	"testing"
)

const helpEnvelopeDocument = "state: mapped\nrows[1]{story,seam,red_signal}:\n  1,command,observed red\nhelp[1]{cmd,why}:\n  bench coverage fixture,read the map\n"

func TestDecodeDocumentReportsBlocksInOrder(t *testing.T) {
	document, err := DecodeDocument(helpEnvelopeDocument)
	if err != nil {
		t.Fatalf("decode envelope: %v", err)
	}
	if want := []string{"state", "rows", "help"}; !reflect.DeepEqual(document.Blocks, want) {
		t.Fatalf("blocks = %q, want %q", document.Blocks, want)
	}
	rows, err := document.Rows("rows")
	if err != nil {
		t.Fatalf("rows: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one row", rows)
	}
	actions, err := document.HelpActions()
	if err != nil {
		t.Fatalf("help actions: %v", err)
	}
	want := []HelpAction{{Cmd: "bench coverage fixture", Why: "read the map"}}
	if !reflect.DeepEqual(actions, want) {
		t.Fatalf("help actions = %#v, want %#v", actions, want)
	}
}

func TestDecodeDocumentAcceptsZeroRowEnvelope(t *testing.T) {
	document, err := DecodeDocument("rows[0]{story,seam,red_signal}:\nhelp[0]{cmd,why}:\n")
	if err != nil {
		t.Fatalf("decode empty envelope: %v", err)
	}
	if want := []string{"rows", "help"}; !reflect.DeepEqual(document.Blocks, want) {
		t.Fatalf("blocks = %q, want %q", document.Blocks, want)
	}
	actions, err := document.HelpActions()
	if err != nil {
		t.Fatalf("help actions: %v", err)
	}
	if len(actions) != 0 {
		t.Fatalf("help actions = %#v, want none", actions)
	}
}

// TestDecodeDocumentRejectsUnstructuredOutput pins what a substring check cannot see:
// bytes outside the blocks, and a help table that is malformed, misschematized, or not
// the document's last block.
func TestDecodeDocumentRejectsUnstructuredOutput(t *testing.T) {
	for name, output := range map[string]string{
		"trailing prose":     helpEnvelopeDocument + "rendered by bench\n",
		"leading prose":      "rendered by bench\n" + helpEnvelopeDocument,
		"interleaved prose":  "state: mapped\nrendered by bench\nhelp[0]{cmd,why}:\n",
		"help row shortfall": "rows[0]{a}:\nhelp[2]{cmd,why}:\n  bench x,why\n",
	} {
		t.Run(name, func(t *testing.T) {
			if document, err := DecodeDocument(output); err == nil {
				t.Fatalf("DecodeDocument(%q) = %#v, want a decode failure", output, document)
			}
		})
	}
	for name, output := range map[string]string{
		"help not terminal":  "help[0]{cmd,why}:\nrows[0]{a}:\n",
		"trailing block":     helpEnvelopeDocument + "note: rendered by bench\n",
		"help schema short":  "rows[0]{a}:\nhelp[1]{cmd}:\n  bench x\n",
		"help schema wide":   "rows[0]{a}:\nhelp[1]{cmd,why,extra}:\n  bench x,why,more\n",
		"help fields swappd": "rows[0]{a}:\nhelp[1]{why,cmd}:\n  why,bench x\n",
		"empty help schema":  "rows[0]{a}:\nhelp[0]{cmd}:\n",
		"empty help cell":    "rows[0]{a}:\nhelp[1]{cmd,why}:\n  bench x,\"\"\n",
	} {
		t.Run(name, func(t *testing.T) {
			document, err := DecodeDocument(output)
			if err != nil {
				t.Fatalf("decode %q: %v", output, err)
			}
			if actions, err := document.HelpActions(); err == nil {
				t.Fatalf("HelpActions(%q) = %#v, want a help envelope failure", output, actions)
			}
		})
	}
}
