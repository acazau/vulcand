package command

import (
	"github.com/mailgun/vulcand/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/mailgun/vulcand/backend"
)

func NewLocationCommand(cmd *Command) cli.Command {
	return cli.Command{
		Name:  "location",
		Usage: "Operations with vulcan locations",
		Subcommands: []cli.Command{
			{
				Name:  "show",
				Usage: "Show location details",
				Flags: []cli.Flag{
					cli.StringFlag{Name: "id", Usage: "id"},
					cli.StringFlag{Name: "host", Usage: "parent host"},
				},
				Action: cmd.printLocationAction,
			},
			{
				Name:  "add",
				Usage: "Add a new location to host",
				Flags: append([]cli.Flag{
					cli.StringFlag{Name: "id", Usage: "id, autogenerated if empty"},
					cli.StringFlag{Name: "host", Usage: "parent host"},
					cli.StringFlag{Name: "path", Usage: " path, will be matched against request's path"},
					cli.StringFlag{Name: "upstream, up", Usage: "upstream id"},
				}, locationOptions()...),
				Action: cmd.addLocationAction,
			},
			{
				Name:   "rm",
				Usage:  "Remove a location from host",
				Action: cmd.deleteLocationAction,
				Flags: []cli.Flag{
					cli.StringFlag{Name: "id", Usage: "id"},
					cli.StringFlag{Name: "host", Usage: "parent host"},
				},
			},
			{
				Name:   "set_upstream",
				Usage:  "Update location upstream",
				Action: cmd.locationUpdateUpstreamAction,
				Flags: []cli.Flag{
					cli.StringFlag{Name: "id", Usage: "id"},
					cli.StringFlag{Name: "host", Usage: "parent host"},
					cli.StringFlag{Name: "up", Usage: "new upstream id"},
				},
			},
			{
				Name:   "set_options",
				Usage:  "Update location options",
				Action: cmd.locationUpdateOptionsAction,
				Flags: append([]cli.Flag{
					cli.StringFlag{Name: "id", Usage: "id, autogenerated if empty"},
					cli.StringFlag{Name: "host", Usage: "parent host"}},
					locationOptions()...),
			},
		},
	}
}

func (cmd *Command) printLocationAction(c *cli.Context) {
	location, err := cmd.client.GetLocation(c.String("host"), c.String("id"))
	if err != nil {
		cmd.printError(err)
		return
	}
	cmd.printLocation(location)
}

func (cmd *Command) addLocationAction(c *cli.Context) {
	options, err := getOptions(c)
	if err != nil {
		cmd.printError(err)
		return
	}
	cmd.printStatus(cmd.client.AddLocationWithOptions(c.String("host"), c.String("id"), c.String("path"), c.String("up"), options))
}

func (cmd *Command) locationUpdateUpstreamAction(c *cli.Context) {
	cmd.printStatus(cmd.client.UpdateLocationUpstream(c.String("host"), c.String("id"), c.String("up")))
}

func (cmd *Command) deleteLocationAction(c *cli.Context) {
	cmd.printStatus(cmd.client.DeleteLocation(c.String("host"), c.String("id")))
}

func (cmd *Command) locationUpdateOptionsAction(c *cli.Context) {
	options, err := getOptions(c)
	if err != nil {
		cmd.printError(err)
		return
	}
	cmd.printStatus(cmd.client.UpdateLocationOptions(c.String("host"), c.String("id"), options))
}

func getOptions(c *cli.Context) (backend.LocationOptions, error) {
	o := backend.LocationOptions{}

	o.Timeouts.Read = c.Duration("readTimeout").String()
	o.Timeouts.Dial = c.Duration("dialTimeout").String()
	o.Timeouts.TlsHandshake = c.Duration("handshakeTimeout").String()

	o.KeepAlive.Period = c.Duration("keepAlivePeriod").String()
	o.KeepAlive.MaxIdleConnsPerHost = c.Int("maxIdleConns")

	o.Limits.MaxMemBodyBytes = int64(c.Int("maxMemBodyKB") * 1024)
	o.Limits.MaxBodyBytes = int64(c.Int("maxBodyKB") * 1024)

	o.FailoverPredicate = c.String("failoverPredicate")
	o.Hostname = c.String("forwardHost")
	o.TrustForwardHeader = c.Bool("trustForwardHeader")

	return o, nil
}

func locationOptions() []cli.Flag {
	return []cli.Flag{
		// Timeouts
		cli.DurationFlag{Name: "readTimeout", Usage: "read timeout"},
		cli.DurationFlag{Name: "dialTimeout", Usage: "dial timeout"},
		cli.DurationFlag{Name: "handshakeTimeout", Usage: "TLS handshake timeout"},

		// Keep-alive parameters
		cli.StringFlag{Name: "keepAlivePeriod", Usage: "keep-alive period"},
		cli.IntFlag{Name: "maxIdleConns", Usage: "maximum idle connections per host"},

		// Location limits
		cli.IntFlag{Name: "maxMemBodyKB", Usage: "maximum request size to cache in memory, in KB"},
		cli.IntFlag{Name: "maxBodyKB", Usage: "maximum request size to allow for a location, in KB"},

		// Misc options
		cli.StringFlag{Name: "failoverPredicate", Usage: "predicate that defines cases when failover is allowed"},
		cli.StringFlag{Name: "forwardHost", Usage: "hostname to set when forwarding a request"},
		cli.BoolFlag{Name: "trustForwardHeader", Usage: "allows copying X-Forwarded-For header value from the original request"},
	}
}
