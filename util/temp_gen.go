package util

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// TempTime is defined to implement interfaces
type TempTime time.Time

// CityTemp represents a temperature for a given city at a given time
type CityTemp struct {
	City string
	Time TempTime
	Temp int
}

// CityTemps is a set of temps with complementary information
type CityTemps struct {
	Temps             []CityTemp
	Min, Max, Average int
}

// TempProvider is the component used to get temperatures
type TempProvider struct {
	// INHERITANCE BY COMPOSITION
	DataUtils
	rand *rand.Rand
}

const timeToStringFormat = time.RFC1123

var nilSample = sample{}
var id = 0

// NewTempProvider is the constructor for TempProvider
func NewTempProvider(dataUtils DataUtils) (TempProvider, error) {
	tempProvider := TempProvider{DataUtils: dataUtils}
	tempProvider.rand = rand.New(rand.NewSource(time.Now().Unix()))
	id++
	return tempProvider, nil
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
	err = tempProvider.saveTemp(temp, city, requestTime)
	if err != nil {
		return 0, err
	}
	return temp, nil
}

// GetForDates provides temperature in a given city for dates returned by time slice channel
// It uses Get underneath
func (tempProvider TempProvider) GetForDates(city string, timesChan chan []time.Time) CityTemps {
	tempsChan := make(chan []CityTemp)
	go func() {
		defer func() {
			close(tempsChan)
		}()
		tempProvider.getForDates(city, timesChan, tempsChan)
	}()

	cityTemps := CityTemps{}
	temps, open := <-tempsChan
	nbTemps := 0
	cityTemps.Min, cityTemps.Max = math.MaxInt32, math.MinInt32
	avg := 0.

	for open {
		nbNewTemps := len(temps)
		newMin, newMax, newAvg := tempProvider.stats(temps)
		avg = (float64(nbTemps)*avg + float64(nbNewTemps)*newAvg) / float64(nbTemps+nbNewTemps)
		if newMin < cityTemps.Min {
			cityTemps.Min = newMin
		}
		if newMax > cityTemps.Max {
			cityTemps.Max = newMax
		}
		cityTemps.Temps = append(cityTemps.Temps, temps...)

		nbTemps += nbNewTemps
		temps, open = <-tempsChan
	}
	cityTemps.Average = int(avg)
	return cityTemps
}

func (tempProvider TempProvider) stats(temps []CityTemp) (min int, max int, avg float64) {
	min, max = math.MaxInt32, math.MinInt32
	average := 0.
	for i, temp := range temps {
		if min > temp.Temp {
			min = temp.Temp
		}
		if max < temp.Temp {
			max = temp.Temp
		}
		average = float64(float64(i)*average+float64(temp.Temp)) / float64(i+1)
	}
	return min, max, average
}

func (tempProvider TempProvider) getForDates(city string, timesChan chan []time.Time, tempsChan chan []CityTemp) {
	times, open := <-timesChan
	var temps []CityTemp
	var err error
	var temp int

	for open {
		temps = make([]CityTemp, 0, len(times))
		for _, time := range times {
			temp, err = tempProvider.Get(city, time)
			if err != nil {
				panic(err)
			}
			temps = append(temps, CityTemp{city, TempTime(time), temp})
		}
		tempsChan <- temps
		times, open = <-timesChan
	}
}

func (tempProvider *TempProvider) generate(city string, requestTime time.Time) (int, error) {
	sample, error := tempProvider.getCityTempSample(city, requestTime)
	if error != nil {
		return 0, error
	}
	min, max := sample.TempRange[0], sample.TempRange[1]
	providerRand := tempProvider.rand
	loTemp := min - providerRand.Intn(2)
	diff := max + providerRand.Intn(3) - loTemp
	generatedTemp := loTemp + providerRand.Intn(diff)
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
	return timeIn2014AtEleven(day, requestTime)
}

func timeIn2014AtEleven(day int, requestTime time.Time) time.Time {
	return time.Date(2014, requestTime.Month(), day, 11, 0, 0, 0, requestTime.Location())
}

// HOW TO FORMAT

func (obj CityTemp) String() string {
	formattedTime := obj.Time.Format(timeToStringFormat)
	return fmt.Sprintf("%s %s %d", obj.City, formattedTime, obj.Temp)
}

// Format provides custom format for time
func (ourTime TempTime) Format(format string) string {
	return time.Time(ourTime).Format(format)
}

// MarshalJSON provides custom JSON marshaller for time
func (ourTime TempTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + ourTime.Format(TimeFormat) + "\""), nil
}
