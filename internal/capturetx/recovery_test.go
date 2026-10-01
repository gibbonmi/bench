package capturetx

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// Each test reaches a crash state the way a crash does: a real Begin publishes the
// generation, and the test then rewrites the manifest to the step the crash stopped at.
// Every verdict comes through the exported functions.

type fixture struct {
	root, common string
	ideas, notes Source
}

func newFixture(t *testing.T) fixture {
	t.Helper()
	root := gittest.Repo(t)
	common, err := benchgit.CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	return fixture{
		root: root, common: common,
		ideas: Source{Name: "capture/IDEAS.md", Path: filepath.Join(root, "capture", "IDEAS.md")},
		notes: Source{Name: "capture/learnings.md", Path: filepath.Join(root, "capture", "learnings.md")},
	}
}

func (f fixture) write(t *testing.T, source Source, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(source.Path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(source.Path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func (f fixture) live(t *testing.T, source Source) string {
	t.Helper()
	data, err := os.ReadFile(source.Path)
	if errors.Is(err, os.ErrNotExist) {
		return ""
	}
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func (f fixture) append(t *testing.T, source Source, line string) {
	t.Helper()
	if err := Append(f.root, source, func(current []byte) ([]byte, error) {
		return append(append([]byte(nil), current...), line...), nil
	}); err != nil {
		t.Fatal(err)
	}
}

// begin seals "idea one\n" and "note one\n" as the open generation.
func (f fixture) begin(t *testing.T) Bundle {
	t.Helper()
	f.write(t, f.ideas, "idea one\n")
	f.write(t, f.notes, "note one\n")
	bundle, err := Begin(f.root, []Source{f.notes, f.ideas})
	if err != nil {
		t.Fatal(err)
	}
	return bundle
}

// crash rewrites the open manifest to the state a crash left on disk.
func (f fixture) crash(t *testing.T, edit func(*manifest)) {
	t.Helper()
	m, err := readManifest(f.common)
	if err != nil {
		t.Fatal(err)
	}
	edit(&m)
	if err := writeManifest(f.common, m); err != nil {
		t.Fatal(err)
	}
}

// crashAborting leaves the state of an Abort that wrote its restore targets and its
// manifest but stopped before it restored a live document.
func (f fixture) crashAborting(t *testing.T, after map[string]string) {
	t.Helper()
	f.crash(t, func(m *manifest) {
		m.State = "aborting"
		for i := range m.Sources {
			body := after[m.Sources[i].Name]
			m.Sources[i].AfterBlob = "after-" + filepath.Base(m.Sources[i].Name)
			m.Sources[i].AfterDigest = digest([]byte(body))
			if err := os.WriteFile(filepath.Join(transactionDir(f.common), m.Sources[i].AfterBlob), []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func (f fixture) requireNoGeneration(t *testing.T) {
	t.Helper()
	bundle, present, err := Current(f.root)
	if err != nil {
		t.Fatal(err)
	}
	if present {
		t.Fatalf("generation %s is still open", bundle.ID)
	}
}

func requireRefusal(t *testing.T, err error, needle string) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), needle) {
		t.Fatalf("error = %v, want a refusal that names %q", err, needle)
	}
}

func TestCrashBeforePublishOpensNoGeneration(t *testing.T) {
	f := newFixture(t)
	f.write(t, f.ideas, "idea one\n")
	staged := filepath.Join(f.common, stagingPrefix+"crash")
	if err := os.MkdirAll(staged, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(staged, "source-0.bin"), []byte("idea one\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	f.requireNoGeneration(t)
	if got := f.live(t, f.ideas); got != "idea one\n" {
		t.Fatalf("live IDEAS = %q, want it unchanged", got)
	}
	if _, err := Begin(f.root, []Source{f.ideas}); err != nil {
		t.Fatalf("Begin after a staged crash: %v", err)
	}
	if got := f.live(t, f.ideas); got != "" {
		t.Fatalf("live IDEAS = %q, want Begin to clear it", got)
	}
	if _, err := os.Stat(staged); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stale staged directory: stat = %v, want recovery to remove it", err)
	}
}

// A crash can stop the removal that ends Commit or Abort after the manifest is gone
// and before the blobs are gone.
func TestCutShortRemovalDoesNotStopBegin(t *testing.T) {
	f := newFixture(t)
	opened := f.begin(t)
	if err := os.Remove(filepath.Join(transactionDir(f.common), manifestName)); err != nil {
		t.Fatal(err)
	}
	f.write(t, f.ideas, "idea two\n")
	bundle, err := Begin(f.root, []Source{f.ideas})
	if err != nil || bundle.ID == opened.ID || bundle.State != "sealed" {
		t.Fatalf("Begin after a cut-short removal = %+v, %v; want a fresh sealed generation", bundle, err)
	}
	sealed, _, err := Sources(f.root, []string{f.ideas.Name})
	if err != nil || string(sealed[f.ideas.Name]) != "idea two\n" {
		t.Fatalf("Sources = %q, %v; want the fresh sealed text", sealed, err)
	}
}

func TestPreparingRecoversForward(t *testing.T) {
	for _, tc := range []struct {
		name    string
		cleared bool
	}{{"no source cleared", false}, {"one source cleared", true}} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			opened := f.begin(t)
			f.write(t, f.ideas, "idea one\n")
			if !tc.cleared {
				f.write(t, f.notes, "note one\n")
			}
			f.crash(t, func(m *manifest) { m.State = "preparing" })
			bundle, present, err := Current(f.root)
			if err != nil || !present || bundle.ID != opened.ID || bundle.State != "sealed" {
				t.Fatalf("Current = %+v, %v, %v; want sealed generation %s", bundle, present, err, opened.ID)
			}
			for _, source := range []Source{f.ideas, f.notes} {
				if got := f.live(t, source); got != "" {
					t.Fatalf("live %s = %q, want recovery to clear it", source.Name, got)
				}
			}
			sealed, _, err := Sources(f.root, []string{f.ideas.Name, f.notes.Name})
			if err != nil || string(sealed[f.ideas.Name]) != "idea one\n" || string(sealed[f.notes.Name]) != "note one\n" {
				t.Fatalf("Sources = %q, %v; want the sealed text", sealed, err)
			}
		})
	}
}

func TestPreparingRefusesChangedLiveSource(t *testing.T) {
	f := newFixture(t)
	f.begin(t)
	f.write(t, f.ideas, "hand edit\n")
	f.crash(t, func(m *manifest) { m.State = "preparing" })
	_, _, err := Current(f.root)
	requireRefusal(t, err, "capture source capture/IDEAS.md changed during cutover")
	if got := f.live(t, f.ideas); got != "hand edit\n" {
		t.Fatalf("live IDEAS = %q, want the edit kept", got)
	}
}

func TestSealedGenerationStaysOpen(t *testing.T) {
	f := newFixture(t)
	opened := f.begin(t)
	f.append(t, f.ideas, "post\n")
	again, err := Begin(f.root, []Source{f.ideas})
	if err != nil || again.ID != opened.ID {
		t.Fatalf("repeated Begin = %+v, %v; want generation %s", again, err, opened.ID)
	}
	bundle, present, err := Current(f.root)
	if err != nil || !present || bundle.State != "sealed" || len(bundle.Sources) != 2 || bundle.Sources[0].Bytes != 9 {
		t.Fatalf("Current = %+v, %v, %v; want two sealed sources of 9 bytes", bundle, present, err)
	}
	sealed, present, err := Sources(f.root, []string{f.ideas.Name})
	if err != nil || !present || string(sealed[f.ideas.Name]) != "idea one\n" {
		t.Fatalf("Sources = %q, %v, %v; want the sealed IDEAS text", sealed, present, err)
	}
	if got := f.live(t, f.ideas); got != "post\n" {
		t.Fatalf("live IDEAS = %q, want the post-cut entry", got)
	}
}

func TestCommittingRetiresGeneration(t *testing.T) {
	f := newFixture(t)
	f.begin(t)
	f.append(t, f.ideas, "post\n")
	f.crash(t, func(m *manifest) { m.State = "committing" })
	f.requireNoGeneration(t)
	if got := f.live(t, f.ideas); got != "post\n" {
		t.Fatalf("live IDEAS = %q, want the post-cut entry kept", got)
	}
}

func TestAbortingRecovery(t *testing.T) {
	for _, tc := range []struct {
		name, ideasAtCrash string
	}{
		{"before the restore", "post\n"},
		{"after one restore", "idea one\npost\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			f.begin(t)
			f.append(t, f.ideas, "post\n")
			f.crashAborting(t, map[string]string{f.ideas.Name: "idea one\npost\n", f.notes.Name: "note one\n"})
			f.write(t, f.ideas, tc.ideasAtCrash)
			f.requireNoGeneration(t)
			if got := f.live(t, f.ideas); got != "idea one\npost\n" {
				t.Fatalf("live IDEAS = %q, want the sealed text once before the post-cut entry", got)
			}
			if got := f.live(t, f.notes); got != "note one\n" {
				t.Fatalf("live learnings = %q, want the sealed text restored", got)
			}
		})
	}
}

func TestAbortingRefusesChangedLiveSource(t *testing.T) {
	f := newFixture(t)
	f.begin(t)
	f.crashAborting(t, map[string]string{f.ideas.Name: "idea one\n", f.notes.Name: "note one\n"})
	f.write(t, f.ideas, "hand edit\n")
	_, _, err := Current(f.root)
	requireRefusal(t, err, "capture source capture/IDEAS.md changed during abort")
	if got := f.live(t, f.ideas); got != "hand edit\n" {
		t.Fatalf("live IDEAS = %q, want the edit kept", got)
	}
}

func TestRecoveryRefusesUnreadableManifest(t *testing.T) {
	for _, tc := range []struct {
		name, refusal string
		seed          func(t *testing.T, f fixture)
	}{
		{"unknown state", `has unknown state "paused"`, func(t *testing.T, f fixture) {
			f.crash(t, func(m *manifest) { m.State = "paused" })
		}},
		{"unsupported schema", "unsupported capture transaction manifest", func(t *testing.T, f fixture) {
			f.crash(t, func(m *manifest) { m.Schema = 2 })
		}},
		{"manifest does not parse", "parse capture transaction manifest", func(t *testing.T, f fixture) {
			if err := os.WriteFile(filepath.Join(transactionDir(f.common), manifestName), []byte("{"), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			f.begin(t)
			tc.seed(t, f)
			composed := false
			err := Append(f.root, f.ideas, func(current []byte) ([]byte, error) {
				composed = true
				return []byte("post\n"), nil
			})
			requireRefusal(t, err, tc.refusal)
			if composed || f.live(t, f.ideas) != "" {
				t.Fatalf("Append wrote %q through a refused recovery", f.live(t, f.ideas))
			}
		})
	}
}

func TestSourcesRefusesIncompleteGeneration(t *testing.T) {
	t.Run("sealed blob fails its digest", func(t *testing.T) {
		f := newFixture(t)
		f.begin(t)
		f.crash(t, func(m *manifest) {
			for _, source := range m.Sources {
				if err := os.WriteFile(filepath.Join(transactionDir(f.common), source.Blob), []byte("tampered\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
		})
		_, _, err := Sources(f.root, []string{f.ideas.Name})
		requireRefusal(t, err, "sealed capture source capture/IDEAS.md failed its digest")
	})
	t.Run("requested source is absent", func(t *testing.T) {
		f := newFixture(t)
		f.begin(t)
		_, _, err := Sources(f.root, []string{f.ideas.Name, "capture/absent.md"})
		requireRefusal(t, err, "does not contain every requested source")
	})
}

func TestCommitAndAbortRefuseGenerationThatIsNotOpen(t *testing.T) {
	f := newFixture(t)
	opened := f.begin(t)
	for _, retire := range []func(string, string) error{Commit, Abort} {
		requireRefusal(t, retire(f.root, "d-000000000000"), "capture transaction d-000000000000 is not open")
		if _, present, err := Current(f.root); err != nil || !present {
			t.Fatalf("a wrong identifier retired generation %s: %v", opened.ID, err)
		}
	}
	if err := Commit(f.root, opened.ID); err != nil {
		t.Fatal(err)
	}
	if err := Commit(f.root, opened.ID); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("second Commit = %v, want a refusal for no open generation", err)
	}
}

func TestRetireWithoutCrash(t *testing.T) {
	for _, tc := range []struct {
		name   string
		retire func(string, string) error
		want   string
	}{
		{"Commit keeps only the post-cut entry", Commit, "post\n"},
		{"Abort puts the sealed text before the post-cut entry", Abort, "idea one\npost\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			f.write(t, f.ideas, "idea one")
			opened, err := Begin(f.root, []Source{f.ideas})
			if err != nil {
				t.Fatal(err)
			}
			f.append(t, f.ideas, "post\n")
			if err := tc.retire(f.root, opened.ID); err != nil {
				t.Fatal(err)
			}
			f.requireNoGeneration(t)
			if got := f.live(t, f.ideas); got != tc.want {
				t.Fatalf("live IDEAS = %q, want %q", got, tc.want)
			}
		})
	}
}
