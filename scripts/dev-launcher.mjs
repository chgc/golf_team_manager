#!/usr/bin/env node

import { spawnSync } from 'node:child_process';
import fs from 'node:fs';
import { fileURLToPath } from 'node:url';
import path from 'node:path';
import process from 'node:process';

const scriptDir = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(scriptDir, '..');

const backendCommand = 'just backend-start';
const frontendCommand = 'just frontend-start';
const requiredBackendEnv = [
  'LINE_CLIENT_ID',
  'LINE_CLIENT_SECRET',
  'LINE_REDIRECT_URI',
  'FRONTEND_URL',
  'JWT_SECRET',
];

function run(command, args) {
  const result = spawnSync(command, args, { stdio: 'inherit' });

  if (result.error) {
    throw result.error;
  }

  if (typeof result.status === 'number' && result.status !== 0) {
    throw new Error(`${command} exited with status ${result.status}`);
  }
}

function hasCommand(command, args = []) {
  const result = spawnSync(command, args, { stdio: 'ignore' });
  return !result.error && result.status === 0;
}

function toAppleScriptString(value) {
  return `"${value.replace(/\\/g, '\\\\').replace(/"/g, '\\"')}"`;
}

function loadRootDotEnvValues() {
  const envPath = path.join(repoRoot, '.env');
  if (!fs.existsSync(envPath)) {
    return {};
  }

  const values = {};
  for (const rawLine of fs.readFileSync(envPath, 'utf8').split(/\r?\n/)) {
    const line = rawLine.trim();
    if (!line || line.startsWith('#')) {
      continue;
    }

    const exportLine = line.startsWith('export ') ? line.slice('export '.length).trim() : line;
    const separatorIndex = exportLine.indexOf('=');
    if (separatorIndex === -1) {
      continue;
    }

    const key = exportLine.slice(0, separatorIndex).trim();
    if (!key) {
      continue;
    }

    let value = exportLine.slice(separatorIndex + 1).trim();
    if (
      value.length >= 2 &&
      ((value.startsWith('"') && value.endsWith('"')) || (value.startsWith("'") && value.endsWith("'")))
    ) {
      value = value.slice(1, -1);
    }

    values[key] = value;
  }

  return values;
}

function warnIfBackendEnvMissing() {
  const dotEnvValues = loadRootDotEnvValues();
  const missingEnv = requiredBackendEnv.filter((name) => !process.env[name] && !dotEnvValues[name]);

  if (missingEnv.length === 0) {
    return;
  }

  console.warn(
    `Warning: backend-start will fail until these environment variables are set in the shell or ${path.join(repoRoot, '.env')}: ${missingEnv.join(', ')}`
  );
}

function launchWindows() {
  if (hasCommand('where.exe', ['wt.exe'])) {
    run('wt.exe', [
      '-w',
      '0',
      'new-tab',
      '-d',
      repoRoot,
      'cmd.exe',
      '/k',
      backendCommand,
      ';',
      'new-tab',
      '-d',
      repoRoot,
      'cmd.exe',
      '/k',
      frontendCommand,
    ]);
    return;
  }

  run('cmd.exe', ['/c', 'start', '', '/d', repoRoot, 'cmd.exe', '/k', backendCommand]);
  run('cmd.exe', ['/c', 'start', '', '/d', repoRoot, 'cmd.exe', '/k', frontendCommand]);
}

function launchMacOs() {
  run('osascript', [
    '-e',
    'tell application "Terminal"',
    '-e',
    'activate',
    '-e',
    `do script "cd " & quoted form of ${toAppleScriptString(repoRoot)} & "; ${backendCommand}"`,
    '-e',
    `do script "cd " & quoted form of ${toAppleScriptString(repoRoot)} & "; ${frontendCommand}"`,
    '-e',
    'end tell',
  ]);
}

warnIfBackendEnvMissing();

switch (process.platform) {
  case 'win32':
    launchWindows();
    break;
  case 'darwin':
    launchMacOs();
    break;
  default:
    throw new Error(`just dev is only supported on Windows and macOS, received ${process.platform}.`);
}
