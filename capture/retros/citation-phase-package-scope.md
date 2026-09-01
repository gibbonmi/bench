# Retro: citation-phase-package-scope (FT281)

## Outcome

FT281 landed at `75781592`. The citation validator now binds a cited file to
one Go test phase's package scope, not just its build tags. It keeps one
execution entry per phase. It uses a real `go list` child for package
expansion. It applies a phase's own env for GOOS, GOARCH, and CGO_ENABLED.

The spec closed at 13 acceptance rows (PS1–PS13). PS13 was a
reviewer-approved addition made during repair. One ticket built the whole
spec.

## Gate-stage timings

Final landing gate, tree `d359e0f`:

| Phase | Verdict | Elapsed |
|---|---|---|
| gofmt | green | 121ms |
| vet | green | 1.1s |
| test | green | 57.0s |
| race | green | 2.7s |
| system | green | 22.1s |
| shellcheck | green | 519ms |

## Ticket-versus-spec-slice and delegate performance

One ticket carried the whole spec (12, later 13, acceptance rows) as a
single charge, not sliced across delegates. The first write-delegate pass
reached roughly 9 of 12 rows with matching tests. It needed a second charge
to close PS9, PS10, and PS12. That charge also added a `CHANGELOG.md` entry
and fixed a regression in an existing test file the delegate had not been
told to expect. A reviewer-approved third charge added a `go list` timeout.
The delegate's own testing had surfaced a hang bug the original charge never
anticipated.

A whole-spec charge finds real bugs beyond its checklist. It still needs at
least one follow-up charge to close every row — plan for two, not one.

## Coordinator catches

The coordinator independently re-ran every delegate's test claim, uncached
and with `-count=1`. It also re-ran `gofmt`, `go vet`, and
`bench coverage --check`, rather than trusting a done-claim. This caught
nothing the delegates had gotten wrong on tests. It did catch the
coordinator's own mistake: landing a fence amendment with `--spec` flipped
`Status` to `implemented` before the ticket's code had landed. The next
`bench preflight build` refusal caught it, fixed before it compounded. Two
review rounds found 13 findings total, most cited by more than one axis
independently. They caught three defects no test-rerunning would have:

- a `bounds` policy redeclare invisible to the fast lane
- a relative-path bug invisible to the documented invocation's normal callers
- a flag-parsing gap invisible without adversarial testing against
  `go help testflag`

## Repair attribution

| Ticket | Repair rounds | Cause per round |
|---|---|---|
| bind-citations-to-phase-package-scope | 2 | delegate-error, delegate-error |

## Agent-experience improvements

### Bench CLI

- Census: assignment `4ea0611e91352257a0cee32e3af43ab4` logged 256 raw
  calls. Its per-verb breakdown was lost before this session read it.
  Captured as `bench learning "ft281-citation-scope census: 256 raw calls,
  breakdown lost"`.
  Feeds: none
- Add a `--wait` option to `bench worktree land` that polls the current gate
  owner's PID instead of refusing immediately. Two concurrent Bench sessions
  on one machine serialize on the gate lock, and the caller has no built-in
  way to wait.
  Feeds: new
- Print a landing's census per-verb breakdown in the `landed{...}` record
  itself, not only the count. A retro should not need a separate pre-land
  read of `~/.bench/census/<pool>/<assignment>` before the record is deleted
  at release.
  Feeds: new
- `bench worktree land --spec <slug>` flips `Status` to `implemented` on any
  landing under that spec's fence. This includes a stand-alone
  fence-amendment landing that touches only `spec.md`. Refuse `--spec` on a
  landing whose diff does not also touch a ticket's `Writes:` paths.
  Alternatively, require a distinct flag for the status-flip landing.
  Feeds: new

### Skills

- State a proportionality rule in `bench-review-implementation`. It should
  say how many repair-scoped rounds re-invoke the full three-agent axis
  fan-out, versus a coordinator's own direct verification. The skill leaves
  this to judgment today.
  Feeds: new

### Process

- A ticket needing an ownership-fence amendment mid-build costs a full
  separate worktree-create, commit, and land cycle each time. This ticket
  needed two. Allow a fence amendment to fold into the same landing as the
  ticket it unblocks when discovered mid-repair.
  Feeds: new
