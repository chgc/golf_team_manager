import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';

import { PlayerFilterStatus, PlayerReadDto, PlayerWriteDto } from '../../../shared/models/domain.models';

export interface PlayerListFilters {
  readonly query?: string;
  readonly status?: PlayerFilterStatus;
}

@Injectable({
  providedIn: 'root',
})
export class PlayersApi {
  private readonly http = inject(HttpClient);
  private readonly resourcePath = '/api/players';

  listPlayers(filters?: PlayerListFilters) {
    let params = new HttpParams();
    if (filters?.query?.trim()) {
      params = params.set('query', filters.query.trim());
    }

    if (filters?.status && filters.status !== 'all') {
      params = params.set('status', filters.status);
    }

    return this.http.get<PlayerReadDto[]>(this.resourcePath, { params });
  }

  createPlayer(payload: PlayerWriteDto) {
    return this.http.post<PlayerReadDto>(this.resourcePath, payload);
  }

  getPlayer(playerId: string) {
    return this.http.get<PlayerReadDto>(`${this.resourcePath}/${playerId}`);
  }

  updatePlayer(playerId: string, payload: PlayerWriteDto) {
    return this.http.patch<PlayerReadDto>(`${this.resourcePath}/${playerId}`, payload);
  }
}
