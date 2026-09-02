---
name: craft-gate
description: How to author the oracle — adding, editing, weakening, or removing a gate check, choosing its fail posture, and proving it bites. Use whenever touching .bench/gate.sh or any project gate, wiring an enforcement hook or guard, or scaffolding a gate in /bench-setup-repo. Reach for this before changing what decides "done".
index: adding, weakening, or removing a gate check / authoring the oracle
---

# Authoring the oracle

The gate is the only authority on done — invariant #1 rests on it. This makes
a gate check the highest-leverage code in a repo: every shift, every delegate,
every review inherits its strength. Write each check as code that earns more
trust than you do.

## Prove it bites

A check you have never seen fail is a hope, not a check. Before a new check
lands, break the thing it guards, watch the gate go red **with the check's own
message**, then revert. A retained kit fixture earns its place through an
ordinary mutation test. That test calls its exact registered owner, requires
the planted diagnostic, restores the subject, and requires that diagnostic to
disappear. A linked repo receives canary inventory validation from Bench and
owns equivalent planted-reason proof in its native tests. Without a retained
fixture, the one observed red is the minimum, stated in the spec or commit.

A new check is complete when you observed one red with the targeted message,
green after the fix, and a red reproducible from the diff. Use a fixture, or
a documented break-it command. The author asks which single edit defeats a
new check while the gate stays green.

## Attribute every failure

A red gate must name its cause without archaeology.

- One distinct message per failure mode; never reuse a message across checks.
- Guard each check on the presence of the surface it tests (skip when
  absent). A minimal fixture then fails only for its planted reason plus
  honest noise. Attribution is by substring, not isolation.
- Record failures and continue rather than exiting on the first, so one run
  reports the whole surface.

## Run the real path

A check exercises the actual command, script, or parser against a fixture and
asserts the external verdict — exit code, output, a file appearing. It never
reimplements the checked logic. A second derivation can disagree with the real one:
green while the product is broken, or red while it is fine.

A claim can live in docs while its enforcement lives in code — a file list, an
index, a deny surface. Generate one from the other, or grade both directions. An
advertisement without enforcement drifts, and an accident deletes an enforcement
nothing advertises. A check on a workflow output, config key, or environment
variable grades the producer and the consumer, or their binding, in the same change.

```
( cd "$tmp" && git init -q && bash "$root/bin/tool.sh" link )
[ -f "$tmp/.bench/BENCH.md" ] || err "fresh link did not install .bench/BENCH.md"
```
Good — runs the real installer in a hermetic repo and asserts the observable
outcome, with a message that names the failure.

```
grep -q 'cp .*BENCH.md' bin/tool.sh || err "link broken"
```
Bad — reimplements the installer as a grep of its source with a cause-free
message. It stays green through a refactor that breaks behavior, and goes
red on a rename that doesn't.

A grep-anchor is still legitimate as a cheap *tripwire* on prose (a command
file must keep naming a skill). It catches deletion, not decay — never
present one as a behavior check.

## Hermetic and fast

Fixtures run in throwaway temp repos with controlled inputs, never against
the live repo's mutable state, the network, or the clock. Same tree, same
verdict, every run. The gate runs on every shift iteration, so its runtime
taxes every loop; keep checks cheap, and push expensive proofs into bounded
fixtures.

## Choose the fail posture out loud

When a check's dependencies can be absent — a linter not installed, a lib or
config missing — what happens next is a decision, not an accident.

- **Fail closed** for enforcement: a guard that cannot load its rules refuses
  the action. An unguarded pass-through is silent de-enforcement.
- **Fail open** only for ergonomics-layer hooks that must never brick the
  workflow. Always add a one-line stderr warning, and enumerate the open
  cases in the header, so a new one is a visible edit.
- **Best-effort** (runs only when a tool is present) is legitimate for
  supplementary lint. Say so in a comment, and when the tool is present but
  the run fails, make sure the check errs rather than silently passes.

## Layer the gate

parse/validity → structure → conformance → behavior contracts → canary inventory.
The final layer validates a non-empty set of accepted fixture bindings; ordinary
native tests own the direct planted-reason proof. A new check joins its layer,
and a check another check depends on runs first.

## Keep the tripwire alive

The gate runs from the working tree, so the agent it grades can edit it. Keep
that edit loud: every retained kit fixture has direct ordinary-test proof, and
an empty, invalid, or unbound canary inventory is red. Linked repos keep the
same division of responsibility: Bench validates their inventory, while their
native tests prove their checks bite. A scaffolded gate's configuration sentinel
keeps it red until the project supplies real checks and bindings. To delete or
weaken this defense, follow the rule below for any weakening; it is never a
quiet step that makes a change pass. The threat this covers is the lazy
shortcut, not a determined adversary — the contract is loudness, not prevention.

## Weakening is a reviewer decision

When a check blocks a change you believe is legitimate, do not edit, relax,
or delete it as a step that makes that change pass. Stop and surface it
instead. A deliberate weakening or removal ships as its
own visible change with reviewer sign-off. That same change updates its
fixture or canary to match, never quietly inside a feature diff. This is invariant #1 read
from the author's side.
