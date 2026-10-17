// Package ownershiptest owns the committed proposal repository that the proposal command
// tests of package preflight and of package ownership share. Only tests import it. It does
// not import package preflight, because that package's own internal tests import it.
package ownershiptest

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/tickets"
)

// kitRegistry is the kit's own binding registry source, read from the module root before
// any test changes the working directory. A test binary starts in its package directory, so
// the module root is the nearest ancestor that holds go.mod.
var kitRegistry, kitRegistryErr = func() ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return os.ReadFile(filepath.Join(dir, filepath.FromSlash(tickets.BindingRegistry)))
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return nil, errors.New("no go.mod above the test package directory")
		}
		dir = parent
	}
}()

// KitRegistry is the kit's own binding registry source. It fails the test when the source
// is unreadable.
func KitRegistry(t *testing.T) []byte {
	t.Helper()
	if kitRegistryErr != nil {
		t.Fatalf("kit registry source: %v", kitRegistryErr)
	}
	return kitRegistry
}

// Cut removes the one occurrence of part from s. It fails the test unless s holds part
// exactly once, so a drifted fixture cannot pass unchanged.
func Cut(t *testing.T, s, part string) string {
	t.Helper()
	if n := strings.Count(s, part); n != 1 {
		t.Fatalf("fixture text holds %q %d times, want once", part, n)
	}
	return strings.Replace(s, part, "", 1)
}

// ProposalRepo is one committed proposal repository over the slug example. The selected
// ticket one.md writes the conformant fence plus Writes, and each Others ticket writes its
// own list. A nil Fence fences every write, so only an explicit Fence separates the spec
// from the ticket union. The spec declares the complete-empty inventory unless
// OwnInventory is set; Tail then closes the spec, so its bytes are part of what the fixture
// admits. Files are committed in the base; a charge source that Files omits takes a
// one-line body. A non-nil Seal publishes a verified kit binary over the kit's binding
// registry, so that registry source is the authenticated kit context; a registry that Files
// omits takes the kit's own.
type ProposalRepo struct {
	Writes       []string
	Others       map[string][]string
	Files        map[string]string
	Fence        []string
	OwnInventory bool
	Seal         Publisher
	Tail         string
}

// KitModule is the module path of a sealed repository, so a sealed fixture's import names a
// package inside it.
const KitModule = "example.com/kit"

// Publisher is the sealed run-binary publication that a test file sets as Seal. A test file
// supplies it, so this non-test package cannot publish a sealed run binary on its own.
type Publisher func(root, staged, executable, manifestDir, version string) error

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
	if r.OwnInventory {
		body = preflighttest.SpecBodyNoInventory(slug, fence...)
	}
	preflighttest.MustWriteFile(t, spec, body+r.Tail)
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
	if r.Seal != nil {
		if _, given := r.Files[tickets.BindingRegistry]; !given {
			preflighttest.MustWriteFile(t, tickets.BindingRegistry, string(KitRegistry(t)))
		}
		preflighttest.MustWriteFile(t, "go.mod", "module "+KitModule+"\n\ngo 1.25\n")
		preflighttest.MustWriteFile(t, "cmd/bench/main.go", "package main\n\nimport _ \""+KitModule+"/internal/tickets\"\n\nfunc main() {}\n")
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
	if r.Seal != nil {
		staged := filepath.Join(t.TempDir(), "bench")
		if err := os.WriteFile(staged, []byte("Bench executable"), 0o755); err != nil {
			t.Fatal(err)
		}
		executable := freshness.PublishedExecutable(root)
		if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := r.Seal(root, staged, executable, filepath.Dir(executable), "1.2.3"); err != nil {
			t.Fatalf("Publish: %v", err)
		}
	}
	return root
}
