import { ComponentFixture, TestBed } from '@angular/core/testing';
import { of, throwError } from 'rxjs';

import { PlayersApi } from '../../../players/data-access/players-api';
import { AdminUsersApi } from '../../data-access/admin-users-api';
import { AdminUserPage } from './admin-user-page';

describe('AdminUserPage', () => {
  let adminUsersApi: jasmine.SpyObj<AdminUsersApi>;
  let playersApi: jasmine.SpyObj<PlayersApi>;
  let component: AdminUserPage;
  let fixture: ComponentFixture<AdminUserPage>;

  beforeEach(async () => {
    adminUsersApi = jasmine.createSpyObj<AdminUsersApi>('AdminUsersApi', ['listUsers', 'updateUser']);
    playersApi = jasmine.createSpyObj<PlayersApi>('PlayersApi', ['listPlayers']);

    adminUsersApi.listUsers.and.returnValue(
      of([
        {
          userId: 'user-1',
          displayName: 'Manager One',
          provider: 'line',
          subject: 'line-manager-123456',
          role: 'manager',
          playerId: 'player-1',
          createdAt: '2026-03-15T10:00:00Z',
          updatedAt: '2026-03-15T10:30:00Z',
        },
        {
          userId: 'user-2',
          displayName: 'Pending Player',
          provider: 'line',
          subject: 'line-player-654321',
          role: 'player',
          createdAt: '2026-03-15T11:00:00Z',
          updatedAt: '2026-03-15T11:30:00Z',
        },
      ]),
    );
    adminUsersApi.updateUser.and.returnValue(
      of({
        userId: 'user-2',
        displayName: 'Pending Player',
        provider: 'line',
        subject: 'line-player-654321',
        role: 'manager',
        playerId: 'player-2',
        createdAt: '2026-03-15T11:00:00Z',
        updatedAt: '2026-03-15T12:00:00Z',
      }),
    );
    playersApi.listPlayers.and.returnValue(
      of([
        {
          id: 'player-1',
          name: 'Alice',
          handicap: 8,
          status: 'active',
          createdAt: '2026-03-15T10:00:00Z',
          updatedAt: '2026-03-15T10:00:00Z',
        },
        {
          id: 'player-2',
          name: 'Bob',
          handicap: 12,
          status: 'active',
          createdAt: '2026-03-15T10:00:00Z',
          updatedAt: '2026-03-15T10:00:00Z',
        },
      ]),
    );

    await TestBed.configureTestingModule({
      imports: [AdminUserPage],
      providers: [
        {
          provide: AdminUsersApi,
          useValue: adminUsersApi,
        },
        {
          provide: PlayersApi,
          useValue: playersApi,
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(AdminUserPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('renders both linked and unlinked account sections', () => {
    const textContent = (fixture.nativeElement as HTMLElement).textContent ?? '';

    expect(component).toBeTruthy();
    expect(textContent).toContain('Unlinked Accounts');
    expect(textContent).toContain('Linked Accounts');
    expect(textContent).toContain('Pending Player');
    expect(textContent).toContain('Manager One');
  });

  it('shows the API error message when loading users fails', async () => {
    adminUsersApi.listUsers.and.returnValue(
      throwError(() => ({
        error: {
          error: {
            message: 'admin users request failed',
          },
        },
      })),
    );

    fixture = TestBed.createComponent(AdminUserPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
    await fixture.whenStable();

    expect((fixture.nativeElement as HTMLElement).textContent).toContain('admin users request failed');
  });

  it('submits role and player changes through the admin API', () => {
    const user = component['users']().find((entry) => entry.userId === 'user-2');
    expect(user).toBeDefined();

    component['selectedRoles'].update((current) => ({
      ...current,
      'user-2': 'manager',
    }));
    component['selectedPlayerIds'].update((current) => ({
      ...current,
      'user-2': 'player-2',
    }));

    component['saveUser'](user!);

    expect(adminUsersApi.updateUser).toHaveBeenCalledWith('user-2', {
      role: 'manager',
      playerId: 'player-2',
    });
    expect(component['feedbackMessage']()).toContain('Updated Pending Player.');
  });
});
