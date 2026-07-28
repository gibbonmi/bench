#!/usr/bin/env bash
# Canary fixture: the generated wrapper forwards args with unquoted $* instead of
# "$@" (story 9), so multi-word and glob args are resplit and expanded. The shim
# arg-passthrough contract must go red.
set -uo pipefail
pick_dir() {
  local d IFS=:
  for d in $PATH; do
    [ -d "$d" ] && [ -w "$d" ] || continue
    case "$d" in *"/nvm/"*|*/.nvm/*|/usr|/usr/*|/bin|/sbin|/opt|/opt/*) continue ;; esac
    printf '%s\n' "$d"; return
  done
  printf '%s\n' "$HOME/.local/bin"
}
case "${1:-}" in
  doctor)
    if [ "${2:-}" = "--fix" ]; then
      dir="$(pick_dir)"; mkdir -p "$dir"
      # exec "$target" $*  — unquoted, the regression.
      printf '#!/usr/bin/env bash\n# bench-shim v1\ntarget=/x\nexec "$target" $*\n' > "$dir/bench"
      chmod +x "$dir/bench"; echo "  wrote shim $dir/bench"; exit 0
    fi
    ;;
esac
