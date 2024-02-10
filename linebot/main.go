package main

import (
	"linebot/props"
	"linebot/server"
	"linebot/transfer"
)

func main() {
	props.LoadEnv()

	transfer.InitLineBot()

	server.StartServer()
}
