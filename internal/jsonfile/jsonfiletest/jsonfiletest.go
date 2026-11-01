// Package jsonfiletest supplies the hostile document family and the committed baseline
// comparison that the strict JSON compatibility tests share.
package jsonfiletest

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/jsonfile"
)

// UnknownKey is the field that the family adds where the schema declares no such field.
const UnknownKey = "extra"

// Seed is one valid object document and the members that Family edits.
type Seed struct {
	// Valid is a valid object document that the caller accepts.
	Valid string
	// Member is one root member of Valid. Wrong is Member with a value of the wrong type.
	Member, Wrong string
	// Nested is one member of an object inside Valid. An empty Nested omits the nested cases.
	Nested string
}

// Family derives the hostile document family from seed. Each name is a stable baseline key.
// It panics when Valid lacks a seed member, because that case would silently equal Valid.
func Family(seed Seed) map[string]string {
	valid, member, wrong := seed.Valid, seed.Member, seed.Wrong
	for _, part := range []string{member, seed.Nested} {
		if !strings.Contains(valid, part) {
			panic(fmt.Sprintf("jsonfiletest: seed member %q is not in %q", part, valid))
		}
	}
	extra := `"` + UnknownKey + `":true`
	unknown := strings.Replace(valid, "{", "{"+extra+",", 1)
	duplicate := strings.Replace(valid, member, member+","+member, 1)
	mistyped := strings.Replace(valid, member, wrong, 1)
	family := map[string]string{
		"framed":                    valid + "\n",
		"unframed":                  valid,
		"surrounding whitespace":    " \n" + valid + " \t\r\n",
		"whitespace only":           " \n\t",
		"empty":                     "",
		"null":                      "null",
		"unknown field":             unknown,
		"duplicate key":             duplicate,
		"syntax error":              strings.TrimSuffix(valid, "}") + ",}",
		"wrong type":                mistyped,
		"second value":              valid + " {}",
		"framed second value":       valid + "\n{}\n",
		"malformed suffix":          valid + " nope",
		"unknown plus suffix":       unknown + " {}",
		"duplicate plus suffix":     duplicate + " {}",
		"unknown plus duplicate":    strings.Replace(duplicate, "{", "{"+extra+",", 1),
		"wrong type plus suffix":    mistyped + " {}",
		"wrong type plus duplicate": strings.Replace(mistyped, wrong, wrong+","+member, 1),
	}
	if seed.Nested != "" {
		family["nested duplicate key"] = strings.Replace(valid, seed.Nested, seed.Nested+","+seed.Nested, 1)
		family["nested unknown field"] = strings.Replace(valid, seed.Nested, seed.Nested+","+extra, 1)
	}
	return family
}

// ReadBaseline strictly decodes the committed baseline at path into stored. It never
// writes the file.
func ReadBaseline(t testing.TB, path string, stored any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := jsonfile.DecodeDocument(data, stored); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
}

// Compare reports each case whose current result differs from its baseline result. It also
// reports each case that only one side holds.
func Compare[T any](t testing.TB, current, stored map[string]T) {
	t.Helper()
	names := make([]string, 0, len(current)+len(stored))
	for name := range current {
		names = append(names, name)
	}
	for name := range stored {
		if _, ok := current[name]; !ok {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	for _, name := range names {
		got, ok := current[name]
		want, recorded := stored[name]
		switch {
		case !ok:
			t.Errorf("%q: the baseline case has no current result", name)
		case !recorded || !reflect.DeepEqual(got, want):
			t.Errorf("%q: current %s, baseline %s", name, render(got), render(want))
		}
	}
}

// render prints a result as JSON, so the two sides of a mismatch read alike.
func render(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return err.Error()
	}
	return string(data)
}
