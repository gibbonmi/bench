// Package axitest owns shared executable-action proof mechanics for AXI tests.
package axitest

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	toon "github.com/toon-format/toon-go"
)

// RecoverHelpCommandArgv returns the argv recovered from one rendered help command.
func RecoverHelpCommandArgv(output string) ([]string, error) {
	start := strings.Index(output, "help[")
	if start < 0 {
		return nil, errors.New("AXI output has no help block")
	}
	decoded, err := toon.DecodeString(output[start:])
	if err != nil {
		return nil, fmt.Errorf("decode help block: %w", err)
	}
	document, ok := decoded.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("decoded help document has type %T", decoded)
	}
	help, ok := document["help"].([]any)
	if !ok || len(help) != 1 {
		return nil, fmt.Errorf("decoded help rows = %#v, want one action", document["help"])
	}
	row, ok := help[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("decoded help row has type %T", help[0])
	}
	command, ok := row["cmd"].(string)
	if !ok || command == "" {
		return nil, errors.New("decoded help row has no command")
	}
	return shellArgv(command)
}

// Argv returns the argv that the action's rendered command expands to in a POSIX shell.
func (a HelpAction) Argv() ([]string, error) {
	return shellArgv(a.Cmd)
}

// shellArgv returns the argv that one rendered help command expands to in a POSIX shell.
func shellArgv(command string) ([]string, error) {
	recovered, err := exec.Command("sh", "-c", "set -- "+command+`; for argument do printf '%s\000' "$argument"; done`).Output()
	if err != nil {
		return nil, fmt.Errorf("execute decoded help command: %w", err)
	}
	if len(recovered) == 0 || recovered[len(recovered)-1] != 0 {
		return nil, errors.New("argv probe returned an unterminated result")
	}
	parts := bytes.Split(recovered[:len(recovered)-1], []byte{0})
	argv := make([]string, len(parts))
	for i := range parts {
		argv[i] = string(parts[i])
	}
	return argv, nil
}

// WalkPages follows one paged read from argv, which omits the leading `bench`, through each
// continuation that its pages offer, and answers the rebuilt content bytes and the page
// count. run executes one argv, makes the caller's own checks, and answers the page document
// without any identity row. decode decodes one content cell in the encoding of its page. A
// nonterminal page must offer an executable successor, and a walk that does not end within
// 1000 pages fails.
func WalkPages(t testing.TB, argv []string, run func(argv []string) string, decode func(encoding, data string) ([]byte, error)) ([]byte, int) {
	t.Helper()
	var rebuilt []byte
	for pages := 1; pages < 1000; pages++ {
		stdout := run(argv)
		document, err := DecodeDocument(stdout)
		if err != nil {
			t.Fatal(err)
		}
		rows, err := document.Rows("page")
		if err != nil {
			t.Fatal(err)
		}
		page := rows[0].(map[string]any)
		contents, err := document.Rows("content")
		if err != nil {
			t.Fatal(err)
		}
		for _, content := range contents {
			data, err := decode(page["encoding"].(string), content.(map[string]any)["data"].(string))
			if err != nil {
				t.Fatal(err)
			}
			rebuilt = append(rebuilt, data...)
		}
		if page["complete"] == true {
			return rebuilt, pages
		}
		next, err := RecoverHelpCommandArgv(stdout)
		if err != nil || len(next) < 2 || next[0] != "bench" {
			t.Fatalf("nonterminal page %q has no executable successor: %v", stdout, err)
		}
		argv = next[1:]
	}
	t.Fatalf("%q did not terminate", argv)
	return nil, 0
}
