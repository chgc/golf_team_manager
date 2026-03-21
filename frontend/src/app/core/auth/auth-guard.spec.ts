import { computed, signal } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { CanActivateFn, provideRouter } from '@angular/router';

import { AuthShell } from './auth-shell';
import { AuthPrincipal } from '../../shared/models/auth.models';
import { authGuard, managerGuard, pendingLinkGuard } from './auth-guard';

describe('auth guards', () => {
  const authStatus = signal<'loading' | 'authenticated' | 'unauthenticated'>('unauthenticated');
  const isUnlinkedPlayer = signal(false);
  const principal = signal<AuthPrincipal | null>(null);
  const executeAuthGuard: CanActivateFn = (...guardParameters) =>
    TestBed.runInInjectionContext(() => authGuard(...guardParameters));
  const executeManagerGuard: CanActivateFn = (...guardParameters) =>
    TestBed.runInInjectionContext(() => managerGuard(...guardParameters));
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
            principal,
          },
        },
      ],
    });

    authStatus.set('unauthenticated');
    isUnlinkedPlayer.set(false);
    principal.set(null);
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
    principal.set({
      displayName: 'Player One',
      provider: 'line',
      role: 'player',
      subject: 'player-1',
      userId: 'user-1',
    });

    expect(executeAuthGuard({} as never, { url: '/' } as never)).toBeTrue();
  });

  it('allows managers through the manager guard', () => {
    authStatus.set('authenticated');
    principal.set({
      displayName: 'Manager One',
      provider: 'line',
      role: 'manager',
      subject: 'manager-1',
      userId: 'user-1',
    });

    expect(executeManagerGuard({} as never, { url: '/admin/users' } as never)).toBeTrue();
  });

  it('redirects non-managers away from manager-only routes', () => {
    authStatus.set('authenticated');
    principal.set({
      displayName: 'Player One',
      provider: 'line',
      role: 'player',
      subject: 'player-1',
      userId: 'user-1',
    });

    const result = executeManagerGuard({} as never, { url: '/admin/users' } as never);

    expect(result.toString()).toBe('/');
  });

  it('redirects linked users away from the pending-link route', () => {
    authStatus.set('authenticated');

    const result = executePendingLinkGuard({} as never, { url: '/auth/pending-link' } as never);

    expect(result.toString()).toBe('/');
  });
});
