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
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
