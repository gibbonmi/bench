package conformance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// planTarget is one row of a release-plan matrix view. The views print os, arch,
// goos, goarch, and runner as tab-separated fields, so one decoder serves the
// shipped view and the proven view alike.
type planTarget struct {
	os, arch, goos, runner string
}

func (target planTarget) name() string { return target.os + "-" + target.arch }

// planTargets reads one matrix view out of the repository's own release plan. The
// scripts under test read the same plan, so the fixtures a test builds always
// describe the targets the scripts actually iterate.
func planTargets(t *testing.T, kit, command string) []planTarget {
	t.Helper()
	out, err := exec.Command("node", filepath.Join(kit, "scripts", "release-plan.mjs"), kit, command).Output()
	if err != nil {
		t.Fatalf("release-plan.mjs %s: %v", command, err)
	}
	var targets []planTarget
	for _, line := range strings.Split(strings.TrimRight(string(out), "\n"), "\n") {
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 5 {
			t.Fatalf("release-plan.mjs %s: row %q has %d fields, want 5", command, line, len(fields))
		}
		targets = append(targets, planTarget{os: fields[0], arch: fields[1], goos: fields[2], runner: fields[4]})
	}
	if len(targets) == 0 {
		t.Fatalf("release-plan.mjs %s returned no rows", command)
	}
	return targets
}

// unprovenTargetName names a target the proof scripts must refuse. It prefers a
// target the plan ships without a proof. A plan that proves every shipped target
// still needs the refusal graded, so the fallback names a target no plan carries.
func unprovenTargetName(t *testing.T, kit string) planTarget {
	t.Helper()
	proven := map[string]bool{}
	for _, target := range planTargets(t, kit, "proof-targets") {
		proven[target.name()] = true
	}
	for _, target := range planTargets(t, kit, "targets") {
		if !proven[target.name()] {
			return target
		}
	}
	return planTarget{os: "openbsd", arch: "riscv64", goos: "openbsd", runner: "openbsd-latest"}
}

// writeProofFile writes the proof record the aggregator accepts for one target.
// The digests only have to be distinct 64-character hexadecimal strings, and the
// musl field carries the value the target's operating system admits.
func writeProofFile(t *testing.T, dir string, target planTarget) {
	t.Helper()
	digest := func(label string) string {
		sum := sha256.Sum256([]byte(label + target.name()))
		return hex.EncodeToString(sum[:])
	}
	musl := "not_applicable"
	if target.os == "linux" {
		musl = "green"
	}
	body, err := json.Marshal(map[string]any{
		"schema_version":    1,
		"target":            target.name(),
		"runner":            target.runner,
		"status":            "green",
		"rebuilt_sha256":    digest("binary"),
		"binary_sha256":     digest("binary"),
		"package_sha256":    digest("package"),
		"archive_sha256":    digest("archive"),
		"musl_status":       musl,
		"operations_status": "green",
		"strip_status":      "green",
		"tools_status":      "green",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, target.name()+".json"), append(body, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

// provenProofDir builds a proof directory holding one accepted proof per proven
// target, minus the targets named in omit.
func provenProofDir(t *testing.T, kit string, omit ...string) (string, []planTarget) {
	t.Helper()
	skipped := map[string]bool{}
	for _, name := range omit {
		skipped[name] = true
	}
	dir := t.TempDir()
	proven := planTargets(t, kit, "proof-targets")
	for _, target := range proven {
		if !skipped[target.name()] {
			writeProofFile(t, dir, target)
		}
	}
	return dir, proven
}

// unpredicatedProofRoot builds a scratch root whose plan proves a target with no
// platform predicate, and it names that target. native-proof.sh resolves its root
// from its own location, so a copy of the script beside a copy of the plan reader
// grades the refusal against a plan the working tree never carries.
func unpredicatedProofRoot(t *testing.T, kit string) (string, planTarget) {
	t.Helper()
	root := t.TempDir()
	for _, relative := range []string{"package.json", "scripts/release-plan.mjs", "scripts/native-proof.sh", "scripts/lib/search.sh"} {
		data, err := os.ReadFile(filepath.Join(kit, filepath.FromSlash(relative)))
		if err != nil {
			t.Fatal(err)
		}
		copied := filepath.Join(root, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(copied), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(copied, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	planBytes, err := os.ReadFile(filepath.Join(kit, "scripts", "release-plan.json"))
	if err != nil {
		t.Fatal(err)
	}
	var plan struct {
		SchemaVersion  int              `json:"schema_version"`
		Targets        []map[string]any `json:"targets"`
		ArchiveEntries []any            `json:"archive_entries"`
	}
	if err := json.Unmarshal(planBytes, &plan); err != nil {
		t.Fatal(err)
	}
	text := func(value any) string { field, _ := value.(string); return field }
	proven := planTarget{}
	for _, target := range plan.Targets {
		if text(target["goos"]) == "linux" {
			continue
		}
		target["native_proof"] = true
		proven = planTarget{os: text(target["os"]), arch: text(target["arch"]), goos: text(target["goos"]), runner: text(target["runner"])}
		break
	}
	// A plan that ships only Linux still needs the refusal graded, so the fallback
	// adds an operating system no plan carries.
	if proven.os == "" {
		proven = planTarget{os: "openbsd", arch: "riscv64", goos: "openbsd", runner: "openbsd-latest"}
		plan.Targets = append(plan.Targets, map[string]any{"os": proven.os, "arch": proven.arch, "goos": proven.goos, "goarch": "riscv64", "runner": proven.runner, "native_proof": true})
	}
	body, err := json.Marshal(plan)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "scripts", "release-plan.json"), append(body, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, proven
}

// runScript executes one release script and returns its merged output. The proof
// scripts report every refusal on standard error, so the caller grades one stream.
func runScript(t *testing.T, root, script string, args ...string) (string, error) {
	t.Helper()
	command := exec.Command("bash", append([]string{filepath.Join(root, "scripts", script)}, args...)...)
	command.Dir = root
	out, err := command.CombinedOutput()
	return string(out), err
}

func TestAggregateNativeProofsVerifiesTheProvenSet(t *testing.T) {
	kit, err := findKitRoot()
	if err != nil {
		t.Fatal(err)
	}
	dir, proven := provenProofDir(t, kit)

	out, err := runScript(t, kit, "aggregate-native-proofs.sh", dir)

	if err != nil {
		t.Fatalf("aggregate rejected the complete proven set: %v\n%s", err, out)
	}
	RequireSubstring(t, out, "native proof set: "+strconv.Itoa(len(proven))+" canonical targets verified", "aggregate report")
}

func TestAggregateNativeProofsFailsOnMissingProvenProof(t *testing.T) {
	kit, err := findKitRoot()
	if err != nil {
		t.Fatal(err)
	}
	proven := planTargets(t, kit, "proof-targets")
	dropped := proven[0]
	dir, _ := provenProofDir(t, kit, dropped.name())

	out, err := runScript(t, kit, "aggregate-native-proofs.sh", dir)

	if err == nil {
		t.Fatalf("aggregate accepted a proof set missing %s:\n%s", dropped.name(), out)
	}
	RequireSubstring(t, out, "native proof set is missing "+dropped.os+"/"+dropped.arch, "missing proof diagnostic")
}

func TestAggregateNativeProofsFailsOnUnprovenProofFile(t *testing.T) {
	kit, err := findKitRoot()
	if err != nil {
		t.Fatal(err)
	}
	dir, _ := provenProofDir(t, kit)
	unproven := unprovenTargetName(t, kit)
	writeProofFile(t, dir, unproven)

	out, err := runScript(t, kit, "aggregate-native-proofs.sh", dir)

	if err == nil {
		t.Fatalf("aggregate accepted a proof file for unproven target %s:\n%s", unproven.name(), out)
	}
	RequireSubstring(t, out, "native proof set does not contain exactly the canonical proof files", "extra proof diagnostic")
}

func TestNativeProofRefusesUnprovenTarget(t *testing.T) {
	kit, err := findKitRoot()
	if err != nil {
		t.Fatal(err)
	}
	unproven := unprovenTargetName(t, kit)
	work := t.TempDir()

	out, err := runScript(t, kit, "native-proof.sh",
		filepath.Join(work, "artifacts"),
		filepath.Join(work, "proof.json"),
		unproven.os, unproven.arch, unproven.runner)

	if err == nil {
		t.Fatalf("native-proof.sh accepted unproven target %s:\n%s", unproven.name(), out)
	}
	RequireSubstring(t, out,
		"native proof: target is not in the canonical platform matrix: "+unproven.name(), "unproven target refusal")
	if _, statErr := os.Stat(filepath.Join(work, "proof.json")); statErr == nil {
		t.Fatalf("native-proof.sh minted a proof for unproven target %s", unproven.name())
	}
}

func TestNativeProofRefusesTargetWithoutPlatformPredicate(t *testing.T) {
	kit, err := findKitRoot()
	if err != nil {
		t.Fatal(err)
	}
	root, proven := unpredicatedProofRoot(t, kit)
	work := t.TempDir()

	out, err := runScript(t, root, "native-proof.sh",
		filepath.Join(work, "artifacts"),
		filepath.Join(work, "proof.json"),
		proven.os, proven.arch, proven.runner)

	if err == nil {
		t.Fatalf("native-proof.sh proved %s, which has no platform predicate:\n%s", proven.name(), out)
	}
	RequireSubstring(t, out,
		"native proof: no platform predicate exists for operating system "+proven.goos, "unverifiable operating system refusal")
	if _, statErr := os.Stat(filepath.Join(work, "proof.json")); statErr == nil {
		t.Fatalf("native-proof.sh minted a proof for %s, which has no platform predicate", proven.name())
	}
}

// nativeProofBindingDiags grades the proof builder script itself. The script resolves
// its target through the proven view, isolates the Linux non-glibc execution, and
// carries no Darwin symbol assertion: a loadable Mach-O always reports symbols, so
// such an assertion can never go green, and a Darwin proof needs a new one.
func nativeProofBindingDiags(proof string) []string {
	var diags []string
	if !strings.Contains(proof, "docker run --rm --network none") {
		diags = append(diags, "native proof does not isolate the Linux non-glibc execution")
	}
	if !strings.Contains(proof, `"$root/scripts/release-plan.mjs" "$root" proof-target `) {
		diags = append(diags, "native proof does not resolve its target through the proven matrix")
	}
	for _, marker := range []string{"Mach-O", "nm -a", "Darwin binary", "darwin-symbols"} {
		if strings.Contains(proof, marker) {
			diags = append(diags, "native proof still carries the unreachable Darwin branch: "+marker)
		}
	}
	return diags
}

func TestNativeProofBindingDiagsBite(t *testing.T) {
	const bound = "musl_status=not_applicable\n" +
		"docker run --rm --network none -v x:/bench:ro alpine:3.20 /bench version\n" +
		"node \"$root/scripts/release-plan.mjs\" \"$root\" proof-target \"$os_name\" \"$arch_name\"\n"
	if diags := nativeProofBindingDiags(bound); len(diags) != 0 {
		t.Fatalf("bound script reported diagnostics: %v", diags)
	}
	for _, mutation := range []struct{ name, old, replacement, want string }{
		{"docker isolation", "docker run --rm --network none", "docker run", "isolate the Linux non-glibc execution"},
		{"proven matrix", "proof-target", "target", "resolve its target through the proven matrix"},
		{"darwin branch", "musl_status=not_applicable", "nm -a \"$rebuild\"", "unreachable Darwin branch"},
	} {
		t.Run(mutation.name, func(t *testing.T) {
			mutated := strings.Replace(bound, mutation.old, mutation.replacement, 1)
			if mutated == bound {
				t.Fatal("mutation did not apply")
			}
			diags := nativeProofBindingDiags(mutated)
			if len(diags) == 0 || !strings.Contains(strings.Join(diags, "\n"), mutation.want) {
				t.Fatalf("mutation did not bite: %v", diags)
			}
		})
	}
}
