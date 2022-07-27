package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"deltagopher/delta"
	signature2 "deltagopher/signature"
	"github.com/urfave/cli/v2"
)

const defaultDeltaFilename = "delta.yml"

func GetDeltaCommand() *cli.Command {
	return &cli.Command{
		Name:        "delta",
		Usage:       "do the doo",
		UsageText:   "doo - does the dooing",
		Description: "no really, there is a lot of dooing to be done",
		ArgsUsage:   "[arrgh]",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "sigfile"},
		},
		BashComplete: func(cCtx *cli.Context) {
			fmt.Fprintf(cCtx.App.Writer, "lipstick\nkiss\nme\nlipstick\nringo\n")
		},
		SkipFlagParsing: false,
		HideHelp:        false,
		Hidden:          false,
		HelpName:        "doo!",
		Action: func(cCtx *cli.Context) error {
			sf, err := ioutil.ReadFile(cCtx.Value("sigfile").(string))
			if err != nil {
				panic(err)
			}
			signature, err := signature2.UnmarshalYAML(sf)
			if err != nil {
				panic(err)
			}

			f, err := os.Open(cCtx.Args().First())
			if err != nil {
				panic(err)
			}
			db := delta.NewDeltaBuilder(signature, f)
			d := db.BuildDelta()

			mrhs, _ := d.MarshalYAML()
			err = ioutil.WriteFile(defaultDeltaFilename, mrhs, 0644)
			if err != nil {
				panic(err)
			}

			j, err := json.MarshalIndent(*d, "", "  ")
			if err != nil {
				log.Fatalf(err.Error())
			}
			fmt.Printf("Report: \n %s\n", string(j))

			return nil
		},
	}
}
