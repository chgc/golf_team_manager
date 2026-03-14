import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';

import { RegistrationReadDto, RegistrationStatusUpdateDto, RegistrationWriteDto } from '../../../shared/models/domain.models';

@Injectable({
  providedIn: 'root',
})
export class RegistrationsApi {
  private readonly http = inject(HttpClient);

  listRegistrations(sessionId: string) {
    return this.http.get<RegistrationReadDto[]>(`/api/sessions/${sessionId}/registrations`);
  }

  createRegistration(sessionId: string, payload: RegistrationWriteDto) {
    return this.http.post<RegistrationReadDto>(`/api/sessions/${sessionId}/registrations`, payload);
  }

  updateRegistration(registrationId: string, payload: RegistrationStatusUpdateDto) {
    return this.http.patch<RegistrationReadDto>(`/api/registrations/${registrationId}`, payload);
  }
}
