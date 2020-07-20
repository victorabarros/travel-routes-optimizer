.DEFAULT_GOAL := help

APP_NAME?=$(shell pwd | xargs basename)
APP_DIR = /go/src/github.com/victorabarros/${APP_NAME}
PWD=$(shell pwd)
DOCKER_BASE_IMAGE=golang:1.14
ROUTES?=./input-file.txt

welcome:
	@echo "\033[33m  _______                                 _       _____                    _             " && sleep .04
	@echo "\033[33m |__   __|                               | |     |  __ \                  | |            " && sleep .04
	@echo "\033[33m    | |     _ __    __ _  __   __   ___  | |     | |__) |   ___    _   _  | |_    ___    " && sleep .04
	@echo "\033[33m    | |    | '__|  / _' | \ \ / /  / _ \ | |     |  _  /   / _ \  | | | | | __|  / _ \   " && sleep .04
	@echo "\033[33m    | |    | |    | (_| |  \ V /  |  __/ | |     | | \ \  | (_) | | |_| | | |_  |  __/   " && sleep .04
	@echo "\033[33m    |_|    |_|     \__,_|   \_/    \___| |_|     |_|  \_\  \___/   \__,_|  \__|  \___| \n" && sleep .04
	@echo "\033[33m                 ____            _     _               _                                 " && sleep .04
	@echo "\033[33m                / __ \          | |   (_)             (_)                                " && sleep .04
	@echo "\033[33m               | |  | |  _ __   | |_   _   _ __ ___    _   ____   ___   _ __             " && sleep .04
	@echo "\033[33m               | |  | | | '_ \  | __| | | | '_ ' _ \  | | |_  /  / _ \ | '__|            " && sleep .04
	@echo "\033[33m               | |__| | | |_) | | |_  | | | | | | | | | |  / /  |  __/ | |               " && sleep .04
	@echo "\033[33m                \____/  | .__/   \__| |_| |_| |_| |_| |_| /___|  \___| |_|               " && sleep .04
	@echo "\033[33m                        | |                                                              " && sleep .04
	@echo "\033[33m                        |_|                                                            \n" && sleep .04

build:
	@rm -rf ./bin/*
	@docker run -it -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		${DOCKER_BASE_IMAGE} go build main.go
	@mv ./main ./bin/
	@docker run -it -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		${DOCKER_BASE_IMAGE} go build app/client/client.go
	@mv ./client ./bin/

clean-containers:
	@docker rm -f ${APP_NAME}-debug ${APP_NAME}-server ${APP_NAME}-client ${APP_NAME}-test
	@docker network rm bexs-net

clean-network:
	@docker network rm bexs-net

create-network:
	@docker network create bexs-net

debug:
	@echo "\e[1m\033[33m\nDebug mode\e[0m"
	@docker run -it -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		-p 8092:8092 --name ${APP_NAME}-debug ${DOCKER_BASE_IMAGE} bash

run: welcome
	@echo "\e[1m\033[33m\nServer up\e[0m"
	@docker run -itd -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		--env-file .env -p 8092:8092 --network bexs-net --name ${APP_NAME}-server \
		${DOCKER_BASE_IMAGE} ./bin/main -routes ${ROUTES}
	@echo "\e[1m\033[33m\nClient:\e[0m"
	@docker run -it -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		--network bexs-net --name ${APP_NAME}-client ${DOCKER_BASE_IMAGE} ./bin/client

run-srv:
	@echo "\e[1m\033[33m\nServer up\e[0m"
	@docker run -itd -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		-p 8092:8092 --network bexs-net --name ${APP_NAME}-server ${DOCKER_BASE_IMAGE} go run main.go -routes ${ROUTES}

run-cli:
	@echo "\e[1m\033[33m\nClient:\e[0m"
	@docker run -it -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		--network bexs-net --name ${APP_NAME}-client ${DOCKER_BASE_IMAGE} go run app/client/client.go

test:
	@echo "\nInitalizing tests."
	@docker run --rm -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		--env-file .env --name ${APP_NAME}-test ${DOCKER_BASE_IMAGE} \
		go test ./... -v -cover -race -coverprofile=c.out
	@rm -rf internal/database/test*.csv app/server/test*.csv

test-log:
	@echo "Writing dev/tests.log"
	@rm -rf dev/tests*.log
	@make test > dev/tests.log
	@echo "Writing dev/tests-summ.log"
	@cat dev/tests.log  | grep "coverage: " > dev/tests-summ.log

test-html-coverage:
	@echo "Building c.out"
	@rm -rf c.out
	@make test
	@go tool cover -html=c.out

format:
	@docker run -v ${PWD}:${APP_DIR} -w ${APP_DIR} golang gofmt -e -l -s -w .

