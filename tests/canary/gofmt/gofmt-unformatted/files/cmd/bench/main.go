package main

// Shipped formatted, and mutated into misformatted at materialization time: an
// unformatted .go file checked into the kit tree would be named by the kit's own gofmt
// phase, so the offence has to be introduced after the fixture leaves the repo.
func main() {
	x := 1
	_ = x
}
