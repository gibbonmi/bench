package axi

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/sanitize"
)

func TestAXIRoadmapContextContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "AXI roadmap context completeness contract failed", testRoadmapContextComplete)
	contract.RunParallel(t, "AXI roadmap context truncation contract failed", testRoadmapContextTruncation)
	contract.RunParallel(t, "AXI roadmap context usage contract failed", testRoadmapContextUsage)
	contract.RunParallel(t, "AXI roadmap context source-state contract failed", testRoadmapContextSourceStates)
	contract.RunParallel(t, "AXI roadmap context fail-closed contract failed", testRoadmapContextFailClosed)
	contract.RunParallel(t, "AXI roadmap context read-only offline contract failed", testRoadmapContextReadOnlyOffline)
	contract.RunParallel(t, "AXI roadmap context unborn-HEAD contract failed", testRoadmapContextUnbornHead)
}

func contextFixture(t *testing.T) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	f.Git("branch", "-M", "main")
	f.WriteFile("ROADMAP.md", "# Roadmap\n\n## Features\n\n**FT1 — one.** Body specs/one/spec.md.\n\n## Recommended sequence\n\n1. `/bench-implement-spec` — one\n")
	f.WriteFile("IDEAS.md", "- 2026-07-10  retain me\n")
	f.WriteFile(".bench/learnings.md", "## 2026-07-10 — lesson  [open]\n- body\n")
	f.WriteFile(".bench/structure.budgets", "")
	f.WriteFile(".bench/structure-accept", "")
	f.WriteFile("specs/one/spec.md", "# One\n\nStatus: staged\nRoadmap: FT1\n")
	f.Git("add", "-A")
	f.Git("-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-qm", "fixture")
	return f
}

func testRoadmapContextComplete(t *testing.T) {
	contract.NoteContractFailure(t, "AXI roadmap context completeness contract failed")
	f := contextFixture(t)
	out := f.Bench("roadmap", "--context")
	out.RequireExit(0)
	if out.Stderr != "" {
		t.Fatalf("stderr not empty: %s", out.Stderr)
	}
	out.RequireContains(out.Stdout, "2,false")
	headers := []string{
		"context[1]{schema,full}:", "sources[7]{source,state,bytes}:",
		"roadmap_rows[1]{id,title,spec,spec_status,external_trigger,body,body_bytes,truncated}:",
		"roadmap_sequence[1]{rank,text,command}:", "ideas[1]{date,text,text_bytes,truncated}:",
		"learnings[1]{date,title,state,body,body_bytes,truncated}:", "structure[0]{kind,path,actual,limit,state,detail}:",
		"specs[1]{slug,status,roadmap_id}:", "spec_history[0]{slug,hash,date,kind,subject}:",
		"git[1]{branch,default_branch,dirty,ahead,behind}:", "git_changes[0]{status,path}:",
		"gate_cache[1]{present,state,pending_status,status,cached_tree,work_tree,timestamp,stale}:",
		"parse_failures[0]{source,reason,raw,raw_bytes,truncated}:",
	}
	last := -1
	for _, h := range headers {
		p := strings.Index(out.Stdout, h)
		if p < 0 {
			t.Fatalf("missing %q\n%s", h, out.Stdout)
		}
		if p <= last {
			t.Fatalf("block order wrong at %q", h)
		}
		last = p
	}
	second := f.Bench("roadmap", "--context")
	if second.Stdout != out.Stdout {
		t.Fatal("unchanged invocations were not byte-identical")
	}
	deep := filepath.Join(f.Root, "a", "b")
	if err := os.MkdirAll(deep, 0755); err != nil {
		t.Fatal(err)
	}
	fromDeep := runBenchInDir(t, f, deep, "roadmap", "--context")
	fromDeep.RequireExit(0)
	if fromDeep.Stdout != out.Stdout {
		t.Fatal("deep-CWD output differs")
	}
}

func testRoadmapContextTruncation(t *testing.T) {
	f := contextFixture(t)
	f.WriteFile("IDEAS.md", "- 2026-07-10  "+strings.Repeat("x", 4097)+"\n")
	short := f.Bench("roadmap", "--context")
	short.RequireExit(0)
	short.RequireContains(short.Stdout, "4097,true")
	full := f.Bench("roadmap", "--context", "--full")
	full.RequireExit(0)
	full.RequireContains(full.Stdout, "4097,false")
}

func testRoadmapContextUsage(t *testing.T) {
	f := contextFixture(t)
	for _, a := range []string{"-h", "--help"} {
		p := f.Bench("roadmap", a)
		p.RequireExit(0)
		p.RequireContains(p.Stdout, "usage:")
		if p.Stderr != "" {
			t.Fatal("help wrote stderr")
		}
	}
	bad := f.Bench("roadmap", "--context", "--wat")
	bad.RequireExit(2)
	bad.RequireContains(bad.Stdout, "usage:")
	if bad.Stderr != "" {
		t.Fatal("usage wrote stderr")
	}
	for _, args := range [][]string{{"roadmap", "--full"}, {"roadmap", "--context", "--context"}} {
		p := f.Bench(args...)
		p.RequireExit(2)
		if p.Stderr != "" {
			t.Fatalf("%v wrote stderr", args)
		}
	}
	// Help wins wherever it appears before `--`, so `--help --full` is a help request
	// rather than a misuse.
	helpFirst := f.Bench("roadmap", "--help", "--full")
	helpFirst.RequireExit(0)
	helpFirst.RequireContains(helpFirst.Stdout, "usage:")
	if helpFirst.Stderr != "" {
		t.Fatal("help wrote stderr")
	}
	no := contract.NewFixture(t, contract.WithNoRepo()).Bench("roadmap", "--context")
	no.RequireExit(1)
	no.RequireContains(no.Stdout, "error:")
}

func testRoadmapContextSourceStates(t *testing.T) {
	f := contract.NewFixture(t)
	f.Git("branch", "-M", "main")
	f.WriteFile("seed", "x")
	f.Git("add", "seed")
	f.Git("-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-qm", "seed")
	absent := f.Bench("roadmap", "--context")
	absent.RequireExit(0)
	absent.RequireContains(absent.Stdout, "ROADMAP.md,absent,0")
	absent.RequireContains(absent.Stdout, "IDEAS.md,absent,0")
	f.WriteFile("ROADMAP.md", "")
	f.WriteFile("IDEAS.md", "")
	empty := f.Bench("roadmap", "--context")
	empty.RequireExit(0)
	empty.RequireContains(empty.Stdout, "ROADMAP.md,empty,0")
	empty.RequireContains(empty.Stdout, "IDEAS.md,empty,0")
	f.WriteFile("ROADMAP.md", "# Roadmap\n\n**FT1 — valid.** body\n\n**broken row\n")
	mixed := f.Bench("roadmap", "--context")
	mixed.RequireExit(0)
	mixed.RequireContains(mixed.Stdout, "ROADMAP.md,malformed")
	mixed.RequireContains(mixed.Stdout, "malformed roadmap row")
}

func testRoadmapContextFailClosed(t *testing.T) {
	f := contextFixture(t)
	fake := filepath.Join(f.Root, "fakebin")
	if err := os.MkdirAll(fake, 0755); err != nil {
		t.Fatal(err)
	}
	real, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	writeExecutableAt(t, fake, "git", "#!/bin/sh\nfor a in \"$@\"; do [ \"$a\" = rev-list ] && exit 17; done\nexec "+sanitize.ShellQuote(real)+" \"$@\"\n")
	p := f.RunEnv(map[string]string{"PATH": fake + string(os.PathListSeparator) + os.Getenv("PATH")}, "bash", filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh"), "roadmap", "--context")
	p.RequireExit(1)
	p.RequireContains(p.Stdout, "error: roadmap context failed")
	if p.Stderr != "" {
		t.Fatalf("git failure wrote stderr: %s", p.Stderr)
	}
	if strings.Contains(p.Stdout, "context[1]") {
		t.Fatal("git failure leaked partial snapshot")
	}
}

func testRoadmapContextReadOnlyOffline(t *testing.T) {
	f := contextFixture(t)
	fake := filepath.Join(f.Root, "sentinels")
	if err := os.MkdirAll(fake, 0755); err != nil {
		t.Fatal(err)
	}
	log := filepath.Join(f.Root, "called.log")
	for _, name := range []string{"bench", "curl", "wget", "gh", "glab", "claude", "codex", "opencode"} {
		writeExecutableAt(t, fake, name, "#!/bin/sh\nprintf '%s\\n' "+sanitize.ShellQuote(name)+" >>\"$SENTINEL_LOG\"\nexit 99\n")
	}
	f.WriteFile(".bench/gate.sh", "#!/bin/sh\nprintf gate >>\"$SENTINEL_LOG\"\nexit 99\n")
	if err := os.Chmod(filepath.Join(f.Root, ".bench", "gate.sh"), 0755); err != nil {
		t.Fatal(err)
	}
	statusBefore := f.Run("git", "status", "--porcelain=v1").Stdout
	headBefore := f.Run("git", "rev-parse", "HEAD").Stdout
	indexBefore, err := os.ReadFile(filepath.Join(f.Root, ".git", "index"))
	if err != nil {
		t.Fatal(err)
	}
	p := f.RunEnv(map[string]string{"PATH": fake + string(os.PathListSeparator) + os.Getenv("PATH"), "SENTINEL_LOG": log}, "bash", filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh"), "roadmap", "--context")
	p.RequireExit(0)
	if _, err := os.Stat(log); !os.IsNotExist(err) {
		t.Fatalf("offline sentinel invoked: %v", err)
	}
	if got := f.Run("git", "status", "--porcelain=v1").Stdout; got != statusBefore {
		t.Fatalf("status changed: %q -> %q", statusBefore, got)
	}
	if got := f.Run("git", "rev-parse", "HEAD").Stdout; got != headBefore {
		t.Fatal("HEAD changed")
	}
	indexAfter, err := os.ReadFile(filepath.Join(f.Root, ".git", "index"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(indexBefore, indexAfter) {
		t.Fatal("index bytes changed")
	}
	gitDir := strings.TrimSpace(f.Run("git", "rev-parse", "--absolute-git-dir").Stdout)
	if _, err := os.Stat(filepath.Join(gitDir, "bench-last-gate")); !os.IsNotExist(err) {
		t.Fatal("gate cache changed")
	}
}

// testRoadmapContextUnbornHead pins the same degradation posture for a repository with no
// commits: the branch HEAD points at is still named, the cells the missing commit makes
// unknowable read unknown, and every other block still arrives.
func testRoadmapContextUnbornHead(t *testing.T) {
	f := contract.NewFixture(t)
	f.Git("checkout", "-q", "-b", "trunk")
	f.WriteFile("ROADMAP.md", "# Roadmap\n\n## Features\n\n**FT1 — one.** Body specs/one/spec.md.\n")
	f.WriteFile("IDEAS.md", "- 2026-07-10  retain me\n")

	out := f.Bench("roadmap", "--context")

	out.RequireExit(0)
	if out.Stderr != "" {
		t.Fatalf("stderr not empty: %s", out.Stderr)
	}
	out.RequireContains(out.Stdout, "git[1]{branch,default_branch,dirty,ahead,behind}:")
	out.RequireContains(out.Stdout, "trunk,unknown,true,unknown,unknown")
	out.RequireContains(out.Stdout, "retain me")
	out.RequireContains(out.Stdout, "roadmap_rows[1]{")
}

// TestAXIRoadmapContextGitUnknown pins the degradation posture at the git block: an
// unresolvable default costs the three cells it makes unknowable, not the snapshot.
func TestAXIRoadmapContextGitUnknown(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contextFixture(t)
	f.Git("branch", "-M", "master")
	f.Git("branch", "feature")

	out := f.Bench("roadmap", "--context")

	out.RequireExit(0)
	out.RequireContains(out.Stdout, "git[1]{branch,default_branch,dirty,ahead,behind}:")
	out.RequireContains(out.Stdout, "master,unknown,false,unknown,unknown")
	out.RequireContains(out.Stdout, "retain me")
	out.RequireContains(out.Stdout, "roadmap_rows[1]{")
}
