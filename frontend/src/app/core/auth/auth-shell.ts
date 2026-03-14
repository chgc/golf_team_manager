import { computed, Injectable, signal } from '@angular/core';

import { AuthPrincipal } from '../../shared/models/auth.models';

@Injectable({
  providedIn: 'root',
})
export class AuthShell {
  private readonly principalState = signal<AuthPrincipal>({
    displayName: 'Demo Manager',
    provider: 'dev_stub',
    role: 'manager',
    subject: 'dev-manager',
    userId: 'dev-manager',
  });

  readonly principal = this.principalState.asReadonly();
  readonly roleLabel = computed(() =>
    this.principal().role === 'manager' ? 'Manager' : 'Player',
  );
  readonly isDevelopmentStub = computed(() => this.principal().provider === 'dev_stub');

  setDevelopmentPrincipal(principal: AuthPrincipal) {
    this.principalState.set(principal);
  }
}
