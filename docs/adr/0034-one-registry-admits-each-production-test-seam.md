# One registry admits each production test seam

Status: accepted

## Decision

A production test seam is a controlled dependency or fault input on a production execution path.
One audit owner derives every seam from the production Go source and its tests, and one registry admits each seam.
The audit reads the graded tree only, and it never reads the working checkout.

The audit derives five kinds of seam:

- a named port, which is a nonempty interface, a named function type, or a struct of function fields that production injects
- a package variable that a test replaces, including function values, collaborator structs, and scalar bounds
- a setter that changes such a variable, resolved through forwarders and aliases to its canonical owner
- a branch on the test runtime, found through the import binding and through forwarding predicates
- an environment control whose value a test sets and whose reader that set reaches

Each admission row names a production owner, an actual top-level consumer test, a concrete failure purpose, and isolation evidence.
An alias row names its canonical seam and repeats none of its facts.
The audit refuses a missing field, a consumer that does not exist or never reaches the seam, a duplicate identity, an unresolved alias, and a row whose seam is gone.
An unresolved environment identity refuses audit success.

A named port with a production producer names a test that drives the real producer.
A named port that production supplies only as nil is a fault port.
A fault port takes the fault-exception fields, and its consumer supplies a non-nil value at the injection site.
A production closure or forwarder that returns only nil or the result of a nil fault counts as nil.

An environment control is a seam only when a test controls it.
An in-process set must reach the reader through production calls, and a child-process entry must set the name.
A set that reaches no reader derives no seam.

## Consequence

A new override, setter, runtime branch, or environment control fails the gate until it has a reviewed row.
A row cannot certify safety by its words alone, because review and the consumer's failure test own the purpose and the isolation.
The audit does not trace a test's reach through a function value or through a child that inherits the environment, so review owns those consumers.
A hook with no consumer stays readable for compatibility until a reviewed retirement or new consumer evidence.

## Considered options

A name pattern such as a test suffix finds wrappers but misses direct replacements and aliases, so the audit resolves declarations and uses.
A fixed list of audited packages hides a new seam outside the list, so named-port derivation covers every production package.
A literal rule that treats every test-set environment name as a seam admits names that control nothing, so derivation follows reach.
