# One lifetime owner and process-user protection

Status: accepted

## Decision

Every production child process starts through one subprocess lifetime owner.
The owner applies the caller's normal-exit policy, either to drain the process
group or to preserve it, and it ends the group with a bounded shutdown when
the caller cancels or a deadline expires. A raw child exit status never
certifies cleanup. The owner reports every unresolved obligation, and the
caller keeps each retained resource. A registered conformance check holds the
census of production callers, and a new caller needs a census row.

A resource that a child can still use carries durable process-user
protection. The protection store for each resource lives under one store root
in temporary storage that is private to the effective user and to the current
boot. No environment value selects that root, and the store name is a digest
of the resource path and its binding. Each launch registers as a user of the
resources it inherits or opens. A start gate holds the child until each
registration is recorded, so a child that dies before the gate opens never
counts as a user.

Removal of a protected resource needs proof that every registered user is
resolved. One safe removal seals the store, observes every user, removes the
resource, and then deletes its store. The last release of a cache hold grades
every user first, and it seals and retires the generation only when every
user is safe. An absent store names no user in the
current boot. An unreadable store or store root retains the resource. Builds,
kit tests, caches, prospective bundles, release evidence, worktrees and shift
scratch use this one protocol.

Worktree exec keeps its separate output wait delay. When that delay expires,
exec closes only its own output, reports the forced close on stderr, and
returns the child's exit code.

## Platform scope

The lifetime owner and the protection store serve Linux and Darwin. Darwin
keeps its own process-group policy and gets no Linux parent-death guarantee.
Native Windows lifetime semantics are out of scope, because the protection
package does not build for Windows.

## Consequence

A killed or interrupted run leaves no resource that a live descendant still
uses deleted under it, and a later operation removes the resource after its
users end. A retired store keeps a record of its resource. When the safe
removal of a directory deletes a resource, the removal also deletes the
retired store of that resource. The last release of a cache hold that no
other process shares also deletes each retired cache store whose cache no
longer exists. A retired store that names a user stays.
