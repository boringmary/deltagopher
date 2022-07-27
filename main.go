package main

import (
	"log"
	"os"

	"deltagopher/cli"
)

const defaultWindow = 3
const defaultSigFilename = "signature.yml"
const defaultDeltaFilename = "delta.yml"

func main() {
	if err := cli.CliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
