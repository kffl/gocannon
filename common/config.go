package common

import "time"

// Config is a struct containing all of the parsed CLI flags and arguments
type Config struct {
	Duration    *time.Duration
	Connections *int
	Timeout     *time.Duration
	Mode        *string
	OutputFile  *string
	Interval    *time.Duration
	Preallocate *int
	Method      *string
	Body        *RawRequestBody
	Headers     *RequestHeaders
	TrustAll    *bool
	Format      *string
	Plugin      *string
	Target      *string
}
