package spec

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/maps"
	"github.com/gibbonmi/bench/internal/tickets"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

const widenCmd = "bench spec widen"

// WidenUsage is the widen grammar, which parsing and public discovery share.
const WidenUsage = widenCmd + " <slug> --ticket <basename> [--decision <reference>] <entry>..."

var widenGrammar = usage.Grammar{
	Cmd:  widenCmd,
	Help: "usage: " + WidenUsage,
	Flags: []usage.Flag{
		{Name: "--ticket", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--decision", HasValue: true, NoEmptyValue: true},
	},
	MinArgs: 2,
	MaxArgs: -1,
}

// writesField is the prefix of the ticket's Writes: line. The edit locates that line
// through it, and tickets.ParseTicket re-reads the result, so a locator that disagrees
// with the grammar refuses instead of writing.
const writesField = "Writes:"

// writesSeparator joins the entries of one Writes: value, as maps.FieldList splits them.
const writesSeparator = ", "

// widenCommand runs `bench spec widen`. It adds each entry inside the signed-off scope to
// the one Writes: line of specs/<slug>/tickets/<basename> and makes no commit. The scope
// comes from the default-branch tip, so a build cannot widen the scope that grades its own
// widening. An entry the ticket already covers reports current; an entry outside the scope
// refuses unless --decision names the reviewer's admission. Any refusal writes nothing for
// the whole request.
func widenCommand(args []string) (string, int) {
	parsed, line, code := usage.Parse(widenGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	slug, basename, decision := parsed.Positionals[0], parsed.Flags["--ticket"], parsed.Flags["--decision"]
	base := RepoBase()
	if base == "" {
		return toon.NotInRepo() + "\n", 1
	}
	if refusal, _ := checkoutRefusal(base); refusal != "" {
		return refusal + "\n", 1
	}
	if !pathSegment(slug) || !pathSegment(basename) || path.Ext(basename) != tickets.Ext {
		return toon.Errorf(fmt.Sprintf("ticket address is not specs/<slug>/tickets/<basename>%s: %q, %q", tickets.Ext, slug, basename),
			"pass a slug and a ticket basename without a path") + "\n", 1
	}
	entries := parsed.Positionals[1:]
	for _, entry := range entries {
		if reason := hostileEntry(entry); reason != "" {
			return toon.Errorf(fmt.Sprintf("entry %q %s", entry, reason), "pass one clean repository path for each entry") + "\n", 1
		}
	}
	rel := ClosedFolderPath(slug) + "/tickets/" + basename
	ticketPath := filepath.Join(base, filepath.FromSlash(rel))
	if err := bounds.RefuseLinks(ticketPath); err != nil {
		return toon.RecordError(rel, bounds.StateWrongType, err.Error()) + "\n", 1
	}
	read := bounds.ClassifyNoFollow(ticketPath)
	if read.State != bounds.StateParsed {
		return toon.RecordError(rel, read.State, read.Reason) + "\n", 1
	}
	at, found := writesLine(read.Data)
	if !found {
		return toon.Errorf(rel+" has no Writes: line", "add the field by hand; widen edits only an existing line") + "\n", 1
	}
	scope, refusal := signedOffScope(base, slug)
	if refusal != "" {
		return refusal + "\n", 1
	}
	own, _ := tickets.ParseTicket(basename, read.Data, nil, "")
	covered := make([]string, 0, len(own.Writes)+len(entries))
	for _, entry := range own.Writes {
		p, _ := tickets.WritesPath(entry)
		covered = append(covered, p)
	}
	var rows [][]string
	var added, refused []string
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		p, _ := tickets.WritesPath(entry)
		state := widenAdded
		switch {
		case coveredBy(covered, p):
			state = widenCurrent
		case coveredBy(scope, p):
		case decision != "":
			state = widenDecided
		default:
			refused = append(refused, entry)
			continue
		}
		if state != widenCurrent {
			added = append(added, entry)
			covered = append(covered, p)
		}
		rows = append(rows, []string{rel, entry, state})
	}
	if len(refused) > 0 {
		return toon.Errorf("outside the signed-off scope: "+strings.Join(refused, writesSeparator),
			fmt.Sprintf("after the reviewer admits the path, rerun %s %s --ticket %s --decision <reference> %s", widenCmd, slug, basename, strings.Join(refused, " "))) + "\n", 1
	}
	table, err := toon.Table("widen", []string{"ticket", "entry", "state"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	if len(added) == 0 {
		return table, 0
	}
	widened, ok := widenLine(read.Data, at, own.Writes, added)
	if !ok {
		return toon.Errorf(rel+" would not parse back to the requested Writes: entries", "edit the Writes: line by hand") + "\n", 1
	}
	info, err := os.Lstat(ticketPath)
	if err == nil {
		err = durablefile.Replace(ticketPath, widened, info.Mode().Perm())
	}
	if err != nil {
		return toon.Errorf(fmt.Sprintf("write %s: %v", rel, err), "check file permissions") + "\n", 1
	}
	return table, 0
}

// The states of one widen row.
const (
	widenAdded   = "added"
	widenDecided = "decided"
	widenCurrent = "current"
)

// pathSegment reports whether name is one representable path element below a folder.
func pathSegment(name string) bool {
	return tickets.RepositoryPath(name) && !strings.Contains(name, "/") && toon.Representable(name)
}

// hostileEntry names why one entry cannot join a Writes: line, or returns "". The
// grammar grades the path, the field list grades the separator, and the line sink
// refuses the line breaks that TOON permits.
func hostileEntry(entry string) string {
	switch {
	case !toon.Representable(entry):
		return "holds a byte that TOON cannot render"
	case strings.ContainsAny(entry, "\r\n"):
		return "holds a line break"
	case len(maps.FieldList(entry, "", "")) != 1:
		return "holds the " + writesSeparator + " separator"
	}
	if _, err := tickets.ParseWrite(entry); err != nil {
		return "is not a Writes: path: " + err.Error()
	}
	return ""
}

// coveredBy reports whether one of the entry paths covers p.
func coveredBy(entries []string, p string) bool {
	return slices.ContainsFunc(entries, func(entry string) bool { return tickets.Covers(entry, p) })
}

// signedOffScope reads the directory prefix of every Writes: entry of the slug's tickets
// at the default-branch tip. The prefix of an entry is the entry when the tip holds it as
// a directory, else its parent directory; a root-level file is its own prefix. A tip
// without a ticket of the slug refuses, so an absent scope never reads as an open one.
func signedOffScope(root, slug string) ([]string, string) {
	branch, ok := git.ResolvedDefault(root)
	if !ok {
		return nil, toon.Errorf("default branch is unresolved", "widen reads the signed-off scope from the default-branch tip")
	}
	tip, err := git.ResolveCommit(root, branch)
	if err != nil {
		return nil, toon.Errorf("default branch tip is unreadable: "+branch, err.Error())
	}
	folder := ClosedFolderPath(slug) + "/tickets/"
	listing, err := git.ListTree(root, "-r", tip, "--", folder)
	if err != nil {
		return nil, toon.Errorf("signed-off tickets are unreadable at "+branch, err.Error())
	}
	var scope []string
	signed := false
	for _, entry := range listing {
		if !entry.IsRegularFile() || !strings.HasPrefix(entry.Path, folder) || path.Ext(entry.Path) != tickets.Ext {
			continue
		}
		data, err := git.ReadControlBlob(root, entry.Object)
		if err != nil {
			return nil, toon.Errorf("signed-off ticket is unreadable: "+entry.Path, err.Error())
		}
		signed = true
		ticket, _ := tickets.ParseTicket(path.Base(entry.Path), data, nil, "")
		for _, raw := range ticket.Writes {
			write, err := tickets.ParseWrite(raw)
			if err != nil {
				continue
			}
			prefix, err := scopePrefix(root, tip, write.Path)
			if err != nil {
				return nil, toon.Errorf("signed-off entry is unreadable at "+branch+": "+write.Path, err.Error())
			}
			scope = append(scope, prefix)
		}
	}
	if !signed {
		return nil, toon.Errorf(branch+" holds no ticket of "+ClosedFolderPath(slug), "land the signed-off tickets before a widening")
	}
	return scope, ""
}

// scopePrefix is the directory prefix of one signed-off entry at tip.
func scopePrefix(root, tip, entry string) (string, error) {
	listing, err := git.ListTree(root, tip, "--", entry)
	if err != nil {
		return "", err
	}
	for _, held := range listing {
		if held.Path == entry && held.Type == "tree" {
			return entry, nil
		}
	}
	if dir := path.Dir(entry); dir != "." {
		return dir, nil
	}
	return entry, nil
}

// writesLine returns the index of the first unfenced Writes: line, the one the ticket
// grammar reads, in the LF split of content.
func writesLine(content []byte) (int, bool) {
	scanned, _ := maps.FieldScan{Table: []maps.FieldSpec{{Name: writesField, Syntax: writesField}}}.Scan(content)
	for i, line := range scanned {
		if line.Field == writesField && line.Diagnostic == "" {
			return i, true
		}
	}
	return 0, false
}

// widenLine appends the added entries to the value of line at and keeps every other
// byte, including a carriage return and the whitespace after the value. A value with no
// entry, such as the none sentinel, gives way to the added entries. The result must
// parse back to the old entries plus the added ones.
func widenLine(content []byte, at int, old, added []string) ([]byte, bool) {
	lines := strings.Split(string(content), "\n")
	body, cr := strings.CutSuffix(lines[at], "\r")
	value := strings.TrimPrefix(body, writesField)
	kept := strings.TrimRightFunc(value, unicode.IsSpace)
	tail := value[len(kept):]
	if len(old) == 0 {
		kept = " "
	} else {
		kept += writesSeparator
	}
	body = writesField + kept + strings.Join(added, writesSeparator) + tail
	if cr {
		body += "\r"
	}
	lines[at] = body
	widened := []byte(strings.Join(lines, "\n"))
	reread, _ := tickets.ParseTicket("", widened, nil, "")
	return widened, slices.Equal(reread.Writes, append(slices.Clip(old), added...))
}
