package publication

import (
	"context"
	"testing"
)

// memoryRegistry is an in-process Registry: publish stores the bytes and every
// integrity query answers from them, so a complete first publication runs with
// no process, no network, and no fixture server — leaving the durable record as
// the only thing under test.
type memoryRegistry struct {
	published map[string][]byte
}

func newMemoryRegistry() *memoryRegistry {
	return &memoryRegistry{published: map[string][]byte{}}
}

func (m *memoryRegistry) Publish(ctx context.Context, name, version, tag string, tarball []byte) (string, error) {
	m.published[name+"@"+version] = tarball
	return sriIntegrity(tarball), nil
}

func (m *memoryRegistry) StageSubmit(ctx context.Context, name, version string, tarball []byte) (string, error) {
	return "", nil
}

func (m *memoryRegistry) Approve(ctx context.Context, stageID string) error { return nil }

func (m *memoryRegistry) Integrity(ctx context.Context, name, version string) (string, bool, error) {
	tarball, ok := m.published[name+"@"+version]
	if !ok {
		return "", false, nil
	}
	return sriIntegrity(tarball), true, nil
}

func (m *memoryRegistry) TagAdd(ctx context.Context, name, tag, version string) error { return nil }

func (m *memoryRegistry) TagRemove(ctx context.Context, name, tag string) error { return nil }

func (m *memoryRegistry) Deprecate(ctx context.Context, name, version, message string) error {
	return nil
}

// TestFirstPublicationRecordsPlatformsBeforeWrapper is row R13: the durable
// record proves every one of the four platform packages was published before
// the wrapper. This is the live owner of the platform-first, wrapper-last
// ordering — the release workflow's step-name check retires with the raw
// publish steps.
func TestFirstPublicationRecordsPlatformsBeforeWrapper(t *testing.T) {
	const version = "9.9.9"
	root := approvedReleaseRoot(t, version)

	record, err := RunFirstPublication(context.Background(), root, version, "public", newMemoryRegistry())
	if err != nil {
		t.Fatalf("first publication: %v", err)
	}

	kinds := map[string]string{}
	for _, provenance := range record.Provenance {
		kinds[provenance.Package] = provenance.Kind
	}
	var publishedKinds []string
	for _, transition := range record.Transitions {
		if transition.Action != "publish" {
			continue
		}
		if transition.Result != "success" {
			t.Fatalf("publish transition for %s is %q, want success", transition.Package, transition.Result)
		}
		publishedKinds = append(publishedKinds, kinds[transition.Package])
	}

	if len(publishedKinds) != 5 {
		t.Fatalf("recorded %d publish transitions, want four platform packages and one wrapper: %v", len(publishedKinds), publishedKinds)
	}
	for i, kind := range publishedKinds[:4] {
		if kind != "platform" {
			t.Fatalf("publish transition %d is a %q package; every platform package must be published before the wrapper: %v", i, kind, publishedKinds)
		}
	}
	if last := publishedKinds[4]; last != "wrapper" {
		t.Fatalf("the last publish transition is a %q package, want the wrapper: %v", last, publishedKinds)
	}
}
