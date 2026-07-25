# The smell baseline — glosses and fixes

One gloss and one fix per smell in `craft-review`'s Standards baseline
(Fowler, *Refactoring* ch. 3). Load this when a hunt needs more than the name.

- **Mysterious Name** — the name needs the code read to be understood; rename
  until it says what it's for.
- **Duplicated Code** — the code-standard defect, two sources for one fact;
  collapse to one.
- **Long Function** — a body that accreted past what one read can hold; extract
  until each piece says one thing.
- **Long Parameter List** — arguments multiplying at a call site; introduce a
  type before they start traveling together.
- **Global Data** — ambient state any code can mutate from anywhere; wrap access
  in one owner.
- **Mutable Data** — shared state read and rewritten from multiple sites;
  narrow who writes, or copy instead of mutate.
- **Feature Envy** — a function mostly reaching into another module's data;
  move it to where the data lives.
- **Data Clumps** — the same few values always traveling together; make them
  one object.
- **Primitive Obsession** — a domain concept passed as a bare string or int;
  introduce a value type.
- **Repeated Switches** — the same type-dispatch repeated across sites;
  collapse to one dispatch.
- **Shotgun Surgery** — one logical change forcing edits in many places;
  gather the pieces into one module.
- **Divergent Change** — one module edited for unrelated reasons; split it by
  reason for change.
- **Lazy Element** — an abstraction whose reason for existing has since left —
  a wrapper forwarding one call, a type with one field; inline it.
- **Speculative Generality** — hooks and parameters for a future nobody
  scheduled; delete until needed.
- **Message Chains** — `a.b().c().d()` walks through the object graph; hide
  the delegation behind the object the caller holds.
- **Middle Man** — a module that mostly forwards; let callers reach the real
  object.
- **Refused Bequest** — an inheritor ignoring most of what it inherits; swap
  inheritance for delegation.
