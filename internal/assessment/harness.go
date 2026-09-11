package assessment

import (
	"encoding/json"
	"fmt"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/jsonfile"
	"os"
)

func collectHarness(r *Run, input HarnessInput) error {
	a, err := mapped(r, input.Mapping)
	if err != nil {
		return err
	}
	if input.SessionID != a.SessionID || input.Epoch < 0 || input.Sequence < 0 || input.EventID == "" {
		return fmt.Errorf("invalid native harness mapping")
	}
	expected := map[string]string{"total_token_usage": CumulativeMode, "last_token_usage": DeltaMode}
	if input.Format != "codex-token-count-v1" || expected[input.Counter] == "" || input.Mode != expected[input.Counter] {
		return fmt.Errorf("unrecognized harness counter semantics")
	}
	if err := regularInput(input.Path); err != nil {
		return err
	}
	read := bounds.ClassifyNoFollow(input.Path)
	if read.State != bounds.StateParsed {
		diagnostic(r, "Codex", input.Path, "native input "+string(read.State))
		return nil
	}
	var document map[string]json.RawMessage
	if err := jsonfile.DecodeDocument(read.Data, &document); err != nil {
		diagnostic(r, "Codex", input.Path, "malformed native input")
		return nil
	}
	var native struct {
		Type    string `json:"type"`
		Payload struct {
			Type string                     `json:"type"`
			Info map[string]json.RawMessage `json:"info"`
		} `json:"payload"`
	}
	if err := json.Unmarshal(read.Data, &native); err != nil || native.Type != "event_msg" || native.Payload.Type != "token_count" {
		diagnostic(r, "Codex", input.Path, "unsupported or unfinished native event")
		return nil
	}
	raw, ok := native.Payload.Info[input.Counter]
	if !ok || string(raw) == "null" {
		diagnostic(r, "Codex", input.Path, "selected counter missing")
		return nil
	}
	var counters struct {
		Input  *int64 `json:"input_tokens"`
		Cached *int64 `json:"cached_input_tokens"`
		Output *int64 `json:"output_tokens"`
	}
	if err := json.Unmarshal(raw, &counters); err != nil {
		diagnostic(r, "Codex", input.Path, "malformed native counters")
		return nil
	}
	event := Event{EventID: input.EventID, SessionID: input.SessionID, Epoch: input.Epoch, Sequence: input.Sequence, Mode: input.Mode, Counter: input.Counter, Reference: Reference{"Codex", input.Path + "#" + input.EventID}, Usage: Usage{InputTotal: counters.Input, InputCached: counters.Cached, Output: counters.Output, TotalSemantics: "inclusive"}}
	for _, old := range a.Usage {
		if old.EventID == event.EventID && old.SessionID == event.SessionID {
			if encoded(old) != encoded(event) {
				return fmt.Errorf("conflicting mapped harness event")
			}
			return nil
		}
	}
	a.Usage = append(a.Usage, event)
	return nil
}
func regularInput(path string) error {
	if err := noLinks(path); err != nil {
		return err
	}
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("native input is not a regular file")
	}
	return nil
}
