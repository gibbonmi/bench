package diff

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
)

func TestProductionCommandRawGitMatrix(t *testing.T) {
	t.Run("committed", func(t *testing.T) {
		root, _, base, _ := seedDivergedRepo(t)
		out, code := runProductionDiff(t, root)
		assertRawLiveSnapshot(t, root, base, out, code)
		head := rawGitText(t, root, "rev-parse", "HEAD")
		commitOut, commitCode := runProductionDiff(t, root, "--commit", head)
		if commitCode != 0 {
			t.Fatalf("commit diff exit = %d:\n%s", commitCode, commitOut)
		}
		for _, absent := range []string{"checkout[", "whitespace[", "staged,unstaged,untracked"} {
			if strings.Contains(commitOut, absent) {
				t.Errorf("immutable commit response contains live fact %q:\n%s", absent, commitOut)
			}
		}
	})

	for _, tc := range []struct {
		name   string
		mutate func(*testing.T, string)
	}{
		{name: "staged", mutate: func(t *testing.T, _ string) {
			mustWriteFile(t, "f.txt", "staged\n")
			runGit(t, "add", "f.txt")
		}},
		{name: "unstaged", mutate: func(t *testing.T, _ string) { mustWriteFile(t, "f.txt", "unstaged\n") }},
		{name: "untracked", mutate: func(t *testing.T, _ string) { mustWriteFile(t, "new.txt", "new\n") }},
		{name: "nested directory", mutate: func(t *testing.T, _ string) {
			if err := os.MkdirAll("new/deep", 0o755); err != nil {
				t.Fatal(err)
			}
			mustWriteFile(t, "new/deep/file.txt", "nested\n")
		}},
		{name: "rename", mutate: func(t *testing.T, _ string) {
			runGit(t, "mv", "f.txt", "renamed.txt")
		}},
		{name: "deletion", mutate: func(t *testing.T, _ string) {
			if err := os.Remove("f.txt"); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "whitespace offense", mutate: func(t *testing.T, _ string) { mustWriteFile(t, "f.txt", "trailing space \n") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, _, base, _ := seedDivergedRepo(t)
			tc.mutate(t, root)
			out, code := runProductionDiff(t, root)
			assertRawLiveSnapshot(t, root, base, out, code)
		})
	}

	t.Run("binary", func(t *testing.T) {
		root, _, base, _ := seedDivergedRepo(t)
		if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte{0, 1, 2, 3}, 0o644); err != nil {
			t.Fatal(err)
		}
		out, code := runProductionDiff(t, root, "--full")
		if code != 0 {
			t.Fatalf("binary full exit = %d:\n%s", code, out)
		}
		raw := rawGitBytes(t, root, false, "diff", base)
		if !bytes.Contains([]byte(out), append([]byte("diff_body:\n"), raw...)) {
			t.Fatalf("tracked binary diff body differs from raw Git:\n%s", out)
		}
	})

	t.Run("hostile filename control refusal", func(t *testing.T) {
		root, _, _, _ := seedDivergedRepo(t)
		mustWriteFile(t, "bad\x1bname", "hostile\n")
		raw := rawGitBytes(t, root, false, "status", "--porcelain=v1", "-z", "--no-renames", "--untracked-files=all")
		if !bytes.Contains(raw, []byte("bad\x1bname")) {
			t.Fatal("raw Git fixture lost the control-bearing path")
		}
		out, code := runProductionDiff(t, root)
		if code != 1 || !strings.Contains(out, "error: unrepresentable TOON cell") {
			t.Fatalf("control path response = (exit %d, %q), want structured refusal", code, out)
		}
	})

	t.Run("clean", func(t *testing.T) {
		root := t.TempDir()
		t.Chdir(root)
		runGit(t, "init", "-q", "-b", "main")
		runGit(t, "config", "user.email", "t@example.com")
		runGit(t, "config", "user.name", "t")
		base := commitFile(t, "clean.txt", "clean\n", "clean")
		out, code := runProductionDiff(t, root)
		assertRawLiveSnapshot(t, root, base, out, code)
		if !strings.HasSuffix(out, "help[0]{cmd,why}:\n") {
			t.Fatalf("clean response lacks honest empty help:\n%s", out)
		}
	})

	t.Run("detached HEAD", func(t *testing.T) {
		root, _, base, _ := seedDivergedRepo(t)
		runGit(t, "checkout", "-q", "--detach")
		out, code := runProductionDiff(t, root)
		assertRawLiveSnapshot(t, root, base, out, code)
		if got := tableRows(t, out, "revision")[0][0]; got != "(detached)" {
			t.Fatalf("detached revision branch = %q", got)
		}
	})

	t.Run("deep cwd", func(t *testing.T) {
		root, _, _, _ := seedDivergedRepo(t)
		deep := filepath.Join(root, "a", "b", "c")
		if err := os.MkdirAll(deep, 0o755); err != nil {
			t.Fatal(err)
		}
		fromRoot, rootCode := runProductionDiff(t, root)
		fromDeep, deepCode := runProductionDiff(t, deep)
		if rootCode != 0 || deepCode != 0 || fromDeep != fromRoot {
			t.Fatalf("deep cwd response differs: root=(%d,%q) deep=(%d,%q)", rootCode, fromRoot, deepCode, fromDeep)
		}
	})

	t.Run("base equals HEAD", func(t *testing.T) {
		root, _, _, _ := seedDivergedRepo(t)
		head := rawGitText(t, root, "rev-parse", "HEAD")
		runGit(t, "config", "branch.feature.benchBase", head)
		mustWriteFile(t, "only-untracked.txt", "new\n")
		out, code := runProductionDiff(t, root)
		assertRawLiveSnapshot(t, root, head, out, code)
		if got := tableRows(t, out, "aggregate")[0][0]; got != "0" {
			t.Fatalf("base=HEAD commits = %q, want 0", got)
		}
	})

	t.Run("still-tree byte idempotency", func(t *testing.T) {
		root, _, _, _ := seedDivergedRepo(t)
		mustWriteFile(t, "new.txt", "new\n")
		first, firstCode := runProductionDiff(t, root, "--full")
		second, secondCode := runProductionDiff(t, root, "--full")
		if firstCode != 0 || secondCode != 0 || first != second {
			t.Fatalf("still-tree invocations differ: first=(%d,%q) second=(%d,%q)", firstCode, first, secondCode, second)
		}
	})

	t.Run("mid-read drift", func(t *testing.T) {
		seedDivergedRepo(t)
		mustWriteFile(t, "f.txt", "dirty 0\n")
		assertRepeatedDrift(t, "dirty content", func(call int) {
			mustWriteFile(t, "f.txt", "dirty "+string(rune('0'+call))+"\n")
		})
	})
}

func TestProductionCommandFullUntrackedKindsUseRawGitBodies(t *testing.T) {
	root, _, _, _ := seedDivergedRepo(t)
	if err := os.WriteFile(filepath.Join(root, "b.bin"), []byte{0, 1, 2}, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteFile(t, "a.txt", "text\n")
	if err := syscall.Mkfifo(filepath.Join(root, "named-pipe"), 0o600); err != nil {
		t.Fatal(err)
	}
	out, code := runProductionDiff(t, root, "--full")
	if code != 0 {
		t.Fatalf("full untracked kinds exit = %d:\n%s", code, out)
	}
	aPatch := rawGitBytes(t, root, true, "diff", "--no-index", "--", "/dev/null", "a.txt")
	bPatch := rawGitBytes(t, root, true, "diff", "--no-index", "--", "/dev/null", "b.bin")
	if !bytes.Contains([]byte(out), aPatch) || !bytes.Contains([]byte(out), bPatch) || bytes.Index([]byte(out), aPatch) > bytes.Index([]byte(out), bPatch) {
		t.Fatalf("untracked patches differ from sorted raw Git bodies:\n%s", out)
	}
	if !bytes.Contains(bPatch, []byte("new file mode 100755")) || !bytes.Contains(bPatch, []byte("Binary files")) {
		t.Fatalf("binary fixture lacks raw mode/binary form:\n%s", bPatch)
	}
	if !strings.Contains(out, "?,named-pipe,fifo") || strings.Contains(out, "diff --git a/named-pipe") {
		t.Fatalf("FIFO was not named without being read:\n%s", out)
	}
	if got := tableRows(t, out, "aggregate")[0][6]; got != "3" {
		t.Fatalf("untracked aggregate omitted the FIFO: got %s, want 3", got)
	}
}

func TestProductionCommandReportsExactSpecialKinds(t *testing.T) {
	root, _, _, _ := seedDivergedRepo(t)
	mustWriteFile(t, "target.txt", "target\n")
	if err := os.Symlink("target.txt", filepath.Join(root, "linked")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("missing", filepath.Join(root, "dangling")); err != nil {
		t.Fatal(err)
	}
	out, code := runProductionDiff(t, root)
	if code != 0 {
		t.Fatalf("special-kind response exit = %d:\n%s", code, out)
	}
	for _, want := range []string{"?,linked,symlink", "?,dangling,dangling-symlink", "?,target.txt,\"\""} {
		if !strings.Contains(out, want) {
			t.Errorf("special-kind response missing %q:\n%s", want, out)
		}
	}
}

func assertRawLiveSnapshot(t *testing.T, root, base, out string, code int) {
	t.Helper()
	if code != 0 {
		t.Fatalf("production diff exit = %d:\n%s", code, out)
	}
	revision := tableRows(t, out, "revision")[0]
	head := rawGitText(t, root, "rev-parse", "HEAD")
	behindAhead := strings.Fields(rawGitText(t, root, "rev-list", "--left-right", "--count", "main..."+head))
	if len(behindAhead) != 2 || revision[2] != behindAhead[1] || revision[3] != behindAhead[0] || revision[4] != base || revision[6] != head {
		t.Fatalf("revision row %v disagrees with raw Git head=%s divergence=%v base=%s", revision, head, behindAhead, base)
	}

	statusRaw := rawGitBytes(t, root, false, "status", "--porcelain=v1", "-z", "--no-renames", "--untracked-files=all")
	statusEntries := independentPorcelain(statusRaw)
	tracked := independentNameStatus(rawGitBytes(t, root, false, "diff", "--name-status", "--no-renames", "-z", base))
	staged, unstaged, untracked := 0, 0, 0
	expectedCheckout := make([][]string, 0, len(statusEntries))
	for _, entry := range statusEntries {
		x, y := entry[0][:1], entry[0][1:]
		if entry[0] == "??" {
			untracked++
		} else {
			if x != " " {
				staged++
			}
			if y != " " {
				unstaged++
			}
		}
		if x == " " {
			x = "-"
		}
		if y == " " {
			y = "-"
		}
		expectedCheckout = append(expectedCheckout, []string{x, y, entry[1]})
	}
	if got := tableRows(t, out, "checkout"); fmt.Sprint(got) != fmt.Sprint(expectedCheckout) {
		t.Fatalf("checkout rows %v disagree with raw porcelain %v", got, expectedCheckout)
	}

	commits, err := strconv.Atoi(rawGitText(t, root, "rev-list", "--count", base+".."+head))
	if err != nil {
		t.Fatal(err)
	}
	insertions, deletions := independentShortstat(rawGitText(t, root, "diff", "--shortstat", base))
	aggregate := tableRows(t, out, "aggregate")[0]
	wantAggregate := []string{strconv.Itoa(commits), strconv.Itoa(len(tracked) + untracked), strconv.Itoa(insertions), strconv.Itoa(deletions), strconv.Itoa(staged), strconv.Itoa(unstaged), strconv.Itoa(untracked)}
	if fmt.Sprint(aggregate) != fmt.Sprint(wantAggregate) {
		t.Fatalf("aggregate %v disagrees with raw Git %v", aggregate, wantAggregate)
	}

	check := exec.Command("git", "-C", root, "diff", "--check", base)
	checkOut, checkErr := check.Output()
	clean := checkErr == nil
	offenses := 0
	if len(checkOut) > 0 {
		offenses = len(bytes.Split(bytes.TrimSuffix(checkOut, []byte("\n")), []byte("\n")))
	}
	wantWhitespace := []string{strconv.FormatBool(clean), strconv.Itoa(offenses)}
	if got := tableRows(t, out, "whitespace")[0]; fmt.Sprint(got) != fmt.Sprint(wantWhitespace) {
		t.Fatalf("whitespace %v disagrees with raw Git %v", got, wantWhitespace)
	}
}

func tableRows(t *testing.T, out, name string) [][]string {
	t.Helper()
	lines := strings.Split(out, "\n")
	prefix := name + "["
	for i, line := range lines {
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		end := strings.Index(line, "]")
		if end == -1 {
			t.Fatalf("malformed %s header: %q", name, line)
		}
		count, err := strconv.Atoi(line[len(prefix):end])
		if err != nil || i+count >= len(lines) {
			t.Fatalf("malformed %s row count: %q", name, line)
		}
		rows := make([][]string, 0, count)
		for _, row := range lines[i+1 : i+1+count] {
			parsed, err := csv.NewReader(strings.NewReader(strings.TrimPrefix(row, "  "))).Read()
			if err != nil {
				t.Fatalf("parse %s row %q: %v", name, row, err)
			}
			rows = append(rows, parsed)
		}
		return rows
	}
	t.Fatalf("missing %s table:\n%s", name, out)
	return nil
}

func rawGitText(t *testing.T, root string, args ...string) string {
	t.Helper()
	return strings.TrimRight(string(rawGitBytes(t, root, false, args...)), "\n")
}

func rawGitBytes(t *testing.T, root string, permitOne bool, args ...string) []byte {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.Output()
	if err == nil {
		return out
	}
	if permitOne {
		var exit *exec.ExitError
		if errorsAs(err, &exit) && exit.ExitCode() == 1 {
			return out
		}
	}
	t.Fatalf("git -C %s %v: %v\n%s", root, args, err, out)
	return nil
}

func errorsAs(err error, target any) bool {
	return errors.As(err, target)
}

func independentPorcelain(raw []byte) [][]string {
	var rows [][]string
	for _, record := range bytes.Split(raw, []byte{0}) {
		if len(record) >= 4 {
			rows = append(rows, []string{string(record[:2]), string(record[3:])})
		}
	}
	return rows
}

func independentNameStatus(raw []byte) [][]string {
	parts := bytes.Split(raw, []byte{0})
	var rows [][]string
	for i := 0; i+1 < len(parts); i += 2 {
		if len(parts[i]) > 0 {
			rows = append(rows, []string{string(parts[i]), string(parts[i+1])})
		}
	}
	return rows
}

func independentShortstat(raw string) (insertions, deletions int) {
	for _, field := range strings.Split(raw, ",") {
		words := strings.Fields(field)
		if len(words) < 2 {
			continue
		}
		n, _ := strconv.Atoi(words[0])
		if strings.HasPrefix(words[1], "insertion") {
			insertions = n
		}
		if strings.HasPrefix(words[1], "deletion") {
			deletions = n
		}
	}
	return
}
