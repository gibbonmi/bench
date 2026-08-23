package handoff

import (
	"path/filepath"
	"strings"

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

// render composes the pin block — everything a cold session needs to know before it acts.
// The document written to disk is this block plus the Shape section, so the two sinks
// cannot state different facts about the same tree.
func render(f facts, state string) string {
	var b strings.Builder
	b.WriteString("# Session handoff\n\n")
	b.WriteString("Repository: `" + f.Repo + "` (" + originField(f.Origin) + ")\n")
	b.WriteString("Path: `" + f.Path + "`\n")
	b.WriteString("Branch: " + f.Branch + " — " + f.Head + ", " + f.Dirty + ", " + f.Unpushed + "\n")
	b.WriteString("Spec: " + specField(f.Specs) + "\n")
	b.WriteString("Gate: " + f.Gate + "\n\n")
	if state == "" {
		b.WriteString(scaffoldGuidance + "\n")
	}
	b.WriteString(stateHeading + "\n\n")
	if state != "" {
		b.WriteString(state + "\n\n")
	}
	b.WriteString(nextHeading + "\n\n" + nextField(f) + "\n")
	return b.String()
}

// document is the file's full text: the pin block plus the Shape section stdout omits.
func document(pin string) string {
	return pin + "\n" + ShapeHeading + "\n\n" + ShapeSection
}

func originField(origin string) string {
	if origin == "" {
		return "origin unknown"
	}
	return "origin `" + origin + "`"
}

func specField(specs []string) string {
	if len(specs) == 0 {
		return "none staged."
	}
	return strings.Join(specs, ", ")
}

// nextField states the next command. A board with signals but no command among them says
// so, and points at the override. This field promises an invocation, and the hints left
// on such a board are not ones. An empty Signal marks an overridden action, rendered
// alone — the board had nothing to do with it. Naming a signal would be a claim
// about a derivation that never ran.
func nextField(f facts) string {
	if f.NoInvocable {
		return "No invocable command derives from the board; name the next step with `--next`."
	}
	if f.Action == "" {
		return "Nothing pending — the board is clean."
	}
	if f.Signal == "" {
		return "`" + f.Action + "`"
	}
	return "`" + f.Action + "` — the board's leading invocable signal (`" + f.Signal + "`)."
}

// validate refuses any field that cannot survive the sink it is about to be composed into,
// before a line is composed. A control byte reaching the rendered block would ride into
// every downstream reader of the artifact, so the refusal is the whole answer.
//
// The sink is a line-structured markdown document, which is stricter than TOON.
// toon.Representable permits tab, newline, and return because the encoder escapes them,
// and none of the three survives here. A newline is the sharp one: a value carrying it
// splits its own field across lines. A `--next` override carrying one can write a
// second `## State` heading that makes every later run refuse the document as ambiguous.
func validate(f facts) error {
	fields := append([]string{f.Repo, f.Origin, f.Path, f.Branch, f.Head, f.Dirty, f.Unpushed, f.Gate, f.Action, f.Signal}, f.Specs...)
	for _, value := range fields {
		if !toon.Representable(value) || strings.ContainsAny(value, "\n\r\t") {
			return refusal{"unrepresentable handoff field",
				"a derived value carries a control byte; fix the branch, spec, or status text that holds it"}
		}
	}
	return nil
}
