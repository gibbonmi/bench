package sessioninspect

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/compatibility"
)

func TestCompatibilityUnknownContext(t *testing.T) {
	session := compatibility.Session{ID: "chat", Epoch: "start", Context: compatibility.Context{
		Interface:         compatibility.CodexDesktop,
		Repository:        compatibility.Fact{Value: "/repo", Source: "git"},
		Environment:       compatibility.Fact{Value: "linux/wsl2", Source: "runtime"},
		ConfigurationHome: compatibility.Fact{Value: "/config", Source: "harness"},
		ActiveRuntime:     compatibility.Fact{Value: "runtime-a", Source: "harness"},
		PolicyProvenance:  compatibility.Fact{Value: "policy-a", Source: "harness"},
	}}
	request := compatibility.SessionRequest{Session: session, Operation: "diagnose", Observations: []compatibility.Observation{{
		Capability: "normal-shell", Operation: "diagnose", Route: "chat-shell", Provenance: "actual-tool", Permission: "normal", Session: session, Success: true,
	}}}
	output, err := compatibilityInspection(request)
	if err != nil || !strings.Contains(output, "normal-shell,ok") {
		t.Fatalf("current observed shell is not usable: %q, %v", output, err)
	}
	request.Session.Context.ActiveRuntime.Value = "unknown"
	output, err = compatibilityInspection(request)
	if err != nil || !strings.Contains(output, "normal-shell,unknown") {
		t.Fatalf("unknown runtime reused prior success: %q, %v", output, err)
	}
	request.Observations[0].Session = request.Session
	output, err = compatibilityInspection(request)
	if err != nil || !strings.Contains(output, "normal-shell,unknown") {
		t.Fatalf("matching unknown contexts reused success: %q, %v", output, err)
	}
}

func TestCompatibilityStartup(t *testing.T) {
	output := inspectCompatibility(t)
	if !strings.HasPrefix(output, "bench: on start or resume") || !strings.Contains(output, "repository-wrapper") {
		t.Fatalf("startup omitted active-interface obligation: %q", output)
	}
}

func TestCompatibilityRetestRequired(t *testing.T) {
	output := inspectCompatibility(t)
	if !strings.Contains(output, "previously failed interface with normal permissions after repair") {
		t.Fatalf("startup omitted failed-interface retest: %q", output)
	}
}

func TestCompatibilityUnknownRepair(t *testing.T) {
	output := inspectCompatibility(t)
	if !strings.Contains(output, "upstream-repair,unresolved") {
		t.Fatalf("startup implied an established upstream repair: %q", output)
	}
}

func TestCompatibilityStartupTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out bytes.Buffer
	if code := Inspect(ctx, &out, t.TempDir()); code != 0 || !strings.Contains(out.String(), "deadline exceeded") {
		t.Fatalf("cancelled startup = %d, %q", code, out.String())
	}
	if strings.Contains(out.String(), ",ok,") {
		t.Fatalf("cancelled inspection qualified a capability: %q", out.String())
	}
}

func inspectCompatibility(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	t.Chdir(root)
	var out bytes.Buffer
	if code := Inspect(context.Background(), &out, root); code != 0 {
		t.Fatalf("informational startup exit = %d", code)
	}
	return out.String()
}
