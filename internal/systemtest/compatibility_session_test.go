//go:build system

package systemtest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/compatibility"
)

func TestCompatibilityResume(t *testing.T) {
	fixture := newCompatibilitySession(t)
	hook := filepath.Join(fixture.repo, ".bench", "hooks", "session-start.sh")
	for _, event := range []string{"startup", "resume"} {
		config := filepath.Join(fixture.homes[compatibility.CodexDesktop], "config.toml")
		if err := os.WriteFile(config, []byte("# "+event+" context\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		result := owner.runWithInput(fixture.repo, fixture.overrides, `{"source":"`+event+`"}`, "bash", hook)
		output := compatibilityOutput(t, result)
		if result.code != 0 || !strings.Contains(output, "on start or resume") || !strings.Contains(output, "repository-wrapper") {
			t.Fatalf("%s omitted active-interface obligation: %#v", event, result)
		}
	}
}

func TestCompatibilityOtherHosts(t *testing.T) {
	for _, host := range []string{"linux/amd64", "darwin/arm64"} {
		report := compatibility.Inspect(compatibility.Input{Context: compatibility.Context{
			Interface: compatibility.CodexCLI, Environment: compatibility.Fact{Value: host, Source: "host fixture"},
		}})
		found := false
		for _, row := range report.Context {
			if row.Field == "execution-environment" && row.Value == host {
				found = true
			}
		}
		if !found {
			t.Fatalf("host context lost: %#v", report)
		}
		session := compatibility.SessionReport(compatibility.SessionRequest{
			Operation: "diagnose", Session: compatibility.Session{Context: compatibility.Context{Environment: compatibility.Fact{Value: host}}},
		})
		for _, row := range session.Checks {
			if row.Check == "supported-environment" && row.State == compatibility.StateFailed {
				t.Fatalf("existing host refused: %#v", session)
			}
		}
	}
}

func TestCompatibilityPermissionConflict(t *testing.T) {
	fixture := newCompatibilitySession(t)
	path := filepath.Join(fixture.homes[compatibility.CodexDesktop], "config.toml")
	policy := []byte("approval_policy = \"untrusted\"\nsandbox_mode = \"read-only\"\n")
	if err := os.WriteFile(path, policy, 0o600); err != nil {
		t.Fatal(err)
	}
	result := fixture.run(t, compatibility.CodexDesktop)
	output := compatibilityOutput(t, result)
	if result.code != 1 || !strings.Contains(output, "effective-configuration,unknown") || !strings.Contains(output, "permission-policy") {
		t.Fatalf("unverified permission context was hidden: %#v", result)
	}
	after, err := os.ReadFile(path)
	if err != nil || string(after) != string(policy) {
		t.Fatalf("inspection broadened configured policy: %q, %v", after, err)
	}
}

func TestCompatibilityPayload(t *testing.T) {
	fixture := newCompatibilitySession(t)
	for _, rel := range []string{".bench/BENCH.md", ".bench/BENCH-reference.md"} {
		installed, err := os.ReadFile(filepath.Join(fixture.repo, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatal(err)
		}
		canonical, err := os.ReadFile(filepath.Join(owner.kit, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatal(err)
		}
		if string(installed) != string(canonical) || !strings.Contains(string(installed), "Session compatibility") {
			t.Fatalf("installed session instructions differ from the payload: %s", rel)
		}
	}
}
