#!/usr/bin/env bash
# Canary support only: a minimal doctor that writes a placeholder shim to the first
# writable non-manager PATH dir (or ~/.local/bin, failing loudly) so the postinstall
# contract's global/.git/non-npm stories can pass and reach the write-failure story.
# The regression under test is in bench-postinstall.sh, not here; no assertion reads
# the shim body. Intentionally identical to the twin in postinstall-guard-bypassed —
# incidental plumbing, not shared knowledge (both are throwaway, neither re-derives
# the real doctor's selection or shim-template logic).
set -uo pipefail
case "${1:-}" in
  doctor)
    [ "${2:-}" = "--fix" ] || exit 1
    d=""; IFS=:
    for cand in $PATH; do
      [ -d "$cand" ] && [ -w "$cand" ] || continue
      case "$cand" in *"/nvm/"*|/usr|/usr/*|/bin|/sbin|/opt|/opt/*) continue ;; esac
      d="$cand"; break
    done
    [ -n "$d" ] || { mkdir -p "$HOME/.local/bin" || exit 1; d="$HOME/.local/bin"; }
    printf '#!/usr/bin/env bash\nexec true\n' > "$d/bench" || exit 1
    chmod +x "$d/bench"; echo "  wrote shim $d/bench"
    ;;
esac
