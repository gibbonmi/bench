package contract

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDoctorContracts(t *testing.T) {
	skipIfSubjectBenchMissing(t)
	t.Run("bench doctor report contract", testDoctorReport)
	t.Run("bench doctor manifest skew contract", testDoctorManifestSkew)
	t.Run("bench doctor --fix write contract", testDoctorFixWrite)
	t.Run("bench doctor --fix spaced-target contract", testDoctorFixSpacedTarget)
	t.Run("bench doctor --fix idempotency contract", testDoctorFixIdempotency)
	t.Run("bench doctor --fix foreign-refuse contract", testDoctorFixForeignRefuse)
	t.Run("bench doctor --fix fallback contract", testDoctorFixFallback)
	t.Run("bench doctor --fix path-notice contract", testDoctorFixPathNotice)
	t.Run("bench doctor shim stale-target contract", testDoctorShimStaleTarget)
	t.Run("bench doctor shim arg-passthrough contract", testDoctorShimArgPassthrough)
	t.Run("bench doctor postinstall contract", testDoctorPostinstall)
	t.Run("bench doctor session-start advisory contract", testDoctorSessionStartAdvisory)
}

type doctorSandbox struct {
	home   string
	nvmBin string
	plain  string
	path   string
	env    map[string]string
}

func testDoctorReport(t *testing.T) {
	f := NewFixture(t)
	sb := newDoctorSandbox(t, f)
	targetPath := filepath.Join(sb.plain, "bench")

	probe := f.BenchEnv(sb.env, "doctor")
	probe.RequireExit(1)
	requirePathAbsent(t, targetPath, "report on a missing shim wrote a file")
	probe.RequireContains(probe.Stdout, "missing")
	probe.RequireContains(probe.Stdout, `rm -f "`+targetPath+`"`)

	f.BenchEnv(sb.env, "doctor", "--fix").RequireExit(0)
	f.BenchEnv(sb.env, "doctor").RequireExit(0)

	mustWriteFile(t, targetPath, "not a bench shim\n", 0o644)
	f.BenchEnv(sb.env, "doctor").RequireExit(1)

	mustWriteFile(t, targetPath, "#!/usr/bin/env bash\n# bench-shim v1 marker\ntarget=/no/such/bench\nexec \"$target\" \"$@\"\n", 0o755)
	probe = f.BenchEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "stale")

	mustWriteFile(t, targetPath, "#!/usr/bin/env bash\n# bench-shim v1 marker\ntarget=$(touch \""+filepath.Join(f.Root, "pwned")+"\")\nexec \"$target\"\n", 0o755)
	_ = f.BenchEnv(sb.env, "doctor")
	requirePathAbsent(t, filepath.Join(f.Root, "pwned"), "report executed a hostile shim's target line")
}

func testDoctorManifestSkew(t *testing.T) {
	f := NewFixture(t)
	sb := newDoctorSandbox(t, f)

	f.WriteFile(".bench/link-manifest.tsv", ".bench/BENCH.md\tabc\n")
	probe := f.BenchEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "version unknown")

	f.WriteFile(".bench/link-manifest.tsv", "#kit\t0.0.0\n.bench/BENCH.md\tabc\n")
	probe = f.BenchEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "version skew")
	probe.RequireContains(probe.Stdout, "0.0.0")
}

func testDoctorFixWrite(t *testing.T) {
	f := NewFixture(t)
	sb := newDoctorSandbox(t, f)
	targetPath := filepath.Join(sb.plain, "bench")

	probe := f.BenchEnv(sb.env, "doctor", "--fix")

	probe.RequireExit(0)
	requirePathExists(t, targetPath, "--fix wrote no shim in the plain PATH dir")
	requirePathAbsent(t, filepath.Join(sb.nvmBin, "bench"), "--fix wrote into the manager-owned nvm dir")
	doctorRequireFileContains(t, targetPath, "bench-shim v1", "shim carries no bench marker")
	doctorRequireFileContains(t, targetPath, filepath.Join(SubjectRoot(t), "bin", "bench.sh"), "shim does not target the resolved CLI")
	probe.RequireContains(probe.Stdout, targetPath)
}

func testDoctorFixSpacedTarget(t *testing.T) {
	f := NewFixture(t, WithSpacePath())
	sb := newDoctorSandbox(t, f)
	kit := filepath.Join(f.Root, "kit")
	copyDoctorKit(t, kit)
	targetPath := filepath.Join(sb.plain, "bench")

	probe := f.RunEnv(sb.env, "bash", filepath.Join(kit, "bin", "bench.sh"), "doctor", "--fix")

	probe.RequireExit(0)
	requirePathExists(t, targetPath, "spaced-path --fix wrote no shim")
	doctorRequireFileContains(t, targetPath, filepath.Join(kit, "bin", "bench.sh"), "shim lost the spaced target path")
	probe = f.Run(targetPath, "doctor")
	probe.RequireContains(probe.Stdout, "shim health")
}

func testDoctorFixIdempotency(t *testing.T) {
	f := NewFixture(t)
	sb := newDoctorSandbox(t, f)
	targetPath := filepath.Join(sb.plain, "bench")

	f.BenchEnv(sb.env, "doctor", "--fix").RequireExit(0)
	before := mustReadFile(t, targetPath)
	probe := f.BenchEnv(sb.env, "doctor", "--fix")

	probe.RequireExit(0)
	doctorRequireEqual(t, mustReadFile(t, targetPath), before, "second --fix rewrote an already-current shim")
	if !strings.Contains(strings.ToLower(probe.Stdout), "no change") {
		t.Fatalf("second --fix did not announce a no-op\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
}

func testDoctorFixForeignRefuse(t *testing.T) {
	noteContractFailure(t, "over a foreign file did not exit 1")
	f := NewFixture(t)
	sb := newDoctorSandbox(t, f)
	foreign := filepath.Join(sb.plain, "bench")

	mustWriteFile(t, foreign, "foreign contents\n", 0o644)
	before := mustReadFile(t, foreign)
	probe := f.BenchEnv(sb.env, "doctor", "--fix")
	probe.RequireExit(1)
	doctorRequireEqual(t, mustReadFile(t, foreign), before, "--fix clobbered a foreign file")
	if !strings.Contains(strings.ToLower(probe.Stderr), "refus") {
		t.Fatalf("--fix did not report the refusal\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}

	mustWriteFile(t, foreign, "", 0o644)
	f.BenchEnv(sb.env, "doctor", "--fix").RequireExit(1)
	info, err := os.Stat(foreign)
	if err != nil {
		t.Fatalf("stat present-but-empty foreign file: %v", err)
	}
	if info.Size() != 0 {
		t.Fatal("--fix wrote over the present-but-empty file")
	}

	mustWriteFile(t, foreign, "#!/usr/bin/env bash\n# bench-shim v1 marker\ntarget="+filepath.Join(SubjectRoot(t), "bin", "bench.sh")+"\nexec \"$target\" \"$@\"", 0o755)
	f.BenchEnv(sb.env, "doctor", "--fix").RequireExit(0)
}

func testDoctorFixFallback(t *testing.T) {
	f := NewFixture(t)
	sb := newDoctorSandbox(t, f)
	env := sb.managerOnlyEnv()

	probe := f.BenchEnv(env, "doctor", "--fix")

	probe.RequireExit(0)
	requirePathExists(t, filepath.Join(sb.home, ".local", "bin", "bench"), "fallback did not write to ~/.local/bin")
	probe.RequireContains(probe.Stdout, "created directory "+filepath.Join(sb.home, ".local", "bin"))
}

func testDoctorFixPathNotice(t *testing.T) {
	f := NewFixture(t)
	sb := newDoctorSandbox(t, f)
	env := sb.managerOnlyEnv()

	probe := f.BenchEnv(env, "doctor", "--fix")

	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "export PATH")
	probe.RequireContains(probe.Stdout, ".local/bin")
	requirePathAbsent(t, filepath.Join(sb.home, ".bashrc"), "--fix edited an rc file")
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

func newDoctorSandbox(t *testing.T, f Fixture) doctorSandbox {
	t.Helper()
	home := filepath.Join(f.Root, "home")
	nvm := filepath.Join(f.Root, "nvm")
	nvmBin := filepath.Join(nvm, "versions", "node", "v22", "bin")
	plain := filepath.Join(f.Root, "plain")
	for _, dir := range []string{home, nvmBin, plain} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("create doctor sandbox dir %s: %v", dir, err)
		}
	}
	path := nvmBin + string(os.PathListSeparator) + plain + string(os.PathListSeparator) + "/usr/bin" + string(os.PathListSeparator) + "/bin"
	sb := doctorSandbox{home: home, nvmBin: nvmBin, plain: plain, path: path}
	sb.env = map[string]string{
		"HOME":    home,
		"SHELL":   "/bin/bash",
		"NVM_DIR": nvm,
		"PATH":    path,
	}
	return sb
}

func (sb doctorSandbox) managerOnlyEnv() map[string]string {
	return map[string]string{
		"HOME":    sb.home,
		"SHELL":   "/bin/bash",
		"NVM_DIR": filepath.Dir(filepath.Dir(filepath.Dir(sb.nvmBin))),
		"PATH":    sb.nvmBin + string(os.PathListSeparator) + "/usr/bin" + string(os.PathListSeparator) + "/bin",
	}
}

func (sb doctorSandbox) with(k, v string) map[string]string {
	out := make(map[string]string, len(sb.env)+1)
	for key, value := range sb.env {
		out[key] = value
	}
	out[k] = v
	return out
}

func copyDoctorKit(t *testing.T, dst string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dst, "bin"), 0o755); err != nil {
		t.Fatalf("create kit bin: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(dst, "dist"), 0o755); err != nil {
		t.Fatalf("create kit dist: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(dst, ".agents", "commands"), 0o755); err != nil {
		t.Fatalf("create kit commands dir: %v", err)
	}
	copyBenchScripts(t, filepath.Join(dst, "bin"))
	doctorCopyFileIfExists(t, filepath.Join(SubjectRoot(t), "dist", "bench"), filepath.Join(dst, "dist", "bench"), 0o755)
	doctorCopyFileIfExists(t, filepath.Join(SubjectRoot(t), "AGENTS.md"), filepath.Join(dst, "AGENTS.md"), 0o644)
}

func copyBenchScripts(t *testing.T, dst string) {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(SubjectRoot(t), "bin"))
	if err != nil {
		t.Fatalf("read kit bin: %v", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sh") {
			continue
		}
		doctorCopyFile(t, filepath.Join(SubjectRoot(t), "bin", entry.Name()), filepath.Join(dst, entry.Name()), 0o755)
	}
}

func doctorCopyFileIfExists(t *testing.T, src, dst string, mode os.FileMode) {
	t.Helper()
	if _, err := os.Stat(src); err != nil {
		if os.IsNotExist(err) {
			return
		}
		t.Fatalf("stat %s: %v", src, err)
	}
	doctorCopyFile(t, src, dst, mode)
}

func doctorCopyFile(t *testing.T, src, dst string, mode os.FileMode) {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read %s: %v", src, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(dst), err)
	}
	if err := os.WriteFile(dst, data, mode); err != nil {
		t.Fatalf("write %s: %v", dst, err)
	}
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

func mustReadFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func mustWriteFile(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func requirePathExists(t *testing.T, path, msg string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

func requirePathAbsent(t *testing.T, path, msg string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Fatal(msg)
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat %s: %v", path, err)
	}
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

func doctorRequireEqual(t *testing.T, got, want, msg string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %q, want %q", msg, got, want)
	}
}

func doctorRequireFileContains(t *testing.T, path, needle, msg string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !strings.Contains(string(data), needle) {
		t.Fatalf("%s: missing %q in %s", msg, needle, path)
	}
}
