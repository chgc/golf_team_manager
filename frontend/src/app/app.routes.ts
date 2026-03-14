import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./features/home/home-page/home-page').then((module) => module.HomePage),
    title: 'Golf Team Manager',
  },
  {
    path: 'players',
    loadComponent: () =>
      import('./features/players/pages/player-list-page/player-list-page').then(
        (module) => module.PlayerListPage,
      ),
    title: 'Players | Golf Team Manager',
  },
  {
    path: 'sessions',
    loadComponent: () =>
      import('./features/sessions/pages/session-list-page/session-list-page').then(
        (module) => module.SessionListPage,
      ),
    title: 'Sessions | Golf Team Manager',
  },
  {
    path: 'registrations',
    loadComponent: () =>
      import('./features/registrations/pages/registration-list-page/registration-list-page').then(
        (module) => module.RegistrationListPage,
      ),
    title: 'Registrations | Golf Team Manager',
  },
  {
    path: '**',
    redirectTo: '',
  },
];
