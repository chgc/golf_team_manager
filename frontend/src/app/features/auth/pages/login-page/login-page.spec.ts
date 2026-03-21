import { computed, signal } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, convertToParamMap, provideRouter } from '@angular/router';

import { AuthShell } from '../../../../core/auth/auth-shell';
import { AuthPrincipal } from '../../../../shared/models/auth.models';
import { LoginPage } from './login-page';

describe('LoginPage', () => {
  let component: LoginPage;
  let fixture: ComponentFixture<LoginPage>;
  const authStatus = signal<'loading' | 'authenticated' | 'unauthenticated'>('unauthenticated');
  const principal = signal<AuthPrincipal | null>(null);
  const rememberPendingRedirect = jasmine.createSpy('rememberPendingRedirect');

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LoginPage],
      providers: [
        provideRouter([]),
        {
          provide: ActivatedRoute,
          useValue: {
            snapshot: {
              queryParamMap: convertToParamMap({ redirect: '/sessions' }),
            },
          },
        },
        {
          provide: AuthShell,
          useValue: {
            principal,
            status: authStatus,
            isAuthenticated: computed(
              () => authStatus() === 'authenticated' && principal() !== null,
            ),
            isLineMode: signal(true),
            isUnlinkedPlayer: signal(false),
            authModeLabel: signal('LINE SSO'),
            roleLabel: signal('Player'),
            getLineLoginUrl: () => 'http://localhost:8080/api/auth/line/login',
            rememberPendingRedirect,
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(LoginPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  beforeEach(() => {
    rememberPendingRedirect.calls.reset();
  });

  it('renders a direct backend login link for line mode', () => {
    const loginLink = fixture.nativeElement.querySelector(
      '.login-page__auth-action',
    ) as HTMLAnchorElement;

    expect(loginLink.getAttribute('href')).toBe('http://localhost:8080/api/auth/line/login');
    expect(loginLink.textContent).toContain('Login with LINE');
  });

  it('remembers the protected redirect before leaving for line auth', () => {
    (component as unknown as { rememberRedirect: () => void }).rememberRedirect();

    expect(rememberPendingRedirect).toHaveBeenCalledWith('/sessions');
  });

  it('shows the continue button for authenticated users', () => {
    authStatus.set('authenticated');
    principal.set({
      displayName: 'Taylor Player',
      provider: 'line',
      role: 'player',
      playerId: 'player-1',
      subject: 'line-user-1',
      userId: 'user-1',
    });
    fixture.detectChanges();

    const continueButton = fixture.nativeElement.querySelector(
      'button.login-page__auth-action',
    ) as HTMLButtonElement | null;

    expect(continueButton?.textContent).toContain('Continue to Golf Team Manager');
  });
});
