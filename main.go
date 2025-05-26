package main

import (
	"flag"
	"fmt"
	"github.com/bhaski-1234/redis-db/config"
)

func initFlags() {
	flag.StringVar(&config.Host, "host", "localhost", "Redis server host")
	flag.IntVar(&config.Port, "port", 6379, "Redis server port")
}

func main() {
	initFlags()
	fmt.Println("Hello World")
}
