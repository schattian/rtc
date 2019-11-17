package cli

import (
	"github.com/spf13/afero"

	"github.com/sebach1/git-crud/schema"
	"github.com/urfave/cli"
)

var registerCmd = cli.Command{
	Name:        "register",
	Aliases:     []string{"reg", "r"},
	Description: "Register your own (or third party) schemas giving a yml path",
	Usage:       "[schema-path]",
	Before:      registerValidate,
	Action:      registerExec,
	Flags:       []cli.Flag{
		// cli.StringFlag{
		// 	Name:     "marshal, msh",
		// 	Value:    "json",
		// 	Required: false,
		// 	Usage:    "marshal format of the fabricated fields",
		// },
	},
}

func registerExec(c *cli.Context) error {
	schFilename := c.Args().First()
	// marshalType := c.String("marshal, msh")

	osFs := afero.NewOsFs()

	decodedSchema, err := schema.FromFilename(schFilename, osFs)
	if err != nil {
		return err
	}

	return nil
}

func registerValidate(c *cli.Context) error {
	if !c.Args().Present() {
		return errNotEnoughArgs
	}
	return nil
}
