package riff

import (
	"github.com/gimke/cart"
	"time"
	"log"
)

func Logger() cart.Handler {

	return func(c *cart.Context, next cart.Next) {
		start := time.Now()
		path := c.Request.URL.Path
		next()
		end := time.Now()
		latency := end.Sub(start)
		method := c.Request.Method
		clientIP := c.ClientIP()
		statusCode := c.Response.Status()

		log.Printf("[CART] status:%d latency:%v ip:%s method:%s path:%s\n",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
