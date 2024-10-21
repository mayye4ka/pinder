build:
	go build -o pinder ./cmd/main.go
run:
	go run main.go
test:
	go test -v -race ./... -coverprofile=coverage.out
genmocks:
	mockgen -source internal/usecase/authenticator/authenticator.go -destination internal/usecase/authenticator/authenticator_mock_test.go -package authenticator
	mockgen -source internal/usecase/service/service.go -destination internal/usecase/service/service_mock_test.go -package service
cover:
	go tool cover -html=coverage.out