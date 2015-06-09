package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ekougs/weather-station/util"
)

const defaultCity = "NYC"

// THE PROGRAM ITSELF

func main() {
	applicationPath := os.Args[0]
	applicationDir := path.Dir(applicationPath)
	var err error
	var timeUtils util.TimeUtils
	var dataUtils util.DataUtils

	dataUtils, err = util.NewDataUtils(applicationDir + "/resources/cities.json")
	if err != nil {
		log.Fatal(err)
	}
	timeUtils, err = util.NewTimeUtils(dataUtils)
	if err != nil {
		log.Fatal(err)
	}

	city, formattedDate, duration := initFlags(timeUtils, dataUtils)

	var date time.Time

	date, err = timeUtils.GetTime(*formattedDate, *city)
	// Error handling is important
	// A method often returns as last return value an error
	ifErrorInformAndLeave(err)

	var tempProvider util.TempProvider
	tempProvider, err = util.NewTempProvider(dataUtils)
	ifErrorInformAndLeave(err)

	*duration = strings.TrimSpace(*duration)
	if "" == *duration {
		var temp int
		temp, err = tempProvider.Get(*city, date)
		ifErrorInformAndLeave(err)
		cityTemp := util.CityTemp{City: *city, Time: util.TempTime(date), Temp: temp}

		var cityTempJSON []byte
		cityTempJSON, err = json.MarshalIndent(cityTemp, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(cityTempJSON))
	} else {
		// Duration has been provided
		datesChan, err := timeUtils.GetDatesForPeriod(date, *duration)
		ifErrorInformAndLeave(err)
		cityTemps := tempProvider.GetForDates(*city, datesChan)

		var cityTempsJSON []byte
		cityTempsJSON, err = json.MarshalIndent(cityTemps, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(cityTempsJSON))
	}
}

// HANDLE ERRORS
func ifErrorInformAndLeave(err error) {
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}
}

// FLAGS INFORMATION AND RETRIEVAL

func initFlags(utils util.TimeUtils, dataUtils util.DataUtils) (city, formattedDate, duration *string) {
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

	flag.Parse()

	return city, formattedDate, duration
}
