package getter

import (
	"os"
	"testing"

	"github.com/sudhanshuraheja/golem/pkg/utils"
)

func TestGet(t *testing.T) {
	url := "http://feeds.arstechnica.com/arstechnica/index"

	g := NewGetter(nil)

	response := g.FetchResponse(Request{
		Path: url,
	})
	utils.OK(t, response.Error)
	utils.Equals(t, 200, response.Code)
	utils.Contains(t, "https://arstechnica.com/", string(response.Data))

	response = g.FetchResponse(Request{
		Path: "http://gotrixsterinmypajamas.com",
	})
	utils.Equals(t, 0, response.Code)
	utils.Contains(t, "dial tcp", response.Error.Error())
}

func TestZipped(t *testing.T) {
	url := "https://github.com/gojekfarm/async-worker/archive/refs/heads/master.zip"

	g := NewGetter(nil)
	response := g.FetchResponse(Request{
		Path: url,
	})
	utils.OK(t, response.Error)
	utils.Equals(t, 200, response.Code)
	utils.Equals(t, "zip", response.ContentType)
	os.Remove(response.DataPath)
}
