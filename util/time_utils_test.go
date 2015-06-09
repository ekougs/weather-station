package util

import (
	"testing"
	"time"
)

var timeDataUtils, _ = NewDataUtils("../resources/cities.json")
var utils, _ = NewTimeUtils(timeDataUtils)

func TestGetTimeForNYCValidBecauseUsingKnownIataCode(t *testing.T) {
	nycTime, error := utils.GetTime("2015-04-02T17:00:00", "NYC")
	if error != nil {
		t.Errorf("Should not have an error '%s'", error)
	}
	expectedFormattedNYCTime := "2015-04-02T17:00:00-04:00"
	actualNYCTime := nycTime.Format(time.RFC3339)
	if expectedFormattedNYCTime != actualNYCTime {
		t.Errorf("Expected time '%s' is different from the actual one '%s'", expectedFormattedNYCTime,
			actualNYCTime)
	}
}

func TestGetTimeForDakarBecauseUsingKnownCityName(t *testing.T) {
	dakarTime, error := utils.GetTime("2015-04-02T17:00:00", "Dakar")
	if error != nil {
		t.Errorf("Should not have an error '%s'", error)
	}
	expectedFormattedDakarTime := "2015-04-02T17:00:00Z"
	actualDakarTime := dakarTime.Format(time.RFC3339)
	if expectedFormattedDakarTime != actualDakarTime {
		t.Errorf("Expected time '%s' is different from the actual one '%s'", expectedFormattedDakarTime,
			actualDakarTime)
	}
}

func TestGetTimeForDummyShouldFailBecauseUsingUnknownCityNameOrIataCode(t *testing.T) {
	execution := func() (interface{}, error) { return utils.GetTime("2015-04-02T17:00:00", "Dummy") }
	assertError(execution, "Should have an error for 'Dummy' city name", t)
}

func TestGetTimeShouldReturnErrorIfFileDoesNotExist(t *testing.T) {
	execution := func() (interface{}, error) { return NewDataUtils("nimportequoi") }
	assertError(execution, "Should have an error when creating component", t)
}

func TestWithDefaultFormattedDate(t *testing.T) {
	formattedTime := time.Date(2015, 4, 8, 16, 0, 0, 0, time.Local).Format(TimeFormat)
	dakarTime, error := utils.GetTime(formattedTime, "DKR")
	if error != nil {
		t.Errorf("Should not have an error '%s'", error)
	}
	expectedFormattedDakarTime := "2015-04-08T16:00:00Z"
	actualDakarTime := dakarTime.Format(TimeFormat)
	if actualDakarTime != expectedFormattedDakarTime {
		t.Errorf("Expected time '%s' is different from the actual one '%s'", expectedFormattedDakarTime,
			actualDakarTime)
	}
}

func TestWithBadDateFormatShouldFail(t *testing.T) {
	execution := func() (interface{}, error) { return utils.GetTime("nimportequoi", "DKR") }
	assertError(execution, "Should have an error when getting 'nimportequoi' time", t)
}

func TestWithDoubleDateShouldFail(t *testing.T) {
	execution := func() (interface{}, error) { return utils.GetTime("2015-04-02T17:00:002015-04-02T17:00:00", "DKR") }
	assertError(execution, "Should have an error when getting double date time", t)
}

func TestGetEndTimeWithWholeDuration(t *testing.T) {
	startTime, err := utils.getStartTime(time.Date(2015, 3, 4, 16, 0, 0, 0, time.Local), "1Y2M3D")
	if err != nil {
		t.Errorf("Should not have an error '%s'", err)
	}
	actualStartTime := startTime.Format(TimeFormat)
	expectedStartTime := "2014-01-01T16:00:00+01:00"
	if expectedStartTime != actualStartTime {
		t.Errorf("Actual start time %s should be equal to %s", actualStartTime, expectedStartTime)
	}
}

func TestGetEndTimeWithPartDuration(t *testing.T) {
	startTime, err := utils.getStartTime(time.Date(2015, 3, 4, 16, 0, 0, 0, time.Local), "2M3D")
	if err != nil {
		t.Errorf("Should not have an error '%s'", err)
	}
	actualStartTime := startTime.Format(TimeFormat)
	expectedStartTime := "2015-01-01T16:00:00+01:00"
	if expectedStartTime != actualStartTime {
		t.Errorf("Actual start time %s should be equal to %s", actualStartTime, expectedStartTime)
	}
}

func TestGetEndTimeWithWrongDuration(t *testing.T) {
	execution := func() (interface{}, error) {
		return utils.getStartTime(time.Date(2015, 3, 4, 16, 0, 0, 0, time.Local), "2M1Y3D")
	}
	assertError(execution, "Should fail because of wrong format", t)
}

func TestGetDatesFor3DaysPeriodAtBeginningOfMonth(t *testing.T) {
	datesChan, err := utils.GetDatesForPeriod(time.Date(2015, 3, 2, 16, 0, 0, 0, time.Local), "3D")
	if err != nil {
		t.Errorf("Should not have an error '%s'", err)
	}
	expectedDates := make([]time.Time, 0, 4)
	expectedDates = append(expectedDates, time.Date(2015, 2, 27, 16, 0, 0, 0, time.Local))
	expectedDates = append(expectedDates, time.Date(2015, 2, 28, 16, 0, 0, 0, time.Local))
	expectedDates = append(expectedDates, time.Date(2015, 3, 1, 16, 0, 0, 0, time.Local))
	expectedDates = append(expectedDates, time.Date(2015, 3, 2, 16, 0, 0, 0, time.Local))
	actualDates := make([]time.Time, 0, 4)

	for {
		dates, open := <-datesChan
		if !open {
			break
		}
		actualDates = append(actualDates, dates...)
	}
	assertDatesEqual(actualDates, expectedDates, t)
}

func TestGetDatesFor5DaysPeriod(t *testing.T) {
	datesChan, err := utils.GetDatesForPeriod(time.Date(2015, 3, 12, 16, 0, 0, 0, time.Local), "5D")
	if err != nil {
		t.Errorf("Should not have an error '%s'", err)
	}
	expectedDates := make([]time.Time, 0, 6)
	expectedDates = append(expectedDates, time.Date(2015, 3, 7, 16, 0, 0, 0, time.Local))
	expectedDates = append(expectedDates, time.Date(2015, 3, 8, 16, 0, 0, 0, time.Local))
	expectedDates = append(expectedDates, time.Date(2015, 3, 9, 16, 0, 0, 0, time.Local))
	expectedDates = append(expectedDates, time.Date(2015, 3, 10, 16, 0, 0, 0, time.Local))
	expectedDates = append(expectedDates, time.Date(2015, 3, 11, 16, 0, 0, 0, time.Local))
	expectedDates = append(expectedDates, time.Date(2015, 3, 12, 16, 0, 0, 0, time.Local))
	actualDates := make([]time.Time, 0, 6)

	for {
		dates, open := <-datesChan
		if !open {
			break
		}
		actualDates = append(actualDates, dates...)
	}
	assertDatesEqual(actualDates, expectedDates, t)
}

func assertError(execution func() (interface{}, error), noErrorFoundMessage string, t *testing.T) {
	_, error := execution()
	if error == nil {
		t.Errorf(noErrorFoundMessage)
	}
}

func assertDatesEqual(actual, expected []time.Time, t *testing.T) {
	allContained, missingElements := containsAll(actual, expected)
	if !allContained {
		t.Errorf("These expected elements %s are missing in actual %s", missingElements, actual)
		return
	}
	allContained, missingElements = containsAll(expected, actual)
	if !allContained {
		t.Errorf("These actual elements %s are missing in expected %s", missingElements, expected)
	}
}

func containsAll(container, contained []time.Time) (bool, []time.Time) {
	containsAll := true
	var missingElements []time.Time
external:
	for _, date1 := range contained {
		for _, date2 := range container {
			if date1.Equal(date2) {
				continue external
			}
		}
		containsAll = false
		missingElements = append(missingElements, date1)
	}
	return containsAll, missingElements
}
