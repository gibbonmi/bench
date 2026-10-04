package landing

import "github.com/gibbonmi/bench/internal/spec"

// TreeReader retains the landing API for the shared spec-folder classifier.
type TreeReader = spec.TreeReader

func TicketsOnly(reader TreeReader, name string) bool  { return spec.TicketsOnly(reader, name) }
func WorkTree(root string) TreeReader                  { return spec.WorkTree(root) }
func CommitTree(root, commit string) TreeReader        { return spec.CommitTree(root, commit) }
func TicketsOnlyFolder(root, name string) bool         { return spec.TicketsOnlyFolder(root, name) }
func TicketsOnlyFolderPath(root, name string) string   { return spec.TicketsOnlyFolderPath(root, name) }
func TicketsOnlyFolders(root string) ([]string, error) { return spec.TicketsOnlyFolders(root) }
func ClosedFolderPath(name string) string              { return spec.ClosedFolderPath(name) }
