package prose

import "strings"

// prepare splits one document into lines and blanks every span that is not prose: the
// frontmatter block, every HTML comment, and every fenced code block. It returns the fault of
// the first unterminated delimiter instead of the lines, because past that delimiter the
// parser cannot tell prose from code. Findings and Paragraphs both prepare here, so a step
// added to this list reaches both of them.
func prepare(doc string) ([]string, *Finding) {
	lines := strings.Split(doc, "\n")
	if f := stripFrontmatter(lines); f != nil {
		return nil, f
	}
	if f := stripComments(lines); f != nil {
		return nil, f
	}
	if f := stripFences(lines); f != nil {
		return nil, f
	}
	return lines, nil
}

// stripFrontmatter blanks a leading YAML block. Frontmatter values are trigger text for
// a loader, not prose for a reader.
func stripFrontmatter(lines []string) *Finding {
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return nil
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			for j := 0; j <= i; j++ {
				lines[j] = ""
			}
			return nil
		}
	}
	return &Finding{Kind: KindFrontmatter, Line: 1}
}

// stripComments removes every HTML comment span in place and keeps the line count, so
// each later finding still names the physical line of the document.
func stripComments(lines []string) *Finding {
	open, openLine := false, 0
	for i, line := range lines {
		var kept strings.Builder
		rest := line
		for {
			if open {
				end := strings.Index(rest, "-->")
				if end < 0 {
					rest = ""
					break
				}
				open, rest = false, rest[end+3:]
				continue
			}
			start := strings.Index(rest, "<!--")
			if start < 0 {
				kept.WriteString(rest)
				break
			}
			kept.WriteString(rest[:start])
			open, openLine, rest = true, i+1, rest[start+4:]
		}
		lines[i] = kept.String()
	}
	if open {
		return &Finding{Kind: KindComment, Line: openLine}
	}
	return nil
}

// stripFences blanks every fenced code block and both of its markers. A closing fence
// uses the same marker character and is at least as long as the opening one.
func stripFences(lines []string) *Finding {
	openLine, marker := 0, ""
	for i, line := range lines {
		found := fenceMarker(strings.TrimLeft(line, " \t"))
		if openLine == 0 {
			if found != "" {
				openLine, marker, lines[i] = i+1, found, ""
			}
			continue
		}
		lines[i] = ""
		if found != "" && found[0] == marker[0] && len(found) >= len(marker) {
			openLine, marker = 0, ""
		}
	}
	if openLine != 0 {
		return &Finding{Kind: KindFence, Line: openLine}
	}
	return nil
}

// fenceMarker returns the leading run of backticks or tildes when the run is a fence
// marker, and an empty string otherwise.
func fenceMarker(trimmed string) string {
	if len(trimmed) < 3 {
		return ""
	}
	c := trimmed[0]
	if c != '`' && c != '~' {
		return ""
	}
	n := 0
	for n < len(trimmed) && trimmed[n] == c {
		n++
	}
	if n < 3 {
		return ""
	}
	return trimmed[:n]
}
