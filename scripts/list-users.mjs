#!/usr/bin/env node

import { spawnSync } from 'node:child_process';
import path from 'node:path';
import process from 'node:process';
import { fileURLToPath } from 'node:url';

const scriptDir = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(scriptDir, '..');
const backendDir = path.join(repoRoot, 'backend');

const usageText = `Usage:
  node scripts/list-users.mjs
  node scripts/list-users.mjs --role manager
  node scripts/list-users.mjs --link-state unlinked
  node scripts/list-users.mjs --role player --link-state linked

Notes:
  - This script lists auth users from the backend database and prints their user IDs.
  - Valid roles: manager, player
  - Valid link states: linked, unlinked`;

function fail(message) {
  console.error(`Error: ${message}\n`);
  console.error(usageText);
  process.exit(1);
}

function printUsage(exitCode = 0) {
  const stream = exitCode === 0 ? process.stdout : process.stderr;
  stream.write(`${usageText}\n`);
  process.exit(exitCode);
}

function requireValue(args, index, flagName) {
  const value = args[index + 1]?.trim();
  if (!value || value.startsWith('--')) {
    fail(`${flagName} requires a value`);
  }

  return value;
}

function parseArgs(rawArgs) {
  const parsed = {
    linkState: '',
    role: '',
  };

  for (let index = 0; index < rawArgs.length; index += 1) {
    const arg = rawArgs[index];
    switch (arg) {
      case '--':
        break;
      case '--help':
      case '-h':
        printUsage(0);
        break;
      case '--link-state':
        parsed.linkState = requireValue(rawArgs, index, arg);
        index += 1;
        break;
      case '--role':
        parsed.role = requireValue(rawArgs, index, arg);
        index += 1;
        break;
      default:
        fail(`unknown argument ${arg}`);
    }
  }

  if (parsed.role !== '' && parsed.role !== 'manager' && parsed.role !== 'player') {
    fail(`unsupported role "${parsed.role}"`);
  }

  if (parsed.linkState !== '' && parsed.linkState !== 'linked' && parsed.linkState !== 'unlinked') {
    fail(`unsupported link-state "${parsed.linkState}"`);
  }

  return parsed;
}

function buildGoArgs(parsedArgs) {
  const goArgs = ['run', './cmd/admin', 'list-users'];

  if (parsedArgs.role) {
    goArgs.push('--role', parsedArgs.role);
  }

  if (parsedArgs.linkState) {
    goArgs.push('--link-state', parsedArgs.linkState);
  }

  return goArgs;
}

const parsedArgs = parseArgs(process.argv.slice(2));
const goArgs = buildGoArgs(parsedArgs);

const result = spawnSync('go', goArgs, {
  cwd: backendDir,
  stdio: 'inherit',
});

if (result.error) {
  console.error(`Failed to run Go admin command: ${result.error.message}`);
  process.exit(1);
}

process.exit(result.status ?? 1);
