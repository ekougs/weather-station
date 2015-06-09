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

func assertError(execution func() (interface{}, error), noErrorFoundMessage string, t *testing.T) {
	_, error := execution()
	if error == nil {
		t.Errorf(noErrorFoundMessage)
	}
}
