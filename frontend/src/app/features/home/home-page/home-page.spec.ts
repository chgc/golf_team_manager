import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';

import { HomePage } from './home-page';

describe('HomePage', () => {
  let component: HomePage;
  let fixture: ComponentFixture<HomePage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HomePage],
      providers: [provideRouter([])],
    })
    .compileComponents();

    fixture = TestBed.createComponent(HomePage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('renders an empty-safe fallback when no shell sections are configured', () => {
    (component as any).shellSections.set([]);
    fixture.detectChanges();

    expect(fixture.nativeElement.textContent).toContain('No feature shortcuts are configured yet.');
  });
});
