package diff

import (
	"strings"
	"testing"
)

func TestProductionCommandDerivesBoundedActionsAndHonestEmptyHelp(t *testing.T) {
	root, _, _, _ := seedDivergedRepo(t)
	live, liveCode := runProductionDiff(t, root)
	if liveCode != 0 || !strings.HasSuffix(live, "help[1]{cmd,why}:\n  bench diff --full,inspect the complete patch\n") {
		t.Fatalf("bounded live help = (exit %d):\n%s", liveCode, live)
	}

	first := rawGitText(t, root, "rev-parse", "HEAD")
	second := commitFile(t, "second.txt", "second\n", "second feature")
	if first == second {
		t.Fatal("fixture commits are not distinct")
	}
	for _, sha := range []string{first, second} {
		out, code := runProductionDiff(t, root, "--commit", sha)
		want := "help[1]{cmd,why}:\n  bench diff --full --commit " + sha + ",inspect the complete patch\n"
		if code != 0 || !strings.HasSuffix(out, want) {
			t.Fatalf("bounded commit %s help = (exit %d):\n%s", sha, code, out)
		}
	}

	for _, args := range [][]string{{"--full"}, {"--commit", first, "--full"}, {"--commit", second, "--full"}} {
		out, code := runProductionDiff(t, root, args...)
		if code != 0 || !strings.HasSuffix(out, "help[0]{cmd,why}:\n") {
			t.Fatalf("complete response %v help = (exit %d):\n%s", args, code, out)
		}
	}
}

func TestDriftHelpRepeatsTheExactRefusedInvocation(t *testing.T) {
	seedDivergedRepo(t)
	mustWriteFile(t, "f.txt", "dirty 0\n")
	previous := snapshotAfterRead
	defer func() { snapshotAfterRead = previous }()
	calls := 0
	snapshotAfterRead = func() {
		calls++
		mustWriteFile(t, "f.txt", "dirty "+string(rune('0'+calls))+"\n")
	}
	out, code := Command([]string{"--full"})
	want := "help[1]{cmd,why}:\n  bench diff --full,retry after the repository stopped moving\n"
	if code != 1 || calls != 2 || !strings.HasSuffix(out, want) {
		t.Fatalf("full drift help = (exit %d, calls %d):\n%s", code, calls, out)
	}
}
