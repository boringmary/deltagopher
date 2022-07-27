package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	signature "deltagopher/signature"
	"github.com/urfave/cli/v2"
)

const defaultSigFilename = "signature.yml"

func GetSignatureCommand() *cli.Command {
	return &cli.Command{
		Name:        "signature",
		Usage:       "do the doo",
		UsageText:   "doo - does the dooing",
		Description: "no really, there is a lot of dooing to be done",
		ArgsUsage:   "[arrgh]",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "block-size"},
			&cli.StringFlag{Name: "filename"},
		},
		BashComplete: func(cCtx *cli.Context) {
			fmt.Fprintf(cCtx.App.Writer, "lipstick\nkiss\nme\nlipstick\nringo\n")
		},
		SkipFlagParsing: false,
		HideHelp:        false,
		Hidden:          false,
		HelpName:        "doo!",
		Action: func(cCtx *cli.Context) error {
			f1, err := os.Open(cCtx.Value("filename").(string))
			defer f1.Close()
			sb := signature.NewSignatureBuilder(f1, cCtx.Value("block-size").(int))
			signature := sb.BuildSignature()

			// Try marshal and save
			mrhs, _ := signature.MarshalYAML()
			err = ioutil.WriteFile(defaultSigFilename, mrhs, 0644)
			if err != nil {
				panic(err)
			}
			return nil
		},
	}
}
