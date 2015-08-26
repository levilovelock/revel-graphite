package rg

import (
	"fmt"

	graphite "github.com/marpaia/graphite-golang"
	"github.com/revel/revel"
	"github.com/robfig/cron"
)

var (
	g    *graphite.Graphite
	host string
	port int
)

func init() {
	revel.OnAppStart(func() {
		host, hostFound := revel.Config.String("revel-graphite.host")

		if !hostFound {
			revel.ERROR.Println("Expecting revel-graphite.host to be set in config, but it wasn't!")
		}

		port, portFound := revel.Config.Int("revel-graphite.port")
		if !portFound {
			revel.ERROR.Println("Expecting revel-graphite.host to be set in config, but it wasn't!")
		}

		var err error
		g, err = graphite.NewGraphite(host, port)
		if err != nil {
			revel.ERROR.Printf("Failed to connect to graphite server with host: '%s' and port: '%d'.\n", host, port)
			return
		} else {
			revel.INFO.Printf("Connected to graphite server '%s:%d' succesfully.\n", host, port)
		}

		c := cron.New()
		c.AddFunc("@every 10s", sessionsMetric)
		c.Start()
	})
}

func sessionsMetric() {
	/* TODO: Fix it so it sends an actual metric
	 * One way to do this would be to create a filter for the plugin which we tell teh users to
	 * add at the end (or at least after sessions for this exact case), then that filter will
	 * check for a session and increment a global counter. Probably the counter will be reset every 10s or
	 * 30s or 60s and then we actually publish the metric to graphite then and there with the past minutes
	 * data or whatever. We would probably need to have a session map to not double count the same sessions
	 * too. If the user wasn't using the default session filter we could make session (and other standard
	 * metrics) opt-out via module-specific-configuration. We could probably do most web-metrics via this
	 * filter approach.
	 */

	err := g.SimpleSend(fmt.Sprintf("revel.%s.sessions", revel.AppName), "0")
	if err != nil {
		revel.ERROR.Printf("Cannot send session stats to graphite server '%s:%d'.\n", host, port)
	}
}
