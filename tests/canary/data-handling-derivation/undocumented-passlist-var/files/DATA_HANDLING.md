# Data handling

Broken canary fixture: the variable listing below documents every internal/env
passlist pattern except `VERTEX_LOCATION`, so the derivation conformance check
must fire that pattern's own targeted diagnostic.

## Environment

<!-- passlist:begin -->

| Pattern | Family | Official documentation |
|---|---|---|
| `PATH` | Process basics | POSIX |
| `HOME` | Process basics | POSIX |
| `USER` | Process basics | POSIX |
| `LOGNAME` | Process basics | POSIX |
| `SHELL` | Process basics | POSIX |
| `TMPDIR` | Process basics | POSIX |
| `TERM` | Process basics | POSIX |
| `COLORTERM` | Process basics | terminal convention |
| `LANG` | Process basics | POSIX |
| `LC_*` | Process basics | POSIX |
| `XDG_*` | Process basics | XDG Base Directory |
| `BENCH_*` | Bench namespace | this repository |
| `ANTHROPIC_*` | Claude Code | https://docs.claude.com/en/docs/claude-code/settings |
| `CLAUDE_CODE_*` | Claude Code | https://docs.claude.com/en/docs/claude-code/settings |
| `CLAUDE_CONFIG_DIR` | Claude Code | https://docs.claude.com/en/docs/claude-code/settings |
| `API_TIMEOUT_MS` | Claude Code | https://docs.claude.com/en/docs/claude-code/settings |
| `CODEX_*` | Codex | https://github.com/openai/codex/blob/main/docs/config.md |
| `RUST_LOG` | Codex | https://github.com/openai/codex/blob/main/docs/config.md |
| `SSL_CERT_FILE` | Codex | https://github.com/openai/codex/blob/main/docs/config.md |
| `OPENCODE_*` | opencode | https://opencode.ai/docs/ |
| `OPENAI_API_KEY` | opencode | https://opencode.ai/docs/ |
| `AWS_*` | Claude Code — Bedrock | https://docs.claude.com/en/docs/claude-code/amazon-bedrock |
| `GOOGLE_*` | Claude Code — Vertex | https://docs.claude.com/en/docs/claude-code/google-vertex-ai |
| `GCLOUD_PROJECT` | Claude Code — Vertex | https://docs.claude.com/en/docs/claude-code/google-vertex-ai |
| `CLOUD_ML_REGION` | Claude Code — Vertex | https://docs.claude.com/en/docs/claude-code/google-vertex-ai |

<!-- passlist:end -->
