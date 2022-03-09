.DEFAULT_GOAL := help

APP_NAME?=$(shell pwd | xargs basename)
APP_DIR = /go/src/github.com/victorabarros/${APP_NAME}
PWD=$(shell pwd)
DOCKER_BASE_IMAGE=golang:1.14
ROUTES?=./input-file.txt

YELLOW=\e[1m\033[33m
COLOR_OFF=\e[0m

welcome:
	@echo "${YELLOW}"
	@echo " _______                                 _       _____                    _             " && sleep .04
	@echo "|__   __|                               | |     |  __ \                  | |            " && sleep .04
	@echo "   | |     _ __    __ _  __   __   ___  | |     | |__) |   ___    _   _  | |_    ___    " && sleep .04
	@echo "   | |    | '__|  / _' | \ \ / /  / _ \ | |     |  _  /   / _ \  | | | | | __|  / _ \   " && sleep .04
	@echo "   | |    | |    | (_| |  \ V /  |  __/ | |     | | \ \  | (_) | | |_| | | |_  |  __/   " && sleep .04
	@echo "   |_|    |_|     \__,_|   \_/    \___| |_|     |_|  \_\  \___/   \__,_|  \__|  \___| \n" && sleep .04
	@echo "                ____            _     _               _                                 " && sleep .04
	@echo "               / __ \          | |   (_)             (_)                                " && sleep .04
	@echo "              | |  | |  _ __   | |_   _   _ __ ___    _   ____   ___   _ __             " && sleep .04
	@echo "              | |  | | | '_ \  | __| | | | '_ ' _ \  | | |_  /  / _ \ | '__|            " && sleep .04
	@echo "              | |__| | | |_) | | |_  | | | | | | | | | |  / /  |  __/ | |               " && sleep .04
	@echo "               \____/  | .__/   \__| |_| |_| |_| |_| |_| /___|  \___| |_|               " && sleep .04
	@echo "                       | |                                                              " && sleep .04
	@echo "                       |_|                                                            \n${COLOR_OFF}" && sleep .04

_build:
	@rm -rf ./bin/client ./bin/main
	@go build main.go && go build app/client/client.go
	@mv ./main ./bin/
	@mv ./client ./bin/

build:
	@docker run -it -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		${DOCKER_BASE_IMAGE} make _build

clean-containers:
	@docker rm -f ${APP_NAME}-debug ${APP_NAME}-server ${APP_NAME}-client ${APP_NAME}-test
	@docker network rm routes-optimizer-net

clean-network:
	@docker network rm routes-optimizer-net

create-network:
	@docker network create routes-optimizer-net

debug:
	@echo "${YELLOW}\nDebug mode${COLOR_OFF}"
	@docker run -it -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		-p 8092:8092 --name ${APP_NAME}-debug ${DOCKER_BASE_IMAGE} bash

run: welcome
	@echo "${YELLOW}\nServer up${COLOR_OFF}"
	@docker run -itd -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		--env-file .env -p 8092:8092 --network routes-optimizer-net --name ${APP_NAME}-server \
		${DOCKER_BASE_IMAGE} ./bin/main -routes ${ROUTES}
	@echo "${YELLOW}\nClient:${COLOR_OFF}"
	@docker run -it -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		--network routes-optimizer-net --name ${APP_NAME}-client ${DOCKER_BASE_IMAGE} ./bin/client

run-srv:
	@echo "${YELLOW}\nServer up${COLOR_OFF}"
	@docker run -itd -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		-p 8092:8092 --network routes-optimizer-net --name ${APP_NAME}-server ${DOCKER_BASE_IMAGE} go run main.go -routes ${ROUTES}

run-cli:
	@echo "${YELLOW}\nClient:${COLOR_OFF}"
	@docker run -it -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		--network routes-optimizer-net --name ${APP_NAME}-client ${DOCKER_BASE_IMAGE} go run app/client/client.go

_test:
	@go test ./... -v -cover -race -coverprofile=c.out
	@rm -rf internal/database/test*.csv app/server/test*.csv

test:
	@echo "\nInitalizing tests."
	@docker run --rm -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		--env-file .env --name ${APP_NAME}-test ${DOCKER_BASE_IMAGE} \
		make _test

test-log:
	@echo "Writing dev/tests.log"
	@[ ! -d "./dev" ] && mkdir dev || : # make dev folder, if not exists
	@rm -rf dev/tests*.log
	@make test > dev/tests.log
	@echo "Writing dev/tests-summ.log"
	@cat dev/tests.log  | grep "coverage: " > dev/tests-summ.log

test-coverage:
	@echo "Building c.out"
	@rm -rf c.out
	@make test
	@go tool cover -html=c.out

format:
	@docker run -v ${PWD}:${APP_DIR} -w ${APP_DIR} golang gofmt -e -l -s -w .

