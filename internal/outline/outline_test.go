package outline

import (
	"reflect"
	"testing"
)

// symbolsCase drives the pure parser seam: feed one buffer for a language and assert
// the exact (line,kind,name) rows, which is where the per-language pattern table's
// correctness lives.
type symbolsCase struct {
	name    string
	path    string
	content string
	want    []Symbol
}

func TestSymbols(t *testing.T) {
	cases := []symbolsCase{
		{
			name: "go func method and type",
			path: "x.go",
			content: "package x\n" +
				"\n" +
				"type Widget struct{}\n" +
				"\n" +
				"func New() *Widget { return nil }\n" +
				"\n" +
				"func (w *Widget) Name() string { return \"\" }\n",
			want: []Symbol{
				{Line: 3, Kind: "type", Name: "Widget"},
				{Line: 5, Kind: "func", Name: "New"},
				{Line: 7, Kind: "func", Name: "Name"},
			},
		},
		{
			name: "shell paren form and function keyword",
			path: "x.sh",
			content: "#!/usr/bin/env bash\n" +
				"route_porcelain() {\n" +
				"  :\n" +
				"}\n" +
				"function run_gate {\n" +
				"  :\n" +
				"}\n",
			want: []Symbol{
				{Line: 2, Kind: "function", Name: "route_porcelain"},
				{Line: 5, Kind: "function", Name: "run_gate"},
			},
		},
		{
			name: "markdown atx headings",
			path: "x.md",
			content: "# Title\n" +
				"\n" +
				"intro prose\n" +
				"\n" +
				"### A sub heading\n" +
				"not a #heading without space\n",
			want: []Symbol{
				{Line: 1, Kind: "heading", Name: "Title"},
				{Line: 5, Kind: "heading", Name: "A sub heading"},
			},
		},
		{
			name: "python def and class",
			path: "x.py",
			content: "import os\n" +
				"\n" +
				"class Widget:\n" +
				"    def render(self):\n" +
				"        return 1\n",
			want: []Symbol{
				{Line: 3, Kind: "class", Name: "Widget"},
				{Line: 4, Kind: "def", Name: "render"},
			},
		},
		{
			name: "js/ts function class and exported const arrow",
			path: "x.ts",
			content: "function plain() {}\n" +
				"class Widget {}\n" +
				"export const make = (a) => a + 1;\n",
			want: []Symbol{
				{Line: 1, Kind: "function", Name: "plain"},
				{Line: 2, Kind: "class", Name: "Widget"},
				{Line: 3, Kind: "const", Name: "make"},
			},
		},
		{
			name: "final line without trailing newline is still scanned",
			path: "x.go",
			content: "package x\n" +
				"func Tail() {}", // no trailing \n
			want: []Symbol{
				{Line: 2, Kind: "func", Name: "Tail"},
			},
		},
		{
			name:    "unknown extension yields no rows",
			path:    "x.rs",
			content: "fn main() {}\n",
			want:    nil,
		},
		{
			name:    "present-but-empty buffer yields no rows",
			path:    "x.go",
			content: "",
			want:    nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Symbols(c.path, []byte(c.content))
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Symbols(%q) = %#v\nwant %#v", c.path, got, c.want)
			}
		})
	}
}
