package models

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestAPIInventoryTimesOutAsUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	t.Cleanup(server.Close)

	started := time.Now()
	row, models := apiInventory("openai", server.URL, "key", openAIHeaders, "set key")
	if elapsed := time.Since(started); elapsed > 12*time.Second {
		t.Fatalf("hung provider returned after %s, want at most 12s", elapsed)
	}
	if row.status != "unavailable" || row.hint != "query failed" {
		t.Fatalf("timeout row = %+v, want unavailable query failed", row)
	}
	if len(models) != 0 {
		t.Fatalf("timeout models = %+v, want none", models)
	}
}

func TestCommandEmitsCodexOpenAIAndAnthropicRows(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "openai-key")
	t.Setenv("ANTHROPIC_API_KEY", "anthropic-key")
	stubCommand(t, func(name string, args ...string) ([]byte, error) {
		if name != "codex" || strings.Join(args, " ") != "debug models" {
			t.Fatalf("unexpected command: %s %v", name, args)
		}
		return []byte(`{"models":[{"slug":"gpt-5.5","visibility":"list"},{"slug":"hidden","visibility":"hidden"},{"slug":"bad value","visibility":"list"}]}`), nil
	})
	stubHTTP(t, func(req *http.Request) (*http.Response, error) {
		switch req.URL.String() {
		case openAIModelsURL:
			if got := req.Header.Get("Authorization"); got != "Bearer openai-key" {
				t.Fatalf("OpenAI Authorization = %q", got)
			}
			return response(200, `{"data":[{"id":"gpt-5.4"},{"id":"bad value"}]}`), nil
		case anthropicModelsURL:
			if got := req.Header.Get("x-api-key"); got != "anthropic-key" {
				t.Fatalf("Anthropic x-api-key = %q", got)
			}
			if got := req.Header.Get("anthropic-version"); got != "2023-06-01" {
				t.Fatalf("Anthropic version = %q", got)
			}
			return response(200, `{"data":[{"id":"claude-opus-4-8"}]}`), nil
		default:
			t.Fatalf("unexpected URL: %s", req.URL.String())
			return nil, nil
		}
	})

	out, code := Command(nil)
	if code != 0 {
		t.Fatalf("code = %d, want 0; out=%s", code, out)
	}
	want := "model_sources[3]{source,freshness,status,hint}:\n" +
		"  codex,live,available,none\n" +
		"  openai,live,available,none\n" +
		"  anthropic,live,available,none\n" +
		"models[3]{source,freshness,id}:\n" +
		"  codex,live,gpt-5.5\n" +
		"  openai,live,gpt-5.4\n" +
		"  anthropic,live,claude-opus-4-8\n"
	if out != want {
		t.Fatalf("Command output =\n%q\nwant\n%q", out, want)
	}
}

func TestCommandFallsBackToCodexBundledAndReportsMissingKeys(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "")
	var calls []string
	stubCommand(t, func(name string, args ...string) ([]byte, error) {
		calls = append(calls, name+" "+strings.Join(args, " "))
		if strings.Join(args, " ") == "debug models" {
			return nil, errors.New("refresh failed")
		}
		if strings.Join(args, " ") == "debug models --bundled" {
			return []byte(`{"models":[{"slug":"gpt-5.4-mini","visibility":"list"}]}`), nil
		}
		t.Fatalf("unexpected command: %s %v", name, args)
		return nil, nil
	})
	stubHTTP(t, func(req *http.Request) (*http.Response, error) {
		t.Fatalf("HTTP should not be called without API keys")
		return nil, nil
	})

	out, code := Command(nil)
	if code != 0 {
		t.Fatalf("code = %d, want 0; out=%s", code, out)
	}
	if got, want := strings.Join(calls, "|"), "codex debug models|codex debug models --bundled"; got != want {
		t.Fatalf("commands = %q, want %q", got, want)
	}
	want := "model_sources[3]{source,freshness,status,hint}:\n" +
		"  codex,bundled,available,none\n" +
		"  openai,unavailable,unavailable,set OPENAI_API_KEY\n" +
		"  anthropic,unavailable,unavailable,set ANTHROPIC_API_KEY\n" +
		"models[1]{source,freshness,id}:\n" +
		"  codex,bundled,gpt-5.4-mini\n"
	if out != want {
		t.Fatalf("Command output =\n%q\nwant\n%q", out, want)
	}
}

func TestCommandReportsSourceFailuresAsUnavailableAtExitZero(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "openai-key")
	t.Setenv("ANTHROPIC_API_KEY", "anthropic-key")
	stubCommand(t, func(name string, args ...string) ([]byte, error) {
		return nil, errors.New("codex failed")
	})
	stubHTTP(t, func(req *http.Request) (*http.Response, error) {
		switch req.URL.String() {
		case openAIModelsURL:
			return response(500, `nope`), nil
		case anthropicModelsURL:
			return response(200, `{"data":[{"id":`), nil
		default:
			t.Fatalf("unexpected URL: %s", req.URL.String())
			return nil, nil
		}
	})

	out, code := Command(nil)
	if code != 0 {
		t.Fatalf("code = %d, want 0; out=%s", code, out)
	}
	for _, want := range []string{
		`codex,unavailable,unavailable,codex debug models unavailable`,
		`openai,unavailable,unavailable,query failed`,
		`anthropic,unavailable,unavailable,query failed`,
		`models[0]{source,freshness,id}:`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, out)
		}
	}
}

func TestCommandKeepsAvailableRowsForEmptyDiscoveryResults(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "openai-key")
	t.Setenv("ANTHROPIC_API_KEY", "anthropic-key")
	stubCommand(t, func(name string, args ...string) ([]byte, error) {
		return []byte(`{"models":[]}`), nil
	})
	stubHTTP(t, func(req *http.Request) (*http.Response, error) {
		switch req.URL.String() {
		case openAIModelsURL, anthropicModelsURL:
			return response(200, `{"data":[]}`), nil
		default:
			t.Fatalf("unexpected URL: %s", req.URL.String())
			return nil, nil
		}
	})

	out, code := Command(nil)
	if code != 0 {
		t.Fatalf("code = %d, want 0; out=%s", code, out)
	}
	want := "model_sources[3]{source,freshness,status,hint}:\n" +
		"  codex,live,available,none\n" +
		"  openai,live,available,none\n" +
		"  anthropic,live,available,none\n" +
		"models[0]{source,freshness,id}:\n"
	if out != want {
		t.Fatalf("Command output =\n%q\nwant\n%q", out, want)
	}
}

func TestCommandSuppressesUnsafeProviderIDs(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "")
	stubCommand(t, func(name string, args ...string) ([]byte, error) {
		return []byte("{\"models\":[" +
			`{"slug":"gpt-5.4","visibility":"list"},` +
			`{"slug":"gpt 5","visibility":"list"},` +
			`{"slug":"gpt\u001b5","visibility":"list"}` +
			"]}"), nil
	})

	out, code := Command(nil)
	if code != 0 {
		t.Fatalf("code = %d, want 0; out=%s", code, out)
	}
	if !strings.Contains(out, "  codex,live,gpt-5.4\n") {
		t.Fatalf("safe id missing from output:\n%s", out)
	}
	if strings.Contains(out, "gpt 5") || strings.Contains(out, "gpt\x1b5") {
		t.Fatalf("unsafe provider id was emitted:\n%q", out)
	}
}

// An unknown argument is rejected with a usage line at exit 2 — the sibling porcelain
// norm — driven directly since the usage rides stdout. Keys are pinned empty so the
// pre-guard path (which reaches inventory) never touches the network.
func TestCommandRejectsUnknownArg(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "")
	out, code := Command([]string{"bogus"})
	if code != 2 {
		t.Fatalf("code = %d, want 2; out=%s", code, out)
	}
	if !strings.Contains(out, "usage") {
		t.Fatalf("missing usage line: %q", out)
	}
}

func TestParseDataIDs(t *testing.T) {
	got, err := parseDataIDs([]byte(`{"data":[{"id":"gpt-5.4"},{"id":"openai/gpt-5"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(got, ",") != "gpt-5.4,openai/gpt-5" {
		t.Fatalf("parseDataIDs = %v", got)
	}
	if _, err := parseDataIDs([]byte(`{"data":[{"id":`)); err == nil {
		t.Fatal("parseDataIDs malformed body: want error")
	}
}

func stubCommand(t *testing.T, fn func(string, ...string) ([]byte, error)) {
	t.Helper()
	old := runCommand
	runCommand = fn
	t.Cleanup(func() { runCommand = old })
}

func stubHTTP(t *testing.T, fn func(*http.Request) (*http.Response, error)) {
	t.Helper()
	old := doHTTP
	doHTTP = fn
	t.Cleanup(func() { doHTTP = old })
}

func response(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
