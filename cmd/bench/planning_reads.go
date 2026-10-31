// The exact planning readers compose the page owner with each source owner.
package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

const (
	linesFlag    = "--lines"
	artifactFlag = "--artifact"
)

// readUsage is the one grammar of `bench read`. Its help, its root help row, and its parser
// read this line.
const readUsage = "bench read <path> [" + spec.SectionFlag + " <literal-H2-title> | " + linesFlag + " <first>:<last>] [" + artifactFlag + "] " + responsebound.CursorUsage

var readGrammar = usage.Grammar{
	Cmd:     usage.CommandName(readUsage),
	Help:    "usage: " + readUsage,
	MinArgs: 1,
	MaxArgs: 1,
	Flags: []usage.Flag{
		{Name: spec.SectionFlag, HasValue: true, NoEmptyValue: true},
		{Name: linesFlag, HasValue: true, NoEmptyValue: true},
		{Name: artifactFlag},
		{Name: responsebound.OffsetFlag, HasValue: true, NoEmptyValue: true},
		{Name: responsebound.SourceFlag, HasValue: true, NoEmptyValue: true},
	},
}

// usageRefusal renders a usage refusal of one parsed exact-reader call: the grammar's help
// line with the reason, at exit 2.
func usageRefusal(grammar usage.Grammar, reason string) (string, int) {
	return grammar.Help + " (" + reason + ")\n", 2
}

// repositoryPath answers the repository-relative path that operand names from the caller's
// directory, or a refusal with hint for a path outside the repository.
func repositoryPath(root, operand, hint string) (string, string, bool) {
	cwd, _ := os.Getwd()
	_, path := canonicalpath.Operand(root, cwd, operand)
	if !filepath.IsLocal(path) {
		return path, toon.Errorf(path+" is outside the repository", hint) + "\n", false
	}
	return path, "", true
}

// readCommand pages the exact bytes of one repository file, one physical-line range or one
// literal visible H2 section of it, or one absolute spill or exported-evidence artifact. A
// read grants the artifact no evidence authority: it verifies no manifest and qualifies
// nothing.
func readCommand(args []string) (string, int) {
	parsed, line, code := usage.Parse(readGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	offset, token, continued, err := responsebound.ParseCursor(parsed.Flags)
	if err != nil {
		return usageRefusal(readGrammar, err.Error())
	}
	_, artifact := parsed.Flags[artifactFlag]
	lines, ranged := parsed.Flags[linesFlag]
	section, sectioned := parsed.Flags[spec.SectionFlag]
	first, last := 0, 0
	switch {
	case sectioned && (ranged || artifact):
		return usageRefusal(readGrammar, spec.SectionFlag+" takes neither "+linesFlag+" nor "+artifactFlag)
	case artifact && ranged:
		return usageRefusal(readGrammar, linesFlag+" and "+artifactFlag+" are mutually exclusive")
	case artifact && !filepath.IsAbs(parsed.Positionals[0]):
		return usageRefusal(readGrammar, artifactFlag+" takes an absolute path")
	case ranged:
		if first, last, err = responsebound.ParseLineRange(lines); err != nil {
			return usageRefusal(readGrammar, err.Error())
		}
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	request := responsebound.ReadRequest{Kind: "file", Scope: root, Offset: offset, Token: token, Continued: continued}
	// The identity keeps the root-relative path, and every command the page offers keeps
	// the caller's own operand, so it names the same file from the caller's directory.
	path := filepath.Clean(parsed.Positionals[0])
	request.Command = append(strings.Fields(readGrammar.Cmd)[1:], path)
	if artifact {
		request.Kind, request.Scope, request.Path = "artifact", "", path
	} else {
		relative, refusal, local := repositoryPath(root, parsed.Positionals[0], "read a repository path, or pass "+artifactFlag+" with an absolute artifact path")
		if !local {
			return refusal, 1
		}
		request.Path = relative
		path = filepath.Join(root, request.Path)
	}
	source, err := responsebound.OpenSource(path, artifact)
	if err != nil {
		return readSourceRefusal(request.Path, err)
	}
	defer source.Close()
	request.Digest, request.SourceSize, request.Payload, request.Size = source.Digest, source.Size, source, source.Size
	switch {
	case ranged:
		selected, err := responsebound.SelectLines(source.Bytes(), first, last)
		if err != nil {
			return toon.Errorf("line range unavailable", err.Error()) + "\n", 1
		}
		request.Kind, request.Selection = "lines", lines
		request.Payload, request.Size = bytes.NewReader(selected), int64(len(selected))
		request.Command = append(request.Command, linesFlag, lines)
	case sectioned:
		selection, err := spec.Section(source.Bytes(), section)
		if err != nil {
			return toon.Errorf("section unavailable", err.Error()) + "\n", 1
		}
		return pageSelection(request, source, selection, spec.SectionFlag, section)
	case artifact:
		request.Command = append(request.Command, artifactFlag)
	}
	return pageSource(request, source)
}

// specShowCommand pages one selection of one spec: the whole file, one literal visible H2
// section, the User stories section, the accepted coverage map, the raw Status: line, or one
// story's typed coverage obligations. The spec owner declares the grammar and selects its
// spans, the coverage owner bounds the map and projects each story, and the page owner renders
// each page. Spec does not import coverage, so this adapter composes them.
func specShowCommand(args []string) (string, int) {
	parsed, line, code := usage.Parse(spec.ShowGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	offset, token, continued, err := responsebound.ParseCursor(parsed.Flags)
	if err != nil {
		return usageRefusal(spec.ShowGrammar, err.Error())
	}
	section, sectioned := parsed.Flags[spec.SectionFlag]
	story, storied := parsed.Flags[spec.StoryFlag]
	number := 0
	switch {
	case sectioned && storied:
		return usageRefusal(spec.ShowGrammar, spec.SectionFlag+" and "+spec.StoryFlag+" are mutually exclusive")
	case storied:
		if number, err = responsebound.ParsePositive(story); err != nil {
			return usageRefusal(spec.ShowGrammar, spec.StoryFlag+" takes "+err.Error())
		}
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	operand := parsed.Positionals[0]
	_, resolved, tried, ok, err := spec.Resolve(root, operand)
	if err != nil {
		return toon.Errorf("spec not readable: "+err.Error(), "fix the file or pass another spec") + "\n", 1
	}
	if !ok {
		return toon.Errorf("spec not found: "+strings.Join(tried, ", "), "pass a spec slug or a path to a spec markdown file") + "\n", 1
	}
	path, refusal, local := repositoryPath(root, resolved, "pass a spec slug or a path to a spec markdown file inside the repository")
	if !local {
		return refusal, 1
	}
	request := responsebound.ReadRequest{Scope: root, Path: path, Offset: offset, Token: token, Continued: continued}
	request.Command = append(strings.Fields(spec.ShowGrammar.Cmd)[1:], operand)
	source, err := responsebound.OpenSource(resolved, false)
	if err != nil {
		return readSourceRefusal(request.Path, err)
	}
	defer source.Close()
	request.Digest, request.SourceSize = source.Digest, source.Size
	if storied || section == spec.CoverageShortcut {
		return pageCoverage(request, source, section, number)
	}
	selection, err := spec.Show(source.Bytes(), section)
	if err != nil {
		return toon.Errorf("section unavailable", err.Error()) + "\n", 1
	}
	var selector []string
	if sectioned {
		selector = []string{spec.SectionFlag, section}
	}
	return pageSelection(request, source, selection, selector...)
}

// pageCoverage pages the coverage owner's selection of an open spec: the accepted coverage map
// for the coverage shortcut, and otherwise the versioned typed projection of the story
// numbered number. A malformed map and a repeated story number refuse before any content.
func pageCoverage(request responsebound.ReadRequest, source *responsebound.Source, section string, number int) (string, int) {
	document, err := coverage.Describe(source.Bytes())
	if err != nil {
		return toon.Errorf("coverage unavailable", err.Error()) + "\n", 1
	}
	if section == spec.CoverageShortcut {
		selection := spec.Selection{Kind: spec.KindCoverage, Selector: section, Present: document.MapPresent, Start: document.Map.Start, End: document.Map.End}
		return pageSelection(request, source, selection, spec.SectionFlag, section)
	}
	projection, err := document.Story(number)
	if err != nil {
		return toon.Errorf("story unavailable", err.Error()) + "\n", 1
	}
	var payload []byte
	if projection.Present {
		payload = projection.Encode()
	}
	selector := strconv.Itoa(number)
	request.Kind, request.Selection, request.Projection = spec.KindStory, selector, coverage.StorySchema
	return pagePayload(request, source, projection.Present, payload, spec.StoryFlag, selector)
}

// specCommand routes `bench spec show` to its exact reader, which composes the spec owner's
// selections with the page owner, and every other spec subcommand to the spec owner.
func specCommand(c Command, args []string) int {
	if len(args) > 0 && args[0] == spec.ShowSubcommand {
		return outputCommand(specShowCommand)(c, args[1:])
	}
	return retireListingCommand(spec.Command)(c, args)
}

// pageSelection pages one selection of an open source, and every command the page offers
// repeats the selector arguments.
func pageSelection(request responsebound.ReadRequest, source *responsebound.Source, selection spec.Selection, selector ...string) (string, int) {
	request.Kind, request.Selection = selection.Kind, selection.Selector
	return pagePayload(request, source, selection.Present, selection.Bytes(source.Bytes()), selector...)
}

// pagePayload pages the selected payload of an open source. The page owner answers an absent
// selection as one terminal page with present false and no content.
func pagePayload(request responsebound.ReadRequest, source *responsebound.Source, present bool, payload []byte, selector ...string) (string, int) {
	request.Absent = !present
	request.Payload, request.Size = bytes.NewReader(payload), int64(len(payload))
	request.Command = append(request.Command, selector...)
	return pageSource(request, source)
}

// pageSource renders the request's page, and refuses a source that changed during the read.
func pageSource(request responsebound.ReadRequest, source *responsebound.Source) (string, int) {
	page, err := responsebound.PageRead(request)
	if err != nil {
		return err.Error() + "\n", 1
	}
	if err := source.Unchanged(); err != nil {
		return readSourceRefusal(request.Path, err)
	}
	return page.Output, 0
}

// readSourceRefusal renders a file-state refusal as the shared record-error line.
func readSourceRefusal(path string, err error) (string, int) {
	var state *responsebound.SourceError
	if errors.As(err, &state) {
		return toon.RecordError(path, state.State, state.Reason) + "\n", 1
	}
	return toon.Errorf(path+" is unreadable", err.Error()) + "\n", 1
}
