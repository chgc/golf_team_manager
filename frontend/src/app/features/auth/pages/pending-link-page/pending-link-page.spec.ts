import { signal } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';

import { AuthShell } from '../../../../core/auth/auth-shell';
import { PendingLinkPage } from './pending-link-page';

describe('PendingLinkPage', () => {
  let component: PendingLinkPage;
  let fixture: ComponentFixture<PendingLinkPage>;
  const logout = jasmine.createSpy('logout');

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PendingLinkPage],
      providers: [
        provideRouter([]),
        {
          provide: AuthShell,
          useValue: {
            principal: signal({
              displayName: 'Unlinked Player',
              provider: 'line',
              role: 'player',
              subject: 'line-user',
              userId: 'user-1',
            }),
            logout,
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(PendingLinkPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
    logout.calls.reset();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
    expect((fixture.nativeElement as HTMLElement).textContent).toContain('Unlinked Player');
  });
});
