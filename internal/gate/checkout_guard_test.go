package gate

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"strings"
	"testing"
)

func TestCheckoutGuardRedsAnAddedPath(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[{"name":"write","argv":["sh","-c","touch stray"]}]}`)
	code, stdout, stderr := runFixturePhases(context.Background(), t, root)
	if code != 1 {
		t.Fatalf("checkout write exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if rows := rowsForPhase(t, stdout, capabilityPhase); len(rows) != 1 || !strings.Contains(rows[0], `"stray"`) {
		t.Fatalf("checkout rows = %q, want one row naming stray", rows)
	}
}

func TestCheckoutGuardRedsAnIgnoredRewrite(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[{"name":"write","argv":["sh","-c","printf after > ignored"]}]}`)
	for name, data := range map[string]string{".gitignore": "ignored\n", "ignored": "prior"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	old := time.Unix(1, 0)
	if err := os.Chtimes(filepath.Join(root, "ignored"), old, old); err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := runFixturePhases(context.Background(), t, root)
	if code != 1 {
		t.Fatalf("ignored rewrite exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if rows := rowsForPhase(t, stdout, capabilityPhase); len(rows) != 1 || !strings.Contains(rows[0], `"ignored"`) {
		t.Fatalf("checkout rows = %q, want one row naming ignored", rows)
	}
}

func TestCheckoutGuardRedsARemovedPath(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[{"name":"remove","argv":["sh","-c","rm stray"]}]}`)
	if err := os.WriteFile(filepath.Join(root, "stray"), []byte("before"), 0o600); err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := runFixturePhases(context.Background(), t, root)
	if code != 1 {
		t.Fatalf("checkout removal exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if rows := rowsForPhase(t, stdout, capabilityPhase); len(rows) != 1 || !strings.Contains(rows[0], `"stray"`) {
		t.Fatalf("checkout rows = %q, want one row naming stray", rows)
	}
}

func TestCheckoutGuardRedsAGitRecord(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[{"name":"write","argv":["sh","-c","touch .git/bench-stray"]}]}`)
	code, stdout, stderr := runFixturePhases(context.Background(), t, root)
	if code != 1 {
		t.Fatalf("git record exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if rows := rowsForPhase(t, stdout, capabilityPhase); len(rows) != 1 || !strings.Contains(rows[0], `".git/bench-stray"`) {
		t.Fatalf("checkout rows = %q, want one row naming .git/bench-stray", rows)
	}
}

func TestCheckoutGuardPermitsTiming(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[{"name":"timing","argv":["sh","-c","printf timing > \"$TIMING\""]}]}`)
	t.Setenv("TIMING", registry.TimingPath(root))
	code, stdout, stderr := runFixturePhases(context.Background(), t, root)
	if code != 0 {
		t.Fatalf("timing write exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if data, err := os.ReadFile(registry.TimingPath(root)); err != nil || string(data) != "timing" {
		t.Fatalf("timing write did not run: %q, %v", data, err)
	}
}

func TestCheckoutGuardPermitsTheCurrentRunLog(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[{"name":"output","argv":["echo","phase output"]}]}`)
	stubGateLogPathIgnored(t)
	var diagnostic bytes.Buffer
	ctx, finish := beginGateRunLog(context.Background(), root, &diagnostic, "dev")
	defer finish(Result{})
	log, _ := ctx.Value(gateRunLogKey{}).(*gateRunLog)
	if log == nil || log.streamFile == nil {
		t.Fatalf("run records not opened: %s", diagnostic.String())
	}
	before, err := log.file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := runFixturePhases(ctx, t, root)
	if code != 0 {
		t.Fatalf("logged run exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	after, err := log.file.Stat()
	if err != nil || after.Size() <= before.Size() {
		t.Fatalf("run log did not grow: %v, %v", after, err)
	}
	if data, err := os.ReadFile(log.streamPath()); err != nil || !strings.Contains(string(data), "phase output") {
		t.Fatalf("run stream did not grow: %q, %v", data, err)
	}
}

func TestCheckoutGuardOrdersChangedPaths(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[{"name":"write","argv":["sh","-c","touch b-stray a-stray"]}]}`)
	code, stdout, stderr := runFixturePhases(context.Background(), t, root)
	if code != 1 {
		t.Fatalf("checkout writes exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	rows := rowsForPhase(t, stdout, capabilityPhase)
	if len(rows) != 2 || !strings.Contains(rows[0], `"a-stray"`) || !strings.Contains(rows[1], `"b-stray"`) {
		t.Fatalf("checkout rows = %q, want a-stray before b-stray", rows)
	}
}

func TestCheckoutGuardRedsUnreadableState(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[{"name":"remove","argv":["sh","-c","rm -rf .git"]}]}`)
	code, stdout, stderr := runFixturePhases(context.Background(), t, root)
	if code != 1 {
		t.Fatalf("unreadable checkout exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if rows := rowsForPhase(t, stdout, capabilityPhase); len(rows) != 1 || !strings.Contains(rows[0], "unreadable checkout state") {
		t.Fatalf("checkout rows = %q, want one unreadable-state row", rows)
	}
}

func TestCheckoutGuardLeavesLinkedPhasesAlone(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[{"name":"write","argv":["sh","-c","touch stray"],"dir":"."}]}`)
	var stdout, stderr bytes.Buffer
	code := phasesCommandAtKitWithSelection(context.Background(), root, t.TempDir(), fixtureSelection(root), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("linked checkout write exit = %d, want 0; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(filepath.Join(root, "stray")); err != nil {
		t.Fatalf("linked root write did not run: %v", err)
	}
}

func TestCheckoutGuardEscapesControlBytes(t *testing.T) {
	name := "stray\n\t\x1b"
	argv, err := json.Marshal([]string{"touch", name})
	if err != nil {
		t.Fatal(err)
	}
	root := fixturePhaseRoot(t, `{"phases":[{"name":"write","argv":`+string(argv)+`}]}`)
	code, stdout, stderr := runFixturePhases(context.Background(), t, root)
	if code != 1 {
		t.Fatalf("control-byte path exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if rows := rowsForPhase(t, stdout, capabilityPhase); len(rows) != 1 || !strings.Contains(rows[0], `stray\n\t\x1b`) {
		t.Fatalf("checkout rows = %q, want the escaped path in one row", rows)
	}
}

func TestCheckoutGuardFollowsTheRootAlias(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[{"name":"write","argv":["sh","-c","touch stray"]}]}`)
	alias := filepath.Join(t.TempDir(), "root alias")
	if err := os.Symlink(root, alias); err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := runFixturePhases(context.Background(), t, alias)
	if code != 1 {
		t.Fatalf("aliased checkout write exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if rows := rowsForPhase(t, stdout, capabilityPhase); len(rows) != 1 || !strings.Contains(rows[0], `"stray"`) {
		t.Fatalf("checkout rows = %q, want one row naming stray", rows)
	}
}

func TestCheckoutGuardDoesNotFollowFileLinks(t *testing.T) {
	outside := filepath.Join(t.TempDir(), "outside")
	if err := os.WriteFile(outside, []byte("before"), 0o600); err != nil {
		t.Fatal(err)
	}
	root := fixturePhaseRoot(t, `{"phases":[{"name":"write","argv":["sh","-c","printf changed > link"]}]}`)
	if err := os.Symlink(outside, filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("absent", filepath.Join(root, "dangling")); err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := runFixturePhases(context.Background(), t, root)
	if code != 0 {
		t.Fatalf("file-link metadata exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if data, err := os.ReadFile(outside); err != nil || string(data) != "changed" {
		t.Fatalf("outside write did not run: %q, %v", data, err)
	}
}

func TestCheckoutGuardCannotHideAnAddedPathInTheIndex(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[{"name":"write","argv":["sh","-c","touch stray; git add stray"]}]}`)
	code, stdout, stderr := runFixturePhases(context.Background(), t, root)
	if code != 1 {
		t.Fatalf("staged checkout write exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if rows := rowsForPhase(t, stdout, capabilityPhase); len(rows) != 1 || !strings.Contains(rows[0], `"stray"`) {
		t.Fatalf("checkout rows = %q, want one row naming stray", rows)
	}
}

func TestCheckoutGuardRedsAnOlderRunLog(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[{"name":"write","argv":["sh","-c","printf changed > \"$OLD_LOG\""]}]}`)
	old := gateLogRecordPath(root, gateLogRunToken(time.Unix(1, 0), 1))
	if err := os.MkdirAll(filepath.Dir(old), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(old, []byte("prior"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("OLD_LOG", old)
	code, stdout, stderr := runFixturePhases(context.Background(), t, root)
	if code != 1 {
		t.Fatalf("older log rewrite exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if rows := rowsForPhase(t, stdout, capabilityPhase); len(rows) != 1 || !strings.Contains(rows[0], filepath.Base(old)) {
		t.Fatalf("checkout rows = %q, want one row naming the older log", rows)
	}
}

func TestCheckoutGuardPermitsRunRecordsThroughARootAlias(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[{"name":"output","argv":["sh","-c","printf lock > \"$LOCK\"; printf owner > \"$OWNER\"; echo phase output"]}]}`)
	alias := filepath.Join(t.TempDir(), "root alias")
	if err := os.Symlink(root, alias); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LOCK", gateLockPath(filepath.Join(root, ".git")))
	t.Setenv("OWNER", gateOwnerPath(filepath.Join(root, ".git")))
	stubGateLogPathIgnored(t)
	var diagnostic bytes.Buffer
	ctx, finish := beginGateRunLog(context.Background(), alias, &diagnostic, "dev")
	defer finish(Result{})
	if gateRunStreamFile(ctx) == nil {
		t.Fatalf("run records not opened: %s", diagnostic.String())
	}
	code, stdout, stderr := runFixturePhases(ctx, t, alias)
	if code != 0 {
		t.Fatalf("aliased run records exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	for _, path := range []string{os.Getenv("LOCK"), os.Getenv("OWNER")} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("record write did not run: %v", err)
		}
	}
}
