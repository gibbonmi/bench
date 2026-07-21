import fs from "node:fs";
import path from "node:path";
import crypto from "node:crypto";
import {readReleasePlan, readReleaseRequirements} from "./release-plan.mjs";

const [root, wrapperDir, artifactsDir, version] = process.argv.slice(2);
if (!root || !wrapperDir || !artifactsDir || !version) throw new Error("usage: build-binary-pins.mjs <root> <wrapper-dir> <artifact-dir> <version>");
const schema = readReleaseRequirements(root).binary_pin_manifest;
if (!schema || schema.schema_version !== 1 || typeof schema.path !== "string" || !schema.path || schema.path.includes("\\") || path.posix.isAbsolute(schema.path) || schema.path.split("/").some(part => !part || part === "." || part === "..") || /[\x00-\x1f\x7f]/.test(schema.path)) throw new Error("binary pin manifest requirement is invalid");
const pins = readReleasePlan(root).targets.map(target => {
  const file = path.join(artifactsDir, `redbench-${target.os}-${target.arch}-${version}.tgz`);
  const stat = fs.lstatSync(file);
  if (!stat.isFile() || stat.isSymbolicLink() || stat.size === 0) throw new Error(`platform tarball is missing or unsafe: ${file}`);
  const integrity = "sha512-" + crypto.createHash("sha512").update(fs.readFileSync(file)).digest("base64");
  return {name: `@redbench/${target.os}-${target.arch}`, version, integrity};
});
pins.sort((a, b) => a.name < b.name ? -1 : a.name > b.name ? 1 : 0);
const output = path.join(wrapperDir, schema.path);
fs.mkdirSync(path.dirname(output), {recursive: true});
fs.writeFileSync(output, JSON.stringify({schema_version: schema.schema_version, pins}, null, 2) + "\n", {mode: 0o644});
