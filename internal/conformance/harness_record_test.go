package conformance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/harnesses"
)

// harnessHooksDir is the one directory a shipped hook config may name. A command that
// names no path under it wires something outside Bench, and the record does not grade it.
const harnessHooksDir = ".bench/hooks/"

// harnessAdaptersDir holds one headless adapter per harness. The check enumerates it, so a
// new adapter cannot ship without a row.
const harnessAdaptersDir = ".bench/adapters"

// harnessGuardScript is the delegation guard. A config that names it holds a deny-capable
// verdict on a spawn, and the record's delegation_guard cell must agree.
const harnessGuardScript = harnessHooksDir + "check-agent-line.sh"

// harnessConfigNames are the two file names a harness uses for its hook config. A config
// under either name that wires a Bench hook script must map to a record row.
var harnessConfigNames = []string{"settings.json", "hooks.json"}

// harnessHookConfig is the shape both shipped configs share: an event names a list of
// groups, and a group carries a matcher plus the commands it runs. The check reads only
// the fields it grades, so a config field Bench does not own cannot red the check.
type harnessHookConfig struct {
	Hooks map[string][]struct {
		Matcher string `json:"matcher"`
		Hooks   []struct {
			Command string `json:"command"`
		} `json:"hooks"`
	} `json:"hooks"`
}

// checkHarnessRecord grades internal/harnesses against the tree in both directions. Disk
// supplies the adapters and the hook configs, and each one must map to a row. Each row
// then supplies its declared events, its adapter path, and its delegation verdict, and
// each one must match what the tree wires. A record that only walks itself stays green
// while a new harness ships ungraded, so the enumeration half is what keeps the record
// honest.
func checkHarnessRecord(root string) []string {
	var diags []string
	diags = append(diags, harnessRecordAdapterDiags(root)...)
	diags = append(diags, harnessRecordConfigDiags(root)...)
	for _, row := range harnesses.Rows {
		diags = append(diags, harnessRecordRowDiags(root, row)...)
	}
	return diags
}

// harnessRecordAdapterDiags names each shipped adapter no row claims. The adapter's base
// name is the harness name, because the record's Headless path is built that way.
func harnessRecordAdapterDiags(root string) []string {
	entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(harnessAdaptersDir)))
	if err != nil {
		return nil
	}
	var diags []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if _, found := harnesses.Lookup(entry.Name()); found {
			continue
		}
		diags = append(diags, fmt.Sprintf("harness-record: %s/%s maps to no harness row; give the record a %q row", harnessAdaptersDir, entry.Name(), entry.Name()))
	}
	sort.Strings(diags)
	return diags
}

// harnessRecordConfigDiags names each shipped hook config no row claims. A dot directory
// at the root holds a config when it carries settings.json or hooks.json that names a
// .bench/hooks/ script. A dot directory that wires nothing of Bench's is not a harness
// config, so the check leaves it alone.
func harnessRecordConfigDiags(root string) []string {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	claimed := map[string]bool{}
	for _, row := range harnesses.Rows {
		if row.HookConfig != "" {
			claimed[row.HookConfig] = true
		}
	}
	var diags []string
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		for _, name := range harnessConfigNames {
			rel := entry.Name() + "/" + name
			classified := bounds.ClassifyNoFollow(filepath.Join(root, entry.Name(), name))
			if classified.State != bounds.StateParsed || !strings.Contains(string(classified.Data), harnessHooksDir) {
				continue
			}
			if claimed[rel] {
				continue
			}
			diags = append(diags, fmt.Sprintf("harness-record: %s wires a %s script and maps to no harness row", rel, harnessHooksDir))
		}
	}
	sort.Strings(diags)
	return diags
}

// harnessRecordRowDiags grades one row against the tree. A refused config path returns at
// once with the refusal alone: the bytes behind a FIFO or a link were never read, so every
// further verdict about that row would rest on nothing.
func harnessRecordRowDiags(root string, row harnesses.Row) []string {
	diags, refused := harnessRecordConfigRowDiags(root, row)
	if refused {
		return diags
	}
	return append(diags, harnessRecordHeadlessDiags(root, row)...)
}

// harnessRecordConfigRowDiags grades a row's hook config. The second result reports a
// refused path. An absent config and an empty one are two verdicts: the first says the row
// names a file the tree does not ship, and the second says the file ships and wires
// nothing.
func harnessRecordConfigRowDiags(root string, row harnesses.Row) ([]string, bool) {
	if row.HookConfig == "" {
		return nil, false
	}
	classified := bounds.ClassifyNoFollow(filepath.Join(root, filepath.FromSlash(row.HookConfig)))
	switch classified.State {
	case bounds.StateAbsent:
		return []string{fmt.Sprintf("harness-record: %s declares hook config %s, and the tree ships no such file", row.Harness, row.HookConfig)}, false
	case bounds.StateEmpty:
		return []string{fmt.Sprintf("harness-record: %s declares hook config %s, and the file is empty", row.Harness, row.HookConfig)}, false
	case bounds.StateParsed:
	default:
		return []string{fmt.Sprintf("harness-record: %s is refused by the no-follow classifier: %s", row.HookConfig, classified.Reason)}, true
	}
	var config harnessHookConfig
	if err := json.Unmarshal(classified.Data, &config); err != nil {
		return []string{fmt.Sprintf("harness-record: %s is not valid JSON: %s", row.HookConfig, err.Error())}, false
	}
	diags := harnessRecordEventDiags(row, config)
	return append(diags, harnessRecordGuardDiags(row, classified.Data)...), false
}

// harnessRecordEventDiags grades the declared events against the wired groups in both
// directions. An overclaim names the event the config does not wire. An underclaim names
// each script the unrecorded group runs, because a script is what a reader must go and
// read to decide whether the row should carry the group.
func harnessRecordEventDiags(row harnesses.Row, config harnessHookConfig) []string {
	wired := harnessRecordWiredGroups(config)
	declared := map[string]bool{}
	for _, event := range row.HookEvents {
		declared[event] = true
	}
	var diags []string
	for _, event := range row.HookEvents {
		if _, found := wired[event]; !found {
			diags = append(diags, fmt.Sprintf("harness-record: %s declares hook event %s, and %s does not wire it", row.Harness, event, row.HookConfig))
		}
	}
	for _, event := range harnessRecordSortedKeys(wired) {
		if declared[event] {
			continue
		}
		for _, script := range wired[event] {
			diags = append(diags, fmt.Sprintf("harness-record: %s omits %s, which %s wires at %s", row.Harness, script, row.HookConfig, event))
		}
	}
	return diags
}

// harnessRecordWiredGroups maps each wired group to the .bench/hooks/ scripts it runs. A
// group key is the event, qualified by its matcher when the matcher selects a tool. The
// record uses the same convention, so a config that wires one event under two matchers
// grades as two groups.
func harnessRecordWiredGroups(config harnessHookConfig) map[string][]string {
	wired := map[string][]string{}
	for event, groups := range config.Hooks {
		for _, group := range groups {
			key := event
			if group.Matcher != "" && group.Matcher != "*" {
				key = event + ":" + group.Matcher
			}
			for _, hook := range group.Hooks {
				if script := harnessRecordScript(hook.Command); script != "" {
					wired[key] = append(wired[key], script)
				}
			}
			if _, found := wired[key]; !found {
				wired[key] = nil
			}
		}
	}
	return wired
}

// harnessRecordScript extracts the repository-relative script path a command names. A
// command reaches the script through $CLAUDE_PROJECT_DIR, ${CLAUDE_PROJECT_DIR}, or a git
// toplevel expansion, so the relative token is the one part every form shares. That is the
// same acceptance rule the guards wiring reader applies.
func harnessRecordScript(command string) string {
	index := strings.Index(command, harnessHooksDir)
	if index < 0 {
		return ""
	}
	rest := command[index+len(harnessHooksDir):]
	end := strings.IndexFunc(rest, func(r rune) bool {
		return !(r == '.' || r == '_' || r == '-' || r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z')
	})
	if end >= 0 {
		rest = rest[:end]
	}
	if rest == "" {
		return ""
	}
	return harnessHooksDir + rest
}

func harnessRecordSortedKeys(wired map[string][]string) []string {
	keys := make([]string, 0, len(wired))
	for key := range wired {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// harnessRecordGuardDiags grades the delegation_guard cell against the guard wiring. An
// unknown cell records no fact, so the tree cannot contradict it.
func harnessRecordGuardDiags(row harnesses.Row, config []byte) []string {
	wired := strings.Contains(string(config), harnessGuardScript)
	switch {
	case row.DelegationGuard.Value == harnesses.Yes && !wired:
		return []string{fmt.Sprintf("harness-record: %s records delegation_guard yes, and %s wires no %s", row.Harness, row.HookConfig, harnessGuardScript)}
	case row.DelegationGuard.Value == harnesses.No && wired:
		return []string{fmt.Sprintf("harness-record: %s records delegation_guard no, and %s wires %s", row.Harness, row.HookConfig, harnessGuardScript)}
	}
	return nil
}

// harnessRecordHeadlessDiags grades the adapter path against disk. The resolving stat is
// deliberate: a dangling symlink names an adapter that cannot run, so it counts as absent
// rather than as a refusal.
func harnessRecordHeadlessDiags(root string, row harnesses.Row) []string {
	if row.Headless == "" {
		return nil
	}
	if exists(filepath.Join(root, filepath.FromSlash(row.Headless))) {
		return nil
	}
	return []string{fmt.Sprintf("harness-record: %s records headless adapter %s, and the tree ships no such entry", row.Harness, row.Headless)}
}

// harnessRecordConfig renders a hook config from its wired groups. A test states the
// groups it means to ship, so no case restates the JSON envelope.
func harnessRecordConfig(groups ...string) string {
	var order []string
	rendered := map[string][]string{}
	for _, group := range groups {
		event, matcher, scripts := harnessRecordGroupParts(group)
		if _, seen := rendered[event]; !seen {
			order = append(order, event)
		}
		var commands []string
		for _, script := range scripts {
			commands = append(commands, fmt.Sprintf(`{"type":"command","command":"$CLAUDE_PROJECT_DIR/%s%s"}`, harnessHooksDir, script))
		}
		rendered[event] = append(rendered[event], fmt.Sprintf(`{"matcher":%q,"hooks":[%s]}`, matcher, strings.Join(commands, ",")))
	}
	var events []string
	for _, event := range order {
		events = append(events, fmt.Sprintf(`%q:[%s]`, event, strings.Join(rendered[event], ",")))
	}
	return `{"hooks":{` + strings.Join(events, ",") + `}}`
}

// harnessRecordGroupParts reads one test group spec: an event, an optional :matcher, and a
// space-separated script list after an equals sign.
func harnessRecordGroupParts(group string) (string, string, []string) {
	key, list, _ := strings.Cut(group, "=")
	event, matcher, _ := strings.Cut(key, ":")
	return event, matcher, strings.Fields(list)
}

// harnessRecordRoot builds a throwaway tree that satisfies every row, then applies the one
// planted fault a case grades. Each case therefore reds on its own fault alone.
func harnessRecordRoot(t *testing.T, fault harnessRecordFault) string {
	t.Helper()
	files := map[string]string{
		".bench/adapters/claude":   "#!/usr/bin/env bash\nexit 0\n",
		".bench/adapters/codex":    "#!/usr/bin/env bash\nexit 0\n",
		".bench/adapters/opencode": "#!/usr/bin/env bash\nexit 0\n",
		".claude/settings.json": harnessRecordConfig(
			"WorktreeCreate=worktree-lifecycle.sh",
			"WorktreeRemove=worktree-lifecycle.sh",
			"SessionStart=session-start.sh",
			"Stop:*=stop.sh",
			"PreToolUse:Bash=block-dangerous-git.sh block-bench-follow-on.sh",
			"PreToolUse:Agent=check-agent-line.sh",
		),
		".codex/hooks.json": harnessRecordConfig(
			"SessionStart=session-start.sh",
			"Stop=stop.sh",
			"PreToolUse:Bash=block-dangerous-git.sh block-bench-follow-on.sh",
		),
	}
	delete(files, fault.drop)
	for rel, content := range fault.write {
		files[rel] = content
	}
	return throwawayRoot{files: files, plants: fault.plant}.build(t)
}

// harnessRecordFault is the one planted fault a case grades: an entry the honest tree
// ships and this root drops, entries it overwrites or adds, and entries a hostile planter
// writes rather than the builder.
type harnessRecordFault struct {
	drop  string
	write map[string]string
	plant map[string]func(*testing.T, string)
}

// TestHarnessRecordIsGreenOverTheHonestTree proves every diagnostic below reds on its own
// planted fault, not on the fixture's shape.
func TestHarnessRecordIsGreenOverTheHonestTree(t *testing.T) {
	if diags := checkHarnessRecord(harnessRecordRoot(t, harnessRecordFault{})); len(diags) != 0 {
		t.Fatalf("the honest tree is not green:\n%s", strings.Join(diags, "\n"))
	}
}

// TestHarnessRecordBites plants one fault per graded direction. Each case names the exact
// subject a reader must go and fix.
func TestHarnessRecordBites(t *testing.T) {
	tests := []struct {
		name  string
		fault harnessRecordFault
		want  string
	}{
		{
			name:  "an adapter outside the record",
			fault: harnessRecordFault{write: map[string]string{".bench/adapters/cursor": "#!/usr/bin/env bash\nexit 0\n"}},
			want:  "harness-record: .bench/adapters/cursor maps to no harness row",
		},
		{
			name:  "a hook config outside the record",
			fault: harnessRecordFault{write: map[string]string{".cursor/hooks.json": harnessRecordConfig("Stop=stop.sh")}},
			want:  "harness-record: .cursor/hooks.json wires a .bench/hooks/ script and maps to no harness row",
		},
		{
			name: "a declared event the config does not wire",
			fault: harnessRecordFault{write: map[string]string{".claude/settings.json": harnessRecordConfig(
				"WorktreeCreate=worktree-lifecycle.sh",
				"WorktreeRemove=worktree-lifecycle.sh",
				"SessionStart=session-start.sh",
				"PreToolUse:Bash=block-dangerous-git.sh block-bench-follow-on.sh",
				"PreToolUse:Agent=check-agent-line.sh",
			)}},
			want: "harness-record: claude declares hook event Stop, and .claude/settings.json does not wire it",
		},
		{
			name: "a wired script the row omits",
			fault: harnessRecordFault{write: map[string]string{".codex/hooks.json": harnessRecordConfig(
				"SessionStart=session-start.sh",
				"Stop=stop.sh",
				"PreToolUse:Bash=block-dangerous-git.sh block-bench-follow-on.sh",
				"PreToolUse:Agent=check-agent-line.sh",
			)}},
			want: "harness-record: codex omits .bench/hooks/check-agent-line.sh, which .codex/hooks.json wires at PreToolUse:Agent",
		},
		{
			name:  "an absent headless adapter",
			fault: harnessRecordFault{drop: ".bench/adapters/codex"},
			want:  "harness-record: codex records headless adapter .bench/adapters/codex, and the tree ships no such entry",
		},
		{
			name: "a dangling headless adapter",
			fault: harnessRecordFault{
				drop:  ".bench/adapters/codex",
				plant: map[string]func(*testing.T, string){".bench/adapters/codex": hostileSkillPlanters["dangling symlink"]},
			},
			want: "harness-record: codex records headless adapter .bench/adapters/codex, and the tree ships no such entry",
		},
		{
			name: "a delegation guard the config contradicts",
			fault: harnessRecordFault{write: map[string]string{".claude/settings.json": harnessRecordConfig(
				"WorktreeCreate=worktree-lifecycle.sh",
				"WorktreeRemove=worktree-lifecycle.sh",
				"SessionStart=session-start.sh",
				"Stop:*=stop.sh",
				"PreToolUse:Bash=block-dangerous-git.sh block-bench-follow-on.sh",
				"PreToolUse:Agent=session-start.sh",
			)}},
			want: "harness-record: claude records delegation_guard yes, and .claude/settings.json wires no .bench/hooks/check-agent-line.sh",
		},
		{
			name:  "an absent hook config",
			fault: harnessRecordFault{drop: ".codex/hooks.json"},
			want:  "harness-record: codex declares hook config .codex/hooks.json, and the tree ships no such file",
		},
		{
			name:  "an empty hook config",
			fault: harnessRecordFault{write: map[string]string{".codex/hooks.json": ""}},
			want:  "harness-record: codex declares hook config .codex/hooks.json, and the file is empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := checkHarnessRecord(harnessRecordRoot(t, tt.fault))
			if !containsDiagnostic(diags, tt.want) {
				t.Fatalf("the planted fault did not bite with %q:\n%s", tt.want, strings.Join(diags, "\n"))
			}
		})
	}
}

// TestHarnessRecordAbsentAndEmptyConfigsDiffer proves the two config verdicts are distinct
// text. A single shared wording would send a reader to write a file that already exists.
func TestHarnessRecordAbsentAndEmptyConfigsDiffer(t *testing.T) {
	absent := harnessRecordCodexConfigDiags(t, checkHarnessRecord(harnessRecordRoot(t, harnessRecordFault{drop: ".codex/hooks.json"})))
	empty := harnessRecordCodexConfigDiags(t, checkHarnessRecord(harnessRecordRoot(t, harnessRecordFault{write: map[string]string{".codex/hooks.json": ""}})))
	if absent == empty {
		t.Fatalf("the absent and the empty config share one diagnostic: %q", absent)
	}
}

func harnessRecordCodexConfigDiags(t *testing.T, diags []string) string {
	t.Helper()
	for _, diag := range diags {
		if strings.Contains(diag, "codex declares hook config") {
			return diag
		}
	}
	t.Fatalf("no codex config diagnostic:\n%s", strings.Join(diags, "\n"))
	return ""
}

// TestHarnessRecordRefusesAFIFOConfig proves the classifier runs before the read. The
// refusal is the row's only verdict, because a check that also graded the events would be
// reporting on bytes it never read.
func TestHarnessRecordRefusesAFIFOConfig(t *testing.T) {
	root := harnessRecordRoot(t, harnessRecordFault{
		drop:  ".codex/hooks.json",
		plant: map[string]func(*testing.T, string){".codex/hooks.json": hostileSkillPlanters["fifo"]},
	})

	diags := checkHarnessRecord(root)

	if !containsDiagnostic(diags, "harness-record: .codex/hooks.json is refused by the no-follow classifier") {
		t.Fatalf("the FIFO config was not refused:\n%s", strings.Join(diags, "\n"))
	}
	for _, diag := range diags {
		if strings.Contains(diag, "codex") && !strings.Contains(diag, "refused by the no-follow classifier") {
			t.Fatalf("the refused row carries a second diagnostic: %q", diag)
		}
	}
}
