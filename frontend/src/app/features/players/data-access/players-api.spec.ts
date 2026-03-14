import { TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';

import { PlayersApi } from './players-api';

describe('PlayersApi', () => {
  let service: PlayersApi;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideHttpClient(), provideHttpClientTesting()],
    });
    service = TestBed.inject(PlayersApi);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
