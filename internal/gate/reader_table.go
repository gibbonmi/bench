package gate

// The reader-table derivation: one pure function from the Go test logs of a whole kit run
// to the test packages that read each covered family, or a refusal.

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

// testLogHeader is the first line of every Go test log.
const testLogHeader = "# test log"

// readerSources are the inputs of the derivation. Packages are the test packages, as
// slash-separated directories relative to Root. Logs holds the log of each package that left
// one. Links maps each tracked symbolic link of the graded tree to its link text.
type readerSources struct {
	Root     string
	Packages []string
	Logs     map[string][]byte
	Links    map[string]string
	Families []string
	Children []childReader
}

// readerTable maps each covered family to the sorted test packages that read it.
type readerTable map[string][]string

// deriveReaderTable answers the reader table of in, or the refusal that names its cause. It
// refuses a listed package with no log or with an unreadable log, and a table in which the
// root-conformance package reads less than every family.
func deriveReaderTable(in readerSources) (readerTable, error) {
	table := readerTable{}
	for _, pkg := range in.Packages {
		data, logged := in.Logs[pkg]
		if !logged {
			return nil, fmt.Errorf("reader table refused: %s: no test log", pkg)
		}
		families, err := loggedFamilies(in, pkg, data)
		if err != nil {
			return nil, fmt.Errorf("reader table refused: %s: %w", pkg, err)
		}
		for _, family := range families {
			table.add(family, pkg)
		}
	}
	for _, family := range in.Families {
		if !slices.Contains(table[family], conformancePackagePath) {
			return nil, fmt.Errorf("reader table refused: %s reads no %s path", conformancePackagePath, family)
		}
	}
	for _, child := range in.Children {
		table.add(child.Family, child.Package)
	}
	return table, nil
}

// add records pkg as a reader of family once, in sorted order.
func (t readerTable) add(family, pkg string) {
	readers := t[family]
	if i, found := slices.BinarySearch(readers, pkg); !found {
		t[family] = slices.Insert(readers, i, pkg)
	}
}

// loggedFamilies answers the covered families that the log data of pkg reads. A relative
// path resolves against the package directory, which each chdir line moves. A path outside
// the root or inside the package's own directory adds nothing, and a tracked link resolves to
// its target before the family match. The root package's directory is the whole root, so the
// own-directory rule does not apply to it.
func loggedFamilies(in readerSources, pkg string, data []byte) ([]string, error) {
	if len(data) == 0 {
		return nil, errors.New("empty test log")
	}
	lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	if lines[0] != testLogHeader {
		return nil, fmt.Errorf("malformed test log: first line %q", lines[0])
	}
	dir := filepath.Join(in.Root, filepath.FromSlash(pkg))
	var families []string
	for _, line := range lines[1:] {
		op, name, _ := strings.Cut(line, " ")
		if name == "" {
			return nil, fmt.Errorf("malformed test log line %q", line)
		}
		switch op {
		case "getenv":
			continue
		case "chdir":
			dir = logPath(dir, name)
			continue
		case "open", "stat":
		default:
			return nil, fmt.Errorf("unknown test log operation %q", op)
		}
		rel, inside := rootRelative(in.Root, logPath(dir, name))
		if !inside || rel == pkg || strings.HasPrefix(rel, pkg+"/") {
			continue
		}
		if family, ok := familyOf(in.Families, linkTarget(in.Links, rel)); ok && !slices.Contains(families, family) {
			families = append(families, family)
		}
	}
	return families, nil
}

// logPath answers the logged name as an absolute path: a relative name resolves against dir.
func logPath(dir, name string) string {
	if filepath.IsAbs(name) {
		return filepath.Clean(name)
	}
	return filepath.Join(dir, name)
}

// rootRelative answers the slash-separated path of p relative to root, and whether p lies
// inside root.
func rootRelative(root, p string) (string, bool) {
	rel, err := filepath.Rel(root, p)
	if err != nil {
		return "", false
	}
	rel = filepath.ToSlash(rel)
	return rel, rel != ".." && !strings.HasPrefix(rel, "../")
}

// linkTarget resolves each tracked link at a prefix of rel to its target, so a read through
// a link counts as a read of the target's family. A target outside the root keeps its `../`
// prefix, which no family claims. The bound ends a cycle of links.
func linkTarget(links map[string]string, rel string) string {
	for range len(links) + 1 {
		resolved := false
		for link, target := range links {
			rest, below := strings.CutPrefix(rel, link)
			if below && (rest == "" || strings.HasPrefix(rest, "/")) {
				if path.IsAbs(target) {
					return target
				}
				rel, resolved = path.Join(path.Dir(link), target)+rest, true
				break
			}
		}
		if !resolved {
			return rel
		}
	}
	return rel
}
