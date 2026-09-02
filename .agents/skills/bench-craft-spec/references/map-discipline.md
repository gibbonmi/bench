# Coverage-map discipline

Charged from `craft-spec` when the author writes or audits an acceptance coverage
map. Each rule below settles one question the map must answer. `craft-spec` keeps
the explore reads, the fence rules, and the review rubric.

## Before the map locks

- Name the enumeration source for a count. Do not copy an implementation-derived
  count into the spec, because the Standards axis grades a restated count.
- A spec that retires a symbol enumerates every production consumer and every
  pinning test before it slices tickets.
- A spec that deletes or moves literal bytes runs one repo-wide search for
  those bytes. The search covers the canary fixtures, `tests/`, and
  `internal/conformance`.
- A compile-flag spec runs one whole gate with its flag, then censuses every red,
  before it slices tickets.
- A claim that something is or is not a tree convention cites the check that
  enforces it. A sample of artifacts is not evidence, because an absence in the
  sample does not show an absent rule.
- The flagged-additions list sits under Further notes before the first review charge.
- The source-sentence-to-row table sits under Further notes before the first review charge.
- Before the first review charge, Further notes carries a pre-review proof for each applicable class and records `none` for each absent class.
- It resolves cited symbols and import edges, quotes source-row clauses with all occurrences, names exact field labels, and lists all changed-function callers.
- When a new owner replaces copies, it names a red-capable row that fails if any copy survives.
- Each canary row and each conformance row traces to its executed root before the coverage map locks.
- The reader sweep lists each named consumer of the decision fact.
- The reader sweep lists each helper that a named consumer calls directly.
- A deeper callee joins the reader sweep only when the callee reads the decision fact.
- Each shared reader in the reader sweep takes an exact ownership fence.
- The reader sweep names the shipped-surface claim words, because `package-core-guard` reds a claim word beside a repo-only path.

## Per row

- A universal claim names its authoritative inventory, its enforcement seam, and
  one omission mutation that turns the seam red.
- A spec justified by preserved behavior gives an exit row that names a
  differential run. The run covers the old and the new implementation over one
  enumerated input family. A pre-existing suite cannot protect behavior it never
  asserted.
- A row that promises output "names the command" pins the exact command token,
  never the command family.
- A row for a new standing check enumerates the in-flight tree states the check
  tolerates. Draw those states from the same spec's own lifecycle, because the
  spec's close step creates one of them.
- A ticket that replaces a blunt check with a precise one gets a row for the parse
  surface the precision creates. A row inherited from the blunt check may not cite
  a surface that did not exist when someone wrote that row.
- A row that names a failure mechanism traces the message to its producer. It also
  confirms the claimed input reaches that producer. Evidence from an operator
  session does not establish reachability inside the tree.
- A row whose seam is the existing tests names the test function, and someone reads
  that function in the same session.
- A row that substitutes a package variable names the venue. A substitution in the
  test process reaches nothing inside a test that drives a real subprocess.
- An ordering promise gets a row where two refusals compete, and the row names the
  message that wins. A row that asserts only an absent mutation cannot tell the
  first read-only position from the last.
- A spec that changes an oracle or a publisher keeps three separate identities: the
  authority, the test subject, and the published subject. A candidate that grades
  and publishes itself is self-attestation.
- An output-shape row bounds or excludes the old output. A row that asserts
  presence alone lets the old output survive.
- An inline TOON header in a row matches the shape of the design table. A design
  that permits `cells[N]` forbids a row that fixes one count.
- A loader-derived row drives the live loader.
- A seam that names two concurrent gate runs in one repository first checks the
  gate execution lock, because the lock refuses the second authorization.
- A transaction-shaped spec gives three rows for its verification failures. The rows are persistence before the oracle runs, interruption inside the oracle, and persistence at the terminal step.
- Each in-scope edge-inventory promise, source promise, and fence-closure promise takes one red-capable row.
- An either-side predicate takes two rows, one side per row. One row that names both sides is not sufficient.
- Each named diagnostic state is addable or mutable in a fixture.

## In the edge inventory

- The edge inventory asks which tests swap a package variable.
- The inventory lists each fail-closed cleanup error and each refusal beside the
  happy paths. Each one gets one red-capable row before the first review charge.
- A kit spec names the audience each behavior serves: this repository, or every repository that links the kit. The inventory walks the absent-versus-empty pair for each directory the spec reads, because the two audiences can want different answers.
- Each excluded edge takes a Won't handle line that names a surviving in-scope caller.

## At ticket slicing

- The last ticket that touches a package carries the invariant for that whole
  package.
- A fence over the public help traces every inventory fixture that help forces.
- The author quotes each pasted operand in the delegate charge.

## At review

- A source that requires real-environment evidence refuses a fixture that only
  simulates that environment.
- A process-group timeout row names a descendant-survival oracle, never elapsed
  time alone.
- The review round demands one row for each listed addition, and it removes each unlisted addition.
