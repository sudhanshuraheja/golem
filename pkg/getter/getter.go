package getter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/config"
	"github.com/sudhanshuraheja/golem/pkg/utils"
)

var (
	defaultTimeout = 10 * time.Second
	fakeUserAgent  = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 YaBrowser/21.11.0 Yowser/2.5 Safari/537.36"
)

// Error ...
type Error struct {
	ErrorString string
}

func (ge *Error) Error() string { return ge.ErrorString }

var (
	// ErrRateLimited ...
	ErrRateLimited = &Error{ErrorString: "rate limited"}
	// ErrRateLimitingDown ...
	ErrRateLimitingDown = &Error{ErrorString: "rate limiting is down"}
	// ErrRateLimitingInvalid ...
	ErrRateLimitingInvalid = &Error{ErrorString: "invalid internal data in rate limiting"}
)

// Getter definition
type Getter interface {
	SetDefaultTimeout(timeout time.Duration)
	SetUserAgent(useragent string)
	FetchResponse(r Request) *Response
}

// Download definition
type getter struct {
	conf      *config.Config
	timeout   time.Duration
	userAgent string
}

// NewGetter ...
func NewGetter(conf *config.Config) Getter {
	userAgent := fakeUserAgent
	timeout := defaultTimeout
	return &getter{conf: conf, timeout: timeout, userAgent: userAgent}
}

func (g *getter) SetDefaultTimeout(timeout time.Duration) {
	g.timeout = timeout
}

func (g *getter) SetUserAgent(useragent string) {
	g.userAgent = useragent
}

func (g *getter) FetchResponse(r Request) *Response {
	var request *http.Request
	var err error

	var response Response

	// Timeout
	timeout := g.timeout
	if r.Timeout > 0 {
		timeout = r.Timeout
	}
	client := &http.Client{Timeout: timeout}

	// Method
	if r.Method == "" {
		r.Method = http.MethodGet
	}

	// Path
	if len(r.Query) > 0 {
		r.Path = fmt.Sprintf("%s%s", r.Path, g.getStringFromQuery(r.Query))
	}

	// HTTP.Request
	var requestBody io.Reader
	if r.SendFormData {
		if len(r.FormDataMap) > 0 {
			r.FormData = g.getURLValuesFromMap(r.FormDataMap)
		}
		requestBody = strings.NewReader(r.FormData.Encode())
	} else if r.SendJSONData {
		if r.JSONInterface != nil {
			r.JSONData, err = g.getJSONBytesFromData(r.JSONInterface)
			if err != nil {
				response.Error = err
				return &response
			}
		}
		requestBody = bytes.NewReader(r.JSONData)
	}
	request, err = http.NewRequest(r.Method, r.Path, requestBody)
	if err != nil {
		response.Error = err
		return &response
	}

	// Headers
	if r.SendJSONData {
		request.Header.Set("Content-Type", "application/json")
	}
	request.Header.Set("User-Agent", g.userAgent)
	for key, value := range r.Headers {
		request.Header.Set(key, value)
	}

	// Make request
	data, err := client.Do(request)
	if err != nil {
		response.Error = err
		return &response
	}
	defer data.Body.Close()

	// Add response code
	response.Code = data.StatusCode
	if data.StatusCode >= 200 && data.StatusCode <= 300 {
		response.Error = nil
	} else {
		response.Error = fmt.Errorf("recieved %d error: %s", data.StatusCode, response.Data)
	}

	// GetContentType
	contentTypes, ok := data.Header["Content-Type"]
	if !ok {
		contentTypes = []string{}
	}
	contentType := utils.GetContentTypeString(contentTypes)
	response.ContentType = contentType

	if r.SaveToDisk || utils.ArrayContains([]string{"zip", "gzip", "pdf", "rar"}, contentType, true) != -1 {
		tempFilePath := utils.GetUUID()
		file, err := os.CreateTemp("", tempFilePath)
		if err != nil {
			logger.Errorf("getter | could not create file: %v", err)
			response.Error = err
			return &response
		}
		_, err = io.Copy(file, data.Body)
		if err != nil {
			logger.Errorf("getter | could not copy to file: %v", err)
			response.Error = err
			return &response
		}
		response.DataPath = file.Name()
		contentType, err := utils.GetFileContentType(file.Name())
		if err != nil {
			logger.Errorf("getter | could not check content type: %v", err)
			response.Error = err
			return &response
		}
		response.ContentType = contentType
	} else {
		// Read body
		body, err := ioutil.ReadAll(data.Body)
		if err != nil {
			response.Error = err
			return &response
		}
		response.Data = body
	}

	return &response
}

//
// Internal Functions
//

// get string from r.Query
func (g *getter) getStringFromQuery(query map[string]string) string {
	list := []string{}
	for key, value := range query {
		item := fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(value))
		list = append(list, item)
	}
	return fmt.Sprintf("?%s", strings.Join(list, "&"))
}

func (g *getter) getURLValuesFromMap(data map[string]string) url.Values {
	output := url.Values{}
	for key, value := range data {
		output.Set(key, value)
	}
	return output
}

func (g *getter) getJSONBytesFromData(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}
