package gitguard

import "testing"

// TestScanWrapperDepthAndGlobalOpts pins two scan properties the verdict table does not.
// Wrapper recursion stops at exactly one level: a wrapper nested inside a wrapper string
// is not re-expanded (the honest-mistake threat model). And git's global options and
// their values are skipped before the subcommand.
func TestScanWrapperDepthAndGlobalOpts(t *testing.T) {
	cases := []struct {
		name string
		cmd  string
		want string
	}{
		{"one-level wrapper scanned", `bash -c 'git push'`, "git push"},
		{"nested wrapper not re-expanded", `bash -c 'sh -c "git push"'`, ""},
		{"global -C with value then verb", "git -C /tmp reset --hard", "git reset --hard"},
		{"global --git-dir= form then verb", "git --git-dir=/x push", "git push"},
		{"global opts, benign verb allowed", "git -C . status --short", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Classify(c.cmd, refYes); got != c.want {
				t.Errorf("Classify(%q) = %q, want %q", c.cmd, got, c.want)
			}
		})
	}
}
