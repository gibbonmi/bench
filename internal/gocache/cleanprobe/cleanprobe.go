// Package cleanprobe owns the second-process probe the cache-lock holder rows drive.
// A POSIX record lock is owned per process, so a clean inside the holder's own process
// could never contend with it. A holder row therefore re-executes its test binary with
// the probe row selected, and the child runs `bench cache clean` and records the verb's
// answer.
//
// Answer and Require are the two halves of one wire format. They stay in this package so
// the writer and the graders cannot drift apart. Both the gate suite and the testreport
// suite import them.
package cleanprobe

import (
	"os"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/gocache"
)

// Env names the file the probe writes its answer to. A child with no entry is inert, so
// the probe row is a no-op inside an ordinary suite run.
const Env = "BENCH_TEST_CACHE_CLEAN_PROBE"

// TB is the subset of testing.TB the probe helpers use. This interface, not testing.TB
// itself, keeps the testing package out of every binary that links this one.
type TB interface {
	Helper()
	Fatal(args ...any)
	Fatalf(format string, args ...any)
}

// Answer runs the probe body. With no answer entry in the environment it does nothing.
// A test file's TestCacheCleanProbe row calls this in one line.
func Answer(t TB) {
	t.Helper()
	answerPath := os.Getenv(Env)
	if answerPath == "" {
		return
	}
	answer, code := gocache.Command([]string{"clean"})
	if err := os.WriteFile(answerPath, []byte(strconv.Itoa(code)+"\n"+answer), 0o600); err != nil {
		t.Fatal(err)
	}
}

// Require grades a probe answer as the refusal a live holder produces. An unheld lock
// lets the clean through, and that is what a missing hold looks like here.
func Require(t TB, answer string) {
	t.Helper()
	code, rest, _ := strings.Cut(answer, "\n")
	if code != "1" || !strings.HasPrefix(rest, "error: cache in use — ") {
		t.Fatalf("clean beside the run = exit %s, %q; want the cache-in-use refusal at exit 1", code, rest)
	}
}
