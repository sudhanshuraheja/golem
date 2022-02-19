package utils

import (
	"fmt"
	"io"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/sudhanshuraheja/golem/pkg/config"
	"github.com/sudhanshuraheja/golem/pkg/logger"
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

// GetConfig for testing
func GetConfig(logLevel int) (*config.Conf, *logger.Logger) {
	conf := config.NewConfig()
	log := logger.NewLogger(conf.LogLevel, true)
	return conf, log
}

// GetEchoContext ...
func GetEchoContext(method, url string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, url, body)
	response := httptest.NewRecorder()
	c := e.NewContext(req, response)
	return c, response
}
