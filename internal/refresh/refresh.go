package refresh

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
)

type RefreshResult struct {
	Status string
	Detail string
}

var refreshTimeout = bounds.GitRefreshTimeout

func Refresh(root string) RefreshResult {
	if bounds.Offline() {
		return RefreshResult{Status: "offline", Detail: "BENCH_OFFLINE=1"}
	}
	cmd := exec.Command("git", "-C", root, "fetch", "-q", "--no-recurse-submodules", "origin")
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	result := bounds.Run(context.Background(), refreshTimeout, cmd)
	if result.Status == bounds.ProcessComplete {
		return RefreshResult{Status: "refreshed", Detail: "none"}
	}
	detail := strings.TrimSpace(string(result.Output))
	if detail == "" && result.Err != nil {
		detail = result.Err.Error()
	}
	if result.Status == bounds.ProcessTimeout {
		detail = fmt.Sprintf("timeout after %s: %s", bounds.GitRefreshTimeout, detail)
	}
	if detail == "" {
		detail = "git fetch failed"
	}
	return RefreshResult{Status: "failed", Detail: sanitize.Preview(detail)}
}

func RefreshedStartRef(root string) string {
	remote := git.RemoteDefaultRef(root)
	if remote != "" && git.OK("-C", root, "rev-parse", "--verify", "--quiet", remote+"^{commit}") {
		return remote
	}
	if local, ok := git.ResolvedDefault(root); ok {
		return local
	}
	return "HEAD"
}

func RenderRefresh(result RefreshResult) string {
	out, err := toon.Table("worktree_refresh", []string{"status", "detail"}, [][]string{{result.Status, result.Detail}})
	if err != nil {
		return "worktree_refresh[1]{status,detail}:\n  failed,unrepresentable refresh detail\n"
	}
	return out
}

// Start is the one entry point that owns the refreshed-start-ref rule. An unrequested
// refresh writes nothing and answers an empty ref. A requested one runs the fetch, renders
// the table, and answers the new start ref only when the fetch succeeded. Every other
// status — offline, failed, timed out — answers empty, so a caller never starts from a
// remote head it did not fetch. Each caller maps the empty answer to its own default.
func Start(root string, requested bool, stdout io.Writer) string {
	if !requested {
		return ""
	}
	result := Refresh(root)
	_, _ = io.WriteString(stdout, RenderRefresh(result))
	if result.Status != "refreshed" {
		return ""
	}
	return RefreshedStartRef(root)
}

func Consume(root string, args []string, stdout io.Writer) ([]string, string) {
	filtered := args[:0]
	requested := false
	for _, arg := range args {
		if arg == "--refresh" {
			requested = true
		} else {
			filtered = append(filtered, arg)
		}
	}
	return filtered, Start(root, requested, stdout)
}
