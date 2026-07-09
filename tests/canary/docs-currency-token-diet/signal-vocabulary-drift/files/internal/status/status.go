//go:build ignore

// This is a canary fixture, not a compiled package: the build constraint keeps the
// parent module's `go vet ./...` / `go build ./...` from descending into it, while
// checkSignalVocabulary still reads it as text (readIfExists ignores build tags).
//
// Planted regression: status.go emits a `roadmap` signal, but CONTEXT.md's signal
// enumeration below names only `gate`. checkSignalVocabulary keys on the row{}
// literals here; the decoy prose in CONTEXT.md mentions "roadmap" outside the
// enumeration to prove the scoped match still fires where a whole-file Contains
// would false-pass.
package status

func rows() {
	_ = row{7, "gate", d, a}
	_ = row{10, "roadmap", d, a}
}
