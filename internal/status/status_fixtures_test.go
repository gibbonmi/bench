// Shared test fixtures for package status: repository setup and gate-cache records.
package status

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/git"
)

func initRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	gitRun(t, root, "init")
	gitRun(t, root, "config", "user.email", "t@example.com")
	gitRun(t, root, "config", "user.name", "t")
	// The capture surfaces moved under a directory the repository root no longer
	// supplies for free. A fixture root that omits it fails every write that used
	// to land beside ROADMAP.md.
	if err := os.MkdirAll(filepath.Join(root, "capture"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

func gitRun(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, err := git.Output(append([]string{"-C", root}, args...)...)
	if err != nil {
		t.Fatalf("git %v: %v (%s)", args, err, out)
	}
	return out
}

// writeFullGateCache installs a ready full-class record naming cachedTree with the given
// verdict. It uses the mode and field set the gate's loader requires of the full class.
func writeFullGateCache(t *testing.T, root, cachedTree, status string) {
	t.Helper()
	gitdir := gitRun(t, root, "rev-parse", "--absolute-git-dir")
	recorded := time.Now().UTC().Truncate(time.Second).Add(-time.Minute).Format(time.RFC3339)
	record := fmt.Sprintf(`{"schema":1,"state":"ready","status":%q,"tree":%q,"oracle":%q,"recorded_at":%q}`+"\n",
		status, cachedTree, strings.Repeat("0", 64), recorded)
	if err := os.WriteFile(filepath.Join(gitdir, git.GateCacheFile), []byte(record), 0o600); err != nil {
		t.Fatal(err)
	}
}

// writePartialGateCache installs a ready partial green naming cachedTree, skipping the
// given components. It uses the mode and exact field set the gate's loader requires of
// the partial class. Each skipped component carries ancestor-form evidence (an identity
// and the time it was authored), the simpler of the two forms validatePartition accepts.
func writePartialGateCache(t *testing.T, root, cachedTree string, skipped ...string) {
	t.Helper()
	gitdir := gitRun(t, root, "rev-parse", "--absolute-git-dir")
	recorded := time.Now().UTC().Truncate(time.Second).Add(-time.Minute).Format(time.RFC3339)
	authoredAt := time.Now().UTC().Truncate(time.Second).Add(-time.Hour).Format(time.RFC3339)
	identity := strings.Repeat("b", 64)

	executedJSON, err := json.Marshal([]string{"core"})
	if err != nil {
		t.Fatal(err)
	}
	skippedJSON, err := json.Marshal(skipped)
	if err != nil {
		t.Fatal(err)
	}
	evidence := make(map[string]map[string]string, len(skipped))
	for _, component := range skipped {
		evidence[component] = map[string]string{"identity": identity, "authored_at": authoredAt}
	}
	evidenceJSON, err := json.Marshal(evidence)
	if err != nil {
		t.Fatal(err)
	}

	record := fmt.Sprintf(`{"schema":1,"state":"ready","status":"green","tree":%q,"oracle":%q,"recorded_at":%q,"executed":%s,"skipped":%s,"skip_evidence":%s}`+"\n",
		cachedTree, strings.Repeat("0", 64), recorded, executedJSON, skippedJSON, evidenceJSON)
	path := filepath.Join(gitdir, git.GateCacheFile)
	if err := os.WriteFile(path, []byte(record), 0o600); err != nil {
		t.Fatal(err)
	}
}
