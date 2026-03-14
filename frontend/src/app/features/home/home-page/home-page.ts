import { ChangeDetectionStrategy, Component, signal } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { RouterLink } from '@angular/router';

interface ShellSection {
  readonly description: string;
  readonly path: string;
  readonly title: string;
}

@Component({
  selector: 'app-home-page',
  imports: [MatButtonModule, MatCardModule, RouterLink],
  templateUrl: './home-page.html',
  styleUrl: './home-page.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class HomePage {
  protected readonly shellSections = signal<readonly ShellSection[]>([
    {
      title: 'Players',
      path: '/players',
      description: 'Prepare player list, form, and validation UX around the shared domain contract.',
    },
    {
      title: 'Sessions',
      path: '/sessions',
      description: 'Align session pages with backend status lifecycle, capacity rules, and reporting flow.',
    },
  ]);
}
