# Hold the whole gate after a red whole landing run and name the suspects

Blocked by: 06-publish-a-reader-landing.md
Writes: internal/gate/reader_hold.go (new), internal/gate/reader_hold_test.go (new), internal/gate/reader_plan.go (new), internal/gate/run_transaction.go, internal/gate/gate.go, internal/gate/authorization/authorization.go, internal/landing/landing_refusal.go (new), internal/landing/reader_refusal_test.go (new), internal/git/git.go
Covers: NG61, NG62, NG63, NG64, NG78

## What to build

A red whole run that carries a landing value writes the reader hold. The hold names the green marker commit of that run, under a file name that `internal/git/git.go` declares. A red reader set writes no hold.

The reader plan of ticket 05 gains case 9: a hold that names the current green marker sends the run to the whole gate. A hold that names another commit does not apply, so a whole green landing retires the hold by moving the marker, with no write.

The same red run lists its suspects: each commit that `git.FirstParentsAfter` names and whose tree `gate.ReaderRecorded` answers true for. `gate.Result` carries the list, and `authorization.Result` carries it to the landing. `refusalMessage`, which ticket 06 moved into `internal/landing/landing_refusal.go`, appends `; reader landings since the green marker: <commits>` when the list is not empty. With an empty list, the message keeps its current text, so the two exact-text tests at `internal/landing/landing_reviewed_test.go:399` and `:430` stay green.

`internal/gate/run_transaction.go` must stay at 400 lines or fewer. Put the hold logic in `reader_hold.go`.

## Acceptance

- [ ] A red whole run with a landing value writes a reader hold that names the run's green marker.
- [ ] A reader hold that names the current green marker sends a covered landing to the whole gate.
- [ ] A reader hold that names an older green marker leaves the reader set selectable.
- [ ] A red whole landing run after two reader-record commits lists both commits as suspects in its result.
- [ ] A refusal whose result carries two suspects names both in the refusal line, and an empty list keeps the current line.
