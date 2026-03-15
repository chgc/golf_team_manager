export type PlayerStatus = 'active' | 'inactive';

export interface PlayerWriteDto {
  name: string;
  handicap: number;
  phone?: string;
  email?: string;
  status: PlayerStatus;
  notes?: string;
}

export interface PlayerReadDto extends PlayerWriteDto {
  id: string;
  createdAt: string;
  updatedAt: string;
}

export type PlayerFilterStatus = 'all' | PlayerStatus;

export type SessionStatus = 'open' | 'closed' | 'confirmed' | 'completed' | 'cancelled';

export interface SessionWriteDto {
  date: string;
  courseName: string;
  courseAddress?: string;
  maxPlayers: number;
  registrationDeadline: string;
  status: SessionStatus;
  notes?: string;
}

export interface SessionReadDto extends SessionWriteDto {
  id: string;
  createdAt: string;
  updatedAt: string;
}

export type RegistrationStatus = 'confirmed' | 'cancelled';

export interface RegistrationWriteDto {
  playerId: string;
  status: RegistrationStatus;
}

export interface RegistrationStatusUpdateDto {
  status: RegistrationStatus;
}

export interface RegistrationReadDto {
  id: string;
  playerId: string;
  sessionId: string;
  status: RegistrationStatus;
  registeredAt: string;
  updatedAt: string;
}

export interface ReservationSummaryPlayerReadDto {
  playerId: string;
  playerName: string;
}

export interface ReservationSummaryReadDto {
  sessionId: string;
  sessionDate: string;
  courseName: string;
  courseAddress: string;
  registrationDeadline: string;
  sessionStatus: SessionStatus;
  confirmedPlayerCount: number;
  estimatedGroups: number;
  summaryText: string;
  confirmedPlayers: ReservationSummaryPlayerReadDto[];
}

export type AdminUserLinkState = 'all' | 'linked' | 'unlinked';
export type AdminUserRoleFilter = 'all' | 'manager' | 'player';

export interface AdminUserReadDto {
  userId: string;
  displayName: string;
  provider: string;
  subject: string;
  role: 'manager' | 'player';
  playerId?: string;
  createdAt: string;
  updatedAt: string;
}

export interface AdminUserUpdateDto {
  role?: 'manager' | 'player';
  playerId?: string | null;
}
