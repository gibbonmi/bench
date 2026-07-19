package conformance

import (
	"os"

	"github.com/gibbonmi/bench/internal/testrepo"
)

func materializeAuthenticatedReleaseProbe(source string) (string, func(), error) {
	root, err := os.MkdirTemp("", "bench-release-probe-*")
	if err != nil {
		return "", func() {}, err
	}
	cleanup := func() { _ = os.RemoveAll(root) }
	if err := testrepo.CommitWorkingTree(source, root); err != nil {
		cleanup()
		return "", func() {}, err
	}
	return root, cleanup, nil
}
