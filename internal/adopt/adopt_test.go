package adopt

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
)

func TestRewriteAgentsBlockEdges(t *testing.T) {
	block := BenchAgentsBlock()
	cases := []struct {
		name    string
		in      string
		wantHas []string
		wantErr string
	}{
		{
			name:    "append to project content without trailing newline",
			in:      "PROJECT",
			wantHas: []string{"PROJECT\n", "<!-- bench:start -->", "<!-- bench:end -->"},
		},
		{
			name: "replace only the unfenced managed block",
			in: strings.Join([]string{
				"before",
				"<!-- bench:start -->",
				"old managed",
				"<!-- bench:end -->",
				"after",
				"",
			}, "\n"),
			wantHas: []string{"before\n", block, "\nafter\n"},
		},
		{
			name: "preserve fenced marker examples",
			in: strings.Join([]string{
				"# Project",
				"```",
				"<!-- bench:start -->",
				"example",
				"<!-- bench:end -->",
				"```",
				"KEEP",
				"",
			}, "\n"),
			wantHas: []string{"example", "KEEP", "<!-- bench:start -->"},
		},
		{
			name: "reject reversed markers",
			in: strings.Join([]string{
				"PROJECT BEFORE",
				"<!-- bench:end -->",
				"PROJECT MIDDLE",
				"<!-- bench:start -->",
				"PROJECT AFTER",
				"",
			}, "\n"),
			wantErr: "malformed",
		},
		{
			name: "reject unclosed fence around markers",
			in: strings.Join([]string{
				"Broken docs:",
				"```",
				"<!-- bench:start -->",
				"<!-- bench:end -->",
				"KEEP",
				"",
			}, "\n"),
			wantErr: "fence",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := RewriteAgentsBlock(c.in)
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("err = %v, want substring %q", err, c.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("RewriteAgentsBlock err = %v", err)
			}
			for _, want := range c.wantHas {
				if !strings.Contains(got, want) {
					t.Fatalf("output missing %q:\n%s", want, got)
				}
			}
			if strings.Count(got, "## Bench") != 1 {
				t.Fatalf("managed block not exactly once:\n%s", got)
			}
		})
	}
}

func TestManifestParseEdges(t *testing.T) {
	dir := t.TempDir()
	missing, err := ReadManifest(filepath.Join(dir, "missing.tsv"))
	if err != nil {
		t.Fatalf("missing manifest err = %v", err)
	}
	if missing.Hash(".bench/BENCH.md") != "" || missing.KitVersion != "" {
		t.Fatalf("missing manifest = %+v, want empty", missing)
	}

	path := filepath.Join(dir, "link-manifest.tsv")
	data := strings.Join([]string{
		"#kit\t0.2.0",
		".bench/BENCH.md\told",
		"#comment\tignored",
		".bench/BENCH.md\tnew",
		".agents/commands/bench-implement-spec.md\tcmdhash",
		"",
	}, "\n")
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := ReadManifest(path)
	if err != nil {
		t.Fatalf("ReadManifest err = %v", err)
	}
	if got.KitVersion != "0.2.0" {
		t.Fatalf("KitVersion = %q, want 0.2.0", got.KitVersion)
	}
	if got.Hash(".bench/BENCH.md") != "new" {
		t.Fatalf("duplicate rel did not use last value: %+v", got)
	}
	if got.Hash("#kit") != "" {
		t.Fatalf("#kit parsed as a file row: %+v", got)
	}
}

func TestAdapterTarget(t *testing.T) {
	cases := map[string]string{
		".claude/commands/bench-implement-spec.md":      "../../.agents/commands/bench-implement-spec.md",
		".claude/skills/bench-craft-seams/SKILL.md":     "../../../.agents/skills/bench-craft-seams/SKILL.md",
		".claude/skills/bench-craft-seams/refs/deep.md": "../../../../.agents/skills/bench-craft-seams/refs/deep.md",
	}
	for rel, want := range cases {
		got, ok := AdapterTarget(rel)
		if !ok {
			t.Fatalf("AdapterTarget(%q) not ok", rel)
		}
		if got != want {
			t.Fatalf("AdapterTarget(%q) = %q, want %q", rel, got, want)
		}
	}
	if _, ok := AdapterTarget(".bench/BENCH.md"); ok {
		t.Fatalf("non-adapter rel returned ok")
	}
}

func TestDoctorDirSelectionAndShimRoundTrip(t *testing.T) {
	home := t.TempDir()
	nvm := filepath.Join(home, ".nvm")
	manager := filepath.Join(nvm, "versions", "node", "v22", "bin")
	blocked := filepath.Join(home, "blocked")
	plain := filepath.Join(home, "plain bin")
	if err := os.MkdirAll(manager, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(blocked, 0o555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(blocked, 0o755)
	if err := os.MkdirAll(plain, 0o755); err != nil {
		t.Fatal(err)
	}
	env := DoctorEnv{
		Home:   home,
		Path:   strings.Join([]string{manager, blocked, plain}, string(os.PathListSeparator)),
		NVMDir: nvm,
	}
	chosen, create := SelectDoctorDir(env)
	if create {
		t.Fatalf("SelectDoctorDir create = true, want false")
	}
	if chosen != plain {
		t.Fatalf("SelectDoctorDir = %q, want %q", chosen, plain)
	}

	target := filepath.Join(home, "kit path [x]", "bin", "bench.sh")
	content := ShimContent(target)
	if !strings.Contains(content, "# bench-target: "+target) {
		t.Fatalf("shim missing literal target comment:\n%s", content)
	}
	if !strings.Contains(content, "target=") || !strings.Contains(content, `exec "$target" "$@"`) {
		t.Fatalf("shim missing executable target assignment:\n%s", content)
	}
	if got := ShimTarget(content); got != target {
		t.Fatalf("ShimTarget = %q, want %q", got, target)
	}

	env.Path = manager
	fallback, create := SelectDoctorDir(env)
	if fallback != filepath.Join(home, ".local", "bin") || !create {
		t.Fatalf("fallback = %q create=%v, want ~/.local/bin create", fallback, create)
	}
}

func TestScaffoldGateUsesCanarySubcommand(t *testing.T) {
	gate := scaffoldGate()
	mustContain := []string{
		"BENCH_SENTINEL",
		seedCanaryPath,
		"[ -e DO-NOT-SHIP ] && err \"example check: DO-NOT-SHIP marker file present\"",
		"bench=\"$(dirname \"$0\")/bin/bench.sh\"; [ -x \"$bench\" ] || bench=bench",
		"\"$bench\" canary \"$root\" || err \"canary sweep failed\"",
		"BENCH_CANARY_INNER",
	}
	for _, want := range mustContain {
		if !strings.Contains(gate, want) {
			t.Fatalf("scaffold gate missing %q:\n%s", want, gate)
		}
	}
	for _, forbidden := range []string{". \"$gate_dir/lib/canary-run.sh\"", "canary runner missing"} {
		if strings.Contains(gate, forbidden) {
			t.Fatalf("scaffold gate still contains retired sourcing API %q:\n%s", forbidden, gate)
		}
	}
}

func TestInitScaffoldsTwoLevelSeedCanary(t *testing.T) {
	root := t.TempDir()
	cmd := exec.Command("git", "-C", root, "init", "-q")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_KIT", filepath.Clean(filepath.Join(wd, "..", "..")))
	t.Chdir(root)

	var stdout, stderr bytes.Buffer
	if code := Init(nil, &stdout, &stderr); code != 0 {
		t.Fatalf("Init exit = %d, stderr:\n%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "scaffolded "+seedCanaryPath) {
		t.Fatalf("Init stdout missing seed canary path %q:\n%s", seedCanaryPath, stdout.String())
	}

	seedDir := filepath.Join(root, filepath.FromSlash(seedCanaryPath))
	if got := readFile(t, filepath.Join(seedDir, "EXPECT")); got != "example check: DO-NOT-SHIP marker file present\n" {
		t.Fatalf("seed EXPECT = %q", got)
	}
	if got := readFile(t, filepath.Join(seedDir, "files", "DO-NOT-SHIP")); !strings.Contains(got, "seed example check") {
		t.Fatalf("seed marker file text = %q", got)
	}
	if _, err := os.Stat(filepath.Join(root, "tests", "canary", "example", "EXPECT")); !os.IsNotExist(err) {
		t.Fatalf("legacy flat seed EXPECT exists or stat failed unexpectedly: %v", err)
	}

	var fixtureCalls []string
	err = canary.Sweep(root, func(call canary.RunCall) canary.RunResult {
		if call.FixtureDir == "" {
			return canary.RunResult{ExitCode: 0, Output: "baseline without seed expectation"}
		}
		fixtureCalls = append(fixtureCalls, call.FixtureDir)
		return canary.RunResult{ExitCode: 1, Output: "example check: DO-NOT-SHIP marker file present"}
	})
	if err != nil {
		t.Fatalf("canary sweep over initialized seed failed: %v", err)
	}
	if len(fixtureCalls) != 1 || fixtureCalls[0] != seedDir {
		t.Fatalf("fixture calls = %#v, want only %q", fixtureCalls, seedDir)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", path, err)
	}
	return string(data)
}
