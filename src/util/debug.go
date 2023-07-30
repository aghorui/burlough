////  go:build !RELEASE

package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// Get file, line number and function name of the place where the grandparent
// function is being called from.
// Equivalent to `snprintf(s, len, "%s:%d %s", __FILE__, __LINE__, __func__);`
func SrcStamp(level int) string {
	pc, file, line, ok := runtime.Caller(level)

	if !ok {
		return "???"
	}

	file = filepath.Base(file)

	funcName := runtime.FuncForPC(pc).Name()

	return fmt.Sprintf("%v:%v %v", file, line, funcName)
}

// Prints ...v to stdout. If v[i] is a struct, its field names will be printed
// out.
func DumpVar(vars ...any) {
	log.Printf("<debug> ")
	for _, v := range vars {
		fmt.Printf("[debug] %+v\n", v)
	}
	fmt.Println()
}

// Dumps a file to stderr.
func DumpFile(fpath string) {

	data, err := os.ReadFile(fpath)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "<debug file='%v' size=%v>\n", fpath, len(data))
	fmt.Fprintln(os.Stderr, string(data))
	fmt.Fprintln(os.Stderr, "</debug>")

}

// Generate an error with srcstamp.
func Error(err error) error {
	return fmt.Errorf("%v: %w; ", SrcStamp(2), err)
}

// Log the error. Convenience function.
func LogErr(err error) {
	log.Printf("Error: %v", fmt.Errorf("%v: %w; ", SrcStamp(2), err))
}