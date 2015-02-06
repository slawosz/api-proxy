package main

import (
	"helpers"
	"server"
)

func main() {
	go helpers.ServeMock() // FIXME: for testing, to be removed
	server.Start()
}
