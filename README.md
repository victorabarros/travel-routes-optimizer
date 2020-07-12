# Bexs-Challenge

Software developer challenge from Bexs bank

## Disclaimer

To software's mission is to assert to client the best option of transfer for desired destiny.

## Development

The programming language choiced was [Golang](https://golang.org/), that's has good performance, 
<!-- scalabity. E o author está desenvolvendo novas habilidades com ela. -->
As good practice using [Docker](https://docs.docker.com/) to manage 
<!-- blablabla -->
<!-- TODO Se a lista for maior que quantas linhas irá estourar a memória e melhor usar um redis? -->

### Requirements

- [Docker](https://docs.docker.com/)
- [Make]()

## How to use

### Starting

Starting the software must arg the .txt/.csv file with catalog* as follow example:

```csv
GRU,BRC,10
BRC,SCL,5
GRU,CDG,75
...
```

```bash
make start
```

### Consuming via* terminal* client

<!-- TODO decrever como usar via terminal -->

### Server

<!-- TODO decrever como usar via api -->
<!-- TODO add postman collection -->

Routes:

### Healthcheck

<!-- TODO liveness readness as k8s pattern (olhar hotel-worker)-->

## References

- https://github.com/ardanlabs/service/wiki
- https://www.dudley.codes/posts/2020.05.19-golang-structure-web-servers
<!-- 
├── app/                    # entry point newcomers gravitate towards when exploring the codebase
|   └── service-api/        # micro-service API for this repository; all HTTP implementation details live here
|       ├── cfg/            # configuration files, usually json or yaml saved in plain text files, as they should be checked into git too
|       ├── middleware/     # for all middleware
|       ├── routes/         # API application’s RESTFul-like surface
|       |   ├── makes/
|       |   |   └── models/**
|       |   ├── create.go
|       |   ├── create_test.go
|       |   ├── get.go
|       |   └── get_test.go
|       ├── webserver/      # contains all shared HTTP structs and interfaces (Broker, configuration, Server, etc)
|       ├── main.go         # bootstrapped (New(), Start())
|       └── routebinds.go   # BindRoutes() function
├── cmd/                    # where any command-line applications belong
|   └── service-tool-x/
├── internal/               # directory that cannot be imported by projects outside of this repo
|   └── service/            # domain logic; it can be imported by service-api
|       └── mock/
└── pkg/                    # packages that are encouraged to be imported by projects outside this repo
    ├── client/             # library for accessing service-api. Other teams can import it without having to write their own
    └── dtos/               # data transfer objects, structs designed for sharing data between packages and encoding/transmitting. /internal/service is responsible for mapping the DTOs to/from its internal models -->


<!-- TODO add gitignore dockerignore -->
