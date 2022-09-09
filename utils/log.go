package utils

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
)

func logFatal(err error, funcName string) {
	log.Printf("%v\n\n", err)
	log.Println("Stack trace:")

	stackLines := strings.Split(string(debug.Stack()), "\n")
	for iLine, line := range stackLines {
		if strings.Contains(line, funcName) {
			stackLines = stackLines[iLine+2:]
			break
		}
	}

	for _, line := range stackLines {
		log.Println(line)
	}
	os.Exit(1)
}

func logError(err error, funcName string) {
	if err != nil {
		logFatal(err, funcName)
	}
}

// LogError prints error and exits if error is not nil
func LogError(err error) {
	logError(err, "utils.LogError")
}

// LogFatalf prints a formatted error and exits
func LogFatalf(format string, args ...any) {
	logFatal(fmt.Errorf(format, args...), "utils.LogFatalf")
}

// Try prints error and exits if error is not nil, or return the original value otherwise
func Try[T any](out T, err error) T {
	logError(err, "utils.Try")
	return out
}

// Try2 prints error and exits if error is not nil, or return the 2 original values otherwise
func Try2[T, U any](out1 T, out2 U, err error) (T, U) {
	logError(err, "utils.Try2")
	return out1, out2
}
