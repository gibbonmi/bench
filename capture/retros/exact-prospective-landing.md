## Outcome

Run `869a2d33cab96f962882762a9ea56fd21c952976248e64fa917f6da6c48dbca3`
promoted exact candidate `616d6d5efdac96d3752266d3e1c9acb47793b0e4` to `main` as
`1254ce6`. The implementation added exact prospective landing, exact-tree gate
authorization and reuse, post-publication reconciliation from the authorized
snapshot, clean spec transition, adapter adoption, and public runtime proofs.
All 25 retained assignments released. Fresh Standards, Spec, and Coverage axes
reported zero findings. The spec was then retired in `126c597`.

The four explicitly parked runtime proof gaps were detached-HEAD CAS loss,
prepublication SIGINT, control-byte path input, and dangling-symlink behavior.
They remained noncritical and outside this landing. Cross-harness review was
omitted by user direction.

## Gate-stage timings

The authoritative promotion gate completed green in approximately four minutes
of observed wall time. Its terminal command did not emit per-stage elapsed
summaries, so no stage duration is reconstructed here. The retained lifecycle
evidence binds the green verdict to candidate `616d6d5e` and review receipt
`edea5c7de4d166ef509df37231a5ed6146706af946fb8c18c477ce06acb1d280`.

The separate required retirement commit also completed green. Timings it emitted
included runtime contracts at 125.405s, conformance-suite at 138.412s,
`internal/gate` tests at 219.120s, `internal/specbuild` tests at 143.872s, and
`internal/worktree` tests at 84.786s. Build, gofmt, vet, test, race,
conformance-suite, conformance, and contract all reported green; 11 capability
skips were classified as four capability and seven environment skips.

## Ticket-versus-spec-slice and delegate performance

The original landing and adapter split respected the two production ownership
fences, but lifecycle and review repairs expanded the run to 25 assignments.
Small exact-fence charges performed best: the higher-tier landing delegate closed
the tracked/spec authorization race and later the late-untracked-directory race
with deterministic mutation evidence; the nested-CWD delegate returned a
single-file real-binary proof and a 63.692s full runtime-package green run.

Broad or documentation-shaped charges were less reliable. Fable and Opus both
stalled on the clean-release write charge, so the coordinator completed that
two-file fence. The first Standards repair delegate produced no diff in its
bounded window. Fable with explicit YOLO permissions then produced the complete
nine-file correction but did not return a final response; the coordinator
validated and lifecycle-integrated the filesystem result. The temporary
normalization ticket itself was circular evidence and had to be deleted through
a self-deleting cleanup assignment before final review.

## Coordinator catches

The coordinator caught the first authorized-snapshot draft still re-reading and
re-transforming the live spec after authorization, and required immutable spec
bytes plus captured Git and filesystem modes. Full package verification then
surfaced the staged-deletion/rename restore fallback. Final Spec review found the
remaining late untracked descendant under a named directory; final Coverage
found the valid nested-CWD real-binary gap.

The lifecycle also exposed prerequisites outside the initial fence: the FT78
symlink fixture defect, clean refresh/checkpoint/integration/release closure, and
the 64-entry operation journal ceiling. The journal was raised to 128 with live
65th-entry and full-bound mutation proofs. Ticket assignment caught an unanchored
`Contracts:` line before delegation. Recomposition later refused a same-file
Standards patch conflict; restoring the patchable base and using a self-deleting
cleanup ticket preserved every checkpoint without discarding work.

## Agent-experience improvements

### Bench CLI

Long `spec build promote` runs should emit phase heartbeats and retained stage
timings; the successful authoritative gate was silent until terminal status.
Recomposition conflicts should identify the checkpoint owner and offer a bounded
already-covered or compatible-same-path reconciliation route. The operation
journal should expose remaining capacity before the hard bound is reached.
Headless delegate wrappers should distinguish “edits produced, response channel
stalled” from “no work produced” so a usable fence is not mistaken for a failed
charge.

### Skills

`craft-tickets` should call out that review-only normalization records must be
folded into substantive tickets or deleted before exact review; a fresh semantic
axis is not an independent red owner. The repair-ticket checklist should require
the four discovery lines, concrete path-or-row integration mappings, and a real
public mutation owner before the ticket is committed. `craft-review` should add
post-authorization directory descendants to its concurrency edge inventory.

### Process

Run deterministic tracked-file, transitioned-spec, and new-untracked-descendant
mutations together when first reviewing an immutable authorization snapshot.
Exercise valid nested-CWD arguments as well as usage failures in the first public
adapter slice. Before a long repair run approaches its journal bound, expose and
resolve capacity as a prerequisite rather than discovering it at checkpoint.
Keep the final three-axis review pinned to one exact candidate; it found two
critical gaps that focused package greens alone could not see.
