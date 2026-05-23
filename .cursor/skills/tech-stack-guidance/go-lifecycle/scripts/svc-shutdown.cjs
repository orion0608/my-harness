#!/usr/bin/env node
'use strict';

const {
  DEFAULT_HOST,
  postServiceShutdown,
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

  const result = await postServiceShutdown(host, port);
  writeResult({
    type: 'svc-shutdown',
    host,
    port,
    url: result.url,
    status: result.status,
    body: result.body,
  });
}

main().catch((err) => fail(err.message));
