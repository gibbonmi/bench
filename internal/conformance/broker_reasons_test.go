package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// landRouteRefusalReasons is the authored list of the five reasons a landing refuses an
// unauthenticated promotion broker on. The wrapper's land route is shell that runs before
// any binary is trusted, so it cannot call the doctor row, and the row cannot call it.
// Two derivations are therefore necessary, and this list is what keeps them one contract:
// a reason dropped from either side reds here. Each entry is the variable-free part of
// the sentence, because the route interpolates a path or a version into every one.
var landRouteRefusalReasons = []string{
	"no promotion-broker manifest at ",
	" is incomplete",
	" does not match installed package ",
	" is not a regular executable",
	" does not match its manifest digest",
}

// landRouteFile refuses the landing. brokerRowFile predicts that refusal in bench doctor.
const (
	landRouteFile = "bin/bench.sh"
	brokerRowFile = "internal/adopt/doctor_rows.go"
)

// checkBrokerReasons reports every authored reason a side no longer carries. It stays
// silent where either file is absent, because only the kit owns both derivations.
func checkBrokerReasons(root string) []string {
	sides := map[string]string{}
	for _, rel := range []string{landRouteFile, brokerRowFile} {
		text := readIfExists(filepath.Join(root, filepath.FromSlash(rel)))
		if text == "" {
			return nil
		}
		sides[rel] = text
	}
	var diags []string
	for _, reason := range landRouteRefusalReasons {
		for _, rel := range []string{landRouteFile, brokerRowFile} {
			if !strings.Contains(sides[rel], reason) {
				diags = append(diags, fmt.Sprintf("%s no longer carries the broker refusal reason %q; the land route and the bench doctor broker row must name the same five reasons", rel, reason))
			}
		}
	}
	return uniqueSorted(diags)
}

// TestBrokerRefusalReasonsBites is the recorded bite proof. The fixture writes both sides
// from the authored list itself, so the fixture carries no second copy of the reasons,
// then drops one reason from each side in turn.
func TestBrokerRefusalReasonsBites(t *testing.T) {
	root := t.TempDir()
	write := func(rel string, reasons []string) {
		t.Helper()
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(strings.Join(reasons, "\n")+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write(landRouteFile, landRouteRefusalReasons)
	write(brokerRowFile, landRouteRefusalReasons)
	if diags := checkBrokerReasons(root); len(diags) != 0 {
		t.Fatalf("both sides complete: want no diagnostics, got %v", diags)
	}

	for _, rel := range []string{landRouteFile, brokerRowFile} {
		dropped := landRouteRefusalReasons[len(landRouteRefusalReasons)-1]
		write(rel, landRouteRefusalReasons[:len(landRouteRefusalReasons)-1])
		if !containsDiagnostic(checkBrokerReasons(root), rel+" no longer carries the broker refusal reason "+fmt.Sprintf("%q", dropped)) {
			t.Fatalf("reason dropped from %s: want a diagnostic naming it, got %v", rel, checkBrokerReasons(root))
		}
		write(rel, landRouteRefusalReasons)
	}
}
