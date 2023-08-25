package r2loxerrors

import (
	"fmt"
	"os"

	"github.com/arturoeanton/go-r2lox/globals"
)

func Errors(line int, message string) string {

	return Report(line, "", message)
}

func Report(line int, where, message string) string {
	report_str := fmt.Sprintf("[line %d] Error%s: %s\n", line, where, message)
	if globals.PrintFlag {
		fmt.Fprintf(os.Stderr, "%s", report_str)
	}
	globals.HasError = true
	return report_str
}
