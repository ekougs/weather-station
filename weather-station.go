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

type tempTime time.Time

type cityTemp struct {
	City string
	Time tempTime
	Temp int
}

const timeToStringFormat = time.RFC1123
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

	city, formattedDate := initFlags(timeUtils, dataUtils)

	var date time.Time

	date, err = timeUtils.GetTime(*formattedDate, *city)
	// Error handling is important
	// A method often returns as last return value an error
	ifErrorInformAndLeave(err)

	var tempProvider util.TempProvider
	tempProvider, err = util.NewTempProvider(dataUtils)
	ifErrorInformAndLeave(err)

	var temp int
	temp, err = tempProvider.Get(*city, date)
	ifErrorInformAndLeave(err)
	obj := cityTemp{*city, tempTime(date), temp}

	var cityTempJSON []byte
	cityTempJSON, err = json.MarshalIndent(obj, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(cityTempJSON))
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

func initFlags(utils util.TimeUtils, dataUtils util.DataUtils) (city, formattedDate *string) {
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

	flag.Parse()

	return city, formattedDate
}

// HOW TO FORMAT

func (obj cityTemp) String() string {
	formattedTime := obj.Time.Format(timeToStringFormat)
	return fmt.Sprintf("%s %s %d", obj.City, formattedTime, obj.Temp)
}

// Format provides custom format for time
func (ourTime tempTime) Format(format string) string {
	return time.Time(ourTime).Format(format)
}

// MarshalJSON provides custom JSON marshaller for time
func (ourTime tempTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + ourTime.Format(util.TimeFormat) + "\""), nil
}
