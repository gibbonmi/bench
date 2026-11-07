# Run both runners through the one inherited registration owner

Blocked by: none
Writes: internal/subprocess/inherited.go, internal/subprocess/startgate_test.go, internal/bounds/bounds.go, internal/runbinary/runbinary.go, internal/runbinary/lifetime_resources_test.go, internal/gate/runner.go, internal/gate/background_launch.go, internal/preprelease/preprelease.go, internal/releasepreflight/command.go, internal/testreport/command.go
Covers: RS16, RS17, RS19, RS20, RS21, RS22, RS23, RS24

## What to build

`RunInherited` takes a `Registration` in place of its variadic local descriptors.
The registration's resources join the inherited union after the inherited resources, as today. A nil protocol means the installed protocol.
The five production callers keep their own signatures. Each wraps its descriptors into the registration at its `RunInherited` call, and the detached launch passes an empty registration. Their test call sites do not change.
`internal/gate/runner.go` keeps a net growth of zero or less, and `internal/testreport/command.go` stays inside its 13 lines of headroom.

The bounds runner keeps its parent check, then launches through `RunInherited` and no longer imports the protection package.
It maps any `RunInherited` error to a start failure. Its cancel request answers an interruption when the parent context ended and a deadline otherwise, each with KILL.
The build-script runner keeps its toolchain check, its environment, and its process attributes, then forwards its `local` registration to `RunInherited`.
Its cancel request answers a deadline when the context expired and an interruption otherwise, with TERM. It keeps its own outcome mapping.

Ticket 05 adds the closure rule that refuses a surviving hand-built copy.

## Acceptance

- [ ] With an empty inherited vector set through `t.Setenv`, `RunInherited` registers the launch on the protocol that its registration names, and the child environment names the local resource.
- [ ] `TestBoundsMalformedResources`, `TestBoundsCancel`, `TestBoundsFixedShutdown`, `TestBuilderCancel`, and `TestBuilderNormalDrain` stay green.
- [ ] A malformed inherited vector refuses the build script before it starts.
- [ ] The live `process-lifetime-closure` check stays green with both runner census entries.
- [ ] Neither runner reads the process vector or builds a `subprocess.Registration`.
- [ ] A launch through `RunInherited` with a nil protocol registers on the installed protocol, and `TestProcessLifetimeCacheSecondProcessClean` stays green.
