package server

type AgentOption func(*EipAgent)

func SetListenPort(port int) AgentOption {
	return func(agent *EipAgent) {
		agent.Port = port
	}
}
