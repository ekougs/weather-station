package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ekougs/weather-station/util"
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
	fmt.Println("Launching HTTP server...")
	http.HandleFunc("/cities", handler("GET", server.handleCityRequest))
	http.ListenAndServe(":1987", nil)
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

func (server WeatherServer) handleCityRequest(writer http.ResponseWriter, request *http.Request) {
	cities, err := server.dataUtils.GetCities()
	if hasErrorWriteResponseAndNotify(writer, err) {
		return
	}
	var citiesJSON []byte
	citiesJSON, err = json.Marshal(cities)
	if hasErrorWriteResponseAndNotify(writer, err) {
		return
	}
	writer.Header().Add(http.CanonicalHeaderKey("content-type"), "application/json")
	fmt.Fprintf(writer, "%s", citiesJSON)
}

func hasErrorWriteResponseAndNotify(responseWriter http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}
