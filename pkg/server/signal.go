package server

import (
	"fmt"
	"os"

	logger "github.com/lucheng0127/kube-eip/pkg/utils/log"
	"golang.org/x/sys/unix"
)

var handledSignals = []os.Signal{
	unix.SIGTERM,
	unix.SIGINT,
}

func handleSignal(sigChan chan os.Signal, agent *EipAgent) {
	sig := <-sigChan
	logger.Info(agent.Ctx, fmt.Sprintf("received signal: %v, stop agent", sig))
	agent.Stop()
}
