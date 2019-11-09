package cli

import (
	"github.com/spf13/afero"

	"github.com/sebach1/git-crud/fabric"
	"github.com/sebach1/git-crud/schema"
	"github.com/urfave/cli"
)

var fabricCmd = cli.Command{
	Name:        "fabric",
	Aliases:     []string{"f"},
	Description: "Create structs from the given schema tagging fields with the given marshal type",
	Usage:       "[schema-path]",
	Before:      fabricValidate,
	Action:      fabricExec,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "marshal, msh",
			Value: "json",
			Usage: "marshal format of the fabricated fields",
		},
	},
}

func fabricExec(c *cli.Context) error {
	schemaName := c.Args().First()
	marshalType := c.String("marshal, msh")

	osFs := afero.NewOsFs()

	decodedSchema, err := schema.FromFilename(schemaName, osFs)
	if err != nil {
		return err
	}

	Fabric := &fabric.Fabric{Schema: decodedSchema}
	err = Fabric.Produce(marshalType, osFs)
	if err != nil {
		return err
	}
	return nil
}

func fabricValidate(c *cli.Context) error {
	if !c.Args().Present() {
		return errNotEnoughArgs
	}
	return nil
}
