#!/usr/bin/env node
'use strict';

const {
  getRegistryRoot,
  loadRegistryEntries,
  writeResult,
  fail,
  requireArg,
  parseArgs,
} = require('./_svc-common.cjs');

async function main() {
  const args = parseArgs(process.argv);
  const appName = requireArg(args, 'app');
  const branchName = requireArg(args, 'branch');
  const registryRoot = getRegistryRoot(args['registry-root']);

  const files = loadRegistryEntries(registryRoot, appName, branchName);
  writeResult({
    type: 'svc-registry-branch',
    appName,
    branchName,
    registryRoot,
    count: files.length,
    files,
  });
}

main().catch((err) => fail(err.message));
