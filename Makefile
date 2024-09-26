run:
	go run main.go
test:
	go test -v -race ./... -coverprofile=coverage.out
genmocks:
	mockgen -source service/service.go -destination service/service_mock_test.go -package service
cover:
	go tool cover -html=coverage.out