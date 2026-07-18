import fs from "node:fs";
import http from "node:http";
import crypto from "node:crypto";

const [store, portFile, requestFile] = process.argv.slice(2);
if (!store || !portFile || !requestFile) throw new Error("usage: offline-registry.mjs <store> <port-file> <request-log>");

const versionPattern = /^(?:redbench|redbench-(?:darwin|linux)-(?:arm64|x64))-(\d+\.\d+\.\d+(?:[+-][0-9A-Za-z.-]+)?)\.tgz$/;
const packages = new Map();
for (const file of fs.readdirSync(store).sort((a, b) => Buffer.compare(Buffer.from(a), Buffer.from(b)))) {
  const match = versionPattern.exec(file);
  if (!match) continue;
  const name = file === `redbench-${match[1]}.tgz` ? "redbench" : `@redbench/${file.slice("redbench-".length, -(`-${match[1]}.tgz`).length)}`;
  const data = fs.readFileSync(`${store}/${file}`);
  packages.set(name, {file, version: match[1], data, integrity: `sha512-${crypto.createHash("sha512").update(data).digest("base64")}`});
}
if (packages.size !== 2) throw new Error(`offline registry store has ${packages.size} packages, want platform and wrapper`);

function packageName(pathname) {
  const raw = decodeURIComponent(pathname);
  const value = raw.replace(/^\//, "");
  if (value.startsWith("@redbench/")) {
    const parts = value.split("/");
    return `${parts[0]}/${parts[1]}`;
  }
  return value.split("/")[0];
}

function metadata(name, host) {
  const packageInfo = packages.get(name);
  if (!packageInfo) return null;
  const encodedName = name.startsWith("@") ? name.replace("/", "%2f") : name;
  const tarball = `http://${host}/${encodedName}/-/${packageInfo.file}`;
  return {
    name,
    "dist-tags": {latest: packageInfo.version},
    versions: {
      [packageInfo.version]: {
        name,
        version: packageInfo.version,
        dist: {tarball, integrity: packageInfo.integrity},
      },
    },
  };
}

const server = http.createServer((request, response) => {
  fs.appendFileSync(requestFile, `${request.method} ${request.url}\n`);
  if (request.method !== "GET") {
    response.writeHead(405);
    response.end();
    return;
  }
  const pathname = new URL(request.url, "http://127.0.0.1").pathname;
  const name = packageName(pathname);
  const packageInfo = packages.get(name);
  if (!packageInfo) {
    response.writeHead(404);
    response.end("not found\n");
    return;
  }
  if (pathname.endsWith(".tgz")) {
    response.writeHead(200, {"content-type": "application/octet-stream", "content-length": packageInfo.data.length});
    response.end(packageInfo.data);
    return;
  }
  const body = JSON.stringify(metadata(name, request.headers.host));
  response.writeHead(200, {"content-type": "application/json", "content-length": Buffer.byteLength(body)});
  response.end(body);
});

server.on("error", (error) => {
  fs.writeFileSync(`${portFile}.error`, `${error.code || "listen-error"}\n`);
  process.exit(1);
});

server.listen(0, "127.0.0.1", () => {
  const address = server.address();
  fs.writeFileSync(`${portFile}.tmp-${process.pid}`, `${address.port}\n`, {mode: 0o644});
  fs.renameSync(`${portFile}.tmp-${process.pid}`, portFile);
});

function stop() {
  server.close(() => process.exit(0));
}
process.on("SIGTERM", stop);
process.on("SIGINT", stop);
