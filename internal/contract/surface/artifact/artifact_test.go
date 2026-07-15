package artifact

import (
	"bytes"
	"debug/elf"
	"debug/macho"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

type artifactPlatform struct {
	OS     string `json:"os"`
	Arch   string `json:"arch"`
	GOArch string `json:"goarch"`
	Runner string `json:"runner"`
}

type wrapperAsset struct {
	Source string `json:"source"`
	Mode   string `json:"mode"`
	Tree   bool   `json:"tree"`
}

func TestDistributableArtifactContracts(t *testing.T) {
	assertWrapperAssetPolicy(t, contract.SubjectRoot(t))
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	root := contract.SubjectRoot(t)
	buildRoot := committedHostileArtifactSource(t, root)
	out := filepath.Join(t.TempDir(), "artifact output [hostile]")
	probe := contract.NewExecFixtureAt(t, root).Run("bash", filepath.Join(buildRoot, "scripts", "build-artifacts.sh"), buildRoot, out)
	probe.RequireExit(0)

	var matrix []artifactPlatform
	contract.ReadJSONFile(t, filepath.Join(root, "scripts", "platforms.json"), &matrix)
	var wrapper struct {
		Version string `json:"version"`
	}
	contract.ReadJSONFile(t, filepath.Join(root, "package.json"), &wrapper)
	files, err := os.ReadDir(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != len(matrix)+1 {
		t.Fatalf("artifact count = %d, want matrix + wrapper = %d", len(files), len(matrix)+1)
	}

	wrapperTar := filepath.Join(out, fmt.Sprintf("redbench-%s.tgz", wrapper.Version))
	assertWrapperArtifact(t, root, wrapperTar, wrapper.Version, matrix)
	for _, platform := range matrix {
		name := fmt.Sprintf("redbench-%s-%s-%s.tgz", platform.OS, platform.Arch, wrapper.Version)
		assertPlatformArtifact(t, filepath.Join(out, name), wrapper.Version, platform)
	}
	assertInstalledArtifactLifecycle(t, out, wrapper.Version)
	contract.NewExecFixtureAt(t, root).Run("bash", filepath.Join(root, "scripts", "smoke-artifacts.sh"), out).RequireExit(0)

	// A signal in the promotion window restores the old complete set. The same
	// hostile-path invocation can then be rerun and replace it safely.
	assertInterruptedArtifactPromotion(t, buildRoot, out, len(matrix)+1)
	contract.NewExecFixtureAt(t, root).Run("bash", filepath.Join(buildRoot, "scripts", "build-artifacts.sh"), buildRoot, out).RequireExit(0)
	files, err = os.ReadDir(out)
	if err != nil || len(files) != len(matrix)+1 {
		t.Fatalf("second artifact build was not safe: files=%d err=%v", len(files), err)
	}

	// A staging failure must leave an existing promoted set untouched.
	broken := filepath.Join(t.TempDir(), "broken source [*]")
	if err := os.MkdirAll(filepath.Join(broken, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	for path, content := range map[string]string{
		"package.json":                `{"version":"1.0.0"}`,
		"scripts/platforms.json":      "[]\n",
		"scripts/wrapper-assets.json": `[{"source":"special","mode":"0644"}]`,
	} {
		if err := os.WriteFile(filepath.Join(broken, path), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := syscall.Mkfifo(filepath.Join(broken, "special"), 0o644); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(out, "promoted-sentinel")
	if err := os.WriteFile(sentinel, []byte("owned"), 0o644); err != nil {
		t.Fatal(err)
	}
	bad := contract.NewExecFixtureAt(t, root).Run("bash", filepath.Join(root, "scripts", "build-artifacts.sh"), broken, out)
	if bad.ExitCode == 0 {
		t.Fatal("incomplete artifact builder unexpectedly succeeded")
	}
	if got, err := os.ReadFile(sentinel); err != nil || string(got) != "owned" {
		t.Fatalf("failed build changed promoted artifacts: %q, %v", got, err)
	}
}

func assertWrapperAssetPolicy(t *testing.T, root string) {
	t.Helper()
	manifest := filepath.Join(root, "scripts", "wrapper-assets.json")
	if _, err := os.Stat(manifest); os.IsNotExist(err) {
		return
	}
	var assets []wrapperAsset
	contract.ReadJSONFile(t, manifest, &assets)
	present := map[string]bool{}
	for _, asset := range assets {
		present[asset.Source] = true
		if asset.Source == "dist/bench" || strings.HasPrefix(asset.Source, "dist/packages/") {
			t.Fatalf("wrapper artifact contains forbidden entry package/%s", asset.Source)
		}
	}
	// Independent omission oracle: these public caller families must be shipped.
	// The manifest remains the production policy; this named expectation exists so
	// deleting a whole caller family cannot make builder and inventory agree falsely.
	for _, required := range []string{"bin/bench.sh", ".bench/adapters", ".bench/hooks", ".bench/lib"} {
		if !present[required] {
			t.Fatalf("wrapper artifact omits required shipped surface %s", required)
		}
	}
}

func assertInstalledArtifactLifecycle(t *testing.T, artifacts, version string) {
	t.Helper()
	tmp := t.TempDir()
	app := filepath.Join(tmp, "installed prefix [x]")
	repo := filepath.Join(tmp, "linked repo [x]")
	home := filepath.Join(tmp, "home [x]")
	stableBin := filepath.Join(home, "bin")
	for _, dir := range []string{app, repo, stableBin} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(app, "package.json"), []byte("{\"private\":true}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	target := runtime.GOOS + "-" + runtime.GOARCH
	if runtime.GOARCH == "amd64" {
		target = runtime.GOOS + "-x64"
	}
	wrapperTar := filepath.Join(artifacts, "redbench-"+version+".tgz")
	nativeTar := filepath.Join(artifacts, "redbench-"+target+"-"+version+".tgz")
	runLifecycle(t, app, nil, "npm", "install", "--ignore-scripts", "--omit=optional", wrapperTar, nativeTar)
	// A second exact-tarball install is a no-op success, not a source-tree fallback.
	runLifecycle(t, app, nil, "npm", "install", "--ignore-scripts", "--omit=optional", wrapperTar, nativeTar)
	wrapper := filepath.Join(app, "node_modules", "redbench", "bin", "bench.sh")
	env := map[string]string{
		"BENCH_HOME": home,
		"HOME":       home,
		"PATH":       stableBin + string(os.PathListSeparator) + os.Getenv("PATH"),
	}
	versionOut := runLifecycle(t, tmp, env, "bash", wrapper, "version")
	if !strings.Contains(versionOut, "benchkit "+version+" (") {
		t.Fatalf("installed version output = %q", versionOut)
	}

	runLifecycle(t, repo, nil, "git", "init", "-q")
	runLifecycle(t, repo, nil, "git", "config", "user.email", "bench@local")
	runLifecycle(t, repo, nil, "git", "config", "user.name", "bench")
	if err := os.WriteFile(filepath.Join(repo, "AGENTS.md"), []byte("project owner text\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runLifecycle(t, repo, nil, "git", "add", "AGENTS.md")
	runLifecycle(t, repo, nil, "git", "commit", "-qm", "seed")
	runLifecycle(t, repo, env, "bash", wrapper, "link")
	manifest, err := os.ReadFile(filepath.Join(repo, ".bench", "link-manifest.tsv"))
	if err != nil || !bytes.Contains(manifest, []byte("#kit\t"+version+"\n")) {
		t.Fatalf("packed link did not stamp kit version: %q, %v", manifest, err)
	}
	agents, err := os.ReadFile(filepath.Join(repo, "AGENTS.md"))
	if err != nil || !bytes.Contains(agents, []byte("project owner text")) {
		t.Fatalf("packed link did not preserve project owner text: %q, %v", agents, err)
	}
	if info, err := os.Stat(filepath.Join(repo, ".bench", "bin", "bench.sh")); err != nil || info.Mode().Perm()&0o111 == 0 {
		t.Fatalf("packed link did not install executable local launcher: %v, %v", info, err)
	}
	runLifecycle(t, repo, env, "bash", wrapper, "init")
	runLifecycle(t, repo, env, "bash", wrapper, "init")
	runLifecycle(t, repo, env, "bash", wrapper, "doctor", "--fix")
	shim := filepath.Join(stableBin, "bench")
	if info, err := os.Stat(shim); err != nil || info.Mode().Perm()&0o111 == 0 {
		t.Fatalf("packed doctor did not generate executable stable shim: %v, %v", info, err)
	}
	assertPackedSetupForwarding(t, tmp, wrapper, shim, env)
	localOut := runLifecycle(t, repo, env, "bash", filepath.Join(repo, ".bench", "bin", "bench.sh"), "version")
	if localOut != versionOut {
		t.Fatalf("linked operation output %q != installed output %q", localOut, versionOut)
	}
	assertPackedEntrySurfaceIdentity(t, repo, env, versionOut)
	runLifecycle(t, repo, env, "bash", wrapper, "link")
	runLifecycle(t, repo, env, "bash", wrapper, "init")
	status := runLifecycle(t, repo, nil, "git", "status", "--short", "--ignored")
	if !strings.Contains(status, "!! .bench/dist/") {
		t.Fatalf("packed linked repo did not retain ignored runtime state:\n%s", status)
	}
	runPackedFreshClone(t, repo, wrapper, shim, version)
	runLifecycle(t, repo, env, "bash", wrapper, "unlink")
	agents, err = os.ReadFile(filepath.Join(repo, "AGENTS.md"))
	if err != nil || string(agents) != "project owner text\n" {
		t.Fatalf("packed unlink changed project owner text: %q, %v", agents, err)
	}
	if _, err := os.Stat(filepath.Join(repo, ".bench", "link-manifest.tsv")); !os.IsNotExist(err) {
		t.Fatalf("packed unlink left link manifest: %v", err)
	}
}

func assertWrapperArtifact(t *testing.T, root, path, version string, matrix []artifactPlatform) {
	t.Helper()
	entries := contract.ReadTarball(t, path)
	expectedModes := map[string]int64{"package/package.json": 0o644}
	var assets []wrapperAsset
	contract.ReadJSONFile(t, filepath.Join(root, "scripts", "wrapper-assets.json"), &assets)
	for _, asset := range assets {
		mode := int64(0)
		fmt.Sscanf(asset.Mode, "%o", &mode)
		source := filepath.Join(root, filepath.FromSlash(asset.Source))
		if !asset.Tree {
			expectedModes["package/"+asset.Source] = mode
			continue
		}
		err := filepath.WalkDir(source, func(name string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			rel, err := filepath.Rel(root, name)
			if err == nil {
				expectedModes["package/"+filepath.ToSlash(rel)] = mode
			}
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	gotNames := make([]string, 0, len(entries))
	for name := range entries {
		gotNames = append(gotNames, name)
	}
	sort.Strings(gotNames)
	wantNames := make([]string, 0, len(expectedModes))
	for name := range expectedModes {
		wantNames = append(wantNames, name)
	}
	sort.Strings(wantNames)
	if !reflect.DeepEqual(gotNames, wantNames) {
		t.Fatalf("wrapper artifact entries differ from allowlist\ngot: %v\nwant: %v", gotNames, wantNames)
	}
	for name, mode := range expectedModes {
		if entries[name].Mode&0o777 != mode {
			t.Fatalf("wrapper artifact mode %s = %04o, want %04o", name, entries[name].Mode&0o777, mode)
		}
		if len(entries[name].Data) == 0 {
			t.Fatalf("wrapper artifact contains empty allowlisted file %s", name)
		}
	}
	for _, forbidden := range []string{"package/dist/bench", "package/scripts/build-artifacts.sh", "package/projects/benchkit.md"} {
		if _, ok := entries[forbidden]; ok {
			t.Fatalf("wrapper artifact contains forbidden entry %s", forbidden)
		}
	}
	var pkg struct {
		Version              string            `json:"version"`
		OptionalDependencies map[string]string `json:"optionalDependencies"`
	}
	if err := json.Unmarshal(entries["package/package.json"].Data, &pkg); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{}
	for _, p := range matrix {
		want["@redbench/"+p.OS+"-"+p.Arch] = version
	}
	if pkg.Version != version || !reflect.DeepEqual(pkg.OptionalDependencies, want) {
		t.Fatalf("wrapper platform dependencies = %v at %s, want %v at %s", pkg.OptionalDependencies, pkg.Version, want, version)
	}
}

func assertPlatformArtifact(t *testing.T, path, version string, platform artifactPlatform) {
	t.Helper()
	entries := contract.ReadTarball(t, path)
	_, hasBinary := entries["package/bin/bench"]
	_, hasPackage := entries["package/package.json"]
	if len(entries) != 2 || !hasBinary || !hasPackage {
		t.Fatalf("%s platform artifact entries = %v", platform.OS+"-"+platform.Arch, reflect.ValueOf(entries).MapKeys())
	}
	binaryBytes := entries["package/bin/bench"].Data
	if len(binaryBytes) == 0 || entries["package/bin/bench"].Mode&0o111 == 0 {
		t.Fatalf("%s platform binary is empty or non-executable", platform.OS+"-"+platform.Arch)
	}
	if !bytes.Contains(binaryBytes, []byte(version)) {
		t.Fatalf("%s platform binary lacks embedded wrapper version %s", platform.OS+"-"+platform.Arch, version)
	}
	if platform.OS == "linux" {
		f, err := elf.NewFile(bytes.NewReader(binaryBytes))
		if err != nil {
			t.Fatalf("%s binary is not ELF: %v", platform.Arch, err)
		}
		wantMachine := elf.EM_X86_64
		if platform.GOArch == "arm64" {
			wantMachine = elf.EM_AARCH64
		}
		if f.Machine != wantMachine {
			t.Fatalf("linux/%s machine = %v, want %v", platform.Arch, f.Machine, wantMachine)
		}
		libs, err := f.ImportedLibraries()
		if err != nil || len(libs) != 0 {
			t.Fatalf("linux/%s binary is dynamic: %v (%v)", platform.Arch, libs, err)
		}
		if f.Section(".symtab") != nil {
			t.Fatalf("linux/%s binary was not stripped", platform.Arch)
		}
	} else {
		f, err := macho.NewFile(bytes.NewReader(binaryBytes))
		if err != nil {
			t.Fatalf("%s binary is not Mach-O: %v", platform.Arch, err)
		}
		wantCPU := macho.CpuAmd64
		if platform.GOArch == "arm64" {
			wantCPU = macho.CpuArm64
		}
		for _, section := range f.Sections {
			if strings.HasPrefix(section.Name, "__debug") {
				t.Fatalf("darwin/%s binary retained debug section %s", platform.Arch, section.Name)
			}
		}
		if f.Cpu != wantCPU {
			t.Fatalf("darwin/%s format=%v, want %v", platform.Arch, f.Cpu, wantCPU)
		}
	}
}
