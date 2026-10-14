package ownership

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/maps"
)

// MetadataLabel is the info string of the one planning-facts fence a spec may carry.
const MetadataLabel = "bench-planning-facts"

// metadataID names every diagnostic of the inventory itself.
const metadataID = "metadata"

// ParseManifest reads the spec's one bench-planning-facts fence. The fence lines come from
// the shared FieldScan fence state, so a labeled marker inside another fence opens nothing.
// The payload is bounded by ControlRecordLimit through the classified byte reader. Every
// grammar fault is one diagnostic on a present manifest, so it keeps the inventory
// incomplete without refusing the closure. A spec with no fence answers a manifest that is
// not present.
func ParseManifest(content []byte) Manifest {
	payload, present, err := metadataPayload(content)
	if !present {
		return Manifest{}
	}
	if err == nil {
		var manifest Manifest
		if manifest, err = decodeManifest(payload, coverage.MapRows(content)); err == nil {
			manifest.Present = true
			return manifest
		}
	}
	return Manifest{Present: true, Diagnostics: []Diagnostic{{ID: metadataID, Reason: err.Error()}}}
}

// metadataPayload is the bounded payload of the one metadata fence. present is false when
// the spec opens no such fence.
func metadataPayload(content []byte) ([]byte, bool, error) {
	lines, _ := maps.FieldScan{}.Scan(content)
	var body []string
	opened := false
	for i := 0; i < len(lines); i++ {
		if lines[i].FenceRole != maps.FenceOpen || lines[i].FenceLabel != MetadataLabel {
			continue
		}
		if opened {
			return nil, true, errors.New(MetadataLabel + " fence appears more than once")
		}
		opened = true
		closed := false
		for i++; i < len(lines) && !closed; i++ {
			switch {
			case lines[i].FenceRole == maps.FenceClose && lines[i].FenceLabel != "":
				return nil, true, errors.New(MetadataLabel + " fence holds another fence marker")
			case lines[i].FenceRole == maps.FenceClose:
				closed = true
			default:
				body = append(body, lines[i].Text)
			}
		}
		if !closed {
			return nil, true, errors.New(MetadataLabel + " fence is not closed")
		}
		i--
	}
	if !opened {
		return nil, false, nil
	}
	read := bounds.ClassifyBytes(bounds.Read(strings.NewReader(strings.Join(body, "\n")), bounds.ControlRecordLimit))
	if read.State != bounds.StateParsed {
		return nil, true, fmt.Errorf("%s payload is %s: %s", MetadataLabel, read.State, read.Reason)
	}
	return read.Data, true, nil
}

// fields decodes one JSON object whose key set must be exactly keys. where prefixes each
// fault, so a fault names the item that holds it.
func fields(raw json.RawMessage, where string, keys ...string) (map[string]json.RawMessage, error) {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(raw, &object); err != nil || object == nil {
		return nil, fmt.Errorf("%sis not a JSON object", where)
	}
	var unknown []string
	for key := range object {
		if !slices.Contains(keys, key) {
			unknown = append(unknown, key)
		}
	}
	if len(unknown) > 0 {
		slices.Sort(unknown)
		return nil, fmt.Errorf("%sunknown field %q", where, unknown[0])
	}
	for _, key := range keys {
		if _, ok := object[key]; !ok {
			return nil, fmt.Errorf("%smissing field %q", where, key)
		}
	}
	return object, nil
}

// decodeInto decodes one field value into target and names the field on a type fault.
func decodeInto(object map[string]json.RawMessage, where, key string, target any) error {
	if err := json.Unmarshal(object[key], target); err != nil {
		return fmt.Errorf("%sinvalid %s: %v", where, key, err)
	}
	return nil
}

// enum refuses a value outside its exact spellings. No case folding applies.
func enum(where, field, value string, allowed ...string) error {
	if !slices.Contains(allowed, value) {
		return fmt.Errorf("%sinvalid %s %q", where, field, value)
	}
	return nil
}
