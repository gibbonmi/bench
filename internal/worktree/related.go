package worktree

import (
	"bytes"
	"cmp"
	"crypto/sha256"
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/canonicalpath"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/tickets"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// relatedSpecFlag selects the related view: the assignments that the commitment relates to
// one spec.
const relatedSpecFlag = "--spec"

// relatedSchema names the related projection that each related page identity binds.
const relatedSchema = "bench-related/1"

var relatedFields = []string{"spec", "matched", "unbound"}

// relatedWorktreeGrammar is the related view's grammar. It declares neither selected-view
// flag, so the selected view and this view refuse each other's flags.
var relatedWorktreeGrammar = usage.Grammar{
	Cmd:  usage.WorktreeList,
	Help: "usage: " + usage.WorktreeListRelated,
	Flags: []usage.Flag{
		{Name: relatedSpecFlag, HasValue: true, NoEmptyValue: true, Required: true},
		{Name: responsebound.OffsetFlag, HasValue: true, NoEmptyValue: true},
		{Name: responsebound.SourceFlag, HasValue: true, NoEmptyValue: true},
	},
}

// relatesWorktrees reports whether args choose the related view.
func relatesWorktrees(args []string) bool {
	return usage.FlagPresent(relatedWorktreeGrammar, args, relatedSpecFlag)
}

// listRelatedWorktrees pages the related projection of one spec, whose path operand names it
// from cwd. The page identity binds
// the complete projection, so a continuation refuses before content when any row, count, or
// action changed since the cursor, even when the commitment data stayed equal. Each offered
// action still resolves and authorizes its target through its own verb.
func listRelatedWorktrees(j joins, run *git.Run, root, cwd string, args []string) (string, int) {
	parsed, line, code := usage.Parse(relatedWorktreeGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	if parsed.EndedFlags {
		return relatedWorktreeGrammar.Help + "\n", 2
	}
	offset, token, continued, err := responsebound.ParseCursor(parsed.Flags)
	if err != nil {
		return relatedWorktreeGrammar.Help + " (" + err.Error() + ")\n", 2
	}
	if !inRepository(root) {
		return toon.NotInRepo() + "\n", 1
	}
	operand := parsed.Flags[relatedSpecFlag]
	specPath, err := relatedSpec(root, cwd, operand)
	if err != nil {
		return toon.Errorf(err.Error(), "pass a spec slug or a path to a spec markdown file inside the repository") + "\n", 1
	}
	// One ledger read supplies both the assignments and the commitment state.
	ledger, err := intent.Read(root)
	if err != nil {
		return assignmentsReadRefusal(), 1
	}
	content, err := relatedProjection(j, run, root, ledger, specPath)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	page, err := responsebound.PageRead(responsebound.ReadRequest{
		Command: append(strings.Fields(usage.WorktreeList)[1:], relatedSpecFlag, operand),
		Kind:    "related", Selection: specPath, Projection: relatedSchema,
		Scope: canonicalRoot(root), Path: specPath,
		Digest: sha256.Sum256(content), SourceSize: int64(len(content)),
		Payload: bytes.NewReader(content), Size: int64(len(content)),
		Offset: offset, Token: token, Continued: continued,
	})
	if err != nil {
		return err.Error() + "\n", 1
	}
	return page.Output, 0
}

// relatedSpec answers the repository-relative spec path that operand names: a slug names its
// live folder spec, and a path names itself from cwd, the caller's directory. The spec need
// not exist, because the commitment bindings, not the file, carry the relationship.
func relatedSpec(root, cwd, operand string) (string, error) {
	specPath := spec.LiveSpecPath(operand)
	// LiveSpecPath keeps a path operand as written and moves a slug into the specs tree, so
	// an unchanged operand is a path.
	if specPath == filepath.ToSlash(filepath.Clean(operand)) {
		_, specPath = canonicalpath.Operand(root, cwd, operand)
	}
	if !filepath.IsLocal(specPath) || path.Base(specPath) != "spec.md" {
		return "", fmt.Errorf("%s is not a spec path inside the repository", specPath)
	}
	return specPath, nil
}

// relatedProjection renders the complete related projection of the spec at specPath: one
// summary row with the matched and unbound counts, the list row of each matched assignment
// in assignment-ID order, and their ordered navigation actions. An assignment matches when
// a scope that the commitment state gives it holds the spec folder or lies within it. An
// assignment with no scope is unbound: it is excluded and counted, because a missing binding
// proves no relationship either way.
func relatedProjection(j joins, run *git.Run, root string, ledger intent.Ledger, specPath string) ([]byte, error) {
	folder := path.Dir(specPath)
	var matched []intent.Assignment
	unbound := 0
	for _, assignment := range ledger.Assignments {
		scopes := commitrepo.RunScopes(ledger.Commitment, assignment.ID)
		switch {
		case len(scopes) == 0:
			unbound++
		case slices.ContainsFunc(scopes, func(scope string) bool { return tickets.Covers(scope, folder) || tickets.Covers(folder, scope) }):
			matched = append(matched, assignment)
		}
	}
	slices.SortFunc(matched, func(a, b intent.Assignment) int { return cmp.Compare(a.ID, b.ID) })
	summary, err := toon.TableTyped(relatedTable, relatedFields, [][]any{{specPath, len(matched), unbound}})
	if err != nil {
		return nil, err
	}
	def, defaultResolved := git.ResolvedDefault(root)
	rows, err := renderListRows(j, run, root, assignmentRows(j, run, root, matched, def, defaultResolved), matched, namesTreeAlone(ledger.Assignments))
	if err != nil {
		return nil, err
	}
	return []byte(summary + rows), nil
}

// namesTreeAlone answers the predicate that a related row's label passes when the read action
// can name that row's tree alone. The label is not the primary keyword, which the tree resolver
// reads before any label. No other active assignment of assignments carries it, because the
// resolver refuses an ambiguous label. And the action owner accepts the read action that
// carries it, so a rejected label drops only its own action from the help.
func namesTreeAlone(assignments []intent.Assignment) func(label string) bool {
	empty, _ := axi.RenderHelp(nil)
	return func(label string) bool {
		if label == usage.PrimaryTreeTarget {
			return false
		}
		if len(matchingAssignments(assignments, func(a intent.Assignment) bool { return a.Label == label && landingActiveState(a.State) })) != 1 {
			return false
		}
		help, err := axi.RenderHelp([]axi.Action{readAction(label)})
		return err == nil && help != empty
	}
}
