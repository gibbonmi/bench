package prose

import "github.com/gibbonmi/bench/internal/markdown"

// faultKinds names the finding that reports each fault of the block reader.
var faultKinds = map[markdown.FaultKind]FindingKind{
	markdown.FrontmatterFault: KindFrontmatter,
	markdown.CommentFault:     KindComment,
	markdown.FenceFault:       KindFence,
}

// prepare reads one document through the block reader and blanks every line that is not prose:
// the frontmatter block, whose values are trigger text for a loader, and every fenced block with
// both of its markers. Each other line keeps its comment-free text, and the line count stays, so
// each finding names the physical line of the document. prepare returns the reader's fault
// instead of the lines, because past an unterminated delimiter the parser cannot tell prose from
// code. Findings and Paragraphs both prepare here, so a step added here reaches both of them.
func prepare(doc string) ([]string, *Finding) {
	read := markdown.Read(doc)
	if read.Fault != nil {
		return nil, &Finding{Kind: faultKinds[read.Fault.Kind], Line: read.Fault.Line}
	}
	lines := make([]string, len(read.Lines))
	for i, line := range read.Lines {
		if line.Class == markdown.Body {
			lines[i] = line.Text
		}
	}
	return lines, nil
}
