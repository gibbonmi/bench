package main

// A type error: valid syntax (so gofmt stays clean) that fails go build, so the
// canary attributes to the build check rather than to formatting. If the gate's
// go-build check rots into an always-pass, this fixture stops biting.
func main() {
	var n int = "not an int"
	_ = n
}
