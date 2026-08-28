package releaseevidence

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// The release plan states a native_proof boolean per target. These rows are the
// partition that statement admits: proven, unproven, absent, and non-boolean.
const (
	provenLinuxRow    = `{"os":"linux","arch":"x64","goos":"linux","goarch":"amd64","runner":"ubuntu-24.04","native_proof":true}`
	unprovenDarwinRow = `{"os":"darwin","arch":"arm64","goos":"darwin","goarch":"arm64","runner":"macos-15","native_proof":false}`
	provenDarwinRow   = `{"os":"darwin","arch":"arm64","goos":"darwin","goarch":"arm64","runner":"macos-15","native_proof":true}`
	silentDarwinRow   = `{"os":"darwin","arch":"arm64","goos":"darwin","goarch":"arm64","runner":"macos-15"}`
	stringDarwinRow   = `{"os":"darwin","arch":"arm64","goos":"darwin","goarch":"arm64","runner":"macos-15","native_proof":"true"}`
)

// repositoryRoot locates the repository root two directories above this package
// (internal/releaseevidence). `go test` always runs with the package directory
// as its working directory, so this is the one way a test here reaches scripts/.
func repositoryRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(filepath.Join(cwd, "..", ".."))
}

// planFixtureRoot writes a scratch root holding the repository's release-plan.mjs
// beside a plan built from targets. release-plan.mjs reads the plan from the root
// it is given, so a fixture grades the reader without touching the working tree.
func planFixtureRoot(t *testing.T, targets string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	reader, err := os.ReadFile(filepath.Join(repositoryRoot(t), "scripts", "release-plan.mjs"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "scripts", "release-plan.mjs"), reader, 0o644); err != nil {
		t.Fatal(err)
	}
	plan := `{"schema_version":1,"targets":[` + targets + `],"archive_entries":[{"path":"bin/bench","mode":"0755","kind":"binary"}]}`
	if err := os.WriteFile(filepath.Join(root, "scripts", "release-plan.json"), []byte(plan), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func runReleasePlan(t *testing.T, root string, arguments ...string) (string, error) {
	t.Helper()
	out, err := exec.Command("node", append([]string{filepath.Join(root, "scripts", "release-plan.mjs"), root}, arguments...)...).Output()
	return string(out), err
}

func planFailure(err error) string {
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return string(exit.Stderr)
	}
	return err.Error()
}

// matrixRows decodes one matrix view into its raw rows, so a test can grade both
// the row set and the field set each row carries.
func matrixRows(t *testing.T, root, command string) []map[string]any {
	t.Helper()
	out, err := runReleasePlan(t, root, command)
	if err != nil {
		t.Fatalf("%s: %v\n%s", command, err, planFailure(err))
	}
	var matrix struct {
		Include []map[string]any `json:"include"`
	}
	if err := json.Unmarshal([]byte(out), &matrix); err != nil {
		t.Fatalf("decode %s: %v", command, err)
	}
	return matrix.Include
}

func matrixTargets(rows []map[string]any) []string {
	names := make([]string, 0, len(rows))
	for _, row := range rows {
		names = append(names, row["os"].(string)+"-"+row["arch"].(string))
	}
	sort.Strings(names)
	return names
}

func TestReleasePlanReaderRejectsUnstatedNativeProof(t *testing.T) {
	for name, targets := range map[string]string{
		"absent native_proof":      silentDarwinRow + "," + provenLinuxRow,
		"non-boolean native_proof": stringDarwinRow + "," + provenLinuxRow,
	} {
		t.Run(name, func(t *testing.T) {
			root := planFixtureRoot(t, targets)
			out, err := runReleasePlan(t, root, "normalized-json")
			if err == nil {
				t.Fatalf("plan with a %s was accepted: %s", name, out)
			}
			if !strings.Contains(planFailure(err), "release plan target is invalid") {
				t.Fatalf("plan with a %s failed for another reason: %s", name, planFailure(err))
			}
		})
	}
}

func TestReleasePlanReaderRejectsEmptyProvenSet(t *testing.T) {
	root := planFixtureRoot(t, unprovenDarwinRow)
	out, err := runReleasePlan(t, root, "normalized-json")
	if err == nil {
		t.Fatalf("plan with no proven target was accepted: %s", out)
	}
	if !strings.Contains(planFailure(err), "release plan has no proven target") {
		t.Fatalf("plan with no proven target failed for another reason: %s", planFailure(err))
	}
}

func TestMatrixViewsSeparateProvenTargetsAndHideTheField(t *testing.T) {
	root := planFixtureRoot(t, unprovenDarwinRow+","+provenLinuxRow)
	shipped, proven := matrixRows(t, root, "matrix-json"), matrixRows(t, root, "proof-matrix-json")
	if got := matrixTargets(shipped); strings.Join(got, " ") != "darwin-arm64 linux-x64" {
		t.Fatalf("matrix-json dropped a shipped target: %v", got)
	}
	if got := matrixTargets(proven); strings.Join(got, " ") != "linux-x64" {
		t.Fatalf("proof-matrix-json carries an unproven target: %v", got)
	}
	for _, row := range append(append([]map[string]any{}, shipped...), proven...) {
		if _, leaked := row["native_proof"]; leaked {
			t.Fatalf("matrix row leaks native_proof into a workflow variable: %v", row)
		}
		if len(row) != 5 {
			t.Fatalf("matrix row carries fields beyond the five transport fields: %v", row)
		}
	}
}

func TestProvenViewsFollowTheFieldNotTheOperatingSystem(t *testing.T) {
	root := planFixtureRoot(t, provenDarwinRow+","+provenLinuxRow)
	if got := matrixTargets(matrixRows(t, root, "proof-matrix-json")); strings.Join(got, " ") != "darwin-arm64 linux-x64" {
		t.Fatalf("a proven Darwin target is absent from proof-matrix-json: %v", got)
	}
	rows, err := runReleasePlan(t, root, "proof-targets")
	if err != nil {
		t.Fatalf("proof-targets: %v\n%s", err, planFailure(err))
	}
	if !strings.Contains(rows, "darwin\tarm64\tdarwin\tarm64\tmacos-15\n") {
		t.Fatalf("a proven Darwin target is absent from proof-targets: %q", rows)
	}
	row, err := runReleasePlan(t, root, "proof-target", "darwin", "arm64")
	if err != nil {
		t.Fatalf("proof-target: %v\n%s", err, planFailure(err))
	}
	if row != "darwin\tarm64\tdarwin\tarm64\tmacos-15\n" {
		t.Fatalf("proof-target does not name a proven Darwin target: %q", row)
	}
}

// TestShippedReleasePlanDecodes grades that the Go reader still decodes the
// repository's own plan. That reader disallows unknown fields, so a plan field
// the struct does not carry reddens every release-evidence path. The plan's target
// rows stay unstated here, because scripts/release-plan.json is their one source.
func TestShippedReleasePlanDecodes(t *testing.T) {
	if _, err := readReleasePlan(repositoryRoot(t)); err != nil {
		t.Fatalf("shipped release plan is undecodable: %v", err)
	}
}
