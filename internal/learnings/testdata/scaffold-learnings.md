# Learnings — usage journal

Append one entry when you deviate from the workflow, or when you make a
judgment call you're unsure about. Append one too when you catch a
should-have-asked in hindsight, or when you catch yourself assembling the
same ad-hoc check twice. That is a codification candidate; name the `bench`
subcommand it wants to be. You capture; the reviewer decides.

`/bench-drain` verdicts every open entry in its reviewed batch diff.
Work-shaped and rule-shaped entries become roadmap items; rule-shaped ones
build later under the synthesis discipline. The rest are dismissed with one
line of why. A resolved entry leaves this file, and its verdict is recorded
in the drain commit. The journal holds open entries only; history lives in
git. Never rewrite a kit rule yourself; that is the whole point of capturing
here instead.

Format per entry:

## <date> - <short title>  [open]
- **What happened:** ...
- **Right behavior:** ...
- **Proposed rule change:** ... (or "none")

An entry leaves this file only via /bench-drain.

<!-- entries below -->
