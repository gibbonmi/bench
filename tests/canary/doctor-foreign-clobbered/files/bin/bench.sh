#!/usr/bin/env bash
# Canary fixture: `doctor --fix` CLOBBERS whatever sits at the target instead of
# refusing a foreign (marker-less) file (story 5). The foreign-refuse contract must go
# red — it expects exit 1 and a byte-identical file; this stub exits 0 and overwrites.
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
      printf '#!/usr/bin/env bash\n# bench-shim v1\ntarget=/x\nexec "$target" "$@"\n' > "$dir/bench"
      chmod +x "$dir/bench"; echo "  wrote shim $dir/bench"; exit 0
    fi
    ;;
esac
