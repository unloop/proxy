package main

import (
	"github.com/lavrs/proxy"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	var err error

	app := cli.NewApp()
	app.Usage = "TCP proxy server"
	app.Version = "0.1.0"

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
				From:    c.String("f"),
				To:      c.String("t"),
				Logging: c.Bool("l"),
				Auth:    []byte(c.String("p")),
			}
			server, err := proxy.NewProxyServer(cfg)
			if err != nil {
				log.Fatal(err)
			}

			err = server.Start()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = cli.ShowAppHelp(c)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
