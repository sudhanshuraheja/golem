package recipes

import (
	"fmt"
	"time"

	"github.com/betas-in/getter"
	"github.com/betas-in/logger"
	"github.com/sudhanshuraheja/golem/pkg/localutils"
)

func Download(log *logger.CLILogger, name, url string) (string, error) {
	log.Info(name).Msgf("%s %s", logger.Cyan("downloading"), localutils.TinyString(url, tiny))

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
		"%s %s %s %s %s",
		logger.GreenBold("downloaded"),
		localutils.TinyString(url, tiny),
		logger.GreenBold("to"),
		localutils.TinyString(response.DataPath, tiny),
		localutils.TimeInSecs(startTime),
	)
	return response.DataPath, nil
}
