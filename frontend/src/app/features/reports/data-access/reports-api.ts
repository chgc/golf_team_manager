import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';

import { ReservationSummaryReadDto } from '../../../shared/models/domain.models';

@Injectable({
  providedIn: 'root',
})
export class ReportsApi {
  private readonly http = inject(HttpClient);
  private readonly resourcePath = '/api/reports/sessions';

  getReservationSummary(sessionId: string) {
    return this.http.get<ReservationSummaryReadDto>(`${this.resourcePath}/${sessionId}/reservation-summary`);
  }
}
