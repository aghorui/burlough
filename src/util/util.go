package util

import (
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var LogLevel = 0

// Regex for creating a URL/filename-safe or "wise" string.
// Does NOT replace whitespace. Use string.Fields and then join using a
// delimiter as needed.
// Initialized through a lambda function.
var SanitizeRegex *regexp.Regexp = func() *regexp.Regexp {
	re, err := regexp.Compile("[\\s!@#$%&()+=`~\\*\\^\\{\\}\\[\\]\\(\\)\\:\\;\\<\\>\\,\\?\\|\\/\\\\]")
	if err != nil {
		panic(err)
	}
	return re
}()

// Convenience function for stripping the file extension from a file
func ExtractFilename(s string) string {
	return strings.TrimSuffix(s, filepath.Ext(s))
}

// Splits a comma separated list into a string slice.
func SplitCommaList(s string) []string {
	if s == "" {
		return nil
	}

	list := strings.Split(s, ",")
	ret := make([]string, 0, len(list))

	for i := range list {
		s := strings.TrimSpace(list[i])
		if s == "" {
			continue
		}
		ret = append(ret, s)
	}

	return ret
}

// Standard timestamp.
func GetStandardTimestampString(t time.Time) string {
	return t.Format("02 January 2006")
}
