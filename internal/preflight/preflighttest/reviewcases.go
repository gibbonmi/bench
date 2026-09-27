package preflighttest

import (
	"fmt"
	"strings"
)

// This file holds the pinned review pairs and their expected file patches, stated by hand
// from the pair trees and Git's path order.

// renameEditBody is the ten-line body the edited rename keeps nine lines of.
var renameEditBody = "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n"

// hostilePaths are the added paths whose bytes need quoting, escaping, or pathspec care:
// spaces, glob characters, a quote, a backslash, non-ASCII bytes, and the tab, newline, and
// return bytes the current renderer permits.
var hostilePaths = []string{
	"notes/space name.txt", "notes/glob*[?].txt", "notes/quote\"d.txt", "notes/back\\slash.txt",
	"notes/café π.txt", "notes/tab\there.txt", "notes/new\nline.txt", "notes/carriage\rreturn.txt",
}

// ReviewShapes is the RE5 edge inventory as one pinned pair: added, modified, deleted, and
// renamed files, a pure rename, an empty addition and deletion, a mode-only change, a binary
// change, a committed symlink, a file-to-symlink change, a missing final newline, content
// that resembles patch headers, and hostile paths. renames sets diff.renames for the pair.
func ReviewShapes(renames bool) ReviewPair {
	pair := ReviewPair{
		Base: map[string]TreeEntry{
			"notes/modify.txt":       Regular("one\ntwo\nthree\n"),
			"notes/delete.txt":       Regular("delete me\n"),
			"notes/rename-pure.txt":  Regular("pure rename body\nline two\nline three\n"),
			"notes/rename-edit.txt":  Regular(renameEditBody),
			"notes/empty-delete.txt": Regular(""),
			"notes/mode.sh":          Regular("#!/bin/sh\n"),
			"notes/binary.bin":       Regular("\x00\x01binary"),
			"notes/type-change":      Regular("regular\n"),
			"notes/no-newline.txt":   Regular("last"),
			"notes/lookalike.txt":    Regular("a\n"),
		},
		Removed: []string{
			"notes/delete.txt", "notes/rename-pure.txt", "notes/rename-edit.txt", "notes/empty-delete.txt",
		},
		Tip: map[string]TreeEntry{
			"notes/modify.txt":       Regular("one\ntwo changed\nthree\n"),
			"notes/renamed-pure.txt": Regular("pure rename body\nline two\nline three\n"),
			"notes/renamed-edit.txt": Regular(strings.Replace(renameEditBody, "five", "FIVE", 1)),
			"notes/mode.sh":          Executable("#!/bin/sh\n"),
			"notes/binary.bin":       Regular("\x00\x02binary"),
			"notes/type-change":      Symlink("modify.txt"),
			"notes/no-newline.txt":   Regular("last line changed"),
			"notes/lookalike.txt":    Regular("a\ndiff --git a/x b/x\n--- a/x\n+++ b/x\nrename to y\n@@ -1 +1 @@\n"),
			"notes/add.txt":          Regular("added\n"),
			"notes/empty-add.txt":    Regular(""),
			"notes/link":             Symlink("add.txt"),
		},
		Config: [][2]string{{"diff.renames", fmt.Sprint(renames)}},
	}
	for _, path := range hostilePaths {
		pair.Tip[path] = Regular("hostile\n")
	}
	return pair
}

// ShapePatch is one expected file identity: its declared path, the pathspecs in patch order,
// and its number of patch sources. Git's raw patch over the pathspecs equals the patches of
// each pathspec path joined in that order.
type ShapePatch struct {
	Path      string
	Pathspecs []string
	Patches   int
}

// OnePatch is an expected identity whose raw patch is the patch of its own path.
func OnePatch(path string) ShapePatch { return ShapePatch{path, []string{path}, 1} }

// RenamedPatch is an expected rename identity: the tip path, found by both spellings.
func RenamedPatch(from, to string) ShapePatch { return ShapePatch{to, []string{from, to}, 1} }

// Selected joins the patch bytes of each pathspec path in order, from bodies keyed by
// declared path, so it equals Git's raw patch over the pathspecs.
func (p ShapePatch) Selected(bodies map[string]string) string {
	var joined strings.Builder
	for _, path := range p.Pathspecs {
		joined.WriteString(bodies[path])
	}
	return joined.String()
}

// ShapePaths is the declared path of every expected patch source, in patch order.
func ShapePaths(patches []ShapePatch) []string {
	var paths []string
	for _, patch := range patches {
		for range patch.Patches {
			paths = append(paths, patch.Path)
		}
	}
	return paths
}

// The diff fragment roles, stated from the spec independently of the review collector, so
// a changed role in the collector turns every reader of these names red.
const (
	DiffPrefixRole = "diff-prefix"
	DiffFileRole   = "diff-file"
	DiffSuffixRole = "diff-suffix"
)

// typeChangePatches is the file-to-symlink identity: two patches for one path.
var typeChangePatches = ShapePatch{"notes/type-change", []string{"notes/type-change"}, 2}

// shapesWithRenames is the patch order of ReviewShapes with rename detection enabled.
var shapesWithRenames = []ShapePatch{
	OnePatch("notes/add.txt"), OnePatch("notes/back\\slash.txt"), OnePatch("notes/binary.bin"), OnePatch("notes/café π.txt"),
	OnePatch("notes/carriage\rreturn.txt"), OnePatch("notes/delete.txt"),
	RenamedPatch("notes/empty-delete.txt", "notes/empty-add.txt"), OnePatch("notes/glob*[?].txt"), OnePatch("notes/link"),
	OnePatch("notes/lookalike.txt"), OnePatch("notes/mode.sh"), OnePatch("notes/modify.txt"), OnePatch("notes/new\nline.txt"),
	OnePatch("notes/no-newline.txt"), OnePatch("notes/quote\"d.txt"),
	RenamedPatch("notes/rename-edit.txt", "notes/renamed-edit.txt"),
	RenamedPatch("notes/rename-pure.txt", "notes/renamed-pure.txt"), OnePatch("notes/space name.txt"),
	OnePatch("notes/tab\there.txt"), typeChangePatches,
}

// shapesWithoutRenames is the patch order of ReviewShapes with rename detection disabled.
var shapesWithoutRenames = []ShapePatch{
	OnePatch("notes/add.txt"), OnePatch("notes/back\\slash.txt"), OnePatch("notes/binary.bin"), OnePatch("notes/café π.txt"),
	OnePatch("notes/carriage\rreturn.txt"), OnePatch("notes/delete.txt"), OnePatch("notes/empty-add.txt"),
	OnePatch("notes/empty-delete.txt"), OnePatch("notes/glob*[?].txt"), OnePatch("notes/link"), OnePatch("notes/lookalike.txt"),
	OnePatch("notes/mode.sh"), OnePatch("notes/modify.txt"), OnePatch("notes/new\nline.txt"), OnePatch("notes/no-newline.txt"),
	OnePatch("notes/quote\"d.txt"), OnePatch("notes/rename-edit.txt"), OnePatch("notes/rename-pure.txt"),
	OnePatch("notes/renamed-edit.txt"), OnePatch("notes/renamed-pure.txt"), OnePatch("notes/space name.txt"),
	OnePatch("notes/tab\there.txt"), typeChangePatches,
}

// ReviewShapePatches is the expected patch table of ReviewShapes under one rename setting,
// stated by hand from the pair's trees and Git's path order.
func ReviewShapePatches(renames bool) []ShapePatch {
	if renames {
		return shapesWithRenames
	}
	return shapesWithoutRenames
}

// reviewDocuments is a pinned pair whose diff holds only the spec and the review record.
func reviewDocuments() ReviewPair {
	return ReviewPair{Tip: map[string]TreeEntry{
		pairSpecPath:         Regular(SpecBody(pairSlug, reviewFenceLines()...) + "\n## Further notes\n\nEdited.\n"),
		"reviews/example.md": Regular("# Review record\n"),
	}}
}

// copyBody is the five-line body the copy pair edits in its source and copies unchanged.
var copyBody = "one\ntwo\nthree\nfour\nfive\n"

// ReviewCase is one pinned review pair and its expected file patches, stated by hand from
// the pair's trees and Git's path order. Stored reports that a stored monolithic response
// named Name exists; a case without one compares its fragments with the current command
// response, and each path with Git's own patch.
type ReviewCase struct {
	Name    string
	Stored  bool
	Pair    ReviewPair
	Patches []ShapePatch
}

// ReviewCases is the one table of pinned review pairs that both the diff seam and the
// evidence seam rebuild. The copy pair holds `copy from` and `copy to` headers, the blank
// pair holds a bare blank context line, and the trailing-space pair holds side headers that
// end with a space before Git's tab.
func ReviewCases() []ReviewCase {
	return []ReviewCase{
		{"review-empty", true, ReviewPair{}, nil},
		{"review-documents", true, reviewDocuments(), []ShapePatch{OnePatch("reviews/example.md"), OnePatch("specs/example/spec.md")}},
		{"review-shapes-renames", true, ReviewShapes(true), ReviewShapePatches(true)},
		{"review-shapes-no-renames", true, ReviewShapes(false), ReviewShapePatches(false)},
		{"review-copies", false, ReviewPair{
			Base:   map[string]TreeEntry{"notes/copy-source.txt": Regular(copyBody)},
			Tip:    map[string]TreeEntry{"notes/copy-source.txt": Regular(strings.Replace(copyBody, "three", "THREE", 1)), "notes/copy-target.txt": Regular(copyBody)},
			Config: [][2]string{{"diff.renames", "copies"}},
		}, []ShapePatch{
			OnePatch("notes/copy-source.txt"),
			{"notes/copy-target.txt", []string{"notes/copy-source.txt", "notes/copy-target.txt"}, 1},
		}},
		{"review-blank-suppressed", false, ReviewPair{
			Base:   map[string]TreeEntry{"notes/blank.txt": Regular("a\n\nb\n")},
			Tip:    map[string]TreeEntry{"notes/blank.txt": Regular("a\n\nc\n")},
			Config: [][2]string{{"diff.suppressBlankEmpty", "true"}},
		}, []ShapePatch{OnePatch("notes/blank.txt")}},
		{"review-trailing-space", false, ReviewPair{
			Base:    map[string]TreeEntry{"notes/gone ": Regular("removed body\n"), "notes/trail ": Regular("old\n")},
			Removed: []string{"notes/gone "},
			Tip:     map[string]TreeEntry{"notes/added ": Regular("fresh text\n"), "notes/trail ": Regular("new\n")},
		}, []ShapePatch{OnePatch("notes/added "), OnePatch("notes/gone "), OnePatch("notes/trail ")}},
	}
}
