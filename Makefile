# Simple Make file for this test project.

deps:
	go get github.com/Sirupsen/logrus \
	github.com/kelseyhightower/envconfig \
	gopkg.in/asaskevich/govalidator.v4

install:
	go install github.com/messageparser

test:
	go test -v github.com/messageparser/parser \
	github.com/messageparser/http \
	github.com/messageparser/clients \

run:
	go run server.go
