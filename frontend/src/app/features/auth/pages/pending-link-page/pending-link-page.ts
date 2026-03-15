import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { Router } from '@angular/router';

import { AuthShell } from '../../../../core/auth/auth-shell';

@Component({
  imports: [MatButtonModule, MatCardModule],
  templateUrl: './pending-link-page.html',
  styleUrl: './pending-link-page.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class PendingLinkPage {
  private readonly authShell = inject(AuthShell);
  private readonly router = inject(Router);

  protected readonly authPrincipal = this.authShell.principal;

  protected logout() {
    this.authShell.logout();
    void this.router.navigateByUrl('/login');
  }
}
