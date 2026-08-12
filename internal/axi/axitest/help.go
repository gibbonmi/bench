// Package axitest owns shared executable-action proof mechanics for AXI tests.
package axitest

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

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
