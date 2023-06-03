package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

var lineRE = regexp.MustCompile(`(?m)^`)

func Indent(s, indent string) string {
	if strings.TrimSpace(s) == "" {
		return s
	}
	return lineRE.ReplaceAllLiteralString(s, indent)
}

// Success prints success message in stdout.
func Success(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, fmt.Sprintf("\n\u001B[0;32m✓\u001B[0m %s\n", msg), args...)
}

// Warn prints warning message in stderr.
func Warn(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("\u001B[0;33m%s\u001B[0m\n", msg), args...)
}

// Fail prints failure message in stderr.
func Fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("\u001B[0;31m✗\u001B[0m %s\n", msg), args...)
}

// Failed prints failure message in stderr and exits.
func Failed(msg string, args ...interface{}) {
	Fail(msg, args...)
	os.Exit(1)
}

// ExitIfError exists with error message if err is not nil.
func ExitIfError(err error) {
	if err == nil {
		return
	}

	msg := fmt.Sprintf("Error: %s", err.Error())

	fmt.Fprintf(os.Stderr, "%s\n", msg)
	os.Exit(1)
}

// Format a date to string, return empty if isZero.
func FormatDate(d time.Time) string {
	if d.IsZero() {
		return ""
	}
	return d.Format(time.RFC3339)
}
