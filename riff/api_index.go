package riff

import (
	"github.com/gimke/cart"
	"github.com/gimke/riff/common"
	"fmt"
)

func Index(r *cart.Router) {
	r.Route("/").GET(func(c *cart.Context) {
		c.Redirect(302, "/api")
	})
}

func ApiIndex(r *cart.Router) {
	a:= Api(0)
	r.Route("/").GET(a.version)
	r.Route("/version").GET(a.version)
	r.Route("/snap").GET(a.snap)

}
type Api int

func (a Api) version(c *cart.Context) {
	version := fmt.Sprintf("Riff version %s, build %s-%s", common.Version, common.GitBranch, common.GitSha)
	c.JSON(200, cart.H{"version": version})
}

func (a Api) snap(c *cart.Context) {
	c.JSON(200, cart.H{"snapshot": riffServer.SnapShot})
}