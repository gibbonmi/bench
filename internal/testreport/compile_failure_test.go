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
	if err := os.WriteFile(filepath.Join(root, "chosen", "broken.go"), []byte("package chosen\nvar _ = missingIdentifier\nvar _ = secondMissingIdentifier\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, full := range []bool{false, true} {
		args := []string{"--package", "chosen", "--run", "^TestChosen$"}
		if full {
			args = append(args, "--full")
		}
		output, code := Command(root, args)
		if code != 1 || !strings.Contains(output, "undefined: missingIdentifier") || strings.Contains(output, "run pattern matched no tests") {
			t.Fatalf("filtered compile failure = (%d, %q), want compiler diagnostic", code, output)
		}
		if strings.Contains(output, "undefined: secondMissingIdentifier") != full {
			t.Fatalf("full=%t second compiler diagnostic mismatch:\n%s", full, output)
		}
	}
}
