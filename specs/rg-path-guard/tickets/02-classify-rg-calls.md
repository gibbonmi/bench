# Classify each rg call as pathless, substituting, or permitted

Blocked by: 01-mark-live-substitution.md
Writes: internal/rgguard (new), internal/conformance/guard_projection_owner_test.go, internal/conformance/guard_projection_rgguard_test.go (new), internal/shellcommand/projection.go, internal/shellcommand/projection_test.go, internal/gitguard/scan.go, internal/benchguard/benchguard.go
Covers: R08, R09, R10, R11, R12, R13, R14, R15, R16, R17, R18, R19, R20, R21, R22, R23, R24, R25, R26, R27, R28, R29, R30, R31, R32, R33, R34, R35, R36, R37, R38, R58

## What to build

A new package `internal/rgguard` decides whether one Bash command holds an `rg` call that the guard refuses. The package starts no process and reads only the command string.

The package exports these values for the hook ticket (`03-add-guard-rg-verb.md`):

- `Classify(command string) Verdict`. The zero `Verdict` permits the command.
- `Verdict.Blocked` and `Verdict.Message() string`. `Message` returns the complete refusal text for a blocked verdict.
- `CommandFromEnvelope`, which reads the PreToolUse envelope through `shellcommand.DecodeEnvelope` and returns the command field or an error, like `benchguard.CommandFromEnvelope` (`internal/benchguard/benchguard.go:25-44`).

`Classify` parses the command with `shellcommand.Parse` and reads each simple command. It finds the executable with `ResolveRoutinePrefix` and the words with `ProjectCommandWords`. A span is an `rg` call when the basename of the resolved executable is `rg` and `Executes` is true.
It reads one `-c` child level through `ShellChildren` and the exec child through `WorktreeExec`. It does not recurse below that level.
`WorktreeExec` requires that its caller has already accepted the executable as Bench (`internal/shellcommand/projection.go:46-48`). The name half of that test has two private copies today: `benchExecutables` (`internal/gitguard/scan.go:76-78`) and the basename test in `isBench` (`internal/benchguard/benchguard.go:309-324`). Add one exported predicate, `shellcommand.IsBenchName(word string) bool`, in `projection.go`: true when the basename is `bench` or `bench.sh`. Test it in `projection_test.go`. Migrate `gitguard` to it and remove its map. Let `benchguard.isBench` call it for its name half and keep its resolver half unchanged. `rgguard` calls it before `WorktreeExec`. No Git or Bench verdict changes.
Do not write a shell-name literal (`sh`, `bash`, `zsh`), a regexp that matches `-c`, or a comparison of adjacent `worktree` and `exec` words. The projection-owner check (`internal/conformance/guard_projection_owner_test.go:104-134`, `:186-193`) reds each of those in a governed package.

For each `rg` call, apply the verdict order of the spec exactly:

1. If the stream is `Unlexed`, permit.
2. If a token of the span, or the operator token at the span end, has `Substitution` true, refuse with the substitution reason. For an exec child, read the marks of the outer span tokens after `--` and of the outer span-end operator, because `WorktreeExec` returns plain argv. Do not join or parse that argv again.
3. If the span is uncertain, permit. The uncertain kinds are `ViaXargs`, a `<(` or `>(` operator at the span end, an inline `RIPGREP_CONFIG_PATH=` assignment, and an option word that the table does not list.
4. If a `|` or `|&` operator comes directly before the span, the span holds an input redirection (a `Redirection` token whose text starts with `<` after an optional descriptor digit), or the call holds an explicit stdin request (`-` as a path, or `-f -`), permit. For a child, the outer span's feed also feeds the first `rg` span of the child.
5. If the call names a no-pattern mode, permit.
6. If the call has a pattern and no path operand, refuse with the pathless reason. Otherwise permit.

The option table is the one source of ripgrep 15.1.0 grammar in the tree. Use a map or a switch, not a regexp. Run `rg --help` (the installed version is 15.1.0) for the source text. The table lists every option word of that help with its arity: the 14 short value letters `e f E m j g d t T A B C M r` with their long pairs, the alias `--maxdepth`, the 21 long-only value options that the spec lists, every boolean option, and every `--no-*` negation that the help names in a heading or in prose.
A long value option takes its value after `=` or from the next word, with no prefix abbreviation. In a short bundle, the first value letter takes the rest of the bundle, or the next word when the bundle ends; a value after `=` loses the `=`. A bundle that holds `h` or `V` is a mode. The modes `--files`, `--type-list`, `-h`, `--help`, `-V`, `--version`, `--pcre2-version`, and `--generate` need no path.
After `-e` or `-f`, every positional is a path. Otherwise the first positional is the pattern. After `--`, every word is a positional.

The package owns both refusal sentences and one message format. The pathless message names the missing path and the `.` repair. The substitution message names command substitution and advises single quotes. Each message ends with `segment=` and the projected words of the call, like `internal/benchguard/benchguard.go:78`.

Append `internal/rgguard` to `guardProjectionGoverned` (`internal/conformance/guard_projection_owner_test.go:29`). That file has 373 of its 400 lines, so put the new R37 test in `internal/conformance/guard_projection_rgguard_test.go`. The test plants a private envelope decode in a temporary `internal/rgguard` tree and expects the owner check to red.

Key the verdict table `TestClassifyRgCalls` and the message table `TestClassifyRgMessages` in `internal/rgguard/rgguard_test.go` by coverage row. The expectations come from ripgrep 15.1.0 help and probes and from POSIX quoting, so record one omission red under the `AGENTS.md` exception: remove one value letter from the option table and record the R13 red, the restore, and the green.

## Acceptance

- [ ] Every Classify row from R08 to R36 and R58 passes with the exact verdict and reason that its row names.
- [ ] The substitution message contains `command substitution` and `single quotes` and ends with `segment=` and the call words (R25).
- [ ] The pathless message for `rg foo` contains `no path` and `` `.` `` and ends with `segment=rg foo` (R28).
- [ ] A private envelope decode planted in `internal/rgguard` reds the projection-owner check (R37).
- [ ] `bench test --check package-core-guard` passes on the real package.
- [ ] No Go package other than `internal/rgguard` declares ripgrep option words (R38, review-owned).
- [ ] `Classify` starts no process and reads no file.
- [ ] `shellcommand.IsBenchName` is the only executable-name test for Bench in the guard packages (`gitguard`, `benchguard`, `rgguard`), and `bench test --package ./internal/gitguard --package ./internal/benchguard` passes unchanged.
