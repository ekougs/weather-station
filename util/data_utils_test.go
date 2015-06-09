package util

import "testing"

func TestCreationShouldFailIfCitiesFileProvidedDoesNotExist(t *testing.T) {
	_, err := NewDataUtils("nimportequoi")
	if err == nil {
		t.Errorf("Should have error as file does not exist")
	}
}
