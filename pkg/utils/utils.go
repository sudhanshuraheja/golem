package utils

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sudhanshuraheja/golem/pkg/log"
)

// GetUUID definition
func GetUUID() string {
	return fmt.Sprintf("%v", uuid.New())
}

// GetShortUUID definition
func GetShortUUID() string {
	return strings.Split(GetUUID(), "-")[0]
}

// IsValidUUID definition
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// ArrayContains ...
func ArrayContains(s []string, e string, exactMatch bool) int {
	for k, a := range s {
		if !exactMatch {
			a = strings.ToLower(strings.TrimSpace(a))
			e = strings.ToLower(strings.TrimSpace(e))
		}
		if a == e {
			return k
		}
	}
	return -1
}

// ArrayDelete ...
func ArrayDelete(s []string, i int) []string {
	if i == -1 {
		return s
	}
	if i > len(s) {
		return s
	}

	return append(s[:i], s[i+1:]...)
}

// SplitAndTrim ...
func SplitAndTrim(s string, separator string) []string {
	output := []string{}
	splits := strings.Split(s, separator)
	for _, split := range splits {
		output = append(output, strings.TrimSpace(split))
	}
	return output
}

// GetInt64FromFloat64 ...
func GetInt64FromFloat64(f float64, decimals int) int64 {
	multiplier := math.Pow(10, float64(decimals))
	return int64(math.Round(f * multiplier))
}

// ElapsedTime prints the execution time
func ElapsedTime(what string, logging bool) func() {
	start := time.Now()
	return func() {
		if logging {
			log.Tracef("%v | %s\n", time.Since(start), what)
		}
	}
}

// GetFileContentType ..
func GetFileContentType(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)
	contentTypeString := GetContentTypeString([]string{contentType})
	return contentTypeString, nil
}

// GetContentTypeString ...
func GetContentTypeString(list []string) string {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
	// https://www.iana.org/assignments/media-types/media-types.xhtml
	types := map[string]string{
		"application/zip":              "zip",
		"application/x-zip-compressed": "zip",
		"application/x-rar-compressed": "rar",
		"application/gzip":             "gz",
		"application/x-gzip":           "gz",
		"application/json":             "json",
		"application/pdf":              "pdf",
		"text/html; charset=UTF-8":     "html",
		"text/plain; charset=utf-8":    "txt",
	}
	for key, val := range types {
		match := ArrayContains(list, key, false)
		if match == -1 {
			continue
		}
		return val
	}
	return strings.Join(list, ";;")
}

func StringPtrValue(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}
