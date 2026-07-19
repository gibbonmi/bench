package publication

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// runRelease runs `bench release <sub> <extraArgs...>` against root, with
// optional env overrides layered on the isolated test environment (used by
// the credential-leak assertion in row 6 to inject a fake NPM_TOKEN).
func runRelease(t *testing.T, root, sub string, extraArgs []string, envOverrides map[string]string) (exitCode int, output string) {
	t.Helper()
	bench := filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")
	env := contract.IsolatedEnv(t, t.TempDir())
	args := append([]string{bench, "release", sub}, extraArgs...)
	cmd := exec.Command("bash", args...)
	cmd.Env = contract.ProcessEnv(env, envOverrides)
	data, err := cmd.CombinedOutput()
	if err == nil {
		return 0, string(data)
	}
	var exitErr *exec.ExitError
	if asExitError(err, &exitErr) {
		return exitErr.ExitCode(), string(data)
	}
	t.Fatalf("release %s did not run: %v\n%s", sub, err, data)
	return -1, string(data)
}

func stagedSubmitArgs(root, version, base string) []string {
	return []string{"--root", root, "--version", version, "--profile", "public", "--registry", base, "--path", "staged"}
}

func httpApprove(t *testing.T, base, stageID string) {
	t.Helper()
	response, err := http.Post(base+"/-/approve/"+stageID, "application/octet-stream", nil)
	if err != nil {
		t.Fatalf("approve %s: %v", stageID, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("approve %s: status %d", stageID, response.StatusCode)
	}
}

// stageIDFor returns the stage_id transition recorded for pkgName, if any.
func stageIDFor(record map[string]any, pkgName string) (string, bool) {
	transitions, _ := record["transitions"].([]any)
	for _, raw := range transitions {
		transition, _ := raw.(map[string]any)
		if transition["package"] == pkgName && transition["action"] == "stage" && transition["result"] == "success" {
			if id, ok := transition["stage_id"].(string); ok && id != "" {
				return id, true
			}
		}
	}
	return "", false
}

// TestStagedSubmitStagesEveryPackageWithoutGoingLive is coverage row 2 (part
// one): the staged path stage-submits every package — platforms and wrapper —
// under auth_mode oidc-stage, and submission alone never makes anything live.
func TestStagedSubmitStagesEveryPackageWithoutGoingLive(t *testing.T) {
	const version = "9.9.10"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	writeApprovedSet(t, root, ordered, nil)
	base, requestFile := startRegistry(t, root, filepath.Join(t.TempDir(), "store"))

	exitCode, output := runRelease(t, root, "submit", stagedSubmitArgs(root, version, base), nil)
	if exitCode != 0 {
		t.Fatalf("staged submit exit=%d:\n%s", exitCode, output)
	}
	if !strings.Contains(output, "approve-platform-packages") {
		t.Fatalf("expected next_action approve-platform-packages, got:\n%s", output)
	}

	lines := requestLines(t, requestFile)
	for _, line := range lines {
		if strings.HasPrefix(line, "PUT ") {
			t.Fatalf("staged submit issued a live PUT; submission must never make a package live:\n%s", strings.Join(lines, "\n"))
		}
	}
	staged := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "STAGE ") {
			staged++
		}
	}
	if staged != len(ordered) {
		t.Fatalf("expected %d STAGE calls (every package), got %d:\n%s", len(ordered), staged, strings.Join(lines, "\n"))
	}

	record := readPublicationRecord(t, root)
	if record["path"] != "public" {
		t.Fatalf("record path = %v, want public", record["path"])
	}
	transitions, _ := record["transitions"].([]any)
	stageCount := 0
	for _, raw := range transitions {
		transition, _ := raw.(map[string]any)
		if transition["action"] == "stage" {
			if transition["auth_mode"] != "oidc-stage" {
				t.Fatalf("stage transition auth_mode = %v, want oidc-stage: %+v", transition["auth_mode"], transition)
			}
			stageCount++
		}
	}
	if stageCount != len(ordered) {
		t.Fatalf("expected %d stage transitions, got %d: %+v", len(ordered), stageCount, transitions)
	}
}

// TestStagedSubmitApprovalHandsOffPlatformsBeforeWrapper is coverage row 2
// (part two): rerunning submit after external platform approvals advances
// next_action from approve-platform-packages to approve-wrapper, and only
// after the wrapper is also approved does it reach promote — the exact
// platform-first-then-wrapper-last approval handoff.
func TestStagedSubmitApprovalHandsOffPlatformsBeforeWrapper(t *testing.T) {
	const version = "9.9.11"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	writeApprovedSet(t, root, ordered, nil)
	base, _ := startRegistry(t, root, filepath.Join(t.TempDir(), "store"))

	exitCode, output := runRelease(t, root, "submit", stagedSubmitArgs(root, version, base), nil)
	if exitCode != 0 {
		t.Fatalf("staged submit (initial stage) exit=%d:\n%s", exitCode, output)
	}
	record := readPublicationRecord(t, root)

	var wrapperName string
	var platformNames []string
	for _, r := range ordered {
		name := npmPackageName(r)
		if r.Kind == "wrapper" {
			wrapperName = name
		} else {
			platformNames = append(platformNames, name)
		}
	}

	// Approve platform packages one at a time; next_action must stay
	// approve-platform-packages until every one of them is approved.
	for i, name := range platformNames {
		stageID, ok := stageIDFor(record, name)
		if !ok {
			t.Fatalf("no stage id recorded for %s: %+v", name, record["transitions"])
		}
		httpApprove(t, base, stageID)
		exitCode, output = runRelease(t, root, "submit", stagedSubmitArgs(root, version, base), nil)
		if exitCode != 0 {
			t.Fatalf("staged submit (approving platform %d) exit=%d:\n%s", i, exitCode, output)
		}
		record = readPublicationRecord(t, root)
		if i < len(platformNames)-1 {
			if !strings.Contains(output, "approve-platform-packages") {
				t.Fatalf("expected next_action approve-platform-packages with platforms still pending, got:\n%s", output)
			}
		}
	}
	if !strings.Contains(output, "approve-wrapper") {
		t.Fatalf("expected next_action approve-wrapper once every platform is approved, got:\n%s", output)
	}

	wrapperStageID, ok := stageIDFor(record, wrapperName)
	if !ok {
		t.Fatalf("no stage id recorded for wrapper %s: %+v", wrapperName, record["transitions"])
	}
	httpApprove(t, base, wrapperStageID)
	exitCode, output = runRelease(t, root, "submit", stagedSubmitArgs(root, version, base), nil)
	if exitCode != 0 {
		t.Fatalf("staged submit (approving wrapper) exit=%d:\n%s", exitCode, output)
	}
	if !strings.Contains(output, "promote") {
		t.Fatalf("expected next_action promote once the complete set is approved, got:\n%s", output)
	}
	record = readPublicationRecord(t, root)
	if record["result"] != "in_progress" {
		t.Fatalf("record result = %v, want in_progress until an explicit promote runs", record["result"])
	}
}

// TestStagedSubmitRejectsWrapperApprovedBeforePlatforms is coverage row 2's
// mutation-direction case: if the wrapper is approved (made live) before every
// platform package, the staged path must detect this as a terminal ordering
// violation rather than silently accepting it.
func TestStagedSubmitRejectsWrapperApprovedBeforePlatforms(t *testing.T) {
	const version = "9.9.12"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	writeApprovedSet(t, root, ordered, nil)
	base, _ := startRegistry(t, root, filepath.Join(t.TempDir(), "store"))

	exitCode, output := runRelease(t, root, "submit", stagedSubmitArgs(root, version, base), nil)
	if exitCode != 0 {
		t.Fatalf("staged submit (initial stage) exit=%d:\n%s", exitCode, output)
	}
	record := readPublicationRecord(t, root)
	var wrapperName string
	for _, r := range ordered {
		if r.Kind == "wrapper" {
			wrapperName = npmPackageName(r)
		}
	}
	stageID, ok := stageIDFor(record, wrapperName)
	if !ok {
		t.Fatalf("no stage id recorded for wrapper %s: %+v", wrapperName, record["transitions"])
	}

	// Simulate the wrapper's 2FA approval happening out of order, before any
	// platform package has been approved.
	httpApprove(t, base, stageID)

	exitCode, output = runRelease(t, root, "submit", stagedSubmitArgs(root, version, base), nil)
	if exitCode == 0 {
		t.Fatalf("staged submit accepted the wrapper being approved before every platform package:\n%s", output)
	}
	record = readPublicationRecord(t, root)
	if record["result"] != "failed" {
		t.Fatalf("record result = %v, want failed after an ordering violation", record["result"])
	}
	transitions, _ := record["transitions"].([]any)
	var flagged bool
	for _, raw := range transitions {
		transition, _ := raw.(map[string]any)
		if transition["package"] == wrapperName && transition["result"] == "order-violation" {
			flagged = true
		}
	}
	if !flagged {
		t.Fatalf("expected an order-violation transition for the wrapper: %+v", transitions)
	}
}

// stagedFlowReadyToPromote drives a staged submit through full approval so
// the record reaches next_action promote.
func stagedFlowReadyToPromote(t *testing.T, root, version, base string) {
	t.Helper()
	exitCode, output := runRelease(t, root, "submit", stagedSubmitArgs(root, version, base), nil)
	if exitCode != 0 {
		t.Fatalf("staged submit (initial stage) exit=%d:\n%s", exitCode, output)
	}
	record := readPublicationRecord(t, root)
	transitions, _ := record["transitions"].([]any)
	for _, raw := range transitions {
		transition, _ := raw.(map[string]any)
		if transition["action"] != "stage" {
			continue
		}
		stageID, _ := transition["stage_id"].(string)
		if stageID != "" {
			httpApprove(t, base, stageID)
		}
	}
	exitCode, output = runRelease(t, root, "submit", stagedSubmitArgs(root, version, base), nil)
	if exitCode != 0 {
		t.Fatalf("staged submit (post-approval) exit=%d:\n%s", exitCode, output)
	}
	if !strings.Contains(output, "promote") {
		t.Fatalf("expected next_action promote after approving every package, got:\n%s", output)
	}
}

func runPromoteArgs(root, version, base string) []string {
	return []string{"--root", root, "--version", version, "--profile", "public", "--registry", base}
}

// TestPromotionMovesPlatformLatestBeforeWrapperLast is coverage row 4: once
// the complete live set reverifies, promote moves platform latest tags first
// and the wrapper's latest tag strictly last.
func TestPromotionMovesPlatformLatestBeforeWrapperLast(t *testing.T) {
	const version = "9.9.13"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	writeApprovedSet(t, root, ordered, nil)
	base, requestFile := startRegistry(t, root, filepath.Join(t.TempDir(), "store"))
	stagedFlowReadyToPromote(t, root, version, base)

	exitCode, output := runRelease(t, root, "promote", runPromoteArgs(root, version, base), nil)
	if exitCode != 0 {
		t.Fatalf("promote exit=%d:\n%s", exitCode, output)
	}

	lines := requestLines(t, requestFile)
	platformTag := regexp.MustCompile(`^DIST-TAG-ADD @redbench/\S+ latest=`)
	wrapperTag := regexp.MustCompile(`^DIST-TAG-ADD redbench latest=`)
	lastPlatform, wrapperIndex := -1, -1
	for i, line := range lines {
		if platformTag.MatchString(line) {
			lastPlatform = i
		}
		if wrapperTag.MatchString(line) {
			wrapperIndex = i
		}
	}
	if lastPlatform == -1 || wrapperIndex == -1 {
		t.Fatalf("expected both platform and wrapper latest DIST-TAG-ADD calls:\n%s", strings.Join(lines, "\n"))
	}
	if wrapperIndex < lastPlatform {
		t.Fatalf("wrapper latest was promoted before a platform package; request log:\n%s", strings.Join(lines, "\n"))
	}

	record := readPublicationRecord(t, root)
	if record["result"] != "success" {
		t.Fatalf("record result = %v, want success after promote", record["result"])
	}
	if !strings.Contains(output, "release-complete") {
		t.Fatalf("expected next_action release-complete, got:\n%s", output)
	}
}

// TestPromotionRefusesBeforeFullReverification is coverage row 4's mutation
// direction: promote must refuse (and issue no dist-tag-add calls at all)
// when the complete live set has not yet reverified — e.g. the wrapper is
// still only staged, not approved.
func TestPromotionRefusesBeforeFullReverification(t *testing.T) {
	const version = "9.9.14"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	writeApprovedSet(t, root, ordered, nil)
	base, requestFile := startRegistry(t, root, filepath.Join(t.TempDir(), "store"))

	exitCode, output := runRelease(t, root, "submit", stagedSubmitArgs(root, version, base), nil)
	if exitCode != 0 {
		t.Fatalf("staged submit exit=%d:\n%s", exitCode, output)
	}
	record := readPublicationRecord(t, root)
	transitions, _ := record["transitions"].([]any)
	var platformsApproved int
	for _, raw := range transitions {
		transition, _ := raw.(map[string]any)
		if transition["action"] != "stage" || transition["package"] == "redbench" {
			continue
		}
		stageID, _ := transition["stage_id"].(string)
		if stageID != "" {
			httpApprove(t, base, stageID)
			platformsApproved++
		}
	}
	if platformsApproved == 0 {
		t.Fatal("expected at least one platform stage id to approve")
	}
	// Wrapper is deliberately left unapproved: promote must refuse.

	exitCode, output = runRelease(t, root, "promote", runPromoteArgs(root, version, base), nil)
	if exitCode == 0 {
		t.Fatalf("promote succeeded before the complete set reverified:\n%s", output)
	}

	lines := requestLines(t, requestFile)
	for _, line := range lines {
		if strings.Contains(line, "DIST-TAG-ADD") && strings.Contains(line, "latest") {
			t.Fatalf("promote issued a latest dist-tag-add before the complete set reverified:\n%s", strings.Join(lines, "\n"))
		}
	}
}

// TestRollbackRemovesCandidateTagsPreservesLatestNoUnpublish is coverage row
// 5: rollback removes the candidate tag from every package, deprecates the
// bad version, never touches "latest", and never issues an unpublish call.
func TestRollbackRemovesCandidateTagsPreservesLatestNoUnpublish(t *testing.T) {
	const version = "9.9.15"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	writeApprovedSet(t, root, ordered, nil)
	base, requestFile := startRegistry(t, root, filepath.Join(t.TempDir(), "store"))

	exitCode, output := runRelease(t, root, "submit", []string{"--root", root, "--version", version, "--profile", "public", "--registry", base}, nil)
	if exitCode != 0 {
		t.Fatalf("first-publication submit exit=%d:\n%s", exitCode, output)
	}

	exitCode, output = runRelease(t, root, "rollback", []string{"--root", root, "--version", version, "--profile", "public", "--registry", base, "--message", "recovering from a bad release"}, nil)
	if exitCode != 0 {
		t.Fatalf("rollback exit=%d:\n%s", exitCode, output)
	}

	lines := requestLines(t, requestFile)
	tag := "candidate-" + version
	tagRemoves, deprecates := 0, 0
	for _, line := range lines {
		if strings.Contains(line, "REJECT UNPUBLISH") {
			t.Fatalf("rollback issued an unpublish-shaped request:\n%s", strings.Join(lines, "\n"))
		}
		if strings.HasPrefix(line, "DIST-TAG-ADD") && strings.Contains(line, "latest") {
			t.Fatalf("rollback touched the latest dist-tag, which it must never do:\n%s", strings.Join(lines, "\n"))
		}
		if strings.HasPrefix(line, fmt.Sprintf("DIST-TAG-RM ")) && strings.Contains(line, tag) {
			tagRemoves++
		}
		if strings.HasPrefix(line, "DEPRECATE ") {
			deprecates++
		}
	}
	if tagRemoves != len(ordered) {
		t.Fatalf("expected %d candidate-tag removals, got %d:\n%s", len(ordered), tagRemoves, strings.Join(lines, "\n"))
	}
	if deprecates != len(ordered) {
		t.Fatalf("expected %d deprecate calls (every package was live), got %d:\n%s", len(ordered), deprecates, strings.Join(lines, "\n"))
	}

	record := readPublicationRecord(t, root)
	if record["result"] != "rolled_back" {
		t.Fatalf("record result = %v, want rolled_back", record["result"])
	}
}

// TestPublicationRecordCarriesCompleteTransitionsWithNoCredential is coverage
// row 6: the durable record references the release-index digest, carries
// every ordered transition with the required fields, and never lets an
// ambient credential (e.g. NPM_TOKEN) leak into the record bytes.
func TestPublicationRecordCarriesCompleteTransitionsWithNoCredential(t *testing.T) {
	const version = "9.9.16"
	const secretToken = "npm_super-secret-token-should-never-leak-93f7"
	root := copyPublicationScripts(t)
	ordered := releasePlanArtifacts(t, root, version)
	writeApprovedSet(t, root, ordered, nil)
	base, _ := startRegistry(t, root, filepath.Join(t.TempDir(), "store"))

	exitCode, output := runRelease(t, root, "submit",
		[]string{"--root", root, "--version", version, "--profile", "public", "--registry", base},
		map[string]string{"NPM_TOKEN": secretToken})
	if exitCode != 0 {
		t.Fatalf("submit exit=%d:\n%s", exitCode, output)
	}

	rawRecord, err := readRawRecord(t, root)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(rawRecord, secretToken) {
		t.Fatalf("publication record leaked the ambient NPM_TOKEN credential:\n%s", rawRecord)
	}

	record := readPublicationRecord(t, root)
	if record["schema_version"] == nil {
		t.Fatal("record is missing schema_version")
	}
	if record["release_index_sha256"] == "" || record["release_index_sha256"] == nil {
		t.Fatal("record is missing release_index_sha256")
	}
	if record["path"] != "first" {
		t.Fatalf("record path = %v, want first", record["path"])
	}
	if record["profile"] != "public" {
		t.Fatalf("record profile = %v, want public", record["profile"])
	}
	if record["result"] != "success" {
		t.Fatalf("record result = %v, want success", record["result"])
	}

	provenance, _ := record["provenance"].([]any)
	if len(provenance) != len(ordered) {
		t.Fatalf("expected %d provenance entries, got %d", len(ordered), len(provenance))
	}
	for _, raw := range provenance {
		entry, _ := raw.(map[string]any)
		if entry["package"] == "" || entry["package"] == nil {
			t.Fatalf("provenance entry missing package: %+v", entry)
		}
		if entry["sha256"] == "" || entry["sha256"] == nil {
			t.Fatalf("provenance entry missing sha256: %+v", entry)
		}
	}

	transitions, _ := record["transitions"].([]any)
	if len(transitions) == 0 {
		t.Fatal("record has no transitions")
	}
	for _, raw := range transitions {
		transition, _ := raw.(map[string]any)
		for _, field := range []string{"package", "version", "action", "auth_mode", "result", "timestamp"} {
			if v, ok := transition[field]; !ok || v == "" {
				t.Fatalf("transition missing required field %q: %+v", field, transition)
			}
		}
	}
}

func readRawRecord(t *testing.T, root string) (string, error) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "dist", "publication", "publication-record.json"))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
