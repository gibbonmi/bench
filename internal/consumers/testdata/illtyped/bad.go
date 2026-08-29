package illtyped

// Bad holds one deliberate type error, so the loader's fail-closed rim is observed
// against the real go tool rather than a hand-built error value.
func Bad() {
	var n int = "not an int"
	_ = n
}
