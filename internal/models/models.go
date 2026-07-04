// Package models ports `bench models`: it lists the Anthropic Models API ids when
// ANTHROPIC_API_KEY is set, else prints the
// no-key guidance. The port drops the curl+python3 dependency for Go net/http; the
// live HTTP call is the untested boundary, so parseIDs and noKeyText hold the logic.
package models

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

const modelsURL = "https://api.anthropic.com/v1/models"

// Command implements `bench models`. With ANTHROPIC_API_KEY set it queries the
// Anthropic Models API and lists the ids, falling back to a failure line on any
// error; without a key it returns the no-key guidance. Exit 0 in all cases.
func Command(args []string) (string, int) {
	key := os.Getenv("ANTHROPIC_API_KEY")
	if key == "" {
		return noKeyText(), 0
	}
	const header = "Available models (Anthropic Models API):\n"
	ids, err := fetchIDs(key)
	if err != nil {
		return header + "  (query failed — check the key, or read your harness model list)\n", 0
	}
	var b strings.Builder
	b.WriteString(header)
	for _, id := range ids {
		b.WriteString("  " + id + "\n")
	}
	return b.String(), 0
}

// fetchIDs performs the live query — the untested boundary — and delegates parsing
// to parseIDs so the JSON shape stays unit-tested.
func fetchIDs(key string) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, modelsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", key)
	req.Header.Set("anthropic-version", "2023-06-01")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("non-2xx response from models API")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseIDs(body)
}

// parseIDs extracts data[].id from the Models API JSON body.
func parseIDs(body []byte) ([]string, error) {
	var payload struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(payload.Data))
	for _, m := range payload.Data {
		ids = append(ids, m.ID)
	}
	return ids, nil
}

// noKeyText returns the guidance shown when ANTHROPIC_API_KEY is unset, matching the
// shell heredoc byte-for-byte including its trailing newline.
func noKeyText() string {
	return "No ANTHROPIC_API_KEY set, so I can't query the model list directly. Discover from\n" +
		"your harness instead, then bind the tiers (cheap / mid / top) in projects/<name>.md:\n" +
		"  - Claude Code: `claude --help`, or the in-app /model picker\n" +
		"  - Codex:       `codex --help`, or its model config\n" +
		"  - or export ANTHROPIC_API_KEY and re-run `bench models`\n"
}
