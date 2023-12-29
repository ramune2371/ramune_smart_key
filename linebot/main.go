package main

import (
	"linebot/dao"
	"linebot/props"
	"linebot/server"
	"linebot/transfer"
)

func main() {
	props.LoadEnv()

	dao.InitDB()
	defer dao.Close()

	transfer.InitLineBot()

	server.StartServer()
}
