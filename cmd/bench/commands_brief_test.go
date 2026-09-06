// The commands --brief probe family: the running executable is graded against the
// checkout it answers for, out of process, in a synthetic checkout under t.TempDir().
package main

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// TestCommandsBriefGradesTheRunningExecutable covers BF1 and BF2. The probe is the one
// verb a session runs to learn whether Bench answers at all, so a stale binary answering
// it is a false pass nothing else in the loop catches. The rows run out of process,
// because the verdict is about the executable that is running, and a test binary is not
// any checkout's publication.
//
// Each row grades a synthetic checkout under t.TempDir(): a git repository that declares
// build inputs and holds the published pair at dist/bench. Nothing here writes into the
// repository under test, which runs its own package sweeps concurrently over the same
// tree. The real bench binary is what answers; only the sources it is graded against are
// synthetic, and the handler resolves that checkout itself from the child's directory.
func TestCommandsBriefGradesTheRunningExecutable(t *testing.T) {
	staged := commandsBriefBinary(t)

	t.Run("a matching seal keeps the three-line answer", func(t *testing.T) {
		root, binary := commandsBriefCheckout(t, staged)
		out, code := commandsBriefProbe(t, root, binary)
		if code != 0 {
			t.Fatalf("commands --brief under a matching seal = exit %d:\n%s", code, out)
		}
		if want, _ := commandsCommand([]string{"--brief"}); out != want {
			t.Fatalf("commands --brief = %q, want the probe answer %q", out, want)
		}
	})

	t.Run("a mismatched seal refuses with the rebuild action", func(t *testing.T) {
		root, binary := commandsBriefCheckout(t, staged)
		commandsBriefBreakSeal(t, binary)
		out, code := commandsBriefProbe(t, root, binary)
		if code == 0 {
			t.Fatalf("commands --brief under a mismatched seal exited 0:\n%s", out)
		}
		if action := freshness.RebuildAction(root); !strings.Contains(out, action) {
			t.Fatalf("commands --brief refusal = %q, want the rebuild action %q", out, action)
		}
	})
}

// commandsBriefBinary compiles the subject once into a temporary directory. The rows need
// the real handler, not the builder's stamping, so a plain build is the whole requirement
// and the publication below binds these bytes to the checkout that grades them.
func commandsBriefBinary(t *testing.T) string {
	t.Helper()
	staged := filepath.Join(t.TempDir(), "bench")
	// The probe subject is a main package built from the tree, so Go would stamp VCS
	// state and walk up from the checkout to any stray .git directory above it. The
	// gate owns the one defense, and the build composes it.
	cmd := exec.Command("go", "build", gate.DisableBuildVCS, "-o", staged, "./cmd/bench")
	cmd.Dir = filepath.Join("..", "..")
	cmd.Env = capability.WithoutEnvironment(os.Environ(), runbinary.Env)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build the probe subject: %v\n%s", err, output)
	}
	return staged
}

// commandsBriefCheckout publishes staged as one synthetic checkout's own dist/bench and
// returns the checkout with the published executable. The checkout carries the three
// things the grading branch reads: a git top level for the child's directory, the
// build-input manifest that scopes the verdict, and a resolvable ./cmd/bench for the
// digest. The returned root is spelled as git reports it, because the refusal's rebuild
// action is built from that spelling.
func commandsBriefCheckout(t *testing.T, staged string) (string, string) {
	t.Helper()
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	kit := filepath.Join(base, "kit")
	for path, body := range map[string]string{
		filepath.Join(kit, "go.mod"):                  "module benchprobe\n\ngo 1.21\n",
		filepath.Join(kit, "cmd", "bench", "main.go"): "package main\n\nfunc main() {}\n",
		filepath.Join(kit, "scripts", "go-build.sh"):  "# the probe checkout's build entry point\n",
		// One auxiliary input is enough; the manifest's own presence is what
		// DeclaresBuildInputs reads, and its entries are what the digest covers.
		filepath.Join(kit, "scripts", "go-build.inputs"): "build_script=scripts/go-build.sh\n",
	} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gitInitRepo(t, kit)
	root, err := git.RootAt(kit)
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(root, "dist", "bench")
	if err := os.MkdirAll(filepath.Dir(binary), 0o755); err != nil {
		t.Fatal(err)
	}
	// Publish consumes the staged bytes by renaming them, so each row publishes its own
	// copy and neither depends on the order the other ran in.
	copyExecutable(t, staged, filepath.Join(base, "staged"))
	if err := freshness.Publish(root, filepath.Join(base, "staged"), binary, filepath.Dir(binary), "probe"); err != nil {
		t.Fatal(err)
	}
	return root, binary
}

// commandsBriefBreakSeal rewrites the source digest the seal records, which is the exact
// state a rebuilt tree beside an unrebuilt binary leaves behind.
func commandsBriefBreakSeal(t *testing.T, binary string) {
	t.Helper()
	body, err := os.ReadFile(binary + ".seal")
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(body, &fields); err != nil {
		t.Fatal(err)
	}
	sources, ok := fields["sources"].(string)
	if !ok || sources == "" {
		t.Fatalf("seal records no source digest: %s", body)
	}
	fields["sources"] = strings.Repeat("0", len(sources))
	edited, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binary+".seal", edited, 0o644); err != nil {
		t.Fatal(err)
	}
}

func commandsBriefProbe(t *testing.T, root, binary string) (string, int) {
	t.Helper()
	cmd := exec.Command(binary, "commands", "--brief")
	cmd.Dir = root
	cmd.Env = capability.WithoutEnvironment(os.Environ(), runbinary.Env)
	output, err := cmd.CombinedOutput()
	code := 0
	if err != nil {
		exit := &exec.ExitError{}
		if !errors.As(err, &exit) {
			t.Fatalf("run %s commands --brief: %v\n%s", binary, err, output)
		}
		code = exit.ExitCode()
	}
	return string(output), code
}
