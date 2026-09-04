package adopt

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/freshness"
)

// TestWriteBrokerManifestBindsTheRunningBrokerIdentity pins the published binding:
// the manifest lands beside the wrapper and binds the running executable's path, the
// given version, this host's platform, and the executable's content digest. This
// installer is the one writer of the platform value, so the assertion grades the
// spelling the field must carry.
func TestWriteBrokerManifestBindsTheRunningBrokerIdentity(t *testing.T) {
	dir := t.TempDir()
	wrapper := filepath.Join(dir, "bin", "bench.sh")
	if err := os.MkdirAll(filepath.Dir(wrapper), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_WRAPPER", wrapper)

	path, broker, err := WriteBrokerManifest("1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	if path != filepath.Join(dir, "bin", brokermanifest.Name) {
		t.Fatalf("manifest path = %s, want beside the wrapper", path)
	}
	fields, err := brokermanifest.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if fields["path"] != broker || !filepath.IsAbs(broker) {
		t.Fatalf("manifest path field = %q, want the absolute broker %q", fields["path"], broker)
	}
	running, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if resolved, err := filepath.EvalSymlinks(running); err == nil {
		running = resolved
	}
	if broker != running {
		t.Fatalf("bound broker = %q, want the running executable %q", broker, running)
	}
	if fields["version"] != "1.2.3" {
		t.Fatalf("manifest version = %q", fields["version"])
	}
	if !regexp.MustCompile(`^(linux|darwin)-(x64|arm64)$`).MatchString(fields["platform"]) {
		t.Fatalf("manifest platform = %q", fields["platform"])
	}
	digest, err := fingerprintPath(broker)
	if err != nil {
		t.Fatal(err)
	}
	if fields["sha256"] != digest {
		t.Fatalf("manifest digest = %q, want %q", fields["sha256"], digest)
	}
}

// TestPublishedBrokerManifestCarriesTheStampedVersion grades the field the land route
// compares against the installed package. A manifest that carries "dev" stops the landing
// at exit 127 with nothing to repair.
func TestPublishedBrokerManifestCarriesTheStampedVersion(t *testing.T) {
	root := t.TempDir()
	writeSealFixture(t, root)
	staged := filepath.Join(root, "staged-bench")
	if err := os.WriteFile(staged, []byte("Bench executable"), 0o755); err != nil {
		t.Fatal(err)
	}
	executable := filepath.Join(root, "dist", "bench")
	if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
		t.Fatal(err)
	}

	// The manifest directory is deliberately not the executable's own, because that is the
	// shape an ordinary build has: the executable lands in dist/ and the manifest lands
	// beside the wrapper the land route resolves.
	manifestDir := filepath.Join(root, "bin")
	if err := os.MkdirAll(manifestDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := freshness.Publish(root, staged, executable, manifestDir, "4.5.6"); err != nil {
		t.Fatal(err)
	}
	fields, err := brokermanifest.Read(filepath.Join(manifestDir, brokermanifest.Name))
	if err != nil {
		t.Fatal(err)
	}
	if fields["version"] != "4.5.6" {
		t.Fatalf("published manifest version = %q, want the stamped version", fields["version"])
	}
	digest, err := fingerprintPath(executable)
	if err != nil {
		t.Fatal(err)
	}
	if fields["sha256"] != digest {
		t.Fatalf("published manifest digest = %q, want the published executable digest %q", fields["sha256"], digest)
	}
}

// writeSealFixture materializes the smallest tree freshness.Digest can resolve: one
// buildable cmd/bench, its module file, and the auxiliary inputs manifest that makes the
// repository declare build inputs at all. The doctor seal row is scoped by that
// declaration, so a fixture without it would leave the row quiet.
func writeSealFixture(t *testing.T, root string) {
	t.Helper()
	for name, body := range map[string]string{
		"go.mod":                  "module example.com/doctorfixture\n\ngo 1.25\n",
		"cmd/bench/main.go":       "package main\n\nfunc main() {}\n",
		"scripts/go-build.sh":     "#!/usr/bin/env bash\n",
		"scripts/go-build.inputs": "build_script=scripts/go-build.sh\n",
	} {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// TestDoctorSealRowNamesTheSourceDigestMismatch grades the row an operator reads before a
// landing. Every other row reports healthy beside a binary built from sources that have
// since changed, so this row must name the mismatch itself and carry the one rebuild
// sentence, not a generic "stale binary" verdict.
func TestDoctorSealRowNamesTheSourceDigestMismatch(t *testing.T) {
	root := t.TempDir()
	runAdoptGit(t, root, "init", "-q")
	if err := os.Mkdir(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeSealFixture(t, root)
	staged := filepath.Join(root, "staged-bench")
	if err := os.WriteFile(staged, []byte("Bench executable"), 0o755); err != nil {
		t.Fatal(err)
	}
	executable := filepath.Join(root, "dist", "bench")
	if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := freshness.Publish(root, staged, executable, filepath.Dir(executable), "1.2.3"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "cmd", "bench", "main.go"), []byte("package main\n\nfunc main() { _ = 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)

	var stdout bytes.Buffer
	if !reportDoctorRows(&stdout) {
		t.Fatalf("doctor rows reported no red beside a stale seal:\n%s", stdout.String())
	}
	got := stdout.String()
	if !containsRow(got, "red", "seal source digest does not match current build inputs") {
		t.Fatalf("doctor rows = %q, want the source-digest mismatch verdict", got)
	}
	resolved := root
	if real, err := filepath.EvalSymlinks(root); err == nil {
		resolved = real
	}
	if want := freshness.RebuildAction(resolved); !strings.Contains(got, want) {
		t.Fatalf("doctor rows = %q, want the rebuild action %q", got, want)
	}
}

// TestDoctorSealRowStaysQuietWithoutAPublishedBinary keeps the row scoped. A repository
// that publishes no dist/bench has no seal to grade, and a red row there would refuse
// every consumer checkout for an artifact it never owned.
func TestDoctorSealRowStaysQuietWithoutAPublishedBinary(t *testing.T) {
	root := t.TempDir()
	runAdoptGit(t, root, "init", "-q")
	if err := os.Mkdir(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeSealFixture(t, root)
	t.Chdir(root)

	var stdout bytes.Buffer
	_ = reportDoctorRows(&stdout)
	if got := stdout.String(); strings.Contains(got, "red: dist/bench") {
		t.Fatalf("doctor rows = %q, want no red seal row without a published binary", got)
	}
}

// TestDoctorBrokerRowNamesEveryLandRouteRefusal grades the row against all five reasons
// the land route refuses on. The route is shell that runs before any binary is trusted,
// so the row is a second derivation by necessity; a row that graded four of the five
// predicates would predict only part of the exit 127 an operator then meets.
func TestDoctorBrokerRowNamesEveryLandRouteRefusal(t *testing.T) {
	for _, tc := range []struct {
		name    string
		arrange func(t *testing.T, bindir, broker string)
		want    string
		plain   bool
	}{
		{
			// The absent manifest is the one reason the row names without going red: the
			// resolved wrapper belongs to the environment, not to the repository under doctor.
			name:    "absent manifest",
			arrange: func(*testing.T, string, string) {},
			want:    "no promotion-broker manifest at ",
			plain:   true,
		},
		{
			name: "incomplete manifest",
			arrange: func(t *testing.T, bindir, broker string) {
				t.Helper()
				writeBrokerFixture(t, bindir, broker, "1.2.3")
				if err := os.WriteFile(filepath.Join(bindir, brokermanifest.Name), []byte("path\t"+broker+"\nversion\t1.2.3\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: " is incomplete",
		},
		{
			name: "version skew",
			arrange: func(t *testing.T, bindir, broker string) {
				t.Helper()
				writeBrokerFixture(t, bindir, broker, "9.9.9")
			},
			want: " does not match installed package ",
		},
		{
			name: "broker not executable",
			arrange: func(t *testing.T, bindir, broker string) {
				t.Helper()
				writeBrokerFixture(t, bindir, broker, "1.2.3")
				if err := os.Chmod(broker, 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: " is not a regular executable",
		},
		{
			name: "digest mismatch",
			arrange: func(t *testing.T, bindir, broker string) {
				t.Helper()
				writeBrokerFixture(t, bindir, broker, "1.2.3")
				if err := os.WriteFile(broker, []byte("replaced by a hand rebuild"), 0o755); err != nil {
					t.Fatal(err)
				}
			},
			want: " does not match its manifest digest",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			bindir, broker := arrangeBrokerInstall(t)
			tc.arrange(t, bindir, broker)

			var stdout bytes.Buffer
			_ = reportDoctorRows(&stdout)
			got := stdout.String()
			verdict := "red"
			if tc.plain {
				verdict = "ok"
			}
			if !containsRow(got, verdict, tc.want) {
				t.Fatalf("doctor rows = %q, want a %s broker row naming %q", got, verdict, tc.want)
			}
		})
	}
}

// TestDoctorBrokerRowAcceptsAnAuthenticatedBroker is the row's green side. Without it the
// five refusals could be satisfied by a row that never accepts anything.
func TestDoctorBrokerRowAcceptsAnAuthenticatedBroker(t *testing.T) {
	bindir, broker := arrangeBrokerInstall(t)
	writeBrokerFixture(t, bindir, broker, "1.2.3")

	var stdout bytes.Buffer
	_ = reportDoctorRows(&stdout)
	got := stdout.String()
	if !strings.Contains(got, "ok: promotion broker authenticated at "+broker) {
		t.Fatalf("doctor rows = %q, want the authenticated broker row", got)
	}
}

// arrangeBrokerInstall builds an installed-kit shape around a bench-touched repository:
// an install root carrying package.json, its bin directory with the wrapper the manifest
// publishes beside, and a broker executable. It returns the bin directory and the broker.
func arrangeBrokerInstall(t *testing.T) (string, string) {
	t.Helper()
	repo := t.TempDir()
	runAdoptGit(t, repo, "init", "-q")
	if err := os.Mkdir(filepath.Join(repo, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(repo)

	install := t.TempDir()
	bindir := filepath.Join(install, "bin")
	if err := os.MkdirAll(bindir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(install, "package.json"), []byte("{\"version\":\"1.2.3\"}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bindir, "bench.sh"), []byte("#!/usr/bin/env bash\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	broker := filepath.Join(bindir, "bench")
	if err := os.WriteFile(broker, []byte("promotion broker bytes"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_WRAPPER", filepath.Join(bindir, "bench.sh"))
	return bindir, broker
}

// containsRow reports whether one rendered doctor line carries both the verdict and the
// reason. Grading the two together is what keeps a red reason from passing on some other
// row's red.
func containsRow(rendered, verdict, reason string) bool {
	for _, line := range strings.Split(rendered, "\n") {
		if strings.HasPrefix(line, "  "+verdict+": ") && strings.Contains(line, reason) {
			return true
		}
	}
	return false
}

func writeBrokerFixture(t *testing.T, bindir, broker, version string) {
	t.Helper()
	if _, _, err := brokermanifest.Write(bindir, broker, version); err != nil {
		t.Fatal(err)
	}
}
