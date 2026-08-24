package lines

import "strings"

// fullBinding binds codex and claude and leaves opencode unadopted. Its claude cells are
// deliberately NOT this repo's own (`fable`/`opus`/`sonnet`), so a renderer that hard-codes
// the kit's live binding fails every message assertion below.
const fullBinding = "BENCH_CODEX_TOP=gpt-5.6-sol\n" +
	"BENCH_CODEX_MID=gpt-5.6-terra\n" +
	"BENCH_CODEX_CHEAP=gpt-5.6-luna\n" +
	"BENCH_CLAUDE_TOP=fable-5\n" +
	"BENCH_CLAUDE_MID=opus-4-8\n" +
	"BENCH_CLAUDE_CHEAP=sonnet-5\n"

// retiredBinding carries only the retired BENCH_TIER_* / BENCH_ALIAS_* schema. The cut is
// hard: it binds nothing.
const retiredBinding = "BENCH_TIER_TOP=gpt-5.6-sol\n" +
	"BENCH_TIER_MID=gpt-5.6-terra\n" +
	"BENCH_TIER_CHEAP=gpt-5.6-luna\n" +
	"BENCH_ALIAS_TOP=fable\n" +
	"BENCH_ALIAS_MID=opus\n" +
	"BENCH_ALIAS_CHEAP=sonnet\n"

func bound(content string) Source {
	return Source{Path: ".bench/lines.env", Exists: true, Content: []byte(content)}
}

func envelope(model string) []byte {
	return []byte(`{"tool_input":{"model":"` + model + `"}}`)
}

// agentEnvelope is the captured Agent PreToolUse shape: tool_input carries description,
// prompt, subagent_type, and an optional model. subagentType and model are raw JSON, so a
// case can drive a non-string subagent_type or omit either field entirely.
func agentEnvelope(subagentType, model string) []byte {
	fields := `"description":"d","prompt":"p"`
	if subagentType != "" {
		fields += `,"subagent_type":` + subagentType
	}
	if model != "" {
		fields += `,"model":` + model
	}
	return []byte(`{"tool_name":"Agent","tool_input":{` + fields + `}}`)
}

func forkEnvelope(model string) []byte {
	return agentEnvelope(`"fork"`, model)
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
