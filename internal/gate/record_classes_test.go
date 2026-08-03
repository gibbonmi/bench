package gate

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// [PC10a-family] No two record classes in the store carry the same field set, and no class
// outside the verdict family borrows a verdict field beyond the store-wide shared ones.
// Exact-field-set validation is what makes the classes mutually unreadable, so two classes
// declaring one set would let a record authored as either be read as the other — and the
// refusals every test above rests on would stop separating them.
func TestStoreRecordClassesStayMutuallyUnreadable(t *testing.T) {
	names := make([]string, 0, len(storeRecordClasses))
	for name := range storeRecordClasses {
		names = append(names, name)
	}
	sort.Strings(names)
	for i, name := range names {
		for _, other := range names[i+1:] {
			if reflect.DeepEqual(storeRecordClasses[name], storeRecordClasses[other]) {
				t.Errorf("the %s and %s classes declare one field set %v", name, other, storeRecordClasses[name])
			}
		}
	}

	for _, name := range names {
		if contains(verdictClasses, name) {
			continue
		}
		if shared := recordClassSharesVerdictFields(storeRecordClasses[name]); len(shared) != 0 {
			t.Errorf("the %s class shares %v with the verdict classes", name, shared)
		}
	}
	// The shared names are shared in fact, not merely permitted: a name listed as shared that
	// some class does not carry would silently widen what every other class may borrow.
	for _, name := range recordSharedFields {
		for _, class := range names {
			if !contains(storeRecordClasses[class], name) {
				t.Errorf("%q is listed as shared but the %s class does not carry it", name, class)
			}
		}
	}
}

// [PC10a-family] verdictClasses has to name one class per ready field set — full and
// partial — or the disjointness loop above silently stops comparing sibling classes
// against whichever one dropped out. This pins the count and the identity of the two
// against the *ReadyFields variables themselves, not against a restated list of names, so
// the only way to satisfy it is for verdictClasses to actually cover both.
func TestVerdictClassesCoverAllReadyFieldSets(t *testing.T) {
	want := [][]string{fullReadyFields, partialReadyFields, checkPartialReadyFields, combinedPartialReadyFields}
	if len(verdictClasses) != len(want) {
		t.Fatalf("verdictClasses = %v (%d classes), want %d — one per ready field set (full, partial)", verdictClasses, len(verdictClasses), len(want))
	}
	matched := make(map[string]bool, len(want))
	for _, name := range verdictClasses {
		fields, ok := storeRecordClasses[name]
		if !ok {
			t.Fatalf("verdictClasses names %q, which storeRecordClasses does not carry", name)
		}
		for i, w := range want {
			if reflect.DeepEqual(fields, w) {
				matched[strconv.Itoa(i)] = true
			}
		}
	}
	if len(matched) != len(want) {
		t.Fatalf("verdictClasses = %v, want its classes' field sets to be exactly {full, reduced, partial}ReadyFields", verdictClasses)
	}
}

// [PC10a-family] Every *ReadyFields variable declared in verdict.go has to be registered in
// readyFieldClasses, or a fourth verdict class can be added there without ever joining the
// enumeration this whole family depends on — exactly the gap the partial class fell into
// before this ticket. This is a source-shape check in the tradition of
// TestTreeSnapshotIsTheOnlyListingParser: it parses verdict.go's own declarations rather than
// trusting a hand-maintained list of names, so it cannot itself go stale the way a restated
// list would.
//
// What it guarantees: no package-level identifier ending in "ReadyFields" can exist in
// verdict.go without appearing as a value in the readyFieldClasses map literal. What it does
// not guarantee: that the class was given the right field set, or wired into validation —
// those are TestVerdictClassesCoverAllThreeReadyFieldSets's and verdict.go's own tests' job.
func TestVerdictReadyFieldsAreAllRegistered(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "verdict.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}

	var declared []string
	registered := map[string]bool{}

	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, name := range vs.Names {
				if strings.HasSuffix(name.Name, "ReadyFields") {
					declared = append(declared, name.Name)
				}
				if name.Name != "readyFieldClasses" || i >= len(vs.Values) {
					continue
				}
				lit, ok := vs.Values[i].(*ast.CompositeLit)
				if !ok {
					t.Fatalf("readyFieldClasses is not a composite literal this scan can read")
				}
				for _, elt := range lit.Elts {
					kv, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					if ident, ok := kv.Value.(*ast.Ident); ok {
						registered[ident.Name] = true
					}
				}
			}
		}
	}

	if len(declared) == 0 {
		t.Fatal("found no *ReadyFields variable in verdict.go — the scan itself is broken")
	}
	sort.Strings(declared)
	for _, name := range declared {
		if !registered[name] {
			t.Errorf("%s is declared in verdict.go but not registered in readyFieldClasses", name)
		}
	}
}
