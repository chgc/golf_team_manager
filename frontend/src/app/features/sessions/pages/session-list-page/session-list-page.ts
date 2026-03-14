import { DatePipe } from '@angular/common';
import { computed, ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { NonNullableFormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';

import { SessionsApi } from '../../data-access/sessions-api';
import { RegistrationsApi } from '../../../registrations/data-access/registrations-api';
import { SessionReadDto, SessionStatus, SessionWriteDto, RegistrationReadDto } from '../../../../shared/models/domain.models';

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
  private readonly formBuilder = inject(NonNullableFormBuilder);
  private readonly registrationsApi = inject(RegistrationsApi);
  private readonly sessionsApi = inject(SessionsApi);

  protected readonly errorMessage = signal<string | null>(null);
  protected readonly isSaving = signal(false);
  protected readonly registrations = signal<RegistrationReadDto[]>([]);
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

  constructor() {
    this.reloadSessions();
  }

  protected showView(view: 'upcoming' | 'history') {
    this.selectedView.set(view);
  }

  protected startCreateSession() {
    this.selectedSessionId.set(null);
    this.registrations.set([]);
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
        },
        error: (error: unknown) => {
          this.isSaving.set(false);
          this.errorMessage.set(extractErrorMessage(error));
        },
      });
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
