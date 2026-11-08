# Name the conflicting pair when a record form gets two exclusive flags

Blocked by: none
Writes: internal/reviewrecord/recordcmd/command.go, internal/reviewrecord/recordcmd/refusal.go, internal/reviewrecord/recordcmd/conflict_test.go (new), internal/reviewrecord/recordcmd/review_test.go, internal/reviewrecord/recordcmd/verification_test.go, internal/reviewrecord/recordcmd/verification_import_test.go, CHANGELOG.md
Covers: RP29, RP30, RP31, RP32, RP54, RP55, RP57

## What to build

Review chunk: RP-C3.

When `together` (`internal/reviewrecord/recordcmd/command.go:172-181`) fails, find the first pair of present flags, in the order of flag declaration, that no admitted flag set holds together. Derive the pairs from the form layout that `admits` reads, and list no pair by hand. When such a pair exists, print `error: conflicting flags — <first> and <second> exclude each other`, then the form's usage line, at exit 2. When no pair exists, keep the bare usage line. A `valid` failure also keeps the bare usage line.

Put the conflict helper in `refusal.go`, and keep the branch in `command.go` to a few lines, because `command.go` has 26 lines of headroom and tickets 02 and 04 also add to it. Use `toon.Errorf` and no route literal, so the RR49 route-bypass scan stays green.

Move four existing expectations from the bare usage line to the conflict line: `review_test.go:158` (`--chunk` and `--fold`), `verification_test.go:367` (`--chunk` and `--final`; move this case out of its bare-usage loop into `conflict_test.go`, because `verification_test.go` has 10 lines of headroom), `verification_import_test.go:128-131` (`--probe-outcome` and `--probe-result`), and `verification_import_test.go:272` (`--excerpt` and `--mutant`). Add one `CHANGELOG.md` entry under `## [Unreleased]` for RP-C3.

## Acceptance

- [ ] `--excerpt` with `--probe-result` refuses at exit 2 and names both flags (RP29).
- [ ] Verification `--chunk` with `--final` refuses at exit 2 and names both flags (RP30).
- [ ] Review `--chunk` with `--fold` refuses at exit 2 and names both flags (RP31).
- [ ] An assignment call with `--predecessor` and no `--trigger` prints the bare usage line with no conflict text (RP32).
- [ ] `--probe-outcome` with `--probe-result` names both flags (RP54).
- [ ] `--excerpt` with `--mutant` names both flags (RP55).
- [ ] A call with `--chunk`, `--final`, `--excerpt`, and `--probe-result` names `--chunk` and `--final` (RP57).

## Verification

Run `bench test --package ./internal/reviewrecord/...` whole, and the RR49 test through `bench test --package ./internal/conformance --run '^TestNoWriteVerbComposesARouteOutsideTheRegistry$'`. Observe each planned row red before its production edit, and record the route of each red. Self-probe: print the bare usage line for every `together` failure, and report that RP29 to RP31 go red.
