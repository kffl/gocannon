package main

import "time"

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Nanosecond)
}
