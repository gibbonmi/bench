# Runtime contracts for the destructive-git guard hook.

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"
  guard="$root/.bench/hooks/block-dangerous-git.sh"
  run_guard() {
    local command="$1" out rc
    out="$(printf '{"tool_input":{"command":"%s"}}\n' "$command" | bash "$guard" 2>&1)" && rc=0 || rc=$?
    [ "$rc" = "2" ] || { printf '%s\n' "$out"; echo "guard allowed $command (exit $rc)"; exit 1; }
    grep -qF 'BLOCKED:' <<<"$out" || { echo "guard block output was not actionable"; exit 1; }
  }
  run_guard 'git push'
  run_guard 'git -C . push'
  run_guard 'git -C /tmp reset --hard'
  run_guard 'git -C . clean -fd'
  run_guard 'git branch -D old-work'
  run_guard 'git rebase main'
  run_guard 'git checkout -- README.md'
  run_guard 'git checkout --pathspec-from-file=paths.txt'
  run_guard 'git checkout --pathspec-from-file paths.txt'
  run_guard 'git restore README.md'
  run_guard 'git restore --pathspec-from-file=paths.txt'
  run_guard 'git restore --pathspec-from-file paths.txt'
  allowed="$(printf '{"tool_input":{"command":"git -C . status --short"}}\n' | bash "$guard" 2>&1)" || { printf '%s\n' "$allowed"; echo "guard blocked harmless git status"; exit 1; }
) || err "block-dangerous-git global-option contract failed"
rm -rf "$tmp"
