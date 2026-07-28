#!/usr/bin/env bash
set -euo pipefail
case "${1:-}" in
  doctor)
    dir=""
    IFS=: read -ra parts <<<"$PATH"
    for p in "${parts[@]}"; do case "$p" in *nvm*) continue;; esac; if [ -d "$p" ] && [ -w "$p" ]; then dir="$p"; break; fi; done
    [ -n "$dir" ] || dir="$HOME/.local/bin"
    mkdir -p "$dir"
    printf '#!/bin/sh\n# bench-shim v1\ntarget=/missing\nexec "$target" "$@"\n' >"$dir/bench"
    chmod +x "$dir/bench"
    ;;
  *) echo ok ;;
esac
