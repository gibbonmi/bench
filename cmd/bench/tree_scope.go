package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/treetarget"
	"github.com/gibbonmi/bench/internal/worktree"
)

// treeScope states which tree a public command leaf serves. Each public definition that is
// not a family declares one, and so does each leaf row of a family. A family and a plumbing
// definition declare none. The zero value means undeclared: such a definition prints no
// tree row and takes no tree target, so a planted test registry keeps its output.
type treeScope uint8

const (
	// scopeTree marks a leaf that reads, grades, or writes tracked content of one checkout.
	scopeTree treeScope = iota + 1
	// scopeRepository marks a leaf that takes no tree target: it reads or writes only Git
	// refs, the worktree ledger, git-ignored capture files, or state under the Bench home,
	// or, as root read does, it reads only the checkout of the caller or an absolute
	// artifact path.
	scopeRepository
)

// treeTargetFlag is the one spelling of the tree-target flag.
const treeTargetFlag = "--in"

// leadsWithTreeTarget reports whether args name a tree target. The flag counts only as the
// first argument after the verb, or after the leaf name of a family leaf. In any other
// position the verb's own grammar reads it.
func leadsWithTreeTarget(args []string) bool {
	return len(args) > 0 && args[0] == treeTargetFlag
}

// leafCall splits a call that names a family leaf into the leaf name and the arguments
// after it. Every other call answers no leaf and args whole.
func (definition commandDefinition) leafCall(args []string) (string, []string) {
	if leaf, ok := leafNamed(definition.Leaves, args); ok {
		return leaf.Name, args[1:]
	}
	return "", args
}

// refusesTreeTarget answers the repository refusal. A repository-scoped call that gets the
// tree-target flag as its first argument prints the usage line of its definition, or of its
// declared child, and does not run, so an ignored target never looks accepted. A family
// declares no scope, so its first argument still routes as a leaf name.
func (definition commandDefinition) refusesTreeTarget(c Command, args []string) bool {
	if definition.scope(args) != scopeRepository || !leadsWithTreeTarget(args) {
		return false
	}
	name := "bench " + definition.Name
	if child := formArgument(args); definition.ChildScopes[child] != 0 {
		name += " " + child
	}
	fmt.Fprintln(c.Stdout, toon.Usage(name, treeTargetFlag))
	return true
}

// runInTreeTarget runs a tree-scoped call that names its tree target as one child in that
// tree, and answers the child's exit and true. A tree-scoped family leaf takes the target
// after its leaf name. Every other call answers false and runs in this process, so a
// repository-scoped child of a tree-scoped definition reaches the repository refusal. The
// running executable goes by its absolute path, because the child directory is another
// tree, where the spelling of argv[0] can name another file. When that path is unknown it
// stays empty, and a child start that needs it fails with its own error.
func (c Command) runInTreeTarget(definition commandDefinition, args []string) (int, bool) {
	leaf, rest := definition.leafCall(args)
	if definition.scope(args) != scopeTree || !leadsWithTreeTarget(rest) {
		return 0, false
	}
	running, _ := os.Executable()
	return treetarget.Run(treetarget.Call{
		Name: definition.Name, Leaf: leaf, Flag: treeTargetFlag, Args: rest[1:],
		Root: boundaryRoot(), Home: worktree.Home(),
		Wrapper: os.Getenv(env.WrapperEnv), Running: running,
		Stdin: c.Stdin, Stdout: c.Stdout, Stderr: c.Stderr,
	}), true
}

// treeTargetInsertion answers the root help text that each row of definition carries after
// the verb name: the tree-target form for a tree-scoped definition, and nothing otherwise.
// The help and the dispatcher read the same scope field.
func (definition commandDefinition) treeTargetInsertion() string {
	if definition.Scope != scopeTree {
		return ""
	}
	return " [" + treeTargetFlag + " " + treetarget.Operand + "]"
}

// rowInsertion answers the tree-target insertion of one root help row of definition. The row
// takes the scope of the child form that its suffix begins with, so the help, the identity
// row, and the refusal read the one ChildScopes declaration.
func (definition commandDefinition) rowInsertion(row helpRow) string {
	child, _, _ := strings.Cut(strings.TrimPrefix(row.Suffix, " "), " ")
	if definition.childScope(child) != scopeTree {
		return ""
	}
	return definition.treeTargetInsertion()
}

// runNamingTarget runs one call of definition. A help form that answers at exit 0 then
// prints the tree-target line.
func (definition commandDefinition) runNamingTarget(c Command, args []string) int {
	exit := definition.run(c, args)
	if exit == 0 && definition.helpForm(args) {
		fmt.Fprint(c.Stdout, definition.treeTargetLine())
	}
	return exit
}

// treeTargetLine answers the line that ends the help of definition: the tree-target form for
// a tree-scoped definition, and nothing otherwise. The verb help and the root help read one
// insertion.
func (definition commandDefinition) treeTargetLine() string {
	if insertion := definition.treeTargetInsertion(); insertion != "" {
		return "tree target: bench " + definition.Name + insertion + "\n"
	}
	return ""
}

// scope answers the scope of one call. A call that names a family leaf takes the scope of
// that leaf, and every other call takes the scope of the child form that it names.
func (definition commandDefinition) scope(args []string) treeScope {
	if leaf, ok := leafNamed(definition.Leaves, args); ok {
		return leaf.Scope
	}
	return definition.childScope(formArgument(args))
}

// childScope answers the scope of one child form: the scope that the definition declares
// for that child, or else the scope of the definition.
func (definition commandDefinition) childScope(child string) treeScope {
	if scope, ok := definition.ChildScopes[child]; ok {
		return scope
	}
	return definition.Scope
}

// treeRow answers the identity row block of one call, or the empty string for a call that
// prints no row. Only a tree-scoped call that is not a help form and that runs inside a
// repository prints the row. The dispatcher computes it before the verb runs, so the row
// names the tree that the verb read. The renderer strips the control bytes of a hostile
// label, so only an encoder refusal that the row cells cannot produce drops the row.
func (definition commandDefinition) treeRow(args []string) string {
	if definition.scope(args) != scopeTree || definition.helpForm(args) {
		return ""
	}
	root := boundaryRoot()
	if root == "" {
		return ""
	}
	return identityRow(treetarget.Identify(root))
}

// identityRow answers the identity row block of one tree. The dispatcher prints it, and a
// tree-scoped reader counts it in its page budget.
func identityRow(identity treetarget.Identity) string {
	row, err := treetarget.Row(identity)
	if err != nil {
		return ""
	}
	return row
}

// shownRow answers the row that a call with exit prints. A grammar refusal at exit 2
// prints only its usage line, so it prints no row.
func shownRow(row string, exit int) string {
	if exit == 2 {
		return ""
	}
	return row
}

// finishExempt prints the row of an exempt call on stderr after the verb returns, so an
// artifact on stdout stays byte-clean. It answers the verb's exit.
func (c Command) finishExempt(row string, exit int) int {
	fmt.Fprint(c.Stderr, shownRow(row, exit))
	return exit
}
