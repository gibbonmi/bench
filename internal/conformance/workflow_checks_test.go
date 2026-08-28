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

// secondSourceMarker names the isolated checkout the reduction retired. Every second
// generation reached it through this directory, so one marker covers the clone, the
// build, the uploads, and the downloads that would restore it.
const secondSourceMarker = "bench-second-source"

// The artifacts job publishes one upload naming both paths, so the upload action
// resolves dist as the archive root. nativeArtifactDownload is the matching consumer
// shape: a download into dist, never into dist/artifacts.
const (
	nativeArtifactName     = "name: ${{ inputs.artifact-prefix || 'bench-runtime' }}-artifacts\n"
	nativeArtifactUpload   = nativeArtifactName + "          path: |\n            dist/artifacts\n            dist/reproducibility.json\n"
	nativeArtifactDownload = nativeArtifactName + "          path: dist\n"
)

// releaseArtifactDownload is the same archive root seen from the tag workflow, whose
// jobs name the release prefix literally.
const releaseArtifactDownload = "name: release-artifacts\n          path: dist\n"

// staleReleaseArtifactDownload is the pre-reduction path, one level below the archive
// root. The mutations below use it to prove each tag job graded separately.
const staleReleaseArtifactDownload = "name: release-artifacts\n          path: dist/artifacts\n"

// replaceOccurrence rewrites the nth occurrence, counting from zero, so a mutation
// reaches exactly one of several jobs that carry byte-identical steps.
func replaceOccurrence(text, old, replacement string, n int) string {
	for offset, index := 0, 0; ; index++ {
		found := strings.Index(text[offset:], old)
		if found < 0 {
			return text
		}
		found += offset
		if index == n {
			return text[:found] + replacement + text[found+len(old):]
		}
		offset = found + len(old)
	}
}

// retiredReproducibilityRecord is the cross-checkout comparison output the reduction
// retired. The name is assembled from two pieces because the sweep below reads this
// source file, and a whole literal would report itself.
const retiredReproducibilityRecord = "workflow-" + "reproducibility.json"

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
	// Both tag jobs consume the shared upload, whose archive root is dist. A stale
	// dist/artifacts download nests the tarballs and fails every tag release.
	for _, job := range []string{"authorize", "publish"} {
		if !strings.Contains(workflowJob(text, job), releaseArtifactDownload) {
			diags = append(diags, "release workflow "+job+" job does not download the artifacts at their archive root")
		}
	}
	return diags
}

// checkRetiredReproducibilityRecord sweeps the source trees for the cross-checkout
// record that retired with the second generation. Nothing writes the file now, so any
// surviving reference sends a reader to evidence that never arrives.
func checkRetiredReproducibilityRecord(root string) []string {
	var diags []string
	for _, top := range []string{".github", "scripts", "internal", "tests", "docs"} {
		_ = filepath.WalkDir(filepath.Join(root, top), func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil || entry.IsDir() || !strings.Contains(readIfExists(path), retiredReproducibilityRecord) {
				return nil
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				rel = path
			}
			diags = append(diags, "retired cross-checkout reproducibility record is still named by "+filepath.ToSlash(rel))
			return nil
		})
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
	if strings.Count(artifacts, "uses: actions/upload-artifact@") != 1 {
		diags = append(diags, "native artifact construction does not publish exactly one upload")
	}
	// The one upload carries the artifacts and the reproducibility record, so the
	// action resolves dist as the archive root.
	if !strings.Contains(artifacts, nativeArtifactUpload) {
		diags = append(diags, "native verification does not hand reproducibility records to evidence finalization")
	}
	proof, smoke := workflowJob(text, "native-proof"), workflowJob(text, "smoke")
	// The workflow builds one generation, so a restored second source in a producing
	// job both returns the cost this reduction removed and reads an upload nothing
	// publishes. Every consuming job downloads at the archive root, because a download
	// into dist/artifacts nests the tarballs one level too deep. One table over the
	// four jobs states which fact each job carries, and names the job that failed.
	for _, job := range []struct {
		name, body       string
		builds, consumes bool
	}{
		{name: "artifacts", body: artifacts, builds: true},
		{name: "native-proof", body: proof, builds: true, consumes: true},
		{name: "evidence", body: evidence, builds: true, consumes: true},
		{name: "smoke", body: smoke, consumes: true},
	} {
		if job.builds && strings.Contains(job.body, secondSourceMarker) {
			diags = append(diags, "native "+job.name+" job rebuilds a second source generation")
		}
		if job.consumes && !strings.Contains(job.body, nativeArtifactDownload) {
			diags = append(diags, "native "+job.name+" job does not download the artifact upload at its archive root")
		}
	}
	if !strings.Contains(smoke, "needs: [preflight, artifacts, evidence]") || !strings.Contains(smoke, "preflight-evidence") || !strings.Contains(smoke, "scripts/smoke-artifacts.sh") {
		diags = append(diags, "native verification does not run smoke from finalized evidence")
	}
	// Shipped targets and proven targets are separate facts. The proof job reads the
	// proven view, so an unproven target starts no runner; smoke keeps the shipped
	// view, so every shipped binary still executes on its own operating system.
	if !strings.Contains(proof, "matrix: ${{ fromJSON(needs.preflight.outputs.proven) }}") {
		diags = append(diags, "native proof matrix does not read the proven targets")
	}
	if !strings.Contains(smoke, "matrix: ${{ fromJSON(needs.preflight.outputs.matrix) }}") {
		diags = append(diags, "native smoke matrix does not read the shipped targets")
	}
	if proof := readIfExists(filepath.Join(root, "scripts", "native-proof.sh")); proof != "" {
		diags = append(diags, nativeProofBindingDiags(proof)...)
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
		// The authorize job downloads first and the publish job last, so the two cases
		// mutate opposite ends of the same repeated block.
		{
			name:   "authorize job downloads below the archive root",
			broken: replaceOccurrence(workflow, releaseArtifactDownload, staleReleaseArtifactDownload, 0),
			want:   "release workflow authorize job does not download the artifacts at their archive root",
			cheat:  releaseArtifactDownload,
		},
		{
			name:   "publish job downloads below the archive root",
			broken: replaceOccurrence(workflow, releaseArtifactDownload, staleReleaseArtifactDownload, 1),
			want:   "release workflow publish job does not download the artifacts at their archive root",
			cheat:  releaseArtifactDownload,
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
	workflowBytes, err := os.ReadFile(filepath.Join(kit, ".github", "workflows", "native-runtime.yml"))
	if err != nil {
		t.Fatal(err)
	}
	workflow := string(workflowBytes)
	path := filepath.Join(root, ".github", "workflows", "native-runtime.yml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	diagnose := func(t *testing.T, broken string) string {
		t.Helper()
		if broken == workflow {
			t.Fatal("native workflow mutation changed nothing")
		}
		if err := os.WriteFile(path, []byte(broken), 0o644); err != nil {
			t.Fatal(err)
		}
		return strings.Join(checkNativeRuntimeWorkflow(root), "\n")
	}
	const (
		artifactsBuild      = "      - name: Build the artifact generation\n        run: bash scripts/build-artifacts.sh . dist/artifacts\n"
		proofRun            = "      - run: bash scripts/native-proof.sh dist/artifacts "
		evidenceFinalize    = "      - name: Finalize the release evidence\n        run: |\n"
		secondSourceClone   = "          git clone --quiet --no-hardlinks . \"$RUNNER_TEMP/bench-second-source\"\n"
		artifactsUploadTail = "            dist/reproducibility.json\n          retention-days: ${{ github.retention_days }}\n"
		secondUpload        = "      - uses: actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02 # v4\n        with:\n          name: ${{ inputs.artifact-prefix || 'bench-runtime' }}-artifacts-second\n          path: dist/artifacts\n"
	)
	staleDownload := nativeArtifactName + "          path: dist/artifacts\n"
	for _, bite := range []struct{ name, broken, want string }{
		{
			name:   "smoke loses its evidence edge",
			broken: strings.Replace(workflow, "needs: [preflight, artifacts, evidence]", "needs: [preflight, artifacts]", 1),
			want:   "native verification does not run smoke from finalized evidence",
		},
		{
			name:   "artifact construction bypasses preflight authorization",
			broken: strings.Replace(workflow, "needs: preflight", "needs: []", 1),
			want:   "native artifact construction or upload bypasses preflight authorization",
		},
		{
			name:   "gate step loses the capability flag",
			broken: strings.Replace(workflow, "\n        env:\n          BENCH_REQUIRE_CAPABILITIES: '1'\n", "\n", 1),
			want:   "native verification does not require capabilities on the gate step",
		},
		{
			name:   "proof matrix reads the shipped targets",
			broken: strings.Replace(workflow, "outputs.proven)", "outputs.matrix)", 1),
			want:   "native proof matrix does not read the proven targets",
		},
		{
			name:   "smoke matrix reads the proven targets",
			broken: strings.Replace(workflow, "outputs.matrix)", "outputs.proven)", 1),
			want:   "native smoke matrix does not read the shipped targets",
		},
		{
			name:   "the reproducibility record leaves the upload",
			broken: strings.Replace(workflow, "            dist/reproducibility.json\n", "", 1),
			want:   "native verification does not hand reproducibility records to evidence finalization",
		},
		// The three consuming jobs carry byte-identical download steps in source order:
		// native-proof, then evidence, then smoke. Each case mutates one of them, and
		// the job-named diagnostic proves the mutation reached the job it claims.
		{
			name:   "native-proof downloads below the archive root",
			broken: replaceOccurrence(workflow, nativeArtifactDownload, staleDownload, 0),
			want:   "native native-proof job does not download the artifact upload at its archive root",
		},
		{
			name:   "evidence downloads below the archive root",
			broken: replaceOccurrence(workflow, nativeArtifactDownload, staleDownload, 1),
			want:   "native evidence job does not download the artifact upload at its archive root",
		},
		{
			name:   "smoke downloads below the archive root",
			broken: replaceOccurrence(workflow, nativeArtifactDownload, staleDownload, 2),
			want:   "native smoke job does not download the artifact upload at its archive root",
		},
		{
			name:   "a second source returns to artifact construction",
			broken: strings.Replace(workflow, artifactsBuild, "      - name: Build the artifact generation\n        run: |\n"+secondSourceClone+"          bash scripts/build-artifacts.sh . dist/artifacts\n", 1),
			want:   "native artifacts job rebuilds a second source generation",
		},
		{
			name:   "a second source returns to the native proof",
			broken: strings.Replace(workflow, proofRun, "      - run: git clone --quiet --no-hardlinks . \"$RUNNER_TEMP/bench-second-source\"\n"+proofRun, 1),
			want:   "native native-proof job rebuilds a second source generation",
		},
		{
			name:   "a second source returns to evidence finalization",
			broken: strings.Replace(workflow, evidenceFinalize, evidenceFinalize+secondSourceClone, 1),
			want:   "native evidence job rebuilds a second source generation",
		},
		{
			name:   "artifact construction publishes a second upload",
			broken: strings.Replace(workflow, artifactsUploadTail, artifactsUploadTail+secondUpload, 1),
			want:   "native artifact construction does not publish exactly one upload",
		},
	} {
		t.Run(bite.name, func(t *testing.T) {
			if diagnostics := diagnose(t, bite.broken); !strings.Contains(diagnostics, bite.want) {
				t.Fatalf("mutation did not bite with %q:\n%s", bite.want, diagnostics)
			}
		})
	}
	// The retired record leaves the whole tree, not one job, so its sweep takes a
	// restored reference in a source file rather than a workflow mutation.
	if err := os.WriteFile(path, workflowBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	stale := filepath.Join(root, "docs", "release-runbook.md")
	if err := os.MkdirAll(filepath.Dir(stale), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stale, []byte("The evidence job writes dist/"+retiredReproducibilityRecord+".\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if diagnostics := strings.Join(checkRetiredReproducibilityRecord(root), "\n"); !strings.Contains(diagnostics, "retired cross-checkout reproducibility record is still named by docs/release-runbook.md") {
		t.Fatalf("restored reproducibility record reference did not bite:\n%s", diagnostics)
	}
}
