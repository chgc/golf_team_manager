import { computed, ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { ActivatedRoute, Router } from '@angular/router';

import { AuthShell } from '../../../../core/auth/auth-shell';

@Component({
  imports: [MatButtonModule, MatCardModule],
  templateUrl: './login-page.html',
  styleUrl: './login-page.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoginPage {
  private readonly authShell = inject(AuthShell);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);

  protected readonly authPrincipal = this.authShell.principal;
  protected readonly authStatus = this.authShell.status;
  protected readonly isAuthenticated = this.authShell.isAuthenticated;
  protected readonly isLineMode = this.authShell.isLineMode;
  protected readonly isUnlinkedPlayer = this.authShell.isUnlinkedPlayer;
  protected readonly lineLoginUrl = computed(() => this.authShell.getLineLoginUrl());
  protected readonly redirectTarget = computed(
    () => normalizeRedirect(this.route.snapshot.queryParamMap.get('redirect')) ?? '/',
  );

  protected rememberRedirect() {
    this.authShell.rememberPendingRedirect(this.redirectTarget());
  }

  protected async continueAfterAuthentication() {
    const redirectUrl = this.isUnlinkedPlayer() ? '/auth/pending-link' : this.redirectTarget();
    await this.router.navigateByUrl(redirectUrl);
  }
}

function normalizeRedirect(value: string | null): string | null {
  if (!value || !value.startsWith('/') || value.startsWith('//')) {
    return null;
  }

  return value;
}

