package gitguard

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gibbonmi/bench/internal/shellcommand"
)

// The scanner walks the token stream command by command (splitting on control
// operators). It strips the honest-mistake prefixes an agent reflexively types
// (env/xargs/timeout/…), and hands each `git <subcommand>` to classify. Wrapper strings
// (`sh|bash|zsh -c '…'`) are re-tokenized and scanned exactly one level deep, by design:
// this is an honest-mistake layer, not an evasion-resistant boundary.

// keywords are shell keywords skipped in command position so the verb after them is
// found (`if git …`, `while git …`).
var keywords = map[string]bool{"if": true, "then": true, "elif": true, "else": true, "do": true, "while": true, "until": true, "!": true, "{": true}

// wrappers whose `-c` string is re-scanned one level deep.
var wrappers = map[string]bool{"sh": true, "bash": true, "zsh": true}

// globalOptsWithArg are `git`'s pre-subcommand options that take a separate-word value,
// so find_subcommand skips the value too and does not mistake it for the subcommand.
var globalOptsWithArg = map[string]bool{"-C": true, "-c": true, "--exec-path": true, "--git-dir": true, "--namespace": true, "--work-tree": true}

var wrapperCFlagRe = regexp.MustCompile(`^-[A-Za-z]*c[A-Za-z]*$`)

// scan returns the deny label for the first destructive git command it finds, or "" if
// none. allowWrapper gates the one-level wrapper recursion.
func scan(stream shellcommand.Stream, chk Checker, allowWrapper bool) string {
	for _, span := range stream.Commands {
		tokens := shellcommand.ProjectCommandWords(stream.Tokens[span.Start:span.End])
		if len(tokens) == 0 {
			continue
		}
		for len(tokens) > 0 && keywords[tokens[0]] {
			tokens = tokens[1:]
		}
		prefix := shellcommand.ResolveRoutinePrefix(tokens)
		j, viaXargs := prefix.Index, prefix.ViaXargs
		if prefix.Executes && j < len(tokens) {
			base := filepath.Base(tokens[j])
			if base == "git" {
				sub, argsStart, ok := FindSubcommand(tokens, j+1, len(tokens))
				if ok {
					if reason := classify(sub, tokens[argsStart:], viaXargs, chk); reason != "" {
						return reason
					}
				}
			} else if allowWrapper && wrappers[base] {
				for k := j + 1; k < len(tokens); k++ {
					if wrapperCFlagRe.MatchString(tokens[k]) {
						if k+1 < len(tokens) {
							inner := shellcommand.Parse(tokens[k+1])
							if reason := scan(inner, chk, false); reason != "" {
								return reason
							}
						}
						break
					}
				}
			}
		}
	}
	return ""
}

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
