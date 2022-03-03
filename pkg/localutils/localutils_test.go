package localutils

import (
	"testing"

	"github.com/betas-in/utils"
)

func TestLocalUtils(t *testing.T) {
	str := StringPtrValue(nil, "hello")
	utils.Test().Equals(t, "hello", str)
}

func TestTinyString(t *testing.T) {
	str := "how to get the last x chars of a string"
	tiny := TinyString(str, 20)
	utils.Test().Equals(t, "how to g⋆⋆⋆⋆a string", tiny)

	str = `how to
	 get the 
	last x chars of
	 a string`
	tiny = TinyString(str, 20)
	utils.Test().Equals(t, "how to g⋆⋆⋆⋆a string", tiny)
}
