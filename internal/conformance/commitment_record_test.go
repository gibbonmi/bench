package conformance

import (
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/intent"
)

// DATA_HANDLING.md lists every commitment field of the intent ledger between these markers.
const (
	commitmentRecordBegin = "<!-- commitment-record:begin -->"
	commitmentRecordEnd   = "<!-- commitment-record:end -->"
)

// commitmentFieldRe matches the dotted field path in the first cell of a listing row.
var commitmentFieldRe = regexp.MustCompile("(?m)^\\|\\s*`([a-z_]+(?:\\.[a-z_]+)*)`\\s*\\|")

// checkCommitmentRecordInventory asserts that DATA_HANDLING.md lists exactly the leaf fields
// of the intent ledger's commitment records, in both directions. The field set comes from
// the ledger types, so a new commitment field turns the gate red until the inventory names
// its contents. A missing or empty listing fails loudly.
func checkCommitmentRecordInventory(root string) []string {
	region, ok := markedRegion(readIfExists(filepath.Join(root, "DATA_HANDLING.md")), commitmentRecordBegin, commitmentRecordEnd)
	if !ok {
		return []string{fmt.Sprintf("DATA_HANDLING.md commitment record region missing: expected the field listing between %s and %s", commitmentRecordBegin, commitmentRecordEnd)}
	}
	documented := map[string]bool{}
	for _, match := range commitmentFieldRe.FindAllStringSubmatch(region, -1) {
		documented[match[1]] = true
	}
	if len(documented) == 0 {
		return []string{"DATA_HANDLING.md commitment record region empty: no `field` rows between the commitment record markers"}
	}
	fields := commitmentRecordFields()
	var diags []string
	for _, field := range fields {
		if !documented[field] {
			diags = append(diags, fmt.Sprintf("DATA_HANDLING.md commitment record: ledger field %q is not documented", field))
		}
	}
	for field := range documented {
		if !slices.Contains(fields, field) {
			diags = append(diags, fmt.Sprintf("DATA_HANDLING.md commitment record: documented field %q is not a ledger field", field))
		}
	}
	sort.Strings(diags)
	return diags
}

// commitmentRecordFields returns the dotted JSON path of each leaf field below every intent
// ledger field that holds a commitment record: a plan receipt, a milestone receipt, or the
// runtime commitment state.
func commitmentRecordFields() []string {
	records := map[reflect.Type]bool{
		reflect.TypeOf(intent.CommitmentReceipt{}): true,
		reflect.TypeOf(intent.MilestoneReceipt{}):  true,
		reflect.TypeOf(intent.CommitmentState{}):   true,
	}
	var fields []string
	ledger := reflect.TypeOf(intent.Ledger{})
	for i := range ledger.NumField() {
		if field := ledger.Field(i); records[recordElement(field.Type)] {
			fields = appendLeafFields(fields, jsonFieldName(field), recordElement(field.Type))
		}
	}
	return fields
}

// appendLeafFields appends the path of each leaf JSON field of a struct type below prefix.
func appendLeafFields(fields []string, prefix string, record reflect.Type) []string {
	for i := range record.NumField() {
		field := record.Field(i)
		path := prefix + "." + jsonFieldName(field)
		if element := recordElement(field.Type); element.Kind() == reflect.Struct {
			fields = appendLeafFields(fields, path, element)
			continue
		}
		fields = append(fields, path)
	}
	return fields
}

// recordElement strips the slice and pointer layers that wrap a record type.
func recordElement(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Slice || t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

func jsonFieldName(field reflect.StructField) string {
	name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
	return name
}
