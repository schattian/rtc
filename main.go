package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/sebach1/git-crud/server"

	"github.com/valyala/fasthttp"

	"github.com/sebach1/git-crud/config"
)

func main() {
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = config.DefaultPort
	}
	Port = fmt.Sprintf(":%s", Port)

	log.Println(fmt.Sprintf("Accepting connections at: %s", Port))
	log.Fatal(fasthttp.ListenAndServe(Port, server.Router))
}
