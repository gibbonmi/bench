package gittest

// IdentityName and IdentityEmail are the canonical default fixture identity: the author
// and committer of a fixture commit whose test does not grade the identity. This file is
// the one source of each value. A generated fixture script reads them from here.
const (
	IdentityName  = "bench test"
	IdentityEmail = "bench@example.invalid"
)

// IdentityConfig returns the git config commands that make the default identity the
// repository identity. A repository that runs each command commits as that identity
// until a test sets its own.
func IdentityConfig() [][]string {
	return [][]string{{"config", "user.name", IdentityName}, {"config", "user.email", IdentityEmail}}
}

// IdentityArgs returns the git options that give one command the default identity. They
// precede the subcommand: gittest.Run(root, append(gittest.IdentityArgs(), "commit", ...)...).
func IdentityArgs() []string {
	return []string{"-c", "user.name=" + IdentityName, "-c", "user.email=" + IdentityEmail}
}
