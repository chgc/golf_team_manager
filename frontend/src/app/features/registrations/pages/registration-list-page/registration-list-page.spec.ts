import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RegistrationListPage } from './registration-list-page';

describe('RegistrationListPage', () => {
  let component: RegistrationListPage;
  let fixture: ComponentFixture<RegistrationListPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RegistrationListPage]
    })
    .compileComponents();

    fixture = TestBed.createComponent(RegistrationListPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
