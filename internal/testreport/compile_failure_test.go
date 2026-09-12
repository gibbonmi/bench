package testreport

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunPatternReportsCompilerDiagnostic(t *testing.T) {
	root := focusedTestModule(t)
	installCannedSelection(t)
	if err := os.WriteFile(filepath.Join(root, "chosen", "broken.go"), []byte("package chosen\nvar _ = missingIdentifier\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	output, code := Command(root, []string{"--package", "chosen", "--run", "^TestChosen$"})
	if code != 1 || !strings.Contains(output, "undefined: missingIdentifier") || strings.Contains(output, "run pattern matched no tests") {
		t.Fatalf("filtered compile failure = (%d, %q), want compiler diagnostic", code, output)
	}
}
