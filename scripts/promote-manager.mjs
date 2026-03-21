#!/usr/bin/env node

import { spawnSync } from 'node:child_process';
import path from 'node:path';
import process from 'node:process';
import { fileURLToPath } from 'node:url';

const scriptDir = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(scriptDir, '..');
const backendDir = path.join(repoRoot, 'backend');

const usageText = `Usage:
  node scripts/promote-manager.mjs --user-id <user-id> [--player-id <player-id>]
  node scripts/promote-manager.mjs --subject <line-subject> [--player-id <player-id>]
  node scripts/promote-manager.mjs --provider line --subject <line-subject> [--player-id <player-id>]

Notes:
  - The target user must complete one LINE login first so the backend creates the users row.
  - --subject without --provider defaults to --provider line.
  - Use either --user-id or --provider/--subject lookup, not both.
  - This script only promotes existing users; it never creates a user.

Examples:
  node scripts/promote-manager.mjs --user-id user-123
  node scripts/promote-manager.mjs --subject U0123456789abcdef
  node scripts/promote-manager.mjs --subject U0123456789abcdef --player-id player-ben`;

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

function parseArgs(rawArgs) {
  const parsed = {
    playerId: '',
    provider: '',
    subject: '',
    userId: '',
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
      case '--player-id':
        parsed.playerId = requireValue(rawArgs, index, arg);
        index += 1;
        break;
      case '--provider':
        parsed.provider = requireValue(rawArgs, index, arg);
        index += 1;
        break;
      case '--subject':
        parsed.subject = requireValue(rawArgs, index, arg);
        index += 1;
        break;
      case '--user-id':
        parsed.userId = requireValue(rawArgs, index, arg);
        index += 1;
        break;
      default:
        fail(`unknown argument ${arg}`);
    }
  }

  if (parsed.subject && !parsed.provider) {
    parsed.provider = 'line';
  }

  const lookupByUserId = parsed.userId !== '';
  const lookupByProviderSubject = parsed.provider !== '' || parsed.subject !== '';

  if (lookupByUserId && lookupByProviderSubject) {
    fail('use either --user-id or --provider/--subject');
  }

  if (!lookupByUserId && !lookupByProviderSubject) {
    fail('either --user-id or --subject is required');
  }

  if (parsed.provider !== '' && parsed.subject === '') {
    fail('--subject is required when --provider is provided');
  }

  if (parsed.provider !== '' && parsed.provider !== 'line') {
    fail(`unsupported provider "${parsed.provider}"`);
  }

  return parsed;
}

function requireValue(args, index, flagName) {
  const value = args[index + 1]?.trim();
  if (!value || value.startsWith('--')) {
    fail(`${flagName} requires a value`);
  }

  return value;
}

function buildGoArgs(parsedArgs) {
  const goArgs = ['run', './cmd/admin', 'promote-user'];

  if (parsedArgs.userId) {
    goArgs.push('--user-id', parsedArgs.userId);
  } else {
    goArgs.push('--provider', parsedArgs.provider, '--subject', parsedArgs.subject);
  }

  if (parsedArgs.playerId) {
    goArgs.push('--player-id', parsedArgs.playerId);
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
