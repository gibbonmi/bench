package usage

import (
	"strings"

	"github.com/gibbonmi/bench/internal/toon"
)

// Flag declares one flag a Grammar accepts: its literal spelling (e.g. "-m",
// "--all") and whether it consumes the following argument as its value.
type Flag struct {
	Name     string
	HasValue bool
}

// Grammar declares the argument shape one subcommand parses through Parse:
// its usage-line command name, its declared help text, the flags it accepts,
// and its positional arity. MaxArgs of -1 means unbounded (variadic).
type Grammar struct {
	Cmd     string
	Help    string
	Flags   []Flag
	MinArgs int
	MaxArgs int
}

// Result is a successful parse: the flags present (an empty string value for
// a boolean flag) and the positionals in argv order.
type Result struct {
	Flags       map[string]string
	Positionals []string
}

// Parse applies g's grammar to args and returns exactly one of three
// outcomes: a populated Result with an empty line and code 0 on a successful
// parse; a help outcome (zero Result, g.Help as the line, code 0) when --help
// or -h appears before "--", or when bare help is the sole argument; or a
// rendered usage error (zero Result, a toon.Usage or toon.MissingArg line,
// code 2). A caller prints a non-empty line and exits with code verbatim;
// Result is meaningful only when line == "".
//
// Rendering is delegated to toon.Usage and toon.MissingArg throughout, so
// this package never re-derives the usage-line shape.
func Parse(g Grammar, args []string) (Result, string, int) {
	// Bare "help" is a request only alone: a variadic grammar's free text or path
	// list can legitimately contain the word, and recognizing it anywhere would
	// print usage and silently discard the rest of the invocation. The
	// flag-spelled forms below are unambiguous, and "--" already escapes them.
	if len(args) == 1 && args[0] == "help" {
		return Result{}, g.Help, 0
	}

	valueFlags := make(map[string]bool, len(g.Flags))
	knownFlags := make(map[string]bool, len(g.Flags))
	for _, f := range g.Flags {
		knownFlags[f.Name] = true
		if f.HasValue {
			valueFlags[f.Name] = true
		}
	}

	result := Result{Flags: map[string]string{}}
	// endedFlags becomes permanent at the first "--"; every later "--" or
	// dash-prefixed token is then an ordinary positional rather than being
	// re-examined as a separator or a flag.
	endedFlags := false
	for i := 0; i < len(args); i++ {
		a := args[i]
		if !endedFlags {
			if a == "--help" || a == "-h" {
				return Result{}, g.Help, 0
			}
			if a == "--" {
				endedFlags = true
				continue
			}
			// A bare "-" is a legal filename, not a flag: only a longer
			// dash-prefixed token is flag-like.
			if a != "-" && strings.HasPrefix(a, "-") {
				if !knownFlags[a] {
					return Result{}, toon.Usage(g.Cmd, a), 2
				}
				// A declared flag given twice is a usage error, not last-one-wins:
				// silently keeping the later value hides a mistyped invocation whose
				// two spellings disagree. The check precedes the value read, so a
				// repeated value flag names the flag instead of consuming another
				// argument, and one rule covers value and boolean flags alike.
				if _, repeated := result.Flags[a]; repeated {
					return Result{}, toon.Usage(g.Cmd, a), 2
				}
				if valueFlags[a] {
					if i+1 >= len(args) {
						return Result{}, toon.MissingArg(g.Cmd, a), 2
					}
					i++
					result.Flags[a] = args[i]
				} else {
					result.Flags[a] = ""
				}
				continue
			}
		}
		// An empty positional names nothing — it is what an unset shell variable
		// expands to inside quotes — and a subcommand that resolves it against the
		// filesystem silently widens to the cwd. Rejecting it here gives every
		// grammar the guard instead of each path-taking subcommand re-deriving it.
		if a == "" {
			return Result{}, toon.Usage(g.Cmd, `""`), 2
		}
		// Trailing garbage is reported on the first excess argument, not a
		// generic message, so a mistyped invocation names the token that
		// broke it.
		if g.MaxArgs >= 0 && len(result.Positionals) >= g.MaxArgs {
			return Result{}, toon.Usage(g.Cmd, a), 2
		}
		result.Positionals = append(result.Positionals, a)
	}

	if len(result.Positionals) < g.MinArgs {
		return Result{}, toon.MissingArg(g.Cmd, "argument"), 2
	}
	return result, "", 0
}
