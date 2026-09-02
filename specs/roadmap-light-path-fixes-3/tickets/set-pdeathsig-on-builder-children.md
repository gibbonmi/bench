# Set Pdeathsig on builder children

Blocked by: derive-the-canonical-path-in-one-leaf-package.md
Writes: internal/runbinary/runbinary.go, internal/runbinary/sysprocattr_linux.go (new), internal/runbinary/sysprocattr_darwin.go (new), internal/runbinary/sysprocattr_other.go (new), internal/runbinary/runbinary_test.go, internal/runbinary/pdeathsig_linux_test.go (new)
Covers: LQ24, LQ25, LQ26

## What to build

Verify the premise first: `canonicalBuild` in internal/runbinary/runbinary.go
sets `Setpgid` and nothing else on the builder child. Then move the attribute
literal into a per-platform function with three build-tagged files, mirroring
`internal/releaseevidence/evidence_exchange_{linux,darwin,other}.go`. The Linux
leg sets `Setpgid` and `Pdeathsig: syscall.SIGKILL`; the other legs set
`Setpgid` alone. Extract the inline parking-builder script that
`TestCanonicalBuilderDrainsDescendantsBeforeReturningSelection` in
internal/runbinary/runbinary_test.go carries into a shared in-package helper,
and add a Linux-tagged test on it. Run the owner as its own process, SIGKILL it, and poll the builder child until `ESRCH` within
`bounds.TestDeadline`. The testreport fixture is unexported and outside the
fence, so do not reach for it.

## Acceptance

- [ ] The Linux test observes the builder child dead after SIGKILL to its owner.
- [ ] `TestCanonicalBuilderDrainsDescendantsBeforeReturningSelection` still passes.
- [ ] The cross-compile conformance check passes for darwin.
- [ ] Self-probe: drop `Pdeathsig` from the Linux leg, and report the new test red.
