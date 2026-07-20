# Security

## Trust model

Bench runs as the invoking user and treats its gate as trusted project code. Its hooks
and guards reduce honest mistakes; they are advisory controls, not a security boundary
against a malicious or compromised process. Security boundaries come from the
harness or operating-system sandbox and from server-enforced repository protections,
including protected branches and required review.

For an inventory of every repository-controlled path that prompt, environment,
file, log, network, cache, and retention data reaches, see `DATA_HANDLING.md`.

## Network egress

Bench makes outbound requests in three situations:

- `bench models` queries the OpenAI and Anthropic model-list APIs when their
  corresponding API keys are present.
- Binary self-repair reads package metadata and an artifact from the npm registry.
  Set `BENCH_NO_REPAIR` to any non-empty value to disable this egress.
- Worktree acquisition makes a best-effort `git fetch origin` before creating a pooled
  worktree.

## Reporting a vulnerability

Report suspected vulnerabilities privately to gibs.mikej@gmail.com. Include the
affected version, impact, reproduction details, and any known workaround. Please allow
reasonable time for investigation and remediation before public disclosure; the
maintainer will acknowledge the report and coordinate a disclosure timeline.
