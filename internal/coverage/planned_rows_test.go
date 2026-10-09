// Tests for the planned-row projection: the row IDs whose seam cell holds the planned
// marker, read through the one marker expression.
package coverage

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
)

// plannedSpec is a spec whose map has the header hdr and one row for each line of rows.
func plannedSpec(hdr string, rows ...string) []byte {
	return []byte("# s\n\n" + stories + "\n### Acceptance coverage map\n" + hdr + strings.Join(rows, "\n") + "\n")
}

func TestPlannedRowsReadsTheMarkerOwner(t *testing.T) {
	spec := plannedSpec(hdrReducedID,
		"| E1 | 1 | b | planned TestA in internal/x/a_test.go | w |",
		"| E2 | 1 | b | `internal/x/b_test.go` (`TestB`) | w |",
		"| E3 | 1 | b | the test is planned for a later ticket | w |",
		"| E4 | 1 | b | `internal/x/c_test.go` (`TestC`) | a planned test in the why cell |",
		"| E5 | 1 | b | review-owned: the planned canary stays green | w |",
	)
	if got, want := PlannedRows(spec), []string{"E1", "E3"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("PlannedRows = %q, want %q", got, want)
	}
}

func TestPlannedRowsKeepsTheWordBoundary(t *testing.T) {
	spec := plannedSpec(hdrReducedID,
		"| E1 | 1 | b | unplanned TestA | w |",
		"| E2 | 1 | b | planned TestB | w |",
	)
	if got, want := PlannedRows(spec), []string{"E2"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("PlannedRows = %q, want %q", got, want)
	}
}

func TestPlannedRowsSkipsANonOptInMap(t *testing.T) {
	spec := plannedSpec(hdrReduced, "| 1 | b | planned TestA in internal/x/a_test.go | w |")
	if got := PlannedRows(spec); len(got) != 0 {
		t.Fatalf("PlannedRows = %q, want no row ID for a map without row IDs", got)
	}
}

// TestPlannedMarkerHasOneOwner walks every non-test Go file that `go list ./...` names in
// this module, and finds the text of the marker expression only in this package. The walk
// matches the expression text, not the word, so a constant that holds the word passes.
func TestPlannedMarkerHasOneOwner(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "list", gate.DisableBuildVCS, "-json", "./...")
	cmd.Dir = root
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("go list: %v: %s", err, stderr.String())
	}
	needle := []byte(plannedRe.String())
	owners := map[string]bool{}
	files := 0
	decoder := json.NewDecoder(bytes.NewReader(out))
	for {
		var pkg struct {
			Dir                               string
			GoFiles, CgoFiles, IgnoredGoFiles []string
		}
		if err := decoder.Decode(&pkg); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			t.Fatalf("decode go list: %v", err)
		}
		for _, name := range append(append(pkg.GoFiles, pkg.CgoFiles...), pkg.IgnoredGoFiles...) {
			if strings.HasSuffix(name, "_test.go") {
				continue
			}
			files++
			data, err := os.ReadFile(filepath.Join(pkg.Dir, name))
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Contains(data, needle) {
				rel, err := filepath.Rel(root, pkg.Dir)
				if err != nil {
					t.Fatal(err)
				}
				owners[filepath.ToSlash(rel)] = true
			}
		}
	}
	if want := map[string]bool{"internal/coverage": true}; files == 0 || !reflect.DeepEqual(owners, want) {
		t.Fatalf("marker expression %q found in %v over %d files, want only internal/coverage", needle, owners, files)
	}
}
