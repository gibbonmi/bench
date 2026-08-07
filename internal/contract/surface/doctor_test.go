package surface

import (
	"github.com/gibbonmi/bench/internal/adopt"
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDoctorReportFixContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench doctor report contract", testDoctorReport)
	contract.RunParallel(t, "bench doctor manifest skew contract", testDoctorManifestSkew)
	contract.RunParallel(t, "bench doctor --fix write contract", testDoctorFixWrite)
	contract.RunParallel(t, "bench doctor --fix spaced-target contract", testDoctorFixSpacedTarget)
	contract.RunParallel(t, "bench doctor --fix idempotency contract", testDoctorFixIdempotency)
	contract.RunParallel(t, "bench doctor --fix foreign-refuse contract", testDoctorFixForeignRefuse)
	contract.RunParallel(t, "bench doctor --fix fallback contract", testDoctorFixFallback)
	contract.RunParallel(t, "bench doctor --fix path-notice contract", testDoctorFixPathNotice)
	contract.RunParallel(t, "bench doctor --fix path-notice already-current contract", testDoctorFixPathNoticeAlreadyCurrent)
	contract.RunParallel(t, "bench doctor pre-push row contract", testDoctorPrePushRow)
	contract.RunParallel(t, "bench doctor pre-push health rendering contract", testDoctorPrePushHealthRendering)
	contract.RunParallel(t, "bench doctor pre-push kit-repo contract", testDoctorPrePushKitRepo)
	contract.RunParallel(t, "bench doctor --fix stale pre-push repair contract", testDoctorFixStalePrePush)
	contract.RunParallel(t, "bench doctor --fix stale pre-push execute-mode contract", testDoctorFixStalePrePushExecuteMode)
}

// managedPrePushBody is a minimal hook carrying the bench:managed-pre-push marker that
// InspectPrePush fingerprints — the one source the doctor pre-push fixtures write for a
// bench-managed hook.
const managedPrePushBody = "#!/usr/bin/env bash\n# bench:managed-pre-push\nexit 0\n"

type doctorSandbox struct {
	home   string
	nvmBin string
	plain  string
	path   string
	env    map[string]string
}

func testDoctorReport(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	targetPath := filepath.Join(sb.plain, "bench")
	// A bench-managed pre-push keeps the backstop row green, so this test's exit codes
	// answer only for shim health — the pre-push row has its own contract below.
	mustWriteFile(t, filepath.Join(f.Root, ".git", "hooks", "pre-push"), currentPrePushBody(t, f), 0o755)
	writeDoctorGreenTree(t, f)

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
	f := contract.NewFixture(t)
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
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	targetPath := filepath.Join(sb.plain, "bench")

	probe := f.BenchEnv(sb.env, "doctor", "--fix")

	probe.RequireExit(0)
	requirePathExists(t, targetPath, "--fix wrote no shim in the plain PATH dir")
	requirePathAbsent(t, filepath.Join(sb.nvmBin, "bench"), "--fix wrote into the manager-owned nvm dir")
	doctorRequireFileContains(t, targetPath, "bench-shim v1", "shim carries no bench marker")
	doctorRequireFileContains(t, targetPath, filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh"), "shim does not target the resolved CLI")
	probe.RequireContains(probe.Stdout, targetPath)
}

func testDoctorFixSpacedTarget(t *testing.T) {
	f := contract.NewFixture(t, contract.WithSpacePath())
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
	f := contract.NewFixture(t)
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
	contract.NoteContractFailure(t, "over a foreign file did not exit 1")
	f := contract.NewFixture(t)
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

	mustWriteFile(t, foreign, "#!/usr/bin/env bash\n# bench-shim v1 marker\ntarget="+filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")+"\nexec \"$target\" \"$@\"", 0o755)
	f.BenchEnv(sb.env, "doctor", "--fix").RequireExit(0)
}

func testDoctorFixFallback(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	env := sb.managerOnlyEnv()

	probe := f.BenchEnv(env, "doctor", "--fix")

	probe.RequireExit(0)
	requirePathExists(t, filepath.Join(sb.home, ".local", "bin", "bench"), "fallback did not write to ~/.local/bin")
	probe.RequireContains(probe.Stdout, "created directory "+filepath.Join(sb.home, ".local", "bin"))
}

func testDoctorFixPathNotice(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	env := sb.managerOnlyEnv()

	probe := f.BenchEnv(env, "doctor", "--fix")

	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "export PATH")
	probe.RequireContains(probe.Stdout, ".local/bin")
	requirePathAbsent(t, filepath.Join(sb.home, ".bashrc"), "--fix edited an rc file")
}

// testDoctorFixPathNoticeAlreadyCurrent covers the shim-already-current branch: a second
// --fix over an unchanged shim still off PATH must repeat the PATH notice, not just the
// first --fix that wrote the shim fresh.
func testDoctorFixPathNoticeAlreadyCurrent(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	env := sb.managerOnlyEnv()

	f.BenchEnv(env, "doctor", "--fix").RequireExit(0)

	probe := f.BenchEnv(env, "doctor", "--fix")

	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "already current")
	probe.RequireContains(probe.Stdout, "export PATH")
	probe.RequireContains(probe.Stdout, ".local/bin")
	requirePathAbsent(t, filepath.Join(sb.home, ".bashrc"), "--fix edited an rc file")
}

// testDoctorPrePushRow walks the pre-push backstop states: a fresh clone's absent hook, a
// user-authored foreign hook (reported, never overwritten), a core.hooksPath divert with no
// managed hook, and the managed/boundary greens. The shim is installed first so the report's
// exit answers only for the pre-push row.
func testDoctorPrePushRow(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	f.BenchEnv(sb.env, "doctor", "--fix").RequireExit(0)
	prepush := filepath.Join(f.Root, ".git", "hooks", "pre-push")
	current := currentPrePushBody(t, f)
	if err := os.Remove(prepush); err != nil {
		t.Fatalf("remove linked pre-push: %v", err)
	}

	// Absent — git does not clone hooks, so a fresh clone drops the backstop: red row that
	// names "pre-push" and the resolved install path, doctor exit 1.
	requirePathAbsent(t, prepush, "fresh fixture unexpectedly carries a pre-push")
	probe := f.BenchEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "pre-push")
	probe.RequireContains(probe.Stdout, prepush)

	// Foreign — a present hook with no managed marker is reported "not bench-managed" and
	// left byte-identical (doctor is read-only; only bench link installs).
	foreign := "#!/usr/bin/env bash\n# a hook the user wrote\nexit 0\n"
	mustWriteFile(t, prepush, foreign, 0o755)
	probe = f.BenchEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "not bench-managed")
	doctorRequireEqual(t, mustReadFile(t, prepush), foreign, "doctor rewrote a foreign pre-push")

	// Managed — the marker fingerprints a bench-managed hook: green row, doctor exit 0.
	mustWriteFile(t, prepush, current, 0o755)
	writeDoctorGreenTree(t, f)
	f.BenchEnv(sb.env, "doctor", "--fix").RequireExit(0)
	f.BenchEnv(sb.env, "doctor").RequireExit(0)

	// Diverted — core.hooksPath pointing at a dir with no managed hook (a spaced path, per
	// the hostile-input checklist) is red "diverted", even though .git/hooks stays managed.
	divert := filepath.Join(f.Root, "diverted hooks")
	contract.Mkdir(t, divert)
	f.Git("config", "core.hooksPath", divert)
	probe = f.BenchEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "diverted")

	// Boundary — a divert dir that DOES carry a managed hook is green, not diverted.
	mustWriteFile(t, filepath.Join(divert, "pre-push"), current, 0o755)
	f.BenchEnv(sb.env, "doctor", "--fix").RequireExit(0)
	f.BenchEnv(sb.env, "doctor").RequireExit(0)
}

func testDoctorFixStalePrePush(t *testing.T) {
	contract.NoteContractFailure(t, "a stale managed pre-push did not name or complete bench doctor --fix repair")
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	current := currentPrePushBody(t, f)
	writeDoctorGreenTree(t, f)
	f.BenchEnv(sb.env, "doctor", "--fix").RequireExit(0)
	prepush := filepath.Join(f.Root, ".git", "hooks", "pre-push")
	stale := strings.Replace(current, "#!/usr/bin/env bash\n", "#!/usr/bin/env bash\n# stale managed hook\n", 1)
	mustWriteFile(t, prepush, stale, 0o755)
	beforeInfo, err := os.Stat(prepush)
	if err != nil {
		t.Fatalf("stat stale pre-push: %v", err)
	}

	plain := f.BenchEnv(sb.env, "doctor")
	plain.RequireExit(1)
	plain.RequireContains(plain.Stdout, "stale: bench-managed pre-push")
	plain.RequireContains(plain.Stdout, "current hook - run bench doctor --fix")
	doctorRequireEqual(t, mustReadFile(t, prepush), stale, "plain doctor rewrote stale pre-push")
	afterPlain, err := os.Stat(prepush)
	if err != nil {
		t.Fatalf("stat stale pre-push after plain doctor: %v", err)
	}
	if afterPlain.Mode() != beforeInfo.Mode() {
		t.Fatalf("plain doctor changed stale pre-push mode: got %o, want %o", afterPlain.Mode(), beforeInfo.Mode())
	}

	first := f.BenchEnv(sb.env, "doctor", "--fix")
	first.RequireExit(0)
	first.RequireContains(first.Stdout, "repaired stale pre-push")
	health := adopt.InspectPrePush(f.Root)
	if health.Currency != adopt.PrePushCurrent {
		t.Fatalf("--fix left managed pre-push %q, want current", health.Currency)
	}
	info, err := os.Stat(prepush)
	if err != nil {
		t.Fatalf("stat repaired pre-push: %v", err)
	}
	if info.Mode().Perm()&0o111 == 0 {
		t.Fatalf("repaired pre-push is not executable: mode %o", info.Mode().Perm())
	}
	repaired := mustReadFile(t, prepush)
	mode := info.Mode()

	second := f.BenchEnv(sb.env, "doctor", "--fix")
	second.RequireExit(0)
	second.RequireContains(second.Stdout, "pre-push already current")
	doctorRequireEqual(t, mustReadFile(t, prepush), repaired, "second --fix rewrote repaired pre-push")
	secondInfo, err := os.Stat(prepush)
	if err != nil {
		t.Fatalf("stat repaired pre-push after second --fix: %v", err)
	}
	if secondInfo.Mode() != mode {
		t.Fatalf("second --fix changed repaired pre-push mode: got %o, want %o", secondInfo.Mode(), mode)
	}

	foreign := "#!/usr/bin/env bash\n# user-owned pre-push\nexit 0\n"
	mustWriteFile(t, prepush, foreign, 0o744)
	foreignInfo, err := os.Stat(prepush)
	if err != nil {
		t.Fatalf("stat foreign pre-push: %v", err)
	}
	refusal := f.BenchEnv(sb.env, "doctor", "--fix")
	refusal.RequireExit(1)
	refusal.RequireContains(refusal.Stderr, "refusing")
	doctorRequireEqual(t, mustReadFile(t, prepush), foreign, "--fix clobbered foreign pre-push")
	foreignAfter, err := os.Stat(prepush)
	if err != nil {
		t.Fatalf("stat foreign pre-push after --fix: %v", err)
	}
	if foreignAfter.Mode() != foreignInfo.Mode() {
		t.Fatalf("--fix changed foreign pre-push mode: got %o, want %o", foreignAfter.Mode(), foreignInfo.Mode())
	}

	if err := os.Remove(prepush); err != nil {
		t.Fatalf("remove pre-push: %v", err)
	}
	f.BenchEnv(sb.env, "doctor", "--fix").RequireExit(0)
	requirePathAbsent(t, prepush, "--fix installed an absent pre-push")
}

// testDoctorFixStalePrePushExecuteMode seeds the stale hook at 0644, the mode a rewrite
// alone leaves untouched: os.WriteFile applies its mode only when creating a file, so a
// repair that only re-renders the bytes hands git back a hook it silently skips. The
// repaired hook must be both current and executable.
func testDoctorFixStalePrePushExecuteMode(t *testing.T) {
	contract.NoteContractFailure(t, "a stale managed pre-push seeded at 0644 was repaired without executable mode")
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	current := currentPrePushBody(t, f)
	writeDoctorGreenTree(t, f)
	f.BenchEnv(sb.env, "doctor", "--fix").RequireExit(0)

	prepush := filepath.Join(f.Root, ".git", "hooks", "pre-push")
	stale := strings.Replace(current, "#!/usr/bin/env bash\n", "#!/usr/bin/env bash\n# stale managed hook\n", 1)
	mustWriteFile(t, prepush, stale, 0o644)
	if err := os.Chmod(prepush, 0o644); err != nil {
		t.Fatalf("strip execute bits from stale pre-push: %v", err)
	}

	fix := f.BenchEnv(sb.env, "doctor", "--fix")
	fix.RequireExit(0)
	fix.RequireContains(fix.Stdout, "repaired stale pre-push")
	doctorRequireEqual(t, mustReadFile(t, prepush), current, "--fix left the repaired pre-push off the current body")
	info, err := os.Stat(prepush)
	if err != nil {
		t.Fatalf("stat repaired pre-push: %v", err)
	}
	if info.Mode().Perm()&0o111 == 0 {
		t.Fatalf("--fix repaired a 0644 stale pre-push without executable mode: got %o", info.Mode().Perm())
	}
}

func currentPrePushBody(t *testing.T, f contract.Fixture) string {
	t.Helper()
	f.Bench("link").RequireExit(0)
	return f.ReadFile(".git/hooks/pre-push")
}

func testDoctorPrePushHealthRendering(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	f.Bench("link").RequireExit(0)
	prepush := f.ReadFile(".git/hooks/pre-push")
	writeDoctorGreenTree(t, f)
	mustWriteFile(t, filepath.Join(f.Root, ".git", "hooks", "pre-push"), prepush, 0o755)
	f.BenchEnv(sb.env, "doctor", "--fix").RequireExit(0)

	f.Git("symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/live")
	live := f.BenchEnv(sb.env, "doctor")
	live.RequireExit(0)
	live.RequireContains(live.Stdout, "live")
	live.RequireContains(live.Stdout, "provenance live")

	f.Git("symbolic-ref", "--delete", "refs/remotes/origin/HEAD")
	baked := f.BenchEnv(sb.env, "doctor")
	baked.RequireExit(0)
	baked.RequireContains(baked.Stdout, "main")
	baked.RequireContains(baked.Stdout, "provenance baked")
}

// testDoctorPrePushKitRepo pins FT10: the same row fires on the kit repo itself. A
// kit-shaped repo carries no link manifest, so a row gated behind a manifest would skip it;
// the row must still fire and name the resolved install path.
func testDoctorPrePushKitRepo(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	f.BenchEnv(sb.env, "doctor", "--fix").RequireExit(0)
	f.WriteFile("AGENTS.md", "# Working agreement\n")
	f.WriteFile(".bench/BENCH.md", "# Bench\n")
	prepush := filepath.Join(f.Root, ".git", "hooks", "pre-push")
	requirePathAbsent(t, prepush, "kit-shaped fixture unexpectedly carries a pre-push")

	probe := f.BenchEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "pre-push")
	probe.RequireContains(probe.Stdout, prepush)
}

func newDoctorSandbox(t *testing.T, f contract.Fixture) doctorSandbox {
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
	doctorCopyFileIfExists(t, filepath.Join(contract.SubjectRoot(t), "dist", "bench"), filepath.Join(dst, "dist", "bench"), 0o755)
	doctorCopyFileIfExists(t, filepath.Join(contract.SubjectRoot(t), "AGENTS.md"), filepath.Join(dst, "AGENTS.md"), 0o644)
}

func copyBenchScripts(t *testing.T, dst string) {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(contract.SubjectRoot(t), "bin"))
	if err != nil {
		t.Fatalf("read kit bin: %v", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sh") {
			continue
		}
		doctorCopyFile(t, filepath.Join(contract.SubjectRoot(t), "bin", entry.Name()), filepath.Join(dst, entry.Name()), 0o755)
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
