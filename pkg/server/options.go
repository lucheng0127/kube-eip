package server

type AgentOption func(*EipAgent)

func setListenPort(port int) AgentOption {
	return func(agent *EipAgent) {
		agent.Port = port
	}
}

func setInternalAddrs(addrs []string) AgentOption {
	return func(mgr *EipAgent) {
		mgr.InternalAddrs = addrs
	}
}

func setSecret(sec string) AgentOption {
	return func(mgr *EipAgent) {
		mgr.Secret = sec
	}
}
