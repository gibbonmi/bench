import fs from "node:fs";
import path from "node:path";
import zlib from "node:zlib";

const [source, destination] = process.argv.slice(2);
if (!source || !destination) throw new Error("usage: write-deterministic-archive.mjs <directory> <archive>");

const byteOrder = (a, b) => Buffer.compare(Buffer.from(a), Buffer.from(b));
const control = value => [...Buffer.from(value)].some(byte => byte < 0x20 || byte === 0x7f);
const entries = [];

function walk(current, relative) {
  const info = fs.lstatSync(current);
  if (info.isSymbolicLink() || (!info.isDirectory() && !info.isFile())) throw new Error(`archive input is not a regular file or directory: ${current}`);
  const name = relative || path.basename(current);
  if (control(name) || name.includes("\\") || path.posix.isAbsolute(name)) throw new Error(`archive path is unsafe: ${name}`);
  if (info.isDirectory()) {
    entries.push({name: `${name}/`, directory: true, mode: 0o755});
    for (const child of fs.readdirSync(current).sort(byteOrder)) walk(path.join(current, child), `${name}/${child}`);
    return;
  }
  if (info.size === 0) throw new Error(`archive input is empty: ${current}`);
  const mode = info.mode & 0o777;
  if (mode !== 0o644 && mode !== 0o755) throw new Error(`archive input has unsafe mode ${mode.toString(8)}: ${current}`);
  entries.push({name, directory: false, mode, data: fs.readFileSync(current)});
}

walk(source, "");
entries.sort((a, b) => byteOrder(a.name, b.name));

function field(value, size, label) {
  const bytes = Buffer.from(value);
  if (bytes.length > size) throw new Error(`${label} is too long`);
  const out = Buffer.alloc(size);
  bytes.copy(out);
  return out;
}

function octal(value, size) {
  const text = value.toString(8).padStart(size - 1, "0");
  if (text.length > size - 1) throw new Error("tar numeric field is too large");
  return field(`${text}\0`, size, "tar numeric field");
}

function header(entry) {
  const out = Buffer.alloc(512);
  let name = entry.name;
  let prefix = "";
  if (Buffer.byteLength(name) > 100) {
    const slash = name.lastIndexOf("/", 155);
    if (slash <= 0 || Buffer.byteLength(name.slice(0, slash)) > 155 || Buffer.byteLength(name.slice(slash + 1)) > 100) throw new Error(`archive path is too long: ${name}`);
    prefix = name.slice(0, slash);
    name = name.slice(slash + 1);
  }
  field(name, 100, "tar path").copy(out, 0);
  octal(entry.mode, 8).copy(out, 100);
  octal(0, 8).copy(out, 108);
  octal(0, 8).copy(out, 116);
  octal(entry.directory ? 0 : entry.data.length, 12).copy(out, 124);
  octal(0, 12).copy(out, 136);
  out.fill(0x20, 148, 156);
  out[156] = entry.directory ? 0x35 : 0x30;
  field("ustar\0", 6, "tar format").copy(out, 257);
  field("00", 2, "tar version").copy(out, 263);
  field(prefix, 155, "tar prefix").copy(out, 345);
  let sum = 0;
  for (const byte of out) sum += byte;
  octal(sum, 8).copy(out, 148);
  return out;
}

const blocks = [];
for (const entry of entries) {
  blocks.push(header(entry));
  if (!entry.directory) {
    blocks.push(entry.data);
    const padding = (512 - (entry.data.length % 512)) % 512;
    if (padding) blocks.push(Buffer.alloc(padding));
  }
}
blocks.push(Buffer.alloc(1024));
const compressed = zlib.gzipSync(Buffer.concat(blocks), {level: 9, mtime: 0});
fs.mkdirSync(path.dirname(destination), {recursive: true});
const temporary = `${destination}.tmp-${process.pid}`;
fs.writeFileSync(temporary, compressed, {mode: 0o644});
fs.renameSync(temporary, destination);
