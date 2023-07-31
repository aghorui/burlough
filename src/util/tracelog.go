//go:build !RELEASE

package util

import "log"

func LogTrace(args ...any) {
	log.Println(args...)
}

func LogTracef(format string, args ...any) {
	log.Printf(format, args...)
}