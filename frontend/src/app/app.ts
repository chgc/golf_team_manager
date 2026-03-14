import { computed, ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatToolbarModule } from '@angular/material/toolbar';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';

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

  protected readonly applicationTitle = 'Golf Team Manager';
  protected readonly authPrincipal = this.authShell.principal;
  protected readonly authRoleLabel = this.authShell.roleLabel;
  protected readonly authModeLabel = computed(() =>
    this.authShell.isDevelopmentStub() ? 'Dev Stub' : 'External Auth',
  );
  protected readonly navigationItems: readonly NavigationItem[] = [
    { label: 'Home', path: '/' },
    { label: 'Players', path: '/players' },
    { label: 'Sessions', path: '/sessions' },
  ];
}
