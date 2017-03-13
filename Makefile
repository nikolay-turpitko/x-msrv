.PHONY: all clean build

SHELL=/bin/bash

export PKG_TYPE=rpm

all: | clean build
	@sudo PKG_TYPE=$(PKG_TYPE) docker-compose up -d --build
	@read -p "Press any key to shut down test... " -n1 -s
	@sudo PKG_TYPE=$(PKG_TYPE) docker-compose down
	@sudo PKG_TYPE=$(PKG_TYPE) docker-compose -f ./docker-compose-build.yml rm -f -v

clean:
	@sudo rm -rf ./$(PKG_TYPE)-deploy/$(PKG_TYPE)/*

build:
	@sudo PKG_TYPE=$(PKG_TYPE) docker-compose -f ./docker-compose-build.yml up --build
	@if [ "$(PKG_TYPE)" == "rpm" ]; then \
		until [ -d ./$(PKG_TYPE)-deploy/$(PKG_TYPE)/x86_64 ]; do sleep 5; done \
	else \
		until [ -f ./$(PKG_TYPE)-deploy/$(PKG_TYPE)/*.deb ]; do sleep 5; done \
	fi
