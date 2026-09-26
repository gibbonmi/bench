// File patch partition of the immutable explicit-pair response for package diff.
package diff

import (
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// FilePatch is one verbatim file patch of a frozen diff body and the repository path that
// identifies it: the tip path of a surviving file, or the base path of a deleted one.
type FilePatch struct {
	Path string
	Body []byte
}

// PatchSnapshot is one explicit-pair full response split at its file patch boundaries.
// Prefix, every patch body in order, and Suffix rebuild the response byte for byte.
type PatchSnapshot struct {
	Prefix, Suffix []byte
	Patches        []FilePatch
}

// PairPatches renders the response that `bench diff --base <commit> --source-tip <commit>
// --full` prints for args, and partitions its body into ordered file patches. It reads the
// same snapshot the command renders and keeps Git's ambient body bytes. A body that does not
// partition into identified patches consistent with the frozen inventory refuses as a whole,
// so a caller never holds a partial patch set.
func PairPatches(args []string) (PatchSnapshot, string, int) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return PatchSnapshot{}, line + "\n", code
	}
	for _, flag := range []string{"--base", "--source-tip", "--full"} {
		if _, ok := parsed.Flags[flag]; !ok {
			return PatchSnapshot{}, toon.MissingArg(grammar.Cmd, flag) + "\n", 2
		}
	}
	if _, hasCommit := parsed.Flags["--commit"]; hasCommit {
		return PatchSnapshot{}, toon.Usage(grammar.Cmd, "--commit") + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return PatchSnapshot{}, toon.NotInRepo() + "\n", 1
	}
	dr, kind, hint := resolvePairRange(root, parsed.Flags["--base"], parsed.Flags["--source-tip"])
	if kind != "" {
		return PatchSnapshot{}, toon.Errorf(kind, hint) + "\n", 1
	}
	response, kind, hint := renderCommitResponse(root, dr, true)
	if kind != "" {
		return PatchSnapshot{}, toon.Errorf(kind, hint) + "\n", 1
	}
	patches, err := partitionPatches(response.body, response.inventory)
	if err != nil {
		return PatchSnapshot{}, toon.Errorf("diff body not partitionable", err.Error()) + "\n", 1
	}
	return PatchSnapshot{Prefix: []byte(response.prefix), Suffix: []byte(response.suffix), Patches: patches}, "", 0
}

// extendedHeaders are the Git extended header lines that may follow a `diff --git` line
// before the patch content.
var extendedHeaders = []string{
	"old mode ", "new mode ", "deleted file mode ", "new file mode ", "copy from ", "copy to ",
	"rename from ", "rename to ", "similarity index ", "dissimilarity index ", "index ",
}

// patchHeader is the path facts one patch header region states.
type patchHeader struct {
	gitBase, gitTip   string
	gitPaths, deleted bool
	sideBase, sideTip string
	from, to          string
}

// identity applies the file identity rule: the base path of a deleted file, else the tip
// path. Explicit rename or copy paths come first, then the `diff --git` paths, then the side
// headers, so a patch without `---` and `+++` lines still names its file.
func (h patchHeader) identity() (path, renamedFrom string) {
	if h.deleted {
		return firstPath(h.from, h.gitPath(h.gitBase), h.sideBase), ""
	}
	return firstPath(h.to, h.gitPath(h.gitTip), h.sideTip), h.from
}

func (h patchHeader) gitPath(path string) string {
	if h.gitPaths {
		return path
	}
	return ""
}

func firstPath(paths ...string) string {
	for _, path := range paths {
		if path != "" {
			return path
		}
	}
	return ""
}

// partitionPatches splits a Git patch body at its `diff --git` lines. It reads path headers
// only in each header region, before any content, so content that resembles a header stays
// content. Every patch path must belong to the inventory, and every inventory path must have
// a patch or a rename that names it.
func partitionPatches(body []byte, inventory []string) ([]FilePatch, error) {
	if len(body) > 0 && body[len(body)-1] != '\n' {
		return nil, fmt.Errorf("the body does not end with a line feed")
	}
	lines := strings.SplitAfter(string(body), "\n")
	lines = lines[:len(lines)-1]
	known, covered := map[string]bool{}, map[string]bool{}
	for _, path := range inventory {
		known[path] = true
	}
	var patches []FilePatch
	for i := 0; i < len(lines); {
		start := i
		fields, ok := strings.CutPrefix(lines[i], "diff --git ")
		if !ok {
			return nil, fmt.Errorf("body line %d does not start a patch", i+1)
		}
		var header patchHeader
		header.gitBase, header.gitTip, header.gitPaths = git.PatchGitPaths(strings.TrimSuffix(fields, "\n"))
		for i++; i < len(lines) && readExtendedHeader(&header, lines[i]); i++ {
		}
		if i < len(lines) && strings.HasPrefix(lines[i], "--- ") {
			if i+1 >= len(lines) || !strings.HasPrefix(lines[i+1], "+++ ") {
				return nil, fmt.Errorf("body line %d opens a side header without its tip side", i+1)
			}
			header.sideBase, header.sideTip = git.PatchSidePath(lines[i][4:]), git.PatchSidePath(lines[i+1][4:])
			i += 2
		} else if i < len(lines) && strings.HasPrefix(lines[i], "Binary files ") {
			i++
		}
		if i < len(lines) && !strings.HasPrefix(lines[i], "diff --git ") && !strings.HasPrefix(lines[i], "@@ ") {
			return nil, fmt.Errorf("body line %d is neither a hunk nor a patch", i+1)
		}
		for ; i < len(lines) && !strings.HasPrefix(lines[i], "diff --git "); i++ {
			if !strings.ContainsRune("@ +-\\\n", rune(lines[i][0])) {
				return nil, fmt.Errorf("body line %d is not patch content", i+1)
			}
		}
		path, renamedFrom := header.identity()
		if !known[path] {
			return nil, fmt.Errorf("patch %d names %q, which the frozen inventory does not hold", len(patches)+1, path)
		}
		covered[path], covered[renamedFrom] = true, true
		patches = append(patches, FilePatch{Path: path, Body: []byte(strings.Join(lines[start:i], ""))})
	}
	for _, path := range inventory {
		if !covered[path] {
			return nil, fmt.Errorf("inventory path %q has no patch", path)
		}
	}
	return patches, nil
}

// readExtendedHeader records one extended header line and reports whether it was one.
func readExtendedHeader(header *patchHeader, line string) bool {
	for _, prefix := range extendedHeaders {
		value, ok := strings.CutPrefix(line, prefix)
		if !ok {
			continue
		}
		value = git.UnquotePath(strings.TrimSuffix(value, "\n"))
		switch prefix {
		case "deleted file mode ":
			header.deleted = true
		case "rename from ", "copy from ":
			header.from = value
		case "rename to ", "copy to ":
			header.to = value
		}
		return true
	}
	return false
}
