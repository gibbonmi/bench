package tickets

import "sort"

// BindingRow binds one package prefix to the files a ticket must co-name when
// it writes into that package. The prefix matches a `Writes:` path exactly, or
// at a `/` segment boundary, so `internal/toon` never claims `internal/toon2`.
type BindingRow struct {
	// Prefix is the bound package, repo-relative and slash-spelled.
	Prefix string
	// Files are the registries the ticket must name beside its own edit.
	Files []string
}

// commandRegistries is the file set every command-package ticket co-names: the
// help projection and its assertion, the envelope cases, the approved AXI query
// list, and the subcommand routing census. A verb that ships without one of
// these reaches the gate with a registry nobody updated.
var commandRegistries = []string{
	"cmd/bench/command_registry.go",
	"cmd/bench/command_registry_test.go",
	"cmd/bench/main_test.go",
	"internal/conformance/axi_query_registry_test.go",
	"internal/conformance/subcommand_routing_test.go",
}

// bindings is the ordered binding registry, sorted by prefix. The command rows
// carry the shared registry set. The three owner rows below them are the
// dispatcher, the renderer, and the terminal-lifecycle owner: established owners
// that recent builds reconstructed by hand.
var bindings = []BindingRow{
	{Prefix: "cmd/bench", Files: commandRegistries},
	{Prefix: "internal/anchors", Files: commandRegistries},
	{Prefix: "internal/consumers", Files: commandRegistries},
	{Prefix: "internal/coverage", Files: commandRegistries},
	{Prefix: "internal/diff", Files: commandRegistries},
	{Prefix: "internal/guards", Files: commandRegistries},
	{Prefix: "internal/harnesses", Files: commandRegistries},
	{Prefix: "internal/learnings", Files: commandRegistries},
	{Prefix: "internal/maps", Files: commandRegistries},
	{Prefix: "internal/roadmap", Files: commandRegistries},
	{Prefix: "internal/terminal", Files: []string{
		"internal/adopt/setup_prompt_test.go",
		"internal/adopt/setup_test.go",
		"internal/gate/command_test.go",
	}},
	{Prefix: "internal/toon", Files: []string{
		"internal/conformance/data_handling_test.go",
		"internal/toon/toon_test.go",
	}},
	{Prefix: "internal/worktree", Files: commandRegistries},
}

// Bindings returns the ordered binding registry.
func Bindings() []BindingRow {
	return append([]BindingRow(nil), bindings...)
}

// BoundFiles returns every file bound to a row whose prefix covers path, sorted
// and deduplicated. A path under no bound package answers none.
func BoundFiles(path string) []string {
	var files []string
	for _, row := range bindings {
		if !prefixCovers(row.Prefix, path) {
			continue
		}
		for _, file := range row.Files {
			if file != path && !holdsString(files, file) {
				files = append(files, file)
			}
		}
	}
	sort.Strings(files)
	return files
}

// prefixCovers reports whether prefix names path itself or a path under it. The
// comparison is at a `/` segment boundary, never a bare string prefix.
func prefixCovers(prefix, path string) bool {
	return path == prefix || (len(path) > len(prefix)+1 && path[:len(prefix)] == prefix && path[len(prefix)] == '/')
}

func holdsString(list []string, want string) bool {
	for _, value := range list {
		if value == want {
			return true
		}
	}
	return false
}
