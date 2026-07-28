package surface

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/packagesurface"
	"github.com/gibbonmi/bench/internal/testrepo"
)

func TestPlatformPackageMetadataDerivesFromRoot(t *testing.T) {
	subject := contract.SubjectRoot(t)
	candidate := filepath.Join(t.TempDir(), "candidate")
	if err := testrepo.CommitWorkingTree(subject, candidate); err != nil {
		t.Fatal(err)
	}
	packagePath := filepath.Join(candidate, "package.json")
	var root map[string]any
	contract.ReadJSONFile(t, packagePath, &root)
	root["author"] = "mutation-proves-derivation"
	data, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(packagePath, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	wrapper, packages := filepath.Join(candidate, "generated-wrapper"), filepath.Join(candidate, "generated-packages")
	for _, platform := range packageReadPlatforms(t, candidate) {
		bin := filepath.Join(packages, platform.OS+"-"+platform.Arch, "bin", "bench")
		if err := os.MkdirAll(filepath.Dir(bin), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(bin, []byte("fixture"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	probe := contract.NewExecFixtureAt(t, candidate).Run("node", filepath.Join(candidate, "scripts", "build-release-evidence.mjs"), candidate, wrapper, packages)
	probe.RequireExit(0)
	for _, platform := range packageReadPlatforms(t, candidate) {
		var got struct{ Repository, Homepage, Bugs, Author string }
		contract.ReadJSONFile(t, filepath.Join(packages, platform.OS+"-"+platform.Arch, "package.json"), &got)
		if got.Repository != root["repository"] || got.Homepage != root["homepage"] || got.Bugs != root["bugs"] || got.Author != "mutation-proves-derivation" {
			t.Fatalf("%s-%s metadata did not derive from root: %+v", platform.OS, platform.Arch, got)
		}
	}
}

func TestPackageContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectFileMissing(t, "scripts/gen-platform-packages.sh")
	tmp := t.TempDir()
	gen := filepath.Join(tmp, "artifacts")
	firstCandidate := packageHostCandidate(t)
	packageRunGenerator(t, firstCandidate, gen).RequireExit(0)
	first := packageArtifactNames(t, gen)
	firstPins := packagePinManifest(t, gen)
	t.Run("platform-package generator failed", func(t *testing.T) {
		wantArtifacts := len(packageReadPlatforms(t, firstCandidate))*2 + 1
		if len(first) != wantArtifacts {
			t.Fatalf("artifact count = %d, want packages plus archives = %d", len(first), wantArtifacts)
		}
	})
	secondCandidate := packageHostCandidate(t)
	packageRunGenerator(t, secondCandidate, gen).RequireExit(0)
	second := packageArtifactNames(t, gen)
	secondPins := packagePinManifest(t, gen)
	t.Run("platform-package generator (2nd run) failed", func(t *testing.T) {
		if len(second) == 0 {
			t.Fatal("second artifact build emitted no tarballs")
		}
	})
	t.Run("platform-package generator is not idempotent", func(t *testing.T) {
		if !reflect.DeepEqual(first, second) {
			t.Fatalf("artifact inventory changed on repack: %v != %v", first, second)
		}
	})
	t.Run("binary pin manifest is not reproducible", func(t *testing.T) {
		if !reflect.DeepEqual(firstPins, secondPins) {
			t.Fatalf("binary pins changed across independent builds\nfirst=%s\nsecond=%s", firstPins, secondPins)
		}
	})
	t.Run("platform-package generator output contract failed", func(t *testing.T) {
		testPackageGeneratorOutputAt(t, secondCandidate, gen)
	})
	contract.RunParallel(t, "npm pack installable-surface contract", testPackageNpmPackInstallableSurface)
}

func testPackageGeneratorOutputAt(t *testing.T, root, gen string) {
	t.Helper()
	matrix := packageReadPlatforms(t, root)
	sourceWrapper := packageReadWrapper(t)
	wrapperEntries := contract.ReadTarball(t, filepath.Join(gen, "redbench-"+sourceWrapper.Version+".tgz"))
	var wrapper packageWrapper
	if err := json.Unmarshal(wrapperEntries["package/package.json"].Data, &wrapper); err != nil {
		t.Fatalf("parse emitted wrapper metadata: %v", err)
	}
	wantOptional := map[string]string{}
	var pins struct {
		SchemaVersion int                                         `json:"schema_version"`
		Pins          []struct{ Name, Version, Integrity string } `json:"pins"`
	}
	pinEntry, ok := wrapperEntries["package/binary-pins.json"]
	if !ok {
		t.Fatal("wrapper package omits binary-pins.json")
	}
	if err := json.Unmarshal(pinEntry.Data, &pins); err != nil {
		t.Fatalf("parse binary pins: %v", err)
	}
	if pins.SchemaVersion != 1 {
		t.Fatalf("binary pin schema = %d", pins.SchemaVersion)
	}
	wantPins := map[string]string{}

	for _, p := range matrix {
		name := "@redbench/" + p.OS + "-" + p.Arch
		wantOptional[name] = wrapper.Version
		path := filepath.Join(gen, "redbench-"+p.OS+"-"+p.Arch+"-"+wrapper.Version+".tgz")
		tarball, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		sum := sha512.Sum512(tarball)
		wantPins[name+"@"+wrapper.Version] = "sha512-" + base64.StdEncoding.EncodeToString(sum[:])
		var got struct {
			Name       string   `json:"name"`
			Version    string   `json:"version"`
			Bin        any      `json:"bin"`
			OS         []string `json:"os"`
			CPU        []string `json:"cpu"`
			Repository string   `json:"repository"`
			Homepage   string   `json:"homepage"`
			Bugs       string   `json:"bugs"`
			Author     string   `json:"author"`
		}
		entries := contract.ReadTarball(t, path)
		if err := json.Unmarshal(entries["package/package.json"].Data, &got); err != nil {
			t.Fatalf("parse %s package metadata: %v", name, err)
		}
		if got.Name != name {
			t.Fatalf("%s: name is %s", name, got.Name)
		}
		if got.Version != wrapper.Version {
			t.Fatalf("%s: version %s != wrapper %s", name, got.Version, wrapper.Version)
		}
		if !reflect.DeepEqual(got.OS, []string{p.OS}) {
			t.Fatalf("%s: os %v", name, got.OS)
		}
		if !reflect.DeepEqual(got.CPU, []string{p.Arch}) {
			t.Fatalf("%s: cpu %v", name, got.CPU)
		}
		if got.Bin == nil {
			t.Fatalf("%s: missing bin field", name)
		}
		if got.Repository != wrapper.Repository || got.Homepage != wrapper.Homepage || got.Bugs != wrapper.Bugs || got.Author != wrapper.Author {
			t.Fatalf("%s project metadata differs from wrapper", name)
		}
	}
	gotPins := map[string]string{}
	for _, pin := range pins.Pins {
		gotPins[pin.Name+"@"+pin.Version] = pin.Integrity
	}
	if !reflect.DeepEqual(gotPins, wantPins) {
		t.Fatalf("wrapper binary pins = %v, want %v", gotPins, wantPins)
	}
	if !reflect.DeepEqual(wrapper.OptionalDependencies, wantOptional) {
		t.Fatalf("wrapper optionalDependencies %v != matrix %v", wrapper.OptionalDependencies, wantOptional)
	}
}

func packagePinManifest(t testing.TB, dir string) []byte {
	t.Helper()
	wrapper := packageReadWrapper(t)
	entries := contract.ReadTarball(t, filepath.Join(dir, "redbench-"+wrapper.Version+".tgz"))
	entry, ok := entries["package/binary-pins.json"]
	if !ok {
		t.Fatal("wrapper omits binary pin manifest")
	}
	return append([]byte(nil), entry.Data...)
}

func packageArtifactNames(t testing.TB, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names
}

func testPackageNpmPackInstallableSurface(t *testing.T) {
	f := contract.NewExecFixtureAt(t, contract.SubjectRoot(t))

	// --ignore-scripts: inspect files[] membership only, not the prepare build (which the
	// git-install probe exercises); running prepare here would rebuild dist/bench and
	// defeat the built/unbuilt determinism this shape check holds.
	out := f.Run("npm", "pack", "--dry-run", "--json", "--ignore-scripts")

	out.RequireExit(0)
	var packs []struct {
		Files []struct {
			Path string `json:"path"`
		} `json:"files"`
	}
	if err := json.Unmarshal([]byte(out.Stdout), &packs); err != nil {
		t.Fatalf("npm pack --dry-run JSON unreadable: %v\nstdout:\n%s\nstderr:\n%s", err, out.Stdout, out.Stderr)
	}
	files := map[string]bool{}
	if len(packs) > 0 {
		for _, file := range packs[0].Files {
			files[file.Path] = true
		}
	}
	buildAssets, err := packagesurface.RequiredBuildPackAssets(contract.SubjectRoot(t))
	if err != nil {
		t.Fatalf("read required npm build inputs: %v", err)
	}
	requiredAssets := append([]string{}, packagesurface.RequiredPackAssets...)
	requiredAssets = append(requiredAssets, buildAssets...)
	for _, required := range requiredAssets {
		if !files[required] {
			t.Fatalf("npm package missing %s", required)
		}
	}
	for _, forbidden := range packagesurface.ForbiddenPackAssets {
		if files[forbidden] {
			t.Fatalf("npm package includes local-only file %s", forbidden)
		}
	}
}

type packageWrapper struct {
	Version              string            `json:"version"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
	Repository           string            `json:"repository"`
	Homepage             string            `json:"homepage"`
	Bugs                 string            `json:"bugs"`
	Author               string            `json:"author"`
}

func packageRunGenerator(t testing.TB, candidate, out string) contract.Probe {
	t.Helper()
	return contract.NewExecFixtureAt(t, contract.SubjectRoot(t)).Run("bash", filepath.Join(candidate, "scripts", "gen-platform-packages.sh"), out)
}

// packageHostCandidate stages an independent generator candidate whose release plan
// carries the host target alone. Dev proves the generator's logic and idempotency, which
// one target demonstrates; full matrix breadth is a release-tier claim. Each call builds
// its own candidate, so two generator runs stay the independent builds the pin
// reproducibility assertion compares.
func packageHostCandidate(t testing.TB) string {
	t.Helper()
	return contract.NarrowReleasePlan(t, contract.SubjectRoot(t), func(matrix contract.ReleasePlanTargets) []contract.ReleaseTarget {
		if len(matrix.Host) == 0 {
			capability.Environment(t, fmt.Sprintf("package contract tests require release plan target for host %s/%s", matrix.GOOS, matrix.GOArch))
		}
		return matrix.Host
	})
}

func packageReadPlatforms(t testing.TB, root string) []contract.ReleaseTarget {
	t.Helper()
	var plan struct {
		Targets []contract.ReleaseTarget `json:"targets"`
	}
	contract.ReadJSONFile(t, filepath.Join(root, "scripts", "release-plan.json"), &plan)
	return plan.Targets
}

func packageReadWrapper(t testing.TB) packageWrapper {
	t.Helper()
	var wrapper packageWrapper
	contract.ReadJSONFile(t, filepath.Join(contract.SubjectRoot(t), "package.json"), &wrapper)
	if wrapper.OptionalDependencies == nil {
		wrapper.OptionalDependencies = map[string]string{}
	}
	return wrapper
}
