package adopt

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/adopt/transaction"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/git"
)

// WriteBrokerManifest publishes the broker manifest beside the resolved wrapper,
// binding the currently running executable as the promotion broker. It returns the
// manifest path and the broker path it bound. The install and repair owner calls it, so
// a repository executable can never become the landing owner.
func WriteBrokerManifest(version string) (string, string, error) {
	return writeBrokerManifest(version, nil)
}

func writeBrokerManifest(version string, guard *adoptionGuard) (string, string, error) {
	broker, err := os.Executable()
	if err != nil {
		return "", "", fmt.Errorf("resolve promotion broker executable: %w", err)
	}
	dest := filepath.Join(filepath.Dir(resolvedWrapper()), brokermanifest.Name)
	change, err := observedChange(dest, "", guard)
	if err != nil {
		return "", "", err
	}
	stage, err := os.MkdirTemp("", "bench-broker-stage-")
	if err != nil {
		return "", "", err
	}
	defer os.RemoveAll(stage)
	path, bound, err := brokermanifest.Write(stage, broker, version)
	if err != nil {
		return "", "", err
	}
	change.Stage = path
	if err := publishChange(change, guard); err != nil {
		return "", "", err
	}
	return dest, bound, nil
}

func lockDoctor(dir string) (*adoptionGuard, error) {
	paths := []string{filepath.Join(dir, "bench"), filepath.Join(filepath.Dir(resolvedWrapper()), brokermanifest.Name)}
	if root, err := git.Root(); err == nil {
		hooks, err := hooksDir(root)
		if err != nil {
			return nil, err
		}
		paths = append(paths, filepath.Join(hooks, "pre-push"))
	}
	return lockObserved(paths)
}

func beginDoctorRepair(guard *adoptionGuard, repair *repairRun, dir string) error {
	for path, observation := range guard.observations {
		if observation.err != nil {
			return observation.err
		}
		file := bounds.ClassifyNoFollow(path)
		if file.State != bounds.StateParsed && file.State != bounds.StateAbsent {
			return fmt.Errorf("refused doctor asset: %s", path)
		}
	}
	if err := validateRepairShim(filepath.Join(dir, "bench")); err != nil {
		return err
	}

	stage, err := os.MkdirTemp("", "bench-doctor-repair-")
	if err != nil {
		return err
	}
	repair.staging = stage
	repair.shim = guard.observations[filepath.Join(dir, "bench")].change.Destination
	guard.repair = repair
	return nil
}

func finishDoctorRepair(guard *adoptionGuard, code int, stderr io.Writer) int {
	repair := guard.repair
	defer os.RemoveAll(repair.staging)
	if code != 0 {
		return code
	}
	repair.store.Lease = guard.Lease
	id, err := repair.store.Apply(repair.pending)
	repair.id = id
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func retainDoctorChange(change transaction.Change, repair *repairRun) error {
	if err := change.Recheck(); err != nil {
		return err
	}
	staged := bounds.ClassifyNoFollow(change.Stage)
	info, err := os.Lstat(change.Stage)
	if err != nil {
		return err
	}
	if staged.State != bounds.StateParsed {
		return fmt.Errorf("invalid staged doctor asset: %s", change.Stage)
	}
	current := bounds.ClassifyNoFollow(change.Destination)
	if current.State != bounds.StateAbsent {
		before, err := os.Lstat(change.Destination)
		if err != nil {
			return err
		}
		canonical := change.Destination == repair.shim || bytes.Equal(current.Data, staged.Data)
		if !canonical && filepath.Base(change.Destination) == brokermanifest.Name {
			canonical = canonicalBroker(change.Destination, current.Data, repair.staging)
		}
		if current.State != bounds.StateParsed || before.Mode() != info.Mode() || !canonical {
			return fmt.Errorf("modified-managed doctor asset: %s", change.Destination)
		}
	}
	path, err := stageBytes(repair.staging, fmt.Sprintf("asset-%d", len(repair.pending)), staged.Data, info.Mode())
	if err != nil {
		return err
	}
	change.Stage = path
	repair.pending = append(repair.pending, change)
	return nil
}

const shimMode os.FileMode = 0o755

func managedShimContent(content string) bool { return strings.Contains(content, shimMarkerID) }

func validateRepairShim(path string) error {
	file := bounds.ClassifyNoFollow(path)
	if file.State == bounds.StateAbsent {
		return nil
	}
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	content := string(file.Data)
	target := ShimTarget(content)
	if file.State == bounds.StateParsed && target != "" && content == ShimContent(target)+"\n" && info.Mode() == shimMode {
		return nil
	}
	reason := "project-owned"
	if managedShimContent(content) {
		reason = "modified-managed"
	}
	return fmt.Errorf("%s shim: %s", reason, path)
}

func canonicalBroker(path string, content []byte, stage string) bool {
	fields, err := brokermanifest.Read(path)
	if err != nil {
		return false
	}
	dir, err := os.MkdirTemp(stage, "broker-preimage-")
	if err != nil {
		return false
	}
	defer os.RemoveAll(dir)
	canonical, _, err := brokermanifest.Write(dir, fields["path"], fields["version"])
	if err != nil {
		return false
	}
	data, err := os.ReadFile(canonical)
	return err == nil && bytes.Equal(data, content)
}
