# Splitting method — when the structure gate fires

`bench shift` checks touched files after the implementation loop is green and
runs a bounded refactor pass; projects may also wire `bench structure` into the
gate to hard-block the PR boundary. Either way, how you split decides whether
you've reduced debt or just moved it:

- **Split along responsibility, not line count.** Find the seam already latent
  in the file — the cluster of functions with one reason to change — and lift
  it into its own module behind a small interface. The line count drops as a
  side effect of a real boundary, never as the goal.
- **Apply the deletion test before every split.** Would extracting this
  cluster *concentrate* complexity behind one interface, or just *move* it to
  another file? Only "concentrates" is a real deepening; "just moves it" is
  fragmentation wearing a refactor's clothes. (The deletion test is borrowed
  from Pocock's `improve-codebase-architecture` — a richer external version of
  the same review if his skills happen to be installed; an optional upgrade,
  never a dependency.)
- **Never fragment to beat the number.** Slicing a cohesive file into `part_a`
  / `part_b` scatters one concept across files with no interface between them.
  If you can't name the responsibility the split isolates, don't split — the
  file may legitimately be one deep module; propose the budgets grant
  `SKILL.md` describes instead. Raising the global cap (`BENCH_MAX_LINES`)
  weakens the check everywhere and is the wrong tool for one file.
- **A crowded directory is an ungrouped module.** Thirty files in one dir
  means several modules are hiding in a flat namespace. Group related files
  into a package with a clear entry point; the package *is* the seam.

Bench's **deep pass** is this skill applied broadly — run `bench structure`
(tighten the budgets if needed) and work the deletion test across every flagged
module; it's fully self-contained, nothing outside the kit required.

The deep-module rule applies throughout: prefer a few files with real
interfaces over many shallow ones. The gate measures the symptom; this is the
cure.
