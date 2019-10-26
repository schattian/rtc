package main

import (
	"errors"

	"github.com/spf13/afero"

	"github.com/sebach1/git-crud/fabric"
	"github.com/sebach1/git-crud/schema"
	"github.com/urfave/cli"
)

var fabricCmd = cli.Command{
	Name:        "fabric",
	Aliases:     []string{"f"},
	Description: "Create structs from the given schema tagging fields with the given marshal type",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "marshal, msh"},
	},
	Usage:  "[schema-path]",
	Before: fabricValidate,
	Action: fabricExec,
}

func fabricExec(c *cli.Context) error {
	schemaName := c.Args().First()
	marshalType := c.String("marshal, msh")
	if marshalType == "" {
		marshalType = "json"
	}

	decodedSchema, err := schema.FromFilename(schemaName)
	if err != nil {
		panic(err)
	}

	Fabric := &fabric.Fabric{Schema: decodedSchema}
	err = Fabric.Produce(marshalType, afero.NewOsFs())
	if err != nil {
		return err
	}
	return nil
}

func fabricValidate(c *cli.Context) error {
	if !c.Args().Present() {
		return errors.New("there are NOT ENOUGH ARGS present")
	}
	return nil
}
