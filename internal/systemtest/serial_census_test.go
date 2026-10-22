//go:build system

package systemtest

// This file owns the census of the system serial table, systemSerialTests. It reuses the
// parse and the reach walk of the parallel census.

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/testcensus"
)

// pendingMigration is the temporary reason of a test whose move into the parallel set is
// unfinished. No row may hold it, so an unfinished move reds by the test's name.
const pendingMigration = "pending migration"

// ownerRepositorySubjects names the only tests whose subject is the owner's repositories.
// The census pins them apart from the serial table, so a third row with the
// ownerRepositories reason reds.
var ownerRepositorySubjects = map[string]bool{"TestOwnerSelectionChild": true, "TestSelectedExecutableComposition": true}

// serialTableCensus reports every serial row that names no top-level test in dir, names
// a test that calls t.Parallel(), holds no reason, or holds pendingMigration. A reason of
// only space counts as no reason, and a reason matches pendingMigration in any letter case
// and with any space around it. It also reports where the ownerRepositories reason and
// the owner's repositories disagree: the reason on a test outside ownerRepositorySubjects
// or on a test that reaches no owner.repos, another reason on a serial test that reaches
// owner.repos, and a subject with no row of that reason.
func serialTableCensus(dir string, serial map[string]string) ([]string, error) {
	pkg, err := parseSystemPackage(dir)
	if err != nil {
		return nil, err
	}
	tests := map[string]systemFunc{}
	for _, fn := range pkg.all {
		if testcensus.IsTopLevelTest(fn.decl) {
			tests[fn.decl.Name.Name] = fn
		}
	}
	var reports []string
	for name, reason := range serial {
		fn, declared := tests[name]
		switch {
		case !declared:
			reports = append(reports, fmt.Sprintf("serial row %s names no top-level test", name))
		case testcensus.CallsParallelDirectly(fn.decl):
			reports = append(reports, fmt.Sprintf("serial row %s names a test that calls t.Parallel()", name))
		default:
			if repos, _ := pkg.reach(fn); repos == "" && reason == ownerRepositories {
				reports = append(reports, fmt.Sprintf("serial row %s holds the owner-repository reason and reaches no owner.repos", name))
			} else if repos != "" && reason != ownerRepositories {
				reports = append(reports, fmt.Sprintf("serial row %s reaches %s without the owner-repository reason", name, repos))
			}
		}
		switch trimmed := strings.TrimSpace(reason); {
		case trimmed == "":
			reports = append(reports, fmt.Sprintf("serial row %s holds no reason", name))
		case strings.EqualFold(trimmed, pendingMigration):
			reports = append(reports, fmt.Sprintf("serial row %s holds the temporary reason %q", name, reason))
		case reason == ownerRepositories:
			if !ownerRepositorySubjects[name] {
				reports = append(reports, fmt.Sprintf("serial row %s holds the owner-repository reason and names no owner-repository subject", name))
			}
		}
	}
	for name := range ownerRepositorySubjects {
		if serial[name] != ownerRepositories {
			reports = append(reports, fmt.Sprintf("owner-repository subject %s holds no owner-repository row", name))
		}
	}
	sort.Strings(reports)
	return reports, nil
}

// TestSystemSerialTableIsExact proves a serial row reds when it names a missing test or a
// parallel test, holds no reason or the temporary reason, or holds an owner-repository
// reason that disagrees with the owner's repositories, and that the live table is exact.
// (Coverage rows LGP31, LGP32, LGP38.)
func TestSystemSerialTableIsExact(t *testing.T) {
	t.Parallel()
	requireCensus(t, func(dir string) ([]string, error) { return serialTableCensus(dir, serialTableCase.table) },
		serialTableCase.body, serialTableCase.want)
	requireReports(t, func(dir string) ([]string, error) { return serialTableCensus(dir, systemSerialTests) }, ".", nil)
}
