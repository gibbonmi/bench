package preflight

import (
	"strings"
	"testing"
)

const proposalRegistryCitation = "internal/tickets/registry_data.go"

func proposalArgs(t *testing.T, root, slug string) []string {
	t.Helper()
	activeAssignment(t, root, root)
	return []string{"build", slug, "--propose-writes", "--ticket", "one.md", "--base", runGit(t, "rev-parse", "main"), "--source-tip", runGit(t, "rev-parse", "HEAD")}
}

func seedProposal(t *testing.T) (string, string) {
	return seedProposalWith(t,
		"internal/example/pinned.go, internal/toon/toon.go",
		map[string]string{"two.md": "tests/canary/example-family/pinning-fixture/EXPECT"})
}

func seedProposalWith(t *testing.T, selectedWrites string, others map[string]string, extraFence ...string) (string, string) {
	t.Helper()
	root := initRepo(t)
	slug := "example"
	mustWriteFile(t, "specs/"+slug+"/spec.md", specBody(slug, extraFence...))
	writeProposalTicket(t, slug, "one.md", "One", "none", selectedWrites)
	for name, writes := range others {
		writeProposalTicket(t, slug, name, strings.TrimSuffix(name, ".md"), "none", writes)
	}
	mustWriteFile(t, "tests/canary/example-family/pinning-fixture/BASE", "internal/example/pinned.go\n")
	mustWriteFile(t, "tests/canary/example-family/pinning-fixture/EXPECT", "fixture\n")
	mustWriteFile(t, "internal/example/pinned.go", "package example\n")
	mustWriteFile(t, "internal/toon/toon.go", "package toon\n")
	mustWriteFile(t, "internal/toon/other.go", "package toon\n")
	mustWriteFile(t, "internal/toon/toon_test.go", "package toon\n")
	mustWriteFile(t, "internal/conformance/data_handling_test.go", "package conformance\n")
	mustWriteFile(t, delegateSkill, "# Delegation skill\n")
	mustWriteFile(t, delegateProcedure, "# Delegation procedure\n")
	mustWriteFile(t, buildPhase, "# Build phase\n")
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "c0")
	runGit(t, "checkout", "-q", "-b", "feature")
	mustWriteFile(t, "internal/example/foo.go", "package example\n")
	runGit(t, "add", "internal/example/foo.go")
	runGit(t, "commit", "-q", "-m", "c1")
	return root, slug
}

func writeProposalTicket(t *testing.T, slug, name, title, blockers, writes string) {
	t.Helper()
	body := strings.Replace(ticketDoc(title, "PF1", "PF2"), "Blocked by: none", "Blocked by: "+blockers, 1)
	body = strings.Replace(body, "Writes: specs", "Writes: "+writes, 1)
	mustWriteFile(t, "specs/"+slug+"/tickets/"+name, body)
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
	activeAssignment(t, root, root)
	charge := append([]string{"build", slug, "--charge", "--ticket", "one.md", "--base", runGit(t, "rev-parse", "main"), "--source-tip", runGit(t, "rev-parse", "HEAD")}, "--full")
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
			fences := []string{"- `tests/canary/`", "- `internal/toon/`", "- `internal/conformance/`"}
			root, slug := seedProposalWith(t, test.writes, map[string]string{"two.md": "tests/canary/example-family/pinning-fixture/EXPECT"}, fences...)
			writeProposalTicket(t, slug, "two.md", "Two", "one.md", "tests/canary/example-family/pinning-fixture/EXPECT")
			runGit(t, "add", "specs/"+slug)
			runGit(t, "commit", "-q", "-m", "approve closure and ordering")
			args := proposalArgs(t, root, slug)
			out, code := Command(args)
			if code != 0 || !strings.Contains(out, "writes_proposal[0]{path,source,fence}") || !strings.Contains(out, "ordering[0]{ticket,other,required}") {
				t.Fatalf("covered closure = (%d):\n%s", code, out)
			}
			repeat, repeatCode := Command(args)
			if repeatCode != 0 || repeat != out {
				t.Fatalf("repeat covered closure = (%d, equal=%t):\n%s", repeatCode, repeat == out, repeat)
			}
			charge := []string{"build", slug, "--charge", "--ticket", "one.md", "--base", runGit(t, "rev-parse", "main"), "--source-tip", runGit(t, "rev-parse", "HEAD"), "--full"}
			if chargeOut, chargeCode := Command(charge); chargeCode != 0 || !strings.Contains(chargeOut, "\"true\"") {
				t.Fatalf("approved charge = (%d):\n%s", chargeCode, chargeOut)
			}
		})
	}
}

// TestWritesProposalRefusalOrder covers DP14: grammar wins over closure output.
func TestWritesProposalRefusalOrder(t *testing.T) {
	for _, body := range []string{"# broken\n", strings.Replace(ticketDoc("One", "PF1", "PF2"), "Writes: specs", "Writes: internal/example/pinned.go\nWrites: specs", 1)} {
		root, slug := seedProposal(t)
		path := "specs/" + slug + "/tickets/one.md"
		mustWriteFile(t, path, body)
		runGit(t, "add", path)
		runGit(t, "commit", "-q", "-m", "malformed ticket")
		out, code := Command(proposalArgs(t, root, slug))
		if code != 1 || !strings.Contains(out, "tickets-parse") || strings.Contains(out, "writes_proposal") {
			t.Fatalf("grammar before closure = (%d):\n%s", code, out)
		}
	}
}
