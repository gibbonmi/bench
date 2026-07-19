const net = require("node:net");
const fs = require("node:fs");

const log = process.env.BENCH_OFFLINE_EGRESS_LOG;
const allowed = process.env.BENCH_OFFLINE_ALLOWED_ORIGIN;
if (!log) throw new Error("offline network sentinel requires BENCH_OFFLINE_EGRESS_LOG");

function endpoint(args) {
  const first = args[0];
  if (typeof first === "object" && first !== null) {
    if (first.path || !first.port) return null;
    return {host: first.host || first.hostname, port: String(first.port)};
  }
  if (typeof first === "number") return {host: args[1], port: String(first)};
  return null;
}

const connect = net.Socket.prototype.connect;
net.Socket.prototype.connect = function (...args) {
  const target = endpoint(args);
  if (target) {
    const origin = `${target.host || "localhost"}:${target.port}`;
    if (origin !== allowed) {
      fs.appendFileSync(log, `${origin}\n`);
      const error = new Error(`offline smoke: denied undeclared egress to ${origin}`);
      error.code = "BENCH_OFFLINE_EGRESS";
      process.nextTick(() => this.emit("error", error));
      return this;
    }
  }
  return connect.apply(this, args);
};

function guardedConnect(...args) {
  return new net.Socket().connect(...args);
}
net.connect = guardedConnect;
net.createConnection = guardedConnect;
