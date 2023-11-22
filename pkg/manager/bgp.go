package manager

import (
	"fmt"
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
	AsignEipToIface(string, net.IP) error
}

type GobgpMgr struct {
	EipMaskLen int
}

type NoBgpMgr struct {
	EipMaskLen int
}

var BgpMgr BgpManager

// TODO(shawnlu): Implement eip arp proxy

func newBgpManager(bgpType string, masklen int) (BgpManager, error) {
	switch bgpType {
	case BGP_TYPE_GOBGP:
		bgpMgr := new(GobgpMgr)
		bgpMgr.EipMaskLen = masklen
		return bgpMgr, nil
	default:
		bgpMgr := new(NoBgpMgr)
		bgpMgr.EipMaskLen = masklen
		return bgpMgr, nil
	}
}

func RegisterBgpManager(bgpType string, masklen int) error {
	if masklen < 0 || masklen > 32 {
		return fmt.Errorf("invalidate eip marklen %d", masklen)
	}

	mgr, err := newBgpManager(bgpType, masklen)
	if err != nil {
		return err
	}

	BgpMgr = mgr
	return nil
}

func (mgr *NoBgpMgr) AsignEipToIface(ifaceName string, addr net.IP) error {
	return iface.AddAddrToIface(ifaceName, &net.IPNet{IP: addr, Mask: net.CIDRMask(mgr.EipMaskLen, 32)})
}

func (mgr *GobgpMgr) AsignEipToIface(ifaceName string, addr net.IP) error {
	return iface.AddAddrToIface(ifaceName, &net.IPNet{IP: addr, Mask: net.CIDRMask(mgr.EipMaskLen, 32)})
}
