package evidencecmd

import (
	"fmt"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
)

// Export writes every verified source of one published artifact into dir and prints one
// exported line. A relative dir resolves against the working directory. write copies each
// file's bytes; a real export passes chargeevidence.WriteAll.
func Export(root, identity, dir string, write chargeevidence.ExportWrite) (string, int) {
	path, err := filepath.Abs(dir)
	if err != nil {
		return toon.Errorf("export directory unavailable", err.Error()) + "\n", 1
	}
	// The response line prints the path raw, so a control byte in it would forge a line.
	if !sanitize.LineSafe(path) {
		return storeRefusal(&chargeevidence.Refusal{Class: chargeevidence.RefuseExport,
			Detail: boundedOperand("export directory with a control byte", path)}), 1
	}
	artifact, refusal, code := OpenEvidence(root, identity)
	if refusal != "" {
		return refusal, code
	}
	defer artifact.Close()
	exported, err := artifact.Export(path, write)
	if err != nil {
		return storeRefusal(err), 1
	}
	return fmt.Sprintf("exported{sources=%d,bytes=%d,dir=%s}\n", exported.Sources, exported.Bytes, exported.Dir), 0
}
