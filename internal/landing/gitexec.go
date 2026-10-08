package landing

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

func updateRef(root, ref, new, old string) error { return run(root, "update-ref", ref, new, old) }

// TipMovedError is the compare-and-swap refusal of a destination that moved between the
// read and the update. A rerun recomposes onto the moved destination, so a caller routes
// the refusal by the type and never by the sentence.
type TipMovedError struct{ Err error }

func (e TipMovedError) Error() string { return e.Err.Error() }
func (e TipMovedError) Unwrap() error { return e.Err }

func destinationUpdateFailure(root, ref, expected string, updateErr error) error {
	actual, err := benchgit.ResolveCommit(root, ref)
	if err != nil {
		return fmt.Errorf("read destination after failed ref update: %w", err)
	}
	if actual != expected {
		return TipMovedError{fmt.Errorf("destination compare-and-swap refused; rerun the landing to recompose onto the moved destination: %w", updateErr)}
	}
	return updateErr
}
func output(root string, args ...string) (string, error) {
	return benchgit.Output(append([]string{"-C", root}, args...)...)
}
func outputCombined(root string, args ...string) (string, error) {
	c := exec.Command("git", append([]string{"-C", root}, args...)...)
	b, err := c.CombinedOutput()
	return strings.TrimSpace(string(b)), err
}
func outputRaw(root string, args ...string) ([]byte, error) {
	return exec.Command("git", append([]string{"-C", root}, args...)...).Output()
}
func blobContent(root, oid string) ([]byte, error) {
	content, err := outputRaw(root, "cat-file", "blob", oid)
	if err != nil {
		return nil, fmt.Errorf("read blob %s: %w", oid, err)
	}
	return content, nil
}
func hashBlob(root string, content []byte) (string, error) {
	c := exec.Command("git", "-C", root, "hash-object", "-w", "--no-filters", "--stdin")
	c.Stdin = bytes.NewReader(content)
	b, err := c.Output()
	return strings.TrimSpace(string(b)), err
}
func run(root string, args ...string) error {
	return exec.Command("git", append([]string{"-C", root}, args...)...).Run()
}
func indexRun(root, idx string, args ...string) error {
	return benchgit.IndexCommand(root, idx, args...).Run()
}
func indexOutput(root, idx string, args ...string) (string, error) {
	b, err := benchgit.IndexCommand(root, idx, args...).Output()
	return strings.TrimSpace(string(b)), err
}
