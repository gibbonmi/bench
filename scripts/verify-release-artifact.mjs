import crypto from "node:crypto";
import fs from "node:fs";
import path from "node:path";

const [indexPath, sumsPath, artifactPath] = process.argv.slice(2);
if (![indexPath, sumsPath, artifactPath].every(Boolean)) throw new Error("usage: verify-release-artifact.mjs <release-index> <SHA256SUMS> <artifact>");
const name = path.basename(artifactPath);
const data = fs.readFileSync(artifactPath);
const digest = crypto.createHash("sha256").update(data).digest("hex");
const index = JSON.parse(fs.readFileSync(indexPath, "utf8"));
const sums = new Map(fs.readFileSync(sumsPath, "utf8").trimEnd().split("\n").map(line => {
  const match = /^([0-9a-f]{64})  ([^/]+)$/.exec(line);
  if (!match) throw new Error("checksum row is malformed");
  return [match[2], match[1]];
}));
if (!Array.isArray(index.artifacts) || !index.artifacts.some(item => item.name === name && item.sha256 === digest) || sums.get(name) !== digest) throw new Error("offline smoke: supplied release evidence does not bind artifact bytes");
