package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

const fakeGoScript = "#!/usr/bin/env bash\nprintf 'pwd=<%s>\\nargc=%d\\n' \"$PWD\" \"$#\" >> \"$BENCH_TEST_RECORD\"\nfor ((i = 1; i <= $#; i++)); do\n  printf 'arg[%d]=<%s>\\n' \"$i\" \"${!i}\" >> \"$BENCH_TEST_RECORD\"\ndone\nprintf '%s\\n' '{\"Action\":\"pass\",\"Package\":\"example/pass\"}'\n"

func TestRuntimeTestReportContracts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	contract.RunParallel(t, "bench test runs fresh Go from the repository root contract", testRuntimeTestRunsFreshGoAtRoot)
	contract.RunParallel(t, "bench test grammar preserves one package argv contract", testRuntimeTestGrammar)
	contract.RunParallel(t, "bench test keeps test-scoped no-test-looking output as pass contract", testRuntimeTestDoesNotMistakeTestLogForNoTests)
	contract.RunParallel(t, "bench test renders sorted terminal package rows contract", testRuntimeTestRendersTerminalPackages)
	contract.RunParallel(t, "bench test renders direct and build failure diagnostics contract", testRuntimeTestRendersFailures)
	contract.RunParallel(t, "bench test selects only direct failure diagnostics contract", testRuntimeTestSelectsFailureDiagnostics)
	contract.RunParallel(t, "bench test reports malformed, empty, and unavailable Go output contract", testRuntimeTestErrorPostures)
}

func testRuntimeTestRunsFreshGoAtRoot(t *testing.T) {
	f := contract.NewFixture(t)
	deep := filepath.Join(f.Root, "deep", "cwd")
	contract.Mkdir(t, deep)
	record := filepath.Join(f.Root, "go-record")
	f.WriteExecutable("go", fakeGoScript)
	env := map[string]string{
		"BENCH_TEST_RECORD": record,
		"PATH":              f.Root + string(os.PathListSeparator) + os.Getenv("PATH"),
	}
	bench := benchPath(t)

	for _, tc := range []struct {
		name, pkg string
		want      []string
	}{
		{name: "default", want: []string{"test", "-json", "-count=1", "./..."}},
		{name: "one package", pkg: "./one package*", want: []string{"test", "-json", "-count=1", "./one package*"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := []string{bench, "test"}
			if tc.pkg != "" {
				args = append(args, tc.pkg)
			}
			out := contract.RunAt(t, f, deep, env, "bash", args...)
			if out.ExitCode != 0 || out.Stderr != "" {
				t.Fatalf("bench test = exit %d stdout %q stderr %q, want success with no stderr", out.ExitCode, out.Stdout, out.Stderr)
			}
			if got := contract.ReadFileAbs(t, record); got != wantGoRecord(f.Root, tc.want...) {
				t.Fatalf("fake Go invocation = %q, want %q", got, wantGoRecord(f.Root, tc.want...))
			}
			if err := os.Remove(record); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func testRuntimeTestGrammar(t *testing.T) {
	root := filepath.Join(t.TempDir(), "repo * with space")
	contract.Mkdir(t, root)
	f := contract.NewFixtureAt(t, root, contract.IsolatedEnv(t, t.TempDir()))
	f.Git("init", "-q")
	deep := filepath.Join(root, "deep")
	contract.Mkdir(t, deep)
	record := filepath.Join(root, "go-record")
	f.WriteExecutable("go", fakeGoScript)
	env := map[string]string{"BENCH_TEST_RECORD": record, "PATH": root + string(os.PathListSeparator) + os.Getenv("PATH")}
	bench := benchPath(t)

	for _, args := range [][]string{{"test", "extra", "again"}, {"test", ""}, {"test", "--unknown"}} {
		out := contract.RunAt(t, f, deep, env, "bash", append([]string{bench}, args...)...)
		if out.ExitCode != 2 || out.Stderr != "" || !strings.HasPrefix(out.Stdout, "usage: bench test [--full] [package]") {
			t.Fatalf("bench %q = exit %d stdout %q stderr %q, want usage exit 2 on stdout", args, out.ExitCode, out.Stdout, out.Stderr)
		}
	}
	help := contract.RunAt(t, f, deep, env, "bash", bench, "test", "--help")
	if help.ExitCode != 0 || help.Stdout != "usage: bench test [--full] [package]\n" || help.Stderr != "" {
		t.Fatalf("bench test --help = exit %d stdout %q stderr %q", help.ExitCode, help.Stdout, help.Stderr)
	}
	leading := contract.RunAt(t, f, deep, env, "bash", bench, "test", "--full", "--", "-leading package*")
	if leading.ExitCode != 0 || leading.Stderr != "" {
		t.Fatalf("bench test --full -- dash-leading package = exit %d stdout %q stderr %q", leading.ExitCode, leading.Stdout, leading.Stderr)
	}
	want := wantGoRecord(root, "test", "-json", "-count=1", "-leading package*")
	if got := contract.ReadFileAbs(t, record); got != want {
		t.Fatalf("fake Go invocation = %q, want %q", got, want)
	}
}

func wantGoRecord(root string, args ...string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "pwd=<%s>\nargc=%d\n", root, len(args))
	for i, arg := range args {
		fmt.Fprintf(&b, "arg[%d]=<%s>\n", i+1, arg)
	}
	return b.String()
}

func testRuntimeTestDoesNotMistakeTestLogForNoTests(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteExecutable("go", "#!/usr/bin/env bash\nprintf '%s\\n' '{\"Action\":\"output\",\"Package\":\"example/pass\",\"Test\":\"TestLogs\",\"Output\":\"[no test files]\\\\n\"}' '{\"Action\":\"pass\",\"Package\":\"example/pass\"}'\n")
	out := f.BenchEnv(map[string]string{"PATH": f.Root + string(os.PathListSeparator) + os.Getenv("PATH")}, "test")
	if out.ExitCode != 0 || !strings.HasPrefix(out.Stdout, "packages[1]{package,status}:\n  example/pass,pass\n") {
		t.Fatalf("test-scoped no-test-looking output = exit %d stdout %q stderr %q, want passing package", out.ExitCode, out.Stdout, out.Stderr)
	}
}

func testRuntimeTestRendersTerminalPackages(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/report\n\ngo 1.25\n")
	f.WriteFile("a/pass_test.go", "package a\n\nimport \"testing\"\n\nfunc TestPass(t *testing.T) {}\n")
	f.WriteFile("b/b.go", "package b\n\nfunc broken( {\n")
	f.WriteFile("m/m.go", "package m\n")
	f.WriteFile("z/fail_test.go", "package z\n\nimport \"testing\"\n\nfunc TestFail(t *testing.T) { t.Fatal(\"deliberate\") }\n")

	out := f.Bench("test")
	if out.ExitCode != 1 || out.Stderr != "" {
		t.Fatalf("bench test mixed packages = exit %d stdout %q stderr %q, want exit 1 with no stderr", out.ExitCode, out.Stdout, out.Stderr)
	}
	want := "packages[4]{package,status}:\n  example.com/report/a,pass\n  example.com/report/b,fail\n  example.com/report/m,no-tests\n  example.com/report/z,fail\n"
	if !strings.HasPrefix(out.Stdout, want) {
		t.Fatalf("package table = %q, want %q", out.Stdout, want)
	}
	if !strings.Contains(out.Stdout, "example.com/report/z,TestFail,") || !strings.Contains(out.Stdout, "deliberate") {
		t.Fatalf("real Go failure table = %q, want TestFail diagnostic", out.Stdout)
	}
	buildLine := ""
	for _, line := range strings.Split(out.Stdout, "\n") {
		if strings.HasPrefix(line, "  example.com/report/b,\"\",") {
			buildLine = line
			break
		}
	}
	if buildLine == "" || strings.Contains(buildLine, "# example.com/report/b") || strings.Contains(buildLine, "no diagnostic emitted") {
		t.Fatalf("real Go build failure = %q, want an actionable compiler diagnostic", buildLine)
	}
}

func testRuntimeTestRendersFailures(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteExecutable("go", "#!/usr/bin/env bash\nprintf '%s\\n' '{\"Action\":\"output\",\"Package\":\"example/fail\",\"Test\":\"TestParent/TestChild\",\"Output\":\"boom child\\n\"}' '{\"Action\":\"fail\",\"Package\":\"example/fail\",\"Test\":\"TestParent/TestChild\"}' '{\"Action\":\"fail\",\"Package\":\"example/fail\",\"Test\":\"TestParent\"}' '{\"Action\":\"fail\",\"Package\":\"example/fail\"}' '{\"Action\":\"output\",\"Package\":\"example/build\",\"Output\":\"undefined: missing\\n\"}' '{\"Action\":\"fail\",\"Package\":\"example/build\"}'\nexit 1\n")

	out := f.BenchEnv(map[string]string{"PATH": f.Root + string(os.PathListSeparator) + os.Getenv("PATH")}, "test")
	want := "packages[2]{package,status}:\n  example/build,fail\n  example/fail,fail\nfailures[2]{package,test,line}:\n  example/build,\"\",\"undefined: missing\"\n  example/fail,TestParent/TestChild,boom child\nskips[0]{package,test,reason}:\n"
	if out.ExitCode != 1 || out.Stderr != "" || out.Stdout != want {
		t.Fatalf("bench test failures = exit %d stdout %q stderr %q, want exit 1 and %q", out.ExitCode, out.Stdout, out.Stderr, want)
	}
}

func testRuntimeTestSelectsFailureDiagnostics(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteExecutable("go", "#!/usr/bin/env bash\nprintf '%s\\n' '{\"Action\":\"fail\",\"Package\":\"example/select\",\"Test\":\"TestNoDiagnostic\"}' '{\"Action\":\"output\",\"Package\":\"example/select\",\"Test\":\"TestParent\",\"Output\":\"parent direct\\n\"}' '{\"Action\":\"output\",\"Package\":\"example/select\",\"Test\":\"TestParent/TestChild\",\"Output\":\"child direct\\n\"}' '{\"Action\":\"fail\",\"Package\":\"example/select\",\"Test\":\"TestParent/TestChild\"}' '{\"Action\":\"fail\",\"Package\":\"example/select\",\"Test\":\"TestParent\"}' '{\"Action\":\"output\",\"Package\":\"example/select\",\"Test\":\"TestStructured\",\"Output\":\"bench-skip kind=environment reason=host absent\\n\"}' '{\"Action\":\"fail\",\"Package\":\"example/select\",\"Test\":\"TestStructured\"}' '{\"Action\":\"fail\",\"Package\":\"example/select\"}'\nexit 1\n")

	out := f.BenchEnv(map[string]string{"PATH": f.Root + string(os.PathListSeparator) + os.Getenv("PATH")}, "test")
	want := "packages[1]{package,status}:\n  example/select,fail\nfailures[4]{package,test,line}:\n  example/select,TestNoDiagnostic,no diagnostic emitted\n  example/select,TestParent,parent direct\n  example/select,TestParent/TestChild,child direct\n  example/select,TestStructured,no diagnostic emitted\nskips[0]{package,test,reason}:\n"
	if out.ExitCode != 1 || out.Stderr != "" || out.Stdout != want {
		t.Fatalf("bench test failure selection = exit %d stdout %q stderr %q, want exit 1 and %q", out.ExitCode, out.Stdout, out.Stderr, want)
	}
}

func testRuntimeTestErrorPostures(t *testing.T) {
	for _, tc := range []struct {
		name, script, want string
	}{
		{name: "malformed", script: "#!/bin/sh\nprintf '{not json}\\n'\n", want: "error: go test output malformed — "},
		{name: "no packages", script: "#!/bin/sh\nexit 0\n", want: "error: go test reported no packages — no package terminal event\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := contract.NewFixture(t)
			f.WriteExecutable("go", tc.script)
			out := f.BenchEnv(map[string]string{"PATH": f.Root + string(os.PathListSeparator) + os.Getenv("PATH")}, "test")
			if out.ExitCode != 1 || out.Stderr != "" || !strings.HasPrefix(out.Stdout, tc.want) {
				t.Fatalf("bench test %s = exit %d stdout %q stderr %q, want structured stdout error", tc.name, out.ExitCode, out.Stdout, out.Stderr)
			}
		})
	}

	f := contract.NewFixture(t)
	out := f.BenchEnv(map[string]string{"PATH": "/usr/bin:/bin"}, "test")
	if out.ExitCode != 1 || out.Stderr != "" || !strings.HasPrefix(out.Stdout, "error: go test failed to start — ") {
		t.Fatalf("bench test without Go = exit %d stdout %q stderr %q, want structured stdout error", out.ExitCode, out.Stdout, out.Stderr)
	}
}
