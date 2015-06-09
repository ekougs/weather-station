package util

import (
	"os"
	"testing"
	"time"
)

var tempDataUtils, _ = NewDataUtils("../resources/cities.json")
var tempProvider, _ = NewTempProvider(tempDataUtils)

func TestMain(m *testing.M) {
	retCode := m.Run()

	tearDown()

	os.Exit(retCode)
}

func tearDown() {
	os.Remove("../resources/DKR.json")
	os.Remove("../resources/NYC.json")
	os.Remove("../resources/PAR.json")
}

func TestShouldReturnARealisticTempForDKR(t *testing.T) {
	location, _ := time.LoadLocation("Africa/Dakar")
	temp, error := tempProvider.generate("DKR", time.Date(2015, 4, 16, 13, 00, 00, 00, location))
	if error != nil {
		t.Error(error)
	}
	// 17 23
	if temp < 15 || temp > 26 {
		t.Errorf("Actual generated temp '%d' should be between 15 and 26", temp)
	}
}

func TestShouldReturnSameTempFor2CallsWithSameParameters(t *testing.T) {
	location, _ := time.LoadLocation("Africa/Dakar")
	var temp1, temp2 int
	var err error
	temp1, err = tempProvider.Get("DKR", time.Date(2015, 4, 21, 13, 00, 00, 00, location))
	if err != nil {
		t.Error(err)
	}
	temp2, err = tempProvider.Get("DKR", time.Date(2015, 4, 21, 13, 00, 00, 00, location))
	if err != nil {
		t.Error(err)
	}
	if temp1 != temp2 {
		t.Errorf("Value should be generated only once. First temp %d should be equal to second %d", temp1, temp2)
	}
}

func TestShouldReturnSameTempFor2CallsWithSameParametersBeforeDayLightSavingTimeIncrement(t *testing.T) {
	location, _ := time.LoadLocation("America/New_York")
	var temp1, temp2 int
	var err error
	temp1, err = tempProvider.Get("NYC", time.Date(2014, 2, 21, 10, 00, 00, 00, location))
	if err != nil {
		t.Error(err)
	}
	temp2, err = tempProvider.Get("NYC", time.Date(2014, 2, 21, 10, 00, 00, 00, location))
	if err != nil {
		t.Error(err)
	}
	if temp1 != temp2 {
		t.Errorf("Value should be generated only once. First temp %d should be equal to second %d", temp1, temp2)
	}
}

func TestShouldReturnSameTempFor2CallsWithSameParametersAfterDayLightSavingTimeIncrement(t *testing.T) {
	location, _ := time.LoadLocation("Europe/Paris")
	var temp1, temp2 int
	var err error
	temp1, err = tempProvider.Get("PAR", time.Date(2013, 4, 16, 10, 00, 00, 00, location))
	if err != nil {
		t.Error(err)
	}
	temp2, err = tempProvider.Get("PAR", time.Date(2013, 4, 16, 10, 00, 00, 00, location))
	if err != nil {
		t.Error(err)
	}
	if temp1 != temp2 {
		t.Errorf("Value should be generated only once. First temp %d should be equal to second %d", temp1, temp2)
	}
}

func TestGetForDatesWithNormalPeriod(t *testing.T) {
	location, _ := time.LoadLocation("Africa/Dakar")
	timesChan, err := utils.GetDatesForPeriod(time.Date(2013, 4, 16, 10, 00, 00, 00, location), "5D")
	if err != nil {
		t.Error(err)
	}
	actualTemps := tempProvider.GetForDates("DKR", timesChan)
	if 6 != len(actualTemps.Temps) {
		t.Errorf("We should get exactly 6 temps for provided period %v", actualTemps)
	}
	if len(actualTemps.Temps) < 1 || actualTemps.Temps[0].City != "DKR" {
		t.Errorf("Temps should be for DKR")
	}
	min, max, avg := tempProvider.stats(actualTemps.Temps)
	if actualTemps.Min != min || actualTemps.Max != max || actualTemps.Average != int(avg) {
		t.Errorf("Computed min %d, max %d or avg %d should not be different from %d, %d, %d", actualTemps.Min, actualTemps.Max, actualTemps.Average, min, max, avg)
	}
}

func TestGetForDatesWithSingleDatePeriod(t *testing.T) {
	location, _ := time.LoadLocation("Africa/Dakar")
	timesChan, err := utils.GetDatesForPeriod(time.Date(2013, 4, 16, 10, 00, 00, 00, location), "0D")
	if err != nil {
		t.Error(err)
	}
	actualTemps := tempProvider.GetForDates("DKR", timesChan)
	if 1 != len(actualTemps.Temps) {
		t.Errorf("We should get exactly 1 temp for provided period %v", actualTemps)
	}
	if len(actualTemps.Temps) < 1 || actualTemps.Temps[0].City != "DKR" {
		t.Errorf("Temps should be for DKR")
	}
	soleTemp := actualTemps.Temps[0].Temp
	if actualTemps.Min != soleTemp || actualTemps.Max != soleTemp || actualTemps.Average != soleTemp {
		t.Errorf("Min %d, max %d or avg %d should not be different from temp requested %d", actualTemps.Min, actualTemps.Max, actualTemps.Average, soleTemp)
	}
}

func TestStats(t *testing.T) {
	temps := []CityTemp{CityTemp{Temp: 1}, CityTemp{Temp: 2}, CityTemp{Temp: 3}, CityTemp{Temp: 4}, CityTemp{Temp: 5}}
	min, max, avg := tempProvider.stats(temps)
	if min != 1 {
		t.Errorf("Min %d should be equal to 1", min)
	}
	if max != 5 {
		t.Errorf("Max %d should be equal to 5", max)
	}
	if avg != 3 {
		t.Errorf("Average %d should be equal to 3", avg)
	}
}
