package logger

import "fmt"

var (
	loggingEnabled = false
)

func Init(enabled bool) {
	loggingEnabled = enabled
}

func LogLn(a ...interface{}) {
	if loggingEnabled {
		fmt.Println(a...)
	}
}

func LogF(msg string, a ...interface{}) {
	if loggingEnabled {
		fmt.Printf(msg, a...)
	}
}
