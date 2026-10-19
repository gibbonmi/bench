package ownership

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/structure"
	"github.com/gibbonmi/bench/internal/tickets"
)

// sizingPremises is every sizing premise kind and its collector. A collector answers a
// reason only on complete supported proof. Every form that it does not prove answers an
// error, so its premise stays unknown.
var sizingPremises = map[string]func(*census, PremiseItem) (string, error){
	"lines":        (*census).lines,
	"glob-members": (*census).globMembers,
	"headroom":     (*census).headroom,
}

// lines grades one physical line count of a tracked regular file at its declared source,
// through the line rule of the structure engine. A path that the source does not track
// has no observed count, so an expected count for it stays an estimate.
func (c *census) lines(item PremiseItem) (string, error) {
	var expected struct {
		Count *int `json:"count"`
	}
	subject, err := decodeLinesSubject(item)
	if err == nil {
		_, err = typed(item.Expected, "expected: ", &expected, "count")
	}
	if err == nil && (expected.Count == nil || *expected.Count < 0) {
		err = errors.New("expected: count is not a nonnegative integer")
	}
	if err == nil && (subject.Path == "" || !tickets.RepositoryPath(subject.Path)) {
		err = fmt.Errorf("subject: path %q is no clean repository path", subject.Path)
	}
	var rev string
	if err == nil {
		rev, err = c.rev(subject.Source, subject.Path)
	}
	var tracked []string
	if err == nil {
		tracked, err = c.trackedAt(rev)
	}
	if err != nil {
		return "", err
	}
	if !slices.Contains(tracked, subject.Path) {
		return "", fmt.Errorf("%s is no tracked member at the source %s %s, so its expected count is an unobserved estimate", subject.Path, subject.Source, rev)
	}
	data, err := git.ReadTreeFile(c.root, rev, subject.Path)
	if err != nil {
		return "", fmt.Errorf("%s does not read at the source %s %s: %w", subject.Path, subject.Source, rev, err)
	}
	at := fmt.Sprintf("%s holds %d line(s) at the source %s %s", subject.Path, structure.Lines(data), subject.Source, rev)
	if structure.Lines(data) != *expected.Count {
		return "", fmt.Errorf("%s, not the expected %d", at, *expected.Count)
	}
	return at, nil
}

// globMembers grades one complete member list of a filepath.Match pattern over the tracked
// paths at its declared source, dot paths included. The pattern must stay inside a clean
// directory scope that holds no metacharacter. Recursive ** syntax, which filepath.Match
// does not support, and malformed syntax stay unproved. A member is an exact path, never a
// pattern.
func (c *census) globMembers(item PremiseItem) (string, error) {
	var expected struct {
		Members []string `json:"members"`
	}
	subject, err := decodeGlobSubject(item)
	if err == nil {
		_, err = typed(item.Expected, "expected: ", &expected, "members")
	}
	if err != nil {
		return "", err
	}
	pattern, scope := subject.Pattern, subject.Scope
	dir, cut := strings.CutSuffix(scope, "/")
	if _, malformed := filepath.Match(pattern, ""); malformed != nil {
		err = fmt.Errorf("pattern %q is not filepath.Match syntax: %v", pattern, malformed)
	}
	switch {
	case strings.Contains(pattern, "**"):
		err = fmt.Errorf("pattern %q uses recursive ** syntax, which filepath.Match does not support", pattern)
	case err != nil:
	case !cut || dir == "" || !tickets.RepositoryPath(dir) || strings.ContainsAny(dir, `*?[\`):
		err = fmt.Errorf("scope %q is no clean directory prefix free of pattern metacharacters", scope)
	case !strings.HasPrefix(pattern, scope) || !tickets.RepositoryPath(pattern):
		err = fmt.Errorf("pattern %q escapes its scope %q", pattern, scope)
	case expected.Members == nil || !slices.IsSorted(expected.Members) || len(slices.Compact(slices.Clone(expected.Members))) != len(expected.Members):
		err = errors.New("expected: members is no complete sorted list of exact paths")
	}
	var rev string
	if err == nil {
		rev, err = c.rev(subject.Source, pattern)
	}
	var tracked []string
	if err == nil {
		tracked, err = c.trackedAt(rev)
	}
	if err != nil {
		return "", err
	}
	var actual []string
	for _, path := range tracked {
		if matched, _ := filepath.Match(pattern, path); matched && strings.HasPrefix(path, scope) {
			actual = append(actual, path)
		}
	}
	if err := pathGaps(expected.Members, actual); err != nil {
		return "", err
	}
	return fmt.Sprintf("the glob %s holds %d tracked member(s) at the source %s %s", pattern, len(actual), subject.Source, rev), nil
}

// headroom grades one directory's current file headroom at the source tip through the
// structure engine's limit, grant, and file-count owners, after an authored plan of exact
// creations and deletions. The checkout must be the source tip with its limit lists
// unchanged, because the engine reads those lists from the checkout.
func (c *census) headroom(item PremiseItem) (string, error) {
	var plan struct {
		Create []string `json:"create"`
		Delete []string `json:"delete"`
	}
	var expected struct {
		Maximum     *int     `json:"maximum"`
		Current     *int     `json:"current"`
		Net         *int     `json:"net"`
		BudgetOwner Location `json:"budget_owner"`
		GrantOwner  Location `json:"grant_owner"`
	}
	subject, err := decodeHeadroomSubject(item)
	if err == nil {
		_, err = typed(subject.Files, "subject: files: ", &plan, "create", "delete")
	}
	var object map[string]json.RawMessage
	if err == nil {
		object, err = typed(item.Expected, "expected: ", &expected, "maximum", "current", "net", "budget_owner", "grant_owner")
	}
	if err == nil {
		err = firstError(location(object, "budget_owner", "expected: ", true), location(object, "grant_owner", "expected: ", true))
	}
	dir := subject.Directory
	switch {
	case err != nil:
	case plan.Create == nil || plan.Delete == nil:
		err = errors.New("subject: files: create and delete are not JSON arrays")
	case expected.Maximum == nil || expected.Current == nil || expected.Net == nil:
		err = errors.New("expected: maximum, current, and net are not integers")
	case dir == "" || !tickets.RepositoryPath(dir):
		err = fmt.Errorf("subject: directory %q is no clean repository directory", dir)
	}
	if err == nil {
		err = c.atTip()
	}
	if err == nil {
		err = c.limitsUnchanged()
	}
	for _, owner := range []Location{expected.BudgetOwner, expected.GrantOwner} {
		if err == nil {
			_, err = c.locate(owner)
		}
	}
	var tracked []string
	if err == nil {
		tracked, err = c.trackedAt(c.source.Tip)
	}
	if err == nil {
		err = planned(dir, tracked, plan.Create, plan.Delete)
	}
	var fit structure.Headroom
	if err == nil {
		fit, err = structure.DirectoryHeadroom(c.root, tracked, dir)
	}
	if err != nil {
		return "", err
	}
	net := len(plan.Create) - len(plan.Delete)
	switch {
	case len(fit.Warnings) > 0:
		return "", fmt.Errorf("a malformed budget row can hide the row of %s: %s", dir, fit.Warnings[0])
	case *expected.Maximum != fit.Limit:
		return "", fmt.Errorf("expected maximum %d, the canonical limit is %d", *expected.Maximum, fit.Limit)
	case *expected.Current != fit.Current:
		return "", fmt.Errorf("expected current %d, %s holds %d source file(s) at the source tip", *expected.Current, dir, fit.Current)
	case *expected.Net != net:
		return "", fmt.Errorf("expected net %d, the plan nets %d", *expected.Net, net)
	case !fit.Fits(net):
		return "", fmt.Errorf("the plan of net %d exceeds the limit %d of %s, which no grant exempts", net, fit.Limit, dir)
	}
	reason := fmt.Sprintf("%s holds %d source file(s) at the source tip, and net %d fits limit %d", dir, fit.Current, net, fit.Limit)
	if fit.Granted {
		reason += " under the accept grant: " + fit.Grant
	}
	return reason, nil
}

// linesSubject is the decoded subject of one lines premise: one file path and its source.
type linesSubject struct {
	Path   string `json:"path"`
	Source string `json:"source"`
}

// decodeLinesSubject decodes the subject of one lines premise.
func decodeLinesSubject(item PremiseItem) (linesSubject, error) {
	var subject linesSubject
	_, err := typed(item.Subject, "subject: ", &subject, "path", "source")
	return subject, err
}

// globSubject is the decoded subject of one glob-members premise: a pattern, the directory
// scope that holds every match, and its source.
type globSubject struct {
	Pattern string `json:"pattern"`
	Scope   string `json:"scope"`
	Source  string `json:"source"`
}

// decodeGlobSubject decodes the subject of one glob-members premise.
func decodeGlobSubject(item PremiseItem) (globSubject, error) {
	var subject globSubject
	_, err := typed(item.Subject, "subject: ", &subject, "pattern", "scope", "source")
	return subject, err
}

// headroomSubject is the decoded subject of one headroom premise: a directory and its
// authored file plan.
type headroomSubject struct {
	Directory string          `json:"directory"`
	Files     json.RawMessage `json:"files"`
}

// decodeHeadroomSubject decodes the subject of one headroom premise.
func decodeHeadroomSubject(item PremiseItem) (headroomSubject, error) {
	var subject headroomSubject
	_, err := typed(item.Subject, "subject: ", &subject, "directory", "files")
	return subject, err
}

// planned refuses a file plan that the source tip cannot hold: each entry must be one
// source file that dir's crowding count holds, named once. A creation must be absent from
// the tip, and a deletion must be one exact tracked path there.
func planned(dir string, tracked, create, remove []string) error {
	var named []string
	for _, entry := range slices.Concat(create, remove) {
		in, source := structure.SourceDirectory(entry)
		switch {
		case !tickets.RepositoryPath(entry) || !source || in != dir:
			return fmt.Errorf("subject: files: %q is no source file of %s", entry, dir)
		case slices.Contains(named, entry):
			return fmt.Errorf("subject: files: %s is named more than once", entry)
		}
		named = append(named, entry)
	}
	for _, entry := range create {
		if slices.Contains(tracked, entry) {
			return fmt.Errorf("subject: files: creation %s is already tracked at the source tip", entry)
		}
	}
	for _, entry := range remove {
		if !slices.Contains(tracked, entry) {
			return fmt.Errorf("subject: files: deletion %s is no tracked source file of %s at the source tip", entry, dir)
		}
	}
	return nil
}

// limitsUnchanged refuses a checkout that changes, adds, or ignores a limit list of the
// structure engine relative to the source tip.
func (c *census) limitsUnchanged() error {
	status, err := git.Raw(append([]string{"-C", c.root, "status", "--porcelain", "-z", "--ignored", "--untracked-files=all", "--"}, structure.LimitSources()...)...)
	switch {
	case err != nil:
		return fmt.Errorf("the limit lists do not read: %w", err)
	case len(status) > 0:
		return fmt.Errorf("the checkout changes a limit list from the source tip: %s", strings.ReplaceAll(strings.TrimSuffix(string(status), "\x00"), "\x00", "; "))
	}
	return nil
}

// trackedAt is every tracked path at rev, dot paths included, read once per commit.
func (c *census) trackedAt(rev string) ([]string, error) {
	if paths, ok := c.tracked[rev]; ok {
		return paths, nil
	}
	listing, err := git.ListTree(c.root, "-r", "--full-tree", rev)
	if err != nil {
		return nil, fmt.Errorf("the tracked paths at %s do not read: %w", rev, err)
	}
	var paths []string
	for _, entry := range listing {
		paths = append(paths, entry.Path)
	}
	if c.tracked == nil {
		c.tracked = map[string][]string{}
	}
	c.tracked[rev] = paths
	return paths, nil
}

// pathGaps names each actual member that listed omits and each listed path that is no
// member. Both sides are exact paths.
func pathGaps(listed, actual []string) error {
	var gaps []string
	for _, member := range actual {
		if !slices.Contains(listed, member) {
			gaps = append(gaps, "omits "+member)
		}
	}
	for _, path := range listed {
		if !slices.Contains(actual, path) {
			gaps = append(gaps, "names "+path+", which is no member")
		}
	}
	if len(gaps) > 0 {
		return fmt.Errorf("the member list %s", strings.Join(gaps, "; "))
	}
	return nil
}
