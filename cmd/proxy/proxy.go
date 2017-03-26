package main

import (
	"github.com/lavrs/proxy"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Usage = "TCP proxy server"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "f, from",
			Usage: "set proxy server port",
		},

		cli.StringFlag{
			Name:  "t, to",
			Usage: "set proxy server redirect port",
		},

		cli.StringFlag{
			Name:  "p, pass",
			Usage: "set proxy server password",
		},

		cli.BoolFlag{
			Name:  "l, log",
			Usage: "enable logging",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.NArg() == 0 {
			cfg := proxy.Proxy{
				From:     c.String("f"),
				To:       c.String("t"),
				Logging:  c.Bool("l"),
				Password: []byte(c.String("p")),
			}
			server := proxy.NewProxyServer(cfg)

			server.Start()
		} else {
			cli.ShowAppHelp(c)
		}
	}

	app.Run(os.Args)
}
