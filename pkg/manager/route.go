package manager

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/lucheng0127/kube-eip/pkg/utils/ctx"
	logger "github.com/lucheng0127/kube-eip/pkg/utils/log"
	"github.com/vishvananda/netlink"
)

const (
	RouteTableIdx  int    = 1001
	RouteTableName string = "eip_route"
	RouteTableFile string = "/etc/iproute2/rt_tables"
	CNIDevName     string = "cni0"
)

type RouteManager interface {
	SetupRoute() error
}

type PolicyRouteMgr struct {
	InternalAddrs []string
	ExternalGWIP  net.IP
	ExternalGWDev string
}

var RouteMgr RouteManager

func newRouteManager(gwIP net.IP, gwDev string, internalAddrs []string) (RouteManager, error) {
	routeMgr := new(PolicyRouteMgr)
	routeMgr.InternalAddrs = internalAddrs
	routeMgr.ExternalGWIP = gwIP
	routeMgr.ExternalGWDev = gwDev

	ctx := ctx.NewTraceContext()
	fr, err := os.ReadFile(RouteTableFile)

	if err != nil {
		return nil, err
	}

	if strings.Contains(string(fr), RouteTableName) {
		logger.Debug(ctx, "eip route table exist, do nothing")
		return routeMgr, nil
	}

	rw, err := os.OpenFile(RouteTableFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error(ctx, err.Error())
	}
	defer rw.Close()

	if _, err := rw.WriteString(fmt.Sprintf("%d    %s\n", RouteTableIdx, RouteTableName)); err != nil {
		logger.Error(ctx, err.Error())
	}

	logger.Info(ctx, "create eip route table succeed")

	return routeMgr, nil
}

func RegisterRouteMgr(gwIP net.IP, gwDev string, internalAddrs []string) error {
	mgr, err := newRouteManager(gwIP, gwDev, internalAddrs)
	if err != nil {
		return err
	}

	RouteMgr = mgr
	return nil
}

func (mgr *PolicyRouteMgr) defaultRoute() (*netlink.Route, error) {
	route := new(netlink.Route)

	_, defaultNet, err := net.ParseCIDR("0.0.0.0/0")
	if err != nil {
		return nil, err
	}

	link, err := netlink.LinkByName(mgr.ExternalGWDev)
	if err != nil {
		return nil, err
	}

	route.Dst = defaultNet
	route.Gw = mgr.ExternalGWIP
	route.LinkIndex = link.Attrs().Index
	route.Table = RouteTableIdx

	return route, nil
}

func (mgr *PolicyRouteMgr) internalRoute(addr string) (*netlink.Route, error) {
	route := new(netlink.Route)
	var gwIP net.IP

	_, internalNet, err := net.ParseCIDR(addr)
	if err != nil {
		return nil, err
	}

	iface, err := net.InterfaceByName(CNIDevName)
	if err != nil {
		return nil, err
	}

	ifaceAddrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	for _, ifaceAddr := range ifaceAddrs {
		if v4Addr := ifaceAddr.(*net.IPNet).IP.To4(); v4Addr != nil {
			gwIP = v4Addr
		}
	}

	route.Dst = internalNet
	route.Gw = gwIP
	route.LinkIndex = iface.Index
	route.Table = RouteTableIdx

	return route, nil
}

func (mgr *PolicyRouteMgr) SetupRoute() error {
	dRoute, err := mgr.defaultRoute()
	if err != nil {
		return err
	}

	if err := netlink.RouteAdd(dRoute); err != nil {
		return err
	}

	for _, addr := range mgr.InternalAddrs {
		iRoute, err := mgr.internalRoute(addr)
		if err != nil {
			return err
		}

		if err := netlink.RouteAdd(iRoute); err != nil {
			return err
		}
	}

	return nil
}
