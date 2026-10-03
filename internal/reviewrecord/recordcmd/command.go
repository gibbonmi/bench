// Package recordcmd is the `bench record` verb. It parses one form's flags, applies the
// refusal order, and prints one TOON row; package reviewrecord derives every field and
// owns the write transaction.
package recordcmd

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

const command = "bench record"

// flag is one flag of a form and the placeholder that its usage line shows. An empty
// placeholder declares a flag without a value, a single-line flag refuses a value that
// holds a control character, and occurs says only whether the flag may repeat.
type flag struct {
	name, placeholder string
	singleLine        bool
	occurs            occurrence
}

// occurrence is how often a flag may appear in one call.
type occurrence int

const (
	once occurrence = iota
	repeated
)

// form is one `bench record` form. Its usage line, its grammar, and its help row all
// derive from this one declaration, so the help cannot advertise another grammar. The
// layout orders the usage terms: each term is a declared flag name or one of the
// grouping marks ( ) [ ] |, and an empty layout lists the flags in order. A repeated
// flag shows as [<flag>]... in the usage line. The layout is the one source for which
// flags a call names together and which flags Parse requires: admits reads its groups.
// valid holds the grammar rules that close a value set, and a nil valid has none.
type form struct {
	name, description string
	flags             []flag
	layout            []string
	valid             func(f form, parsed usage.Result) bool
	run               func(f form, root, spec string, parsed usage.Result) (string, int)
}

// probeFlags are the three flags of a planned probe.
var probeFlags = []string{"--probe-outcome", "--probe-exit-code", "--probe-restore"}

var forms = []form{
	{name: "chunk", description: "write one chunk's frozen pair, digests, and acceptance rows into reviews/<slug>.md",
		flags: []flag{{"--chunk", "<id>", true, once}, {"--base", "<commit>", false, once}, {"--tip", "<commit>", false, once}}, run: chunk},
	{name: "verification", description: "append one planned verification result with its computed digests",
		flags: []flag{{"--chunk", "<id>", true, once}, {"--source", "<commit>", false, once}, {"--final", "", false, once},
			{"--requirement", "<id>", true, once}, {"--id", "<id>", true, once}, {"--performer", "<session>", true, once},
			{"--model", "<model>", true, once}, {"--effort", "<effort>", true, once}, {"--exit-code", "<n>", false, once},
			{"--ref", "<ref>", true, once}, {"--excerpt", "<file>", true, once}, {probeFlags[0], "<verdict>", true, once},
			{probeFlags[1], "<n>", false, once}, {probeFlags[2], "pass|fail", false, once}},
		layout: []string{"(", "--chunk", "[", "--source", "]", "|", "--final", "--source", ")", "--requirement", "--id", "--performer",
			"--model", "--effort", "--exit-code", "--ref", "--excerpt", "[", probeFlags[0], probeFlags[1], probeFlags[2], "]"},
		valid: verificationValid, run: verification},
	{name: "review", description: "append one independent review result to a recorded chunk",
		flags: []flag{{"--chunk", "<id>", true, once}, {"--axis", strings.Join(reviewrecord.Axes(), "|"), false, once}, {"--id", "<id>", true, once},
			{"--performer", "<session>", true, once}, {"--model", "<model>", true, once}, {"--effort", "<effort>", true, once},
			{"--ref", "<ref>", true, once}, {"--excerpt", "<file>", true, once}, {"--finding", "<id>", true, repeated}},
		valid: reviewValid, run: review},
	completionForm,
	amendmentForm,
}

// terms is the layout of f, or its flags in order when the layout is empty.
func (f form) terms() []string {
	layout := f.layout
	if layout == nil {
		for _, flag := range f.flags {
			layout = append(layout, flag.name)
		}
	}
	return layout
}

// suffix is the form's grammar after `bench record`.
func (f form) suffix() string {
	terms := []string{"", f.name, "<slug>"}
	for _, term := range f.terms() {
		if strings.HasPrefix(term, "--") {
			term = f.flag(term).usage()
		}
		terms = append(terms, term)
	}
	return strings.NewReplacer("( ", "(", "[ ", "[", " )", ")", " ]", "]").Replace(strings.Join(terms, " "))
}

// index returns the position of the flag of f named name. A layout that names an
// undeclared flag is a defect in the declaration, so it panics on the first help render.
func (f form) index(name string) int {
	i := slices.IndexFunc(f.flags, func(fl flag) bool { return fl.name == name })
	if i < 0 {
		panic("bench record " + f.name + " layout names the undeclared flag " + name)
	}
	return i
}

// flag returns the flag of f named name.
func (f form) flag(name string) flag {
	return f.flags[f.index(name)]
}

// admits reads the layout terms from index i to the close of the current group. It
// returns each flag set that the terms admit, one bit for each flag of f, and the index
// after the close. A sequence takes one set from each term and a ( ) group takes one
// alternative. A [ ] group and a repeated flag, shown as [<flag>]..., also admit no flag,
// so an optional group admits all of its flags or none.
func (f form) admits(terms []string, i int) ([]uint64, int) {
	alternatives, sequence := []uint64{}, []uint64{0}
	for i < len(terms) {
		term := terms[i]
		i++
		var sets []uint64
		switch term {
		case ")", "]":
			return append(alternatives, sequence...), i
		case "|":
			alternatives, sequence = append(alternatives, sequence...), []uint64{0}
			continue
		case "(", "[":
			sets, i = f.admits(terms, i)
			if term == "[" {
				sets = append(sets, 0)
			}
		default:
			index := f.index(term)
			sets = []uint64{1 << index}
			if f.flags[index].occurs == repeated {
				sets = append(sets, 0)
			}
		}
		product := []uint64{}
		for _, prefix := range sequence {
			for _, set := range sets {
				product = append(product, prefix|set)
			}
		}
		sequence = product
	}
	return append(alternatives, sequence...), i
}

// together reports whether the layout of f admits the set of flags present in values.
func (f form) together(values map[string]string) bool {
	var present uint64
	for i, flag := range f.flags {
		if _, ok := values[flag.name]; ok {
			present |= 1 << i
		}
	}
	sets, _ := f.admits(f.terms(), 0)
	return slices.Contains(sets, present)
}

func (fl flag) usage() string {
	term := strings.TrimSpace(fl.name + " " + fl.placeholder)
	if fl.occurs == repeated {
		term = "[" + term + "]..."
	}
	return term
}

func (f form) grammar() usage.Grammar {
	g := usage.Grammar{Cmd: command + " " + f.name, Help: "usage: " + command + f.suffix(), MinArgs: 1, MaxArgs: 1}
	sets, _ := f.admits(f.terms(), 0)
	always := ^uint64(0)
	for _, set := range sets {
		always &= set
	}
	for i, flag := range f.flags {
		g.Flags = append(g.Flags, usage.Flag{Name: flag.name, HasValue: flag.placeholder != "", NoEmptyValue: flag.placeholder != "", Required: always&(1<<i) != 0, Repeatable: flag.occurs == repeated})
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
	if !f.together(parsed.Flags) || f.valid != nil && !f.valid(f, parsed) {
		return f.grammar().Help + "\n", 2
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
		for _, value := range append([]string{parsed.Flags[flag.name]}, parsed.Repeated[flag.name]...) {
			if flag.singleLine && !sanitize.LineSafe(value) {
				return f.refuse(flag.name + " holds a control character")
			}
		}
	}
	spec := "specs/" + parsed.Positionals[0] + "/spec.md"
	if _, err := reviewrecord.Slug(spec); err != nil {
		return f.refuse(err.Error())
	}
	return f.run(f, root, spec, parsed)
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

func chunk(f form, root, spec string, parsed usage.Result) (string, int) {
	values := parsed.Flags
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
	return row("chunk", []string{"id", "action", "base", "tip", "source_digest", "plan_digest", "rows"},
		[]string{entry.ID, action, entry.Base, entry.Tip, entry.SourceDigest, entry.PlanDigest, strconv.Itoa(len(entry.AcceptanceRows))})
}

// row prints the one output table of a form: its name, its header, and its one row.
func row(name string, header, cells []string) (string, int) {
	table, err := toon.Table(name, header, [][]string{cells})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return table, 0
}

// integer is the base-10 integer value of one flag.
func integer(values map[string]string, name string) (int, bool) {
	n, err := strconv.Atoi(values[name])
	return n, err == nil
}

// chosen reports whether the value of flag name is one of the alternatives of its
// placeholder, so the usage line and the check read one value set.
func chosen(f form, values map[string]string, name string) bool {
	return slices.Contains(strings.Split(f.flag(name).placeholder, "|"), values[name])
}

// verificationValid closes the value sets of the verification form. Command calls it
// only for a set of flags that the layout admits.
func verificationValid(f form, parsed usage.Result) bool {
	_, exitCode := integer(parsed.Flags, "--exit-code")
	_, probe := parsed.Flags[probeFlags[0]]
	_, probeExitCode := integer(parsed.Flags, probeFlags[1])
	return exitCode && (!probe || probeExitCode && chosen(f, parsed.Flags, probeFlags[2]))
}

// native reads the --excerpt file as exact bytes and pairs them with --ref. The read
// never follows a link or opens a special file, and every state other than parsed
// refuses with the record error line for the operand path.
func native(values map[string]string) (reviewrecord.NativeRef, string, int) {
	path := values["--excerpt"]
	c := bounds.ClassifyNoFollow(path)
	if c.State != bounds.StateParsed {
		return reviewrecord.NativeRef{}, toon.RecordError(path, c.State, c.Reason) + "\n", 1
	}
	return reviewrecord.NativeRef{Ref: values["--ref"], Excerpt: string(c.Data), Digest: reviewrecord.Digest(c.Data)}, "", 0
}

// cause is the refusal cause of err. A missing chunk entry names the chunk form, and a
// probe mismatch names the probe flags.
func cause(err error) string {
	switch {
	case errors.Is(err, reviewrecord.ErrNoChunkEntry):
		return err.Error() + "; write it with " + command + " chunk"
	case errors.Is(err, reviewrecord.ErrProbe):
		return err.Error() + "; the probe flags are " + strings.Join(probeFlags, ", ")
	}
	return err.Error()
}

func verification(f form, root, spec string, parsed usage.Result) (string, int) {
	values := parsed.Flags
	ref, out, code := native(values)
	if out != "" {
		return out, code
	}
	call := reviewrecord.VerificationCall{Chunk: values["--chunk"], Requirement: values["--requirement"],
		Evidence: reviewrecord.Evidence{ID: values["--id"], Performer: values["--performer"], Model: values["--model"], Effort: values["--effort"], NativeRef: ref}}
	if _, ok := values["--source"]; ok {
		if call.Source, out, code = commit(f, root, values, "--source"); out != "" {
			return out, code
		}
	}
	call.ExitCode, _ = integer(values, "--exit-code")
	if _, ok := values[probeFlags[0]]; ok {
		exitCode, _ := integer(values, probeFlags[1])
		call.Probe = &reviewrecord.Probe{Outcome: values[probeFlags[0]], ExitCode: exitCode, Restore: values[probeFlags[2]]}
	}
	entry, err := reviewrecord.RecordVerification(root, spec, call)
	if err != nil {
		return f.refuse(cause(err))
	}
	list := "chunk"
	if call.Chunk == "" {
		list = "completion"
	}
	return row("verification", []string{"list", "chunk", "id", "requirement", "role", "outcome", "source_digest", "excerpt_digest"},
		[]string{list, call.Chunk, entry.ID, entry.Requirement, entry.Role, entry.Outcome, entry.SourceDigest, entry.NativeRef.Digest})
}

// reviewValid closes the --axis value set.
func reviewValid(f form, parsed usage.Result) bool {
	return chosen(f, parsed.Flags, "--axis")
}

func review(f form, root, spec string, parsed usage.Result) (string, int) {
	values := parsed.Flags
	ref, out, code := native(values)
	if out != "" {
		return out, code
	}
	call := reviewrecord.ReviewCall{Chunk: values["--chunk"], Axis: values["--axis"], Findings: parsed.Repeated["--finding"],
		Evidence: reviewrecord.Evidence{ID: values["--id"], Performer: values["--performer"], Model: values["--model"], Effort: values["--effort"], NativeRef: ref}}
	entry, err := reviewrecord.RecordReview(root, spec, call)
	if err != nil {
		return f.refuse(cause(err))
	}
	return row("review", []string{"chunk", "id", "axis", "outcome", "supersedes", "source_digest", "excerpt_digest"},
		[]string{call.Chunk, entry.ID, entry.Axis, entry.Outcome, strings.Join(entry.Supersedes, ","), entry.SourceDigest, entry.NativeRef.Digest})
}
