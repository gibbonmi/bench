package testreport

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/tickets"
)

// Floor is the package part of one ticket's check floor: the packages that the ticket
// writes directly, without the reverse-import closure of `--changed`.
type Floor struct {
	// Packages are the directory-relative package patterns, in sorted order.
	Packages []string
	// Embed reports that an entry is an embed file or selects a package that declares
	// an embed pattern.
	Embed bool
}

// WritesFloor derives the floor of one ticket from its raw `Writes:` entries, `(new)`
// markers included. It reads the package index that `--changed` reads, so the two
// never disagree about which package holds a path. The error is the loader error.
func WritesFloor(ctx context.Context, root string, entries []string) (Floor, error) {
	return writesFloorWithLoader(ctx, root, entries, listCurrentPackages)
}

func writesFloorWithLoader(ctx context.Context, root string, entries []string, load func(context.Context, string) ([]listedPackage, error)) (Floor, error) {
	if info, err := os.Stat(filepath.Join(root, "go.mod")); err != nil || !info.Mode().IsRegular() {
		// A root without a Go module owes no package test, so the floor starts no `go list`.
		return Floor{}, nil
	}
	packages, err := load(ctx, root)
	if err != nil {
		return Floor{}, err
	}
	index, err := newPackageIndex(root, packages)
	if err != nil {
		return Floor{}, err
	}
	floor := Floor{}
	selected := map[string]bool{}
	add := func(importPath string) {
		selected[importPath] = true
		if len(index.byImport[importPath].EmbedPatterns) != 0 {
			floor.Embed = true
		}
	}
	everyPackage := false
	for _, entry := range entries {
		entryPath, _ := tickets.WritesPath(entry)
		if isGoMetadata(entryPath) {
			everyPackage = true
			continue
		}
		if strings.HasSuffix(entryPath, ".go") {
			if importPath, ok := index.byDirectory[path.Dir(entryPath)]; ok {
				add(importPath)
			}
		}
		// The index holds only the embed files that exist, so an absent `(new)` asset
		// adds no owner.
		if owner, ok := index.byEmbed[entryPath]; ok {
			add(owner)
			floor.Embed = true
		}
		for directory, importPath := range index.byDirectory {
			if tickets.Covers(entryPath, directory) {
				add(importPath)
			}
		}
	}
	for directory, importPath := range index.byDirectory {
		if selected[importPath] {
			floor.Packages = append(floor.Packages, packagePattern(root, directory))
		}
	}
	if everyPackage {
		floor.Packages = append(floor.Packages, currentPackagePattern)
	}
	sort.Strings(floor.Packages)
	return floor, nil
}
