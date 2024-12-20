package coati2lang

import (
	"fmt"
)

func Errors(line int, message string) string {

	return Report(line, "", message)
}

func Report(line int, where, message string) string {
	report_str := fmt.Sprintf("[line %d] Error%s: %s\n", line, where, message)
	if PrintFlag {
		//fmt.Fprintf(os.Stderr, "%s", report_str)
		panic(report_str)
	}
	HasError = true
	return report_str
}
