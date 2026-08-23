# Subprocess data exposure is bounded by name passlists and documented durability

Repository-controlled subprocesses must not receive data the reviewer never
scoped for them. The posture has four legs, each with a rejected alternative.

1. **Agent subprocesses launch from a static name passlist plus a committed
   opt-in.** The default list is process basics, the loop's own variables, and
   each shipped harness's documented auth/config families. A repo widens it only
   through the committed environment opt-in file. That file fails closed on any
   malformed or unknown content, because a silently-ignored opt-in is indistinguishable
   from a working one. A denylist posture was rejected: unlisted-but-sensitive
   names always slip a denylist.

2. **A default glob never straddles two families.** Each glob covers a namespace
   a single owner controls; a name whose prefix is shared with a foreign family
   is enumerated exactly. This binds every future addition to the default list.
   An exact set must not be "simplified" into a wider glob that admits a name
   nobody documented.

3. **The project gate keeps its manifest-declared closed subject.** The gate's
   environment is part of the verdict identity: PATH plus only the names the
   project manifest declares. A separate gate-class passlist was rejected. It
   re-solved an already-closed problem, strictly wider than what ships, and
   rewiring the verdict path reached into verdict-identity hashing for no added
   closure. A sentinel contract pins the closed subject instead. A gate that
   needs a variable declares it in the manifest, which is both the opt-in and
   the verdict-identity declaration.

4. **Durability is documented, never guessed at.** Passlists filter names, not
   values. A passlisted variable's value may still be sensitive, and that is a
   documented fact, not a defect. Content-based secret detection was rejected as
   a different discipline with its own false-positive policy. Instead, the root
   data-handling inventory describes every repository-controlled prompt,
   environment, file, log, network, cache, and retention path. A conformance
   check derives its variable listing from the same constants the enforcement
   uses, so the advertisement cannot drift from the enforcement.

Two supporting decisions ride on the same posture. The iteration prompt travels
on stdin end-to-end, so it never appears in a process listing. A dual-mode
argv-plus-stdin transition was rejected, because the argv path would stay alive
and testable-green while leaking. Durable records also carry an objective
identifier rather than its text. The full text lives only in a reviewer-only
worktree-lifetime file. Shift commit subjects keep the sanitized
reviewer-authored text instead, because a readable git history is a feature.

## Consequences

A session widening any subprocess environment, adding a default passlist entry,
or proposing secret scanning should cite this record. Changing a leg is a
reviewer decision that supersedes it. This record reopens per leg on new evidence:

- leg 1, a needed name that cannot reasonably ship as an opt-in;
- leg 2, a documented single-owner namespace that the straddle rule wrongly excludes;
- leg 3, evidence the manifest-declared subject leaks in practice;
- leg 4, a real incident that documented durability failed to prevent and scanning would have caught.
