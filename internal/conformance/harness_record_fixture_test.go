package conformance

import (
	"fmt"
	"strings"
	"testing"
)

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
			"PreToolUse:Edit|Write|MultiEdit|NotebookEdit=block-primary-file-write.sh",
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
