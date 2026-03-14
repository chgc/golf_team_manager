import { TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';

import { SessionsApi } from './sessions-api';

describe('SessionsApi', () => {
  let service: SessionsApi;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideHttpClient(), provideHttpClientTesting()],
    });
    service = TestBed.inject(SessionsApi);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
