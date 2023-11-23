package manager

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/lucheng0127/kube-eip/pkg/utils/cmd"
	"github.com/lucheng0127/kube-eip/pkg/utils/ctx"
	"github.com/lucheng0127/kube-eip/pkg/utils/errhandle"
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
	AddEipRule(net.IP) error
	DeleteEipRule(net.IP) error
}

type PolicyRouteMgr struct {
	InternalAddrs []string
	ExternalGWIP  net.IP
	ExternalGWDev string
	mgr           *cmd.CmdMgr
	eipCidr       string
}

var RouteMgr RouteManager

func newRouteManager(gwIP net.IP, gwDev string, eipCidr string, internalAddrs []string) (RouteManager, error) {
	routeMgr := new(PolicyRouteMgr)
	routeMgr.InternalAddrs = internalAddrs
	routeMgr.ExternalGWIP = gwIP
	routeMgr.ExternalGWDev = gwDev
	routeMgr.eipCidr = eipCidr

	mgr, err := cmd.NewCmdMgr("ip")
	if err != nil {
		return nil, err
	}

	routeMgr.mgr = mgr

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

func RegisterRouteMgr(gwIP net.IP, gwDev string, eipCidr string, internalAddrs []string) error {
	mgr, err := newRouteManager(gwIP, gwDev, eipCidr, internalAddrs)
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

func (mgr *PolicyRouteMgr) eipRoute() (*netlink.Route, error) {
	route := new(netlink.Route)

	iface, err := net.InterfaceByName(mgr.ExternalGWDev)
	if err != nil {
		return nil, err
	}

	_, eipNet, err := net.ParseCIDR(mgr.eipCidr)
	if err != nil {
		return nil, err
	}

	route.Dst = eipNet
	route.LinkIndex = iface.Index
	route.Table = RouteTableIdx
	return route, nil
}

func (mgr *PolicyRouteMgr) SetupRoute() error {
	dRoute, err := mgr.defaultRoute()
	if err != nil {
		return err
	}

	if err := netlink.RouteAdd(dRoute); err != nil && !errhandle.IsRouteExistError(err) {
		return err
	}

	for _, addr := range mgr.InternalAddrs {
		iRoute, err := mgr.internalRoute(addr)
		if err != nil {
			return err
		}

		if err := netlink.RouteAdd(iRoute); err != nil && !errhandle.IsRouteExistError(err) {
			return err
		}
	}

	// Add eip network route to route table eip_route
	// For icmp reply traffice, there is a ct record, when the reply
	// traffic go through mangle table in POSTROUTING, because of
	// ct record, the traffic will do nat by ct record, so will
	// not match snat rules in nat table.
	eRoute, err := mgr.eipRoute()
	if err != nil {
		return nil
	}

	if err := netlink.RouteAdd(eRoute); err != nil && !errhandle.IsRouteExistError(err) {
		return err
	}

	return nil
}

func (mgr *PolicyRouteMgr) AddEipRule(vmiIp net.IP) error {
	subcmd := strings.Split(fmt.Sprintf("rule add from %s pref %d lookup %s", vmiIp.String(), RouteTableIdx, RouteTableName), " ")
	_, err := mgr.mgr.Execute(subcmd...)
	if err != nil && !errhandle.IsExistError(err) {
		return err
	}

	return nil
}

func (mgr *PolicyRouteMgr) DeleteEipRule(vmiIp net.IP) error {
	subcmd := strings.Split(fmt.Sprintf("rule del from %s pref %d lookup %s", vmiIp.String(), RouteTableIdx, RouteTableName), " ")
	_, err := mgr.mgr.Execute(subcmd...)
	if err != nil && !errhandle.IsExistError(err) {
		return err
	}

	return nil
}
