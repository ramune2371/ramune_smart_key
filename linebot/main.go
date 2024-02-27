package main

import (
	"linebot/props"
	"linebot/server"
)

func main() {
	props.LoadEnv()

	server.StartServer()
}
