import { ChangeDetectionStrategy, Component } from '@angular/core';
import { MatCardModule } from '@angular/material/card';

interface SessionShellItem {
  readonly description: string;
  readonly title: string;
}

@Component({
  imports: [MatCardModule],
  templateUrl: './session-list-page.html',
  styleUrl: './session-list-page.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SessionListPage {
  protected readonly focusAreas: readonly SessionShellItem[] = [
    {
      title: 'Lifecycle visibility',
      description: 'Upcoming list and detail pages will align with open, closed, confirmed, completed, and cancelled states.',
    },
    {
      title: 'Data access',
      description: 'SessionsApi will own list/create requests against the backend foundation routes.',
    },
    {
      title: 'Follow-up work',
      description: 'Later tasks will add edit, auto-close, manual confirm, and reservation-summary entry points.',
    },
  ];
}
