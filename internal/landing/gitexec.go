package landing

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

func updateRef(root, ref, new, old string) error { return run(root, "update-ref", ref, new, old) }
func destinationUpdateFailure(root, ref, expected string, updateErr error) error {
	actual, err := output(root, "rev-parse", "--verify", ref+"^{commit}")
	if err != nil {
		return fmt.Errorf("read destination after failed ref update: %w", err)
	}
	if actual != expected {
		return fmt.Errorf("destination compare-and-swap refused; rerun the landing to recompose onto the moved destination: %w", updateErr)
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
func outputInput(root string, input []byte, args ...string) (string, error) {
	c := exec.Command("git", append([]string{"-C", root}, args...)...)
	c.Stdin = strings.NewReader(string(input))
	b, err := c.Output()
	return strings.TrimSpace(string(b)), err
}
func run(root string, args ...string) error {
	return exec.Command("git", append([]string{"-C", root}, args...)...).Run()
}
func indexRun(root, idx string, args ...string) error {
	c := exec.Command("git", append([]string{"-C", root}, args...)...)
	c.Env = append(os.Environ(), "GIT_INDEX_FILE="+idx)
	return c.Run()
}
func indexOutputRaw(root, idx string, args ...string) ([]byte, error) {
	c := exec.Command("git", append([]string{"-C", root}, args...)...)
	c.Env = append(os.Environ(), "GIT_INDEX_FILE="+idx)
	return c.Output()
}
func indexOutput(root, idx string, args ...string) (string, error) {
	c := exec.Command("git", append([]string{"-C", root}, args...)...)
	c.Env = append(os.Environ(), "GIT_INDEX_FILE="+idx)
	b, err := c.Output()
	return strings.TrimSpace(string(b)), err
}
