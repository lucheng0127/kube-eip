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
	if err := RegisterIPSetMgr(); err != nil {
		logger.Error(ctx, fmt.Sprintf("registry ipset mgr: %s", err.Error()))
		return err
	}

	err := SetupIpset("k8s_internal_net", internal_net...)
	if err != nil && !errhandle.IsExistError(err) {
		logger.Error(ctx, fmt.Sprintf("setup ipset: %s", err.Error()))
		return err
	}

	if err := RegisterRouteMgr(gwIP, gwDev, internal_net); err != nil {
		logger.Error(ctx, fmt.Sprintf("registry route mgr: %s", err.Error()))
		return err
	}

	err = RouteMgr.SetupRoute()
	if err != nil && !errhandle.IsRouteExistError(err) {
		logger.Error(ctx, fmt.Sprintf("setup policy route: %s", err.Error()))
		return err
	}

	return nil
}
