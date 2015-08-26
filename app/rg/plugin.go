package rg

import (
	"fmt"

	graphite "github.com/marpaia/graphite-golang"
	"github.com/revel/revel"
)

var (
	g    graphite.Graphite
	host string
	port int
)

func init() {
	revel.OnAppStart(func() {
		host, hostFound := revel.Config.String("revel-graphite.host")

		if !hostFound {
			fmt.Println("Expecting revel-graphite.host to be set in config, but it wasn't!")
		}

		port, portFound := revel.Config.Int("revel-graphite.port")
		if !portFound {
			fmt.Println("Expecting revel-graphite.host to be set in config, but it wasn't!")
		}

		g, err := graphite.NewGraphite(host, port)
		if err != nil {
			fmt.Printf("Failed to connect to graphite server with host: '%s' and port: '%d'.\n", host, port)
			return
		}

		g.SimpleSend("test.stat", "20")
	})
}
