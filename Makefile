SHELL := /bin/bash
BUILD_DIR="$(GOPATH)"/bin/weather-station
BUILD_RESOURCES_DIR="$(BUILD_DIR)"/resources

all: pkg_weather_station clean

pkg_weather_station:weather-station pkg_resources
	@echo "Packaging application..."
	@mkdir -p $(BUILD_DIR)
	@cp weather-station $(BUILD_DIR)

pkg_resources:
	@echo "Packaging resources..."
	@mkdir -p $(BUILD_RESOURCES_DIR)
	@cp -r resources/ $(BUILD_RESOURCES_DIR)

weather-station: weather-station.go test
	@echo "Building exec..."
	@go build

test: check_init
	@echo "Running tests..."
	@go test -v ./...

clean:
	@rm weather-station

check_init:
	@if [ ! -d "$(GOPATH)" ]; then echo "Setup your GOPATH, this one does not exist: $(GOPATH)"; exit 1; fi
