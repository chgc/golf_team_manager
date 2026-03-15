import { ChangeDetectionStrategy, Component, inject, OnInit, signal } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';

import { AuthShell } from '../../../../core/auth/auth-shell';

type CompletionState = 'working' | 'failed';

@Component({
  imports: [MatButtonModule, MatCardModule, RouterLink],
  templateUrl: './auth-done-page.html',
  styleUrl: './auth-done-page.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AuthDonePage implements OnInit {
  private readonly authShell = inject(AuthShell);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);

  protected readonly completionState = signal<CompletionState>('working');
  protected readonly errorMessage = signal<string | null>(null);

  async ngOnInit() {
    const token = extractToken(this.route.snapshot.fragment);
    if (!token) {
      this.completionState.set('failed');
      this.errorMessage.set('The login callback did not include an app token.');
      return;
    }

    const isAuthenticated = await this.authShell.completeLineAuthentication(token);
    if (!isAuthenticated) {
      this.completionState.set('failed');
      this.errorMessage.set('Sign-in finished, but the app could not load your session.');
      return;
    }

    await this.router.navigateByUrl(this.authShell.getPostAuthenticationRedirect('/'));
  }
}

function extractToken(fragment: string | null): string | null {
  if (!fragment) {
    return null;
  }

  const fragmentParams = new URLSearchParams(fragment);
  return fragmentParams.get('token');
}
