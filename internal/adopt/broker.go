package adopt

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/brokermanifest"
)

// WriteBrokerManifest publishes the broker manifest beside the resolved wrapper,
// binding the currently running executable as the promotion broker. It returns the
// manifest path and the broker path it bound. The install and repair owner calls it, so
// a repository executable can never become the landing owner.
func WriteBrokerManifest(version string) (string, string, error) {
	broker, err := os.Executable()
	if err != nil {
		return "", "", fmt.Errorf("resolve promotion broker executable: %w", err)
	}
	return brokermanifest.Write(filepath.Dir(resolvedWrapper()), broker, version)
}
