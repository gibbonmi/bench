import fs from "node:fs";
import path from "node:path";
import {pathToFileURL} from "node:url";

const byteOrder = (a, b) => Buffer.compare(Buffer.from(a), Buffer.from(b));

export function readReleasePlan(root) {
  const file = path.join(root, "scripts", "release-plan.json");
  const info = fs.lstatSync(file);
  if (!info.isFile() || info.isSymbolicLink() || info.size === 0) throw new Error("release plan is not a non-empty regular file");
  const plan = JSON.parse(fs.readFileSync(file, "utf8"));
  if (plan?.schema_version !== 1 || !Array.isArray(plan.targets) || plan.targets.length === 0 || !Array.isArray(plan.archive_entries)) throw new Error("release plan cardinality is invalid");
	const seen = new Set();
	for (const target of plan.targets) {
		if (!target || ![target.os, target.arch, target.goos, target.goarch, target.runner].every(value => typeof value === "string" && /^[0-9A-Za-z.-]+$/.test(value))) throw new Error("release plan target is invalid");
    const key = `${target.os}-${target.arch}`;
    if (seen.has(key)) throw new Error(`release plan repeats ${key}`);
    seen.add(key);
  }
  const entries = new Set();
  for (const entry of plan.archive_entries) {
    if (!entry || typeof entry.path !== "string" || !/^(0644|0755)$/.test(entry.mode) || typeof entry.kind !== "string" || entry.path.length === 0 || entry.path.includes("\\") || entry.path.startsWith("/") || entry.path.includes("..") || entries.has(entry.path)) throw new Error("release plan archive inventory is invalid");
    entries.add(entry.path);
  }
  return plan;
}

export function readReleaseRequirements(root) {
	const file = path.join(root, "internal", "releaseevidence", "requirements.json");
	const info = fs.lstatSync(file);
	if (!info.isFile() || info.isSymbolicLink() || info.size === 0) throw new Error("release evidence requirements are unsafe");
	const requirements = JSON.parse(fs.readFileSync(file, "utf8"));
	if (requirements?.schema_version !== 1 || !Array.isArray(requirements.records)) throw new Error("release evidence requirements are invalid");
	return requirements;
}

export function packagedEvidenceRecords(requirements) {
	return requirements.records.filter(record => Object.hasOwn(record, "package_mode"));
}

export function targetFor(plan, os, arch) {
  return plan.targets.find(target => target.os === os && target.arch === arch);
}

export function artifactNames(plan, version) {
	return artifactRecords(plan, version).map(item => item.name);
}

export function artifactRecords(plan, version) {
	return [{name: `redbench-${version}.tgz`, target: "wrapper", kind: "wrapper"}, ...plan.targets.flatMap(target => {
		const name = `${target.os}-${target.arch}`;
		return [
			{name: `redbench-${name}-${version}.tgz`, target: name, kind: "platform"},
			{name: `redbench-${version}-${name}.tar.gz`, target: name, kind: "archive"},
		];
	})].sort((left, right) => byteOrder(left.name, right.name));
}

export function archiveEntries(plan, target, version, packageEvidence) {
  const root = `redbench-${version}-${target}`;
  return plan.archive_entries.flatMap(entry => {
    if (entry.kind === "package_evidence") return packageEvidence.map(item => ({...entry, path: entry.path.replace("{package_evidence}", item.path), source_path: item.path, evidence_key: item.key}));
    return [{...entry, path: entry.path.replaceAll("{version}", version).replaceAll("{target}", target)}];
  }).map(entry => ({...entry, archive_path: `${root}/${entry.path}`})).sort((a, b) => byteOrder(a.path, b.path));
}

export function archiveEntryPath(plan, kind, target, version, packageEvidence) {
	const matches = archiveEntries(plan, target, version, packageEvidence).filter(entry => entry.kind === kind);
	if (matches.length !== 1) throw new Error(`release plan has ${matches.length} archive entries for ${kind}`);
	return matches[0].path;
}

export function archiveEvidencePath(plan, key, target, version, packageEvidence) {
	const matches = archiveEntries(plan, target, version, packageEvidence).filter(entry => entry.kind === "package_evidence" && entry.evidence_key === key);
	if (matches.length !== 1) throw new Error(`release plan has ${matches.length} packaged evidence entries for ${key}`);
	return matches[0].path;
}

export function archiveInventory(root, plan, target, version) {
	const packageEvidence = packagedEvidenceRecords(readReleaseRequirements(root));
	return Object.fromEntries(archiveEntries(plan, target, version, packageEvidence).map(entry => [entry.path, Number.parseInt(entry.mode, 8)]));
}

export function releaseEvidenceNames(root) {
	const names = ["internal/releaseevidence/requirements.json", "scripts/release-plan.json", ...packagedEvidenceRecords(readReleaseRequirements(root)).map(record => record.path)];
  if (names.some(name => typeof name !== "string" || name.length === 0) || new Set(names).size !== names.length) throw new Error("release-bound evidence inventory is invalid");
  return names.sort(byteOrder);
}

if (import.meta.url === pathToFileURL(process.argv[1]).href) {
  const [root, command, ...args] = process.argv.slice(2);
	if (!root || !command) throw new Error("usage: release-plan.mjs <root> <normalized-json|targets|matrix-json|artifact-names|artifact-records|artifact-name|archive-inventory|archive-entry-path|archive-evidence-path|evidence-names|target> [arguments]");
  const plan = readReleasePlan(root);
  if (command === "normalized-json") {
    process.stdout.write(JSON.stringify(plan) + "\n");
  } else if (command === "targets") {
    for (const target of plan.targets) process.stdout.write([target.os, target.arch, target.goos, target.goarch, target.runner].join("\t") + "\n");
  } else if (command === "matrix-json") {
    process.stdout.write(JSON.stringify({include: plan.targets}) + "\n");
	} else if (command === "artifact-names") {
    if (args.length !== 1) throw new Error("artifact-names requires version");
		process.stdout.write(artifactNames(plan, args[0]).join("\n") + "\n");
	} else if (command === "artifact-records") {
		if (args.length !== 1) throw new Error("artifact-records requires version");
		process.stdout.write(JSON.stringify(artifactRecords(plan, args[0])) + "\n");
	} else if (command === "artifact-name") {
		if (args.length !== 3) throw new Error("artifact-name requires version, target, and kind");
		const matches = artifactRecords(plan, args[0]).filter(item => item.target === args[1] && item.kind === args[2]);
		if (matches.length !== 1) throw new Error(`release plan has ${matches.length} artifacts for ${args[1]} ${args[2]}`);
		process.stdout.write(matches[0].name + "\n");
	} else if (command === "archive-inventory") {
		if (args.length !== 2) throw new Error("archive-inventory requires target and version");
		process.stdout.write(JSON.stringify(archiveInventory(root, plan, args[0], args[1])) + "\n");
	} else if (command === "archive-entry-path") {
		if (args.length !== 3) throw new Error("archive-entry-path requires kind, target, and version");
		process.stdout.write(archiveEntryPath(plan, args[0], args[1], args[2], packagedEvidenceRecords(readReleaseRequirements(root))) + "\n");
	} else if (command === "archive-evidence-path") {
		if (args.length !== 3) throw new Error("archive-evidence-path requires requirement key, target, and version");
		process.stdout.write(archiveEvidencePath(plan, args[0], args[1], args[2], packagedEvidenceRecords(readReleaseRequirements(root))) + "\n");
  } else if (command === "evidence-names") {
    process.stdout.write(releaseEvidenceNames(root).join("\n") + "\n");
  } else if (command === "target") {
    if (args.length !== 2) throw new Error("target requires os and arch");
    const target = targetFor(plan, args[0], args[1]);
    if (!target) process.exitCode = 1;
    else process.stdout.write([target.os, target.arch, target.goos, target.goarch, target.runner].join("\t") + "\n");
  } else {
    throw new Error(`unknown release-plan command: ${command}`);
  }
}
