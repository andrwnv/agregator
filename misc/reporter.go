package misc

import "fmt"

func ReportCritical(msg string) {
	panic(fmt.Sprintf("[CRIT]: %s.", msg))
}

func ReportError(msg string) {
	fmt.Println(fmt.Sprintf("[ERR]: %s.", msg))
}

func ReportWarning(msg string) {
	fmt.Println(fmt.Sprintf("[WARN]: %s.", msg))
}

func ReportInfo(msg string) {
	fmt.Println(fmt.Sprintf("[INFO]: %s.", msg))
}
