#!/usr/bin/env node
import { createHash, randomUUID } from "node:crypto";
import { createGunzip } from "node:zlib";
import { access, lstat, mkdir, open, readdir, readFile, rename, rm } from "node:fs/promises";
import { dirname, join } from "node:path";
import { Readable } from "node:stream";

// This script runs before a Go binary exists. These Node-runtime resource facts
// deliberately do not import internal/bounds.
const FETCH_DEADLINE_MS = 60_000;
const DOWNLOAD_LIMIT = 100 * 1024 * 1024;
const DECOMPRESSED_LIMIT = 200 * 1024 * 1024;
// Cross-runtime build constant, parity-pinned to releaseevidence requirements.
const PIN_MANIFEST_PATH = "binary-pins.json";
const FETCH_DEADLINE_LABEL = `${FETCH_DEADLINE_MS / 1000}-second`;
const DOWNLOAD_LIMIT_LABEL = `${DOWNLOAD_LIMIT / (1024 * 1024)} MB`;
const DECOMPRESSED_LIMIT_LABEL = `${DECOMPRESSED_LIMIT / (1024 * 1024)} MB`;

const [mode, kit, pkgName, version, platform] = process.argv.slice(2);
if (!kit || !pkgName || !version || !platform || !["repair", "prune"].includes(mode)) {
  console.error("bench: repair failed: internal launcher arguments missing");
  process.exit(1);
}
if (!/^[0-9]+\.[0-9]+\.[0-9]+(?:[+-][0-9A-Za-z.-]+)?$/.test(version)) {
  console.error(`bench: repair failed: invalid kit version ${JSON.stringify(version)}`);
  process.exit(1);
}

const benchHome = process.env.BENCH_HOME || join(process.env.HOME || "", ".bench");
const cacheRoot = join(benchHome, "cache", "bin");
const target = join(cacheRoot, version, platform, "bench");
const targetDir = dirname(target);
let tempPath;
let interrupted = false;

for (const [signal, code] of [["SIGINT", 130], ["SIGTERM", 143]]) {
  process.once(signal, () => {
    if (interrupted) return;
    interrupted = true;
    const tmp = tempPath;
    const done = tmp ? rm(tmp, { force: true }).catch(() => {}) : Promise.resolve();
    done.finally(() => process.exit(code));
  });
}

try {
  if (mode === "prune") {
    await pruneCache();
  } else {
    await install();
  }
} catch (err) {
  if (tempPath) await rm(tempPath, { force: true }).catch(() => {});
  console.error(`bench: repair failed: ${err.message}`);
  process.exit(1);
}

async function install() {
  if (process.env.BENCH_TEST_REPAIR_START_READY_FILE) await waitForTestRelease(process.env.BENCH_TEST_REPAIR_START_READY_FILE);
  const integrity = await expectedIntegrity();
  const registry = normalizeRegistry(process.env.BENCH_NPM_REGISTRY || process.env.npm_config_registry || "https://registry.npmjs.org");
  const signal = AbortSignal.timeout(testNumber("BENCH_TEST_REPAIR_DEADLINE_MS", FETCH_DEADLINE_MS));
  const metadata = await fetchJSON(`${registry}/${encodePackageName(pkgName)}`, signal);
  const dist = metadata?.versions?.[version]?.dist;
  if (!dist?.tarball) throw new Error(`registry metadata did not include ${pkgName}@${version}`);
  const digest = integrity.slice("sha512-".length, "sha512-".length + 12);
  console.error(`bench: installing ${pkgName}@${version} sha512:${digest}`);
  const tgz = await fetchBytes(dist.tarball, signal, testNumber("BENCH_TEST_REPAIR_DOWNLOAD_LIMIT", DOWNLOAD_LIMIT));
  verifyIntegrity(tgz, integrity);
  const binary = await extractBenchBinary(tgz, testNumber("BENCH_TEST_REPAIR_DECOMPRESSED_LIMIT", DECOMPRESSED_LIMIT));
  await mkdir(targetDir, { recursive: true });
  console.error(`bench: created ${targetDir}`);
  const tmp = join(targetDir, `.bench-${process.pid}-${randomUUID()}.tmp`);
  tempPath = tmp;
  const fh = await open(tmp, "wx", 0o755);
  try {
    await fh.writeFile(binary);
    await fh.chmod(0o755);
    await fh.sync();
  } finally {
    await fh.close();
  }
  console.error(`bench: wrote ${tmp}`);
  if (process.env.BENCH_TEST_REPAIR_READY_FILE) {
    await waitForTestRelease(process.env.BENCH_TEST_REPAIR_READY_FILE);
  }
  if (process.env.BENCH_TEST_REPAIR_FAIL_AFTER_READY === "1") throw new Error("injected failure after repair-ready synchronization");
  await rename(tmp, target);
  tempPath = undefined;
  console.error(`bench: wrote ${target}`);
}

async function waitForTestRelease(ready) {
  const marker = await open(ready, "w");
  await marker.close();
  while (true) {
    try { await access(ready); await new Promise(resolve => setTimeout(resolve, 25)); }
    catch { return; }
  }
}

async function expectedIntegrity() {
  const manifestFile = join(kit, PIN_MANIFEST_PATH);
  let bytes;
  try { bytes = await readFile(manifestFile); }
  catch { throw new Error(`pin manifest ${PIN_MANIFEST_PATH} is missing`); }
  if (bytes.length === 0) throw new Error(`pin manifest ${PIN_MANIFEST_PATH} is empty`);
  let manifest;
  try { manifest = JSON.parse(bytes); }
  catch { throw new Error(`pin manifest ${PIN_MANIFEST_PATH} is unparseable`); }
  if (!manifest || manifest.schema_version !== 1 || !Array.isArray(manifest.pins)) throw new Error(`pin manifest ${PIN_MANIFEST_PATH} has invalid schema`);
  for (const item of manifest.pins) {
    if (!item || typeof item !== "object" || Array.isArray(item) || typeof item.name !== "string" || !item.name || typeof item.version !== "string" || !item.version || typeof item.integrity !== "string" || !/^sha512-[A-Za-z0-9+/=]+$/.test(item.integrity)) throw new Error(`pin manifest ${PIN_MANIFEST_PATH} has a malformed entry`);
  }
  const matches = manifest.pins.filter(item => item.name === pkgName && item.version === version);
  if (matches.length === 0) throw new Error(`pin manifest ${PIN_MANIFEST_PATH} has no entry for ${pkgName}@${version}`);
  if (matches.length !== 1) throw new Error(`pin manifest ${PIN_MANIFEST_PATH} has ambiguous entries for ${pkgName}@${version}`);
  const [pin] = matches;
  return pin.integrity;
}

async function pruneCache() {
  let versions;
  try { versions = await readdir(cacheRoot, { withFileTypes: true }); }
  catch (err) {
    if (err.code === "ENOENT") return console.error("bench: repair prune: no stale cache entries");
    throw err;
  }
  const removed = [];
  for (const versionEntry of versions) {
    const versionPath = join(cacheRoot, versionEntry.name);
    if (versionEntry.isSymbolicLink()) continue;
    if (!versionEntry.isDirectory()) {
      if (versionEntry.isFile()) { await rm(versionPath); removed.push(versionEntry.name); }
      continue;
    }
    if (versionEntry.name !== version) {
      await rm(versionPath, { recursive: true }); removed.push(versionEntry.name); continue;
    }
    for (const platformEntry of await readdir(versionPath, { withFileTypes: true })) {
      if (platformEntry.isSymbolicLink() || platformEntry.name === platform && platformEntry.isDirectory()) continue;
      const stale = join(versionPath, platformEntry.name);
      const stat = await lstat(stale);
      if (stat.isFile() || stat.isDirectory()) { await rm(stale, { recursive: stat.isDirectory() }); removed.push(`${version}/${platformEntry.name}`); }
    }
  }
  if (removed.length === 0) console.error("bench: repair prune: no stale cache entries");
  else for (const entry of removed.sort()) console.error(`bench: repair prune: removed ${entry}`);
}

function normalizeRegistry(value) { return value.replace(/\/+$/, ""); }
function encodePackageName(name) {
  if (name.startsWith("@")) { const [scope, pkg] = name.split("/"); return `${scope}%2f${pkg}`; }
  return encodeURIComponent(name);
}
async function fetchJSON(url, signal) {
  const bytes = await fetchBytes(url, signal, DOWNLOAD_LIMIT);
  try { return JSON.parse(bytes.toString("utf8")); }
  catch (err) { throw new Error(`registry metadata was unparseable: ${err.message}`); }
}
async function fetchBytes(url, signal, limit) {
  let res;
  try { res = await fetch(url, { signal }); }
  catch (err) {
    if (signal.aborted) throw new Error(`${FETCH_DEADLINE_LABEL} total repair deadline exceeded`);
    throw err;
  }
  if (!res.ok) throw new Error(`fetch ${url} returned ${res.status}`);
  if (!res.body) return Buffer.alloc(0);
  const chunks = []; let size = 0;
  try {
    for await (const chunk of Readable.fromWeb(res.body)) {
      size += chunk.length;
      if (size > limit) throw new Error(`${DOWNLOAD_LIMIT_LABEL} download limit exceeded`);
      chunks.push(chunk);
    }
  } catch (err) {
    if (signal.aborted) throw new Error(`${FETCH_DEADLINE_LABEL} total repair deadline exceeded`);
    throw err;
  }
  return Buffer.concat(chunks, size);
}
function verifyIntegrity(bytes, integrity) {
  const match = /^sha512-([A-Za-z0-9+/=]+)$/.exec(integrity);
  if (!match) throw new Error(`unsupported integrity ${integrity}`);
  const got = createHash("sha512").update(bytes).digest("base64");
  if (got !== match[1]) throw new Error("pin manifest integrity mismatch");
}
async function extractBenchBinary(tgz, limit) {
  const tarBytes = await gunzip(tgz, limit);
  let offset = 0;
  while (offset + 512 <= tarBytes.length) {
    const header = tarBytes.subarray(offset, offset + 512);
    if (header.every(b => b === 0)) break;
    const name = readString(header, 0, 100), prefix = readString(header, 345, 155);
    const fullName = prefix ? `${prefix}/${name}` : name;
    const sizeText = readString(header, 124, 12).trim();
    const size = sizeText ? Number.parseInt(sizeText, 8) : 0;
    if (!Number.isFinite(size) || size < 0) throw new Error(`invalid tar size for ${fullName}`);
    offset += 512;
    if (offset + size > tarBytes.length) throw new Error(`truncated tar entry for ${fullName}`);
    const data = tarBytes.subarray(offset, offset + size);
    if (fullName === "package/bin/bench") {
      if (data.length === 0) throw new Error("package/bin/bench was empty");
      return data;
    }
    const paddedSize = Math.ceil(size / 512) * 512;
    if (offset + paddedSize > tarBytes.length) throw new Error(`truncated tar padding for ${fullName}`);
    offset += paddedSize;
  }
  throw new Error("tarball did not contain package/bin/bench");
}
async function gunzip(bytes, limit) {
  const chunks = []; let size = 0;
  for await (const chunk of Readable.from(bytes).pipe(createGunzip())) {
    size += chunk.length;
    if (size > limit) throw new Error(`${DECOMPRESSED_LIMIT_LABEL} decompressed limit exceeded after ${size} bytes`);
    chunks.push(chunk);
  }
  return Buffer.concat(chunks, size);
}
function readString(buf, start, length) {
  const raw = buf.subarray(start, start + length), nul = raw.indexOf(0);
  return raw.subarray(0, nul === -1 ? raw.length : nul).toString("utf8");
}
function testNumber(name, fallback) {
  if (!name.startsWith("BENCH_TEST_") || !process.env[name]) return fallback;
  const value = Number(process.env[name]);
  return Number.isSafeInteger(value) && value > 0 ? value : fallback;
}
