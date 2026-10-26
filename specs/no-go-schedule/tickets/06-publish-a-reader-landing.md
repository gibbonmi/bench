# Publish a reader landing and keep the green marker

Blocked by: 05-plan-the-reader-set-at-a-landing.md
Writes: internal/gate/authorization/authorization.go, internal/gate/authorization/reader_kind_test.go (new), internal/landing/landing.go, internal/landing/landing_refusal.go (new), internal/landing/reader_publish_test.go (new), internal/worktree/land.go, internal/worktree/land_marker.go (new), internal/worktree/land_reader_test.go (new), internal/worktree/land_resume.go, internal/worktree/landingpolicy/landingpolicy.go, internal/worktree/landingpolicy/reader_marker_test.go (new), internal/gate/reader_plan.go, internal/gate/reader_plan_test.go, internal/worktree/parallel_census_test.go
Covers: NG55, NG56, NG57, NG58, NG59, NG60

## What to build

`authorization` gains the kind `reader green`. `AuthorizeWithWriters` reads `gate.Result.Schedule` from ticket 04. It maps a green reader set to `reader green`, with no evidence token, and a red reader set to `candidate`.

`landReviewed` calls `gate.WithLanding` with its destination commit and the green marker, so a reviewed landing is the one caller that can run the reader set. `ReviewedRequest` carries the marker that `landingDestination` already reads. The reviewed landing publishes on `green` and on `reader green`. The worktree-commit landing keeps `green` and `lane pass`.

After a `reader green` publication, the landing skips the green-marker advance at `internal/worktree/land.go:220`. The resume marker policy gains one fact: the published tree holds a reader record. `resumeMarkerFacts` reads it through `gate.ReaderRecorded` from ticket 04. When that fact holds, `ResumeMarker` answers accept, so `land_resume.go:82` advances nothing.

`internal/worktree/land.go` holds 397 lines and `internal/landing/landing.go` holds 392. Move the marker step into `land_marker.go`. Move `refusalMessage` and `refusalExplanation` into `landing_refusal.go`, with no change of text, so ticket 07 edits that file. Keep both old files at 400 lines or fewer.

Measure the premise before the chunk checkpoint. In a disposable clone of the integration source, first land one change that runs the whole gate. Then land one `ROADMAP.md` change that runs the reader set. Run both from the session shape of a real landing. Record both times and the reader set's phase table in the review pickup. If the reader set is not faster than the whole gate, stop and return the build to the spec.

## Acceptance

- [ ] A green reader set attributes the kind `reader green` with no evidence token.
- [ ] A red reader set attributes the kind `candidate`.
- [ ] The reviewed landing publishes on `reader green`, and the worktree-commit landing refuses it.
- [ ] After a `reader green` publication, `refs/bench/green/main` names the prior marker.
- [ ] `ResumeMarker` answers accept for a published tree that holds a reader record.
- [ ] The review pickup records one measured whole landing and one measured reader landing from a disposable clone.
