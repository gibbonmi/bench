#!/usr/bin/env bash
# Canary fixture support: a correct mini-doctor so the postinstall contract's global /
# .git / non-npm stories pass; the planted regression is in bench-postinstall.sh.
set -uo pipefail
resolve_script_path() {
  local s="${BASH_SOURCE[0]:-$0}" d t
  while [ -L "$s" ]; do d="$(cd -P "$(dirname "$s")" && pwd)"; t="$(readlink "$s")"; case "$t" in /*) s="$t" ;; *) s="$d/$t" ;; esac; done
  d="$(cd -P "$(dirname "$s")" && pwd)"; printf '%s/%s\n' "$d" "$(basename "$s")"
}
doctor_fix() {
  local d dir="" IFS=:
  for d in $PATH; do
    [ -d "$d" ] && [ -w "$d" ] || continue
    case "$d" in *"/nvm/"*|*/.nvm/*|/usr|/usr/*|/bin|/sbin|/opt|/opt/*) continue ;; esac
    dir="$d"; break
  done
  [ -n "$dir" ] || { mkdir -p "$HOME/.local/bin" || return 1; dir="$HOME/.local/bin"; }
  printf '#!/usr/bin/env bash\n# bench-shim v1\ntarget=%s\nexec "$target" "$@"\n' "$(resolve_script_path)" > "$dir/bench" || return 1
  chmod +x "$dir/bench"; echo "  wrote shim $dir/bench"
}
case "${1:-}" in
  doctor) [ "${2:-}" = "--fix" ] && { doctor_fix; exit $?; }; exit 1 ;;
esac
