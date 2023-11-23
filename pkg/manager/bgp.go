package manager

import (
	"net"

	"github.com/lucheng0127/kube-eip/pkg/utils/iface"
)

const (
	BGP_TYPE_GOBGP string = "gobgp"
	BGP_TYPE_NONE  string = "none"
)

// BgpManager use to sync ipv4 route to eip router
// There are two type of bgp manager, gobgp or nobgp
// Add a eip arp proxy or add eip to gateway device
// directly.
type BgpManager interface {
}

type GobgpMgr struct {
}

type NoBgpMgr struct {
}

var BgpMgr BgpManager

// TODO(shawnlu): Implement eip arp proxy

func newBgpManager(bgpType string) (BgpManager, error) {
	switch bgpType {
	case BGP_TYPE_GOBGP:
		bgpMgr := new(GobgpMgr)
		return bgpMgr, nil
	default:
		bgpMgr := new(NoBgpMgr)
		return bgpMgr, nil
	}
}

func RegisterBgpManager(bgpType string) error {
	mgr, err := newBgpManager(bgpType)
	if err != nil {
		return err
	}

	BgpMgr = mgr
	return nil
}

func (mgr *NoBgpMgr) AsignEipToIface(ifaceName string, addr net.IP) error {
	return iface.AddAddrToIface(ifaceName, &net.IPNet{IP: addr, Mask: net.CIDRMask(24, 32)})
}

func (mgr *GobgpMgr) AsignEipToIface(ifaceName string, addr net.IP) error {
	return iface.AddAddrToIface(ifaceName, &net.IPNet{IP: addr, Mask: net.CIDRMask(24, 32)})
}
