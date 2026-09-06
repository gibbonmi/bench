package handoff

import (
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/handoffdoc"
	"github.com/gibbonmi/bench/internal/toon"
)

// renderPath renders the git root as a reader elsewhere should see it: abbreviated to `~`
// when it is $HOME or sits beneath it, absolute otherwise. The containment test is by path
// component. A $HOME of /home/a leaves /home/abc absolute rather than turning it into
// the `~bc` a raw prefix match would produce. An unset $HOME abbreviates nothing.
func renderPath(root, home string) string {
	if home == "" || !filepath.IsAbs(home) {
		return root
	}
	root, home = filepath.Clean(root), filepath.Clean(home)
	if root == home {
		return "~"
	}
	rest, ok := strings.CutPrefix(root, strings.TrimSuffix(home, string(filepath.Separator))+string(filepath.Separator))
	if !ok {
		return root
	}
	return "~/" + rest
}

// header composes the block above the first section — the repository facts every section
// shares. It names the checkout the document lives in, so the HEAD it pins is main's and
// not the caller worktree's. A section's own tip is a label line under that section.
//
// The scaffold guidance rides in the header while the owned section has no State. It is
// what the first session in a repo reads, and it goes once that session writes its State.
func header(f facts, state string) string {
	var b strings.Builder
	b.WriteString("# Session handoff\n\n")
	b.WriteString("Repository: `" + f.Repo + "` (" + originField(f.Origin) + ")\n")
	b.WriteString("Path: `" + f.Path + "`\n")
	b.WriteString("Branch: " + f.Branch + " — " + f.Head + ", " + f.Dirty + ", " + f.Unpushed + "\n")
	b.WriteString("Gate: " + f.Gate)
	if state == "" {
		b.WriteString("\n\n" + scaffoldGuidance)
	}
	return b.String()
}

// section is the owned block, ready for the leaf package to place. State and Next both
// arrive from the document the run read, because both can be a reviewer's own words: this
// command derives every other field.
func section(f facts, state, next string) handoffdoc.Section {
	return handoffdoc.Section{Key: f.Key, Fields: f.Pins, Next: next, State: state}
}

// preview is what the command prints: the header and the one section this run wrote. It
// renders through the leaf package, so stdout and the file cannot spell a field two ways.
// The Shape body is left off, because a reader who has the command's output does not need
// the document's own instructions.
func preview(head string, owned handoffdoc.Section) string {
	return (&handoffdoc.Document{Header: head, Sections: []handoffdoc.Section{owned}}).Render()
}

func originField(origin string) string {
	if origin == "" {
		return "origin unknown"
	}
	return "origin `" + origin + "`"
}

// nextField states the next command. A board with signals but no command among them says
// so, and points at the override. This field promises an invocation, and the hints left
// on such a board are not ones. An empty Signal marks an overridden action, rendered
// alone — the board had nothing to do with it. Naming a signal would be a claim
// about a derivation that never ran.
//
// No form ends in a sentence terminator. The value is a label line's, and an unterminated
// label-shaped line is what the prose lane skips as a template field.
func nextField(f facts) string {
	if f.NoInvocable {
		return "No invocable command derives from the board; name the next step with `--next`"
	}
	if f.Action == "" {
		return "Nothing pending — the board is clean"
	}
	if f.Signal == "" {
		return "`" + f.Action + "`"
	}
	return "`" + f.Action + "` — the board's leading invocable signal (`" + f.Signal + "`)"
}

// blankNext answers whether a section's own Next command names nothing, so the board's
// leading signal is the better value. A value the reviewer or an earlier run put there is
// kept byte for byte: a mid-build resume invocation carries flags the board cannot derive,
// and regenerating over it loses them.
//
// The rendered form of a named command is a backticked span, so an empty pair of backticks
// is the same statement as an empty line. A rule that read the backticks as content would
// leave that section routeless forever.
func blankNext(next string) bool {
	trimmed := strings.TrimSpace(next)
	return trimmed == "" || trimmed == "``"
}

// validate refuses any field that cannot survive the sink it is about to be composed into,
// before a line is composed. A control byte reaching the rendered block would ride into
// every downstream reader of the artifact, so the refusal is the whole answer.
//
// The sink is a line-structured markdown document, which is stricter than TOON.
// toon.Representable permits tab, newline, and return because the encoder escapes them,
// and none of the three survives here. A newline is the sharp one: a value carrying it
// splits its own field across lines. A `--next` override carrying one can write a second
// section heading that makes every later run refuse the document it can no longer parse.
func validate(f facts) error {
	fields := []string{f.Repo, f.Origin, f.Path, f.Branch, f.Head, f.Dirty, f.Unpushed, f.Gate, f.Action, f.Signal, f.Key}
	for _, pin := range f.Pins {
		fields = append(fields, pin.Value)
	}
	for _, value := range fields {
		if !toon.Representable(value) || strings.ContainsAny(value, "\n\r\t") {
			return refusal{"unrepresentable handoff field",
				"a derived value carries a control byte; fix the branch, spec, or status text that holds it"}
		}
	}
	return nil
}
