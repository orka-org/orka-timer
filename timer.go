package main

import "time"

type TimeInterval struct {
	Start time.Time
	End   time.Time
}

type CreateTimerRequest struct {
	TimeInterval
	Pauses []TimeInterval
}
