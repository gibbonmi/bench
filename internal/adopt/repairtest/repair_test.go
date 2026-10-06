package repairtest

import (
	"bufio"
	"bytes"
	"github.com/gibbonmi/bench/internal/compatibility/compatibilitytest"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/adopt"
)

func TestCompatibilityManagedRepair(t *testing.T) {
	s := linkedSession(t)
	hook := filepath.Join(s.root, ".codex", "hooks.json")
	before, err := os.ReadFile(hook)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := s.doctor("--fix")
	if code != 1 {
		t.Fatalf("repair must leave unknown live obligations visible: %d, %s, %s", code, stdout, stderr)
	}
	after, err := os.ReadFile(hook)
	if err != nil || !bytes.Equal(after, before) {
		t.Fatalf("canonical hook was not restored: %v; %s %s", err, stdout, stderr)
	}
}

func TestCompatibilityModifiedConflict(t *testing.T) {
	for _, kind := range []string{"payload content", "payload mode", "managed pre-push", "ownership manifest"} {
		t.Run(kind, func(t *testing.T) {
			s := linkedSession(t)
			path := filepath.Join(s.root, ".codex", "hooks.json")
			if kind == "managed pre-push" {
				path = filepath.Join(s.root, ".git", "hooks", "pre-push")
			}
			if kind == "ownership manifest" {
				path = filepath.Join(s.root, ".bench", "link-manifest.tsv")
			}
			original, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			changed := append(append([]byte{}, original...), []byte("\nuser modified this managed asset\n")...)
			if kind == "payload mode" {
				changed = original
			}
			if err := os.WriteFile(path, changed, 0o600); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(path, 0o600); err != nil {
				t.Fatal(err)
			}
			code, stdout, stderr := s.doctor("--fix")
			data, err := os.ReadFile(path)
			info, statErr := os.Stat(path)
			if err != nil || statErr != nil || !bytes.Equal(data, changed) || info.Mode().Perm() != 0o600 {
				t.Fatalf("repair overwrote user bytes or mode: %q, %v, %v", data, err, statErr)
			}
			if code != 1 || !strings.Contains(stdout, "modified-managed") {
				t.Fatalf("modified conflict was not reported with a failed local result: %d, %s %s", code, stdout, stderr)
			}
		})
	}
}

func TestCompatibilityUndoCreated(t *testing.T) {
	s := linkedSession(t)
	hook := filepath.Join(s.root, ".codex", "hooks.json")
	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := s.doctor("--fix")
	if code != 1 || stderr != "" {
		t.Fatalf("repair = %d, %s %s", code, stdout, stderr)
	}
	id := repairID(t, stdout)
	if _, err := os.Stat(hook); err != nil {
		t.Fatalf("repair did not create target: %v", err)
	}
	code, stdout, stderr = s.doctor("--undo", id)
	if code != 1 || stderr != "" {
		t.Fatalf("undo = %d, %s %s", code, stdout, stderr)
	}
	if _, err := os.Lstat(hook); !os.IsNotExist(err) {
		t.Fatalf("undo did not restore absence: %v", err)
	}
}

func TestCompatibilityRepairPrivacy(t *testing.T) {
	s := linkedSession(t)
	secret := "project-private-content:" + s.root
	agents := filepath.Join(s.root, "AGENTS.md")
	if err := os.WriteFile(agents, []byte("# Project\n"+secret+"\n"), 0o640); err != nil {
		t.Fatal(err)
	}
	configHome := os.Getenv("CODEX_HOME")
	if err := os.MkdirAll(configHome, 0o700); err != nil {
		t.Fatal(err)
	}
	credential := "private-credential:" + s.root
	if err := os.WriteFile(filepath.Join(configHome, "credentials"), []byte(credential), 0o600); err != nil {
		t.Fatal(err)
	}
	rawEnv := "private-environment:" + s.root
	t.Setenv("PRIVATE_USER_SECRET", rawEnv)
	code, stdout, stderr := s.doctor("--fix")
	if code != 1 || stderr != "" {
		t.Fatalf("repair = %d, %s %s", code, stdout, stderr)
	}
	_ = repairID(t, stdout)
	data, err := os.ReadFile(agents)
	if err != nil || !bytes.Contains(data, []byte(secret)) || !bytes.Contains(data, []byte(adopt.BenchAgentsBlock())) {
		t.Fatalf("managed block repair lost project text or integration: %q, %v", data, err)
	}
	records := filepath.Join(s.home, "bench-state", "compatibility-repairs")
	err = filepath.WalkDir(records, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if info.Mode().Perm() != 0o700 {
				t.Errorf("repair directory mode = %v: %s", info.Mode(), path)
			}
			return nil
		}
		if info.Mode().Perm() != 0o600 {
			t.Errorf("repair file mode = %v: %s", info.Mode(), path)
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, value := range []string{secret, credential, rawEnv} {
			if bytes.Contains(body, []byte(value)) {
				t.Errorf("repair record copied private content: %s", path)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCompatibilityUndo(t *testing.T) {
	s := linkedSession(t)
	hook := filepath.Join(s.root, ".codex", "hooks.json")
	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
	agents := filepath.Join(s.root, "AGENTS.md")
	before := []byte("# Private project\nretain these bytes without a final newline")
	if err := os.WriteFile(agents, before, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(agents, 0o600); err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := s.doctor("--fix")
	if code != 1 || stderr != "" {
		t.Fatalf("repair = %d, %s %s", code, stdout, stderr)
	}
	if _, err := os.Stat(hook); err != nil {
		t.Fatalf("repair omitted second destination: %v", err)
	}
	id := repairID(t, stdout)
	code, stdout, stderr = s.doctor("--undo", id)
	if code != 1 || stderr != "" {
		t.Fatalf("undo = %d, %s %s", code, stdout, stderr)
	}
	if _, err := os.Lstat(hook); !os.IsNotExist(err) {
		t.Fatalf("undo did not restore second destination absence: %v", err)
	}
	data, err := os.ReadFile(agents)
	info, statErr := os.Stat(agents)
	if err != nil || statErr != nil || !bytes.Equal(data, before) || info.Mode().Perm() != 0o600 {
		t.Fatalf("undo did not restore exact project bytes and mode: %q, %v, %v", data, err, statErr)
	}
}

func TestCompatibilityForeignConflict(t *testing.T) {
	s := newSession(t)
	path := filepath.Join(s.root, ".codex", "hooks.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	before := []byte("foreign project hook declaration\n")
	if err := os.WriteFile(path, before, 0o600); err != nil {
		t.Fatal(err)
	}
	s.link(t, 3)
	code, stdout, stderr := s.doctor("--fix")
	data, err := os.ReadFile(path)
	info, statErr := os.Stat(path)
	if err != nil || statErr != nil || !bytes.Equal(data, before) || info.Mode().Perm() != 0o600 {
		t.Fatalf("repair changed foreign bytes or mode: %q, %v, %v", data, err, statErr)
	}
	if code != 1 || !strings.Contains(stdout, "project-owned") {
		t.Fatalf("foreign conflict not reported: %d, %s %s", code, stdout, stderr)
	}
}

func TestCompatibilityPolicyBoundary(t *testing.T) {
	s := linkedSession(t)
	home := os.Getenv("CODEX_HOME")
	if err := os.MkdirAll(home, 0o700); err != nil {
		t.Fatal(err)
	}
	paths := []string{filepath.Join(home, "config.toml"), filepath.Join(home, "trust.json")}
	before := []byte("sandbox_mode = \"read-only\"\ntrust = \"untrusted\"\n")
	for _, path := range paths {
		if err := os.WriteFile(path, before, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	code, stdout, stderr := s.doctor("--fix")
	for _, path := range paths {
		data, err := os.ReadFile(path)
		info, statErr := os.Stat(path)
		if err != nil || statErr != nil || !bytes.Equal(data, before) || info.Mode().Perm() != 0o600 {
			t.Fatalf("repair changed permission or trust settings: %s, %q, %v, %v", path, data, err, statErr)
		}
	}
	if code != 1 || !strings.Contains(stdout, "security-policy,decision,") {
		t.Fatalf("repair omitted security-policy decision: %d, %s %s", code, stdout, stderr)
	}
}

func TestCompatibilityInterruptionBoundary(t *testing.T) {
	s := linkedSession(t)
	process := exec.Command("cat")
	input, err := process.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	output, err := process.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := process.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { input.Close(); process.Process.Kill(); process.Wait() })
	code, stdout, stderr := s.doctor("--fix")
	if _, err := input.Write([]byte("still running\n")); err != nil {
		t.Fatalf("repair interrupted the supervised process: %v", err)
	}
	line, err := bufio.NewReader(output).ReadString('\n')
	if err != nil || line != "still running\n" {
		t.Fatalf("supervised process did not survive: %q, %v", line, err)
	}
	if code != 1 || !strings.Contains(stdout, "process-interruption,decision,") {
		t.Fatalf("repair omitted process-interruption decision: %d, %s %s", code, stdout, stderr)
	}
}

func TestCompatibilityRuntimeBoundary(t *testing.T) {
	s := linkedSession(t)
	home := os.Getenv("CODEX_HOME")
	if err := os.MkdirAll(home, 0o700); err != nil {
		t.Fatal(err)
	}
	runtime := filepath.Join(home, "private-runtime", "missing-launcher")
	config := []byte("cli_path = \"" + runtime + "\"\n")
	configPath := filepath.Join(home, "config.toml")
	if err := os.WriteFile(configPath, config, 0o600); err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := s.doctor("--fix")
	if _, err := os.Lstat(filepath.Dir(runtime)); !os.IsNotExist(err) {
		t.Fatalf("repair created a private runtime directory: %v", err)
	}
	data, err := os.ReadFile(configPath)
	if err != nil || !bytes.Equal(data, config) {
		t.Fatalf("repair changed private launcher selection: %q, %v", data, err)
	}
	if code != 1 || !strings.Contains(stdout, "private-runtime,decision,") {
		t.Fatalf("repair omitted private-runtime decision: %d, %s %s", code, stdout, stderr)
	}
}

func TestCompatibilityRepairIdempotence(t *testing.T) {
	s := linkedSession(t)
	s.track(t)
	hook := filepath.Join(s.root, ".codex", "hooks.json")
	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := s.doctor("--fix")
	if code != 1 || stderr != "" {
		t.Fatalf("first repair: %d %s %s", code, stdout, stderr)
	}
	if _, err := os.Stat(hook); err != nil {
		t.Fatalf("first repair did not restore the asset: %v", err)
	}
	before := compatibilitytest.SnapshotHome(t, s.root)
	code, stdout, stderr = s.doctor("--fix")
	if code != 1 || stderr != "" {
		t.Fatalf("second repair: %d %s %s", code, stdout, stderr)
	}
	if after := compatibilitytest.SnapshotHome(t, s.root); !reflect.DeepEqual(before, after) {
		for name, prior := range before {
			if after[name] != prior {
				t.Fatalf("repeated repair changed tracked bytes or modes at %s", name)
			}
		}
		t.Fatal("repeated repair changed tracked bytes or modes")
	}
}

func TestCompatibilityUndoConflict(t *testing.T) {
	for _, kind := range []string{"edited content", "equivalent replacement"} {
		t.Run(kind, func(t *testing.T) {
			s := linkedSession(t)
			hook := filepath.Join(s.root, ".codex", "hooks.json")
			if err := os.Remove(hook); err != nil {
				t.Fatal(err)
			}
			_, stdout, stderr := s.doctor("--fix")
			if stderr != "" {
				t.Fatal(stderr)
			}
			id := repairID(t, stdout)
			changed := []byte("user edit after repair\n")
			mode := os.FileMode(0o600)
			if kind == "equivalent replacement" {
				var err error
				changed, err = os.ReadFile(hook)
				if err != nil {
					t.Fatal(err)
				}
				info, err := os.Stat(hook)
				if err != nil {
					t.Fatal(err)
				}
				mode = info.Mode()
			}
			replacement := hook
			if kind == "equivalent replacement" {
				replacement = filepath.Join(s.root, "later-hook")
			}
			if err := os.WriteFile(replacement, changed, mode); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(replacement, mode); err != nil {
				t.Fatal(err)
			}
			if replacement != hook {
				if err := os.Rename(replacement, hook); err != nil {
					t.Fatal(err)
				}
			}
			before, err := os.Stat(hook)
			if err != nil {
				t.Fatal(err)
			}
			code, stdout, stderr := s.doctor("--undo", id)
			if code != 1 || !strings.Contains(stderr, "undo conflict") {
				t.Fatalf("stale undo did not report conflict: %d, %s %s", code, stdout, stderr)
			}
			data, err := os.ReadFile(hook)
			info, statErr := os.Stat(hook)
			if err != nil || statErr != nil || !bytes.Equal(data, changed) || info.Mode() != mode || !os.SameFile(before, info) {
				t.Fatalf("undo changed the later replacement: %q, %v, %v", data, err, statErr)
			}
		})
	}
}
