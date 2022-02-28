package localutils

import (
	"os"
	"strconv"
)

func DetectCI() bool {
	ciString := os.Getenv("CI")
	ci, err := strconv.ParseBool(ciString)
	if err != nil {
		return false
	}
	if ci {
		return true
	}
	return false
}

func StringPtrValue(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}
