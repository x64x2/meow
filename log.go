package main

import (
	"fmt"

	"github.com/x64x2/go-debug"
)

// very simple levelled logging (we don't need anything more complex).
func fatalf(s string, a ...any) { fmt.Fprintf(stderr, "FATAL: "+s+"\n", a...) }
func errorf(s string, a ...any) { fmt.Fprintf(stderr, "ERROR: "+s+"\n", a...) }
func warnf(s string, a ...any)  { fmt.Fprintf(stderr, "WARN: "+s+"\n", a...) }
func infof(s string, a ...any)  { fmt.Fprintf(stderr, "INFO: "+s+"\n", a...) }
func debugf(s string, a ...any) {
	if debug.DEBUG {
		fmt.Fprintf(stderr, "DEBUG: "+s+"\n", a...)
	}
}
