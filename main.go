package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Basu008/EasySplit.git/server"
)

func main() {
	fmt.Println("Started EasySplit systems....")
	s := server.NewServer()
	s.StartServer()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c

	s.StopServer()
	fmt.Print("Server Stopped")
}
