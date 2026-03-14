package service

import "errors"

var (
	ErrPlayerInactive      = errors.New("player is inactive")
	ErrSessionCapacityFull = errors.New("session capacity is full")
	ErrSessionNotOpen      = errors.New("session is not open")
)
