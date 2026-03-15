import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';

import {
  AdminUserLinkState,
  AdminUserReadDto,
  AdminUserRoleFilter,
  AdminUserUpdateDto,
} from '../../../shared/models/domain.models';

export interface AdminUserListFilters {
  readonly linkState?: AdminUserLinkState;
  readonly role?: AdminUserRoleFilter;
}

@Injectable({
  providedIn: 'root',
})
export class AdminUsersApi {
  private readonly http = inject(HttpClient);
  private readonly resourcePath = '/api/admin/users';

  listUsers(filters?: AdminUserListFilters) {
    let params = new HttpParams();

    if (filters?.linkState && filters.linkState !== 'all') {
      params = params.set('linkState', filters.linkState);
    }

    if (filters?.role && filters.role !== 'all') {
      params = params.set('role', filters.role);
    }

    return this.http.get<AdminUserReadDto[]>(this.resourcePath, { params });
  }

  updateUser(userId: string, payload: AdminUserUpdateDto) {
    return this.http.patch<AdminUserReadDto>(`${this.resourcePath}/${userId}`, payload);
  }
}
