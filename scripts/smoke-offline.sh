#!/usr/bin/env bash
# Exercise direct, local-tarball npm, and loopback-registry journeys from one
# extracted archive with empty state and public egress unavailable.
set -euo pipefail

artifacts="${1:?usage: smoke-offline.sh <artifact-dir> <release-evidence-dir>}"
evidence_dir="${2:?usage: smoke-offline.sh <artifact-dir> <release-evidence-dir>}"
source_path="${BASH_SOURCE[0]}"
while [[ -L "$source_path" ]]; do
  source_dir="$(cd "$(dirname "$source_path")" && pwd)"
  link_target="$(readlink "$source_path")"
  [[ "$link_target" == /* ]] && source_path="$link_target" || source_path="$source_dir/$link_target"
done
root="$(cd "$(dirname "$source_path")/.." && pwd)"
# shellcheck source=scripts/lib/search.sh
. "$root/scripts/lib/search.sh"
version="$(node -p 'require(process.argv[1]).version' "$root/package.json")"
host_os="$(case "$(uname -s)" in Darwin) printf darwin ;; Linux) printf linux ;; *) printf unsupported ;; esac)"
host_arch="$(case "$(uname -m)" in arm64|aarch64) printf arm64 ;; x86_64|amd64) printf x64 ;; *) printf unsupported ;; esac)"
target_row="$(node "$root/scripts/release-plan.mjs" "$root" target "$host_os" "$host_arch")" || true
target="$(cut -f1-2 <<< "$target_row" | tr '\t' '-')"
[[ -n "$target" ]] || { printf 'offline smoke: unsupported host %s/%s\n' "$host_os" "$host_arch" >&2; exit 2; }
runtime_target="$(cut -f3,4 <<< "$target_row" | tr '\t' '/')"
[[ -n "$runtime_target" ]] || { printf 'offline smoke: target is absent from canonical matrix: %s\n' "$target" >&2; exit 1; }

wrapper="$artifacts/redbench-${version}.tgz"
native="$artifacts/redbench-${target}-${version}.tgz"
archive="$artifacts/redbench-${version}-${target}.tar.gz"
[[ -f "$evidence_dir/release-index.json" && ! -L "$evidence_dir/release-index.json" && -s "$evidence_dir/release-index.json" && -f "$evidence_dir/SHA256SUMS" && ! -L "$evidence_dir/SHA256SUMS" && -s "$evidence_dir/SHA256SUMS" ]] || { printf 'offline smoke: approved release evidence is missing or unsafe\n' >&2; exit 1; }
for input in "$wrapper" "$native" "$archive"; do
  [[ -f "$input" && ! -L "$input" && -s "$input" ]] || { printf 'offline smoke: required artifact is missing or unsafe: %s\n' "$input" >&2; exit 1; }
done

tmp="$(mktemp -d)"
mkdir -p "$tmp/process-tmp"
export TMPDIR="$tmp/process-tmp"
registry_pid=""
egress_log="$tmp/undeclared-egress"
repair_marker="$tmp/repair-attempt"
codex_marker="$tmp/codex-attempt"
git_marker="$tmp/git-refresh-attempt"
offline_payload="$tmp/ft87-offline-payload.json"
offline_flag="$(printf '%s=%s' BENCH_OFFLINE 1)"
evidence_armed=0
: > "$egress_log"
cleanup() {
  if [[ -n "$registry_pid" ]]; then
    kill "$registry_pid" 2>/dev/null || true
    wait "$registry_pid" 2>/dev/null || true
  fi
  rm -rf "$tmp"
}
write_offline_evidence() {
  local smoke_exit="$1" requested_status=failed attempts_exit=0
  [[ "$smoke_exit" == 0 ]] && requested_status=auto
  node -e '
const fs=require("fs"),[file,flag,repair,codex,git,egress]=process.argv.slice(1);
const present=marker=>marker&&fs.existsSync(marker)?1:0;
const lines=egress&&fs.existsSync(egress)?fs.readFileSync(egress,"utf8").split(/\r?\n/).filter(Boolean):[];
const count=host=>lines.filter(line=>line===host||line.startsWith(host+":")).length;
const payload={flag,journeys:["direct","local-npm","loopback-registry","offline-policy"],operations:[
 {class:"wrapper_binary_repair",observed_attempts:present(repair)},
 {class:"worktree_git_refresh",observed_attempts:present(git)},
 {class:"codex_discovery_subprocess_and_bundled_fallback",observed_attempts:present(codex)},
 {class:"openai_models_request",observed_attempts:count("api.openai.com")},
 {class:"anthropic_models_request",observed_attempts:count("api.anthropic.com")},
]}; fs.writeFileSync(file,JSON.stringify(payload,null,2)+"\n");
' "$offline_payload" "$offline_flag" "$repair_marker" "$codex_marker" "$git_marker" "$egress_log" || return 1
  node "$root/scripts/write-producer-envelope.mjs" "$root" "$evidence_dir" public.ft87.offline_network_control "$(git -C "$root" rev-parse HEAD)" "$version" "$requested_status" "$offline_payload" || return 1
  node -e 'const payload=require(process.argv[1]);process.exit(payload.operations.some(operation=>operation.observed_attempts!==0)?1:0)' "$offline_payload" || attempts_exit=$?
  if [[ "$attempts_exit" != 0 ]]; then
    printf 'offline smoke: sentinel observed a nonzero attempt\n' >&2
    return "$attempts_exit"
  fi
}
finalize() {
  local smoke_exit=$? evidence_exit=0
  trap - EXIT
  if [[ "$evidence_armed" == 1 ]]; then
    write_offline_evidence "$smoke_exit" || evidence_exit=$?
    [[ "$smoke_exit" != 0 ]] || smoke_exit="$evidence_exit"
  fi
  cleanup
  exit "$smoke_exit"
}
interrupt() {
  exit 130
}
trap finalize EXIT
trap interrupt INT TERM HUP

offline_stage() {
  local stage="$1"
  if [[ "${BENCH_OFFLINE_TEST_INTERRUPT_STAGE:-}" == "$stage" ]]; then
    [[ -n "${BENCH_OFFLINE_TEST_STAGE_READY_DIR:-}" ]] || { printf 'offline smoke: interruption stage requires a ready directory\n' >&2; return 1; }
    : > "$BENCH_OFFLINE_TEST_STAGE_READY_DIR/$stage"
    while [[ -e "$BENCH_OFFLINE_TEST_STAGE_READY_DIR/$stage" ]]; do sleep 0.02; done
  fi
}

offline_process() {
  local stage="$1" process
  shift
  offline_stage "$stage"
  process="$(basename "$1")"
  case "$process" in
    go|cargo|cc|gcc|clang|make)
      printf 'offline smoke: forbidden repair or rebuild process attempted: %s\n' "$process" >&2
      return 1
      ;;
  esac
  "$@"
}

native_offline_process() {
  local stage="$1"
  shift
  offline_stage "$stage"
  case "$host_os" in
    linux)
      command -v bwrap >/dev/null 2>&1 || { printf 'offline smoke: native network sentinel is unavailable: bwrap\n' >&2; return 1; }
      local -a writable=(--bind "$tmp" "$tmp")
      if [[ -n "${BENCH_OFFLINE_TEST_NATIVE_MARKER:-}" ]]; then
        writable+=(--bind "$(dirname "$BENCH_OFFLINE_TEST_NATIVE_MARKER")" "$(dirname "$BENCH_OFFLINE_TEST_NATIVE_MARKER")")
      fi
      bwrap --unshare-net --ro-bind / / "${writable[@]}" --dev /dev --proc /proc -- "$@"
      ;;
    darwin)
      command -v sandbox-exec >/dev/null 2>&1 || { printf 'offline smoke: native network sentinel is unavailable: sandbox-exec\n' >&2; return 1; }
      sandbox-exec -p '(version 1) (allow default) (deny network*)' "$@"
      ;;
  esac
}

bundle="$tmp/bundle"
mkdir -p "$bundle"
tar_members="$tmp/archive-members"
tar -tzf "$archive" > "$tar_members" || { printf 'offline smoke: archive tarball is corrupt\n' >&2; exit 1; }
while IFS= read -r member; do
  [[ "$member" == "redbench-${version}-${target}/" || "$member" == "redbench-${version}-${target}/"* ]] || {
    printf 'offline smoke: archive contains an unexpected root: %s\n' "$member" >&2
    exit 1
  }
  [[ "$member" != *"../"* && "$member" != *"/../"* && "$member" != /* && "$member" != *"\\"* ]] || {
    printf 'offline smoke: archive contains an unsafe path: %s\n' "$member" >&2
    exit 1
  }
done < "$tar_members"
tar -xzf "$archive" -C "$bundle"
bundle_root="$bundle/redbench-${version}-${target}"
for required in \
  "$bundle_root/bin/bench" \
  "$bundle_root/packages/redbench-${version}.tgz" \
  "$bundle_root/packages/redbench-${target}-${version}.tgz" \
  "$bundle_root/OFFLINE.md" \
	"$bundle_root/evidence/component-manifest.json" \
  "$bundle_root/evidence/components/wrapper-component-manifest.json" \
  "$bundle_root/evidence/components/platform-component-manifest.json"; do
  [[ -f "$required" && ! -L "$required" && -s "$required" ]] || { printf 'offline smoke: archive inventory is incomplete: %s\n' "$required" >&2; exit 1; }
done
verify_component_manifest() {
  # shellcheck disable=SC2016 # Node template literals are intentionally literal here.
  node -e '
  const fs=require("fs"),crypto=require("crypto"),path=require("path");
  const [root,manifestPath]=process.argv.slice(1), manifest=JSON.parse(fs.readFileSync(manifestPath));
  if(manifest.schema_version!==1||!Array.isArray(manifest.files))throw new Error("archive component manifest is malformed");
  const seen=new Set(); for(const item of manifest.files){if(!item||seen.has(item.path)||item.path==="component-manifest.json"||item.path==="evidence/component-manifest.json")throw new Error("component manifest duplicates an entry");seen.add(item.path);const file=path.join(root,item.path),data=fs.readFileSync(file),mode=(fs.statSync(file).mode&0o777).toString(8);const sha=crypto.createHash("sha256").update(data).digest("hex");if(item.size!==data.length||item.mode!==mode||item.sha256!==sha)throw new Error(`component manifest disagrees with ${item.path}`)}
  const actual=[];const walk=dir=>{for(const entry of fs.readdirSync(dir,{withFileTypes:true})){const rel=path.relative(root,path.join(dir,entry.name));if(entry.isDirectory())walk(path.join(dir,entry.name));else if(entry.isFile()&&rel!==path.relative(root,manifestPath))actual.push(rel);else if(!entry.isFile())throw new Error(`component contains unsafe member ${rel}`)}};walk(root);if(actual.length!==seen.size||actual.some(file=>!seen.has(file)))throw new Error("component manifest does not enumerate every file");
' "$1" "$2"
}
verify_component_manifest "$bundle_root" "$bundle_root/evidence/component-manifest.json" || { printf 'offline smoke: archive component manifest does not verify internal files\n' >&2; exit 1; }
cmp -s "$wrapper" "$bundle_root/packages/redbench-${version}.tgz" || { printf 'offline smoke: wrapper bytes were substituted\n' >&2; exit 1; }
cmp -s "$native" "$bundle_root/packages/redbench-${target}-${version}.tgz" || { printf 'offline smoke: platform bytes were substituted\n' >&2; exit 1; }
native_extract="$tmp/native"
wrapper_extract="$tmp/wrapper"
mkdir -p "$native_extract" "$wrapper_extract"
tar -xzf "$native" -C "$native_extract"
tar -xzf "$wrapper" -C "$wrapper_extract"
verify_component_manifest "$native_extract/package" "$native_extract/package/component-manifest.json" || { printf 'offline smoke: platform component manifest does not verify package files\n' >&2; exit 1; }
verify_component_manifest "$wrapper_extract/package" "$wrapper_extract/package/component-manifest.json" || { printf 'offline smoke: wrapper component manifest does not verify package files\n' >&2; exit 1; }
cmp -s "$native_extract/package/bin/bench" "$bundle_root/bin/bench" || { printf 'offline smoke: archive binary differs from platform package\n' >&2; exit 1; }

offline_process comparison node "$root/scripts/verify-release-artifact.mjs" "$evidence_dir/release-index.json" "$evidence_dir/SHA256SUMS" "$archive" || { printf 'offline smoke: supplied release evidence does not bind archive bytes\n' >&2; exit 1; }

if [[ -n "${BENCH_OFFLINE_TEST_PROCESS:-}" ]]; then
  offline_process repair "$BENCH_OFFLINE_TEST_PROCESS" version
fi

direct_home="$tmp/direct-home"
mkdir -p "$direct_home" "$tmp/empty-path"
direct_version="$(native_offline_process direct env -i HOME="$direct_home" BENCH_HOME="$direct_home/.bench" PATH="$tmp/empty-path" BENCH_OFFLINE=1 BENCH_NO_REPAIR=1 "$bundle_root/bin/bench" version)"
[[ "$direct_version" == "bench ${version} (${runtime_target})" ]] || { printf 'offline smoke: direct version output is wrong: %s\n' "$direct_version" >&2; exit 1; }
native_offline_process direct env -i HOME="$direct_home" BENCH_HOME="$direct_home/.bench" PATH="$tmp/empty-path" BENCH_OFFLINE=1 BENCH_NO_REPAIR=1 "$bundle_root/bin/bench" commands --brief > "$tmp/direct-commands"
bench_search_ere -q -- '^commands --brief$' "$tmp/direct-commands" || { printf 'offline smoke: direct operational probe failed\n' >&2; exit 1; }
rm -rf "$bundle"
[[ ! -e "$bundle" && ! -e "$direct_home/.bench" ]] || { printf 'offline smoke: direct removal left bundle or Bench state residue\n' >&2; exit 1; }
bundle="$tmp/install-bundle"
mkdir -p "$bundle"
tar -xzf "$archive" -C "$bundle"
bundle_root="$bundle/redbench-${version}-${target}"

local_home="$tmp/local-home"
local_prefix="$tmp/local-prefix"
repair_disabled="${BENCH_OFFLINE_REPAIR_DISABLED:-1}"
offline_mode=true
mkdir -p "$local_home" "$local_prefix"
printf '{"private":true}\n' > "$local_prefix/package.json"
[[ "$repair_disabled" == 1 && "$offline_mode" == true ]] || { printf 'offline smoke: fetch, rebuild, or repair control is disabled\n' >&2; exit 1; }
offline_stage installation
HOME="$local_home" BENCH_HOME="$local_home/.bench" BENCH_NO_REPAIR="$repair_disabled" NODE_OPTIONS="--require=$root/scripts/offline-network-sentinel.cjs" BENCH_OFFLINE_EGRESS_LOG="$egress_log" npm_config_cache="$tmp/local-cache" npm_config_offline="$offline_mode" npm_config_registry=http://127.0.0.1:9 npm_config_audit=false npm_config_fund=false npm_config_update_notifier=false npm_config_fetch_retries=0 npm_config_proxy='' npm_config_https_proxy='' NO_PROXY='*' npm install --offline --ignore-scripts --omit=optional --prefix "$local_prefix" "$bundle_root/packages/redbench-${version}.tgz" "$bundle_root/packages/redbench-${target}-${version}.tgz" >/dev/null
[[ ! -s "$egress_log" ]] || { printf 'offline smoke: local npm attempted undeclared egress: %s\n' "$(tr '\n' ' ' < "$egress_log")" >&2; exit 1; }
installed="$local_prefix/node_modules/redbench/bin/bench.sh"
[[ -x "$installed" ]] || { printf 'offline smoke: local npm install did not produce wrapper\n' >&2; exit 1; }

# Exercise every slice-1 suppression through shipped surfaces. These live
# zero-attempt sentinels feed the stable FT87 record after every public journey.
repair_bin="$local_prefix/node_modules/@redbench/${target}/bin/bench"
repair_hold="$repair_bin.offline-smoke"
probe_bin="$tmp/offline-policy-bin"
mkdir -p "$probe_bin"
printf '#!/bin/sh\n: > "$BENCH_OFFLINE_REPAIR_MARKER"\nexit 97\n' > "$probe_bin/node"
chmod +x "$probe_bin/node"
mv "$repair_bin" "$repair_hold"
evidence_armed=1
set +e
repair_output="$(HOME="$local_home" BENCH_HOME="$local_home/.bench" BENCH_OFFLINE=1 BENCH_OFFLINE_REPAIR_MARKER="$repair_marker" PATH="$probe_bin:$PATH" bash "$installed" models 2>&1)"
repair_exit=$?
set -e
mv "$repair_hold" "$repair_bin"
[[ "$repair_exit" == 127 && "$repair_output" == *"repair suppressed by BENCH_OFFLINE=1"* && ! -e "$repair_marker" ]] || { printf 'offline smoke: BENCH_OFFLINE=1 did not suppress wrapper repair before process start (exit=%s output=%s marker=%s)\n' "$repair_exit" "$repair_output" "$([[ -e "$repair_marker" ]] && printf present || printf absent)" >&2; exit 1; }

printf '#!/bin/sh\n: > "$BENCH_OFFLINE_CODEX_MARKER"\nexit 97\n' > "$probe_bin/codex"
printf '#!/bin/sh\nif [ "${1:-}" = -C ] && [ "${3:-}" = fetch ]; then : > "$BENCH_OFFLINE_GIT_MARKER"; fi\nexec "$BENCH_OFFLINE_REAL_GIT" "$@"\n' > "$probe_bin/git"
chmod +x "$probe_bin/codex" "$probe_bin/git"
models_output="$tmp/offline-models"
HOME="$local_home" BENCH_HOME="$local_home/.bench" BENCH_OFFLINE=1 OPENAI_API_KEY=sentinel ANTHROPIC_API_KEY=sentinel BENCH_OFFLINE_CODEX_MARKER="$codex_marker" PATH="$probe_bin:$PATH" native_offline_process slice1-models bash "$installed" models > "$models_output"
for source in codex openai anthropic; do
  grep -q "^  ${source},offline,offline,BENCH_OFFLINE=1$" "$models_output" || { printf 'offline smoke: %s suppression lacks BENCH_OFFLINE=1 evidence\n' "$source" >&2; exit 1; }
done
[[ ! -e "$codex_marker" ]] || { printf 'offline smoke: Codex live or bundled subprocess started under BENCH_OFFLINE=1\n' >&2; exit 1; }

policy_repo="$tmp/offline-policy-repo"
mkdir -p "$policy_repo"
"$(command -v git)" -C "$policy_repo" init -q -b main
printf 'offline\n' > "$policy_repo/tracked"
"$(command -v git)" -C "$policy_repo" add tracked
"$(command -v git)" -C "$policy_repo" -c user.name=bench -c user.email=bench@local commit -qm init
refresh_output="$tmp/offline-refresh"
HOME="$local_home" BENCH_HOME="$tmp/offline-policy-home" BENCH_OFFLINE=1 BENCH_OFFLINE_GIT_MARKER="$git_marker" BENCH_OFFLINE_REAL_GIT="$(command -v git)" PATH="$probe_bin:$PATH" native_offline_process slice1-refresh bash -c 'cd "$1" && exec "$2" worktree create --refresh --request offline-smoke --label offline-smoke' bash "$policy_repo" "$installed" > "$refresh_output"
grep -q '^  offline,BENCH_OFFLINE=1$' "$refresh_output" || { printf 'offline smoke: requested Git refresh lacks BENCH_OFFLINE=1 evidence\n' >&2; exit 1; }
[[ ! -e "$git_marker" ]] || { printf 'offline smoke: requested Git refresh started under BENCH_OFFLINE=1\n' >&2; exit 1; }
case "${BENCH_OFFLINE_TEST_INJECT_SENTINEL:-}" in
  wrapper_binary_repair) : > "$repair_marker" ;;
  worktree_git_refresh) : > "$git_marker" ;;
  codex_discovery_subprocess_and_bundled_fallback) : > "$codex_marker" ;;
  openai_models_request) printf 'api.openai.com:443\n' >> "$egress_log" ;;
  anthropic_models_request) printf 'api.anthropic.com:443\n' >> "$egress_log" ;;
  "") ;;
  *) printf 'offline smoke: unknown injected sentinel class\n' >&2; exit 2 ;;
esac
printf 'offline suppression: repair,git_refresh,codex_live,codex_bundled,openai_http,anthropic_http BENCH_OFFLINE=1 zero_attempts\n'

local_version="$(HOME="$local_home" BENCH_HOME="$local_home/.bench" BENCH_OFFLINE=1 BENCH_NO_REPAIR=1 native_offline_process local-runtime bash "$installed" version)"
[[ "$local_version" == "bench ${version} "* ]] || { printf 'offline smoke: local npm version output is wrong: %s\n' "$local_version" >&2; exit 1; }
HOME="$local_home" BENCH_HOME="$local_home/.bench" BENCH_OFFLINE=1 BENCH_NO_REPAIR=1 native_offline_process local-runtime bash "$installed" help > "$tmp/local-help"
bench_search_ere -q -- '^bench —' "$tmp/local-help" || { printf 'offline smoke: local npm operational probe failed\n' >&2; exit 1; }
offline_stage removal
HOME="$local_home" BENCH_HOME="$local_home/.bench" BENCH_NO_REPAIR=1 NODE_OPTIONS="--require=$root/scripts/offline-network-sentinel.cjs" BENCH_OFFLINE_EGRESS_LOG="$egress_log" npm_config_registry=http://127.0.0.1:9 npm_config_offline=true npm_config_audit=false npm_config_fund=false npm_config_update_notifier=false npm_config_fetch_retries=0 npm_config_proxy='' npm_config_https_proxy='' NO_PROXY='*' npm uninstall --offline --ignore-scripts --prefix "$local_prefix" redbench "@redbench/${target}" >/dev/null
[[ ! -s "$egress_log" ]] || { printf 'offline smoke: local npm uninstall attempted undeclared egress: %s\n' "$(tr '\n' ' ' < "$egress_log")" >&2; exit 1; }
[[ ! -e "$local_prefix/node_modules/redbench" && ! -e "$local_prefix/node_modules/@redbench/${target}" && ! -e "$local_home/.bench" ]] || { printf 'offline smoke: local npm uninstall left package or Bench state residue\n' >&2; exit 1; }

store="$tmp/registry-store"
mkdir -p "$store"
port_file="$tmp/registry-port"
request_file="$tmp/registry-requests"
: > "$request_file"
node "$root/scripts/offline-registry.mjs" "$store" "$port_file" "$request_file" &
registry_pid=$!
offline_stage registry-service
for _ in $(seq 1 100); do [[ -s "$port_file" || -s "$port_file.error" ]] && break; sleep 0.02; done
if [[ ! -s "$port_file" ]]; then
  printf 'offline smoke: loopback registry fixture did not start\n' >&2
  exit 1
fi
registry_home="$tmp/registry-home"
registry_prefix="$tmp/registry-prefix"
mkdir -p "$registry_home" "$registry_prefix"
printf '{"private":true}\n' > "$registry_prefix/package.json"
registry="http://127.0.0.1:$(tr -d '[:space:]' < "$port_file")"
registry_origin="127.0.0.1:${registry##*:}"
# shellcheck disable=SC2016 # Node template literals are intentionally literal here.
node -e '
  const fs=require("fs"),http=require("http");
  const [base,...files]=process.argv.slice(1);
  const put=file=>new Promise((resolve,reject)=>{const data=fs.readFileSync(file);const request=http.request(`${base}/upload/${encodeURIComponent(file.split("/").pop())}`,{method:"PUT",headers:{"content-length":data.length}},response=>{response.resume();response.statusCode===201?resolve():reject(new Error(`upload ${file} returned ${response.statusCode}`));});request.on("error",reject);request.end(data);});
  (async()=>{for(const file of files)await put(file)})().catch(error=>{console.error(error.message);process.exit(1)});
' "$registry" "$native" "$wrapper" || { printf 'offline smoke: loopback registry upload failed\n' >&2; exit 1; }
sha256() {
  node -e 'const fs=require("fs"),crypto=require("crypto");process.stdout.write(crypto.createHash("sha256").update(fs.readFileSync(process.argv[1])).digest("hex"))' "$1"
}
native_sha="$(sha256 "$native")"
wrapper_sha="$(sha256 "$wrapper")"
expected_uploads="$(printf '%s\n' "PUT @redbench/${target}@${version} ${native_sha}" "PUT redbench@${version} ${wrapper_sha}")"
actual_uploads="$(bench_search_ere -- '^PUT ' "$request_file")"
[[ "$actual_uploads" == "$expected_uploads" ]] || { printf 'offline smoke: registry upload order or stored digests are wrong\n' >&2; exit 1; }
HOME="$registry_home" BENCH_HOME="$registry_home/.bench" BENCH_NO_REPAIR="$repair_disabled" NODE_OPTIONS="--require=$root/scripts/offline-network-sentinel.cjs" BENCH_OFFLINE_EGRESS_LOG="$egress_log" BENCH_OFFLINE_ALLOWED_ORIGIN="$registry_origin" npm_config_cache="$tmp/registry-cache" npm_config_registry="$registry" npm_config_audit=false npm_config_fund=false npm_config_update_notifier=false npm_config_fetch_retries=0 npm_config_proxy='' npm_config_https_proxy='' NO_PROXY='*' npm install --ignore-scripts --omit=optional --prefix "$registry_prefix" --registry "$registry" "redbench@${version}" "@redbench/${target}@${version}" >/dev/null
[[ ! -s "$egress_log" ]] || { printf 'offline smoke: registry install attempted undeclared egress: %s\n' "$(tr '\n' ' ' < "$egress_log")" >&2; exit 1; }
bench_search_ere -q -- '^GET ' "$request_file" || { printf 'offline smoke: registry fixture observed no package requests\n' >&2; exit 1; }
registry_installed="$registry_prefix/node_modules/redbench/bin/bench.sh"
[[ -x "$registry_installed" ]] || { printf 'offline smoke: registry install did not produce wrapper\n' >&2; exit 1; }
registry_version="$(HOME="$registry_home" BENCH_HOME="$registry_home/.bench" BENCH_OFFLINE=1 BENCH_NO_REPAIR=1 native_offline_process registry-runtime bash "$registry_installed" version)"
[[ "$registry_version" == "bench ${version} "* ]] || { printf 'offline smoke: registry version output is wrong: %s\n' "$registry_version" >&2; exit 1; }
HOME="$registry_home" BENCH_HOME="$registry_home/.bench" BENCH_OFFLINE=1 BENCH_NO_REPAIR=1 native_offline_process registry-runtime bash "$registry_installed" commands --brief > "$tmp/registry-commands"
bench_search_ere -q -- '^commands --brief$' "$tmp/registry-commands" || { printf 'offline smoke: registry operational probe failed\n' >&2; exit 1; }
HOME="$registry_home" BENCH_HOME="$registry_home/.bench" BENCH_NO_REPAIR=1 NODE_OPTIONS="--require=$root/scripts/offline-network-sentinel.cjs" BENCH_OFFLINE_EGRESS_LOG="$egress_log" npm_config_registry=http://127.0.0.1:9 npm_config_offline=true npm_config_audit=false npm_config_fund=false npm_config_update_notifier=false npm_config_fetch_retries=0 npm_config_proxy='' npm_config_https_proxy='' NO_PROXY='*' npm uninstall --offline --ignore-scripts --prefix "$registry_prefix" redbench "@redbench/${target}" >/dev/null
[[ ! -s "$egress_log" ]] || { printf 'offline smoke: registry uninstall attempted undeclared egress: %s\n' "$(tr '\n' ' ' < "$egress_log")" >&2; exit 1; }
[[ ! -e "$registry_prefix/node_modules/redbench" && ! -e "$registry_prefix/node_modules/@redbench/${target}" && ! -e "$registry_home/.bench" ]] || { printf 'offline smoke: registry uninstall left client residue\n' >&2; exit 1; }
if [[ -n "$registry_pid" ]]; then
  kill "$registry_pid"
  wait "$registry_pid" 2>/dev/null || true
  registry_pid=""
fi
printf 'offline smoke: direct, local-npm, and loopback-registry journeys passed for %s\n' "$target"
