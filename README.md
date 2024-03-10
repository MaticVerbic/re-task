# RE Partners tech task

## Instructions: 

Our customers can order any number of these items through our website, but they will always only
be given complete packs.
1. Only whole packs can be sent. Packs cannot be broken open.
2. Within the constraints of Rule 1 above, send out the least amount of items to fulfil the order.
3. Within the constraints of Rules 1 & 2 above, send out as few packs as possible to fulfil each
   order.
   (Please note, rule #2 takes precedence over rule #3)

Write an application that can calculate the number of packs we need to ship to the customer.
The API must be written in Golang & be usable by a HTTP API (by whichever method you
choose) and show any relevant unit tests.

### Explanation of the algorithm
For the solution I chose the greedy algorithm, but to also be in line with rules 2. and 3. I first calculated the approximation of the smallest buckets
Then ran the greedy algorithm, meaning removing as many packages as it will fit from the biggest to smallest. 
In the end I also took care of the possible missing smallest package when the packages are not multiples of the smallest package. 

### Structure
I chose the standard structure for the golang API service. 

- cmd will contain the apps main package
- api will contain the apps api handlers and request/response models
- internal will deal with the actual logic of things, here we would also add databases and other thing. 
- server will hold the custom middleware and server code and endpoint definitions

I chose the client/consumer interface pattern as shown in the handler package, which defines the interface implemented by Packager repo. 
For the server framework I chose chi, because it's lightweight, and it implements default http handlers, which is good for prototyping. 
Config is parsed by viper, and config.yaml is located in directory root.

#### Config.yaml
Config.yaml is composed of several settings: 
* serverPort: 8080
  * specifies the port on which the server will listen
  * if you wish to change this, make sure you also update it in the Dockerfile and docker-compose file. 
* httpTimeout: 10 # in seconds
  * The context will cancel the request after this amount of time in seconds
* logType: text
  * This simply specifies the formatter used in the service
  * possible options here are `text` and `json` 
* logLevel: debug
  * level of logs that will be output to the std out
* packs: `- 250 - 500 - 1000 - 2000 - 5000` including a new line after each of the numbers


### Tests
Since the task is only supposed to take 2 hours, I chose to only implement unittests on the actual algorithm and no end-to-end testing. 
You can run the tests with `go test ./...` or `go test -v ./...` for verbose option with logs. 

I've also added a benchmark test for the algorithm which can be run using `go test -v ./internal/packing/... -bench .`

### Running the service

While you can also run the service simply by running `go run cmd/main.go`, Docker and docker-compose are also available

#### Docker
To build the docker image run the command `docker build --tag=retask .`, to rebuild run `docker build --tag=retask --no-cache .` ensuring the image gets rebuild fully. 
To run the built image use `docker run -p 8080:8080 retask`, you can also use any other combination of ports, depending on your settings in `config.yaml`.

#### Docker-compose
To build and run from docker compose use `docker-compose up -d` for the first time, subsequent times the image will simply be reused. 
If you wish to rebuild and run `docker-compose up -d --build`