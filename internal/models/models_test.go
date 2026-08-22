package models

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestCommandOfflineSuppressesEveryProviderWithEvidence(t *testing.T) {
	t.Setenv("BENCH_OFFLINE", "1")
	t.Setenv("OPENAI_API_KEY", "present")
	t.Setenv("ANTHROPIC_API_KEY", "present")
	stubCommand(t, func(string, ...string) ([]byte, error) {
		t.Fatal("offline discovery started codex")
		return nil, nil
	})
	stubHTTP(t, func(*http.Request) (*http.Response, error) {
		t.Fatal("offline discovery started HTTP")
		return nil, nil
	})
	for run := 1; run <= 2; run++ {
		out, code := Command(nil)
		if code != 0 {
			t.Fatalf("run %d code = %d; out=%s", run, code, out)
		}
		for _, source := range []string{"codex", "openai", "anthropic"} {
			want := source + ",offline,offline,BENCH_OFFLINE=1"
			if !strings.Contains(out, want) {
				t.Fatalf("run %d missing %q:\n%s", run, want, out)
			}
		}
	}
}

func TestInventoryStartsProvidersConcurrentlyAndRendersStableOrder(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "key")
	t.Setenv("ANTHROPIC_API_KEY", "key")
	var mu sync.Mutex
	started := 0
	timedOut := false
	allStarted := make(chan struct{})
	arrive := func() {
		mu.Lock()
		started++
		if started == 3 {
			close(allStarted)
		}
		mu.Unlock()
		select {
		case <-allStarted:
		case <-time.After(200 * time.Millisecond):
			mu.Lock()
			timedOut = true
			mu.Unlock()
		}
	}
	stubCommand(t, func(string, ...string) ([]byte, error) {
		arrive()
		return []byte(`{"models":[{"slug":"codex-id","visibility":"list"}]}`), nil
	})
	stubHTTP(t, func(req *http.Request) (*http.Response, error) {
		arrive()
		id := "openai-id"
		if req.URL.String() == anthropicModelsURL {
			id = "anthropic-id"
		}
		return response(200, `{"data":[{"id":"`+id+`"}]}`), nil
	})
	out, code := Command(nil)
	if code != 0 || started != 3 || timedOut {
		t.Fatalf("code/started/timedOut = %d/%d/%v; out=%s", code, started, timedOut, out)
	}
	if !strings.Contains(out, "  codex,live,codex-id\n  openai,live,openai-id\n  anthropic,live,anthropic-id\n") {
		t.Fatalf("provider render order changed:\n%s", out)
	}
}

func TestModelBodyLimitAcceptsExactAndRejectsPlusOne(t *testing.T) {
	base := `{"data":[]}`
	for _, tt := range []struct {
		name    string
		size    int
		wantErr bool
	}{
		{name: "exact", size: 5 << 20},
		{name: "plus one", size: (5 << 20) + 1, wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stubHTTP(t, func(*http.Request) (*http.Response, error) {
				return response(200, base+strings.Repeat(" ", tt.size-len(base))), nil
			})
			_, err := fetchDataIDs("https://example.invalid", "key", openAIHeaders)
			if (err != nil) != tt.wantErr {
				t.Fatalf("size %d error = %v, wantErr=%v", tt.size, err, tt.wantErr)
			}
		})
	}
}

func TestProviderTimeoutDoesNotCollapsePeers(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "key")
	t.Setenv("ANTHROPIC_API_KEY", "key")
	oldTimeout := providerTimeout
	providerTimeout = 20 * time.Millisecond
	t.Cleanup(func() { providerTimeout = oldTimeout })
	stubCommand(t, func(string, ...string) ([]byte, error) {
		time.Sleep(30 * time.Millisecond)
		return nil, providerError{status: "timeout", err: context.DeadlineExceeded}
	})
	stubHTTP(t, func(*http.Request) (*http.Response, error) {
		return response(200, `{"data":[{"id":"peer"}]}`), nil
	})
	out, code := Command(nil)
	if code != 0 || !strings.Contains(out, "codex,timeout,timeout,10s provider deadline") || strings.Count(out, ",live,peer") != 2 {
		t.Fatalf("timeout collapsed peers: code=%d\n%s", code, out)
	}
}

func TestCodexLiveAndBundledSharePartiallyConsumedDeadline(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "bundled-started")
	script := filepath.Join(dir, "codex")
	body := `#!/bin/sh
if [ "${3:-}" = "--bundled" ]; then
  : > "$BENCH_CODEX_BUNDLED_MARKER"
  sleep 0.25
  printf '{"models":[{"slug":"bundled","visibility":"list"}]}'
  exit 0
fi
sleep 0.25
exit 1
`
	if err := os.WriteFile(script, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("BENCH_CODEX_BUNDLED_MARKER", marker)
	oldTimeout := providerTimeout
	providerTimeout = 400 * time.Millisecond
	t.Cleanup(func() { providerTimeout = oldTimeout })

	row, models := codexInventory()
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("bundled attempt never started: %v", err)
	}
	if row.status != "timeout" || row.hint != "10s provider deadline" || len(models) != 0 {
		t.Fatalf("partially consumed shared budget = row:%+v models:%+v", row, models)
	}
}

func TestAPIInventoryTimesOutAsUnavailable(t *testing.T) {
	oldTimeout := providerTimeout
	providerTimeout = 20 * time.Millisecond
	t.Cleanup(func() { providerTimeout = oldTimeout })
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	t.Cleanup(server.Close)

	started := time.Now()
	row, models := apiInventory("openai", server.URL, "key", openAIHeaders, "set key")
	if elapsed := time.Since(started); elapsed > time.Second {
		t.Fatalf("hung provider returned after %s, want at most 1s", elapsed)
	}
	if row.status != "timeout" || row.hint != "10s provider deadline" {
		t.Fatalf("timeout row = %+v, want distinct timeout", row)
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

// An unknown argument gives a usage line at exit 2, the sibling porcelain norm. The
// test drives this directly, since the usage rides stdout. Keys are pinned empty, so
// the pre-guard path, which reaches inventory, never touches the network.
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
	runCommand = func(_ context.Context, name string, args ...string) ([]byte, error) { return fn(name, args...) }
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
