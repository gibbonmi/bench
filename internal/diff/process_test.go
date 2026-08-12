package diff

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

func TestProductionDiffProcessHelper(t *testing.T) {
	if os.Getenv("BENCH_DIFF_PROCESS_HELPER") != "1" {
		return
	}
	separator := 0
	for i, arg := range os.Args {
		if arg == "--" {
			separator = i + 1
			break
		}
	}
	out, code := Command(os.Args[separator:])
	_, _ = os.Stdout.WriteString(out)
	os.Exit(code)
}

func runProductionDiff(t *testing.T, root string, args ...string) (string, int) {
	t.Helper()
	argv := append([]string{"-test.run=^TestProductionDiffProcessHelper$", "--"}, args...)
	cmd := exec.Command(os.Args[0], argv...)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "BENCH_DIFF_PROCESS_HELPER=1")
	out, err := cmd.CombinedOutput()
	if err == nil {
		return string(out), 0
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return string(out), exit.ExitCode()
	}
	t.Fatalf("run production diff: %v", err)
	return "", -1
}
