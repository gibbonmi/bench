package subprocess_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/subprocess"
)

// replaceProbeEnv names the file where this test binary, as a replacing process, has a shell
// write the shell's PID.
const replaceProbeEnv = "SUBPROCESS_REPLACE_PROBE"

// init runs before TestMain, so a replacing process runs no test.
func init() {
	out, ok := os.LookupEnv(replaceProbeEnv)
	if !ok {
		return
	}
	err := subprocess.Replace("sh", []string{"-c", `printf %s "$$" > "$0"`, out}, os.Environ())
	fmt.Fprintln(os.Stderr, err)
	os.Exit(3)
}

// TestReplaceKeepsTheProcess requires that the program runs as the process that Replace
// replaces, not as a child of it.
func TestReplaceKeepsTheProcess(t *testing.T) {
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "pid")
	cmd := exec.Command(self)
	cmd.Env = append(os.Environ(), replaceProbeEnv+"="+out)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("replacing process = %v; output:\n%s", err, output)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(data); got != strconv.Itoa(cmd.Process.Pid) {
		t.Fatalf("program PID = %s, want the replaced process %d", got, cmd.Process.Pid)
	}
}

// TestReplaceReportsAProgramThatCannotStart requires that a failed replacement returns to its
// caller with an error that names the program: one that is missing, and an executable file
// that holds no program format.
func TestReplaceReportsAProgramThatCannotStart(t *testing.T) {
	dir := t.TempDir()
	formatless := filepath.Join(dir, "formatless")
	if err := os.WriteFile(formatless, []byte("no program format\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, program := range []string{filepath.Join(dir, "missing"), formatless} {
		if err := subprocess.Replace(program, nil, os.Environ()); err == nil || !strings.Contains(err.Error(), program) {
			t.Errorf("replacement with %s = %v, want an error that names it", program, err)
		}
	}
}
