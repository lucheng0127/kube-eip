package manager

import (
	"fmt"
	"net"

	ectx "github.com/lucheng0127/kube-eip/pkg/utils/ctx"
	logger "github.com/lucheng0127/kube-eip/pkg/utils/log"
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

func (mgr *EipMgr) addPolicyRoute() error {
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
	ctx := ectx.NewTraceContext()
	md, err := parseMD(mgr.ExternalIP.String())
	if err != nil {
		logger.Error(ctx, err.Error())
		return 0, err
	}

	if md == nil {
		md = new(EipMetadata)
	}

	if md.Status == MD_STATUS_FINISHED {
		logger.Debug(ctx, fmt.Sprintf("eip %s alrady applyed, do nothing", mgr.ExternalIP.String()))
		return 0, nil
	}

	if err := mgr.addNat(); err != nil {
		return 2, err
	}

	if err := mgr.addPolicyRoute(); err != nil {
		return 3, err
	}

	if err := mgr.addBgpRoute(); err != nil {
		return 4, err
	}

	md.ExternalIP = mgr.ExternalIP.String()
	md.InternalIP = mgr.InternalIP.String()
	md.Status = MD_STATUS_FINISHED
	md.Phase = 4

	if err := md.dumpMD(); err != nil {
		return 5, err
	}

	return 0, nil
}

func (mgr *EipMgr) UnbindEip() (int, error) {
	ctx := ectx.NewTraceContext()
	md, err := parseMD(mgr.ExternalIP.String())
	if err != nil {
		logger.Error(ctx, err.Error())
		return 0, err
	}

	if md == nil {
		logger.Debug(ctx, fmt.Sprintf("eip %s alrady cleaned, do nothing", mgr.ExternalIP.String()))
		return 0, nil
	}

	if err := mgr.deleteNat(); err != nil {
		return 2, err
	}

	if err := mgr.deleteRoutesAndTable(); err != nil {
		return 3, err
	}

	if err := mgr.deleteBgpRoute(); err != nil {
		return 4, err
	}

	if err := md.deleteMD(); err != nil {
		return 5, err
	}

	return 0, nil
}
