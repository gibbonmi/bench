# Move the working-directory tests

Blocked by: 06-carry-the-call-directory.md
Writes: internal/worktree/identifier_operand_test.go, internal/worktree/resume_reconcile_test.go, internal/worktree/request_token_test.go, internal/worktree/identity_component_test.go, internal/worktree/land_resume_test.go, internal/worktree/land_reauthorize_operand_test.go, internal/worktree/path_identifier_test.go, internal/worktree/lifecycle_policy_test.go, internal/worktree/serial_ceiling_test.go, internal/worktree/parallel_census_test.go
Covers: LGP20, LGP21, LGP22, LGP23, LGP34

## What to build

The worktree tests that changed the process directory pass the call directory instead, and they call `t.Parallel`.

Sixteen movable `chdir` tests sit in eight files:

- `identifier_operand_test.go`: 5 tests (`:18`, `:54`, `:73`, `:181`, `:223`);
- `resume_reconcile_test.go`: 4 tests (`:243`, `:262`, `:290`, `:336`);
- `request_token_test.go`: 2 tests (`:18`, `:35`);
- one test each in `identity_component_test.go:358`, `land_resume_test.go:213`, `land_reauthorize_operand_test.go:15`, `path_identifier_test.go:70`, and `lifecycle_policy_test.go:170`.

Replace each `chdir` with the verb call's directory field from ticket 06. A test whose verb has no joins form cannot take the field. Remove its `chdir` only if the verb reads no process directory, and keep it serial otherwise.

Record each test that moves and each test that stays, with its reason.

`TestLandingHandsBackAPathWithNoCanonicalForm` (`land_resume_refusal_test.go:343`) stays serial, because the real `filepath.Abs` failure is its subject. Put the LGP23 test in `serial_ceiling_test.go`, because `parallel_census_test.go` has 1096 lines and no grant.

Set both pins by the recorded arithmetic: the ceiling falls by the moved tests, and the count rises by the new tests. Update the four-class comment at `parallel_census_test.go:1033-1040` with no net line gain.

Record `bench test --full --package ./internal/worktree` before and after the change, for LGP34.

## Acceptance

- [ ] `TestParallelCensusOnTheLiveTree` passes with every moved test parallel (LGP20).
- [ ] The serial set equals the recorded ceiling, and the test count equals the recorded sum (LGP21, LGP22).
- [ ] The removed-directory test stays serial and keeps its real `filepath.Abs` failure (LGP23).
- [ ] The probe "leave the ceiling at its ticket-06 value" bites.
- [ ] The verification record holds the package time before and after the change (LGP34).
