import { TestBed } from '@angular/core/testing';

import { AuthShell } from './auth-shell';

describe('AuthShell', () => {
  let service: AuthShell;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(AuthShell);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should default to manager development identity', () => {
    expect(service.principal().role).toBe('manager');
    expect(service.isDevelopmentStub()).toBeTrue();
  });
});
