package main

import (
	"linebot/dao"
	"linebot/server"
	"linebot/transfer"
)

func main() {
	dao.InitDB()
	defer dao.Close()

	transfer.InitLineBot()

	server.StartServer()
}
