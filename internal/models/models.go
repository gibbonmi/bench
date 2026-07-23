// Package models ports `bench models`: an advisory inventory of model ids the
// reviewer can use when binding the line in .bench/lines.env.
package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/modelid"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand —
// arity, flag recognition, `--`, and help all come from there rather than a local switch.
var grammar = usage.Grammar{
	Cmd:  "bench models",
	Help: "usage: bench models",
}

const (
	openAIModelsURL    = "https://api.openai.com/v1/models"
	anthropicModelsURL = "https://api.anthropic.com/v1/models"
)

var (
	providerTimeout = bounds.ProviderTimeout
	runCommand      = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		result := bounds.Run(ctx, providerTimeout, exec.Command(name, args...))
		if result.Status != bounds.ProcessComplete {
			return result.Output, providerError{status: string(result.Status), err: result.Err}
		}
		return result.Output, nil
	}
	modelsHTTPClient = &http.Client{}
	doHTTP           = modelsHTTPClient.Do
)

type providerError struct {
	status string
	err    error
}

func (e providerError) Error() string {
	if e.err == nil {
		return e.status
	}
	return e.status + ": " + e.err.Error()
}

type sourceRow struct {
	source, freshness, status, hint string
}

type modelRow struct {
	source, freshness, id string
}

// Command implements `bench models`. Discovery is advisory: every per-source
// failure becomes an unavailable row, and the command exits 0 unless rendering
// its own structured output fails.
func Command(args []string) (string, int) {
	// bench models takes no arguments; any is a misuse the grammar rejects with a usage
	// line at exit 2. That is distinct from the discovery tolerance below (unreachable
	// providers → unavailable rows at exit 0), which the no-arg path keeps unchanged.
	if _, line, code := usage.Parse(grammar, args); line != "" {
		return line + "\n", code
	}
	sources, models := inventory()
	out, err := render(sources, models)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out, 0
}

func inventory() ([]sourceRow, []modelRow) {
	if bounds.Offline() {
		sources := make([]sourceRow, 0, 3)
		for _, source := range []string{"codex", "openai", "anthropic"} {
			sources = append(sources, sourceRow{source: source, freshness: "offline", status: "offline", hint: "BENCH_OFFLINE=1"})
		}
		return sources, nil
	}
	type result struct {
		source sourceRow
		models []modelRow
	}
	results := make([]result, 3)
	jobs := []func() (sourceRow, []modelRow){
		codexInventory,
		func() (sourceRow, []modelRow) {
			return apiInventory("openai", openAIModelsURL, os.Getenv("OPENAI_API_KEY"), openAIHeaders, "set OPENAI_API_KEY")
		},
		func() (sourceRow, []modelRow) {
			return apiInventory("anthropic", anthropicModelsURL, os.Getenv("ANTHROPIC_API_KEY"), anthropicHeaders, "set ANTHROPIC_API_KEY")
		},
	}
	var wg sync.WaitGroup
	for i := range jobs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i].source, results[i].models = jobs[i]()
		}()
	}
	wg.Wait()
	sources := make([]sourceRow, 0, 3)
	var models []modelRow
	for _, result := range results {
		sources = append(sources, result.source)
		models = append(models, result.models...)
	}
	return sources, models
}

func codexInventory() (sourceRow, []modelRow) {
	ctx, cancel := bounds.Context(context.Background(), providerTimeout)
	defer cancel()
	body, err := runCommand(ctx, "codex", "debug", "models")
	freshness := "live"
	if err != nil {
		if ctx.Err() != nil || errorStatus(err) == "timeout" {
			return sourceFailure("codex", "timeout", fmt.Sprintf("%s provider deadline", bounds.ProviderTimeout)), nil
		}
		body, err = runCommand(ctx, "codex", "debug", "models", "--bundled")
		freshness = "bundled"
	}
	if err != nil {
		if ctx.Err() != nil || errorStatus(err) == "timeout" {
			return sourceFailure("codex", "timeout", fmt.Sprintf("%s provider deadline", bounds.ProviderTimeout)), nil
		}
		return unavailable("codex", "codex debug models unavailable"), nil
	}
	ids, err := parseCodexSlugs(body)
	if err != nil {
		return unavailable("codex", "invalid codex model catalog"), nil
	}
	return available("codex", freshness), modelRows("codex", freshness, ids)
}

func apiInventory(source, url, key string, headers func(*http.Request, string), missingHint string) (sourceRow, []modelRow) {
	if key == "" {
		return unavailable(source, missingHint), nil
	}
	ctx, cancel := bounds.Context(context.Background(), providerTimeout)
	defer cancel()
	ids, err := fetchDataIDsContext(ctx, url, key, headers)
	if err != nil {
		switch errorStatus(err) {
		case "timeout":
			return sourceFailure(source, "timeout", fmt.Sprintf("%s provider deadline", bounds.ProviderTimeout)), nil
		case "oversized":
			return sourceFailure(source, "oversized", fmt.Sprintf("%d MiB response limit", bounds.ModelReadLimit>>20)), nil
		}
		return unavailable(source, "query failed"), nil
	}
	return available(source, "live"), modelRows(source, "live", ids)
}

func available(source, freshness string) sourceRow {
	return sourceRow{source: source, freshness: freshness, status: "available", hint: "none"}
}

func unavailable(source, hint string) sourceRow {
	return sourceRow{source: source, freshness: "unavailable", status: "unavailable", hint: hint}
}

func sourceFailure(source, status, hint string) sourceRow {
	return sourceRow{source: source, freshness: status, status: status, hint: hint}
}

func errorStatus(err error) string {
	var p providerError
	if errors.As(err, &p) {
		return p.status
	}
	return ""
}

func modelRows(source, freshness string, ids []string) []modelRow {
	rows := make([]modelRow, 0, len(ids))
	for _, id := range ids {
		if !modelid.SafeToken(id) {
			continue
		}
		rows = append(rows, modelRow{source: source, freshness: freshness, id: id})
	}
	return rows
}

func render(sources []sourceRow, models []modelRow) (string, error) {
	sourceRows := make([][]string, 0, len(sources))
	for _, row := range sources {
		sourceRows = append(sourceRows, []string{row.source, row.freshness, row.status, row.hint})
	}
	modelRows := make([][]string, 0, len(models))
	for _, row := range models {
		modelRows = append(modelRows, []string{row.source, row.freshness, row.id})
	}
	sourceTable, err := toon.Table("model_sources", []string{"source", "freshness", "status", "hint"}, sourceRows)
	if err != nil {
		return "", err
	}
	modelTable, err := toon.Table("models", []string{"source", "freshness", "id"}, modelRows)
	if err != nil {
		return "", err
	}
	return sourceTable + modelTable, nil
}

func fetchDataIDs(url, key string, headers func(*http.Request, string)) ([]string, error) {
	ctx, cancel := bounds.Context(context.Background(), providerTimeout)
	defer cancel()
	return fetchDataIDsContext(ctx, url, key, headers)
}

func fetchDataIDsContext(ctx context.Context, url, key string, headers func(*http.Request, string)) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	headers(req, key)
	resp, err := doHTTP(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, providerError{status: "timeout", err: ctx.Err()}
		}
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("non-2xx response from models API")
	}
	read := bounds.Read(resp.Body, bounds.ModelReadLimit)
	if read.Status == bounds.ReadOversized {
		return nil, providerError{status: "oversized", err: read.Err}
	}
	if read.Status == bounds.ReadFailed {
		return nil, read.Err
	}
	return parseDataIDs(read.Data)
}

func openAIHeaders(req *http.Request, key string) {
	req.Header.Set("Authorization", "Bearer "+key)
}

func anthropicHeaders(req *http.Request, key string) {
	req.Header.Set("x-api-key", key)
	req.Header.Set("anthropic-version", "2023-06-01")
}

func parseCodexSlugs(body []byte) ([]string, error) {
	var payload struct {
		Models []struct {
			Slug       string `json:"slug"`
			Visibility string `json:"visibility"`
		} `json:"models"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(payload.Models))
	for _, m := range payload.Models {
		if m.Visibility == "list" {
			ids = append(ids, m.Slug)
		}
	}
	return ids, nil
}

func parseDataIDs(body []byte) ([]string, error) {
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
