package main

import "fmt"

// A printf verb that disagrees with its argument: this compiles and gofmt accepts it,
// so `go vet` is the only step that can object. If the vet phase rots into an
// always-pass, this stops biting.
func main() {
	fmt.Printf("%d\n", "not a number")
}
