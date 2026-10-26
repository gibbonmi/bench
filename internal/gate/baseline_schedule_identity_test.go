package gate

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

// scheduleRoots answers one runtime root and one identity root for each write. Each
// identity root differs only by its baseline phase manifest, which write leaves at the path
// it receives.
func scheduleRoots(t *testing.T, writes ...func(t *testing.T, manifest string)) (string, []string) {
	t.Helper()
	t.Setenv("BENCH_KIT", "")
	var identityRoots []string
	for _, write := range writes {
		identityRoot := t.TempDir()
		manifest := manifestPath(identityRoot)
		if err := os.MkdirAll(filepath.Dir(manifest), 0o755); err != nil {
			t.Fatal(err)
		}
		write(t, manifest)
		identityRoots = append(identityRoots, identityRoot)
	}
	return oracleIdentityFixture(t), identityRoots
}

// scheduleIdentities answers the grade identity of the runtime root under each identity
// root that scheduleRoots prepares, in order.
func scheduleIdentities(t *testing.T, writes ...func(t *testing.T, manifest string)) []string {
	t.Helper()
	root, identityRoots := scheduleRoots(t, writes...)
	var identities []string
	for _, identityRoot := range identityRoots {
		identities = append(identities, oracleIdentitySubjectFor(t, root, identityRoot, policyVersion, strings.Repeat("a", 40)).Identity)
	}
	return identities
}

func scheduleAbsent(*testing.T, string) {}

func scheduleBytes(body string) func(*testing.T, string) {
	return func(t *testing.T, manifest string) {
		t.Helper()
		if err := os.WriteFile(manifest, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// scheduleLink leaves a link at the manifest to a target outside the identity root that
// holds body, or to a missing target when body is empty.
func scheduleLink(body string) func(*testing.T, string) {
	return func(t *testing.T, manifest string) {
		t.Helper()
		target := filepath.Join(t.TempDir(), "phases.json")
		if body != "" {
			scheduleBytes(body)(t, target)
		}
		if err := os.Symlink(target, manifest); err != nil {
			t.Fatal(err)
		}
	}
}

// GR3, GR4: an absent manifest, an empty manifest, and two manifests with different bytes
// give four grade identities.
func TestGradeIdentityCarriesTheBaselineSchedule(t *testing.T) {
	ids := scheduleIdentities(t, scheduleAbsent, scheduleBytes(""), scheduleBytes(`{"phases":[]}`), scheduleBytes(`{"phases":[1]}`))
	names := []string{"absent", "empty", "first bytes", "second bytes"}
	for i := range ids {
		for j := i + 1; j < len(ids); j++ {
			if ids[i] == ids[j] {
				t.Fatalf("identity = %s for the %s and the %s manifest, want two identities", ids[i], names[i], names[j])
			}
		}
	}
}

// GR5: a live link enters as its target's bytes, so it gives the identity of a regular
// manifest with those bytes, and two targets with different bytes give two identities.
func TestGradeIdentityFollowsTheLoaderRead(t *testing.T) {
	ids := scheduleIdentities(t, scheduleBytes(`{"phases":[]}`), scheduleLink(`{"phases":[]}`), scheduleLink(`{"phases":[1]}`))
	if ids[0] != ids[1] {
		t.Fatalf("identity = %s for the manifest and %s for a link to its bytes, want one identity", ids[0], ids[1])
	}
	if ids[1] == ids[2] {
		t.Fatalf("identity = %s for links to two targets, want two identities", ids[1])
	}
}

// GR6: a dangling link and a FIFO at the manifest are defects that give one grade identity,
// and the subject build returns without opening the FIFO.
func TestGradeIdentityGivesADefectOneValue(t *testing.T) {
	root, identityRoots := scheduleRoots(t, scheduleLink(""), func(t *testing.T, manifest string) {
		if err := syscall.Mkfifo(manifest, 0o644); err != nil {
			t.Fatal(err)
		}
	})
	built := make(chan []subject, 1)
	go func() {
		var subjects []subject
		for _, identityRoot := range identityRoots {
			s, err := buildSubjectForTree(root, identityRoot, policyVersion, strings.Repeat("a", 40))
			if err != nil {
				s.Reason = err.Error()
			}
			subjects = append(subjects, s)
		}
		built <- subjects
	}()
	select {
	case subjects := <-built:
		for _, s := range subjects {
			if !s.Closed {
				t.Fatalf("subject open over a defect manifest: %s", s.Reason)
			}
		}
		if subjects[0].Identity != subjects[1].Identity {
			t.Fatalf("identity = %s for a dangling link and %s for a FIFO, want one identity", subjects[0].Identity, subjects[1].Identity)
		}
	case <-time.After(bounds.TestDeadline(0)):
		t.Fatal("the subject build did not return over a FIFO manifest")
	}
}
