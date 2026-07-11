// Package models ports `bench models`: an advisory inventory of model ids the
// reviewer can use when binding the line in .bench/lines.env.
package models

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gibbonmi/bench/internal/modelid"
	"github.com/gibbonmi/bench/internal/toon"
)

const (
	openAIModelsURL    = "https://api.openai.com/v1/models"
	anthropicModelsURL = "https://api.anthropic.com/v1/models"
	// modelsQueryTimeout bounds provider discovery without treating ordinarily slow
	// model-list responses as unavailable.
	modelsQueryTimeout = 10 * time.Second
)

var (
	runCommand = func(name string, args ...string) ([]byte, error) {
		return exec.Command(name, args...).Output()
	}
	modelsHTTPClient = &http.Client{Timeout: modelsQueryTimeout}
	doHTTP           = modelsHTTPClient.Do
)

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
	if len(args) > 0 {
		// bench models takes no arguments; any is a misuse, rejected with a usage line at
		// exit 2 to match every sibling porcelain's default case. This is distinct from the
		// discovery tolerance below (unreachable providers → unavailable rows at exit 0),
		// which the no-arg path keeps unchanged.
		return toon.Usage("bench models", args[0]) + "\n", 2
	}
	sources, models := inventory()
	out, err := render(sources, models)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out, 0
}

func inventory() ([]sourceRow, []modelRow) {
	var sources []sourceRow
	var models []modelRow

	codexSource, codexModels := codexInventory()
	sources = append(sources, codexSource)
	models = append(models, codexModels...)

	openAISource, openAIModels := apiInventory("openai", openAIModelsURL, os.Getenv("OPENAI_API_KEY"), openAIHeaders, "set OPENAI_API_KEY")
	sources = append(sources, openAISource)
	models = append(models, openAIModels...)

	anthropicSource, anthropicModels := apiInventory("anthropic", anthropicModelsURL, os.Getenv("ANTHROPIC_API_KEY"), anthropicHeaders, "set ANTHROPIC_API_KEY")
	sources = append(sources, anthropicSource)
	models = append(models, anthropicModels...)

	return sources, models
}

func codexInventory() (sourceRow, []modelRow) {
	body, err := runCommand("codex", "debug", "models")
	freshness := "live"
	if err != nil {
		body, err = runCommand("codex", "debug", "models", "--bundled")
		freshness = "bundled"
	}
	if err != nil {
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
	ids, err := fetchDataIDs(url, key, headers)
	if err != nil {
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
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	headers(req, key)
	resp, err := doHTTP(req)
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
	return parseDataIDs(body)
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
