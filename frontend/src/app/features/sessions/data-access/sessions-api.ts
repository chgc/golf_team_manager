import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';

import { SessionReadDto, SessionWriteDto } from '../../../shared/models/domain.models';

@Injectable({
  providedIn: 'root',
})
export class SessionsApi {
  private readonly http = inject(HttpClient);
  private readonly resourcePath = '/api/sessions';

  listSessions() {
    return this.http.get<SessionReadDto[]>(this.resourcePath);
  }

  createSession(payload: SessionWriteDto) {
    return this.http.post<SessionReadDto>(this.resourcePath, payload);
  }
}
