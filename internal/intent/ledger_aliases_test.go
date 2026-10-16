package intent

import (
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent/ledger"
	"github.com/gibbonmi/bench/internal/visibletree"
)

// TestIntentReExportsEveryMovedLedgerName proves every name the leaf package declares
// still resolves through its `intent.` spelling and denotes the leaf declaration. A
// dropped alias reds compilation here, and a renamed or re-declared type reds the type
// identity check. (Coverage row LS11.)
func TestIntentReExportsEveryMovedLedgerName(t *testing.T) {
	types := []struct {
		name       string
		alias, own reflect.Type
	}{
		{"Kind", reflect.TypeOf(Kind("")), reflect.TypeOf(ledger.Kind(""))},
		{"Entry", reflect.TypeOf(Entry{}), reflect.TypeOf(ledger.Entry{})},
		{"Ledger", reflect.TypeOf(Ledger{}), reflect.TypeOf(ledger.Ledger{})},
		{"AssignmentState", reflect.TypeOf(AssignmentState("")), reflect.TypeOf(ledger.AssignmentState(""))},
		{"Recovery", reflect.TypeOf(Recovery{}), reflect.TypeOf(ledger.Recovery{})},
		{"Assignment", reflect.TypeOf(Assignment{}), reflect.TypeOf(ledger.Assignment{})},
		{"CleanupReceipt", reflect.TypeOf(CleanupReceipt{}), reflect.TypeOf(ledger.CleanupReceipt{})},
		{"LandingLease", reflect.TypeOf(LandingLease{}), reflect.TypeOf(ledger.LandingLease{})},
		{"ChunkGap", reflect.TypeOf(ChunkGap{}), reflect.TypeOf(ledger.ChunkGap{})},
		{"LandingOverride", reflect.TypeOf(LandingOverride{}), reflect.TypeOf(ledger.LandingOverride{})},
	}
	for _, row := range types {
		if row.alias != row.own {
			t.Errorf("intent.%s is %v, and the leaf declares %v: the re-export is not an alias", row.name, row.alias, row.own)
		}
	}

	values := []struct {
		name       string
		alias, own any
	}{
		{"LegacySchema", LegacySchema, ledger.LegacySchema},
		{"Schema", Schema, ledger.Schema},
		{"KindShift", KindShift, ledger.KindShift},
		{"KindWorktree", KindWorktree, ledger.KindWorktree},
		{"KindClaudeAgent", KindClaudeAgent, ledger.KindClaudeAgent},
		{"AssignmentRecordSchema", AssignmentRecordSchema, ledger.AssignmentRecordSchema},
		{"RecoveryRefNamespace", RecoveryRefNamespace, ledger.RecoveryRefNamespace},
		{"StateActive", StateActive, ledger.StateActive},
		{"StateCleanupPending", StateCleanupPending, ledger.StateCleanupPending},
		{"StateRecovered", StateRecovered, ledger.StateRecovered},
		{"StateComplete", StateComplete, ledger.StateComplete},
		{"CleanupReceiptSchema", CleanupReceiptSchema, ledger.CleanupReceiptSchema},
		{"ReceiptInFlight", ReceiptInFlight, ledger.ReceiptInFlight},
		{"ReceiptComplete", ReceiptComplete, ledger.ReceiptComplete},
		{"MaxCleanupReceipts", MaxCleanupReceipts, ledger.MaxCleanupReceipts},
		{"ReceiptPhasePlanned", ReceiptPhasePlanned, ledger.ReceiptPhasePlanned},
		{"ReceiptPhasePreserved", ReceiptPhasePreserved, ledger.ReceiptPhasePreserved},
		{"ReceiptPhaseRemoving", ReceiptPhaseRemoving, ledger.ReceiptPhaseRemoving},
		{"ReceiptPhaseRemoved", ReceiptPhaseRemoved, ledger.ReceiptPhaseRemoved},
		{"ReceiptPhaseBranch", ReceiptPhaseBranch, ledger.ReceiptPhaseBranch},
		{"ReceiptPhaseTerminal", ReceiptPhaseTerminal, ledger.ReceiptPhaseTerminal},
	}
	for _, row := range values {
		if row.alias != row.own {
			t.Errorf("intent.%s is %v, and the leaf declares %v", row.name, row.alias, row.own)
		}
	}

	functions := []struct {
		name       string
		alias, own any
	}{
		{"AssignmentBranchPrefix", AssignmentBranchPrefix, ledger.AssignmentBranchPrefix},
		{"AssignmentBranchRef", AssignmentBranchRef, ledger.AssignmentBranchRef},
		{"RecoveryRefPrefix", RecoveryRefPrefix, ledger.RecoveryRefPrefix},
		{"RequestDigest", RequestDigest, ledger.RequestDigest},
		{"ValidIdentity", ValidIdentity, ledger.ValidIdentity},
		{"ValidateAssignment", ValidateAssignment, ledger.ValidateAssignment},
		{"validAssignmentBranchRef", validAssignmentBranchRef, ledger.ValidAssignmentBranchRef},
		{"validEntry", validEntry, ledger.ValidateEntry},
		{"validateCleanupReceipts", validateCleanupReceipts, ledger.ValidateCleanupReceipts},
	}
	for _, row := range functions {
		if reflect.ValueOf(row.alias).Pointer() != reflect.ValueOf(row.own).Pointer() {
			t.Errorf("intent.%s does not name the leaf declaration", row.name)
		}
	}
}

// TestOnlyTheIntentOwnerImportsTheLeafLedger walks the module sources and refuses an
// import of the leaf path outside the intent owner and its admission-policy child. A
// build that rewrites a consumer's import to the leaf adds an edge this test names.
// (Coverage row LS12.)
func TestOnlyTheIntentOwnerImportsTheLeafLedger(t *testing.T) {
	offenders, err := leafLedgerOffenders(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if len(offenders) != 0 {
		t.Errorf("these files import the leaf ledger outside the intent owner: %v", offenders)
	}
}

// TestLeafLedgerCensusGradesOnlyTheVisibleTree is the bite proof for the census scope. A
// consumer import under an ignored path is a local copy no commit carries, so it is not an
// offender. The same import at a visible untracked path is one.
func TestLeafLedgerCensusGradesOnlyTheVisibleTree(t *testing.T) {
	root := gittest.RepoIgnoring(t, ".logs/")
	consumer := []byte("package consumer\n\nimport _ \"github.com/gibbonmi/bench/internal/intent/ledger\"\n")
	write := func(rel string) {
		t.Helper()
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, consumer, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write(".logs/copy/internal/consumer/consumer.go")
	if offenders, err := leafLedgerOffenders(root); err != nil || len(offenders) != 0 {
		t.Fatalf("ignored consumer: offenders %v, %v, want none", offenders, err)
	}
	write("internal/consumer/consumer.go")
	offenders, err := leafLedgerOffenders(root)
	if err != nil || !reflect.DeepEqual(offenders, []string{"internal/consumer/consumer.go"}) {
		t.Fatalf("visible consumer: offenders %v, %v, want only internal/consumer/consumer.go", offenders, err)
	}
}

// leafLedgerOffenders lists the visible Go files under root that import the leaf ledger
// outside the intent owner.
func leafLedgerOffenders(root string) ([]string, error) {
	const leafImport = `"github.com/gibbonmi/bench/internal/intent/ledger"`
	// The admission policy is the intent owner's own child. The lifecycle policy is a
	// pure package whose own census forbids internal/intent, so the schema types it
	// reads must come from the leaf directly.
	owners := map[string]bool{
		"internal/intent":                   true,
		"internal/intent/admissionpolicy":   true,
		"internal/worktree/lifecyclepolicy": true,
	}
	visible, err := visibletree.Files(root)
	if err != nil {
		return nil, err
	}
	var offenders []string
	for _, rel := range visible {
		if !strings.HasSuffix(rel, ".go") || strings.HasPrefix(rel, "dist/") || owners[path.Dir(rel)] {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			return nil, err
		}
		if strings.Contains(string(data), leafImport) {
			offenders = append(offenders, rel)
		}
	}
	return offenders, nil
}
