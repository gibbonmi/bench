package main

// A main package with no func main: `go build` refuses it, while gofmt, `go vet`, and
// `go test` all accept it. That asymmetry is deliberate — it is what makes this
// fixture's EXPECT reachable only from the build phase, so a tree that loses the build
// probe's inputs and falls back to the full inner gate stops biting instead of matching
// from a neighbouring step.
func run() {}

var _ = run
