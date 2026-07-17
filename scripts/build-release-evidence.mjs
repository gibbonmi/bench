import fs from "node:fs";
import path from "node:path";
import crypto from "node:crypto";

const [root, wrapperDir, packagesDir] = process.argv.slice(2);
if (!root || !wrapperDir || !packagesDir) throw new Error("usage: build-release-evidence.mjs <source-root> <wrapper-dir> <packages-dir>");

const matrix = readJSON(path.join(root, "scripts/platforms.json"));
const assets = readJSON(path.join(root, "scripts/wrapper-assets.json"));
const sourcePackage = readJSON(path.join(root, "package.json"));
const seenDestinations = new Set();
const byteOrder = (a, b) => a.name < b.name ? -1 : a.name > b.name ? 1 : 0;

if (!Array.isArray(matrix) || matrix.length !== 4 || !Array.isArray(assets)) throw new Error("release evidence matrix or asset registry is invalid");

function readJSON(file) {
  const stat = fs.lstatSync(file);
  if (!stat.isFile() || stat.isSymbolicLink()) throw new Error(`release evidence input is not a regular file: ${file}`);
  const bytes = fs.readFileSync(file);
  if (bytes.length === 0) throw new Error(`release evidence input is empty: ${file}`);
  return JSON.parse(bytes.toString("utf8"));
}

function rejectControl(value, label) {
  for (const byte of Buffer.from(value)) {
    if (byte < 0x20 || byte === 0x7f) throw new Error(`${label} contains control bytes`);
  }
}

function safeRelative(value, label) {
  if (value.includes("\\")) throw new Error(`${label} is unsafe: ${value}`);
  const normalized = path.posix.normalize(value.replaceAll(path.sep, "/"));
  rejectControl(normalized, label);
  if (!normalized || normalized === "." || normalized.startsWith("../") || normalized.includes("/../") || path.posix.isAbsolute(normalized)) {
    throw new Error(`${label} is unsafe: ${value}`);
  }
  return normalized;
}

function sourcePath(relative) {
  const safe = safeRelative(relative, "release evidence path");
  const file = path.resolve(root, safe);
  const relativeToRoot = path.relative(root, file);
  if (relativeToRoot.startsWith("..") || path.isAbsolute(relativeToRoot)) throw new Error(`release evidence path escapes source root: ${relative}`);
  return file;
}

function modeFrom(value) {
  if (typeof value !== "string" || !/^[0-7]{3,4}$/.test(value)) throw new Error(`release evidence mode is invalid: ${value}`);
  const mode = Number.parseInt(value, 8);
  if ((mode & 0o7000) !== 0 || mode !== 0o644 && mode !== 0o755) throw new Error(`release evidence mode is unsafe: ${value}`);
  return mode;
}

function copyRegular(src, dst, mode, label, destinationRoot = wrapperDir) {
  const stat = fs.lstatSync(src);
  if (!stat.isFile() || stat.isSymbolicLink()) throw new Error(`${label} is not a regular file: ${src}`);
  if (stat.size === 0) throw new Error(`${label} is empty: ${src}`);
  if ((stat.mode & 0o7000) !== 0) throw new Error(`${label} has unsafe mode: ${src}`);
  const bytes = fs.readFileSync(src);
  if ((src.includes(`${path.sep}governance${path.sep}`) || path.basename(src) === "THIRD_PARTY_NOTICES.txt") && bytes[bytes.length - 1] !== 0x0a) {
    throw new Error(`${label} is missing a final newline: ${src}`);
  }
  const rel = safeRelative(path.relative(destinationRoot, dst), "package destination");
  if (seenDestinations.has(rel)) throw new Error(`duplicate normalized package path: ${rel}`);
  seenDestinations.add(rel);
  fs.mkdirSync(path.dirname(dst), {recursive: true});
  fs.writeFileSync(dst, bytes);
  fs.chmodSync(dst, mode);
}

function validateRequiredSources() {
  const required = ["LICENSE", "governance/THIRD_PARTY_NOTICES.txt", "governance/sbom.spdx.json", ...[
    "supported-versions.json", "security-response.json", "dependency-license-change.json", "threat-model.json", "recovery-rollback.json", "support.json",
  ].map(name => `governance/policies/${name}`)];
  for (const relative of required) {
    const source = sourcePath(relative);
    const stat = fs.lstatSync(source);
    if (!stat.isFile() || stat.isSymbolicLink() || stat.size === 0 || (stat.mode & 0o7000) !== 0) throw new Error(`required release evidence source is invalid: ${relative}`);
    if ((relative.startsWith("governance/") || relative === "LICENSE") && fs.readFileSync(source)[stat.size - 1] !== 0x0a) throw new Error(`required release evidence source is missing a final newline: ${relative}`);
  }
}

function copyTree(src, dst, mode, label, destinationRoot = wrapperDir) {
  const stat = fs.lstatSync(src);
  if (!stat.isDirectory() || stat.isSymbolicLink()) throw new Error(`${label} is not a real directory: ${src}`);
  for (const name of fs.readdirSync(src, {withFileTypes: true}).sort(byteOrder)) {
    rejectControl(name.name, `${label} entry`);
    const child = path.join(src, name.name);
    const target = path.join(dst, name.name);
    if (name.isDirectory()) copyTree(child, target, mode, label, destinationRoot);
    else if (name.isFile()) copyRegular(child, target, mode, label, destinationRoot);
    else throw new Error(`${label} contains an unsafe file: ${child}`);
  }
}

function copyAssets(destination) {
  seenDestinations.clear();
  for (const asset of assets) {
    if (!asset || typeof asset.source !== "string") throw new Error("release evidence asset has no source");
    const mode = modeFrom(asset.mode);
    const source = sourcePath(asset.source);
    const target = path.join(destination, safeRelative(asset.source, "package destination"));
    if (asset.tree) copyTree(source, target, mode, "wrapper allowlist tree");
    else copyRegular(source, target, mode, "wrapper allowlist input");
  }
}

function writeJSON(file, value) {
  fs.writeFileSync(file, JSON.stringify(value, null, 2) + "\n", {mode: 0o644});
}

function sha256(bytes) {
  return crypto.createHash("sha256").update(bytes).digest("hex");
}

function packageFiles(dir) {
  const files = [];
  const walk = (current) => {
    for (const name of fs.readdirSync(current, {withFileTypes: true}).sort(byteOrder)) {
      const file = path.join(current, name.name);
      if (name.isDirectory()) walk(file);
      else if (name.isFile() && name.name !== "component-manifest.json") {
        const rel = safeRelative(path.relative(dir, file), "component manifest path");
        const stat = fs.lstatSync(file);
        files.push({path: rel, mode: (stat.mode & 0o777).toString(8), size: stat.size, sha256: sha256(fs.readFileSync(file))});
      } else throw new Error(`component package contains unsafe file: ${file}`);
    }
  };
  walk(dir);
  files.sort((a, b) => a.path < b.path ? -1 : a.path > b.path ? 1 : 0);
  return files;
}

function writeEvidence(dir, name, version, target) {
  const sbom = readJSON(path.join(root, "governance", "sbom.spdx.json"));
  sbom.name = `${name}-release`;
  sbom.documentNamespace = `https://github.com/gibbonmi/bench/releases/sbom/${name}`;
  sbom.packages[0].name = name;
  sbom.packages[0].versionInfo = version;
  sbom.packages[0].SPDXID = "SPDXRef-Package-" + name.replaceAll(/[^A-Za-z0-9.-]/g, "-");
  sbom.relationships[0].relatedSpdxElement = sbom.packages[0].SPDXID;
  writeJSON(path.join(dir, "governance", "sbom.spdx.json"), sbom);
  const component = {schema_version: 1, component: {name, version, target: target || null}, files: packageFiles(dir)};
  writeJSON(path.join(dir, "component-manifest.json"), component);
}

validateRequiredSources();
copyAssets(wrapperDir);
const optionalDependencies = Object.fromEntries(matrix.map(p => [`@redbench/${p.os}-${p.arch}`, sourcePackage.version]));
const wrapperPackage = {...sourcePackage, optionalDependencies};
wrapperPackage.scripts = {...(sourcePackage.scripts || {})};
delete wrapperPackage.scripts.prepare;
wrapperPackage.files = [...assets.map(a => a.source), "component-manifest.json"];
writeJSON(path.join(wrapperDir, "package.json"), wrapperPackage);
writeEvidence(wrapperDir, sourcePackage.name, sourcePackage.version, {os: "all", arch: "all"});

for (const p of matrix) {
  const dir = path.join(packagesDir, `${p.os}-${p.arch}`);
  fs.mkdirSync(path.join(dir, "bin"), {recursive: true});
  seenDestinations.clear();
  copyRegular(path.join(root, "LICENSE"), path.join(dir, "LICENSE"), 0o644, "platform license", dir);
  copyTree(path.join(root, "governance"), path.join(dir, "governance"), 0o644, "platform governance", dir);
  const pkg = {
    name: `@redbench/${p.os}-${p.arch}`,
    version: sourcePackage.version,
    description: `benchkit prebuilt binary for ${p.os}-${p.arch}`,
    bin: "bin/bench",
    os: [p.os],
    cpu: [p.arch],
    license: sourcePackage.license,
  };
  writeJSON(path.join(dir, "package.json"), pkg);
  writeEvidence(dir, pkg.name, pkg.version, {os: p.os, arch: p.arch});
}
