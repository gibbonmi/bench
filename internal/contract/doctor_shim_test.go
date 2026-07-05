package contract

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDoctorShimContracts(t *testing.T) {
	t.Parallel()
	skipIfSubjectBenchMissing(t)
	runParallel(t, "bench doctor shim stale-target contract", testDoctorShimStaleTarget)
	runParallel(t, "bench doctor shim arg-passthrough contract", testDoctorShimArgPassthrough)
	runParallel(t, "bench doctor postinstall contract", testDoctorPostinstall)
	runParallel(t, "bench doctor session-start advisory contract", testDoctorSessionStartAdvisory)
}

func testDoctorShimStaleTarget(t *testing.T) {
	noteContractFailure(t, "stale-target shim printed no remedy")
	f := NewFixture(t)
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
	noteContractFailure(t, "shim mangled the args")
	f := NewFixture(t)
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
	noteContractFailure(t, "postinstall exited nonzero on a write failure")
	f := NewFixture(t)
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
	noteContractFailure(t, "by-path session-start advisory omits the doctor pointer")
	f := NewFixture(t, WithNoRepo())
	repo := filepath.Join(f.Root, "repo")
	if err := os.MkdirAll(filepath.Join(repo, "bin"), 0o755); err != nil {
		t.Fatalf("create repo bin: %v", err)
	}
	copyBenchScripts(t, filepath.Join(repo, "bin"))
	repoFixture := Fixture{t: t, Root: repo, Env: isolatedEnv(t, repo)}
	repoFixture.Run("git", "init", "-q").RequireExit(0)
	hook := filepath.Join(SubjectRoot(t), ".bench", "hooks", "session-start.sh")

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
