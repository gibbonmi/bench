package runtime

import (
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestRuntimeOutlineContracts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	contract.RunParallel(t, "hostile filenames scope literally", testOutlineHostileFilenames)
}

// testOutlineHostileFilenames pins outline's path argument to a literal name against the two
// hostile filename shapes a shell CLI meets: an embedded space, which word-splits if any
// layer between argv and the git pathspec drops its quoting, and a `*`, which selects other
// files if the pathspec loses its `:(literal,top)` prefix. Scoping to the glob-shaped name is
// the discriminator — as a glob that name also matches the space-shaped file, so a
// non-literal chain widens the scope past what the agent asked for.
func testOutlineHostileFilenames(t *testing.T) {
	f := contract.NewFixture(t)
	f.Git("config", "user.email", "bench@local")
	f.Git("config", "user.name", "bench")
	space, glob := "a b.md", "a*b.md"
	f.WriteFile(space, "# Space\n")
	f.WriteFile(glob, "# Glob\n")
	f.CommitAll("hostile names")

	all := f.Bench("outline")
	all.RequireExit(0)
	contract.RequireContains(t, all.Stdout, space)
	contract.RequireContains(t, all.Stdout, glob)

	scoped := f.Bench("outline", glob)
	scoped.RequireExit(0)
	contract.RequireContains(t, scoped.Stdout, glob)
	contract.RequireNotContains(t, scoped.Stdout, space)
}
