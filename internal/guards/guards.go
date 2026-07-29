// Package guards ports `bench guards` (and `--brief`): every deny-capable guard's
// manifest aggregated into a `guards[N]{guard,boundary,denies,wired}:` TOON table, so the
// block surface is learnable without a collision. Guards are discovered by convention
// (each .bench/hooks/*.sh and the installed git pre-push hook), and their manifests are
// parsed from static leading-comment headers. Discovery reads scripts only as data and
// never executes them.
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
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/adopt"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand —
// arity, flag recognition, `--`, and help all come from there rather than a local switch.
var grammar = usage.Grammar{
	Cmd:   "bench guards",
	Help:  "usage: bench guards [--brief]",
	Flags: []usage.Flag{{Name: "--brief"}},
}

type Manifest map[string]string

type candidate struct {
	path, fallback, wired string
	prepush               bool
}

type ScanResult struct {
	Rows                     [][]string
	Status, Inspected, Total string
	Omitted, Reason          string
}

var guardScanTimeout = bounds.GuardScanTimeout

var enumerateGuards = func(ctx context.Context, root string) ([]candidate, error) {
	var candidates []candidate
	hooksDir := filepath.Join(root, ".bench", "hooks")
	entries, err := os.ReadDir(hooksDir)
	if err == nil {
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".sh") {
				continue
			}
			scriptRel := ".bench/hooks/" + e.Name()
			candidates = append(candidates, candidate{path: filepath.Join(hooksDir, e.Name()), fallback: strings.TrimSuffix(e.Name(), ".sh"), wired: wiredHarnesses(root, scriptRel)})
		}
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	candidates = append(candidates, candidate{fallback: "pre-push", wired: "git", prepush: true})
	return candidates, nil
}

var inspectGuard = func(ctx context.Context, root string, c candidate) [][]string {
	if ctx.Err() != nil {
		return nil
	}
	if c.prepush {
		return withWired(prePushRow(root), c.wired)
	}
	if row, emit := guardRow(c.path, c.fallback); emit {
		return [][]string{append(row, c.wired)}
	}
	return nil
}

var requiredManifestKeys = []string{"name", "boundary", "denies", "why"}

func (m Manifest) MissingRequired() []string {
	var missing []string
	for _, key := range requiredManifestKeys {
		if m[key] == "" {
			missing = append(missing, key)
		}
	}
	return missing
}

// HeaderFields reads the static manifest in path's leading comment block. The first
// occurrence of each key wins; an empty first value remains missing/invalid.
func HeaderFields(path string) (Manifest, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("guard manifest is not a regular file")
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseHeader(string(content)), nil
}

func parseHeader(content string) Manifest {
	fields := Manifest{}
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "#") {
			break
		}
		for _, key := range requiredManifestKeys {
			prefix := "# " + key + ": "
			if _, found := fields[key]; !found && strings.HasPrefix(line, prefix) {
				fields[key] = strings.TrimRight(strings.TrimPrefix(line, prefix), " \t\r")
			}
		}
	}
	return fields
}

func manifestField(content, key string) string { return parseHeader(content)[key] }

// guardRow builds the row for one discovered guard, or (nil,false) when the guard is
// informational and excluded. A guard with an unreadable or incomplete header gets a
// definitive `no manifest` row under its fallback name rather than a silent omission.
func guardRow(path, fallback string) ([]string, bool) {
	fields, err := HeaderFields(path)
	if err != nil {
		return []string{fallback, "", "no manifest"}, true
	}
	if len(fields.MissingRequired()) != 0 {
		return []string{fallback, "", "no manifest"}, true
	}
	name, boundary, denies := fields["name"], fields["boundary"], fields["denies"]
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
	return Scan(context.Background(), root).Rows
}

func Scan(ctx context.Context, root string) ScanResult {
	type enumeration struct {
		candidates []candidate
		err        error
	}
	enumerated := make(chan enumeration, 1)
	go func() {
		candidates, err := enumerateGuards(ctx, root)
		enumerated <- enumeration{candidates: candidates, err: err}
	}()
	var candidates []candidate
	select {
	case result := <-enumerated:
		if result.err != nil && ctx.Err() != nil {
			return ScanResult{Status: "incomplete", Inspected: "0", Total: "unknown", Omitted: "unknown", Reason: "timeout"}
		}
		candidates = result.candidates
	case <-ctx.Done():
		<-enumerated
		return ScanResult{Status: "incomplete", Inspected: "0", Total: "unknown", Omitted: "unknown", Reason: "timeout"}
	}
	var rows [][]string
	inspected := 0
	for _, c := range candidates {
		finished := make(chan [][]string, 1)
		go func(candidate candidate) { finished <- inspectGuard(ctx, root, candidate) }(c)
		select {
		case candidateRows := <-finished:
			inspected++
			rows = append(rows, candidateRows...)
		case <-ctx.Done():
			<-finished
			return ScanResult{Rows: rows, Status: "incomplete", Inspected: strconv.Itoa(inspected), Total: strconv.Itoa(len(candidates)), Omitted: strconv.Itoa(len(candidates) - inspected), Reason: "timeout"}
		}
	}
	return ScanResult{Rows: rows, Status: "complete", Inspected: strconv.Itoa(inspected), Total: strconv.Itoa(len(candidates)), Omitted: "0", Reason: "none"}
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
	if !bytes.Contains(content, []byte(adopt.PrePushMarker)) {
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
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	_, brief := parsed.Flags["--brief"]
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	ctx, cancel := bounds.Context(context.Background(), guardScanTimeout)
	defer cancel()
	scan := Scan(ctx, root)
	rows := scan.Rows
	if brief {
		var b strings.Builder
		for _, r := range rows {
			fmt.Fprintf(&b, "%s: %s [wired: %s]\n", r[0], r[2], r[3])
		}
		fmt.Fprintf(&b, "guard_scan: status=%s inspected=%s total=%s omitted=%s reason=%s\n", scan.Status, scan.Inspected, scan.Total, scan.Omitted, scan.Reason)
		b.WriteString("full manifests: bench guards\n")
		return b.String(), 0
	}
	out, err := toon.Table("guards", []string{"guard", "boundary", "denies", "wired"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	meta, err := toon.Table("guard_scan", []string{"status", "inspected", "total", "omitted", "reason"}, [][]string{{scan.Status, scan.Inspected, scan.Total, scan.Omitted, scan.Reason}})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out + meta, 0
}
