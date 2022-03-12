package localutils

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/betas-in/getter"
	"github.com/betas-in/logger"
	"github.com/betas-in/utils"
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

func ArrayPtrValue(a *[]string) []string {
	if a == nil {
		return []string{}
	}
	return *a
}

func Question(clog *logger.CLILogger, where, format string, v ...interface{}) string {
	reader := bufio.NewReader(os.Stdin)

	clog.Announce(where).Msgf(format, v...)
	fmt.Printf("%s", logger.MagentaBold("               "))
	text, err := reader.ReadString('\n')
	if err != nil {
		if err.Error() == "EOF" {
			return err.Error()
		}
		clog.Error(where).Msgf("could not read input: %v", err)
		return ""
	}
	return text
}

func TransferRate(bytes int64, start time.Time) string {
	duration := time.Since(start)
	bytesMB := float64(bytes) / (1024 * 1024)
	durationSec := float64(duration/time.Millisecond) / 1000
	transferRate := bytesMB / durationSec
	return fmt.Sprintf(
		"%sMbps, %sMB, %ssec",
		logger.CyanBold("%.3f", transferRate),
		logger.CyanBold("%.3f", bytesMB),
		logger.CyanBold("%.3f", durationSec),
	)
}

func TimeInSecs(start time.Time) string {
	duration := time.Since(start)
	durationSec := float64(duration/time.Millisecond) / 1000
	return fmt.Sprintf(
		"%ssec",
		logger.CyanBold("%.3f", durationSec),
	)
}

func FileCopy(text string) (string, error) {
	tempFilePath := utils.UUID().Get()
	file, err := os.CreateTemp("", tempFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, strings.NewReader(text))
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

func TinyString(text string, length int) string {
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, "\t", "")

	if len(text) <= length {
		return text
	}

	middle := "----"
	remaining := length - len(middle)

	split := 0
	if remaining%2 == 0 {
		split = remaining / 2
	} else {
		split = (remaining - 1) / 2
		middle = middle[:len(middle)-1]
	}

	start := text[:split]
	end := text[len(text)-split:]

	tiny := fmt.Sprintf("%s%s%s", start, logger.YellowBold(middle), end)
	return tiny
}

func Download(log *logger.CLILogger, name, url string) (string, error) {
	log.Info(name).Msgf("%s %s", logger.Cyan("downloading"), TinyString(url, 100))

	glog := logger.NewLogger(3, true)
	g := getter.NewGetter(glog)

	startTime := time.Now()
	response := g.FetchResponse(getter.Request{
		Path:       url,
		SaveToDisk: true,
	})

	if response.Error != nil {
		return "", response.Error
	}
	if response.Code != 200 {
		return "", fmt.Errorf("received error code for %s: %d", url, response.Code)
	}

	log.Info(name).Msgf(
		"%s %s",
		logger.GreenBold("downloaded!"),
		TinyString(url, 100),
	)

	log.Info(name).Msgf(
		"         %s %s %s",
		logger.GreenBold("to"),
		TinyString(response.DataPath, 100),
		TimeInSecs(startTime),
	)
	return response.DataPath, nil
}

func Base64EncodedRandomNumber(size int) (string, error) {
	randomNumber, err := utils.Crypto().RandomNumberGenerator(size)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(randomNumber), nil
}

func ArrayUnique(ar []string) []string {
	keys := map[string]bool{}
	for _, str := range ar {
		keys[str] = true
	}

	uniques := []string{}
	for k := range keys {
		uniques = append(uniques, k)
	}
	return uniques
}

func StrPtr(s string) *string {
	return &s
}

func Create(file string) (bool, error) {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		f, err := os.Create(file)
		if err != nil {
			return false, err
		}
		defer f.Close()
		return true, nil
	} else if err != nil {
		return false, err
	}
	return false, nil
}

func Touch(file string) (bool, error) {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		f, err := os.Create(file)
		if err != nil {
			return false, err
		}
		defer f.Close()
		return true, nil
	} else if err != nil {
		return false, err
	} else {
		now := time.Now().Local()
		err = os.Chtimes(file, now, now)
		if err != nil {
			return false, err
		}
	}
	return false, nil
}
