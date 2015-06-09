package util

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// TimeUtils provides time utilities functions for this application
type TimeUtils struct {
	// INHERITANCE BY COMPOSITION
	DataUtils
}

type duration struct {
	year, month, day int
}

// TimeFormat is the default time format
const TimeFormat = time.RFC3339

var timeUtilsNil = TimeUtils{}
var withoutTimeZoneRegexp *regexp.Regexp
var withTimeZoneRegexp *regexp.Regexp
var durationRegexp *regexp.Regexp
var durationLengthRegexp *regexp.Regexp
var durationUnitRegexp *regexp.Regexp
var timeNil = time.Time{}

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
	durationRegexp, error = regexp.Compile("(\\d+[YMD])?(\\d+[MD])?(\\d+D)?")
	if error != nil {
		return timeUtilsNil, error
	}
	durationLengthRegexp, error = regexp.Compile("\\d+")
	if error != nil {
		return timeUtilsNil, error
	}
	durationUnitRegexp, error = regexp.Compile("[YMD]")
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
		return timeNil, timeFormatError
	}

	if hasNoTimeZone {
		formattedTime += "Z"
	}
	location, error := utils.getLocation(cityStr)
	if error != nil {
		return timeNil, error
	}
	completeTime, error := time.ParseInLocation(TimeFormat, formattedTime, location)
	if error != nil {
		return timeNil, error
	}
	return getTimeWithoutMinuteSecondNano(completeTime, location), nil
}

// GetDatesForPeriod return a channel receiving by blocks day by day dates from startTime to
// end of period. Expecting duration like 1Y3M2D, 1Y2M, 3M2D or 3D
func (utils TimeUtils) GetDatesForPeriod(endTime time.Time, duration string) (chan []time.Time, error) {
	startTime, err := utils.getStartTime(endTime, duration)
	if err != nil {
		return nil, err
	}

	datesChan := make(chan []time.Time)
	go utils.generateDates(startTime, endTime, datesChan)

	return datesChan, nil
}

func (utils TimeUtils) generateDates(startTime, endTime time.Time, datesChan chan []time.Time) {
	count := 0
	empty := true
	generatedDates := make([]time.Time, 0, 5)
	generatedDate := startTime

	defer func() {
		if !empty {
			datesChan <- generatedDates
		}
		close(datesChan)
	}()

	for generatedDate.Before(endTime) || generatedDate.Equal(endTime) {
		generatedDates = append(generatedDates, generatedDate)
		generatedDate = generatedDate.AddDate(0, 0, 1)
		count++
		empty = false
		if count%5 == 0 {
			datesChan <- generatedDates
			generatedDates = make([]time.Time, 0, 5)
			empty = true
		}
	}
}

func (utils TimeUtils) getStartTime(endTime time.Time, durationString string) (time.Time, error) {
	if !isFoundOnce(durationRegexp, durationString) {
		return timeNil, fmt.Errorf("%s has not the right duration format. Expecting duration like 1Y3M2D, 1Y2M, 3M2D or 3D", durationString)
	}
	duration := duration{}
	durations := durationRegexp.FindAllStringSubmatch(durationString, -1)[0][1:]
	for _, durationString := range durations {
		if "" == durationString {
			// Every capturing group has been handled
			break
		}
		lengthString := durationLengthRegexp.FindStringSubmatch(durationString)[0]
		length, _ := parseToInt(lengthString)
		unit := durationUnitRegexp.FindStringSubmatch(durationString)[0]
		switch unit {
		case "Y":
			duration.year = length
		case "M":
			duration.month = length
		case "D":
			duration.day = length
		}
	}

	return endTime.AddDate(-duration.year, -duration.month, -duration.day), nil
}

func parseToInt(numString string) (int, error) {
	num, err := strconv.ParseInt(numString, 10, 32)
	return int(num), err
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
	return nbOfHits(curRegexp, stringToMatch) == 1
}

func isFoundAtLeastOnce(curRegexp *regexp.Regexp, stringToMatch string) bool {
	return nbOfHits(curRegexp, stringToMatch) >= 1
}

func nbOfHits(curRegexp *regexp.Regexp, stringToMatch string) int {
	return len(curRegexp.FindAllIndex([]byte(stringToMatch), -1))
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
