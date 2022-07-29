package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"deltagopher/delta"
	signature2 "deltagopher/signature"
)

const defaultDeltaFilename = "delta.yml"

func GetDeltaCommand() *cli.Command {
	return &cli.Command{
		Name: "delta",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "sigfile", DefaultText: defaultSigFilename},
			&cli.BoolFlag{Name: "full"},
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

			full := cCtx.Bool("full")
			db := delta.NewDeltaBuilder(signature, f, full)
			d := db.BuildDelta()

			mrhs, _ := d.MarshalYAML()
			err = ioutil.WriteFile(defaultDeltaFilename, mrhs, 0644)
			if err != nil {
				panic(err)
			}

			PrettyReport(d)

			return nil
		},
	}
}

var colorsByOp = map[string]func(format string, a ...interface{}){
	"insert": color.Green,
	"copy":   color.Blue,
	"delete": color.Red,
}

type Rep struct {
	block     *delta.SingleDelta
	operation string
}

func PrettyReport(delta *delta.Delta) {
	sorted := []*Rep{}
	for _, x := range delta.Inserted {
		sorted = append(sorted, &Rep{
			block:     x,
			operation: "insert",
		})
	}
	for _, x := range delta.Copied {
		sorted = append(sorted, &Rep{
			block:     x,
			operation: "copy",
		})
	}
	for _, x := range delta.Deleted {
		sorted = append(sorted, &Rep{
			block:     x,
			operation: "delete",
		})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].block.Start < sorted[j].block.Start
	})

	var gap []string
	for _, item := range sorted {
		fun := colorsByOp[item.operation]
		ran := " [ " + strconv.Itoa(item.block.Start) + "-" + strconv.Itoa(item.block.End) + " ] "
		if item.operation == "insert" {
			fun("+ " + ran + string(item.block.DiffBytes))
		}

		if item.operation == "delete" {
			if len(item.block.DiffBytes) != 0 {
				fun("- " + ran + string(item.block.DiffBytes))
			} else {
				for i := 0; i < item.block.End-item.block.Start; i++ {
					gap = append(gap, "-")
				}
				fun("- " + ran + strings.Join(gap, ""))
			}
		}

		if item.operation == "copy" {
			if len(item.block.DiffBytes) != 0 {
				fun("c " + ran + string(item.block.DiffBytes))
			} else {
				for i := 0; i < item.block.End-item.block.Start; i++ {
					gap = append(gap, "-")
				}
				fun("c " + ran + strings.Join(gap, ""))
			}
		}
		gap = nil
	}

}
