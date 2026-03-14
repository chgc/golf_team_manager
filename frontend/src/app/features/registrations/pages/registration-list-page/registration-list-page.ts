import { ChangeDetectionStrategy, Component } from '@angular/core';
import { MatCardModule } from '@angular/material/card';

interface RegistrationShellItem {
  readonly description: string;
  readonly title: string;
}

@Component({
  imports: [MatCardModule],
  templateUrl: './registration-list-page.html',
  styleUrl: './registration-list-page.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RegistrationListPage {
  protected readonly focusAreas: readonly RegistrationShellItem[] = [
    {
      title: 'Player actions',
      description: 'Reserve space for self-service register and cancel flows once auth is in place.',
    },
    {
      title: 'Manager overrides',
      description: 'Future manager controls can build on the same page structure without breaking the shell.',
    },
    {
      title: 'Data access',
      description: 'RegistrationsApi already aligns with session-scoped backend routes.',
    },
  ];
}
