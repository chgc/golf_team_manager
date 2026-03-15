package service

import "errors"

var (
	ErrLastManagerDemotionForbidden  = errors.New("cannot demote the last manager")
	ErrPlayerAlreadyLinked           = errors.New("player is already linked to another user")
	ErrPlayerNotFound                = errors.New("player was not found")
	ErrPlayerInactive                = errors.New("player is inactive")
	ErrReservationSummaryEmpty       = errors.New("reservation summary has no confirmed players")
	ErrReservationSummaryNotEligible = errors.New("session is not eligible for reservation summary")
	ErrSessionCapacityFull           = errors.New("session capacity is full")
	ErrSessionNotOpen                = errors.New("session is not open")
	ErrSessionReportNotFound         = errors.New("session for report was not found")
	ErrUserNotFound                  = errors.New("user was not found")
)
