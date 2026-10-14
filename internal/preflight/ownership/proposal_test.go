package ownership_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

const slug = "example"

// completeEmpty is a declared-complete planning inventory with no item. Each grammar case
// edits one part of it, so the cases share one valid baseline.
const completeEmpty = `{"version":1,"authority":"authored-spec","completeness":"declared-complete","baseline":"source-base","scope":[],"changes":[],"premises":[],"seams":[],"decisions":[],"manual_claims":[]}`

// fenced is one bench-planning-facts section that holds payload.
func fenced(payload string) string {
	return "\n## Planning claim inventory\n\n```bench-planning-facts\n" + payload + "\n```\n"
}

const (
	specPath   = "specs/" + slug + "/spec.md"
	ticketPath = "specs/" + slug + "/tickets/one.md"
)

// commit records every change in the working tree as the new base on main, then adds one
// authorized change. A later proposal reads the change at HEAD, and no fence needs to
// authorize it.
func commit(t *testing.T, message string) {
	t.Helper()
	preflighttest.RunGit(t, "add", "-A")
	preflighttest.RunGit(t, "commit", "-q", "-m", message)
	preflighttest.RunGit(t, "update-ref", "refs/heads/main", "HEAD")
	preflighttest.MustWriteFile(t, "internal/example/after.go", "package example\n")
	preflighttest.RunGit(t, "add", "internal/example/after.go")
	preflighttest.RunGit(t, "commit", "-q", "-m", "advance "+message)
}

// proposalArgs activates the root's assignment and answers the proposal arguments over
// one.md, pinned to main and HEAD. The base is element 6 and the tip is element 8.
func proposalArgs(t *testing.T, root string) []string {
	t.Helper()
	preflighttest.ActiveAssignment(t, root, root)
	return []string{"build", slug, "--propose-writes", "--ticket", "one.md",
		"--base", preflighttest.RunGit(t, "rev-parse", "main"), "--source-tip", preflighttest.RunGit(t, "rev-parse", "HEAD")}
}

// propose runs the actual proposal command over one.md.
func propose(t *testing.T, root string) (string, int) {
	t.Helper()
	return preflight.Command(proposalArgs(t, root))
}

// mustPropose runs the proposal and fails the test unless it answers.
func mustPropose(t *testing.T, root string) string {
	t.Helper()
	out, code := propose(t, root)
	if code != 0 {
		t.Fatalf("proposal = (%d):\n%s", code, out)
	}
	return out
}

// closureTree holds every closure family the fixed point walks. pinned.go is pinned by a
// fixture. toon.go is bound to two registry files, and a second fixture pins one of them,
// a system-tagged test. cyclic.go is pinned by a fixture whose own directory an anchor
// registry literal names, and that same fixture also pins the anchor registry file, so the
// walk returns to a path it already holds. included.go is pinned through a BASE include.
func closureTree() ownershiptest.ProposalRepo {
	return ownershiptest.ProposalRepo{
		Writes: []string{"internal/example/pinned.go", "internal/toon/toon.go", "internal/example/cyclic.go", "internal/example/included.go"},
		Files: map[string]string{
			"internal/example/pinned.go":                 "package example\n",
			"internal/example/cyclic.go":                 "package example\n",
			"internal/example/included.go":               "package example\n",
			"internal/toon/toon.go":                      "package toon\n",
			"internal/toon/toon_test.go":                 "//go:build system\n\npackage toon\n",
			"internal/conformance/data_handling_test.go": "package conformance\n",
			"tests/canary/fam/pin-fixture/BASE":          "internal/example/pinned.go\n",
			"tests/canary/fam/pin-fixture/EXPECT":        "red\n",
			"tests/canary/fam/registry-fixture/BASE":     "internal/toon/toon_test.go\n",
			"tests/canary/fam/registry-fixture/EXPECT":   "red\n",
			"tests/canary/fam/cycle-fixture/BASE":        "internal/example/cyclic.go\ninternal/anchors/cycle.go\n",
			"tests/canary/fam/cycle-fixture/EXPECT":      "red\n",
			"tests/canary/fam/include-fixture/BASE":      "@tests/canary/shared.list\n",
			"tests/canary/fam/include-fixture/EXPECT":    "red\n",
			"tests/canary/shared.list":                   "internal/example/other.go\ninternal/example/included.go\n",
			"internal/anchors/cycle.go":                  "package anchors\n\nvar cycle = []string{\n\t\"tests/canary/fam/cycle-fixture/EXPECT\",\n}\n",
		},
	}
}

// TestProposalFixedPoint drives the transitive closure through the actual proposal command.
func TestProposalFixedPoint(t *testing.T) {
	root := closureTree().Seed(t)
	out := mustPropose(t, root)
	t.Run("OC1 missing pin", func(t *testing.T) {
		one(t, out, "writes_proposal", map[string]string{"path": "tests/canary/fam/pin-fixture"})
	})
	t.Run("OC3 second-hop pin", func(t *testing.T) {
		one(t, out, "writes_proposal", map[string]string{"path": "tests/canary/fam/registry-fixture"})
		one(t, out, "closure", map[string]string{"ticket": "one.md", "path": "internal/toon/toon.go", "requires": "tests/canary/fam/registry-fixture"})
	})
	t.Run("OC4 cycle", func(t *testing.T) {
		for _, required := range []string{"tests/canary/fam/cycle-fixture", "internal/anchors/cycle.go"} {
			one(t, out, "writes_proposal", map[string]string{"path": required})
			one(t, out, "closure", map[string]string{"ticket": "one.md", "path": "internal/example/cyclic.go", "requires": required})
		}
	})
	t.Run("OC85 second-hop system test", func(t *testing.T) {
		row := one(t, out, "closure_evidence", map[string]string{"path": "internal/toon/toon_test.go", "kind": "kit-pin"})
		if cell(row, "source") != "internal/toon/toon_test.go" || cell(row, "line") != "1" || !strings.Contains(cell(row, "reason"), "BENCH_KIT") {
			t.Fatalf("kit-pin evidence = %v, want the tagged file, its build line, and BENCH_KIT", row)
		}
	})
}

// TestProposalProvenance grades each reason's owner, line, and source identity.
func TestProposalProvenance(t *testing.T) {
	root := closureTree().Seed(t)
	out := mustPropose(t, root)
	t.Run("OC5 two reasons", func(t *testing.T) {
		const base = "tests/canary/fam/cycle-fixture/BASE"
		for line, subject := range map[string]string{"1": "internal/example/cyclic.go", "2": "internal/anchors/cycle.go"} {
			row := one(t, out, "closure_evidence", map[string]string{"path": "tests/canary/fam/cycle-fixture", "kind": "fixture", "line": line})
			if cell(row, "source") != base || !strings.Contains(cell(row, "reason"), subject) {
				t.Errorf("fixture reason at line %s = %v, want %s naming %s", line, row, base, subject)
			}
		}
	})
	t.Run("OC6 anchor holder line", func(t *testing.T) {
		row := one(t, out, "closure_evidence", map[string]string{"path": "internal/anchors/cycle.go", "kind": "anchor"})
		if cell(row, "source") != "internal/anchors/cycle.go" || cell(row, "line") != "4" || !strings.Contains(cell(row, "reason"), "tests/canary/fam/cycle-fixture/EXPECT") {
			t.Fatalf("anchor reason = %v, want the holder, line 4, and its literal", row)
		}
	})
	t.Run("OC7 included BASE source", func(t *testing.T) {
		row := one(t, out, "closure_evidence", map[string]string{"path": "tests/canary/fam/include-fixture", "kind": "fixture"})
		if cell(row, "source") != "tests/canary/shared.list" || cell(row, "line") != "2" {
			t.Fatalf("include reason = %v, want tests/canary/shared.list line 2", row)
		}
	})
	t.Run("OC8 registry owner", func(t *testing.T) {
		root, line := sealedKit(t)
		out := mustPropose(t, root)
		row := one(t, out, "closure_evidence", map[string]string{"path": "internal/toon/toon_test.go", "kind": "registry"})
		if cell(row, "source") != "internal/tickets/registry_data.go" || cell(row, "line") != line {
			t.Fatalf("registry reason = %v, want internal/tickets/registry_data.go line %s", row, line)
		}
	})
	t.Run("OC9 exact source tip", func(t *testing.T) {
		args := proposalArgs(t, root)
		tip := args[8]
		args[8] = preflighttest.RunGit(t, "rev-parse", "--short", "HEAD")
		out, code := preflight.Command(args)
		if code != 0 {
			t.Fatalf("abbreviated tip proposal = (%d):\n%s", code, out)
		}
		one(t, out, "source", map[string]string{"base": args[6], "tip": tip})
	})
}

// kitRegistry is the kit's own binding registry source, read from the package directory
// before any test changes the working directory.
var kitRegistry, kitRegistryErr = os.ReadFile(filepath.Join("..", "..", "tickets", "registry_data.go"))

// sealedKit seeds a repository whose published binary seal verifies, so its binding
// registry source is the caller's authenticated kit context. The registry source is the
// kit's own file. It answers the root and the line of the internal/toon binding row,
// read from those bytes by a plain text search.
func sealedKit(t *testing.T) (string, string) {
	t.Helper()
	registry, err := kitRegistry, kitRegistryErr
	if err != nil {
		t.Fatal(err)
	}
	at := bytes.Index(registry, []byte(`{Prefix: "internal/toon",`))
	if at < 0 {
		t.Fatal("kit registry holds no internal/toon binding row")
	}
	line := fmt.Sprint(bytes.Count(registry[:at], []byte("\n")) + 1)
	kit := ownershiptest.ProposalRepo{
		Writes: []string{"internal/toon/toon.go"},
		Files: map[string]string{
			"go.mod":                                     "module example.com/kit\n\ngo 1.25\n",
			"cmd/bench/main.go":                          "package main\n\nimport _ \"example.com/kit/internal/tickets\"\n\nfunc main() {}\n",
			"internal/tickets/registry_data.go":          string(registry),
			"scripts/go-build.sh":                        "#!/usr/bin/env bash\n",
			freshness.BuildInputsManifest:                freshness.BuildInputLine("build_script", "scripts/go-build.sh"),
			".gitignore":                                 "dist/\n",
			"internal/toon/toon.go":                      "package toon\n",
			"internal/toon/toon_test.go":                 "package toon\n",
			"internal/conformance/data_handling_test.go": "package conformance\n",
		},
	}
	root := kit.Seed(t)
	staged := filepath.Join(t.TempDir(), "bench")
	if err := os.WriteFile(staged, []byte("Bench executable"), 0o755); err != nil {
		t.Fatal(err)
	}
	executable := freshness.PublishedExecutable(root)
	if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := freshness.Publish(root, staged, executable, filepath.Dir(executable), "1.2.3"); err != nil {
		t.Fatalf("Publish: %v", err)
	}
	return root, line
}

// authorityState is every byte a proposal must leave alone: the spec, the ticket, the
// index, the status, and HEAD.
func authorityState(t *testing.T, root string) string {
	t.Helper()
	var state strings.Builder
	for _, path := range []string{specPath, ticketPath} {
		data, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatal(err)
		}
		state.Write(data)
	}
	index, err := os.ReadFile(filepath.Join(root, ".git", "index"))
	if err != nil {
		t.Fatal(err)
	}
	state.Write(index)
	status, err := exec.Command("git", "-C", root, "status", "--porcelain=v1", "--untracked-files=all").Output()
	if err != nil {
		t.Fatal(err)
	}
	state.Write(status)
	state.WriteString(preflighttest.RunGit(t, "rev-parse", "HEAD"))
	return state.String()
}

// TestProposalReadOnly grades that two proposals change no authority byte and agree.
func TestProposalReadOnly(t *testing.T) {
	root := closureTree().Seed(t)
	args := proposalArgs(t, root)
	before := authorityState(t, root)
	first, firstCode := preflight.Command(args)
	second, secondCode := preflight.Command(args)
	if firstCode != 0 || secondCode != 0 || first != second {
		t.Fatalf("repeat proposal = (%d, %d, equal=%t):\n%s", firstCode, secondCode, first == second, first)
	}
	if after := authorityState(t, root); after != before {
		t.Fatal("proposal changed the spec, the ticket, the index, the status, or HEAD")
	}
}

// emptyTree is a linked repository with no anchor registry and no fixture inventory, whose
// ticket writes a guidance path and a source file. tail closes its spec.
func emptyTree(tail string) ownershiptest.ProposalRepo {
	return ownershiptest.ProposalRepo{
		Writes: []string{".agents/x/SKILL.md", "internal/example/pinned.go"},
		Files: map[string]string{
			".agents/x/SKILL.md":         "# X\n",
			"internal/example/pinned.go": "package example\n",
		},
		Tail: tail,
	}
}

// TestProposalEmpty covers the canonical empty inventories and the typed empty response.
func TestProposalEmpty(t *testing.T) {
	out := mustPropose(t, emptyTree(fenced(completeEmpty)).Seed(t))
	t.Run("OC10 no anchor registry", func(t *testing.T) {
		if rows := matching(t, out, "closure", map[string]string{}); len(rows) != 0 {
			t.Fatalf("inventoryless closure = %v, want none", rows)
		}
	})
	t.Run("OC11 no fixture inventory", func(t *testing.T) {
		if rows := matching(t, out, "closure_evidence", map[string]string{}); len(rows) != 0 {
			t.Fatalf("inventoryless evidence = %v, want none", rows)
		}
	})
	t.Run("OC12 typed zero rows", func(t *testing.T) {
		for _, header := range []string{"writes_proposal[0]{path,source,fence}", "ordering[0]{ticket,other,required}", "closure[0]{ticket,path,requires}",
			"closure_evidence[0]{ticket,path,kind,source,line,reason}", "source[1]{base,tip}", "premises[0]{id,state,source,reason}"} {
			if !strings.Contains(out, header) {
				t.Errorf("complete empty proposal lacks %s:\n%s", header, out)
			}
		}
	})
	t.Run("absent manifest stays diagnostic", func(t *testing.T) {
		out := mustPropose(t, emptyTree("").Seed(t))
		row := one(t, out, "premises", map[string]string{"id": "metadata"})
		if cell(row, "state") != "unknown" || !strings.Contains(cell(row, "reason"), "no bench-planning-facts fence") {
			t.Fatalf("absent manifest premise = %v, want an unknown inventory diagnostic", row)
		}
	})
}

// TestProposalUnknown covers collectors that cannot certify their source.
func TestProposalUnknown(t *testing.T) {
	t.Run("OC13 unverified registry source", func(t *testing.T) {
		tr := closureTree()
		tr.Tail = fenced(completeEmpty)
		out := mustPropose(t, tr.Seed(t))
		row := one(t, out, "premises", map[string]string{"id": "kit-registry"})
		if cell(row, "state") != "unknown" {
			t.Fatalf("unverified registry premise = %v, want unknown", row)
		}
		evidence := one(t, out, "closure_evidence", map[string]string{"path": "internal/toon/toon_test.go", "kind": "registry"})
		if cell(evidence, "line") != "0" {
			t.Fatalf("unverified registry evidence = %v, want no source line", evidence)
		}
	})
	t.Run("OC83 unreadable fixture inventory", func(t *testing.T) {
		for name, place := range map[string]func(*testing.T){
			"regular file": func(t *testing.T) {
				if err := os.RemoveAll("tests/canary"); err != nil {
					t.Fatal(err)
				}
				preflighttest.MustWriteFile(t, "tests/canary", "not a directory\n")
			},
			"symlink": func(t *testing.T) {
				if err := os.Rename("tests/canary", "outside-canary"); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink("../outside-canary", "tests/canary"); err != nil {
					t.Fatal(err)
				}
			},
		} {
			t.Run(name, func(t *testing.T) {
				root := closureTree().Seed(t)
				place(t)
				commit(t, "special inventory")
				out, code := propose(t, root)
				refused(t, out, code, "fixture inventory")
			})
		}
	})
	t.Run("OC84 unreadable anchor registry", func(t *testing.T) {
		root := closureTree().Seed(t)
		preflighttest.MustWriteFile(t, "outside.go", "package anchors\n\nvar x = \"tests/canary/fam/cycle-fixture/EXPECT\"\n")
		if err := os.Symlink("../../outside.go", "internal/anchors/linked.go"); err != nil {
			t.Fatal(err)
		}
		commit(t, "linked registry")
		out, code := propose(t, root)
		refused(t, out, code, "anchor registry")
	})
}
