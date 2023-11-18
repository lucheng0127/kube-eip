package main

import (
	"fmt"
	"os"

	"github.com/lucheng0127/kube-eip/pkg/client"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "eipctl",
		Action: client.Launch,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "target",
				Value: "127.0.0.1:6127",
				Usage: "rpc server address, default 127.0.0.1:6127",
			},
			&cli.StringFlag{
				Name:     "eip-ip",
				Usage:    "eip ip address",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "vmi-ip",
				Usage:    "vmi ip address",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "action",
				Usage:    "action, bind or unbind",
				Required: true,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
