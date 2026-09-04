# Rewrite the push verdict for explicit refspecs

Blocked by: none
Writes: internal/gitguard/gitguard.go, internal/gitguard/verdict.go, internal/gitguard/verdict_test.go, internal/gitguard/gitguard_test.go, tests/canary/canonical-path-owner/second-derivation, tests/canary/injected-ports/unregistered-port
Covers: PG1, PG2, PG3, PG4, PG5, PG6, PG7, PG8, PG9, PG10, PG11, PG12, PG13, PG14, PG15, PG16, PG17, PG18, PG19, PG28, PG32, PG33, PG34, PG39

## What to build

Verify the premise first: read the `push` case in `classify`
(internal/gitguard/verdict.go), which returns the label unconditionally today.
Read `Checker`, `denyTable`, `denyLabels`, and `BlockMessage` in
internal/gitguard/gitguard.go. Read the `refYes` and `refNo` fakes and
`TestClassifyVerdicts` in internal/gitguard/verdict_test.go. Read the helpers
`contains`, `hasAny`, `freeArgs`, and `anyShortFlagHas` in verdict.go, and
compose them rather than write new ones.

Give `Checker` three more injected facts: `DefaultBranch func() (string, bool)`,
`CheckedOut func() (string, bool)`, and `BareDestination func() (string, bool)`.
The guard ticket wires all three from internal/git. The fakes in this package
set them by hand. This ticket adds no production caller of internal/git.

Replace the one `push` row in `denyTable` with seven rows. The labels are
`git push to the default branch`, `git push --force`, `git push --delete`,
`git push --all`, `git push --mirror`, `git push --tags`, and
`git push with an unresolved destination`. Give the table an advice column.
The unresolved row carries the sentence
`Name the remote and the branch: git push <remote> <branch>.`, and
`BlockMessage` appends the advice of the label it names. The table stays the one
source of every label and every advice sentence.

Write a `pushVerdict` function for the `push` case. It skips each option that
takes a separate value, so `-o`, `--push-option`, `--repo`, `--receive-pack`,
and `--exec` never give up their value as the remote. It reads the first free
arg as the remote and every later free arg as a refspec.

A refspec that starts with `+` is force. A refspec with an empty source is a
deletion. A refspec `src:dst` targets `dst`, and a refspec `src` targets `src`.
A refspec `HEAD` targets the checked-out branch. When `CheckedOut` reports no
branch, a `HEAD` refspec denies with the unresolved label.

The function strips a `refs/heads/` prefix from the destination, then compares
it to the default branch. A destination under another `refs/` namespace is
never protected. With no refspec, the destination comes from `BareDestination`.
No option exempts a push, and `--dry-run` is not an exemption. An absent
default branch and an absent destination both deny with the unresolved label.

Cover the bare-push composition rows in the sibling guard ticket, not here.

## Acceptance

- [ ] `git push origin topic` with a fake default of `main` returns the allow verdict.
- [ ] `git push origin main`, `git push origin HEAD:main`, and `git push origin refs/heads/main` each return the default-branch label.
- [ ] `git push origin main:topic` and `git push origin refs/tags/main` each return the allow verdict.
- [ ] `git push origin HEAD` returns the allow verdict on a fake checked-out `topic` and the default-branch label on a fake checked-out `main`.
- [ ] `git push origin topic main` returns the default-branch label.
- [ ] `git push origin HEAD` with a fake `CheckedOut` that reports no branch returns the unresolved label.
- [ ] A fake default of `master` denies `git push origin master` and allows `git push origin main`.
- [ ] `git push origin topic` with no resolvable default returns the unresolved label.
- [ ] `git push --dry-run origin main` returns the default-branch label.
- [ ] `git push -f origin topic`, `git push --force origin topic`, `git push -fu origin topic`, `git push --force-with-lease origin topic`, `git push --force-with-lease=topic origin topic`, and `git push origin +topic` each return the force label.
- [ ] `git push --delete origin topic`, `git push -d origin topic`, and `git push origin :topic` each return the delete label.
- [ ] `git push --all origin`, `git push --mirror origin`, and `git push --tags origin` each return their own label.
- [ ] `git push -o ci.skip origin topic` returns the allow verdict.
- [ ] `env git push origin main`, `bash -c 'git push origin main'`, and `git status && git push origin main` each return the default-branch label.
- [ ] `git push origin topic && ls` returns the allow verdict.
- [ ] `BlockMessage` names each of the seven push labels inside backticks.
- [ ] `BlockMessage` for the unresolved label ends with `Name the remote and the branch: git push <remote> <branch>.`
- [ ] Self-probe: read the first free arg as a refspec, and report the `git push origin main` row red.
