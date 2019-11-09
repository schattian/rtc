package cli

import (
	"errors"
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

var errNotEnoughArgs = errors.New("there are NOT ENOUGH ARGS present")
