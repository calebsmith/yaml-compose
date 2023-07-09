SHELL := /bin/bash

build:
	go build -o yaml-compose src/*.go

build_test:
	pip3 install -r requirements.txt
test:
	./run-tests.sh
