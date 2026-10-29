package capability

// Unavailable skips t for a host that cannot perform the operation a security
// assertion in class needs. The caller passes only the operation detail, for example
// the underlying error or the host limitation. Unavailable derives the one reason
// prefix from the class token. Capability writes the record and then
// ends the test, so the class, the line shape, and the strict-mode count stay the
// low-level helper's.
func Unavailable(t TB, class Class, detail string) {
	t.Helper()
	Capability(t, class, string(class)+" capability unavailable: "+detail)
}
