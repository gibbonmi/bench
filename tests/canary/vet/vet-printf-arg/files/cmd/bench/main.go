package main

import "fmt"

// canaryVetSubject is the fixture's own name for the mistyped argument, and vet names
// the offending expression back in its diagnostic. That is what the EXPECT matches: a
// quote of vet's own phrasing would tie the fixture to a Go release's wording, where a
// rewording reports "did not bite" about the toolchain rather than about the phase.
var canaryVetSubject = "not a number"

// A printf verb that disagrees with its argument: this compiles and gofmt accepts it,
// so `go vet` is the only step that can object. If the vet phase rots into an
// always-pass, this stops biting.
func main() {
	fmt.Printf("%d\n", canaryVetSubject)
}
