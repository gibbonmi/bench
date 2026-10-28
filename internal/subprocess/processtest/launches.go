package processtest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/subprocess"
)

// The environment of a helper stand-in: the record of the launches and the release file.
const (
	launchRecordEnv  = "BENCH_PROCESSTEST_LAUNCHES"
	launchReleaseEnv = "BENCH_PROCESSTEST_RELEASE"
)

// launchPoll is the poll interval of a stand-in that waits for its release.
var launchPoll = bounds.FixedWindow(10 * time.Millisecond)

// launchPolls is the bounded count of polls of a stand-in, so a stand-in that no test releases
// still ends.
func launchPolls() int { return int(bounds.TestDeadline(0) / launchPoll) }

// Launches records the runs of a stand-in program that the code under test starts by its path.
// Each run appends its pid and its arguments to the record, and it then runs until the test
// releases the stand-ins. StandIn answers a shell stand-in. A helper of this test binary is a
// stand-in too, when it runs under Env and calls RecordLaunch and then AwaitRelease.
type Launches struct {
	record, release string
}

// Launch is one run of a stand-in: its pid and its arguments.
type Launch struct {
	PID  int
	Args []string
}

// NewLaunches prepares the record of one test. Its cleanup ends each stand-in that started.
func NewLaunches(t *testing.T) *Launches {
	t.Helper()
	dir := t.TempDir()
	l := &Launches{record: filepath.Join(dir, "launches"), release: filepath.Join(dir, "release")}
	t.Cleanup(func() { l.End(t) })
	return l
}

// StandIn answers a shell stand-in. It writes its record line in one append, with a tab before
// each argument, so an argument with a space stays one argument.
func (l *Launches) StandIn() string {
	return "#!/bin/sh\ntab=$(printf '\\t')\nline=$$\nfor arg in \"$@\"; do line=$line$tab$arg; done\n" +
		"printf '%s\\n' \"$line\" >> " + sanitize.ShellQuote(l.record) + "\ni=0\n" +
		"while [ ! -e " + sanitize.ShellQuote(l.release) + " ] && [ \"$i\" -lt " + strconv.Itoa(launchPolls()) + " ]; do sleep " +
		strconv.FormatFloat(launchPoll.Seconds(), 'f', -1, 64) + "; i=$((i+1)); done\n"
}

// Env answers the environment of a helper stand-in.
func (l *Launches) Env() []string {
	return []string{launchRecordEnv + "=" + l.record, launchReleaseEnv + "=" + l.release}
}

// RecordLaunch appends the pid and the arguments of this helper run to the record that its
// environment names.
func RecordLaunch() error {
	record, err := os.OpenFile(os.Getenv(launchRecordEnv), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(record, strings.Join(append([]string{strconv.Itoa(os.Getpid())}, os.Args[1:]...), "\t"))
	return errors.Join(err, record.Close())
}

// AwaitRelease waits until the test releases the stand-ins, for the bounded count of polls.
func AwaitRelease() {
	for i := 0; i < launchPolls(); i++ {
		if _, err := os.Stat(os.Getenv(launchReleaseEnv)); err == nil {
			return
		}
		time.Sleep(launchPoll)
	}
}

// Launched answers each run that started, in order. A line that a run is still appending has no
// final newline yet, so it is not a run yet.
func (l *Launches) Launched(t *testing.T) []Launch {
	t.Helper()
	data, err := os.ReadFile(l.record)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(data), "\n")
	var launches []Launch
	for _, line := range lines[:len(lines)-1] {
		fields := strings.Split(line, "\t")
		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			t.Fatalf("launch record line %q names no pid", line)
		}
		launches = append(launches, Launch{PID: pid, Args: fields[1:]})
	}
	return launches
}

// Await waits for the first run to start and answers it.
func (l *Launches) Await(t *testing.T) Launch {
	t.Helper()
	WaitFor(t, "a stand-in to start", func() bool { return len(l.Launched(t)) > 0 })
	return l.Launched(t)[0]
}

// End releases the stand-ins and waits until each run that started has exited.
func (l *Launches) End(t *testing.T) {
	t.Helper()
	if err := os.WriteFile(l.release, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	WaitFor(t, "each stand-in to exit", func() bool {
		for _, launch := range l.Launched(t) {
			if alive, err := subprocess.Alive(launch.PID); err != nil || alive {
				return false
			}
		}
		return true
	})
}
