package gitguard

import (
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/shellcommand"
)

// The scanner walks the token stream command by command (splitting on control
// operators). It strips the honest-mistake prefixes an agent reflexively types
// (env/xargs/timeout/…), and hands each `git <subcommand>` to classify. The first shell
// wrapper child is scanned exactly one level deep, by design: this is an honest-mistake
// layer, not an evasion-resistant boundary. The child argv of `bench worktree exec` is
// the same one level, because a pool worktree runs git that way.

// keywords are shell keywords skipped in command position so the verb after them is
// found (`if git …`, `while git …`).
var keywords = map[string]bool{"if": true, "then": true, "elif": true, "else": true, "do": true, "while": true, "until": true, "!": true, "{": true}

// globalOptsWithArg are `git`'s pre-subcommand options that take a separate-word value,
// so find_subcommand skips the value too and does not mistake it for the subcommand.
var globalOptsWithArg = map[string]bool{"-C": true, "-c": true, "--exec-path": true, "--git-dir": true, "--namespace": true, "--work-tree": true}

// scan returns the deny label for the first destructive git command it finds, or "" if
// none. allowWrapper gates the one-level wrapper recursion.
func scan(stream shellcommand.Stream, chk Checker, allowWrapper bool) string {
	for _, span := range stream.Commands {
		tokens := shellcommand.ProjectCommandWords(stream.Tokens[span.Start:span.End])
		for len(tokens) > 0 && keywords[tokens[0]] {
			tokens = tokens[1:]
		}
		if reason := scanWords(tokens, chk, allowWrapper, false); reason != "" {
			return reason
		}
	}
	return ""
}

// scanWords classifies one simple command's words. The first shell wrapper child and a
// `bench worktree exec` child are each the one level of recursion allowWrapper gates.
// elsewhere marks words that run in another checkout, so the verdicts that read this
// directory's repository facts do not describe them.
func scanWords(tokens []string, chk Checker, allowWrapper, elsewhere bool) string {
	prefix := shellcommand.ResolveRoutinePrefix(tokens)
	j, viaXargs := prefix.Index, prefix.ViaXargs
	if !prefix.Executes || j >= len(tokens) {
		return ""
	}
	switch {
	case isGit(tokens[j]):
		sub, argsStart, ok := FindSubcommand(tokens, j+1, len(tokens))
		if !ok {
			return ""
		}
		redirected := elsewhere || redirectsRepository(tokens[j+1:argsStart-1])
		return classify(sub, tokens[argsStart:], viaXargs, redirected, chk)
	case !allowWrapper:
		return ""
	case shellcommand.IsBenchName(tokens[j]):
		if exec := shellcommand.WorktreeExec(tokens, j); exec.HasSeparator {
			return scanWords(exec.Child, chk, false, exec.DifferentCheckout)
		}
	default:
		if children := shellcommand.ShellChildren(tokens, j); len(children) > 0 {
			return scan(children[0].Stream, chk, false)
		}
	}
	return ""
}

// isGit reports whether a word names the git executable, bare or by path.
func isGit(word string) bool { return filepath.Base(word) == "git" }

// FindSubcommand skips git's global options (and their values) after the `git` token.
// It returns the subcommand, the index its args start at, and whether one was found.
func FindSubcommand(tokens []string, start, end int) (string, int, bool) {
	j := start
	for j < end {
		current := tokens[j]
		if current == "--" {
			j++
			break
		}
		if globalOptsWithArg[current] {
			j += 2
			continue
		}
		if longOptWithValue(current) {
			j++
			continue
		}
		if strings.HasPrefix(current, "-") {
			j++
			continue
		}
		break
	}
	if j < end {
		return tokens[j], j + 1, true
	}
	return "", 0, false
}

// repositoryRedirects are the global options that move git off the process working
// directory. The facts a verdict reads come from that directory, so a verdict that needs
// repository truth cannot grade a command carrying one of these.
var repositoryRedirects = []string{"-C", "--git-dir", "--work-tree"}

// redirectsRepository reports whether the global segment between `git` and its
// subcommand names another repository, in either the separate-word or the `=value`
// spelling.
func redirectsRepository(globals []string) bool {
	for _, token := range globals {
		for _, opt := range repositoryRedirects {
			if token == opt || strings.HasPrefix(token, opt+"=") {
				return true
			}
		}
	}
	return false
}

// longOptWithValue reports whether current is a `--opt=value` form of a global option
// that takes an argument (so it is one token, not two).
func longOptWithValue(current string) bool {
	for opt := range globalOptsWithArg {
		if strings.HasPrefix(opt, "--") && strings.HasPrefix(current, opt+"=") {
			return true
		}
	}
	return false
}
