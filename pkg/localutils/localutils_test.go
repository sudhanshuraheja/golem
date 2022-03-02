package localutils

import (
	"testing"

	"github.com/betas-in/utils"
)

func TestLocalUtils(t *testing.T) {
	str := StringPtrValue(nil, "hello")
	utils.Test().Equals(t, "hello", str)
}
