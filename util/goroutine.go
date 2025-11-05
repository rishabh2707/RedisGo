package util

import (
	"bytes"
	"fmt"
	"runtime"
)

func GetGoroutineID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// Output looks like: "goroutine 18 [running]:"
	fields := bytes.Fields(buf[:n])
	var id uint64
	fmt.Sscanf(string(fields[1]), "%d", &id)
	return id
}
