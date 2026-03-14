import { DatePipe } from '@angular/common';
import { computed, ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { NonNullableFormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';

import { AuthShell } from '../../../../core/auth/auth-shell';
import { PlayersApi } from '../../../players/data-access/players-api';
import { ReportsApi } from '../../../reports/data-access/reports-api';
import { SessionsApi } from '../../data-access/sessions-api';
import { RegistrationsApi } from '../../../registrations/data-access/registrations-api';
import {
  PlayerReadDto,
  RegistrationReadDto,
  RegistrationStatus,
  ReservationSummaryReadDto,
  SessionReadDto,
  SessionStatus,
  SessionWriteDto,
} from '../../../../shared/models/domain.models';

interface SessionRosterEntry {
  readonly id: string;
  readonly playerId: string;
  readonly playerName: string;
  readonly status: RegistrationStatus;
  readonly updatedAt: string;
}

type ReservationSummaryState = 'hidden' | 'loading' | 'ready' | 'ineligible' | 'empty' | 'error';

@Component({
  imports: [
    MatButtonModule,
    MatCardModule,
    DatePipe,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    ReactiveFormsModule,
  ],
  templateUrl: './session-list-page.html',
  styleUrl: './session-list-page.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SessionListPage {
  private readonly authShell = inject(AuthShell);
  private readonly formBuilder = inject(NonNullableFormBuilder);
  private readonly playersApi = inject(PlayersApi);
  private readonly reportsApi = inject(ReportsApi);
  private readonly registrationsApi = inject(RegistrationsApi);
  private readonly sessionsApi = inject(SessionsApi);

  protected readonly activePlayers = signal<PlayerReadDto[]>([]);
  protected readonly copyFeedbackMessage = signal<string | null>(null);
  protected readonly errorMessage = signal<string | null>(null);
  protected readonly isSaving = signal(false);
  protected readonly allPlayers = signal<PlayerReadDto[]>([]);
  protected readonly registrations = signal<RegistrationReadDto[]>([]);
  protected readonly reservationSummary = signal<ReservationSummaryReadDto | null>(null);
  protected readonly reservationSummaryMessage = signal<string | null>(null);
  protected readonly reservationSummaryState = signal<ReservationSummaryState>('hidden');
  protected readonly selectedSessionId = signal<string | null>(null);
  protected readonly selectedView = signal<'upcoming' | 'history'>('upcoming');
  protected readonly sessions = signal<SessionReadDto[]>([]);

  protected readonly sessionForm = this.formBuilder.group({
    date: ['', [Validators.required]],
    courseName: ['', [Validators.required]],
    courseAddress: '',
    maxPlayers: [4, [Validators.required, Validators.min(1)]],
    registrationDeadline: ['', [Validators.required]],
    status: 'open' as SessionStatus,
    notes: '',
  });
  protected readonly managerRegistrationForm = this.formBuilder.group({
    playerId: ['', [Validators.required]],
  });

  protected readonly authPrincipal = this.authShell.principal;
  protected readonly isManager = computed(() => this.authPrincipal().role === 'manager');
  protected readonly currentPlayerId = computed(() => this.authPrincipal().playerId ?? null);

  protected readonly selectedSession = computed(() => {
    const selectedId = this.selectedSessionId();
    if (!selectedId) {
      return null;
    }

    return this.sessions().find((session) => session.id === selectedId) ?? null;
  });

  protected readonly formModeLabel = computed(() =>
    this.selectedSessionId() ? 'Edit session' : 'Create session',
  );

  protected readonly upcomingSessions = computed(() =>
    this.sessions()
      .filter((session) => isUpcomingSession(session))
      .sort((left, right) => left.date.localeCompare(right.date)),
  );

  protected readonly historySessions = computed(() =>
    this.sessions()
      .filter((session) => !isUpcomingSession(session))
      .sort((left, right) => right.date.localeCompare(left.date)),
  );

  protected readonly visibleSessions = computed(() =>
    this.selectedView() === 'upcoming' ? this.upcomingSessions() : this.historySessions(),
  );

  protected readonly confirmedRegistrationCount = computed(
    () => this.registrations().filter((registration) => registration.status === 'confirmed').length,
  );

  protected readonly remainingSpots = computed(() => {
    const session = this.selectedSession();
    if (!session) {
      return 0;
    }

    return Math.max(session.maxPlayers - this.confirmedRegistrationCount(), 0);
  });

  protected readonly estimatedGroups = computed(() => {
    const confirmedCount = this.confirmedRegistrationCount();
    if (confirmedCount === 0) {
      return 0;
    }

    return Math.ceil(confirmedCount / 4);
  });

  protected readonly rosterEntries = computed<readonly SessionRosterEntry[]>(() => {
    const playersById = new Map(this.allPlayers().map((player) => [player.id, player] as const));

    return this.registrations().map((registration) => ({
      id: registration.id,
      playerId: registration.playerId,
      playerName: playersById.get(registration.playerId)?.name ?? registration.playerId,
      status: registration.status,
      updatedAt: registration.updatedAt,
    }));
  });

  protected readonly currentPlayerRegistration = computed(() => {
    const currentPlayerId = this.currentPlayerId();
    if (!currentPlayerId) {
      return null;
    }

    return this.registrations().find((registration) => registration.playerId === currentPlayerId) ?? null;
  });

  protected readonly availableManagerPlayers = computed(() => {
    const registeredPlayerIds = new Set(this.registrations().map((registration) => registration.playerId));
    return this.activePlayers().filter((player) => !registeredPlayerIds.has(player.id));
  });

  constructor() {
    this.loadPlayers();
    this.reloadSessions();
  }

  protected showView(view: 'upcoming' | 'history') {
    this.selectedView.set(view);
  }

  protected startCreateSession() {
    this.selectedSessionId.set(null);
    this.registrations.set([]);
    this.resetReservationSummary();
    this.sessionForm.reset({
      date: '',
      courseName: '',
      courseAddress: '',
      maxPlayers: 4,
      registrationDeadline: '',
      status: 'open',
      notes: '',
    });
  }

  protected viewSession(sessionId: string) {
    this.loadSession(sessionId, false);
  }

  protected editSession(sessionId: string) {
    this.loadSession(sessionId, true);
  }

  protected submitSession() {
    if (this.sessionForm.invalid) {
      this.sessionForm.markAllAsTouched();
      return;
    }

    this.errorMessage.set(null);
    this.isSaving.set(true);

    const payload = this.toSessionWriteDto();
    const selectedSessionId = this.selectedSessionId();
    const request$ = selectedSessionId
      ? this.sessionsApi.updateSession(selectedSessionId, payload)
      : this.sessionsApi.createSession(payload);

    request$.subscribe({
      next: (session) => {
        this.isSaving.set(false);
        this.upsertSession(session);
        this.selectedSessionId.set(session.id);
        this.patchForm(session);
        this.loadRegistrations(session.id);
        this.loadReservationSummary(session);
      },
      error: (error: unknown) => {
        this.isSaving.set(false);
        this.errorMessage.set(extractErrorMessage(error));
      },
    });
  }

  protected applyStatus(session: SessionReadDto, nextStatus: SessionStatus) {
    this.isSaving.set(true);
    this.errorMessage.set(null);

    this.sessionsApi
      .updateSession(session.id, {
        date: session.date,
        courseName: session.courseName,
        courseAddress: session.courseAddress,
        maxPlayers: session.maxPlayers,
        registrationDeadline: session.registrationDeadline,
        status: nextStatus,
        notes: session.notes,
      })
      .subscribe({
        next: (updatedSession) => {
          this.isSaving.set(false);
          this.upsertSession(updatedSession);
          this.selectedSessionId.set(updatedSession.id);
          this.patchForm(updatedSession);
          this.loadRegistrations(updatedSession.id);
          this.loadReservationSummary(updatedSession);
        },
        error: (error: unknown) => {
          this.isSaving.set(false);
          this.errorMessage.set(extractErrorMessage(error));
        },
      });
  }

  protected registerCurrentPlayer() {
    const selectedSession = this.selectedSession();
    const currentPlayerId = this.currentPlayerId();
    if (!selectedSession || !currentPlayerId) {
      return;
    }

    this.isSaving.set(true);
    this.errorMessage.set(null);

    const existingRegistration = this.currentPlayerRegistration();
    const request$ = existingRegistration
      ? this.registrationsApi.updateRegistration(existingRegistration.id, { status: 'confirmed' })
      : this.registrationsApi.createRegistration(selectedSession.id, {
          playerId: currentPlayerId,
          status: 'confirmed',
        });

    request$.subscribe({
      next: () => {
        this.isSaving.set(false);
        this.loadRegistrations(selectedSession.id);
      },
      error: (error: unknown) => {
        this.isSaving.set(false);
        this.errorMessage.set(extractErrorMessage(error));
      },
    });
  }

  protected leaveCurrentPlayerRegistration() {
    const selectedSession = this.selectedSession();
    const currentRegistration = this.currentPlayerRegistration();
    if (!selectedSession || !currentRegistration) {
      return;
    }

    this.isSaving.set(true);
    this.errorMessage.set(null);
    this.registrationsApi.updateRegistration(currentRegistration.id, { status: 'cancelled' }).subscribe({
      next: () => {
        this.isSaving.set(false);
        this.loadRegistrations(selectedSession.id);
      },
      error: (error: unknown) => {
        this.isSaving.set(false);
        this.errorMessage.set(extractErrorMessage(error));
      },
    });
  }

  protected submitManagerRegistration() {
    const selectedSession = this.selectedSession();
    if (!selectedSession) {
      return;
    }

    if (this.managerRegistrationForm.invalid) {
      this.managerRegistrationForm.markAllAsTouched();
      return;
    }

    this.isSaving.set(true);
    this.errorMessage.set(null);
    this.registrationsApi
      .createRegistration(selectedSession.id, {
        playerId: this.managerRegistrationForm.controls.playerId.value,
        status: 'confirmed',
      })
      .subscribe({
        next: () => {
          this.isSaving.set(false);
          this.managerRegistrationForm.reset({ playerId: '' });
          this.loadRegistrations(selectedSession.id);
        },
        error: (error: unknown) => {
          this.isSaving.set(false);
          this.errorMessage.set(extractErrorMessage(error));
        },
      });
  }

  protected updateRegistrationStatus(registrationId: string, status: RegistrationStatus) {
    const selectedSession = this.selectedSession();
    if (!selectedSession) {
      return;
    }

    this.isSaving.set(true);
    this.errorMessage.set(null);
    this.registrationsApi.updateRegistration(registrationId, { status }).subscribe({
      next: () => {
        this.isSaving.set(false);
        this.loadRegistrations(selectedSession.id);
      },
      error: (error: unknown) => {
        this.isSaving.set(false);
        this.errorMessage.set(extractErrorMessage(error));
      },
    });
  }

  protected copyReservationSummary() {
    const summary = this.reservationSummary();
    if (!summary) {
      return;
    }

    const clipboard = globalThis.navigator?.clipboard;
    if (!clipboard) {
      this.copyFeedbackMessage.set('Clipboard is unavailable.');
      return;
    }

    void clipboard.writeText(summary.summaryText).then(
      () => {
        this.copyFeedbackMessage.set('Summary copied.');
      },
      () => {
        this.copyFeedbackMessage.set('Copy failed. Please try again.');
      },
    );
  }

  private loadSession(sessionId: string, patchForm: boolean) {
    this.errorMessage.set(null);
    this.sessionsApi.getSession(sessionId).subscribe({
      next: (session) => {
        this.upsertSession(session);
        this.selectedSessionId.set(session.id);
        if (patchForm) {
          this.patchForm(session);
        }
        this.loadRegistrations(session.id);
        this.loadReservationSummary(session);
      },
      error: (error: unknown) => {
        this.errorMessage.set(extractErrorMessage(error));
      },
    });
  }

  private loadPlayers() {
    this.playersApi.listPlayers().subscribe({
      next: (players) => {
        this.allPlayers.set(players);
        this.activePlayers.set(players.filter((player) => player.status === 'active'));
      },
      error: (error: unknown) => {
        this.errorMessage.set(extractErrorMessage(error));
      },
    });
  }

  private loadRegistrations(sessionId: string) {
    this.registrationsApi.listRegistrations(sessionId).subscribe({
      next: (registrations) => {
        this.registrations.set(registrations);
      },
      error: (error: unknown) => {
        this.errorMessage.set(extractErrorMessage(error));
      },
    });
  }

  private patchForm(session: SessionReadDto) {
    this.sessionForm.setValue({
      date: toLocalDateTimeInput(session.date),
      courseName: session.courseName,
      courseAddress: session.courseAddress ?? '',
      maxPlayers: session.maxPlayers,
      registrationDeadline: toLocalDateTimeInput(session.registrationDeadline),
      status: session.status,
      notes: session.notes ?? '',
    });
  }

  private reloadSessions() {
    this.sessionsApi.listSessions().subscribe({
      next: (sessions) => {
        this.sessions.set(sessions);

        const selectedSessionId = this.selectedSessionId();
        if (!selectedSessionId) {
          return;
        }

        const selectedSession = sessions.find((session) => session.id === selectedSessionId);
        if (!selectedSession) {
          this.startCreateSession();
          return;
        }

        this.loadRegistrations(selectedSession.id);
        this.loadReservationSummary(selectedSession);
      },
      error: (error: unknown) => {
        this.errorMessage.set(extractErrorMessage(error));
      },
    });
  }

  private toSessionWriteDto(): SessionWriteDto {
    const rawValue = this.sessionForm.getRawValue();
    return {
      date: new Date(rawValue.date).toISOString(),
      courseName: rawValue.courseName.trim(),
      courseAddress: rawValue.courseAddress.trim(),
      maxPlayers: rawValue.maxPlayers,
      registrationDeadline: new Date(rawValue.registrationDeadline).toISOString(),
      status: rawValue.status,
      notes: rawValue.notes.trim(),
    };
  }

  private upsertSession(session: SessionReadDto) {
    this.sessions.update((currentSessions) => {
      const existingIndex = currentSessions.findIndex((item) => item.id === session.id);
      if (existingIndex === -1) {
        return [...currentSessions, session];
      }

      return currentSessions.map((item) => item.id === session.id ? session : item);
    });
  }

  private loadReservationSummary(session: SessionReadDto) {
    this.copyFeedbackMessage.set(null);
    this.reservationSummary.set(null);

    if (!this.isManager()) {
      this.reservationSummaryState.set('hidden');
      this.reservationSummaryMessage.set(null);
      return;
    }

    if (!isReservationSummaryEligibleStatus(session.status)) {
      this.reservationSummaryState.set('ineligible');
      this.reservationSummaryMessage.set('This session is not ready for a reservation summary yet.');
      return;
    }

    this.reservationSummaryState.set('loading');
    this.reservationSummaryMessage.set(null);
    this.reportsApi.getReservationSummary(session.id).subscribe({
      next: (summary) => {
        this.reservationSummary.set(summary);
        this.reservationSummaryState.set('ready');
        this.reservationSummaryMessage.set(null);
      },
      error: (error: unknown) => {
        const errorCode = extractErrorCode(error);
        switch (errorCode) {
          case 'reservation_summary_empty':
            this.reservationSummaryState.set('empty');
            this.reservationSummaryMessage.set('There are no confirmed players yet, so the reservation summary is unavailable.');
            return;
          case 'session_not_eligible_for_report':
            this.reservationSummaryState.set('ineligible');
            this.reservationSummaryMessage.set('This session is not ready for a reservation summary yet.');
            return;
          default:
            this.reservationSummaryState.set('error');
            this.reservationSummaryMessage.set(extractErrorMessage(error));
        }
      },
    });
  }

  private resetReservationSummary() {
    this.copyFeedbackMessage.set(null);
    this.reservationSummary.set(null);
    this.reservationSummaryMessage.set(null);
    this.reservationSummaryState.set('hidden');
  }
}

function isUpcomingSession(session: SessionReadDto): boolean {
  if (session.status === 'cancelled') {
    return false;
  }

  const sessionDate = new Date(session.date);
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  return sessionDate >= today;
}

function toLocalDateTimeInput(value: string): string {
  return value.slice(0, 16);
}

function isReservationSummaryEligibleStatus(status: SessionStatus): boolean {
  return status === 'confirmed' || status === 'completed';
}

function extractApiError(error: unknown): { code?: string; message?: string; details?: string[] } | null {
  if (typeof error !== 'object' || error === null) {
    return null;
  }

  const maybeError = error as { error?: { error?: { message?: string; details?: string[] } } };
  return maybeError.error?.error ?? null;
}

function extractErrorCode(error: unknown): string | null {
  return extractApiError(error)?.code ?? null;
}

function extractErrorMessage(error: unknown): string {
  const apiError = extractApiError(error);
  if (!apiError) {
    return 'Unexpected error';
  }

  if (Array.isArray(apiError.details) && apiError.details.length > 0) {
    return apiError.details.join(' ');
  }

  return apiError.message ?? 'Unexpected error';
}
