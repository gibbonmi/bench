package recordcmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/probe"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/shellcommand"
	"github.com/gibbonmi/bench/internal/usage"
)

var excerptForm = form{name: reviewrecord.FormExcerpt, description: "stage the retained response of this worktree as the excerpt of one requirement of the chunk plan at HEAD",
	flags:  []flag{{"--chunk", "<id>", true, once}, {"--requirement", "<id>", true, once}, {mutantFlag, "<n>", false, once}},
	layout: []string{"--chunk", "--requirement", "[", mutantFlag, "]"},
	valid:  mutantRow, run: excerpt}

// mutantRow closes the --mutant value set of a form whose one grammar rule is that row.
func mutantRow(_ form, parsed usage.Result) bool { return mutantValid(parsed) }

// retainedSource names the retained response in a probe import refusal.
const retainedSource = "the retained response"

// excerpt stages the retained response for one requirement of the chunk plan at HEAD. The
// retained call must be the planned command, or `bench probe` for a requirement with a probe,
// and stageRetained must accept it. It writes no review record and no worktree file.
func excerpt(f form, root, spec string, parsed usage.Result) (string, int) {
	values := parsed.Flags
	head, err := git.ResolveCommit(root, "HEAD", "--quiet", "--end-of-options")
	if err != nil {
		return f.refuse("HEAD names no commit")
	}
	_, plan, err := reviewrecord.ReadSpecAtCommit(root, spec, head)
	if err != nil {
		return f.refuse(err.Error())
	}
	requirement, err := reviewrecord.ChunkRequirement(plan, head, values["--chunk"], values["--requirement"])
	if err != nil {
		return f.refuse(err.Error())
	}
	planned, ok := plannedCall(requirement.Command)
	if !ok {
		return f.refuse(fmt.Sprintf("the planned command %q of requirement %s is not one lexed bench command of words alone", requirement.Command, requirement.ID))
	}
	home := benchhome.Dir()
	retained, err := responsebound.ReadRetained(home, root)
	if errors.Is(err, fs.ErrNotExist) {
		return f.refuse("this worktree has no retained response; run " + requirement.Command + " first")
	}
	if err != nil {
		return f.refuse(err.Error())
	}
	call := append([]string{retained.Verb}, retained.Args...)
	if requirement.Probe != "" {
		if retained.Verb != probe.Verb {
			return f.refuse(fmt.Sprintf("requirement %s plans a probe; the retained call is bench %s, not bench %s", requirement.ID, retained.Verb, probe.Verb))
		}
	} else if !slices.Equal(call, planned) {
		return f.refuse(fmt.Sprintf("the retained call bench %s differs from the planned command %s", strings.Join(call, " "), requirement.Command))
	}
	mutant, _ := integer(values, mutantFlag)
	staged, out := stageRetained(requirement, retained, mutant)
	if out != "" {
		return f.refuse(out)
	}
	if err := responsebound.Stage(home, root, parsed.Positionals[0], values["--chunk"], requirement.ID, staged); err != nil {
		return f.refuse(err.Error())
	}
	return row("excerpt", []string{"chunk", "requirement", "verb", "exit", "head", "bytes"},
		[]string{values["--chunk"], requirement.ID, retained.Verb, strconv.Itoa(retained.Exit), retained.Head, strconv.Itoa(len(retained.Output))})
}

// stageRetained answers the stage of retained for requirement, with mutant as the row of a
// mutant-set answer, or the refusal. The excerpt form and the ticket form apply its checks:
// --mutant needs a probe requirement, the bytes pass the record excerpt bound, and the bytes of
// a probe requirement import as a probe answer.
func stageRetained(requirement reviewrecord.Requirement, retained responsebound.Retained, mutant int) (responsebound.Staged, string) {
	if mutant != 0 && requirement.Probe == "" {
		return responsebound.Staged{}, fmt.Sprintf("%s needs a probe requirement; requirement %s plans no probe", mutantFlag, requirement.ID)
	}
	if graded := bounds.ClassifyBytes(bounds.Read(bytes.NewReader(retained.Output), bounds.ControlRecordLimit)); graded.State != bounds.StateParsed {
		return responsebound.Staged{}, fmt.Sprintf("the retained response is %s: %s", graded.State, graded.Reason)
	}
	if requirement.Probe != "" {
		if _, out := readImported(retained.Output, mutant, retainedSource); out != "" {
			return responsebound.Staged{}, out
		}
	}
	return responsebound.Staged{Retained: retained, Mutant: mutant}, ""
}

// plannedCall answers the words after `bench` of a planned command. Only a lexed parse of
// word tokens alone whose first word is `bench` answers true, so a pipe, a redirection, a
// second command, or a foreign executable never matches a retained call.
func plannedCall(command string) ([]string, bool) {
	stream := shellcommand.Parse(command)
	if stream.Unlexed || len(stream.Tokens) < 2 || stream.Tokens[0].Text != "bench" {
		return nil, false
	}
	words := make([]string, 0, len(stream.Tokens)-1)
	for _, token := range stream.Tokens {
		if token.Kind != shellcommand.Word {
			return nil, false
		}
		words = append(words, token.Text)
	}
	return words[1:], true
}
