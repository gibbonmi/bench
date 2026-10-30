// The exact planning readers compose the page owner with each source owner.
package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

const (
	linesFlag    = "--lines"
	artifactFlag = "--artifact"
)

// readUsage is the one grammar of `bench read`. Its help, its root help row, and its parser
// read this line.
const readUsage = "bench read <path> [" + linesFlag + " <first>:<last>] [" + artifactFlag + "] " + responsebound.CursorUsage

var readGrammar = usage.Grammar{
	Cmd:     usage.CommandName(readUsage),
	Help:    "usage: " + readUsage,
	MinArgs: 1,
	MaxArgs: 1,
	Flags: []usage.Flag{
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

// readCommand pages the exact bytes of one repository file, one physical-line range of it, or
// one absolute spill or exported-evidence artifact. A read grants the artifact no evidence
// authority: it verifies no manifest and qualifies nothing.
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
	first, last := 0, 0
	switch {
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
	request := responsebound.ReadRequest{Kind: "file", Scope: root, Offset: offset, Token: token, Continued: continued, Frame: commandDefinition{Scope: scopeTree}.treeRow(nil)}
	path := filepath.Clean(parsed.Positionals[0])
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
	request.Command = []string{"read", request.Path}
	request.Digest, request.SourceSize, request.Payload, request.Size = source.Digest, source.Size, source, source.Size
	if ranged {
		selected, err := responsebound.SelectLines(source.Bytes(), first, last)
		if err != nil {
			return toon.Errorf("line range unavailable", err.Error()) + "\n", 1
		}
		request.Kind, request.Selection = "lines", lines
		request.Payload, request.Size = bytes.NewReader(selected), int64(len(selected))
		request.Command = append(request.Command, linesFlag, lines)
	}
	if artifact {
		request.Command = append(request.Command, artifactFlag)
	}
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
