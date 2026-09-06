package intent

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent/ledger"
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
	const leafImport = `"github.com/gibbonmi/bench/internal/intent/ledger"`
	// The admission policy is the intent owner's own child. The lifecycle policy is a
	// pure package whose own census forbids internal/intent, so the schema types it
	// reads must come from the leaf directly.
	owners := map[string]bool{
		filepath.Join("internal", "intent"):                      true,
		filepath.Join("internal", "intent", "admissionpolicy"):   true,
		filepath.Join("internal", "worktree", "lifecyclepolicy"): true,
	}
	root := filepath.Join("..", "..")
	var offenders []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if name := entry.Name(); name == ".git" || name == "dist" {
				return fs.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(entry.Name(), ".go") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !strings.Contains(string(data), leafImport) {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if owners[filepath.Dir(rel)] {
			return nil
		}
		offenders = append(offenders, rel)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(offenders) != 0 {
		t.Errorf("these files import the leaf ledger outside the intent owner: %v", offenders)
	}
}
