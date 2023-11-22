package server

import "net"

type AgentOption func(*EipAgent)

func setListenPort(port int) AgentOption {
	return func(agent *EipAgent) {
		agent.Port = port
	}
}

func setExternalGWIP(addr net.IP) AgentOption {
	return func(mgr *EipAgent) {
		mgr.ExternalGWIP = addr
	}
}

func setExternalGEDev(dev string) AgentOption {
	return func(mgr *EipAgent) {
		mgr.ExternalGWDev = dev
	}
}

func setInternalAddrs(addrs []string) AgentOption {
	return func(mgr *EipAgent) {
		mgr.InternalAddrs = addrs
	}
}

func setExternalBgpType(bgpType string) AgentOption {
	return func(mgr *EipAgent) {
		mgr.BgpType = bgpType
	}
}

func setEipMaskLen(len int) AgentOption {
	return func(mgr *EipAgent) {
		mgr.EipMaskLen = len
	}
}
