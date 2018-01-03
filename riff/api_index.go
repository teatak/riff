package riff

import (
	"fmt"
	"github.com/gimke/cart"
	"github.com/gimke/riff/common"
)

type Api struct {
	server *Server
}

func (a *Api) Index(r *cart.Router) {
	r.Route("/").GET(func(c *cart.Context) {
		c.Redirect(302, "/api")
	})
	r.Route("/api",a.apiIndex)
}

func (a *Api) apiIndex(r *cart.Router) {
	r.Route("").GET(a.version)
	r.Route("/version").GET(a.version)
	r.Route("/snap").GET(a.snap)
	r.Route("/nodes").GET(a.nodes)
}

func (a Api) version(c *cart.Context) {
	version := fmt.Sprintf("Cart version %s Riff version %s, build %s-%s", cart.Version, common.Version, common.GitBranch, common.GitSha)
	c.JSON(200, cart.H{
		"version": version,
	})
}

func (a Api) snap(c *cart.Context) {
	c.JSON(200, cart.H{
		"snapshot": a.server.SnapShot,
	})
}

func (a Api) nodes(c *cart.Context) {
	c.JSON(200, cart.H{
		"nodes": a.server.Nodes,
	})
}
