package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ekougs/weather-station/cli"
	"github.com/ekougs/weather-station/server"
	"github.com/ekougs/weather-station/util"
)

const defaultCity = "DKR"

// THE PROGRAM ENTRY

func main() {
	applicationPath := os.Args[0]
	applicationDir := path.Dir(applicationPath)
	var err error
	var timeUtils util.TimeUtils
	var dataUtils util.DataUtils
	var tempProvider util.TempProvider

	dataUtils, err = util.NewDataUtils(applicationDir + "/resources/cities.json")
	cli.IfErrorInformAndLeave(err)

	timeUtils, err = util.NewTimeUtils(dataUtils)
	cli.IfErrorInformAndLeave(err)

	tempProvider, err = util.NewTempProvider(dataUtils)
	cli.IfErrorInformAndLeave(err)

	city, formattedDate, duration, serverMode := initFlags(timeUtils, dataUtils)

	if *serverMode {
		weatherServer := server.NewWeatherServer(tempProvider, timeUtils, dataUtils)
		weatherServer.LaunchServer()
		os.Exit(0)
	}

	var date time.Time

	date, err = timeUtils.GetTime(*formattedDate, *city)
	// Error handling is important
	// A method often returns as last return value an error
	cli.IfErrorInformAndLeave(err)

	*duration = strings.TrimSpace(*duration)
	cliInstrHandler := cli.NewCLIInstructionHandler(*city, date, *duration)
	cliInstrHandler.PrintResponse(tempProvider, timeUtils)
}

// FLAGS INFORMATION AND RETRIEVAL

func initFlags(utils util.TimeUtils, dataUtils util.DataUtils) (city, formattedDate, duration *string, serverMode *bool) {
	// Help making --help and retrieve flags
	cities, err := dataUtils.GetCities()
	if err != nil {
		panic(err)
	}
	citiesToString := make([]string, 0, len(cities))
	for _, city := range cities {
		citiesToString = append(citiesToString, city.String())
	}

	citiesHelpMessage := fmt.Sprintf("IATA code or name for city in : %s", strings.Join(citiesToString, ", "))
	city = flag.String("c", defaultCity, citiesHelpMessage)

	dateExample := "Date format example : 2006-01-02T15:04:05"
	currentHour := utils.GetTimeWithoutMinuteSecondNano(time.Now())
	formattedDate = flag.String("d", currentHour.Format(util.TimeFormat), dateExample)

	durationHelpMessage := "Expecting duration like 1Y3M2D, 1Y2M, 3M2D or 3D"
	duration = flag.String("D", "", durationHelpMessage)

	serverMode = flag.Bool("s", false, "Launch HTTP server on port 1987 and ignore other flags")

	flag.Parse()

	return city, formattedDate, duration, serverMode
}
