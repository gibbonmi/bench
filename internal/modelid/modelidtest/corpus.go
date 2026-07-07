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
