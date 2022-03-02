# Water bucket challenge

Build an application that solves the Water Jug Riddle for dynamic inputs (X, Y, Z). 

The simulation should have a UI to display state changes for each state for each jug (Empty, Full or Partially Full).

You have an X-gallon and a Y-gallon jug that you can fill from a lake. (Assume lake has unlimited amount of water.) By using only an X-gallon and Y-gallon jug (no third jug), measure Z gallons of water.

To find solution run this service and check swagger ui [here](http://localhost:8001/swagger/)

## Build

to build an app just do `make build`

binary is stored in `./bin` directory

## Run tests
`make test`

also linter available with `make lint`

## Run with docker compose
`make run-docker-compose`

## Features
- http public port (configurable, e.g. 8000)
- http debug port (configurable, e.g. 8001) (for metrics/debug/swagger/health-checks/etc)
- graceful shutdown (configurable timeout)
- metrics support
- tracing support (sends span to jaeger-agent, logs trace id)

### Swagger
- serves swagger ui on debug port `http://localhost:8001/swagger/`
- serves swagger.json on debug port `http://localhost:8001/swagger/doc.json`

### Metrics
- serves prometheus metrics on debug port `http://localhost:8001/metrics`
- standard go runtime metrics 
- http handler metrics 
```
chi_request_duration_milliseconds_count{code="Bad Request",method="POST",path="/v1/search-solution",service=""} 1
chi_request_duration_milliseconds_bucket{code="OK",method="POST",path="/v1/search-solution",service="",le="300"} 4
chi_request_duration_milliseconds_bucket{code="OK",method="POST",path="/v1/search-solution",service="",le="1200"} 4
chi_request_duration_milliseconds_bucket{code="OK",method="POST",path="/v1/search-solution",service="",le="5000"} 4
chi_request_duration_milliseconds_bucket{code="OK",method="POST",path="/v1/search-solution",service="",le="+Inf"} 4
chi_request_duration_milliseconds_sum{code="OK",method="POST",path="/v1/search-solution",service=""} 19.034916000000003
chi_request_duration_milliseconds_count{code="OK",method="POST",path="/v1/search-solution",service=""} 4
...
```

### Tracing
to check tracing spans visit jaeger ui [here](http://localhost:16686/)