package preflight

import (
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/tickets"
)

const phaseRefusal = "error: source required: .agents/commands/bench-implement-spec.md "
const sourceNext = " — restore the named canonical source and rerun the exact charge\n"

// TestEvidenceRequiredSourceStates covers CE41 through CE48. Every Git-trackable state
// runs through the public command, using the same setups as the enumerated legacy
// differential. The special kinds that cannot enter Git reach the no-follow source
// adapter below the checkout guard, which would otherwise answer first.
func TestEvidenceRequiredSourceStates(t *testing.T) {
	for _, test := range []struct {
		name, mutation, want string
	}{
		{"CE41 absent", "absent phase", "is absent: "},
		{"CE42 empty", "empty phase", "is empty: "},
		{"CE43 symlink", "symlink phase", "is wrong-type: not a regular file: L---------"},
		{"CE47 directory", "directory phase", "is wrong-type: not a regular file: d---------"},
		{"CE48 control byte", "control byte phase", "contains a byte spec-TOON cannot represent"},
	} {
		for _, route := range chargeRoutes {
			t.Run(test.name+" "+route.name, func(t *testing.T) {
				root, slug := seedConformant(t)
				mutationNamed(t, test.mutation)(t, root, slug, nil)
				out, code := Command(legacyCommitted(t, root, slug, test.name, route.full))
				if want := phaseRefusal + test.want + sourceNext; code != 1 || out != want {
					t.Fatalf("%s = (%d, %q), want (1, %q)", test.name, code, out, want)
				}
				assertNothingPublished(t, root)
			})
		}
	}
	for _, route := range chargeRoutes {
		t.Run("CE48 invalid UTF-8 "+route.name, func(t *testing.T) {
			root, slug := seedConformant(t)
			mustWriteFile(t, buildPhase, "bad \xff\n")
			out, code := Command(legacyCommitted(t, root, slug, "CE48 invalid UTF-8", route.full))
			if want := phaseRefusal + "is malformed: invalid UTF-8" + sourceNext; code != 1 || out != want {
				t.Fatalf("CE48 invalid UTF-8 = (%d, %q), want (1, %q)", code, out, want)
			}
			assertNothingPublished(t, root)
		})
	}
	for _, test := range []struct {
		name, want string
		make       func(t *testing.T) (root, path string)
	}{
		{"CE44 FIFO", "not a regular file: p---------", func(t *testing.T) (string, string) {
			root := shortTempDir(t)
			if err := syscall.Mkfifo(filepath.Join(root, "source.md"), 0o600); err != nil {
				t.Fatal(err)
			}
			return root, "source.md"
		}},
		{"CE45 socket", "not a regular file: S---------", func(t *testing.T) (string, string) {
			root := shortTempDir(t)
			listener, err := net.Listen("unix", filepath.Join(root, "source.md"))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = listener.Close() })
			return root, "source.md"
		}},
		{"CE46 device", "not a regular file: Dc---------", func(*testing.T) (string, string) {
			return "/", "dev/null"
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, path := test.make(t)
			policy := []buildSourceDescriptor{{role: "build-phase", path: fixedSource(path)}}
			inputs, failure := loadBuildSources(root, Facts{SourceTip: "HEAD"}, nil, policy)
			if want := path + " is wrong-type: " + test.want; failure != want || inputs != nil {
				t.Fatalf("%s = (%v, %q), want (nil, %q)", test.name, inputs, failure, want)
			}
		})
	}
}

// TestEvidencePreparationRefusals covers CE49 through CE52 at the public command, using
// the same shared setups as the enumerated legacy differential.
func TestEvidencePreparationRefusals(t *testing.T) {
	for _, route := range chargeRoutes {
		t.Run(route.name, func(t *testing.T) { preparationRefusals(t, route.full) })
	}
	t.Run("CE128 missing tool", func(t *testing.T) {
		root, slug := seedConformant(t)
		args := chargeArgs(t, root, slug, false)
		path := os.Getenv("PATH")
		t.Setenv("PATH", t.TempDir())
		out, code := Command(args)
		if code != 1 || len(out) > 48000 || !strings.HasPrefix(out, "error: ") || strings.Contains(out, "prepared[") {
			t.Fatalf("missing tool = (%d):\n%s", code, out)
		}
		t.Setenv("PATH", path)
		assertNothingPublished(t, root)
	})
}

// chargeRoutes are the two build charge routes that share the refusal pipeline: the legacy
// full projection and the evidence preparation.
var chargeRoutes = []struct {
	name string
	full bool
}{{"legacy full", true}, {"preparation", false}}

// assertNothingPublished proves that a refused preparation left no artifact and no
// temporary pack behind.
func assertNothingPublished(t *testing.T, root string) {
	t.Helper()
	if packs, temps := publishedPacks(t, root), stagedTemps(t, root); len(packs) != 0 || len(temps) != 0 {
		t.Fatalf("refusal left packs %v and temporary packs %v", packs, temps)
	}
}

func preparationRefusals(t *testing.T, full bool) {
	for _, test := range []struct {
		name string
		args func(t *testing.T, root, slug string) []string
		want string
	}{
		{"CE49 dirty checkout", func(t *testing.T, root, slug string) []string {
			return mutationNamed(t, "dirty checkout")(t, root, slug, chargeArgs(t, root, slug, full))
		}, "error: checkout required: source checkout is dirty — commit or remove local changes and rerun the exact charge\n"},
		{"CE50 missing ticket", func(t *testing.T, root, slug string) []string {
			return mutationNamed(t, "missing ticket")(t, root, slug, chargeArgs(t, root, slug, full))
		}, "error: ticket required: selected ticket \"missing.md\" was not found — pass a ticket basename from the spec tickets directory\n"},
		{"CE51 source-tip mismatch", func(t *testing.T, root, slug string) []string {
			return mutationNamed(t, "source-tip mismatch")(t, root, slug, chargeArgs(t, root, slug, full))
		}, "error: preflight required: tip-current: --source-tip <base> is not the derived source tip <tip> — repair tip-current and rerun the exact charge\n"},
		{"CE52 no assignment", func(t *testing.T, root, slug string) []string {
			args := mutationNamed(t, "no assignment")(t, root, slug, nil)
			if full {
				args = append(args, "--full")
			}
			return args
		}, "error: assignment required: active assignment is required — run from the assigned worktree\n"},
		{"CE52 foreign assignment", func(t *testing.T, root, slug string) []string {
			return mutationNamed(t, "foreign assignment")(t, root, slug, chargeArgs(t, root, slug, full))
		}, "error: assignment required: active assignment is required — run from the assigned worktree\n"},
		{"CE52 owned non-active assignment", func(t *testing.T, root, slug string) []string {
			args := chargeArgs(t, root, slug, full)
			canonical, err := canonicalpath.Resolve(root)
			if err != nil {
				t.Fatalf("canonicalpath.Resolve(%q): %v", root, err)
			}
			ownedAssignment(t, root, canonical, intent.StateComplete)
			return args
		}, "error: assignment required: active assignment is required — run from the assigned worktree\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, slug := seedConformant(t)
			args := test.args(t, root, slug)
			out, code := Command(args)
			got := normalizeLegacy(out, root, runGit(t, "rev-parse", "main"), runGit(t, "rev-parse", "HEAD"))
			if code != 1 || got != test.want {
				t.Fatalf("%s = (%d, %q), want (1, %q)", test.name, code, got, test.want)
			}
			assertNothingPublished(t, root)
		})
	}
}

// preparedFixtureBuild gathers the seedConformant build inputs for a direct preparation.
func preparedFixtureBuild(t *testing.T) (root string, facts Facts, entry *tickets.Entry, parsed *tickets.Ticket, args []string) {
	t.Helper()
	root, slug := seedConformant(t)
	args = chargeArgs(t, root, slug, true)
	facts, failure := gatherPinned(root, modeBuild, slug, args[6], args[8], true)
	if failure != nil {
		t.Fatalf("gather: %v", failure)
	}
	entry, parsed, detail, _ := preparationTicket(root, facts, "one.md")
	if detail != "" {
		t.Fatalf("ticket: %s", detail)
	}
	return root, facts, entry, parsed, args
}

// TestEvidenceSourcePolicy is CE114. The prepared inventory equals the independently
// listed build policy, and removing one descriptor removes that source from the pack and
// from the legacy charge, so no second inventory survives in the renderer.
func TestEvidenceSourcePolicy(t *testing.T) {
	want := []string{
		"s1 metadata derived  true",
		"s2 ticket repository specs/example/tickets/one.md true",
		"s3 spec repository specs/example/spec.md true",
		"s4 delegate-skill repository .agents/skills/bench-craft-delegate/SKILL.md true",
		"s5 build-phase repository .agents/commands/bench-implement-spec.md true",
		"s6 delegate-procedure repository .agents/skills/bench-craft-delegate/references/delegation-discipline.md true",
	}
	root, facts, entry, parsed, _ := preparedFixtureBuild(t)
	pack, failure, err := prepareBuildPack(root, facts, entry, parsed, buildSourcePolicy())
	if failure != "" || err != nil {
		t.Fatalf("prepare = (%q, %v)", failure, err)
	}
	var got []string
	for _, s := range pack.Manifest().Sources {
		got = append(got, strings.Join([]string{s.ID, s.Role, s.Kind, s.Path, map[bool]string{true: "true", false: "false"}[s.Required]}, " "))
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("prepared inventory =\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	metadata := pack.Metadata()
	if len(metadata.Charge) != 1 || metadata.Charge[0].Ticket != "s2" || strings.Join(metadata.Checks, ",") != "s2,s5" || strings.Join(metadata.Returns, ",") != "s4,s6" {
		t.Fatalf("policy columns = %+v", metadata)
	}

	reduced := buildSourcePolicy()
	reduced = reduced[:len(reduced)-1]
	pack, failure, err = prepareBuildPack(root, facts, entry, parsed, reduced)
	if failure != "" || err != nil || len(pack.Manifest().Sources) != len(want)-1 {
		t.Fatalf("reduced prepare = (%q, %v)", failure, err)
	}
	out, code := renderChargeWithPolicy(root, facts, Decide(facts), "one.md", true, reduced)
	if code != 0 || strings.Contains(out, delegateProcedure) || !strings.Contains(out, "sources[4]{path,identity}") {
		t.Fatalf("reduced policy charge still lists the omitted source (%d):\n%s", code, out)
	}
}

// TestEvidenceBuildMetadataSchema is CE149. The prepared metadata source equals the
// spec's metadata tables for the fixture, with one row per list member.
func TestEvidenceBuildMetadataSchema(t *testing.T) {
	root, facts, entry, parsed, _ := preparedFixtureBuild(t)
	pack, failure, err := prepareBuildPack(root, facts, entry, parsed, buildSourcePolicy())
	if failure != "" || err != nil {
		t.Fatalf("prepare = (%q, %v)", failure, err)
	}
	const want = "charge[1]{axis,ticket,access}:\n  \"\",s2,write-within-fence\n" +
		"fence[4]{path}:\n  internal/example/\n  reviews/example.md\n  .agents/skills/bench-craft-delegate/\n  .agents/commands/bench-implement-spec.md\n" +
		"writes[1]{path}:\n  specs\n" +
		"coverage[2]{row}:\n  PF1\n  PF2\n" +
		"checks[2]{source}:\n  s2\n  s5\n" +
		"returns[2]{source}:\n  s4\n  s6\n" +
		"shared_evidence[0]{kind,source}:\n" +
		"completion_evidence[0]{record,source_digest,plan_digest,record_state,detail}:\n"
	body, _ := pack.Source("s1")
	if string(body) != want {
		t.Fatalf("build metadata =\n%s\nwant\n%s", body, want)
	}
	manifest := pack.Manifest()
	if first := manifest.Sources[0]; first.Role != "metadata" || first.Kind != "derived" || first.Path != "" || !first.Required {
		t.Fatalf("metadata descriptor = %+v", first)
	}
}

// shortTempDir keeps a Unix socket path under the platform length limit.
func shortTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "ce")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	return dir
}
