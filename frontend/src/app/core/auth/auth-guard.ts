import { CanActivateFn, Router } from '@angular/router';
import { inject } from '@angular/core';

import { AuthShell } from './auth-shell';

export const authGuard: CanActivateFn = (_, state) => {
  const authShell = inject(AuthShell);
  const router = inject(Router);

  if (authShell.status() !== 'authenticated') {
    return buildLoginRedirect(router, state.url);
  }

  if (authShell.isUnlinkedPlayer() && state.url !== '/auth/pending-link') {
    return router.createUrlTree(['/auth/pending-link']);
  }

  return true;
};

export const pendingLinkGuard: CanActivateFn = (_, state) => {
  const authShell = inject(AuthShell);
  const router = inject(Router);

  if (authShell.status() !== 'authenticated') {
    return buildLoginRedirect(router, state.url);
  }

  if (!authShell.isUnlinkedPlayer()) {
    return router.createUrlTree(['/']);
  }

  return true;
};

export const managerGuard: CanActivateFn = (_, state) => {
  const authShell = inject(AuthShell);
  const router = inject(Router);

  if (authShell.status() !== 'authenticated') {
    return buildLoginRedirect(router, state.url);
  }

  if (authShell.isUnlinkedPlayer() && state.url !== '/auth/pending-link') {
    return router.createUrlTree(['/auth/pending-link']);
  }

  if (authShell.principal()?.role !== 'manager') {
    return router.createUrlTree(['/']);
  }

  return true;
};

function buildLoginRedirect(router: Router, redirectUrl: string) {
  return router.createUrlTree(['/login'], {
    queryParams: {
      redirect: redirectUrl,
    },
  });
}
