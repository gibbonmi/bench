package surface

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/packagesurface"
)

func TestPackageContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectFileMissing(t, "scripts/gen-platform-packages.sh")
	tmp := t.TempDir()
	gen := filepath.Join(tmp, "artifacts")
	packageRunGenerator(t, gen).RequireExit(0)
	first := packageArtifactNames(t, gen)
	t.Run("platform-package generator failed", func(t *testing.T) {
		if len(first) != len(packageReadPlatforms(t))+1 {
			t.Fatalf("artifact count = %d, want matrix + wrapper", len(first))
		}
	})
	packageRunGenerator(t, gen).RequireExit(0)
	second := packageArtifactNames(t, gen)
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
	t.Run("platform-package generator output contract failed", func(t *testing.T) {
		testPackageGeneratorOutputAt(t, gen)
	})
	contract.RunParallel(t, "npm pack installable-surface contract", testPackageNpmPackInstallableSurface)
}

func testPackageGeneratorOutputAt(t *testing.T, gen string) {
	t.Helper()
	matrix := packageReadPlatforms(t)
	sourceWrapper := packageReadWrapper(t)
	wrapperEntries := contract.ReadTarball(t, filepath.Join(gen, "redbench-"+sourceWrapper.Version+".tgz"))
	var wrapper packageWrapper
	if err := json.Unmarshal(wrapperEntries["package/package.json"].Data, &wrapper); err != nil {
		t.Fatalf("parse emitted wrapper metadata: %v", err)
	}
	wantOptional := map[string]string{}

	for _, p := range matrix {
		name := "@redbench/" + p.OS + "-" + p.Arch
		wantOptional[name] = wrapper.Version
		path := filepath.Join(gen, "redbench-"+p.OS+"-"+p.Arch+"-"+wrapper.Version+".tgz")
		var got struct {
			Name    string   `json:"name"`
			Version string   `json:"version"`
			Bin     any      `json:"bin"`
			OS      []string `json:"os"`
			CPU     []string `json:"cpu"`
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
	}
	if !reflect.DeepEqual(wrapper.OptionalDependencies, wantOptional) {
		t.Fatalf("wrapper optionalDependencies %v != matrix %v", wrapper.OptionalDependencies, wantOptional)
	}
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
	for _, required := range packagesurface.RequiredPackAssets {
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

type packagePlatform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type packageWrapper struct {
	Version              string            `json:"version"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
}

func packageRunGenerator(t testing.TB, out string) contract.Probe {
	t.Helper()
	return contract.NewExecFixtureAt(t, contract.SubjectRoot(t)).Run("bash", filepath.Join(contract.SubjectRoot(t), "scripts", "gen-platform-packages.sh"), out)
}

func packageReadPlatforms(t testing.TB) []packagePlatform {
	t.Helper()
	var platforms []packagePlatform
	contract.ReadJSONFile(t, filepath.Join(contract.SubjectRoot(t), "scripts", "platforms.json"), &platforms)
	return platforms
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
