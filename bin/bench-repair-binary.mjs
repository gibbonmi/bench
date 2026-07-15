#!/usr/bin/env node
import { createHash, randomUUID } from "node:crypto";
import { createGunzip } from "node:zlib";
import { access, mkdir, open, rename, rm } from "node:fs/promises";
import { dirname, join } from "node:path";
import { Readable } from "node:stream";

const [kit, pkgName, version, platform] = process.argv.slice(2);

if (!kit || !pkgName || !version || !platform) {
  console.error("bench: repair failed: internal launcher arguments missing");
  process.exit(1);
}
if (!/^[0-9]+\.[0-9]+\.[0-9]+(?:[+-][0-9A-Za-z.-]+)?$/.test(version)) {
  console.error(`bench: repair failed: invalid kit version ${JSON.stringify(version)}`);
  process.exit(1);
}

const benchHome = process.env.BENCH_HOME || join(process.env.HOME || "", ".bench");
const target = join(benchHome, "cache", "bin", version, platform, "bench");
const targetDir = dirname(target);
let tempPath;

for (const [signal, code] of [
  ["SIGINT", 130],
  ["SIGTERM", 143],
]) {
  process.once(signal, () => {
    const tmp = tempPath;
    const done = tmp ? rm(tmp, { force: true }).catch(() => {}) : Promise.resolve();
    done.finally(() => process.exit(code));
  });
}

try {
  const registry = normalizeRegistry(
    process.env.BENCH_NPM_REGISTRY || process.env.npm_config_registry || "https://registry.npmjs.org",
  );
  const metadata = await fetchJSON(`${registry}/${encodePackageName(pkgName)}`);
  const dist = metadata?.versions?.[version]?.dist;
  if (!dist?.tarball || !dist?.integrity) {
    throw new Error(`registry metadata did not include ${pkgName}@${version}`);
  }
  const digest = dist.integrity.startsWith("sha512-")
    ? dist.integrity.slice("sha512-".length, "sha512-".length + 12)
    : dist.integrity.slice(0, 12);
  console.error(`bench: installing ${pkgName}@${version} sha512:${digest}`);
  const tgz = await fetchBytes(dist.tarball);
  verifyIntegrity(tgz, dist.integrity);
  const binary = await extractBenchBinary(tgz);
  await mkdir(targetDir, { recursive: true });
  console.error(`bench: created ${targetDir}`);
  const tmp = join(targetDir, `.bench-${process.pid}-${randomUUID()}.tmp`);
  tempPath = tmp;
  try {
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
      const ready = process.env.BENCH_TEST_REPAIR_READY_FILE;
      const marker = await open(ready, "w");
      await marker.close();
      while (true) {
        try {
          await access(ready);
          await new Promise((resolve) => setTimeout(resolve, 25));
        } catch {
          break;
        }
      }
    }
    await rename(tmp, target);
    tempPath = undefined;
    console.error(`bench: wrote ${target}`);
  } catch (err) {
    await rm(tmp, { force: true }).catch(() => {});
    throw err;
  }
} catch (err) {
  await rm(target, { force: true }).catch(() => {});
  console.error(`bench: repair failed: ${err.message}`);
  process.exit(1);
}

function normalizeRegistry(value) {
  return value.replace(/\/+$/, "");
}

function encodePackageName(name) {
  if (name.startsWith("@")) {
    const [scope, pkg] = name.split("/");
    return `${scope}%2f${pkg}`;
  }
  return encodeURIComponent(name);
}

async function fetchJSON(url) {
  const bytes = await fetchBytes(url);
  return JSON.parse(bytes.toString("utf8"));
}

async function fetchBytes(url) {
  const res = await fetch(url);
  if (!res.ok) {
    throw new Error(`fetch ${url} returned ${res.status}`);
  }
  return Buffer.from(await res.arrayBuffer());
}

// The integrity hash comes from the same registry as the tarball, so this
// defends against corruption and transport tampering, not registry compromise —
// that threat is owned by npm provenance/2FA upstream, not this launcher.
function verifyIntegrity(bytes, integrity) {
  const match = /^sha512-([A-Za-z0-9+/=]+)$/.exec(integrity);
  if (!match) {
    throw new Error(`unsupported integrity ${integrity}`);
  }
  const got = createHash("sha512").update(bytes).digest("base64");
  if (got !== match[1]) {
    throw new Error("integrity mismatch");
  }
}

async function extractBenchBinary(tgz) {
  const tarBytes = await gunzip(tgz);
  let offset = 0;
  while (offset + 512 <= tarBytes.length) {
    const header = tarBytes.subarray(offset, offset + 512);
    if (header.every((b) => b === 0)) {
      break;
    }
    const name = readString(header, 0, 100);
    const prefix = readString(header, 345, 155);
    const fullName = prefix ? `${prefix}/${name}` : name;
    const sizeText = readString(header, 124, 12).trim();
    const size = sizeText ? Number.parseInt(sizeText, 8) : 0;
    if (!Number.isFinite(size) || size < 0) {
      throw new Error(`invalid tar size for ${fullName}`);
    }
    offset += 512;
    if (offset + size > tarBytes.length) {
      throw new Error(`truncated tar entry for ${fullName}`);
    }
    const data = tarBytes.subarray(offset, offset + size);
    if (fullName === "package/bin/bench") {
      if (data.length === 0) {
        throw new Error("package/bin/bench was empty");
      }
      return data;
    }
    const paddedSize = Math.ceil(size / 512) * 512;
    if (offset + paddedSize > tarBytes.length) {
      throw new Error(`truncated tar padding for ${fullName}`);
    }
    offset += paddedSize;
  }
  throw new Error("tarball did not contain package/bin/bench");
}

async function gunzip(bytes) {
  const chunks = [];
  for await (const chunk of Readable.from(bytes).pipe(createGunzip())) {
    chunks.push(chunk);
  }
  return Buffer.concat(chunks);
}

function readString(buf, start, length) {
  const raw = buf.subarray(start, start + length);
  const nul = raw.indexOf(0);
  return raw.subarray(0, nul === -1 ? raw.length : nul).toString("utf8");
}
