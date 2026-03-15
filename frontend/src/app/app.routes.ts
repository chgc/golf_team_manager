import { Routes } from '@angular/router';

import { authGuard, pendingLinkGuard } from './core/auth/auth-guard';

export const routes: Routes = [
  {
    path: 'login',
    loadComponent: () =>
      import('./features/auth/pages/login-page/login-page').then((module) => module.LoginPage),
    title: 'Login | Golf Team Manager',
  },
  {
    path: 'auth/done',
    loadComponent: () =>
      import('./features/auth/pages/auth-done-page/auth-done-page').then(
        (module) => module.AuthDonePage,
      ),
    title: 'Completing Sign-in | Golf Team Manager',
  },
  {
    path: 'auth/pending-link',
    canActivate: [pendingLinkGuard],
    loadComponent: () =>
      import('./features/auth/pages/pending-link-page/pending-link-page').then(
        (module) => module.PendingLinkPage,
      ),
    title: 'Account Linking Required | Golf Team Manager',
  },
  {
    path: '',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./features/home/home-page/home-page').then((module) => module.HomePage),
    title: 'Golf Team Manager',
  },
  {
    path: 'players',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./features/players/pages/player-list-page/player-list-page').then(
        (module) => module.PlayerListPage,
      ),
    title: 'Players | Golf Team Manager',
  },
  {
    path: 'sessions',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./features/sessions/pages/session-list-page/session-list-page').then(
        (module) => module.SessionListPage,
      ),
    title: 'Sessions | Golf Team Manager',
  },
  {
    path: '**',
    redirectTo: '',
  },
];
