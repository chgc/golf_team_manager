import { ComponentFixture, TestBed } from '@angular/core/testing';
import { of } from 'rxjs';

import { PlayersApi } from '../../data-access/players-api';
import { PlayerListPage } from './player-list-page';

describe('PlayerListPage', () => {
  let component: PlayerListPage;
  let fixture: ComponentFixture<PlayerListPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PlayerListPage],
      providers: [
        {
          provide: PlayersApi,
          useValue: {
            listPlayers: () => of([]),
            getPlayer: () => of(),
            createPlayer: () => of(),
            updatePlayer: () => of(),
          },
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
});
