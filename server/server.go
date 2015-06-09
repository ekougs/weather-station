package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ekougs/weather-station/util"

	"github.com/gorilla/mux"
)

var stopPwd = "Gophers$ince2014"

// WeatherServer provides HTTP access to all features available from the CLI
type WeatherServer struct {
	tempProvider util.TempProvider
	timeUtils    util.TimeUtils
	dataUtils    util.DataUtils
}

// NewWeatherServer provides a fully configured WeatherServer
func NewWeatherServer(tempProvider util.TempProvider, timeUtils util.TimeUtils, dataUtils util.DataUtils) WeatherServer {
	return WeatherServer{tempProvider, timeUtils, dataUtils}
}

// LaunchServer launch an HTTP server which offers all features offered by CLI
func (server WeatherServer) LaunchServer() {
	router := mux.NewRouter().StrictSlash(true)
	fmt.Println("Launching HTTP server...")
	router.HandleFunc("/cities", server.handleCityRequest).Methods("GET")
	router.HandleFunc("/cities/{city}/temps", server.handleTempRequest).Methods("GET")
	http.ListenAndServe(":1987", router)
}

func handler(supportedMethod string, handler func(writer http.ResponseWriter, request *http.Request)) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		method := request.Method
		if method != supportedMethod {
			http.Error(writer, method+" NOT ALLOWED", http.StatusMethodNotAllowed)
			return
		}
		handler(writer, request)
	}
}

func (server WeatherServer) handleTempRequest(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	cityCode, err := server.dataUtils.GetCityCode(vars["city"])
	if hasErrorWriteResponseAndNotify(writer, err) {
		return
	}

	var date time.Time
	date, err = server.getDateParam(request, cityCode)
	if hasErrorWriteResponseAndNotify(writer, err) {
		return
	}

	duration := request.FormValue("duration")
	if duration != "" {
		server.respondForTempWithDuration(writer, request, cityCode, date, duration)
		return
	}

	var temp int
	temp, err = server.tempProvider.Get(cityCode, date)
	if hasErrorWriteResponseAndNotify(writer, err) {
		return
	}
	writeSucessfulResponse(writer, util.CityTemp{City: cityCode, Time: util.TempTime(date), Temp: temp})
}

func (server WeatherServer) respondForTempWithDuration(writer http.ResponseWriter, request *http.Request, city string, date time.Time, duration string) {
	datesChan, err := server.timeUtils.GetDatesForPeriod(date, duration)
	if hasErrorWriteResponseAndNotify(writer, err) {
		return
	}
	cityTemps := server.tempProvider.GetForDates(city, datesChan)
	writeSucessfulResponse(writer, cityTemps)
}

func (server WeatherServer) getDateParam(request *http.Request, city string) (time.Time, error) {
	parsedDate := request.FormValue("date")
	var date time.Time
	var err error
	if "" == parsedDate {
		date = server.timeUtils.GetTimeWithoutMinuteSecondNano(time.Now())
	} else {
		date, err = server.timeUtils.GetTime(parsedDate, city)
	}
	return date, err
}

func (server WeatherServer) handleCityRequest(writer http.ResponseWriter, request *http.Request) {
	cities, err := server.dataUtils.GetCities()
	if hasErrorWriteResponseAndNotify(writer, err) {
		return
	}
	writeSucessfulResponse(writer, cities)
}

func writeSucessfulResponse(writer http.ResponseWriter, result interface{}) {
	resultJSON, err := json.Marshal(result)
	if hasErrorWriteResponseAndNotify(writer, err) {
		return
	}
	writer.Header().Add(http.CanonicalHeaderKey("content-type"), "application/json")
	fmt.Fprintf(writer, "%s", resultJSON)
}

func hasErrorWriteResponseAndNotify(responseWriter http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}
