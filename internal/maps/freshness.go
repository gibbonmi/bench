package maps

import (
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

// sourcePathPrefix opens a Sources record whose locator is a repository path. A
// URL record gets no freshness row, because no commit date exists for a URL.
const sourcePathPrefix = "- Path: "

// staleRows projects one advisory row per Sources Path the map cites whose own last
// commit is newer than the index's. Freshness is advisory: a stale row never raises
// the exit code and never enters the unresolved count.
//
// An uncommitted index or asset yields no row. A file with no commit has no date to
// compare, and treating the absent date as zero would flag every fresh file.
func staleRows(root, name, indexPath, sources string) [][]any {
	indexTime, dated := lastCommitSeconds(root, indexPath)
	if !dated {
		return nil
	}
	var rows [][]any
	for _, locator := range sourcePathLocators(sources) {
		assetTime, dated := lastCommitSeconds(root, locator)
		if !dated || assetTime <= indexTime {
			continue
		}
		rows = append(rows, []any{name, locator, "source", "stale", "", locator})
	}
	return rows
}

// sourcePathLocators lists the Path locators of one Sources body, in the order the
// map wrote them. It shares sourceLocator with the validity check, so the backtick
// rule has one derivation.
func sourcePathLocators(sources string) []string {
	var locators []string
	for _, line := range nonEmptyLines(sources) {
		if !strings.HasPrefix(line, sourcePathPrefix) {
			continue
		}
		if locator := sourceLocator(strings.TrimPrefix(line, sourcePathPrefix)); locator != "" {
			locators = append(locators, locator)
		}
	}
	return locators
}

// lastCommitSeconds reads one path's last commit time in Unix seconds. dated is
// false when the path has no commit, which is the empty output git prints for an
// untracked or uncommitted file.
func lastCommitSeconds(root, path string) (seconds int64, dated bool) {
	out, err := git.Output("-C", root, "log", "-1", "--format=%ct", "--", path)
	if err != nil {
		return 0, false
	}
	seconds, err = strconv.ParseInt(strings.TrimSpace(out), 10, 64)
	if err != nil {
		return 0, false
	}
	return seconds, true
}
