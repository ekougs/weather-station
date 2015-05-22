package main

import (
	"math/rand"
	"fmt"
	"flag"
	"time"
	"log"
	"encoding/json"
	"os"
)

type TempTime time.Time
type Temp float64

type city_temp struct {
	City string
	Time TempTime
	Temp Temp
}

const time_format = time.RFC3339
const time_to_string_format = time.RFC1123
const default_city = "DKR"

// THE PROGRAM ITSELF

func main() {
	city, formatted_date := init_flags()

	date, error := time.Parse(time_format, *formatted_date)
	// Error handling is important
	// A method often returns as last return value an error
	if error != nil {
		fmt.Printf("Date '%s' format is not recognized.\n", *formatted_date)
		flag.Usage()
		os.Exit(1)
	}

	temp := Temp(rand.Float64() * 5 + 20)
	obj := city_temp{*city, TempTime(date), temp}

	city_temp_json, error := json.MarshalIndent(obj, "", "    ")
	if error != nil {
		log.Fatal(error)
	}
	fmt.Println(string(city_temp_json))
}

// FLAGS INFORMATION AND RETRIEVAL

func init_flags() (city, provided_formatted_date *string) {
	// Help making --help and retrieve flags
	city = flag.String("c", default_city, "IATA code for city")

	date_example := "Date format example : " + time_format
	provided_formatted_date = flag.String("d", time_now_no_minute_second_nano().Format(time_format), date_example)

	flag.Parse()

	date_flag_provided := flag.Lookup("d") != nil
	if date_flag_provided {
		*provided_formatted_date = fmt.Sprintf("%s", *provided_formatted_date)
	}
	return city, provided_formatted_date
}

func time_now_no_minute_second_nano() time.Time {
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	return now
}

// HOW TO FORMAT

func (obj city_temp) String() string {
	formatted_time := obj.Time.Format(time_to_string_format)
	return fmt.Sprintf("%s %s %.1f", obj.City, formatted_time, obj.Temp)
}

func (our_time TempTime) Format(format string) string {
	return time.Time(our_time).Format(format)
}

func (our_time TempTime) MarshalJSON() ([]byte, error) {
	return []byte("\""+our_time.Format(time_format)+"\""), nil
}

func (temp Temp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.2f", temp)), nil
}
