package main

import "github.com/akamensky/argparse"

type ImportCmd struct {
	cmd *argparse.Command

	File *string
}

func SetupImportCmd(parser *argparse.Parser) *ImportCmd {

	imCmd := &ImportCmd{}

	cmd := parser.NewCommand("import", "Import the data from the provided csv file.")
	imCmd.cmd = cmd

	imCmd.File = cmd.StringPositional(&argparse.Options{Required: true, Help: "Path to the file that is gonna be read (in the csv format)."})

	return imCmd
}

func (cmd *ImportCmd) GetFile() string {
	return *cmd.File
}

func (cmd *ImportCmd) Happened() bool {
	return cmd.cmd.Happened()
}
