package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/welel/overwork-tracking/internal"
)

func main() {
	data, err := internal.StartupEnvironment()
	if err != nil {
		fmt.Printf("Can't startup a program: %v\n", err)
		os.Exit(1)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go internal.MainLoop(data)

	<-sigChan
	internal.Shutdown(data)
}
