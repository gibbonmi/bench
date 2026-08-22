package gitguard

import (
	"path/filepath"
	"regexp"
	"strings"
)

// The scanner walks the token stream command by command (splitting on control
// operators). It strips the honest-mistake prefixes an agent reflexively types
// (env/xargs/timeout/…), and hands each `git <subcommand>` to classify. Wrapper strings
// (`sh|bash|zsh -c '…'`) are re-tokenized and scanned exactly one level deep, by design:
// this is an honest-mistake layer, not an evasion-resistant boundary.

// controlOps end a command; after one, the next word is a fresh command position.
var controlOps = map[string]bool{";": true, "&&": true, "||": true, "|": true, "&": true, "(": true, ")": true}

// keywords are shell keywords skipped in command position so the verb after them is
// found (`if git …`, `while git …`).
var keywords = map[string]bool{"if": true, "then": true, "elif": true, "else": true, "do": true, "while": true, "until": true, "!": true, "{": true}

// wrappers whose `-c` string is re-scanned one level deep.
var wrappers = map[string]bool{"sh": true, "bash": true, "zsh": true}

// globalOptsWithArg are `git`'s pre-subcommand options that take a separate-word value,
// so find_subcommand skips the value too and does not mistake it for the subcommand.
var globalOptsWithArg = map[string]bool{"-C": true, "-c": true, "--exec-path": true, "--git-dir": true, "--namespace": true, "--work-tree": true}

var envAssignRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*=`)
var wrapperCFlagRe = regexp.MustCompile(`^-[A-Za-z]*c[A-Za-z]*$`)

// scan returns the deny label for the first destructive git command it finds, or "" if
// none. allowWrapper gates the one-level wrapper recursion.
func scan(tokens []string, chk Checker, allowWrapper bool) string {
	expectCommand := true
	i, n := 0, len(tokens)
	for i < n {
		tok := tokens[i]
		if controlOps[tok] {
			expectCommand = true
			i++
			continue
		}
		if expectCommand && keywords[tok] {
			i++
			continue
		}
		if expectCommand {
			end := commandEnd(tokens, i)
			j, viaXargs := resolvePrefixes(tokens, i, end)
			if j < end {
				base := filepath.Base(tokens[j])
				if base == "git" {
					sub, argsStart, ok := findSubcommand(tokens, j+1, end)
					if ok {
						if reason := classify(sub, tokens[argsStart:end], viaXargs, chk); reason != "" {
							return reason
						}
					}
				} else if allowWrapper && wrappers[base] {
					for k := j + 1; k < end; k++ {
						if wrapperCFlagRe.MatchString(tokens[k]) {
							if k+1 < end {
								inner := tokenize(tokens[k+1])
								if reason := scan(inner, chk, false); reason != "" {
									return reason
								}
							}
							break
						}
					}
				}
			}
			expectCommand = false
			i = end
			continue
		}
		i++
	}
	return ""
}

// commandEnd returns the index of the next control op at or after i (the end of the
// current simple command).
func commandEnd(tokens []string, i int) int {
	j := i
	for j < len(tokens) && !controlOps[tokens[j]] {
		j++
	}
	return j
}

// resolvePrefixes advances past leading env assignments and the honest-mistake command
// wrappers (env/command/nohup/timeout/xargs). It returns the index of the real verb and
// whether an xargs prefix was seen. xargs feeds paths from stdin, so a pathspec-less
// checkout/restore under it is treated as destructive.
func resolvePrefixes(tokens []string, i, end int) (int, bool) {
	viaXargs := false
	for i < end {
		if envAssignRe.MatchString(tokens[i]) {
			i++
			continue
		}
		base := filepath.Base(tokens[i])
		switch base {
		case "env":
			i++
			for i < end && envAssignRe.MatchString(tokens[i]) {
				i++
			}
			continue
		case "command", "nohup":
			i++
			continue
		case "timeout":
			i++
			for i < end && strings.HasPrefix(tokens[i], "-") {
				i++
			}
			if i < end {
				i++
			}
			continue
		case "xargs":
			i++
			viaXargs = true
			for i < end && strings.HasPrefix(tokens[i], "-") {
				i++
			}
			continue
		}
		break
	}
	return i, viaXargs
}

// findSubcommand skips git's global options (and their values) after the `git` token.
// It returns the subcommand, the index its args start at, and whether one was found.
func findSubcommand(tokens []string, start, end int) (string, int, bool) {
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
