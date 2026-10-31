package spec

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/gibbonmi/bench/internal/markdown"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/usage"
)

// ShowSubcommand names the show child. SectionFlag selects one literal visible H2 title or
// one shortcut, and `bench read --section` takes the same flag. StoriesShortcut selects the User stories
// section, and StatusShortcut selects the raw Status: line that the metadata owner chooses.
const (
	ShowSubcommand  = "show"
	SectionFlag     = "--section"
	StoriesShortcut = "stories"
	StatusShortcut  = "status"
	// StoriesSection is the title of the H2 section that holds a spec's user stories.
	StoriesSection = "User stories"
)

// ShowUsage is the one grammar of `bench spec show`. Its help, its root help row, and its
// parser read this line.
const ShowUsage = "bench spec " + ShowSubcommand + " <slug-or-path> [" + SectionFlag + " " + StoriesShortcut + "|" + StatusShortcut + "|<literal-H2-title>] " + responsebound.CursorUsage

// ShowGrammar parses `bench spec show`.
var ShowGrammar = usage.Grammar{
	Cmd:     usage.CommandName(ShowUsage),
	Help:    "usage: " + ShowUsage,
	MinArgs: 1,
	MaxArgs: 1,
	Flags: []usage.Flag{
		{Name: SectionFlag, HasValue: true, NoEmptyValue: true},
		{Name: responsebound.OffsetFlag, HasValue: true, NoEmptyValue: true},
		{Name: responsebound.SourceFlag, HasValue: true, NoEmptyValue: true},
	},
}

// The selection kinds of the spec owner.
const (
	KindFile    = "file"
	KindSection = "section"
	KindStories = "stories"
	KindStatus  = "status"
)

// Selection is one exact span of source bytes, from Start up to but not including End. An
// absent selection has Present false and an empty span.
type Selection struct {
	Kind, Selector string
	Present        bool
	Start, End     int
}

// Bytes answers the selected source bytes.
func (s Selection) Bytes(content []byte) []byte { return content[s.Start:s.End] }

// ErrSelection is the sentinel that each selection refusal wraps. A refusal comes before any
// content.
var ErrSelection = errors.New("selection refused")

// Show selects the spec bytes that one show selector names: the whole file for an empty
// selector, the User stories section for the stories shortcut, the metadata owner's raw
// Status: line for the status shortcut, and otherwise one literal visible H2 section.
func Show(content []byte, selector string) (Selection, error) {
	switch selector {
	case "":
		return Selection{Kind: KindFile, Present: true, End: len(content)}, nil
	case StoriesShortcut:
		selection, err := Section(content, StoriesSection)
		selection.Kind, selection.Selector = KindStories, selector
		return selection, err
	case StatusShortcut:
		selection := Selection{Kind: KindStatus, Selector: selector}
		if line, _ := metadataLines(content); line != nil {
			selection.Present, selection.Start = true, line.Offset
			selection.End = physicalLineEnd(content, line.Offset)
		}
		return selection, nil
	}
	return Section(content, selector)
}

// Section selects the visible H2 section whose title is title: its owning heading line
// through the byte before the next visible H2 heading, or through the end of the content.
// The Markdown owner decides each visible heading, so a heading inside a fence or an HTML
// comment owns nothing. A Markdown fault or a repeated matching heading refuses; no matching
// heading answers an absent selection.
func Section(content []byte, title string) (Selection, error) {
	selection := Selection{Kind: KindSection, Selector: title}
	doc := markdown.Read(string(content))
	if doc.Fault != nil {
		return Selection{}, fmt.Errorf("%w: the Markdown block that opens on line %d is unterminated", ErrSelection, doc.Fault.Line)
	}
	matches, open := 0, false
	for _, line := range doc.Lines {
		if !line.Heading {
			continue
		}
		if open {
			selection.End, open = line.Offset, false
		}
		if line.Section == title {
			matches++
			selection.Start, selection.End, open = line.Offset, len(content), true
		}
	}
	if matches > 1 {
		return Selection{}, fmt.Errorf("%w: %d visible H2 headings are titled %q", ErrSelection, matches, title)
	}
	selection.Present = matches == 1
	return selection, nil
}

// physicalLineEnd answers the offset after the physical line that starts at start, its
// newline included.
func physicalLineEnd(content []byte, start int) int {
	if end := bytes.IndexByte(content[start:], '\n'); end >= 0 {
		return start + end + 1
	}
	return len(content)
}
