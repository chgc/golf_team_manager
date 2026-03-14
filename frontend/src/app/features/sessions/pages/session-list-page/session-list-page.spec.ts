import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SessionListPage } from './session-list-page';

describe('SessionListPage', () => {
  let component: SessionListPage;
  let fixture: ComponentFixture<SessionListPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SessionListPage]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SessionListPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
