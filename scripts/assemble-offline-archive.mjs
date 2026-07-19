import crypto from "node:crypto";
import fs from "node:fs";
import path from "node:path";
import {archiveEntries, packagedEvidenceRecords, readReleasePlan, readReleaseRequirements} from "./release-plan.mjs";

const [root, npmDir, archiveDir, target, version, binary, wrapperExtract, platformExtract] = process.argv.slice(2);
if (![root, npmDir, archiveDir, target, version, binary, wrapperExtract, platformExtract].every(Boolean)) throw new Error("usage: assemble-offline-archive.mjs <root> <npm-dir> <archive-dir> <target> <version> <binary> <wrapper-extract> <platform-extract>");
const plan = readReleasePlan(root);
if (!plan.targets.some(item => `${item.os}-${item.arch}` === target)) throw new Error(`release plan does not contain ${target}`);
const packageEvidence = packagedEvidenceRecords(readReleaseRequirements(root));
const entries = archiveEntries(plan, target, version, packageEvidence);
const rootName = `redbench-${version}-${target}`;

function copyRegular(source, destination, mode) {
  const info = fs.lstatSync(source);
  if (!info.isFile() || info.isSymbolicLink() || info.size === 0) throw new Error(`offline archive source is missing or unsafe: ${source}`);
  fs.mkdirSync(path.dirname(destination), {recursive: true});
  fs.copyFileSync(source, destination);
  fs.chmodSync(destination, Number.parseInt(mode, 8));
}

function instructions() {
  const wrapper = `redbench-${version}.tgz`;
  const platform = `redbench-${target}-${version}.tgz`;
  const archive = `${rootName}.tar.gz`;
  return `# Redbench ${version} offline bundle\n\nThis bundle is for ${target}. Keep the separately supplied release-index.json and SHA256SUMS beside ${archive}.\n\n## Verify supplied release evidence\n\nawk '$2 == "${archive}" { print }' SHA256SUMS | sha256sum -c -\nnode -e 'const fs=require("fs"),crypto=require("crypto");const i=JSON.parse(fs.readFileSync("release-index.json"));const a="${archive}";const d=crypto.createHash("sha256").update(fs.readFileSync(a)).digest("hex");if(!i.artifacts.some(x=>x.name===a&&x.sha256===d))process.exit(1)'\n\n## Direct execution\n\n./bin/bench version\n./bin/bench commands --brief\n\n## Local npm installation\n\nBENCH_NO_REPAIR=1 npm install --offline --ignore-scripts --omit=optional --prefix ./prefix ./packages/${platform} ./packages/${wrapper}\nBENCH_NO_REPAIR=1 ./prefix/node_modules/redbench/bin/bench.sh version\nBENCH_NO_REPAIR=1 ./prefix/node_modules/redbench/bin/bench.sh commands --brief\n\n## Internal registry\n\nnpm publish ./packages/${platform} --registry http://127.0.0.1:4873\nnpm publish ./packages/${wrapper} --registry http://127.0.0.1:4873\nBENCH_NO_REPAIR=1 npm install --ignore-scripts --omit=optional --prefix ./prefix --registry http://127.0.0.1:4873 redbench@${version} @redbench/${target}@${version}\n\n## Removal\n\nnpm uninstall --offline --ignore-scripts --prefix ./prefix redbench @redbench/${target}\nrm -rf ./prefix ./npm-cache\n`;
}

for (const entry of entries) {
  const destination = path.join(archiveDir, entry.path);
  if (entry.kind === "binary") copyRegular(binary, destination, entry.mode);
  else if (entry.kind === "wrapper_tarball") copyRegular(path.join(npmDir, `redbench-${version}.tgz`), destination, entry.mode);
  else if (entry.kind === "platform_tarball") copyRegular(path.join(npmDir, `redbench-${target}-${version}.tgz`), destination, entry.mode);
  else if (entry.kind === "instructions") {
    fs.mkdirSync(path.dirname(destination), {recursive: true});
    fs.writeFileSync(destination, instructions(), {mode: Number.parseInt(entry.mode, 8)});
  } else if (entry.kind === "wrapper_manifest") copyRegular(path.join(wrapperExtract, "package", "component-manifest.json"), destination, entry.mode);
  else if (entry.kind === "platform_manifest") copyRegular(path.join(platformExtract, "package", "component-manifest.json"), destination, entry.mode);
  else if (entry.kind === "package_evidence") copyRegular(path.join(wrapperExtract, "package", entry.path.replace("evidence/", "")), destination, entry.mode);
  else if (entry.kind !== "archive_manifest") throw new Error(`unknown archive entry kind: ${entry.kind}`);
}

const manifestEntry = entries.find(entry => entry.kind === "archive_manifest");
if (!manifestEntry) throw new Error("release plan has no archive component manifest");
const archiveManifestPath = manifestEntry.path;
const files = entries.filter(entry => entry.kind !== "archive_manifest").map(entry => {
  const file = path.join(archiveDir, entry.path);
  const data = fs.readFileSync(file);
  return {path: entry.path, mode: Number.parseInt(entry.mode, 8).toString(8), size: data.length, sha256: crypto.createHash("sha256").update(data).digest("hex")};
}).sort((a, b) => Buffer.compare(Buffer.from(a.path), Buffer.from(b.path)));
const [os, arch] = target.split("-");
fs.mkdirSync(path.dirname(path.join(archiveDir, archiveManifestPath)), {recursive: true});
fs.writeFileSync(path.join(archiveDir, archiveManifestPath), JSON.stringify({schema_version: 1, component: {name: `redbench-offline-${target}`, version, target: {os, arch}}, files}) + "\n", {mode: 0o644});
