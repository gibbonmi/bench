package worktree

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
)

func consumeRefresh(root string, args []string, stdout io.Writer) []string {
	filtered := args[:0]
	requested := false
	for _, arg := range args {
		if arg == "--refresh" {
			requested = true
		} else {
			filtered = append(filtered, arg)
		}
	}
	if requested {
		_, _ = io.WriteString(stdout, RenderRefresh(Refresh(root)))
	}
	return filtered
}

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
		detail = "timeout after 30s: " + detail
	}
	if detail == "" {
		detail = "git fetch failed"
	}
	return RefreshResult{Status: "failed", Detail: sanitize.Preview(detail)}
}

func RenderRefresh(result RefreshResult) string {
	out, err := toon.Table("worktree_refresh", []string{"status", "detail"}, [][]string{{result.Status, result.Detail}})
	if err != nil {
		return "worktree_refresh[1]{status,detail}:\n  failed,unrepresentable refresh detail\n"
	}
	return out
}
