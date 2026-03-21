import { computed, signal } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';

import { App } from './app';
import { AuthShell } from './core/auth/auth-shell';
import { AuthPrincipal, AuthSessionStatus } from './shared/models/auth.models';

describe('App', () => {
  const authStatus = signal<AuthSessionStatus>('authenticated');
  const principal = signal<AuthPrincipal | null>({
    displayName: 'Demo Manager',
    provider: 'line',
    role: 'manager',
    subject: 'dev-manager',
    userId: 'dev-manager',
  });
  const logout = jasmine.createSpy('logout');

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [App],
      providers: [
        provideRouter([]),
        {
          provide: AuthShell,
          useValue: {
            status: authStatus,
            principal,
            roleLabel: computed(() => (principal()?.role === 'manager' ? 'Manager' : 'Player')),
            authModeLabel: signal('LINE SSO'),
            isAuthenticated: computed(() => authStatus() === 'authenticated' && principal() !== null),
            logout,
          },
        },
      ],
    }).compileComponents();

    authStatus.set('authenticated');
    principal.set({
      displayName: 'Demo Manager',
      provider: 'line',
      role: 'manager',
      subject: 'dev-manager',
      userId: 'dev-manager',
    });
    logout.calls.reset();
  });

  it('should create the app', () => {
    const fixture = TestBed.createComponent(App);
    const app = fixture.componentInstance;
    expect(app).toBeTruthy();
  });

  it('should render the application title in the toolbar', () => {
    const fixture = TestBed.createComponent(App);
    fixture.detectChanges();
    const compiled = fixture.nativeElement as HTMLElement;
    expect(compiled.querySelector('.app-title')?.textContent).toContain('Golf Team Manager');
  });

  it('should render the authenticated identity badge and logout action', () => {
    const fixture = TestBed.createComponent(App);
    fixture.detectChanges();
    const compiled = fixture.nativeElement as HTMLElement;
    expect(compiled.querySelector('.identity-badge')?.textContent).toContain('Demo Manager');
    expect(compiled.textContent).toContain('Logout');
  });

  it('shows the admin navigation entry for managers', () => {
    const fixture = TestBed.createComponent(App);
    fixture.detectChanges();

    expect((fixture.nativeElement as HTMLElement).textContent).toContain('Admin Users');
  });

  it('hides the admin navigation entry for non-managers', () => {
    principal.set({
      displayName: 'Demo Player',
      provider: 'line',
      role: 'player',
      subject: 'dev-player',
      userId: 'dev-player',
      playerId: 'player-1',
    });

    const fixture = TestBed.createComponent(App);
    fixture.detectChanges();

    expect((fixture.nativeElement as HTMLElement).textContent).not.toContain('Admin Users');
  });

  it('shows the login action when no principal is active', () => {
    authStatus.set('unauthenticated');
    principal.set(null);

    const fixture = TestBed.createComponent(App);
    fixture.detectChanges();

    expect((fixture.nativeElement as HTMLElement).textContent).toContain('Login');
  });
});
