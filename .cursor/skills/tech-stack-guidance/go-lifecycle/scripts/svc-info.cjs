#!/usr/bin/env node
'use strict';

const {
  DEFAULT_HOST,
  fetchServiceInfo,
  writeResult,
  fail,
  requireArg,
  parseArgs,
} = require('./_svc-common.cjs');

async function main() {
  const args = parseArgs(process.argv);
  const port = Number(requireArg(args, 'port'));
  if (!Number.isInteger(port) || port <= 0) {
    fail('--port must be a positive integer');
  }
  const host = args.host ? String(args.host) : DEFAULT_HOST;

  const result = await fetchServiceInfo(host, port);
  writeResult({
    type: 'svc-info',
    host,
    port,
    url: result.url,
    info: result.info,
  });
}

main().catch((err) => fail(err.message));
