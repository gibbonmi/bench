#!/usr/bin/env bash
# Canary fixture: the generated wrapper has NO missing-target guard (story 8) — it
# execs the target directly, so a moved target fails without the remedy message. The
# shim stale-target contract must go red: no "bench moved" on stderr.
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
      tgt="$(cd "$(dirname "$0")" && pwd)/bench.sh"
      # No `if [ ! -x "$target" ]` guard — the regression.
      printf '#!/usr/bin/env bash\n# bench-shim v1\ntarget=%s\nexec "$target" "$@"\n' "$tgt" > "$dir/bench"
      chmod +x "$dir/bench"; echo "  wrote shim $dir/bench"; exit 0
    fi
    ;;
esac
