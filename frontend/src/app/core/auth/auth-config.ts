import { InjectionToken } from '@angular/core';

import { AuthProvider } from '../../shared/models/auth.models';

declare global {
  interface Window {
    __GTM_AUTH_CONFIG?: Partial<Pick<AuthRuntimeConfig, 'authMode' | 'backendOrigin'>>;
  }
}

export interface AuthRuntimeConfig {
  readonly authMode: AuthProvider;
  readonly backendOrigin: string;
  readonly tokenStorageKey: string;
  readonly redirectStorageKey: string;
}

const defaultTokenStorageKey = 'golf-team-manager.auth-token';
const defaultRedirectStorageKey = 'golf-team-manager.auth-redirect';

export const AUTH_RUNTIME_CONFIG = new InjectionToken<AuthRuntimeConfig>('AUTH_RUNTIME_CONFIG', {
  providedIn: 'root',
  factory: () => loadAuthRuntimeConfig(),
});

export function buildLineLoginUrl(config: AuthRuntimeConfig): string {
  return `${config.backendOrigin}/api/auth/line/login`;
}

function loadAuthRuntimeConfig(): AuthRuntimeConfig {
  const runtimeConfig = globalThis.window?.__GTM_AUTH_CONFIG;
  const locationOrigin = globalThis.location?.origin;

  return {
    authMode: normalizeAuthMode(runtimeConfig?.authMode),
    backendOrigin: normalizeBackendOrigin(runtimeConfig?.backendOrigin, locationOrigin),
    tokenStorageKey: defaultTokenStorageKey,
    redirectStorageKey: defaultRedirectStorageKey,
  };
}

function normalizeAuthMode(value: AuthProvider | undefined): AuthProvider {
  return 'line';
}

function normalizeBackendOrigin(value: string | undefined, locationOrigin: string | undefined): string {
  const trimmedValue = value?.trim();
  if (trimmedValue) {
    return trimmedValue.replace(/\/+$/, '');
  }

  if (!locationOrigin) {
    return 'http://localhost:8080';
  }

  if (locationOrigin.endsWith(':4200')) {
    return 'http://localhost:8080';
  }

  return locationOrigin.replace(/\/+$/, '');
}
