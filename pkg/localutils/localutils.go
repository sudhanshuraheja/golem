package localutils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

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

func Question(clog *logger.CLILogger, where, format string, v ...interface{}) string {
	reader := bufio.NewReader(os.Stdin)

	clog.Announce(where).Msgf(format, v...)
	fmt.Printf("%s", logger.MagentaBold("               "))
	text, err := reader.ReadString('\n')
	if err != nil {
		clog.Error(where).Msgf("could not read input: %v", err)
	}
	return text
}

func TransferRate(bytes int64, start time.Time) string {
	duration := time.Since(start)
	bytesMB := float64(float64(bytes) / (1024 * 1024))
	durationSec := float64(duration/time.Millisecond) / 1000
	transferRate := float64(bytesMB / durationSec)
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
