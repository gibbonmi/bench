package conformance

// Canary fixture: the shared classifier table, kept beside the drifted word test so the
// conformance check has a real declaration to disagree with. The xargs row still claims
// true, which the fixture's library no longer answers. The other rows stay to prove the
// check names the one drifted row rather than reporting the whole table.
//
// The row type is declared here rather than borrowed from the package this file shadows.
// checkGuardClassifierTable reads the table with go/parser and never compiles it, but
// `go vet ./...` compiles every Go file in the tree, overlays included, so an overlay
// carries its own declarations the way every other Go overlay under tests/canary does.

type guardClassifierRow struct {
	command string
	invokes bool
}

var guardClassifierTable = []guardClassifierRow{
	{"bench gate", true},
	{"env -u X bench help", true},
	{"timeout -s KILL -k 1 5 bench help", true},
	{"xargs -- bench help", true},
	{"command -v bench", false},
	{"ls", false},
}
