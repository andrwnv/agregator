package utils

import "fmt"

func ReportCritical(msg string) {
	panic(fmt.Sprintf("[CRIT]: %s.", msg))
}

func ReportError(msg string) {
	panic(fmt.Sprintf("[ERR]: %s.", msg))
}

func ReportWarning(msg string) {
	panic(fmt.Sprintf("[WARN]: %s.", msg))
}

func ReportInfo(msg string) {
	panic(fmt.Sprintf("[INFO]: %s.", msg))
}
