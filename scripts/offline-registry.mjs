import crypto from "node:crypto";
import fs from "node:fs";
import http from "node:http";
import path from "node:path";
import {fileURLToPath} from "node:url";
import {readReleasePlan} from "./release-plan.mjs";

const [store, portFile, requestFile] = process.argv.slice(2);
if (!store || !portFile || !requestFile) throw new Error("usage: offline-registry.mjs <store> <port-file> <request-log>");
const versionPattern = /^(\d+\.\d+\.\d+(?:[+-][0-9A-Za-z.-]+)?)$/;
const tagPattern = /^[a-z0-9]([a-z0-9.-]{0,63})$/i;
const scriptRoot = path.dirname(path.dirname(fileURLToPath(import.meta.url)));
const plannedTargets = readReleasePlan(scriptRoot).targets.map(target => `${target.os}-${target.arch}`);
const sha256 = data => crypto.createHash("sha256").update(data).digest("hex");
const integrityOf = data => `sha512-${crypto.createHash("sha512").update(data).digest("base64")}`;

// In-memory registry state that outlives a single request: dist-tag assignments,
// deprecation messages, and pending (not-yet-approved) staged submissions. This
// fixture models exactly the state a real registry keeps server-side; it is not a
// second package builder and makes no claim about public npm's actual behavior.
const tagState = new Map(); // name -> Map<tag, version>
const deprecations = new Map(); // `${name}@${version}` -> message
const staged = new Map(); // stageId -> {file, data}
let stageCounter = 0;

function packageForFile(file) {
	if (!file.startsWith("redbench-") || !file.endsWith(".tgz")) return null;
	const body = file.slice("redbench-".length, -".tgz".length);
	if (versionPattern.test(body)) return {file, name: "redbench", version: body};
	for (const target of plannedTargets) {
		const prefix = `${target}-`;
		if (!body.startsWith(prefix)) continue;
		const version = body.slice(prefix.length);
		if (versionPattern.test(version)) return {file, name: `@redbench/${target}`, version};
	}
	return null;
}

function packages() {
  const out = new Map();
  for (const file of fs.readdirSync(store).sort((a, b) => Buffer.compare(Buffer.from(a), Buffer.from(b)))) {
    const item = packageForFile(file);
    if (!item) continue;
    const data = fs.readFileSync(path.join(store, file));
    out.set(item.name, {...item, data, sha256: sha256(data), integrity: integrityOf(data)});
  }
  return out;
}

function packageName(pathname) {
  const value = decodeURIComponent(pathname).replace(/^\//, "");
  return value.startsWith("@redbench/") ? value.split("/").slice(0, 2).join("/") : value.split("/")[0];
}

function writeLog(line) { fs.appendFileSync(requestFile, `${line}\n`); }

function readBody(request) {
  return new Promise((resolve, reject) => {
    const chunks = [];
    request.on("data", chunk => chunks.push(chunk));
    request.on("end", () => resolve(Buffer.concat(chunks)));
    request.on("error", reject);
  });
}

function reject(response, log, status, message) {
  writeLog(log);
  response.writeHead(status, {"content-type": "application/json"});
  response.end(JSON.stringify({error: message}) + "\n");
}

function sendJSON(response, status, body) {
  const encoded = JSON.stringify(body);
  response.writeHead(status, {"content-type": "application/json", "content-length": Buffer.byteLength(encoded)});
  response.end(encoded);
}

// Route table for the `/-/...` control-plane paths (dist-tags, deprecate,
// integrity, stage, approve). Kept as one parse function so every accepted
// shape is enumerated once; anything else falls through to malformed rejection.
function parseControlPath(pathname) {
  const value = decodeURIComponent(pathname).replace(/^\//, "");
  let match;
  if ((match = /^-\/package\/(.+)\/dist-tags\/([^/]+)$/.exec(value))) {
    return {kind: "dist-tag", name: match[1], tag: match[2]};
  }
  if ((match = /^-\/package\/(.+)\/deprecate$/.exec(value))) {
    return {kind: "deprecate", name: match[1]};
  }
  if ((match = /^-\/integrity\/(.+)\/([^/]+)$/.exec(value))) {
    return {kind: "integrity", name: match[1], version: match[2]};
  }
  if ((match = /^-\/stage\/(.+)$/.exec(value))) {
    return {kind: "stage", file: match[1]};
  }
  if ((match = /^-\/approve\/([^/]+)$/.exec(value))) {
    return {kind: "approve", stageId: match[1]};
  }
  return null;
}

const server = http.createServer(async (request, response) => {
  const pathname = new URL(request.url, "http://127.0.0.1").pathname;

  if (request.method === "PUT" && pathname.startsWith("/upload/")) {
    const file = decodeURIComponent(pathname.slice("/upload/".length));
    if (!packageForFile(file)) return reject(response, `REJECT PUT ${request.url}`, 400, "invalid package filename");
    const data = await readBody(request);
    if (data.length === 0) return reject(response, `REJECT PUT ${request.url}`, 400, "empty package");
    const temporary = path.join(store, `.${file}.${process.pid}`);
    fs.writeFileSync(temporary, data, {mode: 0o644});
    fs.renameSync(temporary, path.join(store, file));
    const item = packageForFile(file);
    const tag = new URL(request.url, "http://127.0.0.1").searchParams.get("tag");
    if (tag) {
      if (!tagPattern.test(tag)) return reject(response, `REJECT PUT ${request.url}`, 400, "invalid dist-tag");
      if (!tagState.has(item.name)) tagState.set(item.name, new Map());
      tagState.get(item.name).set(tag, item.version);
    }
    writeLog(`PUT ${item.name}@${item.version} ${sha256(data)}${tag ? ` tag=${tag}` : ""}`);
    response.writeHead(201); response.end("stored\n");
    return;
  }

  const control = parseControlPath(pathname);

  if (control?.kind === "stage" && request.method === "PUT") {
    if (!packageForFile(control.file)) return reject(response, `REJECT STAGE ${request.url}`, 400, "invalid package filename");
    const data = await readBody(request);
    if (data.length === 0) return reject(response, `REJECT STAGE ${request.url}`, 400, "empty package");
    stageCounter += 1;
    const stageId = `stage-${stageCounter}-${sha256(data).slice(0, 12)}`;
    staged.set(stageId, {file: control.file, data});
    writeLog(`STAGE ${control.file} ${stageId} ${sha256(data)}`);
    sendJSON(response, 201, {stage_id: stageId});
    return;
  }

  if (control?.kind === "approve" && request.method === "POST") {
    const entry = staged.get(control.stageId);
    if (!entry) return reject(response, `REJECT APPROVE ${request.url}`, 404, "unknown stage id");
    staged.delete(control.stageId);
    const temporary = path.join(store, `.${entry.file}.${process.pid}`);
    fs.writeFileSync(temporary, entry.data, {mode: 0o644});
    fs.renameSync(temporary, path.join(store, entry.file));
    const item = packageForFile(entry.file);
    writeLog(`APPROVE ${item.name}@${item.version} ${control.stageId} ${sha256(entry.data)}`);
    response.writeHead(200); response.end("approved\n");
    return;
  }

  if (control?.kind === "dist-tag" && (request.method === "PUT" || request.method === "DELETE")) {
    if (!tagPattern.test(control.tag)) return reject(response, `REJECT DIST-TAG ${request.url}`, 400, "invalid dist-tag");
    const item = packages().get(control.name);
    if (request.method === "DELETE") {
      const had = tagState.get(control.name)?.delete(control.tag);
      writeLog(`DIST-TAG-RM ${control.name} ${control.tag}`);
      if (!had) return reject(response, `REJECT DIST-TAG-RM ${request.url}`, 404, "tag not set");
      response.writeHead(200); response.end("removed\n");
      return;
    }
    const body = await readBody(request);
    let version;
    try { version = JSON.parse(body.toString("utf8")); } catch { version = body.toString("utf8").trim(); }
    if (!item || item.version !== version) return reject(response, `REJECT DIST-TAG-ADD ${request.url}`, 404, "package version is not live");
    if (!tagState.has(control.name)) tagState.set(control.name, new Map());
    tagState.get(control.name).set(control.tag, version);
    writeLog(`DIST-TAG-ADD ${control.name} ${control.tag}=${version}`);
    response.writeHead(200); response.end("tagged\n");
    return;
  }

  if (control?.kind === "deprecate" && request.method === "POST") {
    const body = await readBody(request);
    let payload;
    try { payload = JSON.parse(body.toString("utf8")); } catch { payload = null; }
    const item = packages().get(control.name);
    if (!payload || typeof payload.version !== "string" || typeof payload.message !== "string" || !item || item.version !== payload.version) {
      return reject(response, `REJECT DEPRECATE ${request.url}`, 400, "invalid deprecation request");
    }
    deprecations.set(`${control.name}@${payload.version}`, payload.message);
    writeLog(`DEPRECATE ${control.name}@${payload.version}`);
    response.writeHead(200); response.end("deprecated\n");
    return;
  }

  if (control?.kind === "integrity" && request.method === "GET") {
    const item = packages().get(control.name);
    writeLog(`GET-INTEGRITY ${control.name}@${control.version}`);
    if (!item || item.version !== control.version) { response.writeHead(404); response.end("not found\n"); return; }
    sendJSON(response, 200, {name: control.name, version: item.version, integrity: item.integrity});
    return;
  }

  if (request.method === "DELETE") {
    // No unpublish endpoint exists. Any DELETE that is not a recognized
    // dist-tag removal is an unpublish-shaped request: log and reject it.
    return reject(response, `REJECT UNPUBLISH ${request.url}`, 405, "unpublish is not supported");
  }

  if (request.method !== "GET") return reject(response, `REJECT ${request.method} ${request.url}`, 400, "malformed request");

  const name = packageName(pathname);
  const item = packages().get(name);
  writeLog(`GET ${request.url}`);
  if (!item) { response.writeHead(404); response.end("not found\n"); return; }
  if (pathname.endsWith(".tgz")) {
    response.writeHead(200, {"content-type": "application/octet-stream", "content-length": item.data.length});
    response.end(item.data);
    return;
  }
  const encoded = name.startsWith("@") ? name.replace("/", "%2f") : name;
  const tarball = `http://${request.headers.host}/${encoded}/-/${item.file}`;
  const distTags = {latest: item.version, ...Object.fromEntries(tagState.get(name) ?? [])};
  const deprecated = deprecations.get(`${name}@${item.version}`);
  const versionInfo = {name, version: item.version, dist: {tarball, integrity: item.integrity}};
  if (deprecated) versionInfo.deprecated = deprecated;
  const body = JSON.stringify({name, "dist-tags": distTags, versions: {[item.version]: versionInfo}});
  response.writeHead(200, {"content-type": "application/json", "content-length": Buffer.byteLength(body)});
  response.end(body);
});

server.on("error", error => { fs.writeFileSync(`${portFile}.error`, `${error.code || "listen-error"}\n`); process.exit(1); });
server.listen(0, "127.0.0.1", () => {
  const address = server.address();
  fs.writeFileSync(`${portFile}.tmp-${process.pid}`, `${address.port}\n`, {mode: 0o644});
  fs.renameSync(`${portFile}.tmp-${process.pid}`, portFile);
});
function stop() { server.close(() => process.exit(0)); }
process.on("SIGTERM", stop);
process.on("SIGINT", stop);
