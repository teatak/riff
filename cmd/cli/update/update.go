package update

import (
	"flag"
	"fmt"
	"github.com/teatak/riff/api"
	"github.com/teatak/riff/cmd/cli"
	"github.com/teatak/riff/cmd/cli/daem"
	"github.com/teatak/riff/cmd/cli/quit"
	"github.com/teatak/riff/common"
	"github.com/teatak/riff/git"
	"math"
	"runtime"
	"strings"
)

const help = `Usage: update riff
`

type cmd struct {
	flags *flag.FlagSet
}

func New() *cmd {
	c := &cmd{}
	c.init()
	return c
}

func (c *cmd) init() {
	c.flags = flag.NewFlagSet("update", flag.ContinueOnError)

	c.flags.Usage = func() {
		fmt.Println(c.Help())
	}
}

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}
	c.Update()
	return 0
}

func (c *cmd) Update() {
	//get version from github
	currentVersion := common.Version

	client := git.GithubClient("", "https://github.com/teatak/riff")
	version, _, err := client.GetRelease("latest")

	if err != nil {
		fmt.Println(err)
	} else {
		if version != "v"+currentVersion {
			fmt.Printf("find new version %v to be update [Y/N]:", version)
			var input string
			_, _ = fmt.Scanln(&input)
			if strings.ToLower(input) == "y" {
				fmt.Print("downloading...")
				zipFile := runtime.GOOS + "_" + runtime.GOARCH + ".zip"
				downloadUrl := "https://github.com/teatak/riff/releases/download/" + version + "/" + zipFile
				file := common.BinDir + "/update/riff/" + version + "/" + zipFile
				dir := common.BinDir
				var progress api.Progress
				progress = func(current, total int32) {
					fmt.Printf("\r%s", strings.Repeat(" ", 45))
					// Return again and print current status of download
					// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
					c := math.Round(float64(current)/1024/1024*100) / 100
					t := math.Round(float64(total)/1024/1024*100) / 100
					s := fmt.Sprintf("%vM/%vM", c, t)
					fmt.Printf("\rdownloading... %s complete", s)
				}
				if err := client.DownloadFile(file, downloadUrl, progress); err != nil {
					fmt.Println()
					fmt.Println(err)
				} else {
					//
					fmt.Println()
					if cli.GetPid() != 0 {
						//quit
						q := quit.New()
						q.Run([]string{})
						//copy
						if err := common.Unzip(file, dir, false); err != nil {
							fmt.Println(err)
						}
						//run
						s := daem.New()
						s.Run([]string{})
					} else {
						//copy
						if err := common.Unzip(file, dir, false); err != nil {
							fmt.Println(err)
						}
					}
				}
			}

		} else {
			fmt.Println("riff is latest version")
		}
	}
}

func (c *cmd) Synopsis() string {
	return "Update Riff"
}

func (c *cmd) Help() string {
	return help
}
