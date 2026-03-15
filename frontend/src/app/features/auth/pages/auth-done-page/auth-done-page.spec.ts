import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, provideRouter, Router } from '@angular/router';

import { AuthShell } from '../../../../core/auth/auth-shell';
import { AuthDonePage } from './auth-done-page';

describe('AuthDonePage', () => {
  let component: AuthDonePage;
  let fixture: ComponentFixture<AuthDonePage>;
  let router: Router;
  const completeLineAuthentication = jasmine
    .createSpy('completeLineAuthentication')
    .and.resolveTo(true);
  const getPostAuthenticationRedirect = jasmine
    .createSpy('getPostAuthenticationRedirect')
    .and.returnValue('/sessions');

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AuthDonePage],
      providers: [
        provideRouter([]),
        {
          provide: ActivatedRoute,
          useValue: {
            snapshot: {
              fragment: 'token=test-token',
            },
          },
        },
        {
          provide: AuthShell,
          useValue: {
            completeLineAuthentication,
            getPostAuthenticationRedirect,
          },
        },
      ],
    }).compileComponents();

    router = TestBed.inject(Router);
    spyOn(router, 'navigateByUrl').and.resolveTo(true);

    fixture = TestBed.createComponent(AuthDonePage);
    component = fixture.componentInstance;
  });

  beforeEach(() => {
    completeLineAuthentication.calls.reset();
    getPostAuthenticationRedirect.calls.reset();
  });

  it('stores the callback token through the auth shell and redirects afterwards', async () => {
    fixture.detectChanges();
    await fixture.whenStable();

    expect(component).toBeTruthy();
    expect(completeLineAuthentication).toHaveBeenCalledWith('test-token');
    expect(getPostAuthenticationRedirect).toHaveBeenCalledWith('/');
    expect(router.navigateByUrl).toHaveBeenCalledWith('/sessions');
  });
});
