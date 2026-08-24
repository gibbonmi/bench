#!/usr/bin/env bash
# shellcheck shell=bash
# npm postinstall script. It installs the PATH shim on a global install. It never
# fails the install: the shim is a convenience, and a wrong read must not damage the
# tree. Every path exits 0.
#
# The script only calls through: the real work happens in `bench doctor --fix` in the
# compiled core. That fix also publishes the promotion-broker manifest beside the
# wrapper, so a fresh global install can land through `bench worktree land`. The script invokes bench.sh by its package-relative path, because
# `bench` may not sit on PATH yet.
#
# Both guards must pass before the script writes anything.
#   - npm_config_global is true — this marks `npm install -g`. npm v6 through v10+ set
#     it to "true". pnpm, yarn, and bun do not set it.
#   - the package root has no .git directory — a git checkout or `npm link` produces a
#     dev tree. This guard excludes both.
#
# Any other condition — a missing env var, a .git directory, a failed probe, a failed
# write — falls through to one advice line, and the script exits 0.
set -u

advice() { printf "bench: run 'bench doctor --fix' to install the PATH shim so login shells resolve bench\n"; }

case "${npm_config_global:-}" in
  true|1|yes) ;;
  *) advice; exit 0 ;;
esac

here="$(cd "$(dirname "$0")/.." && pwd)" || { advice; exit 0; }
if [ -e "$here/.git" ]; then advice; exit 0; fi

# On success, relay the fix's announcements. On any failure, fall through to advice.
if bash "$here/bin/bench.sh" doctor --fix; then :; else advice; fi
exit 0
