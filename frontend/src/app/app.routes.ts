import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./features/home/home-page/home-page').then((module) => module.HomePage),
    title: 'Golf Team Manager',
  },
];
