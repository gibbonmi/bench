package usage

import (
	"strings"

	"github.com/gibbonmi/bench/internal/toon"
)

// Flag declares one flag a Grammar accepts: its literal spelling (e.g. "-m",
// "--all"), whether it consumes the following argument as its value, and
// whether that value may be empty.
//
// NoEmptyValue applies the empty-positional rule to a flag's value: an empty
// string is what an unset shell variable expands to inside quotes, so a flag
// whose value names something — a command, a path, a message — treats it as
// the mistyped invocation it almost always is. Declaring it here keeps the
// rule in the shared parser rather than having each subcommand re-derive it.
type Flag struct {
	Name         string
	HasValue     bool
	NoEmptyValue bool
}

// Grammar declares the argument shape one subcommand parses through Parse:
// its usage-line command name, its declared help text, the flags it accepts,
// its positional arity, and any compatibility rendering it preserves. MaxArgs
// of -1 means unbounded (variadic).
// ReservedPositionalsBeforeTerminator treats that many leading positional slots
// as literal values before recognizing help, flags, or the first terminator.
// HelpOnlyWhenSole narrows help to the sole-argument spelling and, with it, turns
// off the whole flag-recognition pass: the grammar accepts no flags, and "--" is
// an ordinary positional rather than a terminator, so Result.EndedFlags stays
// false.
type Grammar struct {
	Cmd                                 string
	Help                                string
	Flags                               []Flag
	MinArgs                             int
	MaxArgs                             int
	ReservedPositionalsBeforeTerminator int
	HelpOnlyWhenSole                    bool
	// UnquotedEmptyPositional retains an established usage response whose
	// unknown-argument cell was empty instead of the shared quoted marker.
	UnquotedEmptyPositional bool
}

// Result is a successful parse: the flags present (an empty string value for
// a boolean flag) and the positionals in argv order.
type Result struct {
	Flags                       map[string]string
	Positionals                 []string
	PositionalsBeforeTerminator int
	EndedFlags                  bool
}

// Parse applies g's grammar to args and returns exactly one of three
// outcomes: a populated Result with an empty line and code 0 on a successful
// parse; a help outcome (zero Result, g.Help as the line, code 0) when --help
// or -h appears after any reserved positionals and before "--", or when a help
// spelling is the sole argument; or a rendered usage error (zero Result, a
// toon.Usage or toon.MissingArg line, code 2). A caller prints a non-empty line
// and exits with code verbatim; Result is meaningful only when line == "".
//
// Rendering is delegated to toon.Usage and toon.MissingArg throughout, so
// this package never re-derives the usage-line shape.
func Parse(g Grammar, args []string) (Result, string, int) {
	// Help spellings remain requests when invoked alone, even for a grammar that
	// reserves its first positional slot as a literal value. Bare help is otherwise
	// input because variadic commands may legitimately receive it as free text.
	if len(args) == 1 {
		switch args[0] {
		case "--help", "-h":
			return Result{}, g.Help, 0
		case "help":
			return Result{}, g.Help, 0
		}
	}

	valueFlags := make(map[string]bool, len(g.Flags))
	knownFlags := make(map[string]bool, len(g.Flags))
	noEmptyFlags := make(map[string]bool, len(g.Flags))
	for _, f := range g.Flags {
		knownFlags[f.Name] = true
		if f.HasValue {
			valueFlags[f.Name] = true
		}
		if f.NoEmptyValue {
			noEmptyFlags[f.Name] = true
		}
	}

	result := Result{Flags: map[string]string{}}
	// endedFlags becomes permanent at the first "--"; every later "--" or
	// dash-prefixed token is then an ordinary positional rather than being
	// re-examined as a separator or a flag.
	endedFlags := false
	for i := 0; i < len(args); i++ {
		a := args[i]
		reserved := !endedFlags && result.PositionalsBeforeTerminator < g.ReservedPositionalsBeforeTerminator
		if !endedFlags && !reserved && !g.HelpOnlyWhenSole {
			if a == "--help" || a == "-h" {
				return Result{}, g.Help, 0
			}
			if a == "--" {
				endedFlags = true
				result.EndedFlags = true
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
					if args[i] == "" && noEmptyFlags[a] {
						return Result{}, toon.Usage(g.Cmd, a+` ""`), 2
					}
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
			if g.UnquotedEmptyPositional {
				return Result{}, toon.Usage(g.Cmd, ""), 2
			}
			return Result{}, toon.Usage(g.Cmd, `""`), 2
		}
		// Trailing garbage is reported on the first excess argument, not a
		// generic message, so a mistyped invocation names the token that
		// broke it.
		if g.MaxArgs >= 0 && len(result.Positionals) >= g.MaxArgs {
			return Result{}, toon.Usage(g.Cmd, a), 2
		}
		result.Positionals = append(result.Positionals, a)
		if !endedFlags {
			result.PositionalsBeforeTerminator++
		}
	}

	if len(result.Positionals) < g.MinArgs {
		return Result{}, toon.MissingArg(g.Cmd, "argument"), 2
	}
	return result, "", 0
}
