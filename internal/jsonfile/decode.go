// Package jsonfile decodes the strict machine-written JSON files Bench trusts as
// persisted lifecycle evidence.
package jsonfile

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
)

// Decode decodes one persisted JSON value into target.
func Decode(data []byte, target any) error {
	if len(data) == 0 || data[len(data)-1] != '\n' {
		return errors.New("persisted JSON requires a final newline")
	}
	return decodeDocument(data[:len(data)-1], target, nil)
}

// DecodeDocument strictly decodes one complete JSON document without imposing
// the final-newline framing required by Bench-owned persisted records.
func DecodeDocument(data []byte, target any) error {
	if len(data) == 0 {
		return errors.New("JSON document is empty")
	}
	return decodeDocument(data, target, nil)
}

// DecodeExactDocument also requires each struct key to match its declared name.
// Its schema must use named fields rather than anonymous field promotion.
func DecodeExactDocument(data []byte, target any) error {
	if len(data) == 0 {
		return errors.New("JSON document is empty")
	}
	return decodeDocument(data, target, reflect.TypeOf(target))
}

func decodeDocument(body []byte, target any, schema reflect.Type) error {
	scan := json.NewDecoder(bytes.NewReader(body))
	scan.UseNumber()
	if err := scanValue(scan, schema); err != nil {
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

func scanValue(decoder *json.Decoder, schema reflect.Type) error {
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
			child, err := objectField(schema, name)
			if err != nil {
				return err
			}
			if err := scanValue(decoder, child); err != nil {
				return err
			}
		}
	case '[':
		for decoder.More() {
			if err := scanValue(decoder, elementType(schema)); err != nil {
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
