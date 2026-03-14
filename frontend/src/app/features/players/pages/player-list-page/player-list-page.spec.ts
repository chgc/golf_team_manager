import { ComponentFixture, TestBed } from '@angular/core/testing';
import { of, Subject, throwError } from 'rxjs';

import { PlayersApi } from '../../data-access/players-api';
import { PlayerListPage } from './player-list-page';

describe('PlayerListPage', () => {
  let playersApi: jasmine.SpyObj<PlayersApi>;
  let component: PlayerListPage;
  let fixture: ComponentFixture<PlayerListPage>;

  beforeEach(async () => {
    playersApi = jasmine.createSpyObj<PlayersApi>('PlayersApi', ['listPlayers', 'getPlayer', 'createPlayer', 'updatePlayer']);
    playersApi.listPlayers.and.returnValue(of([]));
    playersApi.getPlayer.and.returnValue(of());
    playersApi.createPlayer.and.returnValue(of());
    playersApi.updatePlayer.and.returnValue(of());

    await TestBed.configureTestingModule({
      imports: [PlayerListPage],
      providers: [
        {
          provide: PlayersApi,
          useValue: playersApi,
        },
      ],
    })
    .compileComponents();

    fixture = TestBed.createComponent(PlayerListPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('shows a loading state while the player list request is in flight', () => {
    const playersSubject = new Subject<never[]>();
    playersApi.listPlayers.calls.reset();
    playersApi.listPlayers.and.returnValue(playersSubject.asObservable());

    fixture = TestBed.createComponent(PlayerListPage);
    component = fixture.componentInstance;
    fixture.detectChanges();

    expect(fixture.nativeElement.textContent).toContain('Loading players...');
  });

  it('shows an empty state when the player list is empty', () => {
    expect(fixture.nativeElement.textContent).toContain('No players match the current filters.');
  });

  it('shows the API error message when player loading fails', () => {
    playersApi.listPlayers.calls.reset();
    playersApi.listPlayers.and.returnValue(
      throwError(() => ({
        error: {
          error: {
            message: 'players request failed',
          },
        },
      })),
    );

    fixture = TestBed.createComponent(PlayerListPage);
    component = fixture.componentInstance;
    fixture.detectChanges();

    expect(fixture.nativeElement.textContent).toContain('players request failed');
  });
});
