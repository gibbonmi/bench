# Runtime CLI contracts for `bench doctor`, the generated shim, postinstall, and the
# session-start advisory. The seam is the CLI-subprocess sandbox (the runtime-contracts
# pattern): a fabricated HOME/PATH/env, invoke the real scripts as subprocesses, assert
# exit code + stdout/stderr + filesystem effects. Never asserts internal functions.
#
# Without a CLI in the tree there is nothing to contract-test; skip silently (the
# runtime-contracts fragment already surfaces the missing-CLI signal), so a doctor
# canary's targeted substring never leaks into an empty-fixture run.
[ -f "$root/bin/bench.sh" ] || return 0 2>/dev/null || exit 0

# One source for the sandbox provisioning every doctor body needs: HOME under the
# fixture, an nvm-owned bin (manager territory, writable-but-wrong) and a plain
# writable bin, plus a PATH where the plain dir qualifies and the nvm dir is excluded.
# Inherited by each contract body's subshell; callers override PATH for the
# fallback/manager-only cases.
doctor_sandbox() {   # sets HOME/SHELL/NVM_DIR and: DSB_NVMBIN DSB_PLAIN DSB_PATH
  export HOME="$tmp/home"; mkdir -p "$HOME"
  export SHELL=/bin/bash
  export NVM_DIR="$tmp/nvm"; DSB_NVMBIN="$NVM_DIR/versions/node/v22/bin"; mkdir -p "$DSB_NVMBIN"
  DSB_PLAIN="$tmp/plain"; mkdir -p "$DSB_PLAIN"
  DSB_PATH="$DSB_NVMBIN:$DSB_PLAIN:/usr/bin:/bin"
}

doctor_copy_kit() {
  local kit="$1"
  mkdir -p "$kit/bin" "$kit/dist" "$kit/.agents/commands"
  cp "$root"/bin/*.sh "$kit/bin/"
  cp "$root/dist/bench" "$kit/dist/bench"
  cp "$root/AGENTS.md" "$kit/AGENTS.md"
}

# --- doctor report: health classification + machine-exact removal pair (stories 2, 15)
contract "bench doctor report contract" <<'BODY'
doctor_sandbox
target_path="$DSB_PLAIN/bench"
# missing → exit 1, nothing written, removal pair names the resolved shim path
out="$(PATH="$DSB_PATH" bash "$root/bin/bench.sh" doctor)"; rc=$?
[ "$rc" = "1" ] || { echo "report on a missing shim did not exit 1 (got $rc)"; exit 1; }
[ -e "$target_path" ] && { echo "report on a missing shim wrote a file"; exit 1; }
grep -qF "missing" <<<"$out" || { echo "report did not classify the missing shim"; exit 1; }
grep -qF "rm -f \"$target_path\"" <<<"$out" || { echo "removal pair lacks the resolved shim path"; exit 1; }
# healthy → exit 0
PATH="$DSB_PATH" bash "$root/bin/bench.sh" doctor --fix >/dev/null
PATH="$DSB_PATH" bash "$root/bin/bench.sh" doctor >/dev/null; rc=$?
[ "$rc" = "0" ] || { echo "report on a healthy shim did not exit 0 (got $rc)"; exit 1; }
# foreign → exit 1
rm -f "$target_path"; printf 'not a bench shim\n' > "$target_path"
PATH="$DSB_PATH" bash "$root/bin/bench.sh" doctor >/dev/null; rc=$?
[ "$rc" = "1" ] || { echo "report on a foreign file did not exit 1 (got $rc)"; exit 1; }
# stale → exit 1: a marked shim whose target is gone
printf '#!/usr/bin/env bash\n# bench-shim v1 marker\ntarget=/no/such/bench\nexec "$target" "$@"\n' > "$target_path"
out="$(PATH="$DSB_PATH" bash "$root/bin/bench.sh" doctor)"; rc=$?
[ "$rc" = "1" ] || { echo "report on a stale shim did not exit 1 (got $rc)"; exit 1; }
grep -qF "stale" <<<"$out" || { echo "report did not classify the stale shim"; exit 1; }
# a read-only report must never execute a marker-bearing file's contents
printf '#!/usr/bin/env bash\n# bench-shim v1 marker\ntarget=$(touch "%s/pwned")\nexec "$target"\n' "$tmp" > "$target_path"
PATH="$DSB_PATH" bash "$root/bin/bench.sh" doctor >/dev/null 2>&1 || true
[ ! -e "$tmp/pwned" ] || { echo "report executed a hostile shim's target line"; exit 1; }
BODY

# --- doctor manifest skew: stamp mismatch warns; pre-stamp manifest warns without changing exit
contract "bench doctor manifest skew contract" <<'BODY'
doctor_sandbox
mkdir -p .bench
printf '.bench/BENCH.md\tabc\n' > .bench/link-manifest.tsv
out="$(PATH="$DSB_PATH" bash "$root/bin/bench.sh" doctor)"; rc=$?
[ "$rc" = "1" ] || { echo "pre-stamp manifest changed doctor exit (got $rc)"; exit 1; }
grep -qF "version unknown" <<<"$out" || { echo "pre-stamp manifest did not warn skew-unknown"; exit 1; }
printf '#kit\t0.0.0\n.bench/BENCH.md\tabc\n' > .bench/link-manifest.tsv
out="$(PATH="$DSB_PATH" bash "$root/bin/bench.sh" doctor)"; rc=$?
[ "$rc" = "1" ] || { echo "skewed manifest changed doctor exit (got $rc)"; exit 1; }
grep -qF "version skew" <<<"$out" || { echo "skewed manifest did not report version skew"; exit 1; }
grep -qF "0.0.0" <<<"$out" || { echo "skewed manifest output did not name the stamped version"; exit 1; }
BODY

# --- doctor --fix: picks the first non-manager writable PATH dir, marker+target, announced (stories 3, 6)
contract "bench doctor --fix write contract" <<'BODY'
doctor_sandbox
out="$(PATH="$DSB_PATH" bash "$root/bin/bench.sh" doctor --fix)"; rc=$?
[ "$rc" = "0" ] || { echo "--fix did not exit 0 (got $rc)"; exit 1; }
[ -f "$DSB_PLAIN/bench" ] || { echo "--fix wrote no shim in the plain PATH dir"; exit 1; }
[ -e "$DSB_NVMBIN/bench" ] && { echo "--fix wrote into the manager-owned nvm dir"; exit 1; }
grep -qF "bench-shim v1" "$DSB_PLAIN/bench" || { echo "shim carries no bench marker"; exit 1; }
grep -qF "$root/bin/bench.sh" "$DSB_PLAIN/bench" || { echo "shim does not target the resolved CLI"; exit 1; }
grep -qF "$DSB_PLAIN/bench" <<<"$out" || { echo "--fix did not announce the written path"; exit 1; }
BODY

# --- doctor --fix under a spaced/glob target path: %q quoting re-execs intact (stories 3, 9 edge)
# The spec's edge inventory claims the spaced-path case is covered; this is that coverage.
# The invoked CLI lives under a parent containing a space, so the resolved target the shim
# embeds is spaced — exercising doctor_shim_content's %q quoting end to end.
contract "bench doctor --fix spaced-target contract" --space-path <<'BODY'
doctor_sandbox
kit="$tmp/kit"; doctor_copy_kit "$kit"
out="$(PATH="$DSB_PATH" bash "$kit/bin/bench.sh" doctor --fix)"; rc=$?
[ "$rc" = "0" ] || { echo "spaced-path --fix did not exit 0 (got $rc)"; exit 1; }
[ -f "$DSB_PLAIN/bench" ] || { echo "spaced-path --fix wrote no shim"; exit 1; }
grep -qF "$kit/bin/bench.sh" "$DSB_PLAIN/bench" || { echo "shim lost the spaced target path"; exit 1; }
got="$("$DSB_PLAIN/bench" doctor 2>/dev/null | head -1)"
grep -qF "shim health" <<<"$got" || { echo "spaced-target shim failed to re-exec the CLI"; exit 1; }
BODY

# --- doctor --fix idempotency: second run is an announced no-op, bytes unchanged (story 4)
contract "bench doctor --fix idempotency contract" <<'BODY'
doctor_sandbox
PATH="$DSB_PATH" bash "$root/bin/bench.sh" doctor --fix >/dev/null
before="$(cat "$DSB_PLAIN/bench")"
out="$(PATH="$DSB_PATH" bash "$root/bin/bench.sh" doctor --fix)"; rc=$?
[ "$rc" = "0" ] || { echo "second --fix did not exit 0 (got $rc)"; exit 1; }
[ "$(cat "$DSB_PLAIN/bench")" = "$before" ] || { echo "second --fix rewrote an already-current shim"; exit 1; }
grep -qiF "no change" <<<"$out" || { echo "second --fix did not announce a no-op"; exit 1; }
BODY

# --- doctor --fix collision: refuse a foreign (marker-less) file, byte-identical (story 5)
contract "bench doctor --fix foreign-refuse contract" <<'BODY'
doctor_sandbox
foreign="$DSB_PLAIN/bench"
printf 'foreign contents\n' > "$foreign"
sum_before="$(cksum < "$foreign")"
out="$(PATH="$DSB_PATH" bash "$root/bin/bench.sh" doctor --fix 2>&1)"; rc=$?
[ "$rc" = "1" ] || { echo "--fix over a foreign file did not exit 1 (got $rc)"; exit 1; }
[ "$(cksum < "$foreign")" = "$sum_before" ] || { echo "--fix clobbered a foreign file"; exit 1; }
grep -qiF "refus" <<<"$out" || { echo "--fix did not report the refusal"; exit 1; }
# present-but-empty is marker-less → also refused
: > "$foreign"
PATH="$DSB_PATH" bash "$root/bin/bench.sh" doctor --fix >/dev/null 2>&1; rc=$?
[ "$rc" = "1" ] || { echo "--fix over a present-but-empty file did not refuse (got $rc)"; exit 1; }
[ -s "$foreign" ] && { echo "--fix wrote over the present-but-empty file"; exit 1; }
# a Bench-marked shim with NO trailing newline is recognized, not refused
printf '#!/usr/bin/env bash\n# bench-shim v1 marker\ntarget=%s\nexec "$target" "$@"' "$root/bin/bench.sh" > "$foreign"
PATH="$DSB_PATH" bash "$root/bin/bench.sh" doctor --fix >/dev/null 2>&1; rc=$?
[ "$rc" = "0" ] || { echo "--fix refused a newline-less Bench-marked shim (got $rc)"; exit 1; }
BODY

# --- doctor --fix fallback: only manager dirs on PATH → create ~/.local/bin, announced (story 6)
contract "bench doctor --fix fallback contract" <<'BODY'
doctor_sandbox
out="$(PATH="$DSB_NVMBIN:/usr/bin:/bin" bash "$root/bin/bench.sh" doctor --fix)"; rc=$?
[ "$rc" = "0" ] || { echo "fallback --fix did not exit 0 (got $rc)"; exit 1; }
[ -f "$HOME/.local/bin/bench" ] || { echo "fallback did not write to ~/.local/bin"; exit 1; }
grep -qF "created directory $HOME/.local/bin" <<<"$out" || { echo "fallback did not announce the dir creation"; exit 1; }
BODY

# --- doctor --fix off-PATH notice: print the exact PATH line, edit no rc file (story 7)
contract "bench doctor --fix path-notice contract" <<'BODY'
doctor_sandbox
export SHELL=/bin/bash
out="$(PATH="$DSB_NVMBIN:/usr/bin:/bin" bash "$root/bin/bench.sh" doctor --fix)"
grep -qF "export PATH" <<<"$out" || { echo "off-PATH --fix printed no export line"; exit 1; }
grep -qF ".local/bin" <<<"$out" || { echo "off-PATH --fix line names no ~/.local/bin"; exit 1; }
[ -e "$HOME/.bashrc" ] && { echo "--fix edited an rc file"; exit 1; }
exit 0
BODY

# --- generated shim: missing target fails loud (exit 127 + remedy on stderr) (story 8)
contract "bench doctor shim stale-target contract" <<'BODY'
doctor_sandbox
PATH="$DSB_PATH" bash "$root/bin/bench.sh" doctor --fix >/dev/null
shim="$tmp/shim"; cp "$DSB_PLAIN/bench" "$shim"
sed 's#^target=.*#target=/no/such/bench#' "$DSB_PLAIN/bench" > "$shim"; chmod +x "$shim"
err_out="$(PATH="/usr/bin:/bin" "$shim" help 2>&1 >/dev/null)"; rc=$?
[ "$rc" = "127" ] || { echo "stale-target shim did not exit 127 (got $rc)"; exit 1; }
grep -qF "bench moved" <<<"$err_out" || { echo "stale-target shim printed no remedy"; exit 1; }
BODY

# --- generated shim: multi-word and glob args pass through intact (story 9)
contract "bench doctor shim arg-passthrough contract" <<'BODY'
doctor_sandbox
PATH="$DSB_PATH" bash "$root/bin/bench.sh" doctor --fix >/dev/null
stub="$tmp/stub"; printf '#!/usr/bin/env bash\nfor a in "$@"; do printf "[%%s]" "$a"; done\necho\n' > "$stub"; chmod +x "$stub"
shim="$tmp/shim"; sed "s#^target=.*#target=$stub#" "$DSB_PLAIN/bench" > "$shim"; chmod +x "$shim"
got="$(PATH="/usr/bin:/bin" "$shim" 'a b' '*' c)"
[ "$got" = "[a b][*][c]" ] || { echo "shim mangled the args (got '$got')"; exit 1; }
BODY

# --- postinstall: guard conjunction + never-fail-the-install invariant (stories 10, 11, 12, 13)
contract "bench doctor postinstall contract" <<'BODY'
doctor_sandbox
pkg="$tmp/pkg"; doctor_copy_kit "$pkg"
pin="$pkg/bin/bench-postinstall.sh"
# global + no .git → writes shim (story 10)
out="$(npm_config_global=true PATH="$DSB_PATH" bash "$pin")"; rc=$?
[ "$rc" = "0" ] || { echo "postinstall exited nonzero on the write path (got $rc)"; exit 1; }
[ -f "$DSB_PLAIN/bench" ] || { echo "postinstall did not write the shim on a global install"; exit 1; }
rm -f "$DSB_PLAIN/bench"
# .git present → skip mutation, advice, exit 0 (story 11)
touch "$pkg/.git"
out="$(npm_config_global=true PATH="$DSB_PATH" bash "$pin")"; rc=$?
[ "$rc" = "0" ] || { echo "postinstall exited nonzero under the .git guard (got $rc)"; exit 1; }
[ -e "$DSB_PLAIN/bench" ] && { echo "postinstall mutated with .git present"; exit 1; }
grep -qF "doctor --fix" <<<"$out" || { echo "postinstall printed no advice under the .git guard"; exit 1; }
rm -f "$pkg/.git"
# no npm_config_global → advice, exit 0, no shim (story 12)
out="$(PATH="$DSB_PATH" bash "$pin")"; rc=$?
[ "$rc" = "0" ] || { echo "postinstall exited nonzero without npm_config_global (got $rc)"; exit 1; }
[ -e "$DSB_PLAIN/bench" ] && { echo "postinstall mutated without npm_config_global"; exit 1; }
grep -qF "doctor --fix" <<<"$out" || { echo "postinstall printed no advice without npm_config_global"; exit 1; }
# write failure with guards satisfied → still exit 0 (story 13): no PATH dir qualifies and HOME is read-only
roh="$tmp/rohome"; mkdir -p "$roh"; chmod 0555 "$roh"
rc=0; HOME="$roh" npm_config_global=true PATH="$DSB_NVMBIN:/usr/bin:/bin" bash "$pin" >/dev/null 2>&1 || rc=$?
chmod 0755 "$roh"
[ "$rc" = "0" ] || { echo "postinstall exited nonzero on a write failure (got $rc)"; exit 1; }
BODY

# --- session-start advisory: append the doctor pointer on a by-path resolve (story 14)
contract "bench doctor session-start advisory contract" <<'BODY'
repo="$tmp/repo"; mkdir -p "$repo/bin"; cp "$root"/bin/*.sh "$repo/bin/"
( cd "$repo" && git init -q )
out="$(cd "$repo" && PATH="/usr/bin:/bin" bash "$root/.bench/hooks/session-start.sh" 2>/dev/null | head -1)"
grep -qF "doctor --fix" <<<"$out" || { echo "by-path session-start advisory omits the doctor pointer"; exit 1; }
BODY
