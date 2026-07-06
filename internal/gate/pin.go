package gate

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/terminal"
)

const pinFileName = "bench-gate-pin"

// PinCommand is the human-attended `bench gate pin` porcelain. It refuses non-TTY
// stdin before doing any write, then records HEAD's committed .bench tree beside the
// gate cache for the managed pre-push hook to verify.
func PinCommand(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	return pinCommand(args, stdin, stdout, stderr, terminal.IsTerminal)
}

func pinCommand(args []string, stdin io.Reader, stdout, stderr io.Writer, isTerminal func(io.Reader) bool) int {
	if len(args) != 0 {
		fmt.Fprintln(stderr, "usage: bench gate pin")
		return 2
	}
	if !isTerminal(stdin) {
		fmt.Fprintln(stderr, "error: bench gate pin requires an interactive TTY")
		return 1
	}
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, "not in a git repo")
		return 1
	}
	if !showPinReview(root, stdout, stderr) {
		return 1
	}
	fmt.Fprint(stdout, "Type 'pin .bench' to update the gate pin: ")
	reader := bufio.NewReader(stdin)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		fmt.Fprintln(stderr, "error: could not read confirmation")
		return 1
	}
	if strings.TrimSpace(line) != "pin .bench" {
		fmt.Fprintln(stderr, "bench gate pin declined; no pin written")
		return 1
	}
	if err := writePinFromHead(root); err != nil {
		fmt.Fprintln(stderr, "error: "+err.Error())
		return 1
	}
	fmt.Fprintln(stdout, "bench gate pin updated")
	return 0
}

func showPinReview(root string, stdout, stderr io.Writer) bool {
	if dirtyBench(root) {
		fmt.Fprintln(stderr, "warning: .bench has uncommitted changes; pinning HEAD's committed .bench tree")
	}
	commit := existingPinnedCommit(root)
	if commit == "" {
		fmt.Fprintln(stdout, "Initial gate pin for HEAD:.bench")
		return true
	}
	fmt.Fprintf(stdout, "Diff since pinned gate commit %s:\n", commit)
	cmd := exec.Command("git", "-C", root, "diff", commit+"..HEAD", "--", ".bench")
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run() == nil
}

func dirtyBench(root string) bool {
	out, err := git.Output("-C", root, "status", "--porcelain", "--", ".bench")
	return err == nil && out != ""
}

func existingPinnedCommit(root string) string {
	data, err := os.ReadFile(pinPath(root))
	if err != nil {
		return ""
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) < 2 {
		return ""
	}
	return strings.TrimSpace(lines[1])
}

func writePinFromHead(root string) error {
	tree, err := git.Output("-C", root, "rev-parse", "HEAD:.bench")
	if err != nil || tree == "" {
		return fmt.Errorf("cannot resolve HEAD:.bench")
	}
	commit, err := git.Output("-C", root, "rev-parse", "HEAD")
	if err != nil || commit == "" {
		return fmt.Errorf("cannot resolve HEAD")
	}
	content := fmt.Sprintf("%s\n%s\n%s\n", tree, commit, time.Now().UTC().Format(time.RFC3339))
	return os.WriteFile(pinPath(root), []byte(content), 0o644)
}

func pinPath(root string) string {
	gitdir, err := git.Output("-C", root, "rev-parse", "--absolute-git-dir")
	if err != nil || gitdir == "" {
		return filepath.Join(root, ".git", pinFileName)
	}
	return filepath.Join(gitdir, pinFileName)
}
