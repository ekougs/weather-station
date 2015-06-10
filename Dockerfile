FROM golang:1.4

# Grab some go tools
# - cover
# - vet
# - goconvey (for continuous test)
RUN go get golang.org/x/tools/cmd/cover \
    && go get golang.org/x/tools/cmd/vet \
    && go get github.com/smartystreets/goconvey

COPY . /go/src/github.com/ekougs/weather-station
WORKDIR /go/src/github.com/ekougs/weather-station

EXPOSE 1987 8080

# Download dependencies
RUN go-wrapper download
# Install it
RUN go-wrapper install

ENTRYPOINT ["weather-station"]
