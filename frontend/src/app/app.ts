import { ChangeDetectionStrategy, Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatToolbarModule } from '@angular/material/toolbar';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';

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
  protected readonly applicationTitle = 'Golf Team Manager';
  protected readonly navigationItems: readonly NavigationItem[] = [
    { label: 'Home', path: '/' },
    { label: 'Players', path: '/players' },
    { label: 'Sessions', path: '/sessions' },
    { label: 'Registrations', path: '/registrations' },
  ];
}
