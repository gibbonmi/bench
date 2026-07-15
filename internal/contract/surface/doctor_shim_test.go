package surface

import (
	"github.com/gibbonmi/bench/internal/adopt"
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDoctorShimContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench doctor shim stale-target contract", testDoctorShimStaleTarget)
	contract.RunParallel(t, "bench doctor shim arg-passthrough contract", testDoctorShimArgPassthrough)
	contract.RunParallel(t, "repo-local wrapper forwarding contract failed", testDoctorShimRepoLocalForwarding)
	contract.RunParallel(t, "maintenance command escaped installed target", testDoctorShimMaintenanceMatrix)
	contract.RunParallel(t, "bench doctor postinstall contract", testDoctorPostinstall)
	contract.RunParallel(t, "bench doctor session-start advisory contract", testDoctorSessionStartAdvisory)
}

func testDoctorShimMaintenanceMatrix(t *testing.T) {
	f := contract.NewFixture(t)
	installed := filepath.Join(f.Root, "installed target")
	local := filepath.Join(f.Root, ".bench", "bin", "bench.sh")
	shim := filepath.Join(f.Root, "stable bench")
	mustWriteFile(t, installed, "#!/bin/sh\nprintf 'installed:%s\\n' \"$1\"\nexit 7\n", 0o755)
	mustWriteFile(t, local, "#!/bin/sh\nprintf 'local:%s\\n' \"$1\"\nexit 3\n", 0o755)
	mustWriteFile(t, shim, adopt.ShimContent(installed)+"\n", 0o755)

	for _, command := range []string{"setup", "link", "init", "doctor", "unlink"} {
		probe := f.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, shim, command)
		probe.RequireExit(7)
		if probe.Stdout != "installed:"+command+"\n" {
			t.Fatalf("maintenance command %s escaped installed target: %q", command, probe.Stdout)
		}
	}
	for _, args := range [][]string{{}, {"--context"}, {"status"}, {"unknown"}} {
		probe := f.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, shim, args...)
		probe.RequireExit(3)
		want := "local:\n"
		if len(args) > 0 {
			want = "local:" + args[0] + "\n"
		}
		if probe.Stdout != want {
			t.Fatalf("default shim dispatch %v = %q, want %q", args, probe.Stdout, want)
		}
	}
}

func testDoctorShimRepoLocalForwarding(t *testing.T) {
	contract.NoteContractFailure(t, "repo-local wrapper forwarding contract failed")
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	f.BenchEnv(sb.env, "doctor", "--fix").RequireExit(0)
	global := filepath.Join(f.Root, "global")
	mustWriteFile(t, global, "#!/bin/sh\nprintf 'global:%s\\n' \"$*\"\nexit 7\n", 0o755)
	local := filepath.Join(f.Root, ".bench", "bin", "bench.sh")
	mustWriteFile(t, local, "#!/bin/sh\nprintf 'local:%s|%s\\n' \"$1\" \"$2\"\nexit 2\n", 0o755)
	shim := filepath.Join(f.Root, "shim")
	content := rewriteShimTarget(t, filepath.Join(sb.plain, "bench"), global)
	mustWriteFile(t, shim, content, 0o755)

	probe := f.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, shim, "--context", "--full")
	probe.RequireExit(2)
	doctorRequireEqual(t, probe.Stdout, "local:--context|--full\n", "repo-local wrapper forwarding contract failed")
	if probe.Stderr != "" {
		t.Fatalf("repo-local wrapper wrote stderr: %s", probe.Stderr)
	}

	requireRemove(t, local)
	probe = f.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, shim, "--context", "--full")
	probe.RequireExit(7)
	doctorRequireEqual(t, probe.Stdout, "global:--context --full\n", "missing local wrapper did not fall back")
}

func testDoctorShimStaleTarget(t *testing.T) {
	contract.NoteContractFailure(t, "stale-target shim printed no remedy")
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	f.BenchEnv(sb.env, "doctor", "--fix").RequireExit(0)
	shim := filepath.Join(f.Root, "shim")
	content := rewriteShimTarget(t, filepath.Join(sb.plain, "bench"), "/no/such/bench")
	mustWriteFile(t, shim, content, 0o755)

	probe := f.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, shim, "help")

	probe.RequireExit(127)
	probe.RequireContains(probe.Stderr, "bench moved")
}

func testDoctorShimArgPassthrough(t *testing.T) {
	contract.NoteContractFailure(t, "shim mangled the args")
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	f.BenchEnv(sb.env, "doctor", "--fix").RequireExit(0)
	stub := filepath.Join(f.Root, "stub")
	mustWriteFile(t, stub, "#!/usr/bin/env bash\nfor a in \"$@\"; do printf \"[%s]\" \"$a\"; done\necho\n", 0o755)
	shim := filepath.Join(f.Root, "shim")
	content := rewriteShimTarget(t, filepath.Join(sb.plain, "bench"), stub)
	mustWriteFile(t, shim, content, 0o755)

	probe := f.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, shim, "a b", "*", "c")

	probe.RequireExit(0)
	doctorRequireEqual(t, strings.TrimSpace(probe.Stdout), "[a b][*][c]", "shim mangled the args")
}

func testDoctorPostinstall(t *testing.T) {
	contract.NoteContractFailure(t, "postinstall exited nonzero on a write failure")
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	pkg := filepath.Join(f.Root, "pkg")
	copyDoctorKit(t, pkg)
	pin := filepath.Join(pkg, "bin", "bench-postinstall.sh")
	targetPath := filepath.Join(sb.plain, "bench")

	env := sb.with("npm_config_global", "true")
	probe := f.RunEnv(env, "bash", pin)
	probe.RequireExit(0)
	requirePathExists(t, targetPath, "postinstall did not write the shim on a global install")

	requireRemove(t, targetPath)
	mustWriteFile(t, filepath.Join(pkg, ".git"), "", 0o644)
	probe = f.RunEnv(env, "bash", pin)
	probe.RequireExit(0)
	requirePathAbsent(t, targetPath, "postinstall mutated with .git present")
	probe.RequireContains(probe.Stdout, "doctor --fix")

	requireRemove(t, filepath.Join(pkg, ".git"))
	probe = f.RunEnv(sb.env, "bash", pin)
	probe.RequireExit(0)
	requirePathAbsent(t, targetPath, "postinstall mutated without npm_config_global")
	probe.RequireContains(probe.Stdout, "doctor --fix")

	readOnlyHome := filepath.Join(f.Root, "rohome")
	if err := os.MkdirAll(readOnlyHome, 0o755); err != nil {
		t.Fatalf("create read-only home: %v", err)
	}
	if err := os.Chmod(readOnlyHome, 0o555); err != nil {
		t.Fatalf("chmod read-only home: %v", err)
	}
	defer func() { _ = os.Chmod(readOnlyHome, 0o755) }()
	probe = f.RunEnv(map[string]string{
		"HOME":              readOnlyHome,
		"SHELL":             "/bin/bash",
		"NVM_DIR":           filepath.Join(f.Root, "nvm"),
		"PATH":              sb.nvmBin + ":/usr/bin:/bin",
		"npm_config_global": "true",
	}, "bash", pin)
	probe.RequireExit(0)
}

func testDoctorSessionStartAdvisory(t *testing.T) {
	contract.NoteContractFailure(t, "by-path session-start advisory omits the doctor pointer")
	f := contract.NewFixture(t, contract.WithNoRepo())
	repo := filepath.Join(f.Root, "repo")
	if err := os.MkdirAll(filepath.Join(repo, "bin"), 0o755); err != nil {
		t.Fatalf("create repo bin: %v", err)
	}
	copyBenchScripts(t, filepath.Join(repo, "bin"))
	if err := os.MkdirAll(filepath.Join(repo, "dist"), 0o755); err != nil {
		t.Fatalf("create repo dist: %v", err)
	}
	copyFileTo(t, filepath.Join(contract.SubjectRoot(t), "dist", "bench"), filepath.Join(repo, "dist", "bench"))
	repoFixture := contract.NewFixtureAt(t, repo, contract.IsolatedEnv(t, repo))
	repoFixture.Run("git", "init", "-q").RequireExit(0)
	hook := filepath.Join(contract.SubjectRoot(t), ".bench", "hooks", "session-start.sh")

	probe := repoFixture.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, "bash", hook)

	probe.RequireExit(0)
	probe.RequireContains(doctorFirstLine(probe.Stdout), "doctor --fix")
}

func rewriteShimTarget(t *testing.T, path, target string) string {
	t.Helper()
	lines := strings.Split(mustReadFile(t, path), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "target=") {
			lines[i] = "target=" + shellSingleQuote(target)
			return strings.Join(lines, "\n")
		}
	}
	t.Fatalf("shim has no target= line: %s", path)
	return ""
}

func shellSingleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func requireRemove(t *testing.T, path string) {
	t.Helper()
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		t.Fatalf("remove %s: %v", path, err)
	}
}

func doctorFirstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
