package adopt

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTransactionalLinkAdoptsUnownedAdapterThroughSymlinkParent(t *testing.T) {
	for _, tc := range []struct {
		name       string
		content    string
		target     string
		wantCode   int
		wantAbsent string
	}{
		{name: "converged", content: "same\n", wantCode: 0},
		{name: "divergent", content: "different\n", target: "../adapter-mirror", wantCode: 1, wantAbsent: "accepted.txt"},
		{name: "foreign-identical", content: "same\n", target: "../adapter-mirror", wantCode: 1, wantAbsent: "accepted.txt"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			runAdoptGit(t, root, "init", "-q")
			if err := os.MkdirAll(filepath.Join(root, ".agents", "commands"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, ".agents", "commands", "bench-implement-spec.md"), []byte("same\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			if tc.target != "" {
				if err := os.MkdirAll(filepath.Join(root, "adapter-mirror"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, "adapter-mirror", "bench-implement-spec.md"), []byte(tc.content), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			if err := os.MkdirAll(filepath.Join(root, ".claude"), 0o755); err != nil {
				t.Fatal(err)
			}
			target := "../.agents/commands"
			if tc.target != "" {
				target = tc.target
			}
			if err := os.Symlink(target, filepath.Join(root, ".claude", "commands")); err != nil {
				t.Fatal(err)
			}
			kitAsset := filepath.Join(t.TempDir(), "bench-implement-spec.md")
			if err := os.WriteFile(kitAsset, []byte("same\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			plan := []planEntry{{src: kitAsset, rel: ".claude/commands/bench-implement-spec.md", kind: "adapter"}}
			if tc.wantAbsent != "" {
				plan = append([]planEntry{{rel: tc.wantAbsent, kind: "inline", content: "accepted\n"}}, plan...)
			}
			var stdout, stderr bytes.Buffer
			code, _ := transactionalLink(root, t.TempDir(), "copy", "test", plan, &stdout, &stderr)
			if code != tc.wantCode {
				t.Fatalf("transactionalLink exit = %d, want %d; stderr=%q", code, tc.wantCode, stderr.String())
			}
			if tc.wantCode == 0 {
				if got, err := os.ReadFile(filepath.Join(root, ".agents", "commands", "bench-implement-spec.md")); err != nil || string(got) != tc.content {
					t.Fatalf("resolved adapter content = %q, %v", got, err)
				}
			} else if _, err := os.Lstat(filepath.Join(root, tc.wantAbsent)); !os.IsNotExist(err) {
				t.Fatalf("accepted write was promoted: %v", err)
			}
			if tc.wantCode != 0 && !strings.Contains(stderr.String(), "has a symlink parent directory") {
				t.Fatalf("stderr = %q, want symlink-parent conflict", stderr.String())
			}
		})
	}
}

// consumerRepo is a fresh git repository, made the working directory, that links
// against this kit checkout. It is the one setup for a real link round trip.
func consumerRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runAdoptGit(t, root, "init", "-q")
	t.Setenv("BENCH_KIT", filepath.Clean(filepath.Join(mustGetwd(t), "..", "..")))
	t.Chdir(root)
	return root
}

// linkConsumer runs a real Link in root, requires wantCode, and returns the
// manifest rows the link recorded.
func linkConsumer(t *testing.T, root string, wantCode int) []manifestRow {
	t.Helper()
	var stdout, stderr bytes.Buffer
	if code := Link(nil, &stdout, &stderr, "1.0.0"); code != wantCode {
		t.Fatalf("Link = %d, want %d\nstdout:\n%s\nstderr:\n%s", code, wantCode, stdout.String(), stderr.String())
	}
	m, err := ReadManifest(filepath.Join(root, ".bench", "link-manifest.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	return m.Rows()
}

func runUnlink(t *testing.T, args []string, wantCode int) string {
	t.Helper()
	var stdout, stderr bytes.Buffer
	if code := Unlink(args, &stdout, &stderr); code != wantCode {
		t.Fatalf("Unlink %v = %d, want %d\nstdout:\n%s\nstderr:\n%s", args, code, wantCode, stdout.String(), stderr.String())
	}
	return stdout.String() + stderr.String()
}

func requireAbsent(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		t.Fatalf("%s after unlink: %v, want absent", path, err)
	}
}

// TestUnlinkRemovesOnlyCleanManifestAssets pins the clean round trip: every manifest
// row, the manifest, and the managed hook leave, and project content stays.
func TestUnlinkRemovesOnlyCleanManifestAssets(t *testing.T) {
	root := consumerRepo(t)
	writeFixtureFile(t, filepath.Join(root, "notes.md"), "project notes\n", 0o644)
	writeFixtureFile(t, filepath.Join(root, "AGENTS.md"), "# Team rules\n\nKeep this rule.\n", 0o644)
	rows := linkConsumer(t, root, 0)
	if len(rows) == 0 {
		t.Fatal("Link recorded no manifest rows")
	}

	runUnlink(t, nil, 0)
	for _, row := range rows {
		requireAbsent(t, filepath.Join(root, row.rel))
	}
	requireAbsent(t, filepath.Join(root, ".bench", "link-manifest.tsv"))
	requireAbsent(t, filepath.Join(root, ".git", "hooks", "pre-push"))
	if got := readFile(t, filepath.Join(root, "notes.md")); got != "project notes\n" {
		t.Fatalf("notes.md = %q, want the project bytes", got)
	}
	if got := readFile(t, filepath.Join(root, "AGENTS.md")); !strings.Contains(got, "Keep this rule.") || strings.Contains(got, benchStartMarker) {
		t.Fatalf("AGENTS.md = %q, want the project prose without the managed block", got)
	}
}

// TestUnlinkKeepsModifiedAssetsAndProjectCollisions pins the partial posture: a
// modified managed asset and a project-owned collision keep their bytes, the clean
// rows still leave, and the exit and the residuals table report the partial result.
func TestUnlinkKeepsModifiedAssetsAndProjectCollisions(t *testing.T) {
	root := consumerRepo(t)
	const collision = ".agents/commands/bench-implement-spec.md"
	if err := os.MkdirAll(filepath.Join(root, ".agents", "commands"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFixtureFile(t, filepath.Join(root, collision), "project command\n", 0o644)
	rows := linkConsumer(t, root, 3)
	writeFixtureFile(t, filepath.Join(root, ".bench", "BENCH.md"), "edited by the project\n", 0o644)

	runUnlink(t, []string{"--dry-run"}, 3)
	for _, row := range rows {
		if _, err := os.Lstat(filepath.Join(root, row.rel)); err != nil {
			t.Fatalf("dry run removed %s: %v", row.rel, err)
		}
	}

	out := runUnlink(t, nil, 3)
	if !strings.Contains(out, "residuals[1]{path,reason}:\n  .bench/BENCH.md,modified\n") {
		t.Fatalf("unlink output = %q, want the residuals table with the modified row", out)
	}
	if got := readFile(t, filepath.Join(root, ".bench", "BENCH.md")); got != "edited by the project\n" {
		t.Fatalf(".bench/BENCH.md = %q, want the edited bytes", got)
	}
	if got := readFile(t, filepath.Join(root, collision)); got != "project command\n" {
		t.Fatalf("%s = %q, want the project bytes", collision, got)
	}
	if _, err := os.Lstat(filepath.Join(root, ".bench", "link-manifest.tsv")); err != nil {
		t.Fatalf("link manifest after a partial unlink: %v, want it kept", err)
	}
	for _, row := range rows {
		if row.rel != ".bench/BENCH.md" && row.rel != collision {
			requireAbsent(t, filepath.Join(root, row.rel))
		}
	}
}

// TestUnlinkRefusesAManifestRowOutsideTheRepo plants a hand-edited row that escapes
// the repository and carries the exact fingerprint of the file it names.
func TestUnlinkRefusesAManifestRowOutsideTheRepo(t *testing.T) {
	root := consumerRepo(t)
	rows := linkConsumer(t, root, 0)
	var bench manifestRow
	for _, row := range rows {
		if row.rel == ".bench/BENCH.md" {
			bench = row
		}
	}
	if bench.hash == "" {
		t.Fatal("Link recorded no .bench/BENCH.md row")
	}
	outside := t.TempDir()
	victim := filepath.Join(outside, "BENCH.md")
	writeFixtureFile(t, victim, readFile(t, filepath.Join(root, ".bench", "BENCH.md")), 0o644)
	escape := "../" + filepath.Base(outside) + "/BENCH.md"
	manifest := filepath.Join(root, ".bench", "link-manifest.tsv")
	writeFixtureFile(t, manifest, readFile(t, manifest)+escape+"\t"+bench.hash+"\n", 0o644)

	out := runUnlink(t, nil, 3)
	if !strings.Contains(out, "refused: "+escape+"\n") {
		t.Fatalf("unlink output = %q, want the refused escape row", out)
	}
	if _, err := os.Lstat(victim); err != nil {
		t.Fatalf("file outside the repository after unlink: %v, want it kept", err)
	}
}

// TestUnlinkWithoutManifestExitsOne pins the loud refusal on a repository that unlink
// cannot account for, in place of a silent no-op.
func TestUnlinkWithoutManifestExitsOne(t *testing.T) {
	consumerRepo(t)
	if out := runUnlink(t, nil, 1); !strings.Contains(out, "no link manifest") {
		t.Fatalf("unlink output = %q, want it to name the missing manifest", out)
	}
}
