---
name: craft-design-system
description: How to consume a project's design system when building UI. Use whenever building or modifying a UI component, applying a color/spacing/type/motion value, or composing a component. Reach for this on any UI work — pull values from tokens, compose canonical components, never hardcode or regenerate.
index: any UI work
index-note: your project's design source
---

# Design system

When a project has a design system, it is the source of truth for every **visual**
decision. It lives in a **design source you control** — a separate repo, a package,
or a pinned path (the project profile names which) — materialized as two kinds of
artifact any agent can read:

- **Tokens** — the value scale: color, spacing, type, radius, motion. The only
  legal source of a visual value.
- **Canonical components** — the blessed component inventory. The only legal source
  of a component.

Because the source of truth is committed artifacts, not a tool's live state, this
skill is the *visual* layer and is harness-agnostic. If the project also has an
interaction-layer UI skill (its own `*-ui` skill), that governs behavior; both gate
UI work and neither replaces the other.

## Rules (these are gate checks, not suggestions)

1. **Every visual value comes from a token.** No raw hex, no hardcoded px, no
   literal durations. If the value you need has no token, **stop** — that's a
   design decision, not a build decision. It goes to the design source, not into
   the component.

   ```css
   /* good — the value has a home in the design source */
   color: var(--color-accent);

   /* bad — the hex has no home; when the palette moves, this stays behind */
   color: #3b82f6;
   ```

2. **Compose canonical components; never regenerate or restyle them.** If a needed
   variant doesn't exist, **stop** and get the variant added to the design source
   first, then build against it.

3. **The design system is current-state design documentation.** The ADR rule
   applies: the design source records what the design *is*, not how it changed.
   When it changes, re-publish and re-pin; git holds the history.

4. **A green logic suite is necessary but not sufficient for UI.** The
   design-conformance check (raw-value scan, canonical-import check, plus any
   project-specific rules) runs against the design source and is part of the gate.
   UI that passes tests but violates tokens is not done.

## Authoring a change — seamless across harnesses

This is the part that has to survive a harness switch. When a token or component
variant is missing, you author it in the **design source**, and the authoring
surface is interchangeable — it is never a dependency of the workflow:

- **In a Claude session:** Claude Design is a fast canvas for iterating the token or
  variant visually before you commit it to the design source.
- **In Codex, another model, or any other harness:** edit the design source
  directly — it's tokens and components in files. No Claude Design, no Claude
  required.

Either way the output is identical: a committed change in the design source, which
the project re-pins. Consumption always reads the committed artifacts, so **moving
from Claude + Claude Design to another harness changes only how you draw a new
token — not how the build consumes it, not the composition rules, not the
conformance check.** That is what makes the transition seamless: the design tool is
an authoring convenience layered on top, and the project never reaches into it live.

## Working method for a UI shift

1. Read the design source's tokens and component inventory before generating.
2. If a token or variant is missing, stop — add it in the design source (via Claude
   Design or directly), commit, re-pin, then resume.
3. Compose canonical components; reference tokens for every value.
4. Run the design-conformance gate (and the project's screenshot/interaction loop if
   it has one). Done is green on both, not either.
