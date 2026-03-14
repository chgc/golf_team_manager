import { ComponentFixture, TestBed } from '@angular/core/testing';
import { of } from 'rxjs';

import { AuthShell } from '../../../../core/auth/auth-shell';
import { PlayersApi } from '../../../players/data-access/players-api';
import { RegistrationsApi } from '../../../registrations/data-access/registrations-api';
import { SessionsApi } from '../../data-access/sessions-api';
import { SessionListPage } from './session-list-page';

describe('SessionListPage', () => {
  let component: SessionListPage;
  let fixture: ComponentFixture<SessionListPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SessionListPage],
      providers: [
        {
          provide: SessionsApi,
          useValue: {
            listSessions: () => of([]),
            createSession: () => of(),
            getSession: () => of(),
            updateSession: () => of(),
          },
        },
        {
          provide: RegistrationsApi,
          useValue: {
            listRegistrations: () => of([]),
            createRegistration: () => of(),
            updateRegistration: () => of(),
          },
        },
        {
          provide: PlayersApi,
          useValue: {
            listPlayers: () => of([]),
          },
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

    fixture = TestBed.createComponent(SessionListPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
