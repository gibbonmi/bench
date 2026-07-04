// Package diff ports `bench diff`: the single source of review-base truth for the
// review phase. The recorded `branch.<name>.benchBase` key wins when it names a
// reachable ancestor; otherwise merge-base with the default branch. The preamble
// names which method resolved so a review agent can see a fallback happen, and a
// recorded sha that is unreachable or not an ancestor falls back loudly rather than
// silently diffing from the wrong base. Changed files follow as a
// `files[N]{status,path}:` table with paths escaped exactly once.
package diff

import (
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

// parseNameStatusZ turns `git diff --name-status --no-renames -z` output into
// status/path rows. The NUL framing carries each path raw (never git C-quoted), so a
// path with a space, a non-ASCII byte, or a quote arrives intact and is TOON-escaped a
// single layer downstream.
func parseNameStatusZ(raw []byte) [][]string {
	parts := strings.Split(string(raw), "\x00")
	var rows [][]string
	for i := 0; i+1 < len(parts); i += 2 {
		st, p := parts[i], parts[i+1]
		if st == "" && p == "" {
			continue
		}
		rows = append(rows, []string{st, p})
	}
	return rows
}

func changedFiles(base string) [][]string {
	raw, err := git.Raw("diff", "--name-status", "--no-renames", "-z", base+"...HEAD")
	if err != nil {
		return nil
	}
	return parseNameStatusZ(raw)
}

// Command implements `bench diff`.
func Command(args []string) (string, int) {
	switch {
	case len(args) == 0:
	case args[0] == "-h" || args[0] == "--help":
		return "usage: bench diff\n", 0
	default:
		return toon.Usage("bench diff", args[0]) + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}

	base, method := resolveBase()
	if base == "" {
		def := git.DefaultBranch(root)
		mb, err := git.Output("merge-base", def, "HEAD")
		if err != nil {
			return toon.Errorf("cannot resolve a review base",
				"no merge-base with '"+def+"'; record one with: git config branch.<name>.benchBase <sha>") + "\n", 1
		}
		base = mb
		if method == "" {
			method = "merge-base"
		}
	}

	branch, _ := git.Output("symbolic-ref", "--quiet", "--short", "HEAD")
	branchLabel := branch
	if branchLabel == "" {
		branchLabel = "(detached)"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "branch: %s\n", branchLabel)
	fmt.Fprintf(&b, "base: %s\n", base)
	fmt.Fprintf(&b, "method: %s\n", method)
	tbl, err := toon.Table("files", []string{"status", "path"}, changedFiles(base))
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	b.WriteString(tbl)
	return b.String(), 0
}

// resolveBase returns the recorded-key base and `recorded` when the key names a
// reachable ancestor, or ("", <loud fallback method>) when the key is present but
// unreachable/divergent, or ("","") when there is no usable recorded key.
func resolveBase() (base, method string) {
	branch, _ := git.Output("symbolic-ref", "--quiet", "--short", "HEAD")
	if branch == "" {
		return "", ""
	}
	key, _ := git.Output("config", "branch."+branch+".benchBase")
	if key == "" {
		return "", ""
	}
	switch {
	case !git.OK("cat-file", "-e", key+"^{commit}"):
		return "", "merge-base (recorded sha unreachable)"
	case !git.OK("merge-base", "--is-ancestor", key, "HEAD"):
		return "", "merge-base (recorded sha not an ancestor)"
	default:
		return key, "recorded"
	}
}
