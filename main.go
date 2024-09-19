package main

import (
	"pinder/repository"
	"pinder/server"
	"pinder/service"
)

func main() {
	repo := repository.New(nil)
	service := service.New(repo)
	server := server.New(service)
	_ = server
	// TODO: migrations
}
