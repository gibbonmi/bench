package gitguard

import "testing"

// A heredoc body is data the shell hands to a file, not command text. Writing a file
// whose contents mention a destructive git command must not be refused; the same verb
// outside the body must still be. Same class FT229 fixed for the degraded rim.
func TestClassifyTreatsHeredocBodiesAsData(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{"quoted delimiter", "cat > f <<'EOF'\ngit push --force\nEOF", ""},
		{"double-quoted delimiter", "cat > f <<\"EOF\"\ngit reset --hard\nEOF", ""},
		{"unquoted delimiter", "cat > f <<EOF\ngit push origin main\nEOF", ""},
		{"dash form with tab-indented delimiter", "cat > f <<-EOF\n\tgit push\n\tEOF", ""},
		{"two bodies on one command", "cat > a <<'A'\ngit push\nA\ncat > b <<'B'\ngit clean -fd\nB", ""},
		{"verb after the delimiter line", "cat > f <<'EOF'\nx\nEOF\ngit push", "git push"},
		{"verb on the operator's own line", "git push <<'EOF'\nx\nEOF", "git push"},
		{"herestring keeps its classification", "git push <<< x", "git push"},
		{"unterminated body stays data", "cat > f <<'EOF'\ngit push", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Classify(c.in, refYes); got != c.want {
				t.Errorf("Classify(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
