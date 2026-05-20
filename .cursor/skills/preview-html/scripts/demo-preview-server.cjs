#!/usr/bin/env node
/**
 * preview-html — multi-root daemon
 *
 * One daemon serves many --root directories (different ports).
 * - Per root: 5 min without HTTP requests → stop that root
 * - Any request resets that root's 5 min timer
 * - All roots stopped → 5 min later daemon exits
 *
 * CLI:
 *   node demo-preview-server.cjs --root <dir>
 *   node demo-preview-server.cjs --stop --root <dir>
 *   node demo-preview-server.cjs --stop-all
 *   node demo-preview-server.cjs --status --root <dir>
 *   node demo-preview-server.cjs --list
 *
 * Internal: node demo-preview-server.cjs --daemon
 */

const { spawn } = require('child_process');
const crypto = require('crypto');
const http = require('http');
const fs = require('fs');
const path = require('path');
const os = require('os');

const DAEMON_FILE = path.join(os.tmpdir(), 'preview-html-daemon.json');
const STATE_PREFIX = 'preview-html-';
const STATE_SUFFIX = '.state.json';
const LEGACY_STATE_FILE = path.join(os.tmpdir(), 'preview-html-server.state.json');

const ROOT_IDLE_MS = 5 * 60 * 1000;
const DAEMON_EMPTY_MS = 5 * 60 * 1000;
const TICK_MS = 15 * 1000;
const STOP_WAIT_MS = 3000;
const HEALTH_TIMEOUT_MS = 2000;
const DAEMON_BOOT_TIMEOUT_MS = 15000;

const MIME_TYPES = {
  '.html': 'text/html; charset=utf-8',
  '.htm': 'text/html; charset=utf-8',
  '.css': 'text/css; charset=utf-8',
  '.js': 'application/javascript; charset=utf-8',
  '.mjs': 'application/javascript; charset=utf-8',
  '.json': 'application/json; charset=utf-8',
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.jpeg': 'image/jpeg',
  '.gif': 'image/gif',
  '.svg': 'image/svg+xml',
  '.ico': 'image/x-icon',
  '.webp': 'image/webp',
  '.woff': 'font/woff',
  '.woff2': 'font/woff2',
  '.ttf': 'font/ttf',
  '.map': 'application/json',
};

// ---------- utils ----------

function sleep(ms) {
  return new Promise((r) => setTimeout(r, ms));
}

function emitJson(payload) {
  process.stdout.write(JSON.stringify(payload) + '\n');
}

function canonicalRoot(rootDir) {
  const resolved = path.resolve(rootDir);
  try {
    return fs.realpathSync(resolved);
  } catch {
    return resolved;
  }
}

function rootKey(canonical) {
  return crypto.createHash('sha256').update(canonical).digest('hex').slice(0, 16);
}

function stateFileForRoot(canonical) {
  return path.join(os.tmpdir(), `${STATE_PREFIX}${rootKey(canonical)}${STATE_SUFFIX}`);
}

function readJsonFile(file) {
  try {
    return JSON.parse(fs.readFileSync(file, 'utf8'));
  } catch {
    return null;
  }
}

function writeJsonFile(file, data) {
  fs.writeFileSync(file, JSON.stringify(data), 'utf8');
}

function clearFile(file) {
  try {
    fs.unlinkSync(file);
  } catch {
    /* ignore */
  }
}

function isProcessAlive(pid) {
  if (!pid || pid <= 0) return false;
  try {
    process.kill(pid, 0);
    return true;
  } catch (e) {
    return e.code === 'EPERM';
  }
}

async function stopProcess(pid) {
  if (!isProcessAlive(pid)) return;
  try {
    process.kill(pid, 'SIGTERM');
  } catch {
    return;
  }
  const deadline = Date.now() + STOP_WAIT_MS;
  while (Date.now() < deadline) {
    if (!isProcessAlive(pid)) return;
    await sleep(100);
  }
  try {
    process.kill(pid, 'SIGKILL');
  } catch {
    /* ignore */
  }
}

function pickPort(explicit) {
  if (explicit) return Number(explicit);
  return 49152 + Math.floor(Math.random() * 16383);
}

function makePayload(entry, reused = false) {
  return {
    type: 'preview-html-started',
    reused: Boolean(reused),
    pid: entry.daemonPid,
    port: entry.port,
    host: entry.host,
    url: entry.url,
    root: entry.canonical,
  };
}

function checkHealth(url) {
  return new Promise((resolve) => {
    const req = http.get(url, { timeout: HEALTH_TIMEOUT_MS }, (res) => {
      res.resume();
      resolve(res.statusCode >= 200 && res.statusCode < 500);
    });
    req.on('timeout', () => {
      req.destroy();
      resolve(false);
    });
    req.on('error', () => resolve(false));
  });
}

function readBody(req) {
  return new Promise((resolve, reject) => {
    let data = '';
    req.on('data', (c) => { data += c; });
    req.on('end', () => {
      if (!data) return resolve({});
      try {
        resolve(JSON.parse(data));
      } catch (e) {
        reject(e);
      }
    });
    req.on('error', reject);
  });
}

function sendJson(res, status, obj) {
  const body = JSON.stringify(obj);
  res.writeHead(status, { 'Content-Type': 'application/json; charset=utf-8' });
  res.end(body);
}

// ---------- static file serving ----------

function resolveFile(root, urlPath) {
  const raw = decodeURIComponent((urlPath || '/').split('?')[0]);
  let rel = raw.replace(/^\/+/, '');
  if (!rel || rel === '/') return { type: 'index' };
  if (rel.includes('..') || rel.includes('\\')) return null;
  const abs = path.resolve(root, rel);
  const rootWithSep = root.endsWith(path.sep) ? root : root + path.sep;
  if (abs !== root && !abs.startsWith(rootWithSep)) return null;
  return { type: 'file', abs, rel };
}

function listHtmlFiles(root) {
  try {
    return fs.readdirSync(root)
      .filter((f) => f.endsWith('.html') || f.endsWith('.htm'))
      .map((f) => {
        const fp = path.join(root, f);
        return { name: f, mtime: fs.statSync(fp).mtime.getTime() };
      })
      .sort((a, b) => b.mtime - a.mtime);
  } catch {
    return [];
  }
}

function buildIndexPage(root, files) {
  const items = files.length
    ? files.map((f) => `<li><a href="/${encodeURIComponent(f.name)}">${f.name}</a> <span class="muted">(${new Date(f.mtime).toLocaleString()})</span></li>`).join('\n')
    : '<li class="muted">No .html files in this directory yet.</li>';
  return `<!DOCTYPE html>
<html lang="zh-CN">
<head><meta charset="utf-8"><title>preview-html</title>
<style>body{font-family:system-ui,sans-serif;padding:2rem;max-width:720px;margin:0 auto}h1{font-size:1.25rem}.muted{color:#666;font-size:.9rem}code{background:#f4f4f5;padding:.1em .35em;border-radius:4px}</style>
</head><body><h1>preview-html</h1><p class="muted">Serving: <code>${root.replace(/</g, '')}</code></p><ul>${items}</ul></body></html>`;
}

function serveFile(root, resolved, res) {
  let target = resolved.abs;
  if (fs.existsSync(target) && fs.statSync(target).isDirectory()) {
    const index = path.join(target, 'index.html');
    if (fs.existsSync(index)) target = index;
    else {
      res.writeHead(403, { 'Content-Type': 'text/plain; charset=utf-8' });
      res.end('Directory listing only at /');
      return;
    }
  }
  if (!fs.existsSync(target) || !fs.statSync(target).isFile()) {
    res.writeHead(404, { 'Content-Type': 'text/plain; charset=utf-8' });
    res.end('Not found');
    return;
  }
  const ext = path.extname(target).toLowerCase();
  res.writeHead(200, { 'Content-Type': MIME_TYPES[ext] || 'application/octet-stream' });
  fs.createReadStream(target).pipe(res);
}

function createContentServer(rootDir, onActivity) {
  return http.createServer((req, res) => {
    onActivity();
    if (req.method !== 'GET' && req.method !== 'HEAD') {
      res.writeHead(405, { 'Content-Type': 'text/plain; charset=utf-8' });
      res.end('Method not allowed');
      return;
    }
    const resolved = resolveFile(rootDir, req.url);
    if (!resolved) {
      res.writeHead(403, { 'Content-Type': 'text/plain; charset=utf-8' });
      res.end('Forbidden');
      return;
    }
    if (resolved.type === 'index') {
      const html = buildIndexPage(rootDir, listHtmlFiles(rootDir));
      res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8' });
      if (req.method === 'HEAD') res.end();
      else res.end(html);
      return;
    }
    if (req.method === 'HEAD') {
      const st = fs.existsSync(resolved.abs) ? fs.statSync(resolved.abs) : null;
      if (!st || !st.isFile()) {
        res.writeHead(404);
        res.end();
        return;
      }
      const ext = path.extname(resolved.abs).toLowerCase();
      res.writeHead(200, { 'Content-Type': MIME_TYPES[ext] || 'application/octet-stream' });
      res.end();
      return;
    }
    serveFile(rootDir, resolved, res);
  });
}

// ---------- daemon ----------

function runDaemon() {
  const bindHost = process.env.PREVIEW_HTML_HOST || '127.0.0.1';
  const adminPort = pickPort(process.env.PREVIEW_HTML_ADMIN_PORT);
  const roots = new Map();
  let emptySince = null;

  function persistRootState(entry) {
    writeJsonFile(entry.stateFile, {
      pid: process.pid,
      port: entry.port,
      host: entry.host,
      url: entry.url,
      root: entry.canonical,
      daemon: true,
      startedAt: entry.startedAt,
      lastActivity: entry.lastActivity,
    });
  }

  function removeRoot(key, reason) {
    const entry = roots.get(key);
    if (!entry) return;
    entry.server.close();
    roots.delete(key);
    clearFile(entry.stateFile);
    process.stderr.write(`preview-html: stopped root (${reason}) ${entry.canonical}\n`);
  }

  function touchRoot(key) {
    const entry = roots.get(key);
    if (entry) {
      entry.lastActivity = Date.now();
      persistRootState(entry);
    }
  }

  function listEntries() {
    return Array.from(roots.values()).map((e) => ({
      pid: process.pid,
      port: e.port,
      host: e.host,
      url: e.url,
      root: e.canonical,
      lastActivity: e.lastActivity,
      startedAt: e.startedAt,
    }));
  }

  function ensureRoot(rootPath, opts = {}) {
    const canonical = canonicalRoot(rootPath);
    const key = rootKey(canonical);
    const existing = roots.get(key);
    if (existing) {
      touchRoot(key);
      emptySince = null;
      return makePayload({ ...existing, daemonPid: process.pid }, true);
    }

    const port = pickPort(opts.port);
    const host = opts.host || bindHost;
    const urlHost = host === '127.0.0.1' ? 'localhost' : host;
    const entry = {
      key,
      canonical,
      rootDir: rootPath,
      port,
      host,
      url: `http://${urlHost}:${port}/`,
      stateFile: stateFileForRoot(canonical),
      startedAt: new Date().toISOString(),
      lastActivity: Date.now(),
      server: null,
      daemonPid: process.pid,
    };

    entry.server = createContentServer(rootPath, () => touchRoot(key));
    entry.server.listen(port, host);
    roots.set(key, entry);
    emptySince = null;
    persistRootState(entry);
    process.stderr.write(`preview-html: ${entry.url} root=${canonical}\n`);
    return makePayload(entry, false);
  }

  const admin = http.createServer(async (req, res) => {
    try {
      if (req.method === 'GET' && req.url === '/health') {
        return sendJson(res, 200, { ok: true, roots: roots.size });
      }
      if (req.method === 'GET' && req.url === '/list') {
        return sendJson(res, 200, { type: 'preview-html-list', count: roots.size, instances: listEntries() });
      }
      if (req.method === 'POST' && req.url === '/ensure-root') {
        const body = await readBody(req);
        if (!body.root) return sendJson(res, 400, { error: 'root required' });
        const resolved = resolveRootPath(body.root);
        const payload = ensureRoot(resolved, { port: body.port, host: body.host });
        return sendJson(res, 200, payload);
      }
      if (req.method === 'POST' && req.url === '/stop-root') {
        const body = await readBody(req);
        if (!body.root) return sendJson(res, 400, { error: 'root required' });
        const key = rootKey(canonicalRoot(resolveRootPath(body.root)));
        const had = roots.has(key);
        removeRoot(key, 'manual');
        return sendJson(res, 200, { type: 'preview-html-stopped', root: canonicalRoot(resolveRootPath(body.root)), wasRunning: had });
      }
      if (req.method === 'POST' && req.url === '/stop-all') {
        const stopped = listEntries().map((e) => e.root);
        for (const key of Array.from(roots.keys())) removeRoot(key, 'stop-all');
        clearFile(DAEMON_FILE);
        admin.close(() => process.exit(0));
        return sendJson(res, 200, { type: 'preview-html-stopped-all', count: stopped.length, roots: stopped });
      }
      if (req.method === 'GET' && req.url.startsWith('/status?')) {
        const u = new URL(req.url, 'http://127.0.0.1');
        const rootParam = u.searchParams.get('root');
        if (!rootParam) return sendJson(res, 400, { error: 'root query required' });
        const key = rootKey(canonicalRoot(resolveRootPath(rootParam)));
        const entry = roots.get(key);
        if (!entry) return sendJson(res, 404, { error: 'not running' });
        return sendJson(res, 200, makePayload({ ...entry, daemonPid: process.pid }, true));
      }
      sendJson(res, 404, { error: 'not found' });
    } catch (err) {
      sendJson(res, 500, { error: err.message });
    }
  });

  function resolveRootPath(rootArg) {
    const root = path.resolve(rootArg);
    if (!fs.existsSync(root) || !fs.statSync(root).isDirectory()) {
      throw new Error(`invalid root: ${rootArg}`);
    }
    return root;
  }

  admin.listen(adminPort, bindHost, () => {
    writeJsonFile(DAEMON_FILE, { pid: process.pid, adminPort, host: bindHost });
    process.stderr.write(`preview-html daemon: pid=${process.pid} admin=${adminPort}\n`);
  });

  const tick = setInterval(() => {
    const now = Date.now();
    for (const [key, entry] of roots) {
      if (now - entry.lastActivity >= ROOT_IDLE_MS) {
        removeRoot(key, 'idle-5m');
      }
    }
    if (roots.size === 0) {
      if (!emptySince) {
        emptySince = now;
        process.stderr.write('preview-html daemon: no roots, will exit in 5m if still empty\n');
      } else if (now - emptySince >= DAEMON_EMPTY_MS) {
        clearInterval(tick);
        process.stderr.write('preview-html daemon: all roots closed 5m, exiting\n');
        clearFile(DAEMON_FILE);
        admin.close(() => process.exit(0));
      }
    } else {
      emptySince = null;
    }
  }, TICK_MS);

  function shutdown() {
    clearInterval(tick);
    for (const key of Array.from(roots.keys())) removeRoot(key, 'signal');
    clearFile(DAEMON_FILE);
    admin.close(() => process.exit(0));
  }

  process.on('SIGINT', shutdown);
  process.on('SIGTERM', shutdown);
}

// ---------- daemon client ----------

function readDaemon() {
  return readJsonFile(DAEMON_FILE);
}

async function checkAdmin(adminPort) {
  return new Promise((resolve) => {
    const req = http.get(`http://127.0.0.1:${adminPort}/health`, { timeout: HEALTH_TIMEOUT_MS }, (res) => {
      res.resume();
      resolve(res.statusCode === 200);
    });
    req.on('timeout', () => {
      req.destroy();
      resolve(false);
    });
    req.on('error', () => resolve(false));
  });
}

function apiCall(adminPort, method, apiPath, body) {
  return new Promise((resolve, reject) => {
    const payload = body ? JSON.stringify(body) : '';
    const req = http.request(
      {
        hostname: '127.0.0.1',
        port: adminPort,
        path: apiPath,
        method,
        headers: body
          ? { 'Content-Type': 'application/json', 'Content-Length': Buffer.byteLength(payload) }
          : {},
      },
      (res) => {
        let buf = '';
        res.on('data', (c) => { buf += c; });
        res.on('end', () => {
          try {
            resolve({ status: res.statusCode, data: JSON.parse(buf || '{}') });
          } catch {
            resolve({ status: res.statusCode, data: { raw: buf } });
          }
        });
      }
    );
    req.on('error', reject);
    if (body) req.write(payload);
    req.end();
  });
}

async function waitForDaemon() {
  const deadline = Date.now() + DAEMON_BOOT_TIMEOUT_MS;
  while (Date.now() < deadline) {
    const d = readDaemon();
    if (d && isProcessAlive(d.pid) && (await checkAdmin(d.adminPort))) return d;
    await sleep(200);
  }
  throw new Error('preview-html daemon failed to start');
}

async function ensureDaemon() {
  const d = readDaemon();
  if (d && isProcessAlive(d.pid) && (await checkAdmin(d.adminPort))) return d;

  if (d && isProcessAlive(d.pid)) await stopProcess(d.pid);
  clearFile(DAEMON_FILE);
  cleanupLegacyStateFiles();

  const child = spawn(process.execPath, [__filename, '--daemon'], {
    detached: true,
    stdio: 'ignore',
    windowsHide: true,
  });
  child.unref();
  return waitForDaemon();
}

function cleanupLegacyStateFiles() {
  if (fs.existsSync(LEGACY_STATE_FILE)) clearFile(LEGACY_STATE_FILE);
  try {
    for (const name of fs.readdirSync(os.tmpdir())) {
      if (name.startsWith(STATE_PREFIX) && name.endsWith(STATE_SUFFIX)) {
        const file = path.join(os.tmpdir(), name);
        const st = readJsonFile(file);
        if (st && st.pid && st.pid !== process.pid && !isProcessAlive(st.pid)) {
          clearFile(file);
        }
      }
    }
  } catch {
    /* ignore */
  }
}

// ---------- CLI ----------

function parseArgs(argv) {
  const opts = {
    root: process.env.PREVIEW_HTML_ROOT || '',
    port: process.env.PREVIEW_HTML_PORT || '',
    host: process.env.PREVIEW_HTML_HOST || '127.0.0.1',
    daemon: false,
    stop: false,
    stopAll: false,
    status: false,
    list: false,
    help: false,
  };
  for (let i = 2; i < argv.length; i++) {
    const a = argv[i];
    if (a === '--daemon') { opts.daemon = true; continue; }
    if (a === '--root' && argv[i + 1]) { opts.root = argv[++i]; continue; }
    if (a === '--port' && argv[i + 1]) { opts.port = argv[++i]; continue; }
    if (a === '--host' && argv[i + 1]) { opts.host = argv[++i]; continue; }
    if (a === '--stop') { opts.stop = true; continue; }
    if (a === '--stop-all') { opts.stopAll = true; continue; }
    if (a === '--status') { opts.status = true; continue; }
    if (a === '--list') { opts.list = true; continue; }
    if (a === '--help' || a === '-h') { opts.help = true; continue; }
  }
  return opts;
}

function usage() {
  process.stderr.write(
    'preview-html daemon\n\n' +
    '  node demo-preview-server.cjs --root <directory>\n' +
    '  node demo-preview-server.cjs --status --root <directory>\n' +
    '  node demo-preview-server.cjs --stop --root <directory>\n' +
    '  node demo-preview-server.cjs --stop-all\n' +
    '  node demo-preview-server.cjs --list\n\n' +
    'Idle: per-root 5m without requests → stop root; all roots gone 5m → exit daemon\n'
  );
}

function resolveRootArg(rootArg, required) {
  if (!rootArg) {
    if (required) {
      process.stderr.write('Error: --root or PREVIEW_HTML_ROOT is required.\n');
      usage();
      process.exit(1);
    }
    return null;
  }
  const root = path.resolve(rootArg);
  if (!fs.existsSync(root)) {
    process.stderr.write(`Error: root directory does not exist: ${root}\n`);
    process.exit(1);
  }
  if (!fs.statSync(root).isDirectory()) {
    process.stderr.write(`Error: root is not a directory: ${root}\n`);
    process.exit(1);
  }
  return root;
}

async function main() {
  const opts = parseArgs(process.argv);

  if (opts.daemon) {
    runDaemon();
    return;
  }

  if (opts.help) {
    usage();
    process.exit(0);
  }

  const daemon = await ensureDaemon();

  if (opts.stopAll) {
    const { data } = await apiCall(daemon.adminPort, 'POST', '/stop-all');
    emitJson(data);
    process.exit(0);
  }

  if (opts.list) {
    const { data } = await apiCall(daemon.adminPort, 'GET', '/list');
    emitJson(data);
    process.exit(0);
  }

  if (opts.stop) {
    const root = resolveRootArg(opts.root, true);
    const { data } = await apiCall(daemon.adminPort, 'POST', '/stop-root', { root });
    emitJson(data);
    process.exit(0);
  }

  if (opts.status) {
    const root = resolveRootArg(opts.root, true);
    const canonical = canonicalRoot(root);
    const { status, data } = await apiCall(
      daemon.adminPort,
      'GET',
      `/status?root=${encodeURIComponent(canonical)}`
    );
    if (status === 200) {
      emitJson(data);
      process.exit(0);
    }
    process.stderr.write(`No running server for root: ${canonical}\n`);
    process.exit(1);
  }

  const root = resolveRootArg(opts.root, true);
  const { status, data } = await apiCall(daemon.adminPort, 'POST', '/ensure-root', {
    root,
    port: opts.port || undefined,
    host: opts.host,
  });
  if (status !== 200) {
    process.stderr.write(`${data.error || 'ensure-root failed'}\n`);
    process.exit(1);
  }
  emitJson(data);
  if (data.reused) process.stderr.write(`preview-html: reusing ${data.url}\n`);
}

main().catch((err) => {
  process.stderr.write(`${err.message}\n`);
  process.exit(1);
});
