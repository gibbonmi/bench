// Package commit owns the public command grammar and adapts it to exact landing.
package commit

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// Command runs a path-attributed prospective landing. Help exits 0, grammar errors exit
// 2, operational refusals exit 1, and a commit that published without reconciling its
// checkout exits 3; the landing owner alone composes, authorizes, and publishes the
// prospective tree.
func Command(args []string, stdout, stderr io.Writer) int {
	msg, paths, dryRun, help, usageErr := parseArgs(args)
	if help != "" {
		fmt.Fprintln(stdout, helpText)
		return 0
	}
	if usageErr != "" {
		fmt.Fprintln(stderr, grammar.Help+" ("+usageErr+")")
		return 2
	}
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return 1
	}
	primary, err := git.IsPrimaryCheckout(root)
	if err != nil {
		fmt.Fprintln(stderr, toon.Errorf("checkout identity is unknown", "repair Git metadata, then retry from a Bench worktree"))
		return 1
	}
	if primary {
		fmt.Fprintln(stderr, usage.PrimaryCheckoutRefusal())
		return 1
	}

	// Capture publication identity before reading attributed content. A detached checkout
	// updates literal HEAD; an attached checkout updates its full branch ref.
	destination := "HEAD"
	if out, symbolicErr := git.Raw("-C", root, "symbolic-ref", "-q", "HEAD"); symbolicErr == nil {
		destination = strings.TrimSpace(string(out))
	}
	expectedBytes, expectedErr := git.Raw("-C", root, "rev-parse", "--verify", "HEAD^{commit}")
	if expectedErr != nil {
		fmt.Fprintln(stderr, "error: destination has no commit base")
		return 1
	}

	named := make([]string, 0, len(paths))
	for _, path := range paths {
		rel, relErr := rootRel(root, path)
		if relErr != nil {
			fmt.Fprintf(stderr, "error: cannot resolve path %q relative to repo root: %v\n", path, relErr)
			return 1
		}
		named = append(named, rel)
	}
	if dryRun {
		if err := landing.New().DryRun(context.Background(), landing.Request{
			Root: root, Destination: destination, Expected: strings.TrimSpace(string(expectedBytes)),
			Message: msg, Paths: named, Stdout: stdout, Stderr: stderr,
		}); err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return 1
		}
		fmt.Fprintf(stdout, "dry run: composed %d path(s) authorized green; nothing committed\n", len(named))
		return 0
	}
	if _, err := landing.New().Land(context.Background(), landing.Request{
		Root: root, Destination: destination, Expected: strings.TrimSpace(string(expectedBytes)),
		Message: msg, Paths: named, Stdout: stdout, Stderr: stderr,
	}); err != nil {
		var remainder *landing.PublishedUnreconciledError
		if errors.As(err, &remainder) {
			return publicationRemainder(stdout, remainder)
		}
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "committed %d path(s)\n", len(named))
	return 0
}

// publicationRemainder reports the publication boundary the landing owner reached: the
// commit exists and the checkout does not match it. The record uses the landing verb's
// name{key=value,...} grammar, and its exit code separates this outcome from a refusal
// that published nothing.
func publicationRemainder(stdout io.Writer, remainder *landing.PublishedUnreconciledError) int {
	fmt.Fprintf(stdout, "committed{published_commit=%s,path=%s,next=%s}\n",
		remainder.Commit, sanitize.Controls(remainder.Path), restoreNext(remainder.Commit, remainder.Paths))
	return 3
}

// restoreNext names the one restore that reconciles every named path against the
// published commit. The restore is idempotent, so it covers the paths that already
// reconciled as well as the remainder. A path that is not line-safe takes the landing
// verb's pointer form: quoting would still emit the raw byte into a line-structured
// record, and escaping would name a path that does not exist.
//
// The value is line-safe by construction, so it reaches the record unescaped; the
// sanitizer's backslash escaping would break the quoting a reader pastes.
func restoreNext(commit string, paths []string) string {
	command := "git restore --source=" + commit + " --staged --worktree --"
	quoted := make([]string, 0, len(paths))
	for _, path := range paths {
		if !sanitize.LineSafe(path) {
			return command + " <named-paths>"
		}
		quoted = append(quoted, sanitize.ShellQuote(path))
	}
	return command + " " + strings.Join(quoted, " ")
}

// helpText adds one concrete example and the exit-code meanings the grammar line
// cannot carry. A usage error prints the grammar line alone, so only a help request
// pays for them. The example shows the trailing `-- <path>...` form, so a caller does
// not learn the argument shape by tripping the usage line.
var helpText = grammar.Help + "\n" +
	"example: bench commit -m \"fix: tighten the guard\" -- internal/gitguard/scan.go docs/adr/0007.md\n" +
	"--dry-run: gate the exact composed snapshot and report the verdict; commit nothing\n" +
	"exit 1: refused before publication; nothing was committed\n" +
	"exit 2: grammar error\n" +
	"exit 3: published; the checkout did not reconcile — paste next= to repair"

var grammar = usage.Grammar{
	Cmd:     "bench commit",
	Help:    "usage: bench commit [--dry-run] -m <msg> [--] <path>...",
	Flags:   []usage.Flag{{Name: "-m", HasValue: true}, {Name: "--dry-run"}},
	MaxArgs: -1,
}

func parseArgs(args []string) (msg string, paths []string, dryRun bool, help string, usageErr string) {
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		if code == 0 {
			return "", nil, false, line, ""
		}
		return "", nil, false, "", line
	}
	_, dryRun = parsed.Flags["--dry-run"]
	msg, msgSet := parsed.Flags["-m"]
	if !msgSet {
		return "", nil, false, "", "-m <msg> is required"
	}
	if strings.TrimSpace(msg) == "" {
		return "", nil, false, "", "-m <msg> must not be empty"
	}
	if len(parsed.Positionals) == 0 {
		return "", nil, false, "", "at least one <path> is required"
	}
	return msg, parsed.Positionals, dryRun, "", ""
}

func rootRel(root, arg string) (string, error) {
	abs, err := filepath.Abs(arg)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(rel), nil
}
