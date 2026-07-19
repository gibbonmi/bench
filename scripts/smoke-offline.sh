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
registry_pid=""
egress_log="$tmp/undeclared-egress"
: > "$egress_log"
cleanup() {
  if [[ -n "$registry_pid" ]]; then
    kill "$registry_pid" 2>/dev/null || true
    wait "$registry_pid" 2>/dev/null || true
  fi
  rm -rf "$tmp"
}
trap cleanup EXIT INT TERM HUP

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

node "$root/scripts/verify-release-artifact.mjs" "$evidence_dir/release-index.json" "$evidence_dir/SHA256SUMS" "$archive" || { printf 'offline smoke: supplied release evidence does not bind archive bytes\n' >&2; exit 1; }

direct_home="$tmp/direct-home"
mkdir -p "$direct_home" "$tmp/empty-path"
direct_version="$(env -i HOME="$direct_home" BENCH_HOME="$direct_home/.bench" PATH="$tmp/empty-path" BENCH_NO_REPAIR=1 "$bundle_root/bin/bench" version)"
[[ "$direct_version" == "benchkit ${version} (${runtime_target})" ]] || { printf 'offline smoke: direct version output is wrong: %s\n' "$direct_version" >&2; exit 1; }
env -i HOME="$direct_home" BENCH_HOME="$direct_home/.bench" PATH="$tmp/empty-path" BENCH_NO_REPAIR=1 "$bundle_root/bin/bench" commands --brief > "$tmp/direct-commands"
rg -q '^commands --brief$' "$tmp/direct-commands" || { printf 'offline smoke: direct operational probe failed\n' >&2; exit 1; }
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
HOME="$local_home" BENCH_HOME="$local_home/.bench" BENCH_NO_REPAIR="$repair_disabled" NODE_OPTIONS="--require=$root/scripts/offline-network-sentinel.cjs" BENCH_OFFLINE_EGRESS_LOG="$egress_log" npm_config_cache="$tmp/local-cache" npm_config_offline="$offline_mode" npm_config_registry=http://127.0.0.1:9 npm_config_audit=false npm_config_fund=false npm_config_update_notifier=false npm_config_fetch_retries=0 npm_config_proxy='' npm_config_https_proxy='' NO_PROXY='*' npm install --offline --ignore-scripts --omit=optional --prefix "$local_prefix" "$bundle_root/packages/redbench-${version}.tgz" "$bundle_root/packages/redbench-${target}-${version}.tgz" >/dev/null
[[ ! -s "$egress_log" ]] || { printf 'offline smoke: local npm attempted undeclared egress: %s\n' "$(tr '\n' ' ' < "$egress_log")" >&2; exit 1; }
installed="$local_prefix/node_modules/redbench/bin/bench.sh"
[[ -x "$installed" ]] || { printf 'offline smoke: local npm install did not produce wrapper\n' >&2; exit 1; }
local_version="$(HOME="$local_home" BENCH_HOME="$local_home/.bench" BENCH_NO_REPAIR=1 bash "$installed" version)"
[[ "$local_version" == "benchkit ${version} "* ]] || { printf 'offline smoke: local npm version output is wrong: %s\n' "$local_version" >&2; exit 1; }
HOME="$local_home" BENCH_HOME="$local_home/.bench" BENCH_NO_REPAIR=1 bash "$installed" help > "$tmp/local-help"
rg -q '^bench —' "$tmp/local-help" || { printf 'offline smoke: local npm operational probe failed\n' >&2; exit 1; }
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
actual_uploads="$(rg '^PUT ' "$request_file")"
[[ "$actual_uploads" == "$expected_uploads" ]] || { printf 'offline smoke: registry upload order or stored digests are wrong\n' >&2; exit 1; }
HOME="$registry_home" BENCH_HOME="$registry_home/.bench" BENCH_NO_REPAIR="$repair_disabled" NODE_OPTIONS="--require=$root/scripts/offline-network-sentinel.cjs" BENCH_OFFLINE_EGRESS_LOG="$egress_log" BENCH_OFFLINE_ALLOWED_ORIGIN="$registry_origin" npm_config_cache="$tmp/registry-cache" npm_config_registry="$registry" npm_config_audit=false npm_config_fund=false npm_config_update_notifier=false npm_config_fetch_retries=0 npm_config_proxy='' npm_config_https_proxy='' NO_PROXY='*' npm install --ignore-scripts --omit=optional --prefix "$registry_prefix" --registry "$registry" "redbench@${version}" "@redbench/${target}@${version}" >/dev/null
[[ ! -s "$egress_log" ]] || { printf 'offline smoke: registry install attempted undeclared egress: %s\n' "$(tr '\n' ' ' < "$egress_log")" >&2; exit 1; }
rg -q '^GET ' "$request_file" || { printf 'offline smoke: registry fixture observed no package requests\n' >&2; exit 1; }
registry_installed="$registry_prefix/node_modules/redbench/bin/bench.sh"
[[ -x "$registry_installed" ]] || { printf 'offline smoke: registry install did not produce wrapper\n' >&2; exit 1; }
registry_version="$(HOME="$registry_home" BENCH_HOME="$registry_home/.bench" BENCH_NO_REPAIR=1 bash "$registry_installed" version)"
[[ "$registry_version" == "benchkit ${version} "* ]] || { printf 'offline smoke: registry version output is wrong: %s\n' "$registry_version" >&2; exit 1; }
HOME="$registry_home" BENCH_HOME="$registry_home/.bench" BENCH_NO_REPAIR=1 bash "$registry_installed" commands --brief > "$tmp/registry-commands"
rg -q '^commands --brief$' "$tmp/registry-commands" || { printf 'offline smoke: registry operational probe failed\n' >&2; exit 1; }
HOME="$registry_home" BENCH_HOME="$registry_home/.bench" BENCH_NO_REPAIR=1 NODE_OPTIONS="--require=$root/scripts/offline-network-sentinel.cjs" BENCH_OFFLINE_EGRESS_LOG="$egress_log" npm_config_registry=http://127.0.0.1:9 npm_config_offline=true npm_config_audit=false npm_config_fund=false npm_config_update_notifier=false npm_config_fetch_retries=0 npm_config_proxy='' npm_config_https_proxy='' NO_PROXY='*' npm uninstall --offline --ignore-scripts --prefix "$registry_prefix" redbench "@redbench/${target}" >/dev/null
[[ ! -s "$egress_log" ]] || { printf 'offline smoke: registry uninstall attempted undeclared egress: %s\n' "$(tr '\n' ' ' < "$egress_log")" >&2; exit 1; }
[[ ! -e "$registry_prefix/node_modules/redbench" && ! -e "$registry_prefix/node_modules/@redbench/${target}" && ! -e "$registry_home/.bench" ]] || { printf 'offline smoke: registry uninstall left client residue\n' >&2; exit 1; }
if [[ -n "$registry_pid" ]]; then
  kill "$registry_pid"
  wait "$registry_pid" 2>/dev/null || true
  registry_pid=""
fi
printf 'offline smoke: direct, local-npm, and loopback-registry journeys passed for %s\n' "$target"
