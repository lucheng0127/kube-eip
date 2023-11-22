package manager

import (
	"fmt"
	"net"

	"github.com/lucheng0127/kube-eip/pkg/utils/ctx"
	"github.com/lucheng0127/kube-eip/pkg/utils/errhandle"
	logger "github.com/lucheng0127/kube-eip/pkg/utils/log"
)

func RegisterManagers(gwIP net.IP, gwDev string, internal_net ...string) error {
	ctx := ctx.NewTraceContext()

	// Registry ipset mgr
	if err := RegisterIPSetMgr(); err != nil {
		logger.Error(ctx, fmt.Sprintf("registry ipset mgr: %s", err.Error()))
		return err
	}

	err := IpsetMgr.SetupIpset(SET_INTERNAL, "hash:net", internal_net...)
	if err != nil && !errhandle.IsExistError(err) {
		logger.Error(ctx, fmt.Sprintf("setup ipset: %s", err.Error()))
		return err
	}

	setName := []string{SET_EIP, SET_VMI}
	for _, name := range setName {
		err := IpsetMgr.SetupIpset(name, "hash:ip")
		if err != nil && !errhandle.IsExistError(err) {
			logger.Error(ctx, fmt.Sprintf("setup ipset: %s", err.Error()))
			return err
		}
	}

	// Registry route mgr
	if err := RegisterRouteMgr(gwIP, gwDev, internal_net); err != nil {
		logger.Error(ctx, fmt.Sprintf("registry route mgr: %s", err.Error()))
		return err
	}

	err = RouteMgr.SetupRoute()
	if err != nil && !errhandle.IsRouteExistError(err) {
		logger.Error(ctx, fmt.Sprintf("setup policy route: %s", err.Error()))
		return err
	}

	// Registry nat mgr
	if err := RegisterIptablesMgr(); err != nil {
		logger.Error(ctx, fmt.Sprintf("registry nat mgr: %s", err.Error()))
		return err
	}

	err = IptablesNatMgr.SetupChains(SET_EIP, SET_VMI)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("setup iptables chains: %s", err.Error()))
		return err
	}

	return nil
}
