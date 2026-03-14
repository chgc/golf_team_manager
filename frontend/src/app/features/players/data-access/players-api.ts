import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';

import { PlayerReadDto, PlayerWriteDto } from '../../../shared/models/domain.models';

@Injectable({
  providedIn: 'root',
})
export class PlayersApi {
  private readonly http = inject(HttpClient);
  private readonly resourcePath = '/api/players';

  listPlayers() {
    return this.http.get<PlayerReadDto[]>(this.resourcePath);
  }

  createPlayer(payload: PlayerWriteDto) {
    return this.http.post<PlayerReadDto>(this.resourcePath, payload);
  }
}
