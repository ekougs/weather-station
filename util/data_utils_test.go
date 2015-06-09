package util

import (
	"sort"
	"testing"
)

func TestCreationShouldFailIfCitiesFileProvidedDoesNotExist(t *testing.T) {
	_, err := NewDataUtils("nimportequoi")
	if err == nil {
		t.Errorf("Should have error as file does not exist")
	}
}

func TestGetCitiesList(t *testing.T) {
	var utils, err = NewDataUtils("../resources/cities.json")
	if err != nil {
		t.Error(err)
	}
	var cities Cities
	cities, err = utils.GetCities()
	if err != nil {
		t.Error(err)
	}
	expectedCities := Cities{City{"Dakar", "DK"}, City{"New York", "NYC"}, City{"Paris", "PAR"}}
	assertEquals(cities, expectedCities, t)
}

func assertEquals(actual, expected Cities, t *testing.T) {
	sort.Sort(actual)
	sort.Sort(expected)
	for _, city := range actual {
		if !contains(expected, city) {
			t.Errorf("City %s has not been found in expected.", city.Name)
		}
	}
	for _, city := range expected {
		if !contains(actual, city) {
			t.Errorf("City %s has not been found in actual.", city.Name)
		}
	}
}

func contains(cities Cities, city City) bool {
	numberOfCitiesExpected := len(cities)
	index := sort.Search(numberOfCitiesExpected, func(i int) bool {
		return city.Name <= cities[i].Name
	})
	return index >= 0 && index < numberOfCitiesExpected && city.Name == cities[index].Name
}
