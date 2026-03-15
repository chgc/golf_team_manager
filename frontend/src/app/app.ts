import { ChangeDetectionStrategy, Component, computed, inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatToolbarModule } from '@angular/material/toolbar';
import { Router, RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';

import { AuthShell } from './core/auth/auth-shell';

interface NavigationItem {
  readonly label: string;
  readonly path: string;
}

@Component({
  selector: 'app-root',
  imports: [MatButtonModule, MatToolbarModule, RouterLink, RouterLinkActive, RouterOutlet],
  templateUrl: './app.html',
  styleUrl: './app.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class App {
  private readonly authShell = inject(AuthShell);
  private readonly router = inject(Router);

  protected readonly applicationTitle = 'Golf Team Manager';
  protected readonly authStatus = this.authShell.status;
  protected readonly isAuthenticated = this.authShell.isAuthenticated;
  protected readonly authPrincipal = this.authShell.principal;
  protected readonly authRoleLabel = this.authShell.roleLabel;
  protected readonly authModeLabel = this.authShell.authModeLabel;
  protected readonly navigationItems = computed<readonly NavigationItem[]>(() => {
    const items: NavigationItem[] = [
      { label: 'Home', path: '/' },
      { label: 'Players', path: '/players' },
      { label: 'Sessions', path: '/sessions' },
    ];

    if (this.authPrincipal()?.role === 'manager') {
      items.push({ label: 'Admin Users', path: '/admin/users' });
    }

    return items;
  });

  protected logout() {
    this.authShell.logout();
    void this.router.navigateByUrl('/login');
  }
}
