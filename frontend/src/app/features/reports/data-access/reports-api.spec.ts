import { TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';

import { ReportsApi } from './reports-api';

describe('ReportsApi', () => {
  let service: ReportsApi;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideHttpClient(), provideHttpClientTesting()],
    });
    service = TestBed.inject(ReportsApi);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
