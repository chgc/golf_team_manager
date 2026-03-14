import { computed, ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { NonNullableFormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';

import { PlayersApi } from '../../data-access/players-api';
import { PlayerFilterStatus, PlayerReadDto, PlayerStatus, PlayerWriteDto } from '../../../../shared/models/domain.models';

@Component({
  imports: [
    MatButtonModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    ReactiveFormsModule,
  ],
  templateUrl: './player-list-page.html',
  styleUrl: './player-list-page.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class PlayerListPage {
  private readonly formBuilder = inject(NonNullableFormBuilder);
  private readonly playersApi = inject(PlayersApi);

  protected readonly players = signal<PlayerReadDto[]>([]);
  protected readonly allPlayers = signal<PlayerReadDto[]>([]);
  protected readonly errorMessage = signal<string | null>(null);
  protected readonly isLoadingPlayers = signal(false);
  protected readonly isSaving = signal(false);
  protected readonly selectedPlayerId = signal<string | null>(null);

  protected readonly filterForm = this.formBuilder.group({
    query: '',
    status: 'all' as PlayerFilterStatus,
  });

  protected readonly playerForm = this.formBuilder.group({
    name: ['', [Validators.required]],
    handicap: [0, [Validators.required, Validators.min(0), Validators.max(54), halfStepValidator]],
    phone: '',
    email: ['', [Validators.email]],
    status: 'active' as PlayerStatus,
    notes: '',
  });

  protected readonly duplicateNameWarning = computed(() => {
    const currentName = this.playerForm.controls.name.value.trim().toLocaleLowerCase();
    if (!currentName) {
      return null;
    }

    const selectedId = this.selectedPlayerId();
    const duplicateExists = this.allPlayers().some(
      (player) =>
        player.id !== selectedId && player.name.trim().toLocaleLowerCase() === currentName,
    );

    return duplicateExists ? 'Duplicate name detected. Saving is still allowed because UUID is the real identity.' : null;
  });

  protected readonly formModeLabel = computed(() =>
    this.selectedPlayerId() ? 'Edit player' : 'Create player',
  );

  constructor() {
    this.reloadPlayers();
    this.reloadWarningSource();
  }

  protected applyFilters() {
    this.reloadPlayers();
  }

  protected editPlayer(playerId: string) {
    this.errorMessage.set(null);
    this.playersApi.getPlayer(playerId).subscribe({
      next: (player) => {
        this.selectedPlayerId.set(player.id);
        this.playerForm.setValue({
          name: player.name,
          handicap: player.handicap,
          phone: player.phone ?? '',
          email: player.email ?? '',
          status: player.status,
          notes: player.notes ?? '',
        });
      },
      error: (error: unknown) => {
        this.errorMessage.set(extractErrorMessage(error));
      },
    });
  }

  protected resetFilters() {
    this.filterForm.setValue({ query: '', status: 'all' });
    this.reloadPlayers();
  }

  protected resetPlayerForm() {
    this.selectedPlayerId.set(null);
    this.playerForm.reset({
      name: '',
      handicap: 0,
      phone: '',
      email: '',
      status: 'active',
      notes: '',
    });
  }

  protected submitPlayer() {
    if (this.playerForm.invalid) {
      this.playerForm.markAllAsTouched();
      return;
    }

    const payload = this.toPlayerWriteDto();
    this.isSaving.set(true);
    this.errorMessage.set(null);

    const selectedPlayerId = this.selectedPlayerId();
    const request$ = selectedPlayerId
      ? this.playersApi.updatePlayer(selectedPlayerId, payload)
      : this.playersApi.createPlayer(payload);

    request$.subscribe({
      next: () => {
        this.isSaving.set(false);
        this.resetPlayerForm();
        this.reloadPlayers();
        this.reloadWarningSource();
      },
      error: (error: unknown) => {
        this.isSaving.set(false);
        this.errorMessage.set(extractErrorMessage(error));
      },
    });
  }

  protected togglePlayerStatus(player: PlayerReadDto) {
    const nextStatus: PlayerStatus = player.status === 'active' ? 'inactive' : 'active';
    this.isSaving.set(true);
    this.errorMessage.set(null);

    this.playersApi
      .updatePlayer(player.id, {
        name: player.name,
        handicap: player.handicap,
        phone: player.phone,
        email: player.email,
        status: nextStatus,
        notes: player.notes,
      })
      .subscribe({
        next: () => {
          this.isSaving.set(false);
          this.reloadPlayers();
          this.reloadWarningSource();
        },
        error: (error: unknown) => {
          this.isSaving.set(false);
          this.errorMessage.set(extractErrorMessage(error));
        },
      });
  }

  private reloadPlayers() {
    const { query, status } = this.filterForm.getRawValue();
    this.isLoadingPlayers.set(true);
    this.playersApi.listPlayers({ query, status }).subscribe({
      next: (players) => {
        this.isLoadingPlayers.set(false);
        this.players.set(players);
      },
      error: (error: unknown) => {
        this.isLoadingPlayers.set(false);
        this.errorMessage.set(extractErrorMessage(error));
      },
    });
  }

  private reloadWarningSource() {
    this.playersApi.listPlayers().subscribe({
      next: (players) => {
        this.allPlayers.set(players);
      },
      error: (error: unknown) => {
        this.errorMessage.set(extractErrorMessage(error));
      },
    });
  }

  private toPlayerWriteDto(): PlayerWriteDto {
    const rawValue = this.playerForm.getRawValue();
    return {
      ...rawValue,
      name: rawValue.name.trim(),
      email: rawValue.email.trim(),
      notes: rawValue.notes.trim(),
      phone: rawValue.phone.trim(),
    };
  }
}

function halfStepValidator(control: { value: number | null }) {
  const value = control.value;
  if (value === null) {
    return null;
  }

  return Number.isInteger(value * 2) ? null : { halfStep: true };
}

function extractErrorMessage(error: unknown): string {
  if (typeof error !== 'object' || error === null) {
    return 'Unexpected error';
  }

  const maybeError = error as { error?: { error?: { message?: string; details?: string[] } } };
  const apiError = maybeError.error?.error;
  if (!apiError) {
    return 'Unexpected error';
  }

  if (Array.isArray(apiError.details) && apiError.details.length > 0) {
    return apiError.details.join(' ');
  }

  return apiError.message ?? 'Unexpected error';
}
