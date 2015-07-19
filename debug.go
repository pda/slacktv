package slacktv

import (
	"fmt"
	"os"
)

var debugEnabled bool

func init() {
	if os.Getenv("DEBUG") != "" {
		debugEnabled = true
	}
}

// dbg
func dbg(format string, args ...interface{}) {
	if debugEnabled {
		format = fmt.Sprint("[DEBUG] ", format, "\n")
		fmt.Fprintf(os.Stderr, format, args...)
	}
}
