package canary

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

func canaryFixture(root, family, name string) string {
	return filepath.Join(root, "tests", "canary", family, name)
}

func mappedFamily(t *testing.T) string {
	t.Helper()
	families := registry.Families()
	if len(families) == 0 {
		t.Fatal("conformance family registry is empty")
	}
	return families[0]
}

func mkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}

func write(t *testing.T, path, content string) {
	t.Helper()
	mkdir(t, filepath.Dir(path))
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
