'use strict';

const fs = require('fs');
const path = require('path');
const os = require('os');
const http = require('http');

const DEFAULT_HOST = '127.0.0.1';
const KEEPALIVE_STALE_MS = 2 * 60 * 1000;
const INSTANCE_KEY_RE = /^(.+)\+([0-9a-f]{12})\.json$/i;

function getRegistryRoot(override) {
  if (override) return path.resolve(override);
  return path.join(os.homedir(), '.harness-services');
}

function sanitizeBranchFileName(name) {
  return String(name)
    .replace(/\//g, '-')
    .replace(/\\/g, '-')
    .replace(/:/g, '-')
    .replace(/\*/g, '-')
    .replace(/\?/g, '-')
    .replace(/"/g, '-')
    .replace(/</g, '-')
    .replace(/>/g, '-')
    .replace(/\|/g, '-');
}

function parseArgs(argv) {
  const args = {};
  for (let i = 2; i < argv.length; i++) {
    const token = argv[i];
    if (!token.startsWith('--')) continue;
    const key = token.slice(2);
    const next = argv[i + 1];
    if (!next || next.startsWith('--')) {
      args[key] = true;
      continue;
    }
    args[key] = next;
    i++;
  }
  return args;
}

function readJsonFile(filePath) {
  const raw = fs.readFileSync(filePath, 'utf8');
  return JSON.parse(raw);
}

function parseRegistryFileName(fileName) {
  const match = fileName.match(INSTANCE_KEY_RE);
  if (!match) return null;
  return {
    branchName: match[1],
    instanceKey: match[2],
  };
}

function listRegistryFiles(registryRoot, appName, branchName) {
  const appDir = path.join(registryRoot, appName);
  if (!fs.existsSync(appDir)) {
    return [];
  }

  const prefix = branchName ? `${sanitizeBranchFileName(branchName)}+` : null;
  return fs
    .readdirSync(appDir)
    .filter((name) => name.endsWith('.json') && !name.endsWith('.lock'))
    .filter((name) => (prefix ? name.startsWith(prefix) : true))
    .map((name) => path.join(appDir, name))
    .sort();
}

function isKeepaliveStale(lastKeepalive) {
  if (!lastKeepalive) return true;
  const ts = Date.parse(lastKeepalive);
  if (Number.isNaN(ts)) return true;
  return Date.now() - ts > KEEPALIVE_STALE_MS;
}

function enrichRegistry(registry) {
  if (!registry || !Array.isArray(registry.instances)) {
    return registry;
  }
  return {
    ...registry,
    instances: registry.instances.map((inst) => ({
      ...inst,
      keepaliveStale: isKeepaliveStale(inst.lastKeepalive),
    })),
  };
}

function loadRegistryEntries(registryRoot, appName, branchName) {
  const files = listRegistryFiles(registryRoot, appName, branchName);
  return files.map((filePath) => {
    const base = path.basename(filePath);
    const parsed = parseRegistryFileName(base);
    const registry = enrichRegistry(readJsonFile(filePath));
    return {
      path: filePath,
      branchName: parsed?.branchName ?? null,
      instanceKey: parsed?.instanceKey ?? null,
      registry,
    };
  });
}

function writeResult(payload) {
  process.stdout.write(`${JSON.stringify(payload, null, 2)}\n`);
}

function fail(message, code = 1) {
  writeResult({ type: 'svc-error', error: message });
  process.exit(code);
}

function requireArg(args, name) {
  const value = args[name];
  if (!value || value === true) {
    fail(`missing required argument: --${name}`);
  }
  return value;
}

function httpRequest(method, host, port, reqPath, timeoutMs = 8000) {
  return new Promise((resolve, reject) => {
    const req = http.request(
      {
        method,
        hostname: host,
        port,
        path: reqPath,
        headers: { Accept: 'application/json' },
        timeout: timeoutMs,
      },
      (res) => {
        let body = '';
        res.setEncoding('utf8');
        res.on('data', (chunk) => {
          body += chunk;
        });
        res.on('end', () => {
          let json = null;
          if (body) {
            try {
              json = JSON.parse(body);
            } catch {
              json = body;
            }
          }
          resolve({ status: res.statusCode, body: json });
        });
      }
    );
    req.on('timeout', () => {
      req.destroy(new Error(`request timeout after ${timeoutMs}ms`));
    });
    req.on('error', reject);
    req.end();
  });
}

async function fetchServiceInfo(host, port) {
  const url = `http://${host}:${port}/__service/info`;
  const result = await httpRequest('GET', host, port, '/__service/info');
  if (result.status < 200 || result.status >= 300) {
    throw new Error(`GET ${url} returned status ${result.status}`);
  }
  return { url, info: result.body };
}

async function postServiceShutdown(host, port) {
  const url = `http://${host}:${port}/__service/shutdown`;
  const result = await httpRequest('POST', host, port, '/__service/shutdown');
  if (result.status < 200 || result.status >= 300) {
    throw new Error(`POST ${url} returned status ${result.status}`);
  }
  return { url, status: result.status, body: result.body };
}

module.exports = {
  DEFAULT_HOST,
  getRegistryRoot,
  sanitizeBranchFileName,
  parseArgs,
  loadRegistryEntries,
  writeResult,
  fail,
  requireArg,
  fetchServiceInfo,
  postServiceShutdown,
};
