.DEFAULT_GOAL := help

APP_NAME?=$(shell pwd | xargs basename)
APP_DIR = /go/src/github.com/victorabarros/${APP_NAME}
PWD=$(shell pwd)

welcome:
	@echo "\033[33m  _______                                 _       _____                    _             " && sleep .05
	@echo "\033[33m |__   __|                               | |     |  __ \                  | |            " && sleep .05
	@echo "\033[33m    | |     _ __    __ _  __   __   ___  | |     | |__) |   ___    _   _  | |_    ___    " && sleep .05
	@echo "\033[33m    | |    | '__|  / _' | \ \ / /  / _ \ | |     |  _  /   / _ \  | | | | | __|  / _ \   " && sleep .05
	@echo "\033[33m    | |    | |    | (_| |  \ V /  |  __/ | |     | | \ \  | (_) | | |_| | | |_  |  __/   " && sleep .05
	@echo "\033[33m    |_|    |_|     \__,_|   \_/    \___| |_|     |_|  \_\  \___/   \__,_|  \__|  \___| \n" && sleep .05
	@echo "\033[33m                 ____            _     _               _                                 " && sleep .05
	@echo "\033[33m                / __ \          | |   (_)             (_)                                " && sleep .05
	@echo "\033[33m               | |  | |  _ __   | |_   _   _ __ ___    _   ____   ___   _ __             " && sleep .05
	@echo "\033[33m               | |  | | | '_ \  | __| | | | '_ ' _ \  | | |_  /  / _ \ | '__|            " && sleep .05
	@echo "\033[33m               | |__| | | |_) | | |_  | | | | | | | | | |  / /  |  __/ | |               " && sleep .05
	@echo "\033[33m                \____/  | .__/   \__| |_| |_| |_| |_| |_| /___|  \___| |_|               " && sleep .05
	@echo "\033[33m                        | |                                                              " && sleep .05
	@echo "\033[33m                        |_|                                                            \n" && sleep .05

clean-up:
	@docker rm -f ${APP_NAME}

debug: welcome clean-up
	@echo "\e[1m\033[32m\nDebug mode\e[0m"
	docker run -it -v ${PWD}:${APP_DIR} -w ${APP_DIR} \
		-p 8092:8092 --name ${APP_NAME} golang:1.13 bash
