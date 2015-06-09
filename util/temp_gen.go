package util

import (
	"fmt"
	"math/rand"
	"time"
)

// TempProvider is the component used to get temperatures
type TempProvider struct {
	// INHERITANCE BY COMPOSITION
	DataUtils
	seed int64
}

var nilSample = sample{}

// NewTempProvider is the constructor for TempProvider
func NewTempProvider(dataUtils DataUtils) (TempProvider, error) {
	return TempProvider{DataUtils: dataUtils}, nil
}

// Get provides temperature in a given city at a given time. It generates
// and stores it locally if it doesn't exist otherwise it returns the stored
// value
func (tempProvider TempProvider) Get(city string, requestTime time.Time) (int, error) {
	temp, err := tempProvider.getTemp(city, requestTime)
	if err == nil {
		return temp, nil
	}
	temp, err = tempProvider.generate(city, requestTime)
	if err != nil {
		return 0, err
	}
	err = tempProvider.setTemp(temp, city, requestTime)
	if err != nil {
		return 0, err
	}
	return temp, nil
}

func (tempProvider *TempProvider) generate(city string, requestTime time.Time) (int, error) {
	sample, error := tempProvider.getCityTempSample(city, requestTime)
	if error != nil {
		return 0, error
	}
	min, max := sample.TempRange[0], sample.TempRange[1]
	seed := time.Now().Unix()
	if tempProvider.seed != seed {
		tempProvider.seed = seed
		rand.Seed(seed)
	}
	loTemp := min - rand.Intn(2)
	diff := max + rand.Intn(3) - loTemp
	generatedTemp := loTemp + rand.Intn(diff)
	return generatedTemp, nil
}

func (tempProvider TempProvider) getCityTempSample(city string, requestTime time.Time) (sample, error) {
	citiesData, error := tempProvider.getCitiesData()
	if error != nil {
		return nilSample, error
	}
	for _, tempSample := range citiesData {
		if tempSample.Name != city && tempSample.Code != city {
			continue
		}
		return getCityTempSample(tempSample, requestTime)
	}
	return nilSample, fmt.Errorf("Should always be able to find a sample for %s and %s", city, requestTime)
}

func getCityTempSample(tempSample cityData, requestTime time.Time) (sample, error) {
	sampleTime := getSampleTime(requestTime)
	for _, sample := range tempSample.Samples {
		// THIS IS REALLY IMPORTANT CANNOT COMPARE WITH ==
		if sampleTime.Equal(sample.Time) {
			return sample, nil
		}
	}
	return sample{}, fmt.Errorf("Should always be able to find a sample for %s", requestTime)
}

func getSampleTime(requestTime time.Time) time.Time {
	day := requestTime.Day()
	var i = int(day / 10)
	// NO TERNARY OPERATOR IN GOLANG
	if day%10 != 0 && i == 0 {
		day = 1
	} else if day%10 != 0 {
		day = i * 10
	}
	return timeAtEleven(day, requestTime)
}

func timeAtEleven(day int, requestTime time.Time) time.Time {
	return time.Date(requestTime.Year(), requestTime.Month(), day, 11, 0, 0, 0, requestTime.Location())
}
