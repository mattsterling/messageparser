# Simple Make file for this test project.

deps:
	go get github.com/Sirupsen/logrus
	go get github.com/kelseyhightower/envconfig

install:
	go install github.com/messageparser

test:
	go test

run: 
	go run server.go
