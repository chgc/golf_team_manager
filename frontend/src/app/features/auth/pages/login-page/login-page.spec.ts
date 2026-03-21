import { computed, signal } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, convertToParamMap, provideRouter } from '@angular/router';

import { AuthShell } from '../../../../core/auth/auth-shell';
import { LoginPage } from './login-page';

describe('LoginPage', () => {
  let component: LoginPage;
  let fixture: ComponentFixture<LoginPage>;
  const authStatus = signal<'loading' | 'authenticated' | 'unauthenticated'>('unauthenticated');
  const principal = signal(null);
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
            isAuthenticated: computed(() => authStatus() === 'authenticated' && principal() !== null),
            isLineMode: signal(true),
            isUnlinkedPlayer: signal(false),
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
    const loginLink = fixture.nativeElement.querySelector('a') as HTMLAnchorElement;

    expect(loginLink.href).toBe('http://localhost:8080/api/auth/line/login');
  });

  it('remembers the protected redirect before leaving for line auth', () => {
    (component as unknown as { rememberRedirect: () => void }).rememberRedirect();

    expect(rememberPendingRedirect).toHaveBeenCalledWith('/sessions');
  });
});

