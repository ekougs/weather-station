package util

import (
	"time"
	"os"
	"encoding/json"
	"fmt"
	"regexp"
)

type City struct {
	Name     string `json:"name"`
	IataCode string `json:"iata_code"`
	IanaTZ   string `json:"iana_timezone"`
}

type Cities []City

type TimeUtils struct {
	cities_file string
}

const TIME_FORMAT = time.RFC3339
var time_utils_nil_value = TimeUtils{}
var without_time_zone_regexp *regexp.Regexp
var with_time_zone_regexp *regexp.Regexp

func New(cities_file string) (TimeUtils, error) {
	var error error
	without_time_zone_regexp, error = regexp.Compile("\\d{4}-\\d{2}-\\d{2}T\\d{2}\\:\\d{2}\\:\\d{2}")
	if error != nil {
		return time_utils_nil_value, error
	}
	with_time_zone_regexp, error = regexp.Compile("\\d{4}-\\d{2}-\\d{2}T\\d{2}\\:\\d{2}\\:\\d{2}Z|\\+|\\-\\d{2}\\:\\d{2}")
	if error != nil {
		return time_utils_nil_value, error
	}
	if _, err := os.Stat(cities_file); os.IsNotExist(err) {
		return time_utils_nil_value, fmt.Errorf("No such file or directory: %s", cities_file)
	}
	return TimeUtils{cities_file}, nil;
}

var cities Cities
var zero_value_time = time.Time{}

func (utils TimeUtils) GetTime(formatted_time, city_str string) (time.Time, error) {
	has_no_time_zone, time_format_error := has_no_time_zone(formatted_time)
	if time_format_error != nil {
		return zero_value_time, time_format_error
	}

	if has_no_time_zone {
		formatted_time += "Z"
	}
	location, error := utils.get_location(city_str)
	if error != nil {
		return zero_value_time, error
	}
	complete_time, error := time.ParseInLocation(TIME_FORMAT, formatted_time, location)
	if error != nil {
		return zero_value_time, error
	}
	return get_time_without_minute_second_nano(complete_time, location), nil
}

func has_no_time_zone(formatted_time string) (bool, error) {
	if is_found_once(with_time_zone_regexp, formatted_time) {
		return false, nil
	}
	if is_found_once(without_time_zone_regexp, formatted_time) {
		return true, nil
	}
	without_time_example := without_time_zone_regexp.FindString(time.RFC3339)
	return false, fmt.Errorf("Provided time should be like %s or %s", without_time_example, TIME_FORMAT)
}

func is_found_once(cur_regexp *regexp.Regexp, string_to_match string) bool {
	return len(cur_regexp.FindAllIndex([]byte(string_to_match), -1)) == 1
}

func (TimeUtils) GetTimeWithoutMinuteSecondNano(input time.Time) time.Time {
	return get_time_without_minute_second_nano(input, input.Location());
}

func get_time_without_minute_second_nano(input time.Time, location *time.Location) time.Time {
	input = time.Date(input.Year(), input.Month(), input.Day(), input.Hour(), 0, 0, 0, location)
	return input
}

func (utils TimeUtils) get_location(city_str string) (*time.Location, error) {
	iana_timezone, error := utils.get_iana_timezone(city_str)
	if error != nil {
		return nil, error
	}
	return time.LoadLocation(iana_timezone)
}

func (utils TimeUtils) get_iana_timezone(city_str string) (string, error) {
	cities, error := utils.get_cities()
	if error != nil {
		return "", error
	}

	for _, city := range cities {
		if city.Name == city_str || city.IataCode == city_str {
			return city.IanaTZ, nil
		}
	}

	return "", fmt.Errorf("We have no data for city '%s'.", city_str)
}

func (utils TimeUtils) get_cities() (Cities, error) {
	if cities != nil {
		return cities, nil
	}

	cities_json_file_reader, error := os.Open(utils.cities_file)
	if error != nil {
		return nil, error
	}

	cities_json_decoder := json.NewDecoder(cities_json_file_reader)
	cities := Cities{}
	error = cities_json_decoder.Decode(&cities)
	if error != nil {
		return nil, error
	}

	return cities, nil
}
