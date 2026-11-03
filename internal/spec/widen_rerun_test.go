package spec

import (
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
)

// TestWidenRerunArguments: the rerun line spells the slug, the ticket, and every entry as
// one always-quoted shell argument each, and keeps the fixed flags and the <reference>
// placeholder bare. The exact line reds a bypassed owner or a former bare token. A POSIX
// shell then reads the line back, with r1 in place of the placeholder only in the recovery
// input because a shell reads a bare <reference> as a redirection, so a lost or split
// entry reds the recovered argv.
func TestWidenRerunArguments(t *testing.T) {
	entries := []string{"internal/x/w.go", "internal/other/z.go (new)", "internal/other/a b.go", "internal/other/it's.go"}
	root := retireWorktree(t, retirePrimary(t, "s t", "Status: staged\n", map[string]string{"specs/s t/tickets/a b.md": widenTicket("internal/x/y.go")}))
	out, code := runSpec(t, root, append([]string{"widen", "s t", "--ticket", "a b.md"}, entries...)...)
	_, rerun, found := strings.Cut(out, "rerun ")
	rerun, _, _ = strings.Cut(rerun, "\n")
	want := `bench spec widen 's t' --ticket 'a b.md' --decision <reference> 'internal/x/w.go' 'internal/other/z.go (new)' 'internal/other/a b.go' 'internal/other/it'\''s.go'`
	if code != 1 || !found || rerun != want {
		t.Fatalf("widen = (%d, %q), want exit 1 and the rerun %q", code, out, want)
	}
	argv, err := axitest.HelpAction{Cmd: strings.Replace(rerun, " <reference> ", " r1 ", 1)}.Argv()
	wantArgv := append([]string{"bench", "spec", "widen", "s t", "--ticket", "a b.md", "--decision", "r1"}, entries...)
	if err != nil || !slices.Equal(argv, wantArgv) {
		t.Fatalf("argv of %q = (%q, %v), want %q", rerun, argv, err, wantArgv)
	}
}
