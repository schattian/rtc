package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		fabricCmd,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
