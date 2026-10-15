// Package ownershiptest owns the committed proposal repository that the proposal command
// tests of package preflight and of package ownership share. Only tests import it. It does
// not import package preflight, because that package's own internal tests import it.
package ownershiptest

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

var kitRegistry = func() []byte {
	for _, rel := range []string{"../tickets/registry_data.go", "../../tickets/registry_data.go"} {
		if data, err := os.ReadFile(rel); err == nil {
			return data
		}
	}
	return nil
}()

// ProposalRepo is one committed proposal repository over the slug example. The selected
// ticket one.md writes the conformant fence plus Writes, and each Others ticket writes its
// own list. A nil Fence fences every write, so only an explicit Fence separates the spec
// from the ticket union. Tail closes the spec, so its bytes are part of what the fixture
// admits. Files are committed in the base; a charge source that Files omits takes a
// one-line body.
type ProposalRepo struct {
	Writes []string
	Others map[string][]string
	Files  map[string]string
	Fence  []string
	Tail   string
}

// Seed commits the repository as the base on main, then adds one authorized change on a
// feature branch, and returns the root, which is also the working directory.
func (r ProposalRepo) Seed(t *testing.T) string {
	t.Helper()
	const slug = "example"
	root := preflighttest.StartRepo(t)
	others := make([]string, 0, len(r.Others))
	for name := range r.Others {
		others = append(others, name)
	}
	slices.Sort(others)
	fence := r.Fence
	if fence == nil {
		for _, write := range r.Writes {
			fence = append(fence, "- `"+write+"`")
		}
		for _, name := range others {
			for _, write := range r.Others[name] {
				fence = append(fence, "- `"+write+"`")
			}
		}
	}
	spec := "specs/" + slug + "/spec.md"
	body := preflighttest.SpecBody(slug, fence...)
	if r.Tail != "" {
		body = strings.Replace(body, "\n## Planning claim inventory\n\n```bench-planning-facts\n{\"version\":1,\"authority\":\"authored-spec\",\"completeness\":\"declared-complete\",\"baseline\":\"source-base\",\"scope\":[],\"changes\":[],\"premises\":[],\"seams\":[],\"decisions\":[],\"manual_claims\":[]}\n```\n", "", 1) + r.Tail
	}
	preflighttest.MustWriteFile(t, spec, body)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md", preflighttest.WritesTicketDoc("One", append(preflighttest.FenceWrites(preflighttest.ConformantFence), r.Writes...), "PF1", "PF2"))
	for _, name := range others {
		preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/"+name, preflighttest.WritesTicketDoc(strings.TrimSuffix(name, ".md"), r.Others[name], "PF1", "PF2"))
	}
	for path, body := range r.Files {
		preflighttest.MustWriteFile(t, path, body)
	}
	for source, body := range map[string]string{chargesource.DelegateSkill: "# Delegation skill\n", chargesource.DelegateProcedure: "# Delegation procedure\n", chargesource.BuildPhase: "# Build phase\n"} {
		if _, given := r.Files[source]; !given {
			preflighttest.MustWriteFile(t, source, body)
		}
	}
	_, hasRegistry := r.Files["internal/tickets/registry_data.go"]
	seal := hasRegistry || (len(r.Others) > 0 && len(r.Writes) >= 4)
	if seal {
		if !hasRegistry {
			preflighttest.MustWriteFile(t, "internal/tickets/registry_data.go", string(kitRegistry))
		}
		preflighttest.MustWriteFile(t, "go.mod", "module example.com/kit\n\ngo 1.25\n")
		preflighttest.MustWriteFile(t, "cmd/bench/main.go", "package main\n\nimport _ \"example.com/kit/internal/tickets\"\n\nfunc main() {}\n")
		preflighttest.MustWriteFile(t, "scripts/go-build.sh", "#!/usr/bin/env bash\n")
		preflighttest.MustWriteFile(t, freshness.BuildInputsManifest, freshness.BuildInputLine("build_script", "scripts/go-build.sh"))
		preflighttest.MustWriteFile(t, ".gitignore", "dist/\n")
	}
	commitmenttest.SeedAdmission(t, root, spec)
	preflighttest.RunGit(t, "add", ".")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c0")
	preflighttest.RunGit(t, "checkout", "-q", "-b", "feature")
	preflighttest.MustWriteFile(t, "internal/example/foo.go", "package example\n")
	preflighttest.RunGit(t, "add", "internal/example/foo.go")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c1")
	if seal {
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
	}
	return root
}
