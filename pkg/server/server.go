package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/lucheng0127/kube-eip/pkg/protoc/binding"
	"github.com/lucheng0127/kube-eip/pkg/utils/ctx"
	logger "github.com/lucheng0127/kube-eip/pkg/utils/log"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

type EipAgent struct {
	Port   int
	RpcSvc binding.EipAgentServer
	Ctx    context.Context
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
	agent.Ctx = ctx.NewTraceContext()

	for _, opt := range opts {
		opt(agent)
	}

	return agent
}

func Launch(cCtx *cli.Context) error {
	// Init agent
	setLogger(cCtx.String("log-level"))
	agent := NewAgent(
		SetListenPort(cCtx.Int("port")),
	)

	// Signal handle
	sigChan := make(chan os.Signal, 1024)
	signal.Notify(sigChan, handledSignals...)
	go handleSignal(sigChan, agent)

	// Serve
	return agent.Serve()
}
