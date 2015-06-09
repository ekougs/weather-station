package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/ekougs/weather-station/util"
)

type TempTime time.Time
type Temp float64

type CityTemp struct {
	City string
	Time TempTime
	Temp Temp
}

const time_to_string_format = time.RFC1123
const default_city = "NYC"

// THE PROGRAM ITSELF

func main() {
	application_path := os.Args[0]
	application_dir := path.Dir(application_path)
	var err error
	var utils util.TimeUtils
	utils, err = util.New(application_dir + "/resources/cities.json")
	if err != nil {
		log.Fatal(err)
	}

	city, formatted_date := init_flags(utils)

	var date time.Time

	date, err = utils.GetTime(*formatted_date, *city)
	// Error handling is important
	// A method often returns as last return value an error
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}

	temp := Temp(rand.Float64()*5 + 20)
	obj := CityTemp{*city, TempTime(date), temp}

	var city_temp_json []byte
	city_temp_json, err = json.MarshalIndent(obj, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(city_temp_json))
}

// FLAGS INFORMATION AND RETRIEVAL

func init_flags(utils util.TimeUtils) (city, formatted_date *string) {
	// Help making --help and retrieve flags
	city = flag.String("c", default_city, "IATA code for city")

	date_example := "Date format example : " + util.TIME_FORMAT
	current_hour := utils.GetTimeWithoutMinuteSecondNano(time.Now())
	formatted_date = flag.String("d", current_hour.Format(util.TIME_FORMAT), date_example)

	flag.Parse()

	return city, formatted_date
}

// HOW TO FORMAT

func (obj CityTemp) String() string {
	formatted_time := obj.Time.Format(time_to_string_format)
	return fmt.Sprintf("%s %s %.1f", obj.City, formatted_time, obj.Temp)
}

func (our_time TempTime) Format(format string) string {
	return time.Time(our_time).Format(format)
}

func (our_time TempTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + our_time.Format(util.TIME_FORMAT) + "\""), nil
}

func (temp Temp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.2f", temp)), nil
}
