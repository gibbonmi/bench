package registry

import (
	"path"
	"slices"
	"strings"
)

// RouteBypassTest is RR49, the ordinary test that grades each write-verb literal against the
// refusal route registry.
const RouteBypassTest = "TestNoWriteVerbComposesARouteOutsideTheRegistry"

// RouteBypassPackages are the packages in which a write verb raises or prints a refusal.
// The scan reads each one's own directory and no subdirectory, so a package that this list
// omits goes unscanned.
var RouteBypassPackages = []string{
	"internal/commit",
	"internal/commitment",
	"internal/commitment/commitcmd",
	"internal/commitment/repository",
	"internal/gate",
	"internal/gate/authorization",
	"internal/landing",
	"internal/reviewrecord/recordcmd",
	"internal/worktree",
	"internal/worktree/landingpolicy",
}

// RouteBypassScanned reports whether the RR49 scan reads the slash-separated file rel: a
// production Go file directly inside a scanned package.
func RouteBypassScanned(rel string) bool {
	return path.Ext(rel) == ".go" && !strings.HasSuffix(rel, "_test.go") && slices.Contains(RouteBypassPackages, path.Dir(rel))
}
