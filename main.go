package main

import (
	"flag"
	"fmt"

	"github.com/bhaski-1234/redis-db/config"
	"github.com/bhaski-1234/redis-db/server"
)

func initFlags() {
	flag.StringVar(&config.Host, "host", "localhost", "Redis server host")
	flag.IntVar(&config.Port, "port", 6379, "Redis server port")
}

func main() {
	initFlags()
	server := server.NewServer()
	if err := server.Start(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
	defer server.Close()
	fmt.Printf("Server started on %s:%d\n", config.Host, config.Port)
}
