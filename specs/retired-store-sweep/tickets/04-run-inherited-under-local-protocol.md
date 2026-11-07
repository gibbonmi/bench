# Run both runners through the one inherited registration owner

Blocked by: none
Writes: internal/subprocess/inherited.go, internal/subprocess/startgate_test.go, internal/bounds/bounds.go, internal/runbinary/runbinary.go, internal/runbinary/lifetime_resources_test.go, internal/gate/runner.go, internal/gate/background_launch.go, internal/preprelease/preprelease.go, internal/releasepreflight/command.go, internal/testreport/command.go, internal/conformance/process_lifetime_test.go, internal/conformance/process_lifetime_closure_test.go, internal/conformance/process_lifetime_copy_test.go (new)
Covers: RS16, RS17, RS18, RS19, RS20, RS21, RS22, RS23, RS24, RS29

## What to build

`RunInherited` takes a `Registration` in place of its variadic local descriptors.
The registration's resources join the inherited union after the inherited resources, as today. A nil protocol means the installed protocol.
The five production callers keep their own signatures. Each wraps its descriptors into the registration at its `RunInherited` call, and the detached launch passes an empty registration. Their test call sites do not change.
`internal/gate/runner.go` keeps a net growth of zero or less, and `internal/testreport/command.go` stays inside its 13 lines of headroom.

The bounds runner keeps its parent check, then launches through `RunInherited` and no longer imports the protection package.
It maps any `RunInherited` error to a start failure. Its cancel request answers an interruption when the parent context ended and a deadline otherwise, each with KILL.
The build-script runner keeps its toolchain check, its environment, and its process attributes, then forwards its `local` registration to `RunInherited`.
Its cancel request answers a deadline when the context expired and an interruption otherwise, with TERM. It keeps its own outcome mapping.

The closure check adds the hand-built registration rule to the file diagnostics in `process_lifetime_closure_test.go`.
The rule counts `protection.Inherited()`, `protection.InheritedFrom(os.Environ())`, and `protection.ParseVector(os.Getenv(protection.Env))`. It refuses a production function outside `internal/subprocess` that calls a counted reader and builds a `subprocess.Registration`, in its own body or in a same-package function that it calls.
The closure prefilter adds the reader tokens, so it parses each file that spells a counted reader.
The synthetic fixtures go in the new `process_lifetime_copy_test.go`: one per counted reader, one with the `Registration` in a same-package helper, and one uncounted reader that stays green.

## Acceptance

- [ ] With an empty inherited vector set through `t.Setenv`, `RunInherited` registers the launch on the protocol that its registration names, and the child environment names the local resource.
- [ ] `TestBoundsMalformedResources`, `TestBoundsCancel`, `TestBoundsFixedShutdown`, `TestBuilderCancel`, and `TestBuilderNormalDrain` stay green.
- [ ] A malformed inherited vector refuses the build script before it starts.
- [ ] The closure rule names each synthetic hand-built copy and stays green on the uncounted reader.
- [ ] The live `process-lifetime-closure` check stays green with both runner census entries.
- [ ] Neither runner reads the process vector or builds a `subprocess.Registration`.
