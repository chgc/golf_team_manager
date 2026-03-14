import { ChangeDetectionStrategy, Component } from '@angular/core';
import { MatCardModule } from '@angular/material/card';

interface PlayerShellItem {
  readonly description: string;
  readonly title: string;
}

@Component({
  imports: [MatCardModule],
  templateUrl: './player-list-page.html',
  styleUrl: './player-list-page.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class PlayerListPage {
  protected readonly focusAreas: readonly PlayerShellItem[] = [
    {
      title: 'Data access',
      description: 'Use PlayersApi as the single HTTP boundary for player list and create flows.',
    },
    {
      title: 'Feature scope',
      description: 'Reserve room for list, search/filter, create/edit, and active/inactive status flows.',
    },
    {
      title: 'Validation',
      description: 'Mirror backend handicap and duplicate-name guidance in future reactive forms.',
    },
  ];
}
