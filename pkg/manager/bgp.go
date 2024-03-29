package manager

// BgpManager use to sync ipv4 route to eip router
// There are two type of bgp manager, gobgp or nobgp
// Add a eip arp proxy or add eip to gateway device
// directly.
type BgpManager interface {
}

type GobgpMgr struct {
}

var BgpMgr BgpManager

func newBgpManager() (BgpManager, error) {
	bgpMgr := new(GobgpMgr)
	return bgpMgr, nil
}

func RegisterBgpManager() error {
	mgr, err := newBgpManager()
	if err != nil {
		return err
	}

	BgpMgr = mgr
	return nil
}
