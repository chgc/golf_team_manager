import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';

import { AuthShell } from './auth-shell';
import { AUTH_RUNTIME_CONFIG } from './auth-config';

describe('AuthShell', () => {
  let service: AuthShell;
  let httpTestingController: HttpTestingController;

  beforeEach(() => {
    localStorage.clear();
    sessionStorage.clear();
  });

  afterEach(() => {
    httpTestingController.verify();
    localStorage.clear();
    sessionStorage.clear();
  });

  it('bootstraps immediately from /api/auth/me in dev stub mode', async () => {
    configureTestingModule('dev_stub');

    const initializePromise = service.initialize();
    const request = httpTestingController.expectOne('/api/auth/me');
    expect(request.request.method).toBe('GET');
    request.flush({
      displayName: 'Demo Manager',
      provider: 'dev_stub',
      role: 'manager',
      subject: 'dev-manager',
      userId: 'dev-manager',
    });

    await initializePromise;

    expect(service.status()).toBe('authenticated');
    expect(service.principal()?.displayName).toBe('Demo Manager');
    expect(service.isDevelopmentStub()).toBeTrue();
  });

  it('skips /api/auth/me in line mode when no local token exists', async () => {
    configureTestingModule('line');

    await service.initialize();

    httpTestingController.expectNone('/api/auth/me');
    expect(service.status()).toBe('unauthenticated');
    expect(service.token()).toBeNull();
  });

  it('clears expired line tokens before bootstrap', async () => {
    localStorage.setItem('test.auth-token', createToken(-60));
    configureTestingModule('line');

    await service.initialize();

    httpTestingController.expectNone('/api/auth/me');
    expect(localStorage.getItem('test.auth-token')).toBeNull();
    expect(service.status()).toBe('unauthenticated');
  });

  it('clears the stored line token when /api/auth/me returns 401', async () => {
    localStorage.setItem('test.auth-token', createToken(300));
    configureTestingModule('line');

    const initializePromise = service.initialize();
    const request = httpTestingController.expectOne('/api/auth/me');
    request.flush({}, { status: 401, statusText: 'Unauthorized' });

    await initializePromise;

    expect(localStorage.getItem('test.auth-token')).toBeNull();
    expect(service.status()).toBe('unauthenticated');
  });

  function configureTestingModule(authMode: 'dev_stub' | 'line') {
    TestBed.resetTestingModule();
    TestBed.configureTestingModule({
      providers: [
        provideHttpClient(),
        provideHttpClientTesting(),
        {
          provide: AUTH_RUNTIME_CONFIG,
          useValue: {
            authMode,
            backendOrigin: 'http://localhost:8080',
            tokenStorageKey: 'test.auth-token',
            redirectStorageKey: 'test.auth-redirect',
          },
        },
      ],
    });

    service = TestBed.inject(AuthShell);
    httpTestingController = TestBed.inject(HttpTestingController);
  }
});

function createToken(offsetSeconds: number): string {
  const expirationSeconds = Math.floor(Date.now() / 1000) + offsetSeconds;
  const payload = {
    exp: expirationSeconds,
  };

  return ['header', encodeTokenPart(payload), 'signature'].join('.');
}

function encodeTokenPart(payload: object): string {
  return btoa(JSON.stringify(payload))
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=+$/, '');
}
