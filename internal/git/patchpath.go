package git

import (
	"strconv"
	"strings"
)

// This file is the one decoder of patch path syntax. It is pure: it reads header text and
// runs no process, so every caller that reads a patch header shares one spelling rule.

// UnquotePath undoes the C-style quoting git applies to a patch header path. Git always
// escapes a double quote, a backslash, and a control character, whatever core.quotePath
// says, so such a path arrives wrapped in double quotes with C escapes inside; a control
// byte comes as a three-digit octal escape. A field that is not quoted, or whose escapes do
// not parse, is returned as it came: an unmatched path drops its own row, which is a better
// answer than a guessed spelling.
func UnquotePath(field string) string {
	if len(field) < 2 || field[0] != '"' || field[len(field)-1] != '"' {
		return field
	}
	body := field[1 : len(field)-1]
	var out strings.Builder
	for i := 0; i < len(body); i++ {
		if body[i] != '\\' {
			out.WriteByte(body[i])
			continue
		}
		i++
		if i >= len(body) {
			return field
		}
		switch c := body[i]; c {
		case '\\', '"':
			out.WriteByte(c)
		case 'a':
			out.WriteByte('\a')
		case 'b':
			out.WriteByte('\b')
		case 'f':
			out.WriteByte('\f')
		case 'n':
			out.WriteByte('\n')
		case 'r':
			out.WriteByte('\r')
		case 't':
			out.WriteByte('\t')
		case 'v':
			out.WriteByte('\v')
		default:
			if i+3 > len(body) {
				return field
			}
			octal, err := strconv.ParseUint(body[i:i+3], 8, 8)
			if err != nil {
				return field
			}
			out.WriteByte(byte(octal))
			i += 2
		}
	}
	return out.String()
}

// PatchSidePath decodes the path of a `---` or `+++` header field and strips the a/ or b/
// display prefix git prints. A missing side reads as /dev/null and keeps the empty string, so
// an added or deleted file has one nameless side. The quoting wraps the prefix too, so it is
// undone first. Git ends the field with a tab when the path holds a space, so only the line
// ending and that one tab are removed: a space at the end of an unquoted path is path text.
func PatchSidePath(field string) string {
	field = UnquotePath(strings.TrimSuffix(strings.TrimSuffix(field, "\n"), "\t"))
	if field == "/dev/null" {
		return ""
	}
	return stripDisplayPrefix(field)
}

// PatchGitPaths decodes the base and tip paths of a `diff --git` header, given the text after
// that marker, and strips each display prefix. A quoted field ends at its closing quote. Two
// unquoted fields can hold spaces, so the header splits only where both halves name one path;
// otherwise ok is false and the caller reads the explicit rename, copy, or side headers.
func PatchGitPaths(fields string) (base, tip string, ok bool) {
	first, second, ok := splitGitFields(fields)
	if !ok {
		return "", "", false
	}
	return stripDisplayPrefix(UnquotePath(first)), stripDisplayPrefix(UnquotePath(second)), true
}

func splitGitFields(fields string) (first, second string, ok bool) {
	if strings.HasPrefix(fields, `"`) {
		end := quotedEnd(fields)
		if end < 0 || end+1 >= len(fields) || fields[end+1] != ' ' {
			return "", "", false
		}
		return fields[:end+1], fields[end+2:], true
	}
	if strings.HasSuffix(fields, `"`) {
		for i := strings.Index(fields, ` "`); i >= 0; {
			if quotedEnd(fields[i+1:]) == len(fields)-i-2 {
				return fields[:i], fields[i+1:], true
			}
			next := strings.Index(fields[i+1:], ` "`)
			if next < 0 {
				break
			}
			i += next + 1
		}
		return "", "", false
	}
	if len(fields)%2 == 0 || fields[len(fields)/2] != ' ' {
		return "", "", false
	}
	first, second = fields[:len(fields)/2], fields[len(fields)/2+1:]
	if stripDisplayPrefix(first) != stripDisplayPrefix(second) {
		return "", "", false
	}
	return first, second, true
}

// quotedEnd returns the index of the quote that closes the quoted field at the start of s,
// or -1 when the field never closes.
func quotedEnd(s string) int {
	for i := 1; i < len(s); i++ {
		switch s[i] {
		case '\\':
			i++
		case '"':
			return i
		}
	}
	return -1
}

func stripDisplayPrefix(path string) string {
	if len(path) > 2 && (path[:2] == "a/" || path[:2] == "b/") {
		return path[2:]
	}
	return path
}
