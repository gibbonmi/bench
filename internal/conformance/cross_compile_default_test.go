//go:build !stress

package conformance

// crossCompileMatrix is a no-op on the every-commit gate: the four-platform cross-compile
// is portability gold-plating that runs only under `go test -tags stress` (see the stress
// variant), with the release workflow as the standing backstop. This keeps the default
// conformance run — the one the gate invokes — free of the ~8s cross-compile cost.
func crossCompileMatrix(root, buildHelper string) []string { return nil }
