package cli

import (
	"io/ioutil"
	"log"
	"os"

	"deltagopher/delta"
	"deltagopher/signature"
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
	Flags: []cli.Flag{
		&cli.IntFlag{Name: "block-size"},
		&cli.BoolFlag{Name: "full"},
	},
	Action: func(cCtx *cli.Context) error {
		oldFile := cCtx.Args().Get(0)
		newFile := cCtx.Args().Get(1)

		f1, err := os.Open(oldFile)
		defer f1.Close()
		sb := signature.NewSignatureBuilder(f1, cCtx.Value("block-size").(int), cCtx.Bool("full"))
		signature := sb.BuildSignature()

		// Try marshal and save
		mrhs, _ := signature.MarshalYAML()
		err = ioutil.WriteFile(defaultSigFilename, mrhs, 0644)
		if err != nil {
			panic(err)
		}

		f, err := os.Open(newFile)
		if err != nil {
			panic(err)
		}

		full := cCtx.Bool("full")
		db := delta.NewDeltaBuilder(signature, f, full)
		d := db.BuildDelta()

		mrhs, _ = d.MarshalYAML()
		err = ioutil.WriteFile(defaultDeltaFilename, mrhs, 0644)
		if err != nil {
			panic(err)
		}

		PrettyReport(d)
		return nil

	},
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
