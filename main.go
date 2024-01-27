package main

import "go/src/adapters/web/server"

func main() {
	server.MakeNewWebServer().Serve()
}
