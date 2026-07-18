#!/usr/bin/env bash
# Exercise direct, local-tarball npm, and loopback-registry journeys from one
# extracted archive with empty state and public egress unavailable.
set -euo pipefail

artifacts="${1:?usage: smoke-offline.sh <artifact-dir>}"
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
target="$(node -e 'const rows=require(process.argv[1]); const [os,arch]=process.argv.slice(2); const row=rows.find(item=>item.os===os&&item.arch===arch); if(row) process.stdout.write(`${row.os}-${row.arch}`)' "$root/scripts/platforms.json" "$host_os" "$host_arch")"
[[ -n "$target" ]] || { printf 'offline smoke: unsupported host %s/%s\n' "$host_os" "$host_arch" >&2; exit 2; }
runtime_target="$(node -e 'const rows=require(process.argv[1]); const target=process.argv[2]; const row=rows.find(item=>`${item.os}-${item.arch}`===target); if(row) process.stdout.write(`${row.goos}/${row.goarch}`)' "$root/scripts/platforms.json" "$target")"
[[ -n "$runtime_target" ]] || { printf 'offline smoke: target is absent from canonical matrix: %s\n' "$target" >&2; exit 1; }

wrapper="$artifacts/redbench-${version}.tgz"
native="$artifacts/redbench-${target}-${version}.tgz"
archive="$artifacts/redbench-${version}-${target}.tar.gz"
for input in "$wrapper" "$native" "$archive"; do
  [[ -f "$input" && ! -L "$input" && -s "$input" ]] || { printf 'offline smoke: required artifact is missing or unsafe: %s\n' "$input" >&2; exit 1; }
done

tmp="$(mktemp -d)"
registry_pid=""
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
tar -tzf "$archive" > "$tar_members" || { printf 'offline smoke: archive cannot be read\n' >&2; exit 1; }
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
  "$bundle_root/evidence/components/wrapper-component-manifest.json" \
  "$bundle_root/evidence/components/platform-component-manifest.json"; do
  [[ -f "$required" && ! -L "$required" && -s "$required" ]] || { printf 'offline smoke: archive inventory is incomplete: %s\n' "$required" >&2; exit 1; }
done
cmp -s "$wrapper" "$bundle_root/packages/redbench-${version}.tgz" || { printf 'offline smoke: wrapper bytes were substituted\n' >&2; exit 1; }
cmp -s "$native" "$bundle_root/packages/redbench-${target}-${version}.tgz" || { printf 'offline smoke: platform bytes were substituted\n' >&2; exit 1; }
native_extract="$tmp/native"
mkdir -p "$native_extract"
tar -xzf "$native" -C "$native_extract" package/bin/bench
cmp -s "$native_extract/package/bin/bench" "$bundle_root/bin/bench" || { printf 'offline smoke: archive binary differs from platform package\n' >&2; exit 1; }

if [[ -f "$artifacts/../preflight/release-index.json" && -f "$artifacts/../preflight/SHA256SUMS" ]]; then
  node -e '
    const fs=require("fs"), crypto=require("crypto");
    const [indexPath,sumsPath,artifactDir,archiveName]=process.argv.slice(1);
    const index=JSON.parse(fs.readFileSync(indexPath));
    const sums=new Map(fs.readFileSync(sumsPath,"utf8").trim().split("\n").map(line=>line.split("  ")));
    const data=fs.readFileSync(`${artifactDir}/${archiveName}`);
    const digest=crypto.createHash("sha256").update(data).digest("hex");
    if (!index.artifacts.some(item=>item.name===archiveName && item.sha256===digest) || sums.get(archiveName)!==digest) throw new Error("archive digest is not bound by the supplied release evidence");
  ' "$artifacts/../preflight/release-index.json" "$artifacts/../preflight/SHA256SUMS" "$artifacts" "$(basename "$archive")" || { printf 'offline smoke: supplied release evidence does not bind archive bytes\n' >&2; exit 1; }
fi

direct_home="$tmp/direct-home"
mkdir -p "$direct_home" "$tmp/empty-path"
direct_version="$(env -i HOME="$direct_home" BENCH_HOME="$direct_home/.bench" PATH="$tmp/empty-path" BENCH_NO_REPAIR=1 "$bundle_root/bin/bench" version)"
[[ "$direct_version" == "benchkit ${version} (${runtime_target})" ]] || { printf 'offline smoke: direct version output is wrong: %s\n' "$direct_version" >&2; exit 1; }
env -i HOME="$direct_home" BENCH_HOME="$direct_home/.bench" PATH="$tmp/empty-path" BENCH_NO_REPAIR=1 "$bundle_root/bin/bench" commands --brief > "$tmp/direct-commands"
rg -q '^commands --brief$' "$tmp/direct-commands" || { printf 'offline smoke: direct operational probe failed\n' >&2; exit 1; }

local_home="$tmp/local-home"
local_prefix="$tmp/local-prefix"
mkdir -p "$local_home" "$local_prefix"
printf '{"private":true}\n' > "$local_prefix/package.json"
HOME="$local_home" BENCH_HOME="$local_home/.bench" BENCH_NO_REPAIR=1 npm_config_cache="$tmp/local-cache" npm_config_offline=true npm_config_registry=http://127.0.0.1:9 npm install --offline --ignore-scripts --omit=optional --prefix "$local_prefix" "$bundle_root/packages/redbench-${version}.tgz" "$bundle_root/packages/redbench-${target}-${version}.tgz" >/dev/null
installed="$local_prefix/node_modules/redbench/bin/bench.sh"
[[ -x "$installed" ]] || { printf 'offline smoke: local npm install did not produce wrapper\n' >&2; exit 1; }
local_version="$(HOME="$local_home" BENCH_HOME="$local_home/.bench" BENCH_NO_REPAIR=1 bash "$installed" version)"
[[ "$local_version" == "benchkit ${version} "* ]] || { printf 'offline smoke: local npm version output is wrong: %s\n' "$local_version" >&2; exit 1; }
HOME="$local_home" BENCH_HOME="$local_home/.bench" BENCH_NO_REPAIR=1 bash "$installed" help > "$tmp/local-help"
rg -q '^bench —' "$tmp/local-help" || { printf 'offline smoke: local npm operational probe failed\n' >&2; exit 1; }
HOME="$local_home" BENCH_HOME="$local_home/.bench" npm_config_registry=http://127.0.0.1:9 npm_config_offline=true npm uninstall --offline --ignore-scripts --prefix "$local_prefix" redbench "@redbench/${target}" >/dev/null
[[ ! -e "$local_prefix/node_modules/redbench" && ! -e "$local_home/.bench/cache" ]] || { printf 'offline smoke: local npm uninstall left package or Bench cache residue\n' >&2; exit 1; }

store="$tmp/registry-store"
mkdir -p "$store"
cp "$native" "$store/"
printf '%s\n' "@redbench/${target}@${version}" > "$tmp/registry-upload-order"
cp "$wrapper" "$store/"
printf '%s\n' "redbench@${version}" >> "$tmp/registry-upload-order"
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
HOME="$registry_home" BENCH_HOME="$registry_home/.bench" BENCH_NO_REPAIR=1 npm_config_cache="$tmp/registry-cache" npm_config_registry="$registry" npm_config_audit=false npm_config_fund=false npm_config_update_notifier=false npm_config_proxy= npm_config_https_proxy= NO_PROXY='*' npm install --ignore-scripts --omit=optional --prefix "$registry_prefix" --registry "$registry" "redbench@${version}" "@redbench/${target}@${version}" >/dev/null
cmp -s "$tmp/registry-upload-order" <(printf '%s\n' "@redbench/${target}@${version}" "redbench@${version}") || { printf 'offline smoke: registry upload was not platform-first and wrapper-last\n' >&2; exit 1; }
rg -q '^GET ' "$request_file" || { printf 'offline smoke: registry fixture observed no package requests\n' >&2; exit 1; }
registry_installed="$registry_prefix/node_modules/redbench/bin/bench.sh"
[[ -x "$registry_installed" ]] || { printf 'offline smoke: registry install did not produce wrapper\n' >&2; exit 1; }
registry_version="$(HOME="$registry_home" BENCH_HOME="$registry_home/.bench" BENCH_NO_REPAIR=1 bash "$registry_installed" version)"
[[ "$registry_version" == "benchkit ${version} "* ]] || { printf 'offline smoke: registry version output is wrong: %s\n' "$registry_version" >&2; exit 1; }
HOME="$registry_home" BENCH_HOME="$registry_home/.bench" npm_config_registry=http://127.0.0.1:9 npm_config_offline=true npm uninstall --offline --ignore-scripts --prefix "$registry_prefix" redbench "@redbench/${target}" >/dev/null
[[ ! -e "$registry_prefix/node_modules/redbench" && ! -e "$registry_home/.bench/cache" ]] || { printf 'offline smoke: registry uninstall left client residue\n' >&2; exit 1; }
if [[ -n "$registry_pid" ]]; then
  kill "$registry_pid"
  wait "$registry_pid" 2>/dev/null || true
  registry_pid=""
fi
printf 'offline smoke: direct, local-npm, and loopback-registry journeys passed for %s\n' "$target"
