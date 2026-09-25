# 10. Apply the 4096-byte bound to each bounded response

Blocked by: 2-bound-every-public-response.md
Writes: internal/bounds/bounds.go, tests/canary/package-core-guard/bounds-duplicate-owner, internal/responsebound/, internal/systemtest/exec_bound_test.go
Covers: BO63, BO64, BO65, BO74, BO75, BO76

## What to build

Add the byte value 4096 to the production policy registry in `internal/bounds`, beside the line value. The owner treats a response as over-bound when it has more than 10 lines or more than 4096 bytes. The owner starts the spill at line 11 or byte 4097.

The owner derives the line cut as the byte value divided by the line value, 409 bytes. A projected line longer than the cut, without its newline, keeps its first 409 bytes. If that prefix ends inside a well-formed UTF-8 rune, cut back to the start of that rune. Otherwise, cut at 409 bytes. The spill line reports the number of cut lines in `cut_lines`.

After the spill starts, each retained head and tail line holds at most 409 bytes in memory. The owner tests build their 4096, 4097, and 409 fixtures from the registry constants. Do not edit `specs/session-context-queries`. BO73 is a final reconciliation check, not a task of this ticket.

## Acceptance

- [ ] A 3-line response of 5000 bytes spills with `omitted_lines=0`, and its printed bytes stay within 4096.
- [ ] A 10-line response of exactly 4096 bytes prints unchanged and creates no spill file, and one more byte spills.
- [ ] A head line of 999 bytes of 3-byte runes prints 408 bytes, and the spill line holds `cut_lines=1`.
- [ ] An exec child that prints one 64 MiB line with no newline returns its exit code, and the spill file holds the 64 MiB.
- [ ] After one 64 MiB write with no newline, the owner retains at most 4096 plus 9 times 409 bytes.
- [ ] The bounds-policy guard stays green, and no line that this ticket adds to a Go file outside `internal/bounds` states 4096 or 409.
