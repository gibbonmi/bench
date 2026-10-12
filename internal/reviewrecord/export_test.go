package reviewrecord

import (
	"io/fs"
	"testing"
)

// The external replacement test drives the real write transaction with these.
var (
	Write        = write
	Absent       = absent
	RenderRecord = renderRecord
)

const RenderSpec = renderSpec

// BindReplacement binds replacement as the replacement operation of every review-record
// file write for the test t, and the cleanup of t restores the former binding.
func BindReplacement(t testing.TB, replacement func(path string, data []byte, mode fs.FileMode) error) {
	t.Helper()
	original := replaceFile
	replaceFile = replacement
	t.Cleanup(func() { replaceFile = original })
}
