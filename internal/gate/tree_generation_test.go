package gate

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type recordingTreeSource struct {
	treeSource
	trees, listings, reads int
	err                    error
}

func (s *recordingTreeSource) tree() (string, error) { s.trees++; return s.treeSource.tree() }
func (s *recordingTreeSource) list(tree string) (string, error) {
	s.listings++
	return s.treeSource.list(tree)
}

func (s *recordingTreeSource) blob(object string) ([]byte, error) {
	s.reads++
	if s.err != nil {
		return nil, s.err
	}
	return s.treeSource.blob(object)
}

func TestTreeGenerationOneTreeAndListingPerGeneration(t *testing.T) {
	root, tree := treeGenerationFixture(t)
	source := &recordingTreeSource{treeSource: prospectiveTreeSource{root: root, treeID: tree}}
	generation, err := captureTreeGeneration(source)
	if err != nil {
		t.Fatal(err)
	}
	if generation.tree != tree || source.trees != 1 || source.listings != 1 {
		t.Fatalf("generation source work = tree:%q trees:%d listings:%d, want supplied tree and one source/listing", generation.tree, source.trees, source.listings)
	}
}

func TestTreeGenerationSourceFaultDoesNotProduceGeneration(t *testing.T) {
	root, tree := treeGenerationFixture(t)
	for _, tc := range []struct {
		name   string
		source treeSource
	}{
		{"unavailable tree", unavailableTreeSource{}},
		{"listing", failingTreeSource{treeSource: prospectiveTreeSource{root: root, treeID: tree}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if generation, err := captureTreeGeneration(tc.source); err == nil || generation != nil {
				t.Fatalf("capture = (%#v, %v), want refusal", generation, err)
			}
		})
	}
}

func TestTreeGenerationRejectsMalformedListing(t *testing.T) {
	root, tree := treeGenerationFixture(t)
	source := malformedListingTreeSource{treeSource: prospectiveTreeSource{root: root, treeID: tree}}
	if generation, err := captureTreeGeneration(source); err == nil || generation != nil {
		t.Fatalf("malformed listing capture = (%#v, %v), want refusal without a generation", generation, err)
	}
}

type malformedListingTreeSource struct{ treeSource }

func (malformedListingTreeSource) list(string) (string, error) {
	return "100644 blob 0123456789abcdef path-without-tab\x00", nil
}

func TestTreeGenerationRejectsNonBlobRequest(t *testing.T) {
	root, tree := treeGenerationFixture(t)
	source := &recordingTreeSource{treeSource: prospectiveTreeSource{root: root, treeID: tree}}
	generation, err := captureTreeGeneration(source)
	if err != nil {
		t.Fatal(err)
	}
	if data, err := generation.blob(treeEntry{Path: "nested", Metadata: "040000 tree " + tree}); err == nil || data != nil {
		t.Fatalf("non-blob request = (%q, %v), want refusal", data, err)
	}
	if source.reads != 0 {
		t.Fatalf("non-blob source reads = %d, want zero", source.reads)
	}
}

type unavailableTreeSource struct{}

func (unavailableTreeSource) tree() (string, error)       { return "", errBlobUnavailable }
func (unavailableTreeSource) list(string) (string, error) { return "", nil }
func (unavailableTreeSource) blob(string) ([]byte, error) { return nil, nil }

type failingTreeSource struct{ treeSource }

func (failingTreeSource) list(string) (string, error) { return "", errBlobUnavailable }

func TestTreeGenerationProspectiveSourceReadsSuppliedTree(t *testing.T) {
	root, tree := treeGenerationFixture(t)
	if err := os.WriteFile(filepath.Join(root, "shared"), []byte("moved\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "ambient-only"), []byte("ambient\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	working, err := captureWorkingTree(root)
	if err != nil {
		t.Fatal(err)
	}
	prospective, err := captureProspectiveTree(root, tree)
	if err != nil {
		t.Fatal(err)
	}
	prospectiveEntry, prospectiveOK := prospective.entry("shared")
	workingEntry, workingOK := working.entry("shared")
	if !prospectiveOK || !workingOK {
		t.Fatalf("shared entries = prospective:%t working:%t, want both", prospectiveOK, workingOK)
	}
	prospectiveBlob, err := prospective.blob(prospectiveEntry)
	if err != nil {
		t.Fatal(err)
	}
	workingBlob, err := working.blob(workingEntry)
	if err != nil {
		t.Fatal(err)
	}
	if string(prospectiveBlob) != "shared\n" || string(workingBlob) != "moved\n" {
		t.Fatalf("shared blobs = prospective:%q working:%q, want supplied and moved trees", prospectiveBlob, workingBlob)
	}
	if _, found := prospective.entry("ambient-only"); found {
		t.Fatal("prospective source included ambient-only path")
	}
	if _, found := working.entry("ambient-only"); !found {
		t.Fatal("working source omitted ambient-only path")
	}
}

func TestTreeGenerationMemoizesSharedBlobByObject(t *testing.T) {
	root, tree := treeGenerationFixture(t)
	source := &recordingTreeSource{treeSource: prospectiveTreeSource{root: root, treeID: tree}}
	generation, err := captureTreeGeneration(source)
	if err != nil {
		t.Fatal(err)
	}
	first, ok := generation.entry("shared")
	if !ok {
		t.Fatal("shared entry is absent")
	}
	second, ok := generation.entry("nested/same")
	if !ok {
		t.Fatal("nested shared entry is absent")
	}
	for _, entry := range []treeEntry{first, second, first} {
		if _, err := generation.blob(entry); err != nil {
			t.Fatal(err)
		}
	}
	firstView, err := generation.blob(first)
	if err != nil {
		t.Fatal(err)
	}
	firstView[0] = 'X'
	secondView, err := generation.blob(second)
	if err != nil {
		t.Fatal(err)
	}
	if string(secondView) != "shared\n" {
		t.Fatalf("second blob view = %q, want immutable shared content", secondView)
	}
	if source.reads != 1 {
		t.Fatalf("blob reads = %d, want 1 for one object", source.reads)
	}
}

func TestTreeGenerationMemoizesBlobFailure(t *testing.T) {
	root, tree := treeGenerationFixture(t)
	source := &recordingTreeSource{treeSource: prospectiveTreeSource{root: root, treeID: tree}, err: errBlobUnavailable}
	generation, err := captureTreeGeneration(source)
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := generation.entry("shared")
	if !ok {
		t.Fatal("shared entry is absent")
	}
	first, err := generation.blob(entry)
	if !errors.Is(err, errBlobUnavailable) || first != nil {
		t.Fatalf("first blob result = (%q, %v), want unavailable failure", first, err)
	}
	second, err := generation.blob(entry)
	if !errors.Is(err, errBlobUnavailable) || second != nil {
		t.Fatalf("second blob result = (%q, %v), want unavailable failure", second, err)
	}
	if source.reads != 1 {
		t.Fatalf("blob reads = %d, want 1 after a cached failure", source.reads)
	}
}

func TestTreeGenerationCapturesDoNotShareCacheState(t *testing.T) {
	root, tree := treeGenerationFixture(t)
	first, err := captureProspectiveTree(root, tree)
	if err != nil {
		t.Fatal(err)
	}
	second, err := captureProspectiveTree(root, tree)
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatal("captures returned one generation instance")
	}
	entry, ok := first.entry("shared")
	if !ok {
		t.Fatal("shared entry is absent")
	}
	if _, err := first.blob(entry); err != nil {
		t.Fatal(err)
	}
	if len(second.blobs) != 0 {
		t.Fatalf("second generation cache = %#v, want empty independent cache", second.blobs)
	}
}

func treeGenerationFixture(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	if err := os.MkdirAll(filepath.Join(root, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "shared"), []byte("shared\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(filepath.Join(root, "shared"), filepath.Join(root, "nested", "same")); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-qm", "tree")
	return root, gitOutput(t, root, "rev-parse", "HEAD^{tree}")
}

var errBlobUnavailable = errors.New("blob unavailable")
