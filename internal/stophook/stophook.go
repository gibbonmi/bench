// Package stophook is the ported Stop hook — the completion oracle that decides
// whether an armed shift may stop. It exists because a stop that blocks a green
// gate or allows a red one is the worst class of bug: message wording, exit codes,
// and the no-forged-verdict cache guarantee are load-bearing, so they live in one
// unit-tested place rather than scattered through .bench/hooks/stop.sh.
//
// The pure seams (Active, Tail, BlockMessage) are table-tested; Run does the I/O
// (exec the gate, resolve the git dir, write the verdict cache) and is exercised
// end-to-end by the shell gate contracts. The verdict cache is only ever written
// with a real tree hash — a "none" hash is refused loudly rather than forged, so
// the cache can never claim a verdict it cannot key to a tree.
package stophook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/git"
)

// treeHashRE is the shape a real git tree hash must match before it is written to
// the verdict cache. Anything else (notably the literal "none" that internal/git
// returns on failure) is refused, so recordGate never forges a tree.
var treeHashRE = regexp.MustCompile(`^[0-9a-f]+$`)

// blockHeader is the three-line preamble of the BLOCKED message, ending in a
// newline so the gate-output tail follows on its own lines. One source for the
// exact wording the shell hook emitted.
const blockHeader = "BLOCKED: the gate is red, so this shift is not done.\n" +
	"Fix the failing checks at the seam — do not weaken or skip a check —\n" +
	"then stop again. Gate output:\n"

// Active reports whether the Stop envelope's stop_hook_active field renders
// Python-style as "True" — JSON boolean true or the JSON string "True". Claude
// Code sets it when it re-invokes the hook after a prior block, so honoring it
// breaks the stop loop. Everything else (false, absent, a number, another string,
// invalid JSON) is not active, so a malformed envelope fails toward enforcement.
func Active(stdin []byte) bool {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(stdin, &m); err != nil {
		return false
	}
	raw, ok := m["stop_hook_active"]
	if !ok {
		return false
	}
	var b bool
	if err := json.Unmarshal(raw, &b); err == nil {
		return b
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s == "True"
	}
	return false
}

// Tail returns the last n newline-separated lines of output rejoined with "\n",
// mirroring `tail -n 30`. A single trailing newline is not counted as an extra
// blank line, so the tail of "a\nb\n" is "a\nb" rather than a spurious blank. When
// output has n or fewer lines it is returned unchanged.
func Tail(output string, n int) string {
	lines := strings.Split(output, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) <= n {
		return strings.Join(lines, "\n")
	}
	return strings.Join(lines[len(lines)-n:], "\n")
}

// BlockMessage is the refusal returned to the agent when the gate is red: the
// fixed three-line header followed by the last 30 lines of the gate output. One
// source for the message the hook writes to stderr.
func BlockMessage(gateOutput string) string {
	return blockHeader + Tail(gateOutput, 30)
}

// Run is the Stop-hook orchestration. When the shift is not armed, or the envelope
// says the hook is already active, it allows the stop (exit 0) without touching the
// cache. Otherwise it runs `<wrapper> gate`, records the verdict, and returns 0 on
// green or 2 (after writing the BLOCKED message to stderr) on red.
func Run(stdin []byte, wrapper string, armed bool, stderr io.Writer) int {
	if !armed {
		return 0
	}
	if Active(stdin) {
		return 0
	}

	var buf bytes.Buffer
	cmd := exec.Command(wrapper, "gate")
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	output := buf.String()

	rc := 0
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) && ee.ExitCode() != 0 {
			rc = ee.ExitCode()
		} else {
			rc = 1 // failed to run, or a zero-coded error: treat as red
		}
	}

	status := "green"
	if rc != 0 {
		status = "red"
	}
	recordGate(status, stderr)

	if rc == 0 {
		return 0
	}
	fmt.Fprint(stderr, BlockMessage(output)+"\n")
	return 2
}

// recordGate writes "<status> <tree hash> <iso8601>\n" to <git-dir>/bench-last-gate.
// It resolves the repo root and git dir through internal/git and computes the tree
// hash the same way the gate cache key does. If the tree hash is not a real hash it
// writes nothing and warns — the no-forged-verdict guarantee. Cache resolution or a
// failed write is best-effort: it degrades to no cache rather than blocking the stop.
func recordGate(status string, stderr io.Writer) {
	root, err := git.Root()
	if err != nil {
		return
	}
	gitdir, err := git.Output("-C", root, "rev-parse", "--absolute-git-dir")
	if err != nil {
		return
	}
	tree := git.TreeHash(root)
	if !treeHashRE.MatchString(tree) {
		fmt.Fprintln(stderr, "WARNING: bench tree-hash unavailable — not recording a gate verdict (no forged tree).")
		return
	}
	line := status + " " + tree + " " + time.Now().UTC().Format("2006-01-02T15:04:05Z") + "\n"
	_ = os.WriteFile(filepath.Join(gitdir, "bench-last-gate"), []byte(line), 0644)
}
