package worktree

import (
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/otelrecord"
)

// TestWorktreeSeamsMatchTheRegistry covers OT20 at the registry. The verbs share one
// span symbol, so the conformance check cannot see a seam constant that no registry row
// names, or a registry row whose verb opens no span. The package's own seam constants
// and the registry are authored apart, so set equality is the reconciliation.
func TestWorktreeSeamsMatchTheRegistry(t *testing.T) {
	t.Parallel()
	instrumented := append([]string{otelLandingSeam}, otelVerbSeams...)
	var registered []string
	for _, entry := range otelrecord.Registry {
		if strings.HasPrefix(entry.Seam, "worktree.") {
			registered = append(registered, entry.Seam)
		}
	}
	sort.Strings(instrumented)
	sort.Strings(registered)
	if strings.Join(instrumented, ",") != strings.Join(registered, ",") {
		t.Fatalf("worktree seams: package %v, registry %v", instrumented, registered)
	}
}
