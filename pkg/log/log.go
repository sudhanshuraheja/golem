package log

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
)

var logLevel = 3

func SetLogLevel(level *string) {
	if level == nil {
		logLevel = 3
		return
	}

	switch *level {
	case "TRACE":
		logLevel = 6
	case "DEBUG":
		logLevel = 5
	case "INFO":
		logLevel = 4
	case "WARN":
		logLevel = 3
	case "ERROR":
		logLevel = 2
	case "FATAL":
		logLevel = 1
	case "PANIC":
		logLevel = 0
	}
}

func Tracef(format string, v ...interface{}) {
	if logLevel >= 6 {
		fmt.Printf(format+"\n", v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if logLevel >= 5 {
		fmt.Printf(format+"\n", v...)
	}
}

func Infof(format string, v ...interface{}) {
	if logLevel >= 4 {
		code := color.New(color.FgCyan)
		code.Printf(format+"\n", v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if logLevel >= 3 {
		code := color.New(color.FgYellow, color.Bold)
		code.Printf(format+"\n", v...)
	}
}

func MinorSuccessf(format string, v ...interface{}) {
	if logLevel >= 2 {
		code := color.New(color.FgGreen)
		code.Printf(format+"\n", v...)
	}
}

func Successf(format string, v ...interface{}) {
	if logLevel >= 2 {
		code := color.New(color.FgGreen, color.Bold)
		code.Printf(format+"\n", v...)
	}
}

func Announcef(format string, v ...interface{}) {
	if logLevel >= 2 {
		code := color.New(color.FgCyan, color.Bold)
		code.Printf(format+"\n", v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if logLevel >= 2 {
		code := color.New(color.FgRed, color.Bold)
		code.Printf(format+"\n", v...)
	}
}

func Fatalf(format string, v ...interface{}) {
	if logLevel >= 1 {
		code := color.New(color.FgRed, color.BgWhite)
		code.Printf(format+"\n", v...)
	}
}

func Panicf(format string, v ...interface{}) {
	if logLevel >= 0 {
		code := color.New(color.FgRed, color.BgWhite)
		code.Printf(format+"\n", v...)
	}
}

func Dump(v interface{}) {
	spew.Dump(v)
}
