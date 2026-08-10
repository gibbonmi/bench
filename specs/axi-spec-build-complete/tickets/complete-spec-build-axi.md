# Complete spec build AXI

Blocked by: none
Ownership fence: `internal/axi/`, `internal/specbuild/`, `internal/spec/build.go`, `cmd/bench/specbuild.go`, `cmd/bench/specbuild_test.go`, `cmd/bench/testdata/`, `internal/conformance/`, `projects/benchkit.md`
Integration surfaces: action owner and owner tests→`internal/axi/`; shell serialization→existing `internal/sanitize/sanitize.go` exercised by AX1; TOON encoding→existing `internal/toon/toon.go` exercised by AX1 and AX2; operation catalog and root grammar→`internal/spec/build.go`; lifecycle facts, refusals, retained-run enumeration, matrix ledger, and fixtures→`internal/specbuild/`; production route, response renderers, CLI tests, and paired fixtures→`cmd/bench/specbuild.go`, `cmd/bench/specbuild_test.go`, and `cmd/bench/testdata/`; unchanged command registration→existing `cmd/bench/main.go` exercised by AX2 and AX4; conformance registry, executable binding, and AXI check→`internal/conformance/`; advertised check input→`projects/benchkit.md`; unchanged shell operation advertisement→existing `bin/bench.sh` exercised by AX7
Contracts: none crosses
Closure: AX1/shell-kind, AX1/phase-kind, AX1/literal-token, AX1/fixed-value, AX1/open-placeholder, AX1/prose-refusal, AX1/quoted-hostile-value, AX1/newline-refusal, AX1/control-refusal, AX1/empty-help, AX2/operation-order, AX2/operation-cell, AX2/state-cell, AX2/known-slug, AX2/known-assignment, AX2/known-ticket, AX2/known-receipt, AX2/known-fingerprint, AX2/open-future-input, AX2/next-first-action, AX2/next-empty, AX3/refusal-class, AX3/exact-remedy, AX3/orchestration-inspection, AX4/retained-row, AX4/definitive-empty, AX4/record-suffix, AX4/lock-skip, AX4/foreign-diagnostic, AX4/nonregular-diagnostic, AX4/symlink-diagnostic, AX4/malformed-diagnostic, AX4/unreadable-diagnostic, AX4/healthy-row-survives, AX4/entry-cap, AX4/at-least-marker, AX4/repeatable-bytes, AX4/help-word, AX4/help-long, AX4/help-short, AX4/root-home, AX5/abandon-apply, AX5/reclaim-apply, AX5/fixed-slug, AX5/fixed-fingerprint, AX5/abandon-preimage, AX5/reclaim-preimage, AX5/stale-abandon-replan, AX5/stale-reclaim-replan, AX5/spent-reclaim-replan, AX6/operation-totality, AX6/state-totality, AX6/root-help-totality, AX6/plan-apply-totality, AX6/old-bytes, AX6/new-bytes, AX6/help-delta, AX6/next-delta, AX6/home-delta, AX6/unnamed-byte-refusal, AX7/operation-axis, AX7/state-axis, AX7/unclassified-cell, AX7/structured-stdout, AX7/success-exit, AX7/refusal-exit, AX7/usage-exit, AX7/operation-help, AX7/root-help, AX7/help-envelope, AX7/registry-binding, AX7/profile-binding, AX8/useful-action-removal, AX8/known-value-placeholder, AX8/unknown-value-guess, AX8/dropped-fixed-flag, AX8/stale-fingerprint, AX8/prose-advertisement, AX8/unquoted-hostile-value

## What to build

Land the spec's one atomic migration as one independently-green vertical ticket. Add the typed `internal/axi` action owner; carry typed lifecycle outcomes, refusal preconditions, and authority facts through the real service; replace prose `next` values with the first derived action or empty; append the help envelope to all nine operations; add the bounded retained-run family home and shared root help spellings; pin every intentional output delta in one paired fixture ledger; and register the registry-derived AXI conformance check. Keep all existing primary bytes, streams, exits, malformed-argv behavior, and gate payloads except the spec's three named delta classes.

The single ticket is required because thinner cuts strand the named atomic compatibility oracle: an action-only prefactor leaves an unconsumed exported carrier, operation batches leave the constructor-closed matrix and paired-fixture totality red, the family root changes the parser-visible bytes before the fixture ledger authorizes the delta, and conformance cannot land before every operation/state cell is classified. The ticket owns every changed producer and consumer, so no value crosses its ownership fence.

## Acceptance

- [ ] [AX1] (covers SB1) the typed action owner distinguishes shell and harness-phase templates, literal tokens, fixed values, and open placeholders; refuses prose and non-single-line values; quotes fixed hostile argv; and renders useful or honest zero-row help.
- [ ] [AX2] (covers SB2) every cell of the accessor-derived nine-operation and typed-state matrix appends help from lifecycle facts, carries every known identifier and only open future input, and derives `next` from the first action or empty.
- [ ] [AX3] (covers SB3) every typed refusal constructor retains its precondition class through rendering and advertises the exact satisfying command, while orchestration-only states advertise full status rather than prose.
- [ ] [AX4] (covers SB4) the no-operation family home and all shared help spellings return exit 0; the bounded durable enumerator renders every healthy run, a definitive empty, honest overflow, and named hostile-entry diagnostics without hiding healthy rows.
- [ ] [AX5] (covers SB5) abandon and reclaim plans emit exactly one apply action with fixed slug and current fingerprint, and stale or spent authorization emits only the exact re-plan action.
- [ ] [AX6] (covers SB6) one registry-derived ledger pairs independently authored old/new bytes for every operation/state/root/help/authority case, names only the approved delta classes, and rejects every unnamed byte change.
- [ ] [AX7] (covers SB7) one registered conformance check derives the production operation/state axes and enforces structured stdout, exit 0/1/2, shared help spellings, help envelopes, total classification, executable binding, and profile agreement.
- [ ] [AX8] (covers SB8) each of the seven action mutation classes independently turns its focused owner test red.

## Handoff ledger

| spec account | owner or disposition |
|---|---|
| stories 1-4 | AX1/AX2/AX6/AX7/AX8; AX3/AX7; AX4/AX6/AX7; AX5/AX6 |
| SB1-SB8 | AX1, AX2, AX3, AX4, AX5, AX6, AX7, AX8 respectively |
| action, render, refusal, home, plan/apply, paired-fixture, and conformance seams | exact paths in `Ownership fence:` and `Integration surfaces:`; all changed producers and consumers share this ticket |
| error and malformed input | AX3 typed refusals; AX7 preserves usage/2 and structured stdout |
| empty and absent input | AX4 definitive empty home and AX2 empty lifecycle projections |
| boundary values | AX2 one/many states; AX4 cap and honest overflow |
| interrupted or partial state | AX5 spent reclaim receipt and fresh re-plan action |
| rerun idempotency | AX4 repeated fresh-service home bytes; AX2 no-op exit 0 |
| process-boundary lifecycle | AX2/AX4 fixtures reload durable state in fresh services |
| hostile environment | AX1 argv quoting/refusal and AX4 hostile state entries |
| command self-observation | AX2 help reads typed facts without mutating lifecycle state |
| special files and dangling symlinks | AX4 stats before reading and emits bounded diagnostics |
| Won't handle | ambient `bench status` projection and generated skill remain out of scope |
| approved fences | every path from spec lines 104-106 is owned above; gate/lifecycle authority and all other families stay unchanged inputs |

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AX1/shell-kind | construct a shell command as a harness phase | the pure action-kind test | construct both kinds and require shell identity for the shell command |
| AX1/phase-kind | construct a harness phase as a shell command | the pure action-kind test | construct both kinds and require canonical `/bench-*` phase identity |
| AX1/literal-token | swap two literal argv tokens | the pure argv-order test | construct a multi-token command and require exact order |
| AX1/fixed-value | replace a fixed slug with `<slug>` | the fixed-value test | render the known slug and require its literal value |
| AX1/open-placeholder | quote or guess an open placeholder | the open-placeholder test | render unknown future input and require bare `<name>` |
| AX1/prose-refusal | admit `release assignment 123` as executable | the prose-refusal test | construct the label and require typed refusal |
| AX1/quoted-hostile-value | bypass `sanitize.ShellQuote` for `build demo` | the hostile-value test | render the fixed value and require copy-paste-safe argv |
| AX1/newline-refusal | allow a newline-bearing fixed value | the newline-refusal test | construct `build\ndemo` and require refusal before render |
| AX1/control-refusal | allow a fixed ESC byte | the control-refusal test | construct the control-bearing value and require refusal |
| AX1/empty-help | omit `help[0]{command}:` | the empty-help test | render no actions and require exact terminal bytes |
| AX2/operation-order | remove or reorder one accessor member | the operation-accessor test outside `internal/spec` | call the exported accessor through the approved CLI test surface and require all nine in production order without mutability leakage |
| AX2/operation-cell | drop one operation from the matrix | the real-service matrix test | enumerate the accessor and require a classification set per operation |
| AX2/state-cell | add a typed state without a disposition | the real-service matrix test | enumerate typed classes and require fixture or exact not-applicable per cell |
| AX2/known-slug | replace a known slug with `<slug>` | the carried-facts matrix test | drive the state and require its concrete slug |
| AX2/known-assignment | drop the assignment ID | the carried-facts matrix test | drive assigned state and require its concrete assignment |
| AX2/known-ticket | drop the ticket basename | the carried-facts matrix test | drive assignment readiness and require the concrete ticket |
| AX2/known-receipt | replace a known receipt with `<receipt>` | the carried-facts matrix test | drive receipt-bearing state and require its concrete path |
| AX2/known-fingerprint | replace an available fingerprint with a placeholder | the carried-facts matrix test | drive the state and require the exact fixed fingerprint |
| AX2/open-future-input | guess a request or evidence value not yet created | the open-input matrix test | drive pre-creation state and require the open placeholder |
| AX2/next-first-action | derive `next` independently from help | the paired renderer test | render a state and require byte identity with the first help command |
| AX2/next-empty | leave orchestration prose in an empty-action `next` | the terminal renderer test | require empty `next` and honest zero-row help |
| AX3/refusal-class | collapse two refusal constructors into one text class | the refusal-closure test | drive both through the real service and require distinct typed identities |
| AX3/exact-remedy | advertise retry/assign for stale exact-green evidence | the class-exact remedy test | drive the class and require `bench gate --fresh` |
| AX3/orchestration-inspection | advertise `release assignment <id>` as executable | the orchestration-state test | drive release/delegate pending and require exact full-status action |
| AX4/retained-row | omit one valid record | the retained-run fixture | write two valid records and require both rows |
| AX4/definitive-empty | return silence or usage/2 for no state directory | the empty-home fixture | invoke the family root and require zero rows, help, and exit 0 |
| AX4/record-suffix | accept a foreign suffix as a record | the hostile-entry fixture | add a foreign file and require a diagnostic, not a run |
| AX4/lock-skip | render the owner's lock as a run or foreign error | the lock fixture | add a valid record and lock companion and require one run plus the named skip class |
| AX4/foreign-diagnostic | silently skip a foreign regular file | the foreign-entry fixture | add one and require its diagnostic |
| AX4/nonregular-diagnostic | read a FIFO | the bounded special-file fixture | create a FIFO under timeout and require a non-regular diagnostic without blocking |
| AX4/symlink-diagnostic | follow a symlinked record | the symlink fixture | link a valid record and require a symlink diagnostic, not a duplicate run |
| AX4/malformed-diagnostic | skip malformed JSON | the malformed fixture | add malformed bytes and require their diagnostic beside healthy rows |
| AX4/unreadable-diagnostic | abort on one unreadable record | the unreadable fixture | deny read where supported and require its diagnostic plus healthy rows |
| AX4/healthy-row-survives | return on the first hostile entry | the mixed-state fixture | combine hostile and healthy entries and require every row/disposition |
| AX4/entry-cap | enumerate beyond the named cap | the cap fixture | create cap-plus-one entries and require bounded output |
| AX4/at-least-marker | report an exact count after truncation | the cap fixture | exceed the cap and require the at-least marker |
| AX4/repeatable-bytes | retain filesystem iteration order | the idempotency fixture | invoke twice in fresh services and require byte equality |
| AX4/help-word | reject root `help` | the CLI root-help fixture | invoke `bench spec build help` and require catalog/0 |
| AX4/help-long | reject root `--help` | the CLI root-help fixture | invoke `bench spec build --help` and require catalog/0 |
| AX4/help-short | reject root `-h` | the CLI root-help fixture | invoke `bench spec build -h` and require catalog/0 |
| AX4/root-home | restore missing-argument usage/2 | the old/new root fixture | invoke no operation and require live home/help/0 |
| AX5/abandon-apply | omit the abandon apply action | the abandon plan renderer test | plan through the real service and require one exact `--apply` action |
| AX5/reclaim-apply | emit two or zero reclaim actions | the reclaim plan renderer test | plan through the real service and require exactly one action |
| AX5/fixed-slug | replace target with `<slug>` | the authority fixture | plan for `build demo` and require the quoted fixed slug |
| AX5/fixed-fingerprint | replace fingerprint with `<fingerprint>` | the authority fixture | require byte identity between plan fingerprint and apply argv |
| AX5/abandon-preimage | drop one worktree or ref fact from `abandonmentFacts` | the public abandon authority test | vary that input, plan twice, and require distinct fingerprints plus rejection of the first on apply |
| AX5/reclaim-preimage | drop one ref identity or disposition from `reclamationFacts` | the public reclaim authority test | vary that input, plan twice, and require distinct fingerprints plus rejection of the first on apply |
| AX5/stale-abandon-replan | carry a rejected abandon fingerprint | the stale-abandon fixture | apply drifted authority and require fresh plan action without `--apply` |
| AX5/stale-reclaim-replan | carry a rejected reclaim fingerprint | the stale-reclaim fixture | apply drifted authority and require fresh plan action without `--apply` |
| AX5/spent-reclaim-replan | advertise reuse after partial reclaim | the interrupted-reclaim fixture | force partial refusal, retain spent receipt, and require only fresh plan action |
| AX6/operation-totality | remove one operation from the fixture ledger | the ledger-totality test | join accessor members to cases and require no missing member |
| AX6/state-totality | omit one typed cell | the ledger-totality test | enumerate operation/state cells and require fixture or exact disposition |
| AX6/root-help-totality | remove one root/help spelling | the family-variant test | enumerate root, `help`, `--help`, and `-h` and require each identity |
| AX6/plan-apply-totality | omit one authority variant | the authority-variant test | enumerate plan/apply/stale/spent cases and require each identity |
| AX6/old-bytes | derive old bytes from new production output | the fixture-independence test | load checked-in old bytes and require their ledger digest |
| AX6/new-bytes | generate expected bytes at assertion time | the fixture-independence test | load checked-in new bytes and compare real output to them |
| AX6/help-delta | allow `next` change in a help-only case | the named-delta comparison | mutate `next` and require out-of-delta red |
| AX6/next-delta | allow another primary cell in a next-only case | the named-delta comparison | mutate it and require out-of-delta red |
| AX6/home-delta | classify an operation as family-home replacement | the named-delta comparison | change the class and require identity mismatch |
| AX6/unnamed-byte-refusal | alter punctuation outside an allowed block | the byte-preservation comparison | mutate the subject and require exact unexpected-byte red |
| AX7/operation-axis | remove one accessor member from conformance | the AXI disclosure check | enumerate production and require every identity |
| AX7/state-axis | substitute a hand-maintained state list | the typed-axis identity test | compare conformance cells with lifecycle classifications |
| AX7/unclassified-cell | add a constructor without disposition | the AXI disclosure check | enumerate the owner and require named unclassified-cell red |
| AX7/structured-stdout | route one refusal only to stderr | the command observation harness | drive it and require structured stdout |
| AX7/success-exit | change success/no-op to exit 1 | the taxonomy test | execute fixtures and require exit 0 |
| AX7/refusal-exit | change lifecycle refusal to exit 2 | the taxonomy test | execute typed refusal and require exit 1 |
| AX7/usage-exit | accept malformed argv as lifecycle refusal | the taxonomy test | execute unknown operation/missing flag and require exit 2 |
| AX7/operation-help | drop one operation help spelling | the grammar-derived help test | execute every operation/spelling pair and require catalog/0 |
| AX7/root-help | drop one family-root spelling | the root-help test | execute all three and require identical catalog/0 |
| AX7/help-envelope | omit help from one applicable cell | the AXI disclosure check | execute every classified cell and require terminal help table |
| AX7/registry-binding | omit or misbind the executable entry | existing conformance meta-tests | run `go test ./internal/conformance` and require identity/tier/input/order/binding agreement |
| AX7/profile-binding | omit or drift the profile row | existing profile-agreement test | compare registry with `projects/benchkit.md` and require the current input source |
| AX8/useful-action-removal | remove one useful action | the useful-action mutation test | mutate finished derivation and require missing-command red |
| AX8/known-value-placeholder | replace a fixed value with a placeholder | the known-value mutation test | mutate and require literal identity red |
| AX8/unknown-value-guess | replace an open placeholder with a guess | the placeholder mutation test | mutate and require open-input red |
| AX8/dropped-fixed-flag | remove `--full` from a fixed action | the fixed-flag mutation test | mutate and require complete invocation red |
| AX8/stale-fingerprint | reuse superseded authority | the fingerprint mutation test | substitute prior fingerprint and require mismatch red |
| AX8/prose-advertisement | mark orchestration prose executable | the prose mutation test | mutate kind and require construction refusal |
| AX8/unquoted-hostile-value | bypass shell quoting | the hostile serialization mutation test | mutate serializer and require exact command red |
