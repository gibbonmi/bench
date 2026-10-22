//go:build system

package systemtest

// This file holds the cases of the system censuses. A case plants one test file and
// requires the exact reports of a census over it. Each census test also runs its census
// on the live package.

import (
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/testcensus"
)

// requireCensus runs census over one planted test file, whose body starts on line 5,
// and requires exactly want.
func requireCensus(t *testing.T, census func(string) ([]string, error), body string, want []string) {
	t.Helper()
	dir := testcensus.Plant(t, map[string]string{"planted_test.go": "package systemtest\n\nimport \"testing\"\n\n" + body})
	requireReports(t, census, dir, want)
}

// requireReports runs census over dir and requires exactly want.
func requireReports(t *testing.T, census func(string) ([]string, error), dir string, want []string) {
	t.Helper()
	reports, err := census(dir)
	if err != nil {
		t.Fatalf("census %s: %v", dir, err)
	}
	if !slices.Equal(reports, want) {
		t.Fatalf("census %s =\n%s\nwant\n%s", dir, strings.Join(reports, "\n"), strings.Join(want, "\n"))
	}
}

// TestSystemParallelCensus proves the census refuses a test with no t.Parallel() call and
// no serial reason, and a parallel test that reaches the owner's repositories or names a
// ledger field outside an owner writer, and that the live package passes. (Coverage rows
// LGP25, LGP29.)
func TestSystemParallelCensus(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, body string
		serial     map[string]string
		want       []string
	}{
		{"no call and no reason", "func TestQuiet(t *testing.T) {}\n", nil,
			[]string{"planted_test.go:5: TestQuiet calls no t.Parallel() and holds no serial reason"}},
		{"serial reason", "func TestQuiet(t *testing.T) {}\n", map[string]string{"TestQuiet": pendingMigration}, nil},
		{"owner repository", "func TestShare(t *testing.T) { t.Parallel(); _ = owner.repos[0] }\n", nil,
			[]string{"planted_test.go:5: TestShare calls t.Parallel() and reaches owner.repos"}},
		{"owner repository through a method", "func TestShare(t *testing.T) { t.Parallel(); first() }\n\nfunc first() string { return owner.pick() }\n\nfunc (o *systemOwner) pick() string { return o.repos[0] }\n", nil,
			[]string{"planted_test.go:5: TestShare calls t.Parallel() and reaches o.repos through first"}},
		{"ledger field", "func TestBump(t *testing.T) { t.Parallel(); owner.starts++ }\n", nil,
			[]string{"planted_test.go:5: TestBump calls t.Parallel() and names owner.starts"}},
		{"ledger field through a helper", "func TestPeek(t *testing.T) { t.Parallel(); peek() }\n\nfunc peek() int { return len(owner.seen) }\n", nil,
			[]string{"planted_test.go:5: TestPeek calls t.Parallel() and names owner.seen through peek"}},
		{"owner writer", "func TestRun(t *testing.T) { t.Parallel(); owner.bump() }\n\nfunc (o *systemOwner) bump() { o.mu.Lock(); o.starts++; o.mu.Unlock() }\n", nil, nil},
		{"ledger field through a helper that an owner writer also runs", "func TestMix(t *testing.T) { t.Parallel(); owner.bump(); peek() }\n\nfunc (o *systemOwner) bump() { o.mu.Lock(); o.starts++; peek(); o.mu.Unlock() }\n\nfunc peek() int { return len(owner.seen) }\n", nil,
			[]string{"planted_test.go:5: TestMix calls t.Parallel() and names owner.seen through peek"}},
		{"owner repository through an owner writer", "func TestStore(t *testing.T) { t.Parallel(); owner.store() }\n\nfunc (o *systemOwner) store() { o.mu.Lock(); o.starts++; _ = o.repos[0]; o.mu.Unlock() }\n", nil,
			[]string{"planted_test.go:5: TestStore calls t.Parallel() and reaches o.repos through store"}},
		{"owner repository through a method value", "func TestPick(t *testing.T) { t.Parallel(); f := owner.pick; f() }\n\nfunc (o *systemOwner) pick() string { return o.repos[0] }\n", nil,
			[]string{"planted_test.go:5: TestPick calls t.Parallel() and reaches o.repos through pick"}},
		{"owner repository through a function value", "func TestFirst(t *testing.T) { t.Parallel(); f := first; f() }\n\nfunc first() string { return owner.repos[0] }\n", nil,
			[]string{"planted_test.go:5: TestFirst calls t.Parallel() and reaches owner.repos through first"}},
		{"owner repository through a package-level function literal", "func TestLiteral(t *testing.T) { t.Parallel(); _ = pickRepo() }\n\nvar pickRepo = func() string { return owner.repos[0] }\n", nil,
			[]string{"planted_test.go:5: TestLiteral calls t.Parallel() and reaches owner.repos through pickRepo"}},
		{"owner repository through a non-owner method", "func TestHold(t *testing.T) { t.Parallel(); var h holder; h.first() }\n\ntype holder struct{}\n\nfunc (holder) first() string { return owner.repos[0] }\n", nil,
			[]string{"planted_test.go:5: TestHold calls t.Parallel() and reaches owner.repos through first"}},
		{"local names", "func TestLocal(t *testing.T) { t.Parallel(); seen := map[string]bool{\"starts\": true}; seen[\"terminal\"] = true }\n", nil, nil},
	} {
		t.Run(row.name, func(t *testing.T) {
			requireCensus(t, func(dir string) ([]string, error) { return systemParallelCensus(dir, row.serial) }, row.body, row.want)
		})
	}
	requireReports(t, func(dir string) ([]string, error) { return systemParallelCensus(dir, systemSerialTests) }, ".", nil)
}

// TestLedgerWritesHoldTheOwnerLock proves every ledger write outside lockExempt holds the
// owner mutex: an unguarded write reds by name, and the live package passes. (Coverage
// row LGP28.)
func TestLedgerWritesHoldTheOwnerLock(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, body string
		want       []string
	}{
		{"unguarded", "func (o *systemOwner) bump() { o.starts++ }\n",
			[]string{"planted_test.go:5: bump writes o.starts outside the owner mutex"}},
		{"after the unlock", "func record(x int) {\n\towner.mu.Lock()\n\towner.mu.Unlock()\n\towner.seen = append(owner.seen, x)\n}\n",
			[]string{"planted_test.go:8: record writes owner.seen outside the owner mutex"}},
		{"guarded", "func (o *systemOwner) mark(k string) {\n\to.mu.Lock()\n\tdefer o.mu.Unlock()\n\to.terminal[k] = true\n}\n", nil},
		{"exempt", "func TestOwnerSelectionChild(t *testing.T) { owner.starts = 0 }\n", nil},
	} {
		t.Run(row.name, func(t *testing.T) {
			requireCensus(t, ledgerLockCensus, row.body, row.want)
		})
	}
	requireReports(t, ledgerLockCensus, ".", nil)
}

// TestSystemSerialTableIsExact proves a serial row that names a missing test or a
// parallel test, or that holds no reason, reds, and that the live table is exact.
// (Coverage row LGP31.)
func TestSystemSerialTableIsExact(t *testing.T) {
	t.Parallel()
	table := map[string]string{"TestSerial": "", "TestGone": pendingMigration, "TestMain": pendingMigration, "TestQuick": pendingMigration}
	requireCensus(t, func(dir string) ([]string, error) { return serialTableCensus(dir, table) },
		"func TestMain(m *testing.M) {}\n\nfunc TestSerial(t *testing.T) {}\n\nfunc TestQuick(t *testing.T) { t.Parallel() }\n",
		[]string{
			"serial row TestGone names no top-level test",
			"serial row TestMain names no top-level test",
			"serial row TestQuick names a test that calls t.Parallel()",
			"serial row TestSerial holds no reason",
		})
	requireReports(t, func(dir string) ([]string, error) { return serialTableCensus(dir, systemSerialTests) }, ".", nil)
}
