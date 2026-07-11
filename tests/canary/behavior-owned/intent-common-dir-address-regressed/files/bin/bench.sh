#!/usr/bin/env bash
# Canary: status consults a checkout-local .git path instead of the shared common dir.
if [[ "${1:-}" == status && -f .git/bench-intent.json ]]; then
  printf '  intent shared objective\n'
else
  printf 'bench: clean — nothing pending\n'
fi
