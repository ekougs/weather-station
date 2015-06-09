package util

import (
	"testing"
	"time"
)

var utils, _ = New("../resources/cities.json")

func TestGetTimeForNYCValidBecauseUsingKnownIataCode(t *testing.T) {
	nyc_time, error := utils.GetTime("2015-04-02T17:00:00", "NYC")
	if error != nil {
		t.Errorf("Should not have an error '%s'", error)
	}
	expected_formatted_nyc_time := "2015-04-02T17:00:00-04:00"
	actual_nyc_time := nyc_time.Format(time.RFC3339)
	if expected_formatted_nyc_time != actual_nyc_time {
		t.Errorf("Expected time '%s' is different from the actual one '%s'", expected_formatted_nyc_time,
			actual_nyc_time)
	}
}

func TestGetTimeForDakarBecauseUsingKnownCityName(t *testing.T) {
	dakar_time, error := utils.GetTime("2015-04-02T17:00:00", "Dakar")
	if error != nil {
		t.Errorf("Should not have an error '%s'", error)
	}
	expected_formatted_dakar_time := "2015-04-02T17:00:00Z"
	actual_nyc_time := dakar_time.Format(time.RFC3339)
	if expected_formatted_dakar_time != actual_nyc_time {
		t.Errorf("Expected time '%s' is different from the actual one '%s'", expected_formatted_dakar_time,
			actual_nyc_time)
	}
}

func TestGetTimeForDummyShouldFailBecauseUsingUnknownCityNameOrIataCode(t *testing.T) {
	execution := func() (interface{}, error) { return utils.GetTime("2015-04-02T17:00:00", "Dummy") }
	assert_error(execution, "Should have an error for 'Dummy' city name", t)
}

func TestGetTimeShouldReturnErrorIfFileDoesNotExist(t *testing.T) {
	execution := func() (interface{}, error) { return New("nimportequoi") }
	assert_error(execution, "Should have an error when creating component", t)
}

func TestWithDefaultFormattedDate(t *testing.T) {
	formatted_time := time.Date(2015, 4, 8, 16, 0, 0, 0, time.Local).Format(TIME_FORMAT)
	dakar_time, error := utils.GetTime(formatted_time, "DKR")
	if error != nil {
		t.Errorf("Should not have an error '%s'", error)
	}
	expected_formatted_dakar_time := "2015-04-08T16:00:00Z"
	actual_dakar_time := dakar_time.Format(TIME_FORMAT)
	if actual_dakar_time != expected_formatted_dakar_time {
		t.Errorf("Expected time '%s' is different from the actual one '%s'", expected_formatted_dakar_time,
			actual_dakar_time)
	}
}

func TestWithBadDateFormatShouldFail(t *testing.T) {
	execution := func() (interface{}, error) { return utils.GetTime("nimportequoi", "DKR") }
	assert_error(execution, "Should have an error when getting 'nimportequoi' time", t)
}

func TestWithDoubleDateShouldFail(t *testing.T) {
	execution := func() (interface{}, error) { return utils.GetTime("2015-04-02T17:00:002015-04-02T17:00:00", "DKR") }
	assert_error(execution, "Should have an error when getting double date time", t)
}

func assert_error(execution func() (interface{}, error), no_error_found_message string, t *testing.T) {
	_, error := execution()
	if error == nil {
		t.Errorf(no_error_found_message)
	}
}


