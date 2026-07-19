package artifact

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

type smokeFile struct {
	data []byte
	mode int64
}

func TestOfflineSmokeRunsThePublicJourneyAndAttributesMutations(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/smoke-offline.sh")
	artifacts, evidence := makeOfflineJourneyFixture(t)
	run := func(tb testing.TB, artifactDir, evidenceDir string, env ...map[string]string) contract.Probe {
		fixture := contract.NewExecFixtureAt(tb, root)
		if len(env) != 0 {
			return fixture.RunEnv(env[0], "bash", filepath.Join(root, "scripts", "smoke-offline.sh"), artifactDir, evidenceDir)
		}
		return fixture.Run("bash", filepath.Join(root, "scripts", "smoke-offline.sh"), artifactDir, evidenceDir)
	}
	t.Run("baseline", func(t *testing.T) { run(t, artifacts, evidence).RequireExit(0) })

	t.Run("missing-platform", func(t *testing.T) {
		copy := cloneSmokeDir(t, artifacts)
		if err := os.Remove(filepath.Join(copy, "redbench-linux-x64-0.2.0.tgz")); err != nil {
			t.Fatal(err)
		}
		assertOfflineSmokeRed(t, run(t, copy, evidence), "offline smoke: required artifact is missing or unsafe")
	})
	t.Run("wrong-target", func(t *testing.T) {
		copy, copyEvidence := makeOfflineJourneyFixtureForTarget(t, "darwin/amd64")
		assertOfflineSmokeRed(t, run(t, copy, copyEvidence), "offline smoke: direct version output is wrong")
	})
	t.Run("corrupt-tarball", func(t *testing.T) {
		copy := cloneSmokeDir(t, artifacts)
		if err := os.WriteFile(filepath.Join(copy, "redbench-0.2.0-linux-x64.tar.gz"), []byte("corrupt\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		assertOfflineSmokeRed(t, run(t, copy, evidence), "offline smoke: archive tarball is corrupt")
	})
	t.Run("checksum", func(t *testing.T) {
		copy := cloneSmokeDir(t, evidence)
		if err := os.WriteFile(filepath.Join(copy, "SHA256SUMS"), []byte(fmt.Sprintf("%064d  redbench-0.2.0-linux-x64.tar.gz\n", 0)), 0o644); err != nil {
			t.Fatal(err)
		}
		assertOfflineSmokeRed(t, run(t, artifacts, copy), "offline smoke: supplied release evidence does not bind archive bytes")
	})
	t.Run("repair-control", func(t *testing.T) {
		assertOfflineSmokeRed(t, run(t, artifacts, evidence, map[string]string{"BENCH_OFFLINE_REPAIR_DISABLED": "0"}), "offline smoke: fetch, rebuild, or repair control is disabled")
	})
	t.Run("residue", func(t *testing.T) {
		fake := t.TempDir()
		realNPM, err := filepath.Abs(filepath.Join("/usr", "bin", "npm"))
		if err != nil {
			t.Fatal(err)
		}
		wrapper := fmt.Sprintf("#!/bin/bash\nreal=%q\nprefix=\nprevious=\nfor arg in \"$@\"; do [[ \"$previous\" == --prefix ]] && prefix=\"$arg\"; previous=\"$arg\"; done\n\"$real\" \"$@\"\nrc=$?\nif [[ \"${1:-}\" == uninstall && -n \"$prefix\" ]]; then mkdir -p \"$prefix/node_modules/redbench\"; fi\nexit \"$rc\"\n", realNPM)
		if err := os.WriteFile(filepath.Join(fake, "npm"), []byte(wrapper), 0o755); err != nil {
			t.Fatal(err)
		}
		assertOfflineSmokeRed(t, run(t, artifacts, evidence, map[string]string{"PATH": fake + string(os.PathListSeparator) + os.Getenv("PATH")}), "offline smoke: local npm uninstall left package or Bench state residue")
	})
	t.Run("undeclared-registry-egress", func(t *testing.T) {
		fake := t.TempDir()
		wrapper := "#!/bin/bash\nif [[ \"${1:-}\" == install && \"${npm_config_offline:-}\" != true ]]; then printf 'registry.npmjs.org:443\\n' >> \"$BENCH_OFFLINE_EGRESS_LOG\"; fi\nexec /usr/bin/npm \"$@\"\n"
		if err := os.WriteFile(filepath.Join(fake, "npm"), []byte(wrapper), 0o755); err != nil {
			t.Fatal(err)
		}
		assertOfflineSmokeRed(t, run(t, artifacts, evidence, map[string]string{"PATH": fake + string(os.PathListSeparator) + os.Getenv("PATH")}), "offline smoke: registry install attempted undeclared egress")
	})
}

func TestOfflineNetworkSentinelDeniesUndeclaredEgress(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/offline-network-sentinel.cjs")
	log := filepath.Join(t.TempDir(), "egress.log")
	probe := contract.NewExecFixtureAt(t, root).RunEnv(map[string]string{
		"NODE_OPTIONS":             "--require=" + filepath.Join(root, "scripts", "offline-network-sentinel.cjs"),
		"BENCH_OFFLINE_EGRESS_LOG": log,
	}, "node", "-e", `const socket=require("node:net").connect({host:"registry.npmjs.org",port:443}); socket.on("error", error=>{console.error(error.message);process.exit(7)})`)
	if probe.ExitCode == 0 {
		t.Fatal("offline network sentinel permitted undeclared egress")
	}
	probe.RequireContains(probe.Stderr, "offline smoke: denied undeclared egress to registry.npmjs.org:443")
	data, err := os.ReadFile(log)
	if err != nil || string(data) != "registry.npmjs.org:443\n" {
		t.Fatalf("offline network sentinel log = %q, %v", data, err)
	}
}

func makeOfflineJourneyFixture(t *testing.T) (string, string) {
	return makeOfflineJourneyFixtureForTarget(t, "linux/amd64")
}

func makeOfflineJourneyFixtureForTarget(t *testing.T, runtimeTarget string) (string, string) {
	t.Helper()
	artifacts, evidence := filepath.Join(t.TempDir(), "artifacts"), filepath.Join(t.TempDir(), "evidence")
	if err := os.MkdirAll(artifacts, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(evidence, 0o755); err != nil {
		t.Fatal(err)
	}
	bench := []byte(fmt.Sprintf("#!/bin/bash\ncase \"${1:-}\" in version) printf 'benchkit 0.2.0 (%s)\\n';; commands) [[ \"${2:-}\" == --brief ]] && printf 'commands --brief\\n';; *) printf 'bench — fixture\\n';; esac\n", runtimeTarget))
	wrapper := []byte("#!/bin/bash\nscript=$(readlink -f \"$0\")\nexec \"$(dirname \"$script\")/../../@redbench/linux-x64/bin/bench\" \"$@\"\n")
	nativeFiles := map[string]smokeFile{"bin/bench": {bench, 0o755}, "package.json": {[]byte("{\"name\":\"@redbench/linux-x64\",\"version\":\"0.2.0\",\"bin\":{\"bench\":\"bin/bench\"}}\n"), 0o644}}
	wrapperFiles := map[string]smokeFile{"bin/bench.sh": {wrapper, 0o755}, "package.json": {[]byte("{\"name\":\"redbench\",\"version\":\"0.2.0\",\"bin\":{\"bench\":\"bin/bench.sh\"}}\n"), 0o644}}
	nativeFiles["component-manifest.json"] = smokeManifest(t, nativeFiles)
	wrapperFiles["component-manifest.json"] = smokeManifest(t, wrapperFiles)
	nativeTar := writeSmokeTar(t, filepath.Join(artifacts, "redbench-linux-x64-0.2.0.tgz"), "package/", nativeFiles)
	wrapperTar := writeSmokeTar(t, filepath.Join(artifacts, "redbench-0.2.0.tgz"), "package/", wrapperFiles)
	rootName := "redbench-0.2.0-linux-x64/"
	archiveFiles := map[string]smokeFile{
		"bin/bench": {bench, 0o755}, "packages/redbench-linux-x64-0.2.0.tgz": {nativeTar, 0o644}, "packages/redbench-0.2.0.tgz": {wrapperTar, 0o644},
		"OFFLINE.md": {[]byte("fixture offline instructions\n"), 0o644}, "evidence/components/wrapper-component-manifest.json": {wrapperFiles["component-manifest.json"].data, 0o644}, "evidence/components/platform-component-manifest.json": {nativeFiles["component-manifest.json"].data, 0o644},
	}
	archiveFiles["evidence/component-manifest.json"] = smokeManifest(t, archiveFiles)
	archivePath := filepath.Join(artifacts, "redbench-0.2.0-linux-x64.tar.gz")
	archive := writeSmokeTar(t, archivePath, rootName, archiveFiles)
	digest := fmt.Sprintf("%x", sha256.Sum256(archive))
	index, _ := json.Marshal(map[string]any{"artifacts": []map[string]string{{"name": filepath.Base(archivePath), "sha256": digest}}})
	if err := os.WriteFile(filepath.Join(evidence, "release-index.json"), append(index, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(evidence, "SHA256SUMS"), []byte(digest+"  "+filepath.Base(archivePath)+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return artifacts, evidence
}

func smokeManifest(t *testing.T, files map[string]smokeFile) smokeFile {
	t.Helper()
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	items := make([]map[string]any, 0, len(names))
	for _, name := range names {
		file := files[name]
		items = append(items, map[string]any{"path": name, "mode": fmt.Sprintf("%o", file.mode), "size": len(file.data), "sha256": fmt.Sprintf("%x", sha256.Sum256(file.data))})
	}
	data, err := json.Marshal(map[string]any{"schema_version": 1, "files": items})
	if err != nil {
		t.Fatal(err)
	}
	return smokeFile{append(data, '\n'), 0o644}
}

func writeSmokeTar(t *testing.T, destination, prefix string, files map[string]smokeFile) []byte {
	t.Helper()
	file, err := os.Create(destination)
	if err != nil {
		t.Fatal(err)
	}
	gz := gzip.NewWriter(file)
	writer := tar.NewWriter(gz)
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		item := files[name]
		if err := writer.WriteHeader(&tar.Header{Name: prefix + name, Mode: item.mode, Size: int64(len(item.data)), Typeflag: tar.TypeReg}); err != nil {
			t.Fatal(err)
		}
		if _, err := writer.Write(item.data); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(destination)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func cloneSmokeDir(t *testing.T, source string) string {
	t.Helper()
	destination := t.TempDir()
	entries, err := os.ReadDir(source)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		data, err := os.ReadFile(filepath.Join(source, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(destination, entry.Name()), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return destination
}

func assertOfflineSmokeRed(t *testing.T, probe contract.Probe, message string) {
	t.Helper()
	if probe.ExitCode == 0 {
		t.Fatalf("mutation passed public offline smoke: %s", message)
	}
	probe.RequireContains(probe.Stderr, message)
}
