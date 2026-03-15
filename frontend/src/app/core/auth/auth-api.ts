import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';

import { AuthPrincipal } from '../../shared/models/auth.models';

@Injectable({
  providedIn: 'root',
})
export class AuthApi {
  private readonly http = inject(HttpClient);

  getCurrentPrincipal(): Observable<AuthPrincipal> {
    return this.http.get<AuthPrincipal>('/api/auth/me');
  }
}
