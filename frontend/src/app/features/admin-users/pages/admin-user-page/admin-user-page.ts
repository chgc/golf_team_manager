import { DatePipe } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectChange, MatSelectModule } from '@angular/material/select';

import { PlayersApi } from '../../../players/data-access/players-api';
import { AdminUsersApi } from '../../data-access/admin-users-api';
import { AdminUserReadDto, AdminUserUpdateDto, PlayerReadDto } from '../../../../shared/models/domain.models';
import { AuthRole } from '../../../../shared/models/auth.models';

@Component({
  imports: [DatePipe, MatButtonModule, MatCardModule, MatFormFieldModule, MatSelectModule],
  templateUrl: './admin-user-page.html',
  styleUrl: './admin-user-page.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AdminUserPage {
  private readonly adminUsersApi = inject(AdminUsersApi);
  private readonly playersApi = inject(PlayersApi);

  protected readonly errorMessage = signal<string | null>(null);
  protected readonly feedbackMessage = signal<string | null>(null);
  protected readonly isLoadingPlayers = signal(false);
  protected readonly isLoadingUsers = signal(false);
  protected readonly players = signal<PlayerReadDto[]>([]);
  protected readonly savingUserIds = signal<Record<string, boolean>>({});
  protected readonly selectedPlayerIds = signal<Record<string, string>>({});
  protected readonly selectedRoles = signal<Record<string, AuthRole>>({});
  protected readonly users = signal<AdminUserReadDto[]>([]);

  protected readonly linkedUsers = computed(() => this.users().filter((user) => Boolean(user.playerId)));
  protected readonly unlinkedUsers = computed(() => this.users().filter((user) => !user.playerId));
  protected readonly managerCount = computed(() => this.users().filter((user) => user.role === 'manager').length);

  constructor() {
    this.loadPlayers();
    this.loadUsers();
  }

  protected onRoleChange(userId: string, event: MatSelectChange) {
    this.selectedRoles.update((current) => ({
      ...current,
      [userId]: event.value as AuthRole,
    }));
  }

  protected onPlayerChange(userId: string, event: MatSelectChange) {
    this.selectedPlayerIds.update((current) => ({
      ...current,
      [userId]: event.value as string,
    }));
  }

  protected selectedRole(userId: string) {
    return this.selectedRoles()[userId] ?? 'player';
  }

  protected selectedPlayerId(userId: string) {
    return this.selectedPlayerIds()[userId] ?? '';
  }

  protected hasPendingChanges(user: AdminUserReadDto) {
    return (
      this.selectedRole(user.userId) !== user.role ||
      this.selectedPlayerId(user.userId) !== (user.playerId ?? '')
    );
  }

  protected isSaving(userId: string) {
    return this.savingUserIds()[userId] === true;
  }

  protected reload() {
    this.loadPlayers();
    this.loadUsers();
  }

  protected saveUser(user: AdminUserReadDto) {
    const payload = this.buildUpdatePayload(user);
    if (!payload) {
      return;
    }

    this.errorMessage.set(null);
    this.feedbackMessage.set(null);
    this.setSavingState(user.userId, true);

    this.adminUsersApi.updateUser(user.userId, payload).subscribe({
      next: (updatedUser) => {
        this.replaceUser(updatedUser);
        this.syncSelections(updatedUser);
        this.feedbackMessage.set(`Updated ${updatedUser.displayName}.`);
        this.setSavingState(user.userId, false);
      },
      error: (error: HttpErrorResponse) => {
        this.errorMessage.set(readApiError(error, `Failed to update ${user.displayName}.`));
        this.setSavingState(user.userId, false);
      },
    });
  }

  protected unlinkUser(user: AdminUserReadDto) {
    this.selectedPlayerIds.update((current) => ({
      ...current,
      [user.userId]: '',
    }));
    this.saveUser(user);
  }

  protected maskSubject(subject: string) {
    if (subject.length <= 8) {
      return subject;
    }

    return `${subject.slice(0, 4)}...${subject.slice(-4)}`;
  }

  private buildUpdatePayload(user: AdminUserReadDto): AdminUserUpdateDto | null {
    const selectedRole = this.selectedRole(user.userId);
    const selectedPlayerId = this.selectedPlayerId(user.userId);
    const payload: AdminUserUpdateDto = {};

    if (selectedRole !== user.role) {
      payload.role = selectedRole;
    }

    if (selectedPlayerId !== (user.playerId ?? '')) {
      payload.playerId = selectedPlayerId ? selectedPlayerId : null;
    }

    return Object.keys(payload).length === 0 ? null : payload;
  }

  private loadPlayers() {
    this.isLoadingPlayers.set(true);

    this.playersApi.listPlayers({ status: 'active' }).subscribe({
      next: (players) => {
        this.players.set(players);
        this.isLoadingPlayers.set(false);
      },
      error: (error: HttpErrorResponse) => {
        this.errorMessage.set(readApiError(error, 'Failed to load players.'));
        this.isLoadingPlayers.set(false);
      },
    });
  }

  private loadUsers() {
    this.isLoadingUsers.set(true);

    this.adminUsersApi.listUsers().subscribe({
      next: (users) => {
        this.users.set(users);
        this.primeSelections(users);
        this.isLoadingUsers.set(false);
      },
      error: (error: HttpErrorResponse) => {
        this.errorMessage.set(readApiError(error, 'Failed to load users.'));
        this.isLoadingUsers.set(false);
      },
    });
  }

  private primeSelections(users: readonly AdminUserReadDto[]) {
    const roles: Record<string, AuthRole> = {};
    const playerIds: Record<string, string> = {};

    for (const user of users) {
      roles[user.userId] = user.role;
      playerIds[user.userId] = user.playerId ?? '';
    }

    this.selectedRoles.set(roles);
    this.selectedPlayerIds.set(playerIds);
  }

  private replaceUser(updatedUser: AdminUserReadDto) {
    this.users.update((users) =>
      users.map((user) => (user.userId === updatedUser.userId ? updatedUser : user)),
    );
  }

  private syncSelections(user: AdminUserReadDto) {
    this.selectedRoles.update((current) => ({
      ...current,
      [user.userId]: user.role,
    }));
    this.selectedPlayerIds.update((current) => ({
      ...current,
      [user.userId]: user.playerId ?? '',
    }));
  }

  private setSavingState(userId: string, isSaving: boolean) {
    this.savingUserIds.update((current) => ({
      ...current,
      [userId]: isSaving,
    }));
  }
}

function readApiError(error: HttpErrorResponse, fallbackMessage: string) {
  return (
    (error.error as { error?: { message?: string } } | null)?.error?.message ?? fallbackMessage
  );
}
