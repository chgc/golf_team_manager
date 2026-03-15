import { computed, signal } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { CanActivateFn, provideRouter } from '@angular/router';

import { AuthShell } from './auth-shell';
import { authGuard, pendingLinkGuard } from './auth-guard';

describe('auth guards', () => {
  const authStatus = signal<'loading' | 'authenticated' | 'unauthenticated'>('unauthenticated');
  const isUnlinkedPlayer = signal(false);
  const executeAuthGuard: CanActivateFn = (...guardParameters) =>
    TestBed.runInInjectionContext(() => authGuard(...guardParameters));
  const executePendingLinkGuard: CanActivateFn = (...guardParameters) =>
    TestBed.runInInjectionContext(() => pendingLinkGuard(...guardParameters));

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        provideRouter([]),
        {
          provide: AuthShell,
          useValue: {
            status: authStatus,
            isUnlinkedPlayer: computed(() => isUnlinkedPlayer()),
          },
        },
      ],
    });

    authStatus.set('unauthenticated');
    isUnlinkedPlayer.set(false);
  });

  it('redirects unauthenticated users to the login page', () => {
    const result = executeAuthGuard({} as never, { url: '/sessions' } as never);

    expect(result.toString()).toBe('/login?redirect=%2Fsessions');
  });

  it('redirects authenticated but unlinked players to the pending-link page', () => {
    authStatus.set('authenticated');
    isUnlinkedPlayer.set(true);

    const result = executeAuthGuard({} as never, { url: '/players' } as never);

    expect(result.toString()).toBe('/auth/pending-link');
  });

  it('allows linked authenticated users through the auth guard', () => {
    authStatus.set('authenticated');

    expect(executeAuthGuard({} as never, { url: '/' } as never)).toBeTrue();
  });

  it('redirects linked users away from the pending-link route', () => {
    authStatus.set('authenticated');

    const result = executePendingLinkGuard({} as never, { url: '/auth/pending-link' } as never);

    expect(result.toString()).toBe('/');
  });
});
