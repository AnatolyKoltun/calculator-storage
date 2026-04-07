package main

import (
	"github.com/AnatolyKoltun/calculator-storage/database"
	"github.com/AnatolyKoltun/calculator-storage/message_broker"
	"github.com/AnatolyKoltun/calculator-storage/rpc"
)

func main() {
	var StreamNats message_broker.HandleNats

	defer database.Close()
	defer StreamNats.Close()

	database.Connect()
	StreamNats.CreateStreamNats()
	rpc.RunningGrpc()
}
