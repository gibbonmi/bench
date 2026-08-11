# Stream retained-run enumeration into a bounded projection

Blocked by: none
Ownership fence: `internal/specbuild/state.go`, `internal/specbuild/runs_test.go`
Integration surfaces: retained-state directory streaming and bounded projection→`internal/specbuild/state.go`; external directory-reader test double and creation-order fixtures→`internal/specbuild/runs_test.go`; existing family-home rendering→existing `cmd/bench/specbuild.go` exercised unchanged by BE2 through BE5
Contracts: entry batches cross the external directory reader→the bounded accumulator in `internal/specbuild/state.go`, asserted by BE1 against the real reader contract and its test double; bounded healthy rows, diagnostic representatives, and overflow cross `internal/specbuild/state.go`→the existing family-home renderer, asserted by BE2 through BE5
Closure: BE1/bounded-batch-request, BE1/no-whole-directory-read, BE2/creation-order-invariant, BE2/legacy-byte-projection, BE3/healthy-selection, BE4/foreign-representative, BE4/lock-representative, BE4/symlink-representative, BE4/nonregular-representative, BE4/malformed-representative, BE4/unreadable-representative, BE4/nondigest-representative, BE5/honest-at-least-marker

## What to build

Close the buildable P3/C3 defect without contradicting SB4. `Service.Runs`
streams the retained-state directory in bounded batches instead of materializing
the whole directory, classifies the full stream so hostile entries cannot hide
healthy runs, and retains only the bounded state needed to reproduce the current
deterministic projection. Identical entry sets created in different orders emit
identical runs, diagnostic representatives, and overflow marker bytes. Runtime
remains linear in the entries classified; the bounded contract is memory and
output size, not a promise to ignore entries the spec requires it to classify.

## Acceptance

- [ ] [BE1] (covers SB4) (P3, C3) retained-state enumeration requests bounded batches from the directory reader and never requests the whole directory in one call.
- [ ] [BE2] (covers SB4) (P3, C3) identical retained-state entry sets created in different orders produce the same bytes as the current sorted projection.
- [ ] [BE3] (covers SB4) (P3, C3) bounded accumulation retains the lexically selected healthy runs even when hostile entries precede or follow them.
- [ ] [BE4] (covers SB4) (P3, C3) bounded accumulation retains the lexical representative of each existing diagnostic class: foreign, lock, symlink, non-regular, malformed, unreadable, and non-digest name.
- [ ] [BE5] (covers SB4) (P3, C3) any stream larger than the named entry cap emits exactly the existing `at_least_<cap+1>` entry-cap diagnostic without increasing output size.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| BE1/bounded-batch-request | request all remaining entries from the directory reader | the injected directory-reader contract test | drive `Service.Runs` through a reader that records every requested batch size and require each request to be positive and bounded |
| BE1/no-whole-directory-read | restore `os.ReadDir` materialization of the complete directory | the injected directory-reader contract test | drive more than one batch and require the production scan to continue through the bounded reader seam rather than bypass it |
| BE2/creation-order-invariant | retain filesystem iteration order in the bounded projection | the paired creation-order fixture | create the same mixed entry set in forward and reverse order, invoke fresh services, and require byte-identical rendered projections |
| BE2/legacy-byte-projection | change which bounded rows or diagnostics the current sorted projection selects | the legacy-equivalence fixture | feed a mixed over-cap entry set to the current reference projection and the streaming projection and require exact equality |
| BE3/healthy-selection | let hostile entries consume every retained-run slot | the hostile-before-and-after healthy fixture | place a healthy run beyond hostile entries in each creation order and require the healthy row to survive |
| BE4/foreign-representative | drop or replace the lexical foreign-entry representative | the mixed diagnostic fixture | stream multiple foreign entries and require the current lexical representative |
| BE4/lock-representative | drop or replace the lexical lock-entry representative | the mixed diagnostic fixture | stream multiple regular lock companions and require the current lexical representative |
| BE4/symlink-representative | drop or replace the lexical symlink representative | the mixed diagnostic fixture | stream multiple symlink entries and require the current lexical representative |
| BE4/nonregular-representative | drop or replace the lexical non-regular representative | the mixed diagnostic fixture | stream multiple supported non-regular entries and require the current lexical representative |
| BE4/malformed-representative | drop or replace the lexical malformed-record representative | the mixed diagnostic fixture | stream multiple malformed records and require the current lexical representative |
| BE4/unreadable-representative | drop or replace the lexical unreadable-record representative | the mixed diagnostic fixture | stream multiple deterministically unreadable records and require the current lexical representative |
| BE4/nondigest-representative | drop or replace the lexical non-digest-name representative | the mixed diagnostic fixture | stream multiple valid records under noncanonical names and require the current lexical representative |
| BE5/honest-at-least-marker | omit the overflow marker or derive it from retained rows instead of entries seen | the streaming overflow fixture | stream cap-plus-one and many-more subjects, require the same single `at_least_<cap+1>` marker, and require bounded output in both cases |
