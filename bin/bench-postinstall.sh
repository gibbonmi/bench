#!/usr/bin/env bash
# shellcheck shell=bash
# npm postinstall — install the PATH shim on a global install, but NEVER fail the
# install. The shim is a convenience; a wrong read must be non-destructive either
# way, so every path exits 0. Thin pass-through: the real work is `bench doctor
# --fix` in the compiled core, invoked package-relative because `bench` may not be
# on PATH yet.
#
# Guards (both required to mutate):
#   - npm_config_global truthy — this is `npm install -g`, set to "true" across npm
#     v6–v10+ and absent under pnpm/yarn/bun.
#   - the package root has no .git — a git checkout or `npm link` is a dev tree; the
#     .git guard covers both.
# Any other condition (env absent, .git present, probe failure, write failure) falls
# through to one advice line and exits 0.
set -u

advice() { printf "bench: run 'bench doctor --fix' to install the PATH shim so login shells resolve bench\n"; }

case "${npm_config_global:-}" in
  true|1|yes) ;;
  *) advice; exit 0 ;;
esac

here="$(cd "$(dirname "$0")/.." && pwd)" || { advice; exit 0; }
if [ -e "$here/.git" ]; then advice; exit 0; fi

# Relay the fix's announcements on success; fall through to advice on any failure.
if bash "$here/bin/bench.sh" doctor --fix; then :; else advice; fi
exit 0
