#!/usr/bin/env bash
# SessionStart hook: print the ambient dashboard when a session opens cold.
# A thin wrapper over `bench status` — the single renderer the user also runs on demand
# (one source of truth). Never blocks the session: outside a repo, or on any error, it
# prints nothing and exits 0. Shared across harnesses; the .claude adapter wires it under
# hooks.SessionStart, and any AGENTS.md harness can point its own start hook here.
set -uo pipefail
bench status 2>/dev/null || true
