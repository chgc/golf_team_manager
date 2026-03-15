import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { HttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { signal } from '@angular/core';
import { TestBed } from '@angular/core/testing';

import { AuthShell } from './auth-shell';
import { authInterceptor } from './auth-interceptor';

describe('authInterceptor', () => {
  let httpClient: HttpClient;
  let httpTestingController: HttpTestingController;
  const token = signal<string | null>(null);

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        provideHttpClient(withInterceptors([authInterceptor])),
        provideHttpClientTesting(),
        {
          provide: AuthShell,
          useValue: {
            token,
          },
        },
      ],
    });

    httpClient = TestBed.inject(HttpClient);
    httpTestingController = TestBed.inject(HttpTestingController);
    token.set(null);
  });

  afterEach(() => {
    httpTestingController.verify();
  });

  it('adds a bearer token when one is available', () => {
    token.set('test-token');

    httpClient.get('/api/players').subscribe();

    const request = httpTestingController.expectOne('/api/players');
    expect(request.request.headers.get('Authorization')).toBe('Bearer test-token');
    request.flush([]);
  });

  it('leaves requests unchanged when no token is present', () => {
    httpClient.get('/api/players').subscribe();

    const request = httpTestingController.expectOne('/api/players');
    expect(request.request.headers.has('Authorization')).toBeFalse();
    request.flush([]);
  });
});
