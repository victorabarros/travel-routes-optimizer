# Bexs-Challenge

Software developer challenge from Bexs bank

## Disclaimer

To software's mission is to assert the cheapest travel for desired destination.

## Development

The programming language choiced was [Golang](https://golang.org/).
As good practice using [Docker](https://docs.docker.com/).
<!-- TODO code climate -->

### Requirements

- [Docker](https://docs.docker.com/)
- [GNU make](https://www.gnu.org/software/make/)

## How to run

#### write `.env` file

Similar to [.env.example](./.env.example)

#### clean any "garbage"

```bash
make clean-containers
make clean-network
```

### starting

To start the software must arg the .txt/.csv file with catalog* as follow example:

```csv
GRU,BRC,10
BRC,SCL,5
GRU,CDG,75
...
```

Args

```bash
./main --help

Starting Service
Usage of /tmp/go-build314365757/b001/exe/main:
  -routes string
        travel routes file (default "./input-file.txt")
```

```bash
make build
make run ROUTES=./input-file.txt
```

### Terminal interface client

<!-- TODO decrever como usar via terminal -->

### Server

<!-- TODO decrever como usar via api -->
<!-- TODO add postman collection -->

Routes:

### Healthcheck

`/liveness`
<!-- TODO liveness readness as k8s pattern-->

### Tests

```bash
make test-html-coverage
```

## References

- https://github.com/ardanlabs/service/wiki
- https://www.dudley.codes/posts/2020.05.19-golang-structure-web-servers

<!-- TODO add dockerignore -->
