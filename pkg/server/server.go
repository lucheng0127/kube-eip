package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/erikdubbelboer/gspt"
	webserver "github.com/lucheng0127/kube-eip/internal/server"
	"github.com/lucheng0127/kube-eip/pkg/manager"
	"github.com/lucheng0127/kube-eip/pkg/protoc/binding"
	"github.com/lucheng0127/kube-eip/pkg/utils/ctx"
	logger "github.com/lucheng0127/kube-eip/pkg/utils/log"
	"github.com/lucheng0127/kube-eip/pkg/utils/validator"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

type EipAgent struct {
	Port          int
	RpcSvc        binding.EipAgentServer
	Ctx           context.Context
	InternalAddrs []string
}

func setLogger(level string) {
	var logLevel logrus.Level

	switch level {
	case "info":
		logLevel = logrus.InfoLevel
	case "debug":
		logLevel = logrus.DebugLevel
	case "warn":
		logLevel = logrus.WarnLevel
	case "error":
		logLevel = logrus.ErrorLevel
	default:
		logLevel = logrus.InfoLevel
	}

	logger.SetLevel(logLevel)
}
func (agent *EipAgent) Serve() error {
	// Launch grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", agent.Port))
	if err != nil {
		return err
	}

	gSvc := grpc.NewServer()
	binding.RegisterEipAgentServer(gSvc, &GrpcServer{})

	logger.Info(agent.Ctx, fmt.Sprintf("rpc server run on port %d", agent.Port))
	gSvc.Serve(lis)
	return nil
}

func (agent *EipAgent) Stop() {
	logger.Info(agent.Ctx, "stop agent")
	// Do something clean up
	os.Exit(0)
}

func NewAgent(opts ...AgentOption) *EipAgent {
	agent := new(EipAgent)
	agent.InternalAddrs = make([]string, 2)
	agent.Ctx = ctx.NewTraceContext()

	for _, opt := range opts {
		opt(agent)
	}

	return agent
}

func Launch(cCtx *cli.Context) error {
	// Hide process arguments, since they may contain
	// too many infos
	gspt.SetProcTitle(os.Args[0] + " daemon")

	// Init agent
	gwIP := validator.ValidateIPv4(cCtx.String("gateway-ip"))
	if gwIP == nil {
		return errors.New("invalidate gateway ip address")
	}
	gwDev := cCtx.String("gateway-dev")
	bgpType := cCtx.String("bgp-type")
	eipCidr := cCtx.String("eip-cidr")
	arpPoisoning := cCtx.Bool("arp-poisoning")

	setLogger(cCtx.String("log-level"))
	agent := NewAgent(
		setListenPort(cCtx.Int("port")),
		setInternalAddrs(cCtx.StringSlice("internal-net")),
	)

	// Setup manager
	if err := manager.RegisterManagers(gwIP, gwDev, bgpType, eipCidr, arpPoisoning, agent.InternalAddrs...); err != nil {
		return err
	}

	// Signal handle
	sigChan := make(chan os.Signal, 1024)
	signal.Notify(sigChan, handledSignals...)
	go handleSignal(sigChan, agent)

	// Run web server
	wport := validator.ValidatePort(cCtx.Int("webport"))
	if wport != 0 {
		go webserver.Serve(wport)
	}

	// Serve
	return agent.Serve()
}
