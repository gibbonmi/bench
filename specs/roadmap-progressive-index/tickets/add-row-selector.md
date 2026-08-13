# Add the fail-closed --row detail selector

Blocked by: add-context-envelope.md
Writes: internal/roadmap

## What to build

`bench roadmap --context --row <ID,...>` returns complete untruncated bodies
for exactly the named rows. An ID is well-formed when it matches the parser's
row-start recognizer (`[A-Za-z]+[0-9]+`), single-sourced from `ParseDocument`
— never a second FT-only regex. Well-formed but absent is unsatisfied intent:
exit 1 naming the missing ID with no rows emitted, including when mixed with
present IDs. Malformed IDs, an empty `--row` value, and `--row` with `--full`
refuse as usage, exit 2. Duplicates dedupe. The member's usage string becomes
the single unified form naming bare, `--context [--full]`, and `--row`, and
the selector document carries the terminal help block the blocker ticket
established.

## Acceptance

- [ ] Selected rows carry complete bodies with true `body_bytes`; unselected
      rows do not ride the result (covers PI6).
- [ ] Absent ID → exit 1 naming it, zero rows, mixed requests included
      (covers PI7).
- [ ] Malformed ID, empty value, and `--row` + `--full` → usage, exit 2
      (covers PI8).
- [ ] Each grammar's help resolves to the unified usage string naming all
      three forms (covers PI13); `--row` output decodes as one TOON document
      with terminal help (covers PI20).
