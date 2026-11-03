package commit

import (
	"strings"

	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/usage"
)

// LaneClause names what a commit grades: the declared lane, or the gate when the
// project declares no lane. The `bench help` row and the --dry-run help line read it.
const LaneClause = "run the declared lane (or the gate when no lane is declared)"

// helpText adds one concrete example, the deletion route, and the exit-code meanings
// the grammar line cannot carry. A usage error prints the grammar line alone, so only a
// help request pays for them. The example shows the trailing `-- <path>...` form, so a
// caller does not learn the argument shape by tripping the usage line. The command
// layer that runs the chain adds the --preflight-build line.
var helpText = grammar.Help + "\n" +
	"example: bench commit -m \"fix: tighten the guard\" -- internal/gitguard/scan.go docs/adr/0007.md\n" +
	"deleted path: name a file or folder deleted from the worktree and the commit publishes its deletion; no git rm is needed\n" +
	"--dry-run: " + LaneClause + " on the exact composed snapshot and report the outcome; commit nothing\n" +
	"exit 1: refused before publication; nothing was committed\n" +
	"exit 2: grammar error\n" +
	"exit 3: published; the checkout did not reconcile — paste " + refusalroute.NextField + "= to repair"

var grammar = usage.Grammar{
	Cmd:  "bench commit",
	Help: "usage: bench commit [--dry-run | " + PreflightBuildGroup + "] -m <msg> [--] <path>...",
	Flags: []usage.Flag{{Name: "-m", HasValue: true}, {Name: "--dry-run"},
		{Name: PreflightBuildFlag, HasValue: true, NoEmptyValue: true},
		{Name: BaseFlag, HasValue: true, NoEmptyValue: true}},
	MaxArgs: -1,
}

// PreflightBuildFlag names the spec slug whose build the command layer grades at the
// candidate commit before publication.
const PreflightBuildFlag = "--preflight-build"

// BaseFlag names the base that the command layer passes to the chained build preflight.
const BaseFlag = "--base"

// The placeholder names of the chain values. The grammar line and a composed chain read them.
const (
	slugOperand = "slug"
	baseOperand = "commit"
)

// PreflightBuildGroup is the argument shape of the chain: the slug and the optional base.
// The usage line, the `bench help` row, and the chain help line read it.
const PreflightBuildGroup = PreflightBuildFlag + " <" + slugOperand + "> [" + BaseFlag + " <" + baseOperand + ">]"

// ChainArgs composes the chain flags of a commit call: the slug, and the base when it is
// set. Each value is rendered by refusalroute.Arg under the placeholder name of the grammar line.
func ChainArgs(slug, base string) []string {
	args := []string{PreflightBuildFlag, refusalroute.Arg(slugOperand, slug)}
	if base != "" {
		args = append(args, BaseFlag, refusalroute.Arg(baseOperand, base))
	}
	return args
}

// ChainValues answers the slug and the base of a commit call that parses. A call that does
// not parse answers two empty values.
func ChainValues(args []string) (slug, base string) {
	req, _, _ := parseRequest(args)
	return req.preflightBuild, req.base
}

// HelpRowSuffix is the argument shape of the `bench help` row.
const HelpRowSuffix = " -m <msg> [" + PreflightBuildGroup + "] <path>..."

// request is one parsed commit call.
type request struct {
	msg            string
	paths          []string
	dryRun         bool
	preflightBuild string
	base           string
}

func parseRequest(args []string) (req request, help string, usageErr string) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		if code == 0 {
			return request{}, line, ""
		}
		return request{}, "", line
	}
	_, req.dryRun = parsed.Flags["--dry-run"]
	req.preflightBuild = parsed.Flags[PreflightBuildFlag]
	// A dry run writes no candidate commit, so the chain would have no commit to grade.
	if req.dryRun && req.preflightBuild != "" {
		return request{}, "", "--dry-run excludes " + PreflightBuildFlag
	}
	// The base belongs to the chained build preflight, so without the chain nothing reads it.
	base, baseSet := parsed.Flags[BaseFlag]
	if baseSet && req.preflightBuild == "" {
		return request{}, "", BaseFlag + " requires " + PreflightBuildFlag
	}
	req.base = base
	msg, msgSet := parsed.Flags["-m"]
	if !msgSet {
		return request{}, "", "-m <msg> is required"
	}
	if strings.TrimSpace(msg) == "" {
		return request{}, "", "-m <msg> must not be empty"
	}
	if len(parsed.Positionals) == 0 {
		return request{}, "", "at least one <path> is required"
	}
	req.msg, req.paths = msg, parsed.Positionals
	return req, "", ""
}
