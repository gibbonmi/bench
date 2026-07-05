package surface

import (
	"github.com/gibbonmi/bench/internal/contract"
	"strings"
	"testing"
)

func TestLinkMarkerFenceContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench link malformed marker contract failed", testLinkMalformedMarker)
	contract.RunParallel(t, "bench link fenced-marker contract failed", testLinkFencedMarker)
	contract.RunParallel(t, "bench link unclosed-fence contract failed", testLinkUnclosedFence)
}

func testLinkMalformedMarker(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("AGENTS.md", "PROJECT BEFORE\n<!-- bench:end -->\nPROJECT MIDDLE\n<!-- bench:start -->\nPROJECT AFTER\n")

	probe := f.Bench("link")

	if probe.ExitCode == 0 {
		t.Fatal("link succeeded despite reversed Bench managed block markers")
	}
	probe.RequireContains(strings.ToLower(probe.Stderr+probe.Stdout), "malformed")
	requireFixtureFileContains(t, f, "AGENTS.md", "PROJECT AFTER", "malformed marker failure still rewrote project-owned text")
}

func testLinkFencedMarker(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("AGENTS.md", "# Project rules\n\nHow Bench marks its block:\n\n```\n<!-- bench:start -->\nmanaged content example\n<!-- bench:end -->\n```\n\nKEEP-ME project text.\n")

	linkOK(t, f)

	requireFixtureFileContains(t, f, "AGENTS.md", "KEEP-ME project text.", "fenced-marker link lost project text")
	requireFixtureFileContains(t, f, "AGENTS.md", "managed content example", "fenced example content was rewritten")
	linkOK(t, f)
	requireFixtureFileContains(t, f, "AGENTS.md", "managed content example", "relink consumed the fenced example")
	requireLiteralCount(t, f, "AGENTS.md", "## Bench", 1, "fenced markers caused duplicate managed blocks")
}

func testLinkUnclosedFence(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("AGENTS.md", "# Project rules\n\nBroken docs with an unclosed fence:\n\n```\n<!-- bench:start -->\n<!-- bench:end -->\n\nKEEP-ME text after the unclosed fence.\n")

	probe := f.Bench("link")

	if probe.ExitCode == 0 {
		t.Fatal("link succeeded despite an unclosed fence around Bench markers")
	}
	probe.RequireContains(strings.ToLower(probe.Stderr+probe.Stdout), "fence")
	requireFixtureFileContains(t, f, "AGENTS.md", "KEEP-ME text after the unclosed fence.", "unclosed-fence failure rewrote project text")
	requireLiteralCount(t, f, "AGENTS.md", "## Bench", 0, "unclosed-fence link still installed a managed block")
}
