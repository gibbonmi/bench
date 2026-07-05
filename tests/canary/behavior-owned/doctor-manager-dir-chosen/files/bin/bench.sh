#!/usr/bin/env bash
# Canary fixture: `doctor --fix` writes into the FIRST writable PATH dir with no
# manager-owned exclusion (story 6) — so it lands in the planted nvm dir, not the plain
# one. The --fix write contract must go red: the shim never appears in the plain dir.
set -uo pipefail
pick_dir() {
  local d IFS=:
  for d in $PATH; do
    [ -d "$d" ] && [ -w "$d" ] || continue
    printf '%s\n' "$d"; return   # no manager-owned exclusion — the regression
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
