package prepared

import fixture "github.com/gibbonmi/bench/internal/contract/surface/artifact/internal/fixture"

import (
	"errors"
	"io/fs"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
)

type sharedArtifactSet struct {
	sourceRoot  string
	outputDir   string
	entryCount  int
	fingerprint map[string]string
}

type sharedArtifactSetState struct {
	once      sync.Once
	set       *sharedArtifactSet
	builds    int
	directory string
}

// Artifact package tests stay sequential; adding t.Parallel would race shared-set
// consumers against the mutation probes.
var packageSharedArtifactSet sharedArtifactSetState

func (state *sharedArtifactSetState) resolve() (*sharedArtifactSet, error) {
	if state.set == nil {
		return nil, errors.New("shared artifact set staging failed in an earlier test")
	}
	return state.set, nil
}

func requireSharedArtifactSet(t *testing.T) *sharedArtifactSet {
	t.Helper()
	if os.Geteuid() == 0 {
		capability.Capability(t, capability.Privilege, "root bypasses file mode write protection; shared artifact set read-only writes are unobservable")
	}
	packageSharedArtifactSet.once.Do(func() {
		packageSharedArtifactSet.stage(t)
	})
	set, err := packageSharedArtifactSet.resolve()
	if err != nil {
		t.Fatal(err)
	}
	if packageSharedArtifactSet.builds != 1 {
		t.Fatalf("shared artifact set builds = %d, want exactly one", packageSharedArtifactSet.builds)
	}
	if err := set.verify(); err != nil {
		t.Fatal(err)
	}
	return set
}

func (state *sharedArtifactSetState) stage(t *testing.T) {
	t.Helper()
	root := contract.SubjectRoot(t)
	directory, err := os.MkdirTemp("", "bench shared artifact set [*]")
	if err != nil {
		t.Fatal(err)
	}
	state.directory = directory
	source := fixture.CommittedHostileArtifactSourceIn(t, directory, root)
	output := filepath.Join(directory, "promoted artifacts [*]")
	state.build(t, root, source, output)
	fingerprint, err := fixture.PromotedArtifactDigestMap(output)
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{source, output} {
		if err := makeSharedArtifactTreeReadOnly(path); err != nil {
			t.Fatal(err)
		}
	}
	state.set = &sharedArtifactSet{
		sourceRoot:  source,
		outputDir:   output,
		entryCount:  len(fingerprint),
		fingerprint: fingerprint,
	}
}

func (state *sharedArtifactSetState) build(t *testing.T, root, source, output string) {
	t.Helper()
	state.builds++
	if err := appendSharedArtifactSetBuildLog(); err != nil {
		t.Fatal(err)
	}
	contract.NewExecFixtureAt(t, root).Run("bash", filepath.Join(source, "scripts", "build-artifacts.sh"), source, output).RequireExit(0)
}

func appendSharedArtifactSetBuildLog() error {
	if log := os.Getenv("BENCH_TEST_SHARED_SET_BUILD_LOG"); log != "" {
		file, err := os.OpenFile(log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		if _, err := file.WriteString("build\n"); err != nil {
			_ = file.Close()
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}
	}
	return nil
}

func makeSharedArtifactTreeReadOnly(root string) error {
	return chmodSharedArtifactTree(root, 0o444, 0o555)
}

func chmodSharedArtifactTree(root string, fileMode, directoryMode fs.FileMode) error {
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		mode := fileMode
		if entry.IsDir() {
			mode = directoryMode
		}
		if err := os.Chmod(path, mode); err != nil {
			return err
		}
		return nil
	})
}

func (state *sharedArtifactSetState) cleanup() error {
	if state.directory == "" {
		return nil
	}
	if err := chmodSharedArtifactTree(state.directory, 0o644, 0o755); err != nil {
		return err
	}
	return os.RemoveAll(state.directory)
}

func (set *sharedArtifactSet) verify() error {
	fingerprint, err := fixture.PromotedArtifactDigestMap(set.outputDir)
	if err != nil || !maps.Equal(fingerprint, set.fingerprint) {
		return errors.New("shared artifact set mutated")
	}
	return nil
}

func (set *sharedArtifactSet) promotedTarball() (string, error) {
	entries, err := os.ReadDir(set.outputDir)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		if entry.Type().IsRegular() && (strings.HasSuffix(entry.Name(), ".tgz") || strings.HasSuffix(entry.Name(), ".tar.gz")) {
			return filepath.Join(set.outputDir, entry.Name()), nil
		}
	}
	return "", errors.New("shared artifact set has no promoted tarball")
}

func TestSharedArtifactSetFailsClosedAfterEarlierStagingFailure(t *testing.T) {
	var state sharedArtifactSetState
	state.once.Do(func() {})
	_, err := state.resolve()
	if err == nil || err.Error() != "shared artifact set staging failed in an earlier test" {
		t.Fatalf("consumed shared artifact state error = %v, want exact earlier-staging failure", err)
	}
}

func TestSharedArtifactSetBuildIsLazy(t *testing.T) {
	log := filepath.Join(t.TempDir(), "shared artifact builds")
	command := exec.Command(os.Args[0], "-test.run=^(TestArtifactPromotionIsAtomicAndExclusive|TestArtifactSourceStagesCommittedHostPlan|TestSharedCacheBuildPromotesNoRecord|TestSharedCacheBuildRestoresRecordOnInterruptedPromotion|TestOfflineArchiveProjection|TestPackedArtifactRunsSetupOfflineFromASpacedPrefix)$", "-test.v")
	command.Env = contract.ProcessEnv(nil, map[string]string{"BENCH_TEST_SHARED_SET_BUILD_LOG": log})
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("fabricated-fixture-only subprocess failed: %v\n%s", err, output)
	}
	builds, err := os.ReadFile(log)
	if err != nil {
		t.Fatalf("read shared artifact build log: %v", err)
	}
	if got := strings.Count(string(builds), "build\n"); got != 1 {
		t.Fatalf("six prepared sharers built the shared artifact set %d times, want 1", got)
	}
}

func TestSharedArtifactSetAttributesMutation(t *testing.T) {
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	shared := requireSharedArtifactSet(t)
	path, err := shared.promotedTarball()
	if err != nil {
		t.Fatal(err)
	}
	original, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(original, 'x'), 0o644); err != nil {
		t.Fatal(err)
	}
	verifyErr := shared.verify()
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0o444); err != nil {
		t.Fatal(err)
	}
	if verifyErr == nil || verifyErr.Error() != "shared artifact set mutated" {
		t.Fatalf("mutated shared artifact set error = %v, want exact attribution", verifyErr)
	}
}

func TestSharedArtifactSetIsReadOnly(t *testing.T) {
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	shared := requireSharedArtifactSet(t)
	tarball, err := shared.promotedTarball()
	if err != nil {
		t.Fatal(err)
	}
	for name, path := range map[string]string{
		"promoted tarball": tarball,
		"staged source":    filepath.Join(shared.sourceRoot, "LICENSE"),
	} {
		t.Run(name, func(t *testing.T) {
			original, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0)
			if err == nil {
				_, writeErr := file.WriteString("write probe")
				closeErr := file.Close()
				if restoreErr := os.WriteFile(path, original, 0o644); restoreErr != nil {
					t.Fatalf("restore writable probe target: %v", restoreErr)
				}
				t.Fatalf("write into shared artifact set succeeded: write=%v close=%v", writeErr, closeErr)
			}
			if !errors.Is(err, fs.ErrPermission) {
				t.Fatalf("write into shared artifact set error = %v, want permission error", err)
			}
		})
	}
}
