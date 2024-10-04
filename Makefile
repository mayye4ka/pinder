build:
	go build
run:
	go run main.go
test:
	go test -v -race ./... -coverprofile=coverage.out
genmocks:
	mockgen -source authenticator/authenticator.go -destination authenticator/authenticator_mock_test.go -package authenticator
	mockgen -source service/service.go -destination service/service_mock_test.go -package service
cover:
	go tool cover -html=coverage.out