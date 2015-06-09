package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

type sample struct {
	Time      time.Time `json:"date"`
	TempRange []int     `json:"temp_range"`
}

type samples []sample

type cityData struct {
	Name    string  `json:"name"`
	Code    string  `json:"iata_code"`
	IanaTZ  string  `json:"iana_timezone"`
	Samples samples `json:"sample_temps"`
}

type citiesData []cityData

// City represents city's name and IATA code
type City struct {
	Name string `json:"name"`
	Code string `json:"iata_code"`
}

// Cities is a set of City
type Cities []City

// DataUtils is the component used to manipulate local data
type DataUtils struct {
	citiesFile  string
	cities      citiesData
	tempsByCity map[string]temps
}

// Temp represents a temperature at a given time
type Temp struct {
	Time time.Time
	Temp int
}

type temps []Temp

var dataUtilsNil = DataUtils{}

// NewDataUtils is the constructor for DataUtils
// citiesFile provided must exist
func NewDataUtils(citiesFile string) (DataUtils, error) {
	if _, err := os.Stat(citiesFile); os.IsNotExist(err) {
		return dataUtilsNil, fmt.Errorf("No such file or directory: %s", citiesFile)
	}
	return DataUtils{citiesFile: citiesFile, tempsByCity: make(map[string]temps)}, nil
}

// GetCities returns all cities handled by the application
func (utils DataUtils) GetCities() (Cities, error) {
	citiesData, err := utils.getCitiesData()
	if err != nil {
		return Cities{}, nil
	}
	cities := Cities{}
	for _, cityData := range citiesData {
		cities = append(cities, City{cityData.Name, cityData.Code})
	}
	return cities, nil
}

func (city City) String() string {
	return "(" + city.Code + " or " + city.Name + ")"
}

func (utils DataUtils) getCitiesData() (citiesData, error) {
	if utils.cities != nil {
		return utils.cities, nil
	}

	citiesJSONDecoder, error := getJSONDecoder(utils.citiesFile)
	if error != nil {
		return nil, error
	}
	cities := citiesData{}
	error = citiesJSONDecoder.Decode(&cities)
	if error != nil {
		return nil, error
	}
	utils.cities = cities
	return cities, nil
}

func (utils DataUtils) getTemp(city string, requestTime time.Time) (int, error) {
	cityFile, err := utils.getCityFileName(city)
	if err != nil {
		return 0, err
	}
	if !fileExists(cityFile) {
		err := fmt.Errorf("No such file or directory: %s", cityFile)
		return 0, err
	}

	var decoder *json.Decoder
	decoder, err = getJSONDecoder(cityFile)
	if err != nil {
		return 0, err
	}
	temps := temps{}
	decoder.Decode(&temps)
	utils.tempsByCity[cityFile] = temps
	for _, temp := range temps {
		if temp.Time.Equal(requestTime) {
			return temp.Temp, nil
		}
	}

	return 0, fmt.Errorf("Value has not been generated yet for %s, %s", city, requestTime)
}

func (utils DataUtils) saveTemp(temp int, city string, requestTime time.Time) error {
	cityFile, err := utils.getCityFileName(city)
	if err != nil {
		return err
	}
	var cityTemps temps
	if !fileExists(cityFile) {
		if _, err := os.Create(cityFile); err != nil {
			return err
		}
		cityTemps = temps{}
	} else {
		cityTemps = utils.tempsByCity[cityFile]
	}
	cityTemps = append(cityTemps, Temp{requestTime, temp})
	utils.tempsByCity[cityFile] = cityTemps
	var encoder *json.Encoder
	encoder, err = getJSONEncoder(cityFile)
	if err != nil {
		return err
	}
	err = encoder.Encode(cityTemps)
	if err != nil {
		return err
	}
	return nil
}

func (utils DataUtils) getCityFileName(city string) (string, error) {
	cityCode, err := utils.getCityCode(city)
	if err != nil {
		return "", err
	}
	cityFileName, err := utils.getResourceFileName(cityCode + ".json")
	if err != nil {
		return "", err
	}
	return cityFileName, nil
}

// GetCityCode used only by HTTP server
func (utils DataUtils) GetCityCode(city string) (string, error) {
	citiesData, err := utils.getCitiesData()
	if err != nil {
		return "", err
	}
	for _, cityData := range citiesData {
		if strings.EqualFold(cityData.Code, city) || strings.EqualFold(cityData.Name, city) {
			return cityData.Code, nil
		}
	}
	return "", fmt.Errorf("Have not found city %s", city)
}

func (utils DataUtils) getCityCode(city string) (string, error) {
	citiesData, err := utils.getCitiesData()
	if err != nil {
		return "", err
	}
	for _, cityData := range citiesData {
		if cityData.Code == city || cityData.Name == city {
			return cityData.Code, nil
		}
	}
	return "", fmt.Errorf("Have not found city %s", city)
}

func (utils DataUtils) getResourceFileName(resourceName string) (string, error) {
	resourcesDir := path.Dir(utils.citiesFile)
	return resourcesDir + "/" + resourceName, nil
}

func getJSONDecoder(fileLocation string) (*json.Decoder, error) {
	jsonFileReader, error := os.Open(fileLocation)
	if error != nil {
		return nil, error
	}

	jsonDecoder := json.NewDecoder(jsonFileReader)
	return jsonDecoder, nil
}

func getJSONEncoder(fileLocation string) (*json.Encoder, error) {
	jsonFileWriter, err := os.Create(fileLocation)
	if err != nil {
		return nil, err
	}
	jsonEncoder := json.NewEncoder(jsonFileWriter)
	return jsonEncoder, nil
}

func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	return true
}

// Len provided to implement sort.Interface
func (cities Cities) Len() int {
	return len(cities)
}

// Less provided to implement sort.Interface
func (cities Cities) Less(i, j int) bool {
	return cities[i].Name <= cities[j].Name
}

// Swap provided to implement sort.Interface
func (cities Cities) Swap(i, j int) {
	cities[i], cities[j] = cities[j], cities[i]
}
