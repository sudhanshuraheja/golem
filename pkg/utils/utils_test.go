package utils

import (
	"testing"
)

func TestUtils(t *testing.T) {
	u := GetUUID()
	Equals(t, true, u != "")
	Equals(t, true, IsValidUUID(u))

	u2 := GetShortUUID()
	Equals(t, true, u2 != "")

	str := " a-b - c-d - e-f "
	s := SplitAndTrim(str, "-")
	Equals(t, 0, ArrayContains(s, "a", true))
	Equals(t, 1, ArrayContains(s, "b", true))
	Equals(t, 2, ArrayContains(s, "c", true))
	Equals(t, 3, ArrayContains(s, "d", true))
	Equals(t, 4, ArrayContains(s, "e", true))
	Equals(t, 5, ArrayContains(s, "f", true))

	s = ArrayDelete(s, 0)
	Equals(t, -1, ArrayContains(s, "a", true))

	s = ArrayDelete(s, 3)
	Equals(t, 0, ArrayContains(s, "b", true))
	Equals(t, -1, ArrayContains(s, "e", true))

	f := 131.49897
	i := GetInt64FromFloat64(f, 2)
	Equals(t, int64(13150), i)
}

func TestGetFileContentType(t *testing.T) {
	file := "testing.go"
	contentType, err := GetFileContentType(file)
	OK(t, err)
	Equals(t, "txt", contentType)

	// file = "../testdata/test.zip"
	// contentType, err = GetFileContentType(file)
	// OK(t, err)
	// Equals(t, "zip", contentType)

	// file = "../testdata/account.tar.gz"
	// contentType, err = GetFileContentType(file)
	// OK(t, err)
	// Equals(t, "gz", contentType)
}
