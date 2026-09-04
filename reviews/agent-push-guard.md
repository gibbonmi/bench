# Review pickup: agent-push-guard

Frozen base `5044b03dedd6a598fea6e6dc885f761e02def193`, reviewed tip
`063223038c07e8f72652a1277db0c8ccf790a3af`. Three axes ran at `opus` low.
Raw findings: 11. Repair targets after collapse: 7. Six targets are repaired
on the source. One finding stays open for the reviewer.

## Standards

Count: 1 open. Worst: the guard wiring is derived twice.

- `internal/gitguard/checker_junction_test.go` (`realChecker`) composes the same five facts that `guardGit` composes in `cmd/bench/main.go`. A sixth fact added to one side leaves the other side silent. A collapse needs an exported constructor and a new import edge from `internal/gitguard` to `internal/git`, which the spec's proof checklist excludes. Disposition: ask-user. Reported, not repaired.

## Spec

Count: 0 open. The redirect denial, the PG28 seam, the `--force-if-includes` Won't-handle line, and the `git push main` edge are repaired on the source and recorded in the spec.

## Coverage

Count: 0 open. The `xargs` denial, the `@` synonym, the `heads/` prefix, and the redirect are repaired on the source with rows PG42 to PG46.
