// Package modelidtest holds the shared safe-token corpus used by tests at
// different seams. It is intentionally separate from production validation.
package modelidtest

type RejectedToken struct {
	Name  string
	Value string
}

var AcceptedTokens = []string{
	"9_model",
	"gpt-5.4",
	"gpt-5.3-codex-spark",
	"openai/gpt-5",
	"openai/gpt_5",
	"anthropic:claude@2026+beta",
}

var RejectedTokens = []RejectedToken{
	{Name: "empty", Value: ""},
	{Name: "whitespace", Value: "gpt 5"},
	{Name: "control", Value: "gpt\x1b5"},
	{Name: "dollar", Value: "gpt$5"},
	{Name: "backtick", Value: "gpt`5"},
	{Name: "semicolon", Value: "gpt;5"},
	{Name: "pipe", Value: "gpt|5"},
	{Name: "ampersand", Value: "gpt&5"},
	{Name: "glob", Value: "gpt*5"},
	{Name: "quote", Value: "gpt\"5"},
	{Name: "tilde", Value: "gpt~5"},
	{Name: "leading-dash", Value: "-gpt-5"},
}

// NewlineRejectedTokens holds the newline-class rejects: a trailing "\n", an embedded
// "\n", and a CR. These rejects pin the grammar's `$` anchor as end-of-text, Go's
// default `\z` semantics, not multiline.
//
// NewlineRejectedTokens lives in its own slice. Only the direct SafeToken seam
// consumes it. It must not join RejectedTokens.
//
// RejectedTokens has a second consumer, the conformance line-binding check. This check
// writes each value into a temp lines.env and parses it through lines.TierValue.
// TierValue truncates at the first newline and strips a trailing CR. So a
// newline-bearing value would reach SafeToken already sanitized. The "is not a safe
// model token" diagnostic would then never fire. The check would then false-green or
// false-red.
var NewlineRejectedTokens = []RejectedToken{
	{Name: "trailing-newline", Value: "gpt\n"},
	{Name: "embedded-newline", Value: "gpt\n5"},
	{Name: "carriage-return", Value: "gpt\r5"},
}
