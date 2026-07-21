import fs from "node:fs";
import path from "node:path";
import crypto from "node:crypto";
import {readReleaseRequirements} from "./release-plan.mjs";

const [root, evidenceRoot, key, sourceCommit, packageVersion, requestedStatus, payloadFile] = process.argv.slice(2);
if (![root, evidenceRoot, key, sourceCommit, packageVersion, requestedStatus, payloadFile].every(Boolean)) {
  throw new Error("usage: write-producer-envelope.mjs <root> <evidence-root> <key> <source-commit> <package-version> <status|auto> <payload-json>");
}
const matches = readReleaseRequirements(root).records?.filter(record => record?.key === key) ?? [];
if (matches.length !== 1) throw new Error(`producer requirement ${key} is absent or ambiguous`);
const requirement = matches[0];
if (requirement.producer !== true || typeof requirement.owner !== "string" || !requirement.owner || typeof requirement.schema !== "string" || !requirement.schema) throw new Error(`producer requirement ${key} is invalid`);
const relative = requirement.path;
if (typeof relative !== "string" || !relative || relative.includes("\\") || path.posix.isAbsolute(relative) || relative.split("/").some(part => !part || part === "." || part === "..") || /[\x00-\x1f\x7f]/.test(relative)) throw new Error(`producer requirement ${key} path is unsafe`);
const resolvedRoot = path.resolve(evidenceRoot);
const output = path.resolve(resolvedRoot, ...relative.split("/"));
const outputRelative = path.relative(resolvedRoot, output);
if (!outputRelative || outputRelative.startsWith(`..${path.sep}`) || outputRelative === ".." || path.isAbsolute(outputRelative)) throw new Error(`producer requirement ${key} path escapes evidence root`);
const payload = JSON.parse(fs.readFileSync(payloadFile, "utf8"));
const attempts = Array.isArray(payload.operations) ? payload.operations.map(operation => operation?.observed_attempts) : [];
if (attempts.some(value => !Number.isSafeInteger(value) || value < 0)) throw new Error("producer payload has invalid observed attempt count");
const status = requestedStatus === "auto" ? attempts.length > 0 && attempts.every(value => value === 0) ? "satisfied" : "failed" : requestedStatus;
const canonicalPayload = JSON.stringify(sortValue(payload));
const digest = crypto.createHash("sha256").update(canonicalPayload).digest("hex");
const envelope = {
  schema_version: 1,
  key,
  owner: requirement.owner,
  schema: requirement.schema,
  identity: {source_commit: sourceCommit, package_version: packageVersion},
  status,
  reason: status === "satisfied" ? "" : "offline proof did not satisfy every sentinel and journey",
  payload,
  digest,
};
fs.mkdirSync(path.dirname(output), {recursive: true});
fs.writeFileSync(output, JSON.stringify(envelope, null, 2) + "\n", {mode: 0o644});

function sortValue(value) {
  if (Array.isArray(value)) return value.map(sortValue);
  if (value && typeof value === "object") return Object.fromEntries(Object.keys(value).sort().map(key => [key, sortValue(value[key])]));
  return value;
}
