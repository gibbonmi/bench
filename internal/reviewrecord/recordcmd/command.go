// Package recordcmd is the `bench record` verb. It parses one form's flags, applies the
// refusal order, and prints one TOON row; package reviewrecord derives every field and
// owns the write transaction.
package recordcmd

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

const command = "bench record"

// flag is one flag of a form and the placeholder that its usage line shows. A single-line
// flag refuses a value that holds a control character.
type flag struct {
	name, placeholder string
	singleLine        bool
}

// form is one `bench record` form. Its usage line, its grammar, and its help row all
// derive from this one declaration, so the help cannot advertise another grammar.
type form struct {
	name, description string
	flags             []flag
	run               func(f form, root, spec string, values map[string]string) (string, int)
}

var forms = []form{
	{name: "chunk", description: "write one chunk's frozen pair, digests, and acceptance rows into reviews/<slug>.md",
		flags: []flag{{"--chunk", "<id>", true}, {"--base", "<commit>", false}, {"--tip", "<commit>", false}}, run: chunk},
}

// suffix is the form's grammar after `bench record`.
func (f form) suffix() string {
	terms := []string{"", f.name, "<slug>"}
	for _, flag := range f.flags {
		terms = append(terms, flag.name, flag.placeholder)
	}
	return strings.Join(terms, " ")
}

func (f form) grammar() usage.Grammar {
	g := usage.Grammar{Cmd: command + " " + f.name, Help: "usage: " + command + f.suffix(), MinArgs: 1, MaxArgs: 1}
	for _, flag := range f.flags {
		g.Flags = append(g.Flags, usage.Flag{Name: flag.name, HasValue: true, NoEmptyValue: true, Required: true})
	}
	return g
}

func (f form) refuse(cause string) (string, int) {
	return toon.Errorf(command+" "+f.name+" refused", cause) + "\n", 1
}

// HelpRow is one public help row: the text after `bench record` and its description.
type HelpRow struct {
	Suffix, Description string
}

// HelpRows projects the forms for the root help inventory.
func HelpRows() []HelpRow {
	rows := make([]HelpRow, len(forms))
	for i, f := range forms {
		rows[i] = HelpRow{Suffix: f.suffix(), Description: f.description}
	}
	return rows
}

func usageText() string {
	var b strings.Builder
	for _, f := range forms {
		b.WriteString(f.grammar().Help + "\n")
	}
	return b.String()
}

// Command runs `bench record <form> <slug> …` in the worktree at root. An empty root
// means that the caller is outside a repository. A grammar error exits 2, and every
// content refusal exits 1 in the refusal order of the verb.
func Command(root string, args []string) (string, int) {
	family := usage.Grammar{Cmd: command, Help: strings.TrimSuffix(usageText(), "\n"), HelpOnlyWhenSole: true, MaxArgs: -1}
	if _, line, code := usage.Parse(family, args); line != "" {
		return line + "\n", code
	}
	if len(args) == 0 {
		return usageText(), 2
	}
	var selected *form
	for i := range forms {
		if forms[i].name == args[0] {
			selected = &forms[i]
		}
	}
	if selected == nil {
		return toon.Usage(command, args[0]) + "\n", 2
	}
	f := *selected
	parsed, line, code := usage.Parse(f.grammar(), args[1:])
	if line != "" {
		return line + "\n", code
	}
	if root == "" {
		return toon.NotInRepo() + "\n", 1
	}
	primary, err := git.IsPrimaryCheckout(root)
	if err != nil {
		return f.refuse(err.Error())
	}
	if primary {
		return usage.PrimaryCheckoutRefusal() + "\n", 1
	}
	for _, flag := range f.flags {
		if flag.singleLine && strings.IndexFunc(parsed.Flags[flag.name], unicode.IsControl) >= 0 {
			return f.refuse(flag.name + " holds a control character")
		}
	}
	spec := "specs/" + parsed.Positionals[0] + "/spec.md"
	if _, err := reviewrecord.Slug(spec); err != nil {
		return f.refuse(err.Error())
	}
	return f.run(f, root, spec, parsed.Flags)
}

// commit resolves the revision of one flag to its full commit ID.
func commit(f form, root string, values map[string]string, name string) (string, string, int) {
	id, err := git.ResolveCommit(root, values[name], "--quiet", "--end-of-options")
	if err != nil {
		out, code := f.refuse(fmt.Sprintf("%s %q names no commit", name, values[name]))
		return "", out, code
	}
	return id, "", 0
}

func chunk(f form, root, spec string, values map[string]string) (string, int) {
	base, out, code := commit(f, root, values, "--base")
	if out != "" {
		return out, code
	}
	tip, out, code := commit(f, root, values, "--tip")
	if out != "" {
		return out, code
	}
	action, entry, err := reviewrecord.RecordChunk(root, spec, values["--chunk"], base, tip)
	if err != nil {
		return f.refuse(err.Error())
	}
	table, err := toon.Table("chunk", []string{"id", "action", "base", "tip", "source_digest", "plan_digest", "rows"},
		[][]string{{entry.ID, action, entry.Base, entry.Tip, entry.SourceDigest, entry.PlanDigest, strconv.Itoa(len(entry.AcceptanceRows))}})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return table, 0
}
