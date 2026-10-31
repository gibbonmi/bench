// The exact planning readers compose the page owner with each source owner.
package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/canonicalpath"
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

// readUsageRefusal renders a usage refusal of a parsed `bench read` call.
func readUsageRefusal(reason string) (string, int) {
	return readGrammar.Help + " (" + reason + ")\n", 2
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
		return readUsageRefusal(err.Error())
	}
	_, artifact := parsed.Flags[artifactFlag]
	lines, ranged := parsed.Flags[linesFlag]
	section, sectioned := parsed.Flags[spec.SectionFlag]
	first, last := 0, 0
	switch {
	case sectioned && (ranged || artifact):
		return readUsageRefusal(spec.SectionFlag + " takes neither " + linesFlag + " nor " + artifactFlag)
	case artifact && ranged:
		return readUsageRefusal(linesFlag + " and " + artifactFlag + " are mutually exclusive")
	case artifact && !filepath.IsAbs(parsed.Positionals[0]):
		return readUsageRefusal(artifactFlag + " takes an absolute path")
	case ranged:
		if first, last, err = responsebound.ParseLineRange(lines); err != nil {
			return readUsageRefusal(err.Error())
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
		cwd, _ := os.Getwd()
		_, request.Path = canonicalpath.Operand(root, cwd, parsed.Positionals[0])
		if !filepath.IsLocal(request.Path) {
			return toon.Errorf(request.Path+" is outside the repository", "read a repository path, or pass "+artifactFlag+" with an absolute artifact path") + "\n", 1
		}
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
// section, the User stories section, or the raw Status: line. The spec owner declares the
// grammar and selects the bytes, and the page owner renders each page.
func specShowCommand(args []string) (string, int) {
	parsed, line, code := usage.Parse(spec.ShowGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	offset, token, continued, err := responsebound.ParseCursor(parsed.Flags)
	if err != nil {
		return spec.ShowGrammar.Help + " (" + err.Error() + ")\n", 2
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
	request := responsebound.ReadRequest{Scope: root, Path: spec.RelTo(root, resolved), Offset: offset, Token: token, Continued: continued}
	request.Command = append(strings.Fields(spec.ShowGrammar.Cmd)[1:], operand)
	source, err := responsebound.OpenSource(resolved, false)
	if err != nil {
		return readSourceRefusal(request.Path, err)
	}
	defer source.Close()
	request.Digest, request.SourceSize = source.Digest, source.Size
	section, sectioned := parsed.Flags[spec.SectionFlag]
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

// specCommand routes `bench spec show` to its exact reader, which composes the spec owner's
// selections with the page owner, and every other spec subcommand to the spec owner.
func specCommand(c Command, args []string) int {
	if len(args) > 0 && args[0] == spec.ShowSubcommand {
		return outputCommand(specShowCommand)(c, args[1:])
	}
	return retireListingCommand(spec.Command)(c, args)
}

// absenceFields are the columns of an absent selection, which no page represents.
var absenceFields = []string{"kind", "selector", "present"}

// pageSelection pages one present selection of an open source, and every command the page
// offers repeats the selector arguments. A first read of an absent selection answers its
// explicit absence at exit 0. A continuation of an absent selection pages the empty
// selection, so the page owner refuses its cursor, because no page of that source offered it.
func pageSelection(request responsebound.ReadRequest, source *responsebound.Source, selection spec.Selection, selector ...string) (string, int) {
	if !selection.Present && !request.Continued {
		if err := source.Unchanged(); err != nil {
			return readSourceRefusal(request.Path, err)
		}
		absence, err := toon.TableTyped("selection", absenceFields, [][]any{{selection.Kind, selection.Selector, false}})
		if err != nil {
			return toon.RenderError(err) + "\n", 1
		}
		help, _ := axi.RenderHelp(nil)
		return absence + help, 0
	}
	selected := selection.Bytes(source.Bytes())
	request.Kind, request.Selection = selection.Kind, selection.Selector
	request.Payload, request.Size = bytes.NewReader(selected), int64(len(selected))
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
