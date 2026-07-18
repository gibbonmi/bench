package artifact

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestArtifactPromotionIsAtomicAndExclusive(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	source := committedHostileArtifactSource(t, root)
	prepared, expected := prepareArtifactGeneration(t, source)
	output := filepath.Join(t.TempDir(), "artifact output")
	assertConcurrentFirstArtifactPromotion(t, source, prepared, output, expected)
	assertInterruptedArtifactPromotion(t, source, prepared, output, expected)
}

func assertConcurrentFirstArtifactPromotion(t *testing.T, root, prepared, output string, expected int) {
	t.Helper()
	ready := filepath.Join(t.TempDir(), "winner-ready")
	winner := exec.Command("bash", filepath.Join(root, "scripts", "build-artifacts.sh"), root, output)
	winner.Env = promotionTestEnv(prepared, ready)
	if err := winner.Start(); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(30 * time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("concurrent artifact builders did not reach promotion")
		}
		time.Sleep(20 * time.Millisecond)
	}
	loser := exec.Command("bash", filepath.Join(root, "scripts", "build-artifacts.sh"), root, output)
	loser.Env = promotionTestEnv(prepared, ready)
	if err := loser.Run(); err == nil {
		t.Fatal("concurrent artifact builder did not fail closed")
	}
	if err := os.Remove(ready); err != nil {
		t.Fatal(err)
	}
	if err := winner.Wait(); err != nil {
		t.Fatalf("artifact winner failed: %v", err)
	}
	entries, err := os.ReadDir(output)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != expected {
		t.Fatalf("concurrent artifact promotion entries=%d, want one complete %d-file generation", len(entries), expected)
	}
	for _, entry := range entries {
		if !entry.Type().IsRegular() || (!strings.HasSuffix(entry.Name(), ".tgz") && !strings.HasSuffix(entry.Name(), ".tar.gz")) {
			t.Fatalf("concurrent artifact promotion left nested or non-tar output: %s", entry.Name())
		}
	}
}
