package jsonfile

import (
	"errors"
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

const testDocumentBody = `{"name":"x","child":{"value":1}}`

// A refusal case names the category that errors.Is must find. A case with no category
// names the message part of its refusal, and a case with neither is an acceptance.
func TestDocumentStructuralContract(t *testing.T) {
	cases := []struct {
		name, body string
		target     any
		category   error
		errorPart  string
	}{
		{name: "nested duplicate in an array element", body: `[{"name":"x","name":"y"}]`, target: new([]testDocument), category: ErrDuplicateKey},
		{name: "repeated key with an escaped spelling", body: `{"name":"x","name":"y"}`, target: new(testDocument), category: ErrDuplicateKey},
		{name: "unknown nested field", body: `{"name":"x","child":{"value":1,"extra":2}}`, target: new(testDocument), errorPart: `unknown field "extra"`},
		{name: "second object", body: testDocumentBody + ` {}`, target: new(testDocument), category: ErrTrailingData},
		{name: "second array", body: testDocumentBody + ` []`, target: new(testDocument), category: ErrTrailingData},
		{name: "second string", body: testDocumentBody + ` "s"`, target: new(testDocument), category: ErrTrailingData},
		{name: "second number", body: testDocumentBody + ` 7`, target: new(testDocument), category: ErrTrailingData},
		{name: "second true", body: testDocumentBody + ` true`, target: new(testDocument), category: ErrTrailingData},
		{name: "second false", body: testDocumentBody + ` false`, target: new(testDocument), category: ErrTrailingData},
		{name: "second null", body: testDocumentBody + ` null`, target: new(testDocument), category: ErrTrailingData},
		{name: "malformed word suffix", body: testDocumentBody + ` nope`, target: new(testDocument), category: ErrTrailingData},
		{name: "stray closing brace", body: testDocumentBody + `}`, target: new(testDocument), category: ErrTrailingData},
		{name: "stray comma", body: testDocumentBody + ` ,`, target: new(testDocument), category: ErrTrailingData},
		{name: "whitespace after the value", body: " \n" + testDocumentBody + " \t\r\n ", target: new(testDocument)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := DecodeDocument([]byte(tc.body), tc.target)
			switch {
			case tc.category != nil:
				if !errors.Is(err, tc.category) {
					t.Fatalf("DecodeDocument error = %v, want category %v", err, tc.category)
				}
			case tc.errorPart != "":
				if err == nil || !strings.Contains(err.Error(), tc.errorPart) {
					t.Fatalf("DecodeDocument error = %v, want containing %q", err, tc.errorPart)
				}
			case err != nil:
				t.Fatalf("DecodeDocument error = %v, want acceptance", err)
			}
		})
	}
}

func TestDocumentFramingCompatibility(t *testing.T) {
	t.Run("record refuses a missing final newline", func(t *testing.T) {
		// The final byte is whitespace, so a decoder without the framing check reads a valid value.
		if err := Decode([]byte(testDocumentBody+" "), new(testDocument)); err == nil || !strings.Contains(err.Error(), "final newline") {
			t.Fatalf("Decode error = %v, want the final-newline refusal", err)
		}
	})
	t.Run("document accepts a missing final newline", func(t *testing.T) {
		var got testDocument
		if err := DecodeDocument([]byte(testDocumentBody), &got); err != nil {
			t.Fatalf("DecodeDocument error = %v, want acceptance", err)
		}
		if got.Name != "x" || got.Child.Value != 1 {
			t.Fatalf("DecodeDocument = %#v, want the typed value", got)
		}
	})
	t.Run("exact document refuses a case alias", func(t *testing.T) {
		alias := []byte(`{"Name":"x","child":{"value":1}}`)
		if err := DecodeDocument(alias, new(testDocument)); err != nil {
			t.Fatalf("DecodeDocument error = %v, want the case alias accepted", err)
		}
		if err := DecodeExactDocument(alias, new(testDocument)); err == nil || !strings.Contains(err.Error(), `unknown exact JSON field "Name"`) {
			t.Fatalf("DecodeExactDocument error = %v, want the case alias refused", err)
		}
	})
}

func TestDocumentErrorCategories(t *testing.T) {
	decoders := map[string]func([]byte, any) error{
		"Decode":              Decode,
		"DecodeDocument":      DecodeDocument,
		"DecodeExactDocument": DecodeExactDocument,
	}
	bodies := []struct {
		name, body    string
		category, not error
	}{
		{"duplicate", `{"name":"x","name":"y","child":{"value":1}}` + "\n", ErrDuplicateKey, ErrTrailingData},
		{"trailing value", testDocumentBody + " {}\n", ErrTrailingData, ErrDuplicateKey},
		{"trailing bytes", testDocumentBody + " nope\n", ErrTrailingData, ErrDuplicateKey},
	}
	for api, decode := range decoders {
		for _, tc := range bodies {
			t.Run(api+"/"+tc.name, func(t *testing.T) {
				err := decode([]byte(tc.body), new(testDocument))
				if !errors.Is(err, tc.category) || errors.Is(err, tc.not) {
					t.Fatalf("%s error = %v, want category %v and not %v", api, err, tc.category, tc.not)
				}
			})
		}
	}
}
