package conformance

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

type workflowTriggerShape struct {
	pullRequest, pushBranches, mainBranch bool
}

// releaseSubmitInvocation is the whole publish command, binary path included. The
// verifier CI runs must be the one it compiled from the tag's checkout. The
// adapter and provenance selections make the run reach real npm.
const releaseSubmitInvocation = `dist/bench release submit --version "${GITHUB_REF_NAME#v}" --profile public --path first --adapter npm --provenance --registry https://registry.npmjs.org`

// The workflow seam is the job body: these checks parse job ownership before
// asserting that evidence follows all native proof rows.
func checkReleaseWorkflow(root string) []string {
	if !exists(filepath.Join(root, "scripts", "release-plan.json")) {
		return nil
	}
	wf := filepath.Join(root, ".github", "workflows", "release.yml")
	if !exists(wf) {
		return []string{"release workflow missing (.github/workflows/release.yml)"}
	}
	text := readIfExists(wf)
	var diags []string
	if !regexp.MustCompile(`(?m)^\s*tags:`).MatchString(text) {
		diags = append(diags, "release workflow does not trigger on tags")
	}
	if !strings.Contains(readIfExists(filepath.Join(root, "scripts", "build-artifacts.sh")), "scripts/release-plan.mjs") {
		diags = append(diags, "artifact builder does not derive targets from the canonical release plan")
	}
	for message, anchor := range map[string]string{"release workflow does not run full publish preflight": "scripts/release-preflight.sh --mode publish", "release workflow does not publish through bench release submit": releaseSubmitInvocation, "release workflow does not require capabilities on the gate step": requireCapabilitiesGateStep("--mode publish --profile public")} {
		if !strings.Contains(text, anchor) {
			diags = append(diags, message)
		}
	}
	// Publication authority is the state machine alone. A raw publish bypasses the
	// resumable, digest-verified path. Promotion is the reviewer's attended act.
	for message, anchor := range map[string]string{"release workflow publishes with raw npm publish": "npm publish", "release workflow promotes from CI": "release promote"} {
		if strings.Contains(text, anchor) {
			diags = append(diags, message)
		}
	}
	// This check scopes to the publish job, not the file. The authorize job uploads
	// publish-preflight-evidence itself, so a file-wide match would pass on its bytes.
	publish := workflowJob(text, "publish")
	for message, anchor := range map[string]string{"release workflow publish job does not download the publish preflight evidence": "name: publish-preflight-evidence\n", "release workflow publish job does not upload the publication record": "name: publication-record\n"} {
		if !strings.Contains(publish, anchor) {
			diags = append(diags, message)
		}
	}
	if job := workflowJob(text, "verify"); !strings.Contains(job, "uses: ./.github/workflows/native-runtime.yml") || !strings.Contains(job, "artifact-prefix: release") {
		diags = append(diags, "release workflow does not compose shared native verification")
	}
	return diags
}

func checkNativeRuntimeWorkflow(root string) []string {
	workflow := filepath.Join(root, ".github", "workflows", "native-runtime.yml")
	if !exists(workflow) {
		return []string{"native verification workflow missing (.github/workflows/native-runtime.yml)"}
	}
	text, diags := readIfExists(workflow), []string{}
	triggers := nativeWorkflowTriggers(text)
	for label, present := range map[string]bool{"pull requests": triggers.pullRequest, "default branch pushes": triggers.pushBranches, "main branch": triggers.mainBranch} {
		if !present {
			diags = append(diags, "native verification workflow does not include "+label)
		}
	}
	for label, anchor := range map[string]string{"canonical matrix": "scripts/release-plan.mjs", "derived matrix": "fromJSON(needs.preflight.outputs.matrix)", "full preflight": "scripts/release-preflight.sh --mode verify", "native proof builder": "scripts/native-proof.sh", "native proof aggregate": "scripts/aggregate-native-proofs.sh", "native proof evidence": "dist/native-proofs"} {
		if !strings.Contains(text, anchor) {
			diags = append(diags, "native verification workflow does not include "+label)
		}
	}
	preflight, artifacts := workflowJob(text, "preflight"), workflowJob(text, "artifacts")
	if !strings.Contains(preflight, "scripts/release-preflight.sh --mode verify --phase gate") || !strings.Contains(artifacts, "needs: preflight") || !strings.Contains(artifacts, "scripts/build-artifacts.sh") || !strings.Contains(artifacts, "uses: actions/upload-artifact@") {
		diags = append(diags, "native artifact construction or upload bypasses preflight authorization")
	}
	if !strings.Contains(preflight, requireCapabilitiesGateStep("--mode verify --phase gate")) {
		diags = append(diags, "native verification does not require capabilities on the gate step")
	}
	evidence := workflowJob(text, "evidence")
	if !strings.Contains(evidence, "needs: [artifacts, native-proof]") || !strings.Contains(evidence, "scripts/release-preflight.sh --mode verify") {
		diags = append(diags, "native verification does not finalize evidence after every native proof")
	}
	if !strings.Contains(artifacts, `git clone --quiet --no-hardlinks . "$second_source"`) || !strings.Contains(artifacts, `"$second_source/scripts/build-artifacts.sh" "$second_source" "$second_source/dist/artifacts"`) || !strings.Contains(evidence, `bash scripts/compare-artifacts.sh dist/artifacts "$second_source/dist/artifacts" dist/workflow-reproducibility.json . "$second_source" dist/preflight "$second_source/dist/preflight"`) {
		diags = append(diags, "native verification does not compare independently finalized evidence")
	}
	for _, handoff := range []struct{ name, upload, download string }{
		{name: "first", upload: "path: dist/reproducibility.json", download: "path: dist"},
		{name: "second", upload: "path: ${{ runner.temp }}/bench-second-source/dist/reproducibility.json", download: "path: ${{ runner.temp }}/bench-second-source/dist"},
	} {
		upload := "name: ${{ inputs.artifact-prefix || 'bench-runtime' }}-reproducibility-" + handoff.name + "\n          " + handoff.upload + "\n"
		download := "name: ${{ inputs.artifact-prefix || 'bench-runtime' }}-reproducibility-" + handoff.name + "\n          " + handoff.download + "\n"
		if !strings.Contains(artifacts, upload) || !strings.Contains(evidence, download) {
			diags = append(diags, "native verification does not hand reproducibility records to evidence finalization")
			break
		}
	}
	if job := workflowJob(text, "smoke"); !strings.Contains(job, "needs: [preflight, artifacts, evidence]") || !strings.Contains(job, "preflight-evidence") || !strings.Contains(job, "scripts/smoke-artifacts.sh") {
		diags = append(diags, "native verification does not run smoke from finalized evidence")
	}
	if proof := readIfExists(filepath.Join(root, "scripts", "native-proof.sh")); proof != "" && !strings.Contains(proof, "docker run --rm --network none") {
		diags = append(diags, "native proof does not isolate the Linux non-glibc execution")
	}
	return diags
}

// requireCapabilitiesGateStep is the step shape both release paths must carry: the
// strict capability flag bound to the same step that runs the gate. The check matches
// the env key against the run line, not the file, so the flag provably reaches the
// gate. A class silently skipped on a fully capable runner must fail the release, not
// ship.
//
// The env block follows the run line rather than preceding it. Canary fixtures mutate
// these workflows by anchoring on the exact bytes of a preflight run line, six-space
// indent included. A leading env block re-indents that line into a continuation, and
// those anchors stop occurring.
func requireCapabilitiesGateStep(preflightArgs string) string {
	return "run: bash scripts/release-preflight.sh " + preflightArgs + "\n        env:\n          BENCH_REQUIRE_CAPABILITIES: '1'\n"
}

func workflowJob(workflow, name string) string {
	needle := "  " + name + ":\n"
	start := strings.Index(workflow, needle)
	if start < 0 {
		return ""
	}
	rest, end := workflow[start+len(needle):], len(workflow[start+len(needle):])
	for _, candidate := range []string{"\n  preflight:\n", "\n  artifacts:\n", "\n  native-proof:\n", "\n  evidence:\n", "\n  smoke:\n", "\n  authorize:\n", "\n  publish:\n"} {
		if offset := strings.Index(rest, candidate); offset >= 0 && offset < end {
			end = offset
		}
	}
	return rest[:end]
}

func nativeWorkflowTriggers(text string) workflowTriggerShape {
	var shape workflowTriggerShape
	inOn, inPush, inBranches := false, false, false
	for _, raw := range strings.Split(text, "\n") {
		line, trimmed := strings.TrimRight(raw, " \t\r"), strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		if indent == 0 {
			inOn, inPush, inBranches = line == "on:", false, false
			continue
		}
		if !inOn {
			continue
		}
		if indent == 2 {
			inPush, inBranches = trimmed == "push:", false
			if trimmed == "pull_request:" {
				shape.pullRequest = true
			}
			continue
		}
		if inPush && indent == 4 && trimmed == "branches:" {
			shape.pushBranches, inBranches = true, true
			continue
		}
		if inBranches && indent == 6 && trimmed == "- main" {
			shape.mainBranch = true
		}
	}
	return shape
}

// TestReleaseWorkflowPublicationBites proves every publication-authority diagnostic
// red-capable. Each case mutates the live release workflow into a temp root whose
// only other content is the release plan the check gates on.
func TestReleaseWorkflowPublicationBites(t *testing.T) {
	kit, err := findKitRoot()
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	plan, err := os.ReadFile(filepath.Join(kit, "scripts", "release-plan.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "scripts", "release-plan.json"), plan, 0o644); err != nil {
		t.Fatal(err)
	}
	workflowBytes, err := os.ReadFile(filepath.Join(kit, ".github", "workflows", "release.yml"))
	if err != nil {
		t.Fatal(err)
	}
	workflow := string(workflowBytes)
	path := filepath.Join(root, ".github", "workflows", "release.yml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	diagnose := func(t *testing.T, broken string) string {
		t.Helper()
		if broken == workflow {
			t.Fatal("release workflow mutation changed nothing")
		}
		if err := os.WriteFile(path, []byte(broken), 0o644); err != nil {
			t.Fatal(err)
		}
		return strings.Join(checkReleaseWorkflow(root), "\n")
	}
	// The publish job's own bytes must carry the evidence handoff, so both
	// scoping cases leave the artifact name present elsewhere in the file.
	for _, bite := range []struct {
		name, broken, want, cheat string
	}{
		{
			name:   "raw npm publish returns",
			broken: workflow + "      - run: npm publish dist/artifacts/redbench-0.0.0.tgz --access public\n",
			want:   "release workflow publishes with raw npm publish",
		},
		{
			name:   "submit invocation drops the npm adapter",
			broken: strings.Replace(workflow, "--adapter npm ", "", 1),
			want:   "release workflow does not publish through bench release submit",
		},
		{
			name:   "publish job stops downloading preflight evidence",
			broken: strings.Replace(workflow, "          name: publish-preflight-evidence\n          path: dist/preflight\n", "", 1),
			want:   "release workflow publish job does not download the publish preflight evidence",
			cheat:  "name: publish-preflight-evidence",
		},
		{
			name:   "publish job stops uploading the publication record",
			broken: strings.Replace(strings.Replace(workflow, "          name: publication-record\n", "", 1), "          name: publish-preflight-evidence\n          path: dist/preflight/\n", "          name: publication-record\n          path: dist/preflight/\n", 1),
			want:   "release workflow publish job does not upload the publication record",
			cheat:  "name: publication-record",
		},
		{
			name:   "CI promotes",
			broken: workflow + "      - run: dist/bench release promote --version \"${GITHUB_REF_NAME#v}\"\n",
			want:   "release workflow promotes from CI",
		},
	} {
		t.Run(bite.name, func(t *testing.T) {
			if bite.cheat != "" && !strings.Contains(bite.broken, bite.cheat) {
				t.Fatalf("mutation removed %q file-wide, so it does not exercise the unscoped-contains cheat", bite.cheat)
			}
			if diagnostics := diagnose(t, bite.broken); !strings.Contains(diagnostics, bite.want) {
				t.Fatalf("mutation did not bite with %q:\n%s", bite.want, diagnostics)
			}
		})
	}
}

func TestNativeWorkflowEvidenceEdgeBites(t *testing.T) {
	kit, err := findKitRoot()
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	for _, relative := range []string{"scripts/release-plan.json", "scripts/native-proof.sh"} {
		data, err := os.ReadFile(filepath.Join(kit, filepath.FromSlash(relative)))
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(root, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	workflow, err := os.ReadFile(filepath.Join(kit, ".github", "workflows", "native-runtime.yml"))
	if err != nil {
		t.Fatal(err)
	}
	broken := strings.Replace(string(workflow), "needs: [preflight, artifacts, evidence]", "needs: [preflight, artifacts]", 1)
	if broken == string(workflow) {
		t.Fatal("native workflow mutation did not remove the evidence edge")
	}
	path := filepath.Join(root, ".github", "workflows", "native-runtime.yml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}
	if diagnostics := strings.Join(checkNativeRuntimeWorkflow(root), "\n"); !strings.Contains(diagnostics, "native verification does not run smoke from finalized evidence") {
		t.Fatalf("removed evidence edge did not bite:\n%s", diagnostics)
	}
	broken = strings.Replace(string(workflow), "needs: preflight", "needs: []", 1)
	if broken == string(workflow) {
		t.Fatal("native workflow mutation did not bypass preflight authorization")
	}
	if err := os.WriteFile(path, []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}
	if diagnostics := strings.Join(checkNativeRuntimeWorkflow(root), "\n"); !strings.Contains(diagnostics, "native artifact construction or upload bypasses preflight authorization") {
		t.Fatalf("bypassed preflight authorization did not bite:\n%s", diagnostics)
	}
	broken = strings.Replace(string(workflow), "\n        env:\n          BENCH_REQUIRE_CAPABILITIES: '1'\n", "\n", 1)
	if broken == string(workflow) {
		t.Fatal("native workflow mutation did not strip the capability flag from the gate step")
	}
	if err := os.WriteFile(path, []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}
	if diagnostics := strings.Join(checkNativeRuntimeWorkflow(root), "\n"); !strings.Contains(diagnostics, "native verification does not require capabilities on the gate step") {
		t.Fatalf("gate step detached from the capability flag did not bite:\n%s", diagnostics)
	}
	broken = strings.Replace(string(workflow), ` "$second_source/dist/preflight"`, "", 1)
	if broken == string(workflow) {
		t.Fatal("native workflow mutation did not remove the second finalized-evidence operand")
	}
	if err := os.WriteFile(path, []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}
	if diagnostics := strings.Join(checkNativeRuntimeWorkflow(root), "\n"); !strings.Contains(diagnostics, "native verification does not compare independently finalized evidence") {
		t.Fatalf("removed finalized-evidence comparison did not bite:\n%s", diagnostics)
	}
	broken = strings.Replace(string(workflow), "path: dist/reproducibility.json", "path: dist/wrong-reproducibility.json", 1)
	if broken == string(workflow) {
		t.Fatal("native workflow mutation did not break reproducibility upload path")
	}
	if err := os.WriteFile(path, []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}
	if diagnostics := strings.Join(checkNativeRuntimeWorkflow(root), "\n"); !strings.Contains(diagnostics, "native verification does not hand reproducibility records to evidence finalization") {
		t.Fatalf("broken reproducibility upload path did not bite:\n%s", diagnostics)
	}
	broken = strings.Replace(string(workflow), "name: ${{ inputs.artifact-prefix || 'bench-runtime' }}-reproducibility-first\n          path: dist\n      - uses: actions/download-artifact@", "name: ${{ inputs.artifact-prefix || 'bench-runtime' }}-reproducibility-first\n          path: wrong-dist\n      - uses: actions/download-artifact@", 1)
	if broken == string(workflow) {
		t.Fatal("native workflow mutation did not break reproducibility download path")
	}
	if err := os.WriteFile(path, []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}
	if diagnostics := strings.Join(checkNativeRuntimeWorkflow(root), "\n"); !strings.Contains(diagnostics, "native verification does not hand reproducibility records to evidence finalization") {
		t.Fatalf("broken reproducibility download path did not bite:\n%s", diagnostics)
	}
}
