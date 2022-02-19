package utils

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/davecgh/go-spew/spew"
)

// OK fails the test if an err is not nil.
func OK(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("➤ %s:%d ➤ Error was expected to be nil \n", filepath.Base(file), line)
		fmt.Print("--- RECIEVED --- \n")
		spew.Dump(err)
		fmt.Print("---------------- \n")
		tb.FailNow()
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, expected, received interface{}) {
	if !reflect.DeepEqual(expected, received) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("➤ %s:%d ➤ DeepEqual failed\n", filepath.Base(file), line)
		fmt.Print(diff.LineDiff(spew.Sdump(expected), spew.Sdump(received)))
		fmt.Print("\n")
		tb.Fail()
	}
}

// NotEmptyString Equals fails the test if exp is not equal to act.
func NotEmptyString(tb testing.TB, received string) {
	if received == "" {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("➤ %s:%d ➤ String is empty\n", filepath.Base(file), line)
		fmt.Print("\n")
		tb.Fail()
	}
}

// Contains checks two strings
func Contains(tb testing.TB, expected, errorString string) {
	if !strings.Contains(errorString, expected) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("➤ %s:%d ➤ Error string is unexpected\n➤➤➤ Error: %s\n➤➤➤ ShouldHave: %s\n", filepath.Base(file), line, errorString, expected)
		tb.Fail()
	}
}
