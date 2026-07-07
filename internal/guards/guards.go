// Package guards ports `bench guards` (and `--brief`): every deny-capable guard's
// manifest aggregated into a `guards[N]{guard,boundary,denies,wired}:` TOON table, so the
// block surface is learnable without a collision. Guards are discovered by convention
// (each .bench/hooks/*.sh and the installed git pre-push hook); each guard's --describe
// is read under a hard time bound so a hook that
// ignores --describe cannot stall aggregation, and an unmanaged pre-push is never
// executed — running an unknown hook's body just to read a manifest is the collision
// this surface avoids.
//
// The `wired` cell names which harness configs actually reference a hook script, so the
// deny surface the reader sees matches the hooks that can fire here: it is derived (never
// declared) by scanning .claude/settings.json and .codex/hooks.json for the script's
// relative path token, the same substring convention conformance uses. The pre-push hook
// is wired through git rather than a harness config, so its wired cell is the constant
// `git` and its install posture (managed/unmanaged/not installed) stays in the denies
// column.
package guards

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

// describeTimeout bounds each guard's --describe. Replaces the shell's coreutils
// `timeout 5` / watchdog with exec.CommandContext. The deadline kills the hook's
// whole process group (its own group via Setpgid), so a grandchild (a `sleep`)
// dies with it instead of holding the stdout pipe; waitGrace is the backstop that
// forces Wait to return even if something in the group survives the kill.
const (
	describeTimeout = 5 * time.Second
	waitGrace       = 3 * time.Second
)

// guardDescribe runs `bash <path> --describe` with stdin on the null device under the
// time bound, returning its stdout and exit code — or 124 when the bound trips.
func guardDescribe(path string) (string, int) {
	ctx, cancel := context.WithTimeout(context.Background(), describeTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", path, "--describe")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	cmd.WaitDelay = waitGrace
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", 124
	}
	if err == nil {
		return out.String(), 0
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return out.String(), ee.ExitCode()
	}
	return out.String(), 1
}

// manifestField pulls one `key: value` field from a manifest blob; "" if absent.
func manifestField(out, key string) string {
	prefix := key + ": "
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimPrefix(line, prefix)
		}
	}
	return ""
}

// guardRow builds the row for one discovered guard, or (nil,false) when the guard is
// informational (`denies: nothing`) and excluded. A guard that does not answer
// --describe (nonzero, timeout, or a manifest missing any of the three fields) gets a
// definitive `no manifest` row under its fallback name rather than a silent omission.
func guardRow(path, fallback string) ([]string, bool) {
	out, rc := guardDescribe(path)
	switch {
	case rc == 124:
		return []string{fallback, "", "no manifest (timed out)"}, true
	case rc != 0:
		return []string{fallback, "", "no manifest"}, true
	}
	name := manifestField(out, "name")
	boundary := manifestField(out, "boundary")
	denies := manifestField(out, "denies")
	if name == "" || boundary == "" || denies == "" {
		return []string{fallback, "", "no manifest"}, true
	}
	if denies == "nothing (informational)" {
		return nil, false
	}
	return []string{name, boundary, denies}, true
}

// Rows discovers every guard and returns its row: the hook scripts and the git
// pre-push hook, in that order. The shift-adapter line enforcement lives in `bench
// resolve-model` (internal/lines), not a `.sh` manifest — `bench guards` aggregates
// hook scripts only and grows no second answerer, so that enforcement has no row here.
func Rows(root string) [][]string {
	var rows [][]string
	add := func(path, fallback, wired string) {
		if r, emit := guardRow(path, fallback); emit {
			rows = append(rows, append(r, wired))
		}
	}
	hooksDir := filepath.Join(root, ".bench", "hooks")
	if entries, err := os.ReadDir(hooksDir); err == nil {
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".sh") {
				continue
			}
			scriptRel := ".bench/hooks/" + e.Name()
			add(filepath.Join(hooksDir, e.Name()), strings.TrimSuffix(e.Name(), ".sh"), wiredHarnesses(root, scriptRel))
		}
	}
	// The pre-push guard is wired through git, not a harness config, so its wired
	// cell is the constant `git`; its install posture stays in the denies column.
	return append(rows, withWired(prePushRow(root), "git")...)
}

// wiredHarnesses names the harness configs that reference scriptRel — its relative
// path token — scanning .claude/settings.json and .codex/hooks.json with the same
// substring convention conformance uses. An absent config contributes nothing (a repo
// without .codex/ cannot wire Codex); an unparseable config scans as not-wired, because
// guards is a read-only wiring reporter and the JSON-validity conformance family owns
// malformedness. Returns "claude", "codex", the comma-joined "claude,codex", or the
// definitive "none" — never a blank cell.
func wiredHarnesses(root, scriptRel string) string {
	token := []byte(scriptRel)
	var wired []string
	for _, cfg := range []struct{ name, rel string }{
		{"claude", filepath.Join(".claude", "settings.json")},
		{"codex", filepath.Join(".codex", "hooks.json")},
	} {
		content, err := os.ReadFile(filepath.Join(root, cfg.rel))
		if err != nil || !json.Valid(content) {
			continue
		}
		if bytes.Contains(content, token) {
			wired = append(wired, cfg.name)
		}
	}
	if len(wired) == 0 {
		return "none"
	}
	return strings.Join(wired, ",")
}

// withWired appends the wired cell to every row, keeping the pre-push wiring channel
// (git) out of the config-scan path that the hook-script rows take.
func withWired(rows [][]string, wired string) [][]string {
	for i := range rows {
		rows[i] = append(rows[i], wired)
	}
	return rows
}

// prePushRow resolves the installed git pre-push hook. A managed hook (carrying the
// bench marker) is read for its manifest; a marker-less foreign hook is reported
// `unmanaged` and never executed; an absent hook is a definitive `not installed`.
func prePushRow(root string) [][]string {
	notInstalled := [][]string{{"pre-push", "", "not installed"}}
	hooksGit, err := git.Output("-C", root, "rev-parse", "--git-path", "hooks")
	if err != nil || hooksGit == "" {
		return notInstalled
	}
	if !filepath.IsAbs(hooksGit) {
		hooksGit = filepath.Join(root, hooksGit)
	}
	prepush := filepath.Join(hooksGit, "pre-push")
	if !fileExists(prepush) {
		return notInstalled
	}
	content, _ := os.ReadFile(prepush)
	if !bytes.Contains(content, []byte("bench:managed-pre-push")) {
		return [][]string{{"pre-push", "", "unmanaged (no manifest)"}}
	}
	if r, emit := guardRow(prepush, "pre-push"); emit {
		return [][]string{r}
	}
	return nil
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

// Command implements `bench guards [--brief]`. --brief emits one plain line per
// deny-capable guard plus exactly one footer — the surface session-start injects.
func Command(args []string) (string, int) {
	brief := false
	switch {
	case len(args) == 0:
	case args[0] == "--brief":
		brief = true
	case args[0] == "-h" || args[0] == "--help":
		return "usage: bench guards [--brief]\n", 0
	default:
		return toon.Usage("bench guards", args[0]) + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	rows := Rows(root)
	if brief {
		var b strings.Builder
		for _, r := range rows {
			fmt.Fprintf(&b, "%s: %s [wired: %s]\n", r[0], r[2], r[3])
		}
		b.WriteString("full manifests: bench guards\n")
		return b.String(), 0
	}
	out, err := toon.Table("guards", []string{"guard", "boundary", "denies", "wired"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out, 0
}
