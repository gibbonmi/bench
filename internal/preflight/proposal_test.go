package preflight

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

const proposalRegistryCitation = "internal/tickets/registry_data.go"

func proposalArgs(t *testing.T, root, slug string) []string {
	t.Helper()
	preflighttest.ActiveAssignment(t, root, root)
	return []string{"build", slug, "--propose-writes", "--ticket", "one.md", "--base", preflighttest.RunGit(t, "rev-parse", "main"), "--source-tip", preflighttest.RunGit(t, "rev-parse", "HEAD")}
}

func seedProposal(t *testing.T) (string, string) {
	return seedProposalWith(t,
		"internal/example/pinned.go, internal/toon/toon.go",
		map[string]string{"two.md": "tests/canary/example-family/pinning-fixture/EXPECT"})
}

// seedProposalWith fences every entry the tickets write, so only extraFence can leave the
// fence and the ticket union apart. The selected ticket also writes the conformant fence.
func seedProposalWith(t *testing.T, selectedWrites string, others map[string]string, extraFence ...string) (string, string) {
	t.Helper()
	root := preflighttest.StartRepo(t)
	slug := "example"
	fence := append([]string{}, extraFence...)
	for _, writes := range append([]string{selectedWrites}, mapValues(others)...) {
		for _, entry := range strings.Split(writes, ", ") {
			fence = append(fence, "- `"+entry+"`")
		}
	}
	preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", preflighttest.SpecBody(slug, fence...))
	writeProposalTicket(t, slug, "one.md", "One", "none", append(preflighttest.FenceWrites(preflighttest.ConformantFence), strings.Split(selectedWrites, ", ")...)...)
	for name, writes := range others {
		writeProposalTicket(t, slug, name, strings.TrimSuffix(name, ".md"), "none", strings.Split(writes, ", ")...)
	}
	preflighttest.MustWriteFile(t, "tests/canary/example-family/pinning-fixture/BASE", "internal/example/pinned.go\n")
	preflighttest.MustWriteFile(t, "tests/canary/example-family/pinning-fixture/EXPECT", "fixture\n")
	preflighttest.MustWriteFile(t, "internal/example/pinned.go", "package example\n")
	preflighttest.MustWriteFile(t, "internal/toon/toon.go", "package toon\n")
	preflighttest.MustWriteFile(t, "internal/toon/other.go", "package toon\n")
	preflighttest.MustWriteFile(t, "internal/toon/toon_test.go", "package toon\n")
	preflighttest.MustWriteFile(t, "internal/conformance/data_handling_test.go", "package conformance\n")
	preflighttest.MustWriteFile(t, chargesource.DelegateSkill, "# Delegation skill\n")
	preflighttest.MustWriteFile(t, chargesource.DelegateProcedure, "# Delegation procedure\n")
	preflighttest.MustWriteFile(t, chargesource.BuildPhase, "# Build phase\n")
	preflighttest.RunGit(t, "add", ".")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c0")
	preflighttest.RunGit(t, "checkout", "-q", "-b", "feature")
	preflighttest.MustWriteFile(t, "internal/example/foo.go", "package example\n")
	preflighttest.RunGit(t, "add", "internal/example/foo.go")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c1")
	return root, slug
}

func writeProposalTicket(t *testing.T, slug, name, title, blockers string, writes ...string) {
	t.Helper()
	body := strings.Replace(preflighttest.WritesTicketDoc(title, writes, "PF1", "PF2"), "Blocked by: none", "Blocked by: "+blockers, 1)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/"+name, body)
}

func mapValues(m map[string]string) []string {
	values := make([]string, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	return values
}

func assertProposalRow(t *testing.T, out, path, source, fence string) {
	t.Helper()
	want := "\n  " + path + "," + source + "," + fence + "\n"
	if strings.Count(out, want) != 1 {
		t.Errorf("proposal row count for %s = %d, want 1:\n%s", path, strings.Count(out, want), out)
	}
}

// TestWritesProposalClosure covers DP10 through the public command.
func TestWritesProposalClosure(t *testing.T) {
	cases := []struct {
		name, writes string
		fence        []string
		want         [][3]string
	}{
		{"fixture-only-covered", "internal/example/pinned.go", []string{"- `tests/canary/example-family/pinning-fixture/`"}, [][3]string{{"tests/canary/example-family/pinning-fixture", "fixture tests/canary/example-family/pinning-fixture", "covered"}}},
		{"registry-only-expansion", "internal/toon/toon.go", nil, [][3]string{{"internal/conformance/data_handling_test.go", "registry " + proposalRegistryCitation, "spec fence expansion required"}, {"internal/toon/toon_test.go", "registry " + proposalRegistryCitation, "spec fence expansion required"}}},
		{"fixture-and-registry", "internal/example/pinned.go, internal/toon/toon.go", nil, [][3]string{{"tests/canary/example-family/pinning-fixture", "fixture tests/canary/example-family/pinning-fixture", "spec fence expansion required"}, {"internal/conformance/data_handling_test.go", "registry " + proposalRegistryCitation, "spec fence expansion required"}, {"internal/toon/toon_test.go", "registry " + proposalRegistryCitation, "spec fence expansion required"}}},
		{"duplicate-registry-triggers", "internal/toon/toon.go, internal/toon/other.go", nil, [][3]string{{"internal/conformance/data_handling_test.go", "registry " + proposalRegistryCitation, "spec fence expansion required"}, {"internal/toon/toon_test.go", "registry " + proposalRegistryCitation, "spec fence expansion required"}}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			root, slug := seedProposalWith(t, test.writes, nil, test.fence...)
			out, code := Command(proposalArgs(t, root, slug))
			if code != 0 || !strings.Contains(out, "writes_proposal[") {
				t.Fatalf("proposal = (%d):\n%s", code, out)
			}
			for _, row := range test.want {
				assertProposalRow(t, out, row[0], row[1], row[2])
			}
		})
	}

	root, slug := seedProposal(t)
	preflighttest.ActiveAssignment(t, root, root)
	charge := []string{"build", slug, "--charge", "--ticket", "one.md", "--base", preflighttest.RunGit(t, "rev-parse", "main"), "--source-tip", preflighttest.RunGit(t, "rev-parse", "HEAD")}
	if chargeOut, chargeCode := Command(charge); chargeCode != 1 || !strings.Contains(chargeOut, "fixture-closure") {
		t.Errorf("charge after proposal = (%d):\n%s", chargeCode, chargeOut)
	}
}

func TestWritesProposalLinkedRepoNeedsNoKitRegistrySource(t *testing.T) {
	root, slug := seedProposalWith(t, "internal/toon/toon.go", nil)
	out, code := Command(proposalArgs(t, root, slug))
	if code != 0 || !strings.Contains(out, "registry "+proposalRegistryCitation) {
		t.Fatalf("linked repository proposal = (%d):\n%s", code, out)
	}
}

// TestWritesProposalAlreadyCovered covers DP13 without adding a repeat request.
func TestWritesProposalAlreadyCovered(t *testing.T) {
	for _, test := range []struct{ name, writes string }{
		{"exact", "internal/example/pinned.go, tests/canary/example-family/pinning-fixture, internal/toon/toon.go, internal/conformance/data_handling_test.go, internal/toon/toon_test.go"},
		{"prefix", "internal/example/pinned.go, tests/canary, internal/toon, internal/conformance"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, slug := seedProposalWith(t, test.writes, map[string]string{"two.md": "tests/canary/example-family/pinning-fixture/EXPECT"})
			writeProposalTicket(t, slug, "two.md", "Two", "one.md", "tests/canary/example-family/pinning-fixture/EXPECT")
			preflighttest.RunGit(t, "add", "specs/"+slug)
			preflighttest.RunGit(t, "commit", "-q", "-m", "approve closure and ordering")
			args := proposalArgs(t, root, slug)
			out, code := Command(args)
			if code != 0 || !strings.Contains(out, "writes_proposal[0]{path,source,fence}") || !strings.Contains(out, "ordering[0]{ticket,other,required}") {
				t.Fatalf("covered closure = (%d):\n%s", code, out)
			}
			repeat, repeatCode := Command(args)
			if repeatCode != 0 || repeat != out {
				t.Fatalf("repeat covered closure = (%d, equal=%t):\n%s", repeatCode, repeat == out, repeat)
			}
			charge := []string{"build", slug, "--charge", "--ticket", "one.md", "--base", preflighttest.RunGit(t, "rev-parse", "main"), "--source-tip", preflighttest.RunGit(t, "rev-parse", "HEAD")}
			if chargeOut, chargeCode := Command(charge); chargeCode != 0 || !strings.HasPrefix(chargeOut, "prepared[1]") {
				t.Fatalf("approved charge = (%d):\n%s", chargeCode, chargeOut)
			}
		})
	}
}

// TestWritesProposalRefusalOrder covers DP14: grammar wins over closure output.
func TestWritesProposalRefusalOrder(t *testing.T) {
	for _, body := range []string{"# broken\n", strings.Replace(preflighttest.TicketDoc("One", "PF1", "PF2"), "Covers:", "Writes: internal/example/pinned.go\nCovers:", 1)} {
		root, slug := seedProposal(t)
		path := "specs/" + slug + "/tickets/one.md"
		preflighttest.MustWriteFile(t, path, body)
		preflighttest.RunGit(t, "add", path)
		preflighttest.RunGit(t, "commit", "-q", "-m", "malformed ticket")
		out, code := Command(proposalArgs(t, root, slug))
		if code != 1 || !strings.Contains(out, "tickets-parse") || strings.Contains(out, "writes_proposal") {
			t.Fatalf("grammar before closure = (%d):\n%s", code, out)
		}
	}
}
