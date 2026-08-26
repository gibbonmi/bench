// Package guards ports `bench guards` (and `--brief`). It aggregates every deny-capable
// guard's manifest into a `guards[N]{guard,boundary,denies,wired}:` TOON table, so the
// block surface is learnable without a collision. Guards are discovered by convention:
// each .bench/hooks/*.sh and the installed git pre-push hook. Their manifests are
// parsed from static leading-comment headers. Discovery reads scripts only as data and
// never executes them.
//
// The `wired` cell names which harness configs actually reference a hook script, so the
// deny surface the reader sees matches the hooks that can fire here. This cell is
// derived, never declared, by scanning the hook config of every harnesses record row
// that names one for the script's relative path token, the same substring convention
// conformance uses.
// The pre-push hook is wired through git rather than a harness config, so its wired
// cell is the constant `git`, and its install posture (managed/unmanaged/not
// installed) stays in the denies column.
package guards

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/adopt"
	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/harnesses"
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
		return [][]string{append(row, "", "", "", c.wired)}
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

// WiresScript reports whether a harness hook config wires the script at script, its
// repository-relative path token. A command reaches the script through
// $CLAUDE_PROJECT_DIR, a brace form, or a git toplevel expansion, so the relative token
// is the one part every form shares, and a substring test is the rule. An unparseable
// config wires nothing, because this predicate is a read-only wiring reader and the
// JSON-validity conformance family owns malformedness. This function is the one owner of
// the rule: the guards report and the harness-record conformance check both call it.
func WiresScript(config []byte, script string) bool {
	return json.Valid(config) && bytes.Contains(config, []byte(script))
}

// wiredHarnesses names the harness configs that reference scriptRel, its relative
// path token. It scans the hook config of every harnesses record row that names one.
// A row with no hook config
// contributes nothing, so a new harness joins the report as one record row. An absent
// config contributes nothing too; a repo without .codex/ cannot wire Codex. An
// unparseable config scans as not-wired, because guards is a read-only wiring reporter
// and the JSON-validity conformance family owns malformedness. The names sort, so the
// cell stays stable against the record's own row order. The result is one name, the
// comma-joined names, or the definitive "none", never a blank cell.
func wiredHarnesses(root, scriptRel string) string {
	var wired []string
	for _, row := range harnesses.Rows {
		if row.HookConfig == "" {
			continue
		}
		content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(row.HookConfig)))
		if err != nil {
			continue
		}
		if WiresScript(content, scriptRel) {
			wired = append(wired, row.Harness)
		}
	}
	if len(wired) == 0 {
		return "none"
	}
	sort.Strings(wired)
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

// prePushRow renders the hook-health record and keeps the generic manifest parser for
// the managed hook's guard fields.
func prePushRow(root string) [][]string {
	health := adopt.InspectPrePush(root)
	if health.State == adopt.PrePushManaged {
		if row, emit := guardRow(health.Path, "pre-push"); emit {
			return [][]string{append(row, health.Branch, string(health.Provenance), string(health.Currency))}
		}
		return nil
	}
	if health.State == adopt.PrePushAbsent {
		return [][]string{{"pre-push", "", "not installed", "", "", string(health.Currency)}}
	}
	return [][]string{{"pre-push", "", "unmanaged (no manifest)", "", "", string(health.Currency)}}
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
			if r[0] == "pre-push" {
				fmt.Fprintf(&b, "%s: %s [branch: %s provenance: %s currency: %s wired: %s]\n", r[0], r[2], r[3], r[4], r[5], r[6])
				continue
			}
			fmt.Fprintf(&b, "%s: %s [wired: %s]\n", r[0], r[2], r[6])
		}
		fmt.Fprintf(&b, "guard_scan: status=%s inspected=%s total=%s omitted=%s reason=%s\n", scan.Status, scan.Inspected, scan.Total, scan.Omitted, scan.Reason)
		b.WriteString("full manifests: bench guards\n")
		return b.String(), 0
	}
	out, err := toon.Table("guards", []string{"guard", "boundary", "denies", "branch", "provenance", "currency", "wired"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	meta, err := toon.Table("guard_scan", []string{"status", "inspected", "total", "omitted", "reason"}, [][]string{{scan.Status, scan.Inspected, scan.Total, scan.Omitted, scan.Reason}})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	helpActions := []axi.Action(nil)
	if scan.Status == "complete" {
		helpActions = actionsForRows(rows)
	}
	help, err := axi.RenderHelp(helpActions)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out + meta + help, 0
}

func actionsForRows(rows [][]string) []axi.Action {
	actions := make([]axi.Action, 0, len(rows))
	for _, row := range rows {
		if len(row) < 7 {
			continue
		}
		if row[5] == string(adopt.PrePushStale) || row[6] == "none" {
			actions = append(actions, axi.ExecutableInvocation("repair "+row[0], axi.KnownArgument("link")))
		}
	}
	return actions
}
