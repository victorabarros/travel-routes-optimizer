# travel routes optimizer

[![Maintainability](https://api.codeclimate.com/v1/badges/a99a88d28ad37a79dbf6/maintainability)](https://codeclimate.com/github/victorabarros/travel-routes-optimizer) [![Go Report Card](https://goreportcard.com/badge/github.com/victorabarros/travel-routes-optimizer)](https://goreportcard.com/report/github.com/victorabarros/travel-routes-optimizer) [![MIT Licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/victorabarros/travel-routes-optimizer/master/LICENSE)
<!-- TODO circleci, codeclimate test coverage, license link -->

Software developer challenge

## Disclaimer

To software's mission is to assert the cheapest travel for desired destination.

## Development

The programming language choiced was [Golang](https://golang.org/).
As good practice using [Docker](https://docs.docker.com/).

### Requirements

- [Docker](https://docs.docker.com/)
- [GNU make](https://www.gnu.org/software/make/)

## How to run

### csv

To start the software needs a file with possibles routes.
Example:

```csv
GRU,BRC,10
BRC,SCL,5
GRU,CDG,75
...
```

#### write `.env` file

Similar to [.env.example](./.env.example)

#### clean any "garbage"

```bash
make clean-containers
make clean-network
```

#### set and build

```bash
make create-network
make build
```

#### flags

```bash
./main --help

Starting Service
Usage of /tmp/go-build314365757/b001/exe/main:
  -routes string
        travel routes file (default "./input-file.txt")
```

#### start

```bash
make run ROUTES=./input-file.txt
```

## Terminal interface client

Enter `ORG-DES` format.
Example:

```bash
GRU-ORL
```

Answer:

```bash
best route: BRC - SCL - GRU - ORL > $35.00
```

### Server

### Healthcheck

Based on [k8s best practices](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)

- `/healthz`
- `/started`

### Search

To find the cheapest transfer travel option.

- `/routes?origin=GRU&destination=ORL`

### Insert

To insert new route. It will be perssistent on data inputed file.

- `/routes`
  - `Method: Post`
  - `Body: {"origin": "GRU", "destination": "BRC", "price": 10}`

## Tests

```bash
make test-coverage
```
