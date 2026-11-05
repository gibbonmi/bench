package outline

import (
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
)

// TestOutlineHelpArguments is S41: the directory summary, the path-scoped listing, the
// declaration candidates, and the same-file line routes quote each known argument, and a
// future slot keeps its placeholder. Each command is authored by hand. A POSIX shell reads
// each command back, or only its known prefix when the command holds a future slot.
func TestOutlineHelpArguments(t *testing.T) {
	root := outlineRepo(t)
	writeOutlineFile(t, root, "my dir's/a.go", "package a\n\nfunc Get() {}\n\nfunc init() {}\n\nfunc init() {}\n")
	writeOutlineFile(t, root, "b/b.go", "package b\n\nfunc Get() {}\n")
	gitAddOutline(t, root)
	const dir = `'my dir'\''s'`
	for _, row := range []struct {
		name string
		args []string
		want []string
		argv [][]string
	}{
		{name: "directory summary", want: []string{"bench 'outline' 'b'", "bench 'outline' " + dir},
			argv: [][]string{{"bench", "outline", "b"}, {"bench", "outline", "my dir's"}}},
		{name: "path listing", args: []string{"my dir's"}, want: []string{"bench 'outline' " + dir + " '--symbol' <name>"},
			argv: [][]string{{"bench", "outline", "my dir's", "--symbol"}}},
		{name: "declaration candidates", args: []string{"--symbol", "Get"},
			want: []string{"bench 'outline' 'b/b.go' '--symbol' 'Get'", `bench 'outline' 'my dir'\''s/a.go' '--symbol' 'Get'`},
			argv: [][]string{{"bench", "outline", "b/b.go", "--symbol", "Get"}, {"bench", "outline", "my dir's/a.go", "--symbol", "Get"}}},
		{name: "same-file line routes", args: []string{"--symbol", "init"},
			want: []string{`bench 'read' 'my dir'\''s/a.go' '--lines' '5:5'`, `bench 'read' 'my dir'\''s/a.go' '--lines' '7:7'`},
			argv: [][]string{{"bench", "read", "my dir's/a.go", "--lines", "5:5"}, {"bench", "read", "my dir's/a.go", "--lines", "7:7"}}},
	} {
		t.Run(row.name, func(t *testing.T) {
			out, _ := Command(row.args)
			axitest.RequireHelpCommands(t, out, row.want, row.argv)
		})
	}
}
