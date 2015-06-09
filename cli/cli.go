package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ekougs/weather-station/util"
)

type instructionParams struct {
	city, duration string
	date           time.Time
}

// InstructionHandler handles an instruction from CLI
type InstructionHandler struct {
	instructionParams
}

// NewCLIInstructionHandler creates a cli handler
func NewCLIInstructionHandler(city string, date time.Time, duration string) InstructionHandler {
	params := instructionParams{city, duration, date}
	return InstructionHandler{params}
}

// PrintResponse pretty print in JSON format a response matching parameters
// of the instruction handler
func (handler InstructionHandler) PrintResponse(tempProvider util.TempProvider, timeUtils util.TimeUtils) {
	var data interface{}
	if "" == handler.duration {
		data = handler.getTemp(tempProvider)
	} else {
		data = handler.getTemps(tempProvider, timeUtils)
	}
	handler.prettyPrintJSON(data)
}

func (handler InstructionHandler) getTemp(tempProvider util.TempProvider) util.CityTemp {
	city, date := handler.city, handler.date
	temp, err := tempProvider.Get(city, date)
	IfErrorInformAndLeave(err)
	return util.CityTemp{City: city, Time: util.TempTime(date), Temp: temp}
}

func (handler InstructionHandler) getTemps(tempProvider util.TempProvider, timeUtils util.TimeUtils) util.CityTemps {
	city, date, duration := handler.city, handler.date, handler.duration
	datesChan, err := timeUtils.GetDatesForPeriod(date, duration)
	IfErrorInformAndLeave(err)
	return tempProvider.GetForDates(city, datesChan)
}

func (handler InstructionHandler) prettyPrintJSON(object interface{}) {
	prettyJSON, err := json.MarshalIndent(object, "", "    ")
	IfErrorInformAndLeave(err)
	fmt.Println(string(prettyJSON))
}

// IfErrorInformAndLeave is an utility function to print error, show information
// and leave program with error status
func IfErrorInformAndLeave(err error) {
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}
}
