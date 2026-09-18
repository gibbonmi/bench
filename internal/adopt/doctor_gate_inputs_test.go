package adopt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDoctorGateInputsRowNamesEachUnsuppliedName is the FT315 doctor row. A gate input
// manifest that declares a name this process does not supply is red, and the row names
// that name alone. A manifest whose every name is supplied says nothing.
func TestDoctorGateInputsRowNamesEachUnsuppliedName(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, environment string
		wantOK            bool
		wantNamed         string
	}{
		{name: "supplied", environment: `["HOME"]`, wantOK: true},
		{name: "unsupplied", environment: `["FT315_DOCTOR_UNSUPPLIED", "HOME"]`, wantNamed: "FT315_DOCTOR_UNSUPPLIED"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			if err := os.Mkdir(filepath.Join(root, ".bench"), 0o755); err != nil {
				t.Fatal(err)
			}
			manifest := `{"schema":1,"closure":"local","environment":` + tc.environment + `,"paths":[],"tools":[]}` + "\n"
			if err := os.WriteFile(filepath.Join(root, ".bench", "gate-inputs.json"), []byte(manifest), 0o644); err != nil {
				t.Fatal(err)
			}
			ok, message := evalGateInputsRow(root)
			if tc.wantOK {
				if !ok || message != "" {
					t.Fatalf("supplied row = (%v, %q), want a silent ok", ok, message)
				}
				return
			}
			if ok || !strings.Contains(message, "supply: "+tc.wantNamed+" (") || strings.Contains(message, "HOME") {
				t.Fatalf("unsupplied row = (%v, %q), want red naming %s alone", ok, message, tc.wantNamed)
			}
		})
	}
}
