// Package jsonfile decodes the strict machine-written JSON files Bench trusts as
// persisted lifecycle evidence.
package jsonfile

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// Decode decodes one persisted JSON value into target.
func Decode(data []byte, target any) error {
	if len(data) == 0 || data[len(data)-1] != '\n' {
		return errors.New("persisted JSON requires a final newline")
	}
	body := data[:len(data)-1]
	scan := json.NewDecoder(bytes.NewReader(body))
	scan.UseNumber()
	if err := scanValue(scan); err != nil {
		return err
	}
	if err := rejectTrailing(scan); err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode persisted JSON: %w", err)
	}
	return rejectTrailing(decoder)
}

func scanValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("decode persisted JSON: %w", err)
	}
	delim, compound := token.(json.Delim)
	if !compound {
		return nil
	}
	switch delim {
	case '{':
		seen := map[string]bool{}
		for decoder.More() {
			token, err := decoder.Token()
			if err != nil {
				return fmt.Errorf("decode persisted JSON object: %w", err)
			}
			name, ok := token.(string)
			if !ok {
				return errors.New("decode persisted JSON object: field name is not a string")
			}
			if seen[name] {
				return fmt.Errorf("duplicate object field %q", name)
			}
			seen[name] = true
			if err := scanValue(decoder); err != nil {
				return err
			}
		}
	case '[':
		for decoder.More() {
			if err := scanValue(decoder); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("decode persisted JSON: unexpected delimiter %q", delim)
	}
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("decode persisted JSON: %w", err)
	}
	return nil
}

func rejectTrailing(decoder *json.Decoder) error {
	if _, err := decoder.Token(); errors.Is(err, io.EOF) {
		return nil
	} else if err != nil {
		return fmt.Errorf("trailing JSON data: %w", err)
	}
	return errors.New("trailing JSON value")
}
