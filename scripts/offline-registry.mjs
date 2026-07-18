import crypto from "node:crypto";
import fs from "node:fs";
import http from "node:http";
import path from "node:path";

const [store, portFile, requestFile] = process.argv.slice(2);
if (!store || !portFile || !requestFile) throw new Error("usage: offline-registry.mjs <store> <port-file> <request-log>");
const versionPattern = /^(?:redbench|redbench-(?:darwin|linux)-(?:arm64|x64))-(\d+\.\d+\.\d+(?:[+-][0-9A-Za-z.-]+)?)\.tgz$/;
const sha256 = data => crypto.createHash("sha256").update(data).digest("hex");

function packageForFile(file) {
  const match = versionPattern.exec(file);
  if (!match) return null;
  const name = file === `redbench-${match[1]}.tgz` ? "redbench" : `@redbench/${file.slice("redbench-".length, -(`-${match[1]}.tgz`).length)}`;
  return {file, name, version: match[1]};
}

function packages() {
  const out = new Map();
  for (const file of fs.readdirSync(store).sort((a, b) => Buffer.compare(Buffer.from(a), Buffer.from(b)))) {
    const item = packageForFile(file);
    if (!item) continue;
    const data = fs.readFileSync(path.join(store, file));
    out.set(item.name, {...item, data, sha256: sha256(data), integrity: `sha512-${crypto.createHash("sha512").update(data).digest("base64")}`});
  }
  return out;
}

function packageName(pathname) {
  const value = decodeURIComponent(pathname).replace(/^\//, "");
  return value.startsWith("@redbench/") ? value.split("/").slice(0, 2).join("/") : value.split("/")[0];
}

function writeLog(line) { fs.appendFileSync(requestFile, `${line}\n`); }

const server = http.createServer((request, response) => {
  const pathname = new URL(request.url, "http://127.0.0.1").pathname;
  if (request.method === "PUT" && pathname.startsWith("/upload/")) {
    const file = decodeURIComponent(pathname.slice("/upload/".length));
    if (!packageForFile(file)) { response.writeHead(400); response.end("invalid package filename\n"); return; }
    const chunks = [];
    request.on("data", chunk => chunks.push(chunk));
    request.on("end", () => {
      const data = Buffer.concat(chunks);
      if (data.length === 0) { response.writeHead(400); response.end("empty package\n"); return; }
      const temporary = path.join(store, `.${file}.${process.pid}`);
      fs.writeFileSync(temporary, data, {mode: 0o644});
      fs.renameSync(temporary, path.join(store, file));
      const item = packageForFile(file);
      writeLog(`PUT ${item.name}@${item.version} ${sha256(data)}`);
      response.writeHead(201); response.end("stored\n");
    });
    return;
  }
  if (request.method !== "GET") { writeLog(`${request.method} ${request.url}`); response.writeHead(405); response.end(); return; }
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
  const body = JSON.stringify({name, "dist-tags": {latest: item.version}, versions: {[item.version]: {name, version: item.version, dist: {tarball, integrity: item.integrity}}}});
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
