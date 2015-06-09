# Dummy Weather

A Go language simple demo application simulating a weather station for some cities.
It allows you to discover step by step new Go features as application grows.
Each step is a tag "Step \# : description of step with features introduced"

## Setup instructions:
[Set your GOPATH](https://golang.org/doc/code.html#GOPATH)  
Clone the repository in `$GOPATH/src/com.github/ekougs`  

To build the application:
###Run `make` command.
It will install the application in `$GOPATH/bin/weather-station/`.

To launch: `$GOPATH/bin/weather-station/weather-station` (--help for more information)  

###Using with docker
There is a ``Dockerfile`` in this repository, so you do not need to have golang installed on your computer to try this out. To build and run it, do the following :

```bash
$ docker build -t weather-station .
# […]
$ docker run -ti --rm weather-station                        # To get the default
$ docker run -ti --rm -p 1987:1987 weather-station -s        # To get the server
$ docker run -ti --rm --entrypoint /bin/bash weather-station # To get in there with a shell
```

If you want to develop on it, and use [goconvey](http://goconvey.co/), that's easy too :

```
$ docker run -ti --rm \
             --entrypoint goconvey \
             -p 8080:8080 \
             -v $PWD:/go/src/github.com/ekougs/weather-station \
             weather-station \
             -host="0.0.0.0"
```
