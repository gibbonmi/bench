package runtime

import (
	"os"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/status"
)

func provePresentation(t *testing.T, variant string) {
	switch variant {
	case "html":
		proveSelfContainedDashboard(t)
	case "axi":
		proveAXIGateCache(t)
	case "hostile-status", "hostile-dashboard", "hostile-roadmap":
		proveHostileProjection(t, strings.TrimPrefix(variant, "hostile-"))
	default:
		proveSignalSeverity(t, variant)
	}
}

func proveSignalSeverity(t *testing.T, state string) {
	p := newProjectionFixture(t, state)
	wantSeverity := map[string]int{"red": 0, "locked-pending": 1, "interrupted-pending": 2, "invalid": 3, "unavailable": 3, "stale": 7}
	wantDetail := map[string]string{"red": "red", "locked-pending": "locked-pending", "interrupted-pending": "interrupted-pending", "invalid": "invalid verdict", "unavailable": "verdict unavailable", "stale": "stale"}
	gateSignals := 0
	for _, signal := range status.Signals(p.f.Root) {
		if signal.Name != "gate" {
			continue
		}
		gateSignals++
		if signal.Severity != wantSeverity[state] || !strings.Contains(signal.Detail, wantDetail[state]) {
			t.Fatalf("%s gate signal = %+v, want severity %d and detail %q", state, signal, wantSeverity[state], wantDetail[state])
		}
	}
	show := state != "absent" && state != "reusable-green"
	if (gateSignals == 1) != show {
		t.Fatalf("%s gate signal count = %d, show = %v", state, gateSignals, show)
	}
	statusOut := runProjectionSurface(t, p, "status")
	dashboardOut := runProjectionSurface(t, p, "dashboard")
	if show {
		if !strings.Contains(statusOut, "gate") || !strings.Contains(dashboardOut, wantDetail[state]) {
			t.Fatalf("%s signal missing from a public human format", state)
		}
	} else if strings.Contains(statusOut, "gate       ") || strings.Contains(dashboardOut, "<td>gate</td>") {
		t.Fatalf("%s rendered a gate signal without a signal", state)
	}
}

func proveSelfContainedDashboard(t *testing.T) {
	p := newProjectionFixture(t, "reusable-green")
	page := runProjectionSurface(t, p, "dashboard")
	for _, literal := range []string{"<!DOCTYPE html>", "<html", "<head>", "<style>", "</style>", "<body>", "</body>", "</html>"} {
		if !strings.Contains(page, literal) {
			t.Fatalf("dashboard HTML missing %q", literal)
		}
	}
	if !strings.HasPrefix(page, "<!DOCTYPE html>") || !strings.HasSuffix(strings.TrimSpace(page), "</html>") {
		t.Fatal("dashboard is not one complete HTML document")
	}
	for _, external := range []string{"<link ", "<script src=", "http://", "https://"} {
		if strings.Contains(page, external) {
			t.Fatalf("dashboard is not self-contained: contains %q", external)
		}
	}
}

func proveAXIGateCache(t *testing.T) {
	p := newProjectionFixture(t, "locked-pending")
	out := runProjectionSurface(t, p, "roadmap")
	want := "gate_cache[1]{present,state,pending_status,status,cached_tree,work_tree,timestamp,stale}:"
	if strings.Count(out, want) != 1 {
		t.Fatalf("AXI gate-cache schema count = %d, want 1\n%s", strings.Count(out, want), out)
	}
	if !strings.Contains(out, "true,pending,locked-pending") {
		t.Fatalf("AXI gate-cache row lost typed pending fields:\n%s", out)
	}
}

func proveHostileProjection(t *testing.T, surface string) {
	p := newProjectionFixture(t, "invalid")
	hostile := "SECRET-FT78\x1b<script>alert(1)</script>\x00\x7f\n"
	if err := os.WriteFile(p.cache, []byte(hostile), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(p.cache, 0o600); err != nil {
		t.Fatal(err)
	}
	out := runProjectionSurface(t, p, surface)
	for _, raw := range []string{"SECRET-FT78", "<script>", "alert(1)", "\x1b", "\x00", "\x7f"} {
		if strings.Contains(out, raw) {
			t.Fatalf("%s leaked hostile cache bytes %q", surface, raw)
		}
	}
}
