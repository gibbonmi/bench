package env

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// allowResult is the parsed contents of a .bench/env.allow file. Patterns are
// stored in their validated on-disk form (an exact name, or a PREFIX* glob).
// This lets them be appended directly to the default passlist and matched with
// the same matchesAny helper.
type allowResult struct {
	agent []string
}

// namePattern is the portable environment-name set: an ASCII letter or
// underscore, then any run of ASCII letters, digits, or underscores. It is
// applied to an exact entry whole, and to a glob entry's prefix (the part
// before its trailing *).
var namePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// parseAllowFile reads and parses path. A missing file is not an error — it
// means defaults only. A present-but-unreadable file, or one that fails
// parseAllow's grammar, is fail-closed: an error naming the problem, never a
// silent fall-back to defaults.
func parseAllowFile(path string) (allowResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return allowResult{}, nil
		}
		return allowResult{}, fmt.Errorf(".bench/env.allow: %w", err)
	}
	return parseAllow(string(data))
}

// utf8BOM is the UTF-8 byte-order mark. A file that opens with it is rejected by
// name. This avoids the generic "entry before any section header" error the BOM
// bytes would otherwise trigger. The BOM is not Unicode whitespace, so TrimSpace
// leaves it attached to the first line. Both paths fail closed, but a named
// reason beats a misleading one.
const utf8BOM = "\ufeff"

// parseAllow implements the env.allow grammar. The grammar is optional and
// line-oriented. It allows # comments and blank lines, and an [agent] section
// header (the only known section). One entry appears per line, either an
// exact name or a PREFIX* glob.
//
// Any violation is rejected with an error naming the 1-indexed line and the
// reason. The first three reasons: an entry before any section header, an
// unknown section name, and a bare *. A stale [gate] section counts as
// unknown, since the gate opt-in lives in the manifest, not this file.
//
// Two more reasons: a glob that is not a single trailing *, or an entry
// with a / or =. Another is a character outside the portable
// environment-name set.
//
// A last reason is a leading UTF-8 byte-order mark. A present-but-empty
// file yields an allowResult with no entries, which is not an error.
func parseAllow(data string) (allowResult, error) {
	if strings.HasPrefix(data, utf8BOM) {
		return allowResult{}, errors.New(".bench/env.allow:1: file begins with a UTF-8 byte-order mark (BOM); save it as UTF-8 without a BOM")
	}
	var result allowResult
	section := ""
	for i, raw := range strings.Split(data, "\n") {
		lineNum := i + 1
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") {
			name, ok := parseSectionHeader(line)
			if !ok {
				return allowResult{}, fmt.Errorf(".bench/env.allow:%d: unknown section name %q", lineNum, line)
			}
			section = name
			continue
		}
		if section == "" {
			return allowResult{}, fmt.Errorf(".bench/env.allow:%d: entry %q before any section header", lineNum, line)
		}
		entry, err := validateEntry(line)
		if err != nil {
			return allowResult{}, fmt.Errorf(".bench/env.allow:%d: %w", lineNum, err)
		}
		if section == "agent" {
			result.agent = append(result.agent, entry)
		}
	}
	return result, nil
}

// parseSectionHeader reports the section name for a line that starts with "[",
// and false if it is not exactly "[agent]" — the only known section. A stale
// "[gate]" is rejected here like any other unknown section: the gate opt-in
// lives in the manifest, not this file.
func parseSectionHeader(line string) (string, bool) {
	if !strings.HasSuffix(line, "]") {
		return "", false
	}
	name := line[1 : len(line)-1]
	if name == "agent" {
		return name, true
	}
	return "", false
}

// validateEntry checks one non-comment, non-header line against the grammar.
// It returns the line unchanged (it is already in the PREFIX* or exact-name
// form Build's matcher expects), or a reason it was rejected.
func validateEntry(raw string) (string, error) {
	if raw == "*" {
		return "", errors.New("bare wildcard entry")
	}
	if strings.Contains(raw, "/") {
		return "", fmt.Errorf("entry %q contains '/'", raw)
	}
	if strings.Contains(raw, "=") {
		return "", fmt.Errorf("entry %q contains '='", raw)
	}
	if strings.Contains(raw, "*") {
		if !strings.HasSuffix(raw, "*") || strings.Count(raw, "*") > 1 {
			return "", fmt.Errorf("entry %q is not a single trailing glob", raw)
		}
		prefix := raw[:len(raw)-1]
		if !namePattern.MatchString(prefix) {
			return "", fmt.Errorf("entry %q contains a character outside the portable environment-name set", raw)
		}
		return raw, nil
	}
	if !namePattern.MatchString(raw) {
		return "", fmt.Errorf("entry %q contains a character outside the portable environment-name set", raw)
	}
	return raw, nil
}
