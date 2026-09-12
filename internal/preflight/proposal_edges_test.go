package preflight

import (
	"strconv"
	"strings"
	"testing"
)

type proposalTicketSpec struct {
	name, blockers, writes string
}

func setProposalGraph(t *testing.T, slug string, specs ...proposalTicketSpec) {
	t.Helper()
	names := make([]string, 0, len(specs))
	for _, spec := range specs {
		writeProposalTicket(t, slug, spec.name, strings.TrimSuffix(spec.name, ".md"), spec.blockers, spec.writes)
		names = append(names, spec.name)
	}
	replanSpec(t, slug, names...)
	runGit(t, "add", "specs/"+slug)
	runGit(t, "commit", "-q", "-m", "set proposal graph")
}

func assertOrdering(t *testing.T, out string, pairs ...string) {
	t.Helper()
	wantHeader := "ordering[" + strconv.Itoa(len(pairs)) + "]{ticket,other,required}"
	if !strings.Contains(out, wantHeader) {
		t.Fatalf("ordering count = %d:\n%s", len(pairs), out)
	}
	for _, pair := range pairs {
		row := "\n  one.md," + pair + ",approved ordering edge required\n"
		if strings.Count(out, row) != 1 {
			t.Errorf("ordering pair %s count = %d, want 1:\n%s", pair, strings.Count(out, row), out)
		}
	}
}

// TestWritesProposalEdges covers DP12 through real parsed ticket graphs.
func TestWritesProposalEdges(t *testing.T) {
	t.Run("no overlap", func(t *testing.T) {
		root, slug := seedProposalWith(t, "internal/example/pinned.go", map[string]string{"two.md": "specs"})
		out, code := Command(proposalArgs(t, root, slug))
		if code != 0 {
			t.Fatalf("proposal = (%d):\n%s", code, out)
		}
		assertOrdering(t, out)
	})

	t.Run("unordered fixture overlap", func(t *testing.T) {
		root, slug := seedProposalWith(t, "internal/example/pinned.go", map[string]string{"two.md": "tests/canary/example-family/pinning-fixture/EXPECT"})
		out, code := Command(proposalArgs(t, root, slug))
		if code != 0 {
			t.Fatalf("proposal = (%d):\n%s", code, out)
		}
		assertOrdering(t, out, "two.md")
	})

	t.Run("multiple unordered pairs", func(t *testing.T) {
		root, slug := seedProposalWith(t, "internal/example/pinned.go", map[string]string{
			"two.md":   "tests/canary/example-family/pinning-fixture/EXPECT",
			"three.md": "tests/canary/example-family/pinning-fixture/BASE",
		})
		out, code := Command(proposalArgs(t, root, slug))
		if code != 0 {
			t.Fatalf("proposal = (%d):\n%s", code, out)
		}
		assertOrdering(t, out, "three.md", "two.md")
	})

	t.Run("directory prefix overlap", func(t *testing.T) {
		root, slug := seedProposalWith(t, "internal/example/pinned.go", map[string]string{"two.md": "tests/canary/example-family"})
		out, code := Command(proposalArgs(t, root, slug))
		if code != 0 {
			t.Fatalf("proposal = (%d):\n%s", code, out)
		}
		assertOrdering(t, out, "two.md")
	})

	for _, test := range []struct {
		name string
		one  string
		two  string
	}{
		{"selected blocks other", "two.md", "none"},
		{"other blocks selected", "none", "one.md"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, slug := seedProposal(t)
			setProposalGraph(t, slug,
				proposalTicketSpec{"one.md", test.one, "internal/example/pinned.go"},
				proposalTicketSpec{"two.md", test.two, "tests/canary/example-family/pinning-fixture/EXPECT"})
			out, code := Command(proposalArgs(t, root, slug))
			if code != 0 {
				t.Fatalf("proposal = (%d):\n%s", code, out)
			}
			assertOrdering(t, out)
		})
	}

	t.Run("transitive blocker", func(t *testing.T) {
		root, slug := seedProposal(t)
		setProposalGraph(t, slug,
			proposalTicketSpec{"one.md", "middle.md", "internal/example/pinned.go"},
			proposalTicketSpec{"middle.md", "two.md", "specs"},
			proposalTicketSpec{"two.md", "none", "tests/canary/example-family/pinning-fixture/EXPECT"})
		out, code := Command(proposalArgs(t, root, slug))
		if code != 0 {
			t.Fatalf("proposal = (%d):\n%s", code, out)
		}
		assertOrdering(t, out)
	})
}
