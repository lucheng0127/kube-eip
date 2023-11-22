package iface

import (
	"net"

	"github.com/vishvananda/netlink"
)

func AddAddrToIface(ifaceName string, addr *net.IPNet) error {
	link, err := netlink.LinkByName(ifaceName)
	if err != nil {
		return err
	}

	ifaceAddr := &netlink.Addr{IPNet: addr}

	if err = netlink.AddrAdd(link, ifaceAddr); err != nil {
		return err
	}

	return nil
}
