package jsonfile

import (
	"strings"
	"testing"
)

type testDocument struct {
	Name  string `json:"name"`
	Child struct {
		Value int `json:"value"`
	} `json:"child"`
}

func TestDecodePersistedJSONContract(t *testing.T) {
	t.Run("typed value with legacy whitespace", func(t *testing.T) {
		var got testDocument
		if err := Decode([]byte(" \n { \"name\" : \"legacy\", \"child\" : { \"value\" : 7 } } \t\n"), &got); err != nil {
			t.Fatal(err)
		}
		if got.Name != "legacy" || got.Child.Value != 7 {
			t.Fatalf("Decode = %#v, want typed legacy value", got)
		}
	})

	invalid := []struct {
		name, body, errorPart string
	}{
		{"missing final newline", `{"name":"x","child":{"value":1}}`, "final newline"},
		{"duplicate field", "{\"name\":\"x\",\"name\":\"y\",\"child\":{\"value\":1}}\n", "duplicate object field"},
		{"nested duplicate field", "{\"name\":\"x\",\"child\":{\"value\":1,\"value\":2}}\n", "duplicate object field"},
		{"unknown field", "{\"name\":\"x\",\"extra\":true,\"child\":{\"value\":1}}\n", "unknown field"},
		{"trailing value", "{\"name\":\"x\",\"child\":{\"value\":1}} {}\n", "trailing JSON"},
		{"trailing bytes", "{\"name\":\"x\",\"child\":{\"value\":1}} nope\n", "trailing JSON"},
		{"wrong field type", "{\"name\":7,\"child\":{\"value\":1}}\n", "cannot unmarshal"},
	}
	for _, tc := range invalid {
		t.Run(tc.name, func(t *testing.T) {
			var got testDocument
			err := Decode([]byte(tc.body), &got)
			if err == nil || !strings.Contains(err.Error(), tc.errorPart) {
				t.Fatalf("Decode error = %v, want containing %q", err, tc.errorPart)
			}
		})
	}
}
