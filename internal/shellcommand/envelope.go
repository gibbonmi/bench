package shellcommand

import (
	"encoding/json"
	"fmt"
)

// EnvelopeForm selects the standard decoder target shape for a hook envelope.
// A form selects decoder semantics only, not required fields or guard authority.
type EnvelopeForm uint8

const (
	// CommandStruct decodes tool_input into a nested struct with a folded command key.
	CommandStruct EnvelopeForm = iota
	// InputMap decodes tool_input into a map of raw values with exact keys.
	InputMap
	// CwdInputMap decodes a cwd string and tool_input into a map of raw values.
	CwdInputMap
)

// Envelope holds the facts that one form decodes from a hook envelope.
// Callers read only the accessors for their selected form.
type Envelope struct {
	command string
	cwd     string
	input   map[string]json.RawMessage
}

// DecodeEnvelope decodes data once with encoding/json into the target shape of form.
// Folded root keys, ordered merges, null rules, and the earliest type error are the decoder's own.
// Any error returns the zero Envelope.
func DecodeEnvelope(data []byte, form EnvelopeForm) (Envelope, error) {
	switch form {
	case CommandStruct:
		var e struct {
			ToolInput struct {
				Command string `json:"command"`
			} `json:"tool_input"`
		}
		if err := json.Unmarshal(data, &e); err != nil {
			return Envelope{}, err
		}
		return Envelope{command: e.ToolInput.Command}, nil
	case InputMap:
		var e struct {
			ToolInput map[string]json.RawMessage `json:"tool_input"`
		}
		if err := json.Unmarshal(data, &e); err != nil {
			return Envelope{}, err
		}
		return Envelope{input: e.ToolInput}, nil
	case CwdInputMap:
		var e struct {
			Cwd       string                     `json:"cwd"`
			ToolInput map[string]json.RawMessage `json:"tool_input"`
		}
		if err := json.Unmarshal(data, &e); err != nil {
			return Envelope{}, err
		}
		return Envelope{cwd: e.Cwd, input: e.ToolInput}, nil
	}
	return Envelope{}, fmt.Errorf("unknown envelope form %d", form)
}

// Command returns the decoded tool_input command of a CommandStruct envelope.
func (e Envelope) Command() string { return e.command }

// Cwd returns the decoded cwd of a CwdInputMap envelope.
func (e Envelope) Cwd() string { return e.cwd }

// InputString decodes the tool_input map value at the exact key name as a string.
// A missing key returns present=false; a null value returns an empty present string.
// It requires no field and rejects no empty string or NUL byte.
func (e Envelope) InputString(name string) (value string, present bool, err error) {
	raw, ok := e.input[name]
	if !ok {
		return "", false, nil
	}
	if err := json.Unmarshal(raw, &value); err != nil {
		return "", true, err
	}
	return value, true, nil
}
