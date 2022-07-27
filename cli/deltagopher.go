package cli

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var CliApp = &cli.App{
	Name:     "deltagopher",
	HelpName: "dg",
	Version:  "v1",
	Authors: []*cli.Author{
		&cli.Author{
			Name:  "Maria Oleksik",
			Email: "oleksik.maria@gmail.com",
		},
	},
	Usage:     "",
	UsageText: "",
	ArgsUsage: "",
	Commands: []*cli.Command{
		GetSignatureCommand(),
		GetDeltaCommand(),
	},
}

func main() {
	if err := CliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
