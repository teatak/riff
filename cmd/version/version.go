package version

import (
	"fmt"
	"github.com/gimke/riff/common"
)

type cmd struct {
	version string
}

func New(version string) *cmd {
	return &cmd{version: version}
}

func (c *cmd) Run(_ []string) int {
	fmt.Printf("Riff version %s, build %s-%s\n", c.version, common.GitBranch, common.GitSha)
	return 0
}

func (c *cmd) Synopsis() string {
	return "Prints the Riff version"
}

func (c *cmd) Help() string {
	return ""
}
