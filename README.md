# Bexs-Challenge

Software developer challenge from Bexs bank

## Disclaimer

To software's mission is to assert to client the best option of transfer for desired destination.

## Development

The programming language choiced was [Golang](https://golang.org/), that's has good performance, 
<!-- scalabity. E o author está desenvolvendo novas habilidades com ela. -->
As good practice using [Docker](https://docs.docker.com/) to manage 
<!-- TODO Se a lista for maior que quantas linhas irá estourar a memória e melhor usar um redis? -->
<!-- TODO code climate -->

### Requirements

- [Docker](https://docs.docker.com/)
- [GNU make](https://www.gnu.org/software/make/)

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
make run-dev
```

### Consuming via* terminal* interface (client)

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
