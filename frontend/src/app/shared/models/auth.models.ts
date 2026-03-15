export type AuthRole = 'manager' | 'player';
export type AuthProvider = 'dev_stub' | 'line';
export type AuthSessionStatus = 'loading' | 'authenticated' | 'unauthenticated';

export interface AuthPrincipal {
  displayName: string;
  playerId?: string;
  provider: AuthProvider;
  role: AuthRole;
  subject: string;
  userId: string;
}
