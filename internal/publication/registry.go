package publication

import "context"

// Registry is the npm-registry port the state machine drives. Two adapters
// implement it. FixtureRegistry works over HTTP against the hermetic
// offline-registry.mjs fixture, used by the gate. NPMCLIRegistry shells the
// real `npm` CLI, for the runbook path only, never exercised by the gate.
//
// Neither adapter, nor any caller, ever puts a credential into the durable
// record or the evidence trail. Auth material stays in the process
// environment/npm config and is referenced only by an opaque auth_mode
// label.
type Registry interface {
	// Publish uploads tarball for name@version under the given non-default
	// dist-tag in one call (direct publish — no staging). It returns the
	// registry's own integrity (SHA-512 SRI) for the now-live package.
	Publish(ctx context.Context, name, version, tag string, tarball []byte) (registryIntegrity string, err error)

	// StageSubmit uploads tarball to a staging area that is not yet live and
	// returns an opaque stage id. Nothing is publicly visible until Approve.
	StageSubmit(ctx context.Context, name, version string, tarball []byte) (stageID string, err error)

	// Approve promotes a previously staged submission to live.
	Approve(ctx context.Context, stageID string) error

	// Integrity performs a read-only query of the registry's live integrity
	// for name@version. live is false when the version is not (yet) published.
	Integrity(ctx context.Context, name, version string) (integrity string, live bool, err error)

	// TagAdd points tag at an already-live version.
	TagAdd(ctx context.Context, name, tag, version string) error

	// TagRemove removes tag from name. It is never used to unpublish a
	// version — there is no unpublish operation in this port.
	TagRemove(ctx context.Context, name, tag string) error

	// Deprecate attaches a deprecation message to a live version.
	Deprecate(ctx context.Context, name, version, message string) error
}
