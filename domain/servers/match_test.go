package servers

import (
	"testing"

	"github.com/betas-in/utils"
)

func TestMatcher(t *testing.T) {
	// String array
	list := []string{"one", "two", "three", "four", "five"}

	m := NewMatch("name", "contains", "one")
	found, err := m.CompareStringArray(list)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, true, found)

	m = NewMatch("name", "not-contains", "six")
	found, err = m.CompareStringArray(list)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, true, found)

	m = NewMatch("name", "=", "six")
	_, err = m.CompareStringArray(list)
	utils.Test().Contains(t, err.Error(), "operators are supported")

	// String
	str := "abcdefghijklmnopqrstuvwxyz"

	m = NewMatch("name", "=", "cdef")
	found, err = m.CompareString(str)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, false, found)

	m = NewMatch("name", "!=", "cdef")
	found, err = m.CompareString(str)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, true, found)

	m = NewMatch("name", "like", "cdef")
	found, err = m.CompareString(str)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, true, found)

	m = NewMatch("name", "contains", "cdef")
	_, err = m.CompareString(str)
	utils.Test().Contains(t, err.Error(), "operators are supported")

	// Int
	count := 50

	m = NewMatch("count", "=", "50")
	found, err = m.CompareInt(count)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, true, found)

	m = NewMatch("count", "!=", "49")
	found, err = m.CompareInt(count)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, true, found)

	m = NewMatch("count", ">", "49")
	found, err = m.CompareInt(count)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, true, found)

	m = NewMatch("count", ">=", "49")
	found, err = m.CompareInt(count)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, true, found)

	m = NewMatch("count", "<", "51")
	found, err = m.CompareInt(count)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, true, found)

	m = NewMatch("count", "<=", "51")
	found, err = m.CompareInt(count)
	utils.Test().Nil(t, err)
	utils.Test().Equals(t, true, found)

	m = NewMatch("count", "contains", "51")
	_, err = m.CompareInt(count)
	utils.Test().Contains(t, err.Error(), "operators are supported")
}
