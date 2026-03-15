import { HttpErrorResponse } from '@angular/common/http';
import { computed, inject, Injectable, signal } from '@angular/core';
import { firstValueFrom } from 'rxjs';

import { AuthPrincipal, AuthSessionStatus } from '../../shared/models/auth.models';
import { AuthApi } from './auth-api';
import { AUTH_RUNTIME_CONFIG, buildLineLoginUrl } from './auth-config';

@Injectable({
  providedIn: 'root',
})
export class AuthShell {
  private readonly authApi = inject(AuthApi);
  private readonly config = inject(AUTH_RUNTIME_CONFIG);
  private readonly initializedState = signal(false);
  private readonly principalState = signal<AuthPrincipal | null>(null);
  private readonly statusState = signal<AuthSessionStatus>('loading');
  private readonly tokenState = signal<string | null>(this.readStoredToken());

  readonly authMode = this.config.authMode;
  readonly principal = this.principalState.asReadonly();
  readonly status = this.statusState.asReadonly();
  readonly isInitialized = this.initializedState.asReadonly();
  readonly token = this.tokenState.asReadonly();
  readonly isAuthenticated = computed(
    () => this.status() === 'authenticated' && this.principal() !== null,
  );
  readonly isUnauthenticated = computed(() => this.status() === 'unauthenticated');
  readonly roleLabel = computed(() =>
    this.principal()?.role === 'manager' ? 'Manager' : 'Player',
  );
  readonly authModeLabel = computed(() =>
    this.isDevelopmentStub() ? 'Dev Stub' : 'LINE SSO',
  );
  readonly isDevelopmentStub = computed(() => this.config.authMode === 'dev_stub');
  readonly isLineMode = computed(() => this.config.authMode === 'line');
  readonly isUnlinkedPlayer = computed(() => {
    const principal = this.principal();
    return this.status() === 'authenticated' && principal?.role === 'player' && !principal.playerId;
  });

  async initialize(): Promise<void> {
    if (this.initializedState()) {
      return;
    }

    if (this.isDevelopmentStub()) {
      await this.refreshPrincipal();
      this.initializedState.set(true);
      return;
    }

    if (!this.requireUsableToken()) {
      this.principalState.set(null);
      this.statusState.set('unauthenticated');
      this.initializedState.set(true);
      return;
    }

    await this.refreshPrincipal();
    this.initializedState.set(true);
  }

  async refreshPrincipal(): Promise<boolean> {
    if (this.isLineMode() && !this.requireUsableToken()) {
      this.principalState.set(null);
      this.statusState.set('unauthenticated');
      return false;
    }

    this.statusState.set('loading');

    try {
      const principal = await firstValueFrom(this.authApi.getCurrentPrincipal());
      this.principalState.set(principal);
      this.statusState.set('authenticated');
      return true;
    } catch (error) {
      this.handlePrincipalRefreshError(error);
      return false;
    }
  }

  async retryDevelopmentBootstrap(): Promise<boolean> {
    return this.refreshPrincipal();
  }

  async completeLineAuthentication(token: string): Promise<boolean> {
    this.storeToken(token);
    return this.refreshPrincipal();
  }

  logout(): void {
    this.clearToken();
    this.clearPendingRedirect();
    this.principalState.set(null);
    this.statusState.set('unauthenticated');
  }

  getLineLoginUrl(): string {
    return buildLineLoginUrl(this.config);
  }

  rememberPendingRedirect(redirectUrl: string | null | undefined): void {
    const normalizedRedirect = normalizeRedirectUrl(redirectUrl);
    if (!normalizedRedirect) {
      this.clearPendingRedirect();
      return;
    }

    writeToStorage(sessionStorage, this.config.redirectStorageKey, normalizedRedirect);
  }

  consumePendingRedirect(): string | null {
    const redirectUrl = normalizeRedirectUrl(
      readFromStorage(sessionStorage, this.config.redirectStorageKey),
    );
    this.clearPendingRedirect();
    return redirectUrl;
  }

  getPostAuthenticationRedirect(fallbackUrl: string): string {
    if (this.isUnlinkedPlayer()) {
      return '/auth/pending-link';
    }

    return normalizeRedirectUrl(this.consumePendingRedirect()) ?? fallbackUrl;
  }

  private handlePrincipalRefreshError(error: unknown): void {
    if (this.isLineMode() && isUnauthorized(error)) {
      this.clearToken();
    }

    this.principalState.set(null);
    this.statusState.set('unauthenticated');
  }

  private readStoredToken(): string | null {
    const storedToken = readFromStorage(localStorage, this.config.tokenStorageKey);
    if (!storedToken) {
      return null;
    }

    if (isExpiredToken(storedToken)) {
      removeFromStorage(localStorage, this.config.tokenStorageKey);
      return null;
    }

    return storedToken;
  }

  private requireUsableToken(): string | null {
    const token = this.readStoredToken();
    this.tokenState.set(token);
    return token;
  }

  private storeToken(token: string): void {
    const normalizedToken = token.trim();
    if (!normalizedToken || isExpiredToken(normalizedToken)) {
      this.clearToken();
      return;
    }

    writeToStorage(localStorage, this.config.tokenStorageKey, normalizedToken);
    this.tokenState.set(normalizedToken);
  }

  private clearToken(): void {
    removeFromStorage(localStorage, this.config.tokenStorageKey);
    this.tokenState.set(null);
  }

  private clearPendingRedirect(): void {
    removeFromStorage(sessionStorage, this.config.redirectStorageKey);
  }
}

function isUnauthorized(error: unknown): boolean {
  return error instanceof HttpErrorResponse && error.status === 401;
}

function normalizeRedirectUrl(value: string | null | undefined): string | null {
  if (!value) {
    return null;
  }

  const trimmedValue = value.trim();
  if (!trimmedValue.startsWith('/') || trimmedValue.startsWith('//')) {
    return null;
  }

  return trimmedValue;
}

function isExpiredToken(token: string): boolean {
  const expirationSeconds = readTokenExpiration(token);
  if (expirationSeconds === null) {
    return true;
  }

  return expirationSeconds <= Math.floor(Date.now() / 1000);
}

function readTokenExpiration(token: string): number | null {
  const parts = token.split('.');
  if (parts.length < 2) {
    return null;
  }

  try {
    const payload = JSON.parse(decodeBase64Url(parts[1])) as { exp?: unknown };
    return typeof payload.exp === 'number' ? payload.exp : null;
  } catch {
    return null;
  }
}

function decodeBase64Url(value: string): string {
  const normalizedValue = value.replace(/-/g, '+').replace(/_/g, '/');
  const paddingLength = (4 - (normalizedValue.length % 4)) % 4;
  return globalThis.atob(normalizedValue.padEnd(normalizedValue.length + paddingLength, '='));
}

function readFromStorage(storage: Storage, key: string): string | null {
  try {
    return storage.getItem(key);
  } catch {
    return null;
  }
}

function writeToStorage(storage: Storage, key: string, value: string) {
  try {
    storage.setItem(key, value);
  } catch {
    // Ignore storage write failures so auth bootstrap can continue in restricted environments.
  }
}

function removeFromStorage(storage: Storage, key: string) {
  try {
    storage.removeItem(key);
  } catch {
    // Ignore storage removal failures so logout can continue in restricted environments.
  }
}
