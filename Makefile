SHELL=/bin/bash

all:
	@sudo rm -rf ./deploy/rpms/x86_64
	@sudo docker-compose -f ./docker-compose-build.yml up --build
	@until [ -d ./deploy/rpms/x86_64 ]; do sleep 5; done
	@sudo docker-compose up -d --build
	@read -p "Press any key to shut down test... " -n1 -s
	@sudo docker-compose down
	@sudo docker-compose -f ./docker-compose-build.yml rm -f -v
