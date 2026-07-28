# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, catch a should-have-asked in hindsight, or catch yourself assembling
the same ad-hoc check a second time (a codification candidate — name the `bench`
subcommand it wants to be). You capture; the reviewer
decides. `/bench-what-next` verdicts every open entry in its reviewed batch
diff: work-shaped and rule-shaped entries become roadmap items (rule-shaped ones
built later under the synthesis discipline), the rest are dismissed with one line
of why. A resolved entry leaves this file, and its verdict is recorded in the
drain commit. The journal holds open entries only; history lives in git. Never
rewrite a kit rule yourself — that is the whole point of capturing here instead.

Format per entry. Heading: `## YYYY-MM-DD — short title  [open]`

- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

## 2026-07-28 — decision maps should cite the sources they were researched from  [open]

- **What happened:** The closed `decisions/slice-unit.md` map carried claims
  from six read sources (the vendored to-tickets family, two tree files) with
  no pointers, until the reviewer asked for them — a fresh session picking up
  the map would have re-derived which files hold the evidence.
- **Right behavior:** A map's Handoff cites the primary sources each decision
  rests on, with one clause of what each fed, so pickup is one hop; paths
  flagged as drift-prone rather than trusted.
- **Proposed rule change:** Add a `## Sources` section to the
  `/bench-shape-idea` map template — read sources with what-each-fed, required
  when any ticket's answer rests on files outside the map. Consistent with the
  staged warrant rule (spec story 15): the map is where external-source claims
  become reviewer decisions, so the warrant belongs on the artifact, not only
  in conversation.
