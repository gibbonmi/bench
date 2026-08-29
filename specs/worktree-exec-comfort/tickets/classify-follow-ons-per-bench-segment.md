# Classify follow-ons per Bench segment

Blocked by: none
Writes: internal/benchguard/benchguard.go, internal/benchguard/benchguard_test.go, internal/shellcommand/shellcommand.go, internal/census/census_test.go, internal/systemtest/bench_follow_on_test.go, .bench/BENCH.md, CONTEXT.md

## What to build

`Classify` changes from a stream-wide scan to a span-scoped rule. The
classifier finds the first Bench-headed simple command and reads its span. When
that command is `bench worktree exec`, the classifier allows a heredoc
redirection inside the span. It also allows a control operator after the span
when the operator is `;` or `&&` and every later simple command is non-Bench.
The classifier refuses everything else, and any other Bench-headed command
keeps the stream-wide rule.

`Classify` returns the verdict, the Bench segment's projected words, and the
adjacent operator token. The refusal message is the fixed sentence, one space,
then `segment=<words> operator=<token>`. The words join by single spaces. The
token is a redirection inside the span when one exists, else the control
operator before the span, else the one after. A redirection prints as its
source text, with the fd digits, the operator, and the operand joined without
spaces. The exec allowance covers a heredoc only; a here-string `<<<` refuses.

The refusal's first line stays
`BLOCKED: Bench response is bounded, complete, and self-contained`, because the
hook's pin and the system suite read that sentence. The census keeps skipping a
Bench-headed call, so an exec follow-on records no entry.

This ticket also carries the guidance for stories 42 and 43, because both
sentences describe this classifier and no other ticket owns them. The operating
guide's contract sentence in `.bench/BENCH.md` names the exec exception, and
`CONTEXT.md`'s `follow-on` term states the same exception. The review round
grades both against the rows below.

Package invariant for `internal/benchguard`: one rule decides the verdict and
the named token together, and the fixed sentence has one source. A second
verdict path, or a second copy of the sentence, drifts from the hook's pin.

## Acceptance

- [ ] G1: `bench worktree exec L -- cp a b; cp b a` and
      `bench worktree exec L -- cp a b && rg -n x b` are allowed.
- [ ] G2: `cp a b && bench worktree exec L -- true` is refused.
- [ ] G3: `bench worktree exec L -- true | cat`, `... || echo x`, `... &`,
      `... > out`, and `... <<< x` are each refused.
- [ ] G4: `bench worktree exec L -- true; bench maps` is refused.
- [ ] G5: `bench worktree exec L -- cat <<'EOF'` with the body line
      `bench gate` is allowed.
- [ ] G6: `bench gate <<'EOF'` with the body line `input` is refused.
- [ ] G7: `cat <<'EOF'` with two non-bench body lines is allowed.
- [ ] G8: the refusal for `cat a && echo x; bench maps` ends with
      `segment=bench maps operator=;`.
- [ ] G9: the refusal for `bench gate 2>&1` ends with
      `segment=bench gate operator=2>&1`.
- [ ] G10: the refusal for `cp a b && bench worktree exec L -- true | cat`
      names `operator=&&`.
- [ ] G11: every refusal's stderr still starts with
      `BLOCKED: Bench response is bounded, complete, and self-contained`.
- [ ] G12: the hook process allows `bench worktree exec label -- cp a b; cp b a`
      and `bench worktree exec label -- cat <<'EOF'` with a body.
- [ ] G13: `bench worktree exec L -- sed -i x <pool path>; cp a b` records no
      census entry.
- [ ] G14: `bench gate; cp a b` is refused.
