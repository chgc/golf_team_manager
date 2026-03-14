import { ComponentFixture, TestBed } from '@angular/core/testing';
import { of, throwError } from 'rxjs';

import { AuthShell } from '../../../../core/auth/auth-shell';
import { PlayersApi } from '../../../players/data-access/players-api';
import { RegistrationsApi } from '../../../registrations/data-access/registrations-api';
import { ReportsApi } from '../../../reports/data-access/reports-api';
import { SessionsApi } from '../../data-access/sessions-api';
import { SessionListPage } from './session-list-page';

describe('SessionListPage', () => {
  const confirmedSession = {
    id: 'session-1',
    date: '2026-10-01T08:00:00Z',
    courseName: 'Sunrise Golf Club',
    courseAddress: '',
    maxPlayers: 8,
    registrationDeadline: '2026-09-25T23:59:00Z',
    status: 'confirmed' as const,
    notes: '',
    createdAt: '2026-09-01T00:00:00Z',
    updatedAt: '2026-09-01T00:00:00Z',
  };
  const openSession = {
    ...confirmedSession,
    status: 'open' as const,
  };
  const summary = {
    sessionId: confirmedSession.id,
    sessionDate: confirmedSession.date,
    courseName: confirmedSession.courseName,
    courseAddress: '',
    registrationDeadline: confirmedSession.registrationDeadline,
    sessionStatus: 'confirmed' as const,
    confirmedPlayerCount: 2,
    estimatedGroups: 1,
    summaryText: 'Session: 2026-10-01T08:00:00Z\nCourse: Sunrise Golf Club\nAddress: N/A\nRoster:\n- Alice\n- Bob',
    confirmedPlayers: [
      { playerId: 'player-1', playerName: 'Alice' },
      { playerId: 'player-2', playerName: 'Bob' },
    ],
  };

  let component: SessionListPage;
  let fixture: ComponentFixture<SessionListPage>;
  let playersApi: jasmine.SpyObj<PlayersApi>;
  let registrationsApi: jasmine.SpyObj<RegistrationsApi>;
  let reportsApi: jasmine.SpyObj<ReportsApi>;
  let sessionsApi: jasmine.SpyObj<SessionsApi>;

  beforeEach(async () => {
    playersApi = jasmine.createSpyObj<PlayersApi>('PlayersApi', ['listPlayers']);
    registrationsApi = jasmine.createSpyObj<RegistrationsApi>('RegistrationsApi', [
      'listRegistrations',
      'createRegistration',
      'updateRegistration',
    ]);
    reportsApi = jasmine.createSpyObj<ReportsApi>('ReportsApi', ['getReservationSummary']);
    sessionsApi = jasmine.createSpyObj<SessionsApi>('SessionsApi', [
      'listSessions',
      'createSession',
      'getSession',
      'updateSession',
    ]);

    playersApi.listPlayers.and.returnValue(of([]));
    registrationsApi.listRegistrations.and.returnValue(of([]));
    registrationsApi.createRegistration.and.returnValue(of());
    registrationsApi.updateRegistration.and.returnValue(of());
    reportsApi.getReservationSummary.and.returnValue(of(summary));
    sessionsApi.listSessions.and.returnValue(of([]));
    sessionsApi.createSession.and.returnValue(of(confirmedSession));
    sessionsApi.getSession.and.returnValue(of(confirmedSession));
    sessionsApi.updateSession.and.returnValue(of(confirmedSession));

    await TestBed.configureTestingModule({
      imports: [SessionListPage],
      providers: [
        {
          provide: SessionsApi,
          useValue: sessionsApi,
        },
        {
          provide: RegistrationsApi,
          useValue: registrationsApi,
        },
        {
          provide: PlayersApi,
          useValue: playersApi,
        },
        {
          provide: ReportsApi,
          useValue: reportsApi,
        },
        {
          provide: AuthShell,
          useValue: {
            principal: () => ({ role: 'manager' }),
          },
        },
      ],
    })
    .compileComponents();
  });

  function createComponent() {
    fixture = TestBed.createComponent(SessionListPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }

  it('should create', () => {
    createComponent();
    expect(component).toBeTruthy();
  });

  it('loads the reservation summary for confirmed sessions', () => {
    sessionsApi.listSessions.and.returnValue(of([confirmedSession]));
    sessionsApi.getSession.and.returnValue(of(confirmedSession));

    createComponent();
    (component as any).viewSession(confirmedSession.id);
    fixture.detectChanges();

    expect(reportsApi.getReservationSummary).toHaveBeenCalledWith(confirmedSession.id);
    expect(fixture.nativeElement.textContent).toContain('Copy-ready summary');
    expect(fixture.nativeElement.textContent).toContain('Alice');
    expect(fixture.nativeElement.textContent).toContain('Bob');
  });

  it('shows an inline ineligible hint without calling the summary API', () => {
    sessionsApi.listSessions.and.returnValue(of([openSession]));
    sessionsApi.getSession.and.returnValue(of(openSession));

    createComponent();
    (component as any).viewSession(openSession.id);
    fixture.detectChanges();

    expect(reportsApi.getReservationSummary).not.toHaveBeenCalled();
    expect(fixture.nativeElement.textContent).toContain('This session is not ready for a reservation summary yet.');
    expect(fixture.nativeElement.querySelector('.error-card')).toBeNull();
  });

  it('keeps reservation summary empty state inside the summary card', () => {
    sessionsApi.listSessions.and.returnValue(of([confirmedSession]));
    sessionsApi.getSession.and.returnValue(of(confirmedSession));
    reportsApi.getReservationSummary.and.returnValue(
      throwError(() => ({
        error: {
          error: {
            code: 'reservation_summary_empty',
            message: 'reservation summary requires at least one confirmed player',
          },
        },
      })),
    );

    createComponent();
    (component as any).viewSession(confirmedSession.id);
    fixture.detectChanges();

    expect(fixture.nativeElement.textContent).toContain('There are no confirmed players yet, so the reservation summary is unavailable.');
    expect(fixture.nativeElement.querySelector('.error-card')).toBeNull();
  });

  it('copies the backend summary text through the clipboard API', async () => {
    const writeText = jasmine.createSpy().and.returnValue(Promise.resolve());
    Object.defineProperty(window.navigator, 'clipboard', {
      configurable: true,
      value: { writeText },
    });

    sessionsApi.listSessions.and.returnValue(of([confirmedSession]));
    sessionsApi.getSession.and.returnValue(of(confirmedSession));

    createComponent();
    (component as any).viewSession(confirmedSession.id);
    fixture.detectChanges();

    (component as any).copyReservationSummary();
    await fixture.whenStable();
    fixture.detectChanges();

    expect(writeText).toHaveBeenCalledWith(summary.summaryText);
    expect(fixture.nativeElement.textContent).toContain('Summary copied.');
  });
});
