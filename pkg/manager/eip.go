package manager

import (
	"net"
)

type EipManager interface {
	BindEip() (int, error)
	UnbindEip() (int, error)
}

type EipMgr struct {
	IPSetMgr IpsetManager
	NatMgr   NatManager
	RouteMgr RouteManager
	BgpMgr   BgpManager

	ExternalIP net.IP
	InternalIP net.IP
}

func (mgr *EipMgr) addNat() error {
	return nil
}

func (mgr *EipMgr) deleteNat() error {
	return nil
}

func (mgr *EipMgr) addRoutesAndTable() error {
	return nil
}

func (mgr *EipMgr) deleteRoutesAndTable() error {
	return nil
}

func (mgr *EipMgr) addBgpRoute() error {
	return nil
}

func (mgr *EipMgr) deleteBgpRoute() error {
	return nil
}

func (mgr *EipMgr) BindEip() (int, error) {
	// TODO(shawnlu): Due to operator will update CRD phase,
	// so maybe bind or unbind will be call multi times.
	// if exist no need bind multi times.
	if err := mgr.addNat(); err != nil {
		return 2, err
	}

	if err := mgr.addRoutesAndTable(); err != nil {
		return 3, err
	}

	if err := mgr.addBgpRoute(); err != nil {
		return 4, err
	}

	return 0, nil
}

func (mgr *EipMgr) UnbindEip() (int, error) {
	// TODO(shawnlu): Due to operator will update CRD phase,
	// so maybe bind or unbind will be call multi times.
	// if exist no need bind multi times.
	if err := mgr.deleteNat(); err != nil {
		return 2, err
	}

	if err := mgr.deleteRoutesAndTable(); err != nil {
		return 3, err
	}

	if err := mgr.deleteBgpRoute(); err != nil {
		return 4, err
	}

	return 4, nil
}
