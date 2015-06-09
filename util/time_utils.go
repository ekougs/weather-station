package util

import (
	"fmt"
	"regexp"
	"time"
)

// TimeUtils provides time utilities functions for this application
type TimeUtils struct {
	// INHERITANCE BY COMPOSITION
	DataUtils
}

// TimeFormat is the default time format
const TimeFormat = time.RFC3339

var timeUtilsNil = TimeUtils{}
var withoutTimeZoneRegexp *regexp.Regexp
var withTimeZoneRegexp *regexp.Regexp
var zeroValueTime = time.Time{}

// NewTimeUtils is the constructor for TimeUtils
func NewTimeUtils(dataUtils DataUtils) (TimeUtils, error) {
	var error error
	withoutTimeZoneRegexp, error = regexp.Compile("\\d{4}-\\d{2}-\\d{2}T\\d{2}\\:\\d{2}\\:\\d{2}")
	if error != nil {
		return timeUtilsNil, error
	}
	withTimeZoneRegexp, error = regexp.Compile("\\d{4}-\\d{2}-\\d{2}T\\d{2}\\:\\d{2}\\:\\d{2}Z|\\+|\\-\\d{2}\\:\\d{2}")
	if error != nil {
		return timeUtilsNil, error
	}
	return TimeUtils{dataUtils}, nil
}

// GetTime provides the time (with or without timezone) for a given city
// formattedTime ex :
func (utils TimeUtils) GetTime(formattedTime, cityStr string) (time.Time, error) {
	hasNoTimeZone, timeFormatError := hasNoTimeZone(formattedTime)
	if timeFormatError != nil {
		return zeroValueTime, timeFormatError
	}

	if hasNoTimeZone {
		formattedTime += "Z"
	}
	location, error := utils.getLocation(cityStr)
	if error != nil {
		return zeroValueTime, error
	}
	completeTime, error := time.ParseInLocation(TimeFormat, formattedTime, location)
	if error != nil {
		return zeroValueTime, error
	}
	return getTimeWithoutMinuteSecondNano(completeTime, location), nil
}

func hasNoTimeZone(formattedTime string) (bool, error) {
	if isFoundOnce(withTimeZoneRegexp, formattedTime) {
		return false, nil
	}
	if isFoundOnce(withoutTimeZoneRegexp, formattedTime) {
		return true, nil
	}
	withoutTimeExample := withoutTimeZoneRegexp.FindString(time.RFC3339)
	return false, fmt.Errorf("Provided time should be like %s or %s", withoutTimeExample, TimeFormat)
}

func isFoundOnce(curRegexp *regexp.Regexp, stringToMatch string) bool {
	return len(curRegexp.FindAllIndex([]byte(stringToMatch), -1)) == 1
}

// GetTimeWithoutMinuteSecondNano returns the time you provide without minutes, seconds or nanos
func (TimeUtils) GetTimeWithoutMinuteSecondNano(input time.Time) time.Time {
	return getTimeWithoutMinuteSecondNano(input, input.Location())
}

func getTimeWithoutMinuteSecondNano(input time.Time, location *time.Location) time.Time {
	input = time.Date(input.Year(), input.Month(), input.Day(), input.Hour(), 0, 0, 0, location)
	return input
}

func (utils TimeUtils) getLocation(cityStr string) (*time.Location, error) {
	ianaTimezone, error := utils.getIANATimezone(cityStr)
	if error != nil {
		return nil, error
	}
	return time.LoadLocation(ianaTimezone)
}

func (utils TimeUtils) getIANATimezone(cityStr string) (string, error) {
	cities, error := utils.getCitiesData()
	if error != nil {
		return "", error
	}

	for _, city := range cities {
		if city.Name == cityStr || city.Code == cityStr {
			return city.IanaTZ, nil
		}
	}

	return "", fmt.Errorf("No data for city '%s'.", cityStr)
}
