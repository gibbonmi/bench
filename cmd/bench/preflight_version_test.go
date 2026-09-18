package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
	toonlib "github.com/toon-format/toon-go"
)

func TestPreflightReviewChargeUsesCurrentVersion(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	t.Chdir(root)
	writeAXIFixture(t, filepath.Join(root, "go.mod"), "module example.com/versiontest\n\ngo 1.25\n")
	const spec = "# Example\n\nStatus: staged\n\n## User stories\n\n" +
		"1. As a reviewer, I want evidence, so that I can inspect it.\n\n" +
		"### Acceptance coverage map\n\n" +
		"| row | story | behavior | seam | why it catches the failure |\n" +
		"|---|---|---|---|---|\n" +
		"| PF1 | 1 | evidence is complete | command | catches omitted evidence |\n\n" +
		"## Ownership fences\n\n- `target/`\n- `outside/`\n- `reviews/example.md`\n\n" +
		// The review charge grades the completion plan the checkpoint later reads,
		// so a staged spec states it.
		"## Completion plan\n\n```bench-completion-plan\n" +
		`{"version":1,"chunks":[{"id":"c1","tickets":["one.md"],` +
		`"verification":[{"id":"tests","command":"go test ./..."}]}],` +
		`"final_verification":[{"id":"acceptance","command":"go test ./..."}]}` +
		"\n```\n"
	writeAXIFixture(t, filepath.Join(root, "specs/example/spec.md"), spec)
	writeAXIFixture(t, filepath.Join(root, "specs/example/tickets/one.md"), `# One

Blocked by: none
Writes: target
Covers: PF1

## What to build

Build it.

## Acceptance

- [ ] It is built.
`)
	for path, body := range map[string]string{
		".agents/skills/bench-craft-review/SKILL.md":                              "# Review skill\n",
		".agents/commands/bench-review-implementation.md":                         "# Review phase\n",
		".agents/skills/bench-craft-delegate/SKILL.md":                            "# Delegate skill\n",
		".agents/skills/bench-craft-delegate/references/delegation-discipline.md": "# Delegate procedure\n",
		"target/target.go": "package target\n\nfunc Changed() int { return 0 }\n",
		"outside/user.go": "package outside\n\nimport \"example.com/versiontest/target\"\n\n" +
			"func Use() int { return target.Changed() }\n",
	} {
		writeAXIFixture(t, filepath.Join(root, filepath.FromSlash(path)), body)
	}
	runAXIGit(t, "-C", root, "add", ".")
	runAXIGit(t, "-C", root, "commit", "-q", "-m", "base")
	base := trimmedAXIGit(t, root, "rev-parse", "HEAD")
	runAXIGit(t, "-C", root, "checkout", "-q", "-b", "feature")
	writeAXIFixture(t, filepath.Join(root, "target/target.go"), "package target\n\nfunc Changed() int { return 1 }\n")
	runAXIGit(t, "-C", root, "add", "target/target.go")
	runAXIGit(t, "-C", root, "commit", "-q", "-m", "source")
	tip := trimmedAXIGit(t, root, "rev-parse", "HEAD")
	const assignment = "00000000000000000000000000000001"
	err := intent.PutAssignment(root, intent.Assignment{
		Schema: intent.AssignmentRecordSchema, ID: assignment,
		OwnerID: "00000000000000000000000000000002",
		Request: intent.RequestDigest("preflight-version"), Label: "preflight-version",
		Start: tip, Branch: intent.AssignmentBranchRef("00000000000000000000000000000002", assignment),
		Worktree: root, State: intent.StateActive,
	})
	if err != nil {
		t.Fatal(err)
	}

	stdout := runPreflight(t, []string{
		"preflight", "review", "example", "--charge", "--base", base, "--source-tip", tip,
	})
	prepared := tableValue(t, decodeVersionTOON(t, stdout), "prepared")
	if len(prepared) != 1 {
		t.Fatalf("prepared rows = %d:\n%s", len(prepared), stdout)
	}
	identity, ok := prepared[0]["evidence"].(string)
	if !ok {
		t.Fatalf("prepared evidence = %#v", prepared[0]["evidence"])
	}
	// The manifest opens the default stream, so the producer rows arrive before any
	// source page. Every generated capture must name the running executable's version.
	manifest := readManifestStream(t, identity)
	producers := tableValue(t, decodeVersionTOON(t, manifest), "producers")
	if len(producers) == 0 {
		t.Fatalf("manifest declares no producer:\n%s", manifest)
	}
	for _, row := range producers {
		if row["version"] != version {
			t.Fatalf("producer %#v, want current version %q", row, version)
		}
	}
}

// runPreflight runs one preflight invocation through the real command and returns its
// stdout, so the assertion grades the version the executable itself carries.
func runPreflight(t *testing.T, args []string) string {
	t.Helper()
	var stdout, stderr bytes.Buffer
	if code := (Command{Stdout: &stdout, Stderr: &stderr}).Run(args); code != 0 || stderr.Len() != 0 {
		t.Fatalf("%v = (%d, stderr=%q):\n%s", args, code, stderr.String(), stdout.String())
	}
	return stdout.String()
}

// readManifestStream follows the default stream from its first read until the manifest
// fragments end, and returns the reconstructed manifest bytes.
func readManifestStream(t *testing.T, identity string) string {
	t.Helper()
	var manifest strings.Builder
	args := []string{"preflight", "evidence", identity}
	for i := 0; i < 10000; i++ {
		rows := tableValue(t, decodeVersionTOON(t, runPreflight(t, args)), "page")
		if len(rows) != 1 {
			t.Fatalf("page rows = %d", len(rows))
		}
		row := rows[0]
		if row["stream"] != "manifest" {
			return manifest.String()
		}
		content, ok := row["content"].(string)
		if !ok {
			t.Fatalf("page content = %#v", row["content"])
		}
		manifest.WriteString(content)
		next, _ := row["next"].(string)
		if next == "" {
			return manifest.String()
		}
		args = strings.Fields(next)[1:]
	}
	t.Fatal("manifest traversal did not end")
	return ""
}

func trimmedAXIGit(t *testing.T, root string, args ...string) string {
	t.Helper()
	all := append([]string{"-C", root}, args...)
	return string(bytes.TrimSpace([]byte(runAXIGit(t, all...))))
}

func decodeVersionTOON(t *testing.T, output string) map[string]any {
	t.Helper()
	decoded, err := toonlib.DecodeString(output)
	if err != nil {
		t.Fatalf("decode TOON: %v\n%s", err, output)
	}
	document, ok := decoded.(map[string]any)
	if !ok {
		t.Fatalf("decoded = %T, want object", decoded)
	}
	return document
}

func tableValue(t *testing.T, document map[string]any, name string) []map[string]any {
	t.Helper()
	values, ok := document[name].([]any)
	if !ok {
		t.Fatalf("%s = %T, want table", name, document[name])
	}
	rows := make([]map[string]any, len(values))
	for i, value := range values {
		row, ok := value.(map[string]any)
		if !ok {
			t.Fatalf("%s row %d = %T, want object", name, i, value)
		}
		rows[i] = row
	}
	return rows
}

func evidenceValue(t *testing.T, document map[string]any, source string) string {
	t.Helper()
	for _, row := range tableValue(t, document, "evidence") {
		if row["source"] == source {
			value, ok := row["content"].(string)
			if !ok {
				t.Fatalf("%s evidence = %T, want string", source, row["content"])
			}
			return value
		}
	}
	t.Fatalf("evidence omitted %s", source)
	return ""
}
