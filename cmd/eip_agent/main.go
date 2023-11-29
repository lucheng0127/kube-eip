package main

import (
	"fmt"
	"os"

	"github.com/lucheng0127/kube-eip/pkg/server"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "EipAgent",
		Action: server.Launch,
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:  "port",
				Value: 6127,
				Usage: "agent port that rpc listen on",
			},
			&cli.StringFlag{
				Name:  "log-level",
				Value: "info",
				Usage: "log level, default info",
			},
			&cli.StringSliceFlag{
				Name:  "internal-net",
				Usage: "networks that exclude from nat",
			},
			&cli.StringFlag{
				Name:     "gateway-ip",
				Usage:    "externel network gateway ip address",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "gateway-dev",
				Usage:    "externel network gateway device",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "bgp-type",
				Value: "none",
				Usage: "bgp manager type, default is none, gobgp is avaliable",
			},
			&cli.StringFlag{
				Name:     "eip-cidr",
				Required: true,
				Usage:    "eip network cidr",
			},
			&cli.BoolFlag{
				Name:  "arp-poisoning",
				Value: false,
				Usage: "whether use arp poisoning to make a arp reply for eip, default is false, when enable will not add eip to external network gateway device",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
