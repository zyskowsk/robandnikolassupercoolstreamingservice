package main

import "fmt"

var logLevelDebug = false

func LogDebug(format string, a ...interface{}) (int, error) {
	if logLevelDebug {
		return fmt.Printf(format, a)
	}
	return 0, nil
}

func LoglnDebug(line interface{}) (int, error) {
	if logLevelDebug {
		return fmt.Println(line)
	}
	return 0, nil
}

func Log(format string, a ...interface{}) (int, error) {
	return fmt.Printf(format, a)
}

func Logln(line interface{}) (int, error) {
	return fmt.Println(line)
}

func SetLogLevelDebug() {
	logLevelDebug = true
}
