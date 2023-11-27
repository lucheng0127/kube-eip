package manager

import (
	"context"
	"fmt"
	"net"

	ectx "github.com/lucheng0127/kube-eip/pkg/utils/ctx"
	"github.com/lucheng0127/kube-eip/pkg/utils/errhandle"
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

func (mgr *EipMgr) addToSet(ctx context.Context) error {
	if err := mgr.IPSetMgr.AddIPToSet(SET_EIP, mgr.ExternalIP); err != nil && !errhandle.IsIpsetExistError(err) {
		logger.Error(ctx, fmt.Sprintf("add eip %s to ipset %s: %s", mgr.ExternalIP.String(), SET_EIP, err.Error()))
		return err
	}

	if err := mgr.IPSetMgr.AddIPToSet(SET_VMI, mgr.InternalIP); err != nil && !errhandle.IsIpsetExistError(err) {
		logger.Error(ctx, fmt.Sprintf("add vmi ip %s to ipset %s: %s", mgr.InternalIP.String(), SET_VMI, err.Error()))
		return err
	}
	return nil
}

func (mgr *EipMgr) deleteFromSet(ctx context.Context) error {
	if err := mgr.IPSetMgr.DeleteFromSet(SET_EIP, mgr.ExternalIP); err != nil && !errhandle.IsIpsetItemNotExist(err) {
		logger.Error(ctx, fmt.Sprintf("delete eip %s from ipset %s: %s", mgr.ExternalIP.String(), SET_EIP, err.Error()))
		return err
	}

	if err := mgr.IPSetMgr.DeleteFromSet(SET_VMI, mgr.InternalIP); err != nil && !errhandle.IsIpsetItemNotExist(err) {
		logger.Error(ctx, fmt.Sprintf("delete vmi ip %s from ipset %s: %s", mgr.InternalIP.String(), SET_VMI, err.Error()))
		return err
	}
	return nil
}

func (mgr *EipMgr) addNat(ctx context.Context) error {
	if err := mgr.NatMgr.AddPreroutingRule(mgr.ExternalIP, mgr.InternalIP); err != nil {
		logger.Error(ctx, fmt.Sprintf("add dnat eip %s to vmi ip %s: %s", mgr.ExternalIP.String(), mgr.InternalIP.String(), err.Error()))
		return err
	}

	if err := mgr.NatMgr.AddPostroutingRule(mgr.ExternalIP, mgr.InternalIP, SET_INTERNAL); err != nil {
		logger.Error(ctx, fmt.Sprintf("add snat cmi ip %s to eip %s: %s", mgr.InternalIP.String(), mgr.ExternalIP.String(), err.Error()))
		return err
	}
	return nil
}

func (mgr *EipMgr) deleteNat(ctx context.Context) error {
	if err := mgr.NatMgr.DeletePreroutingRule(mgr.ExternalIP, mgr.InternalIP); err != nil && !errhandle.IsIptablesRuleNotExist(err) {
		logger.Error(ctx, fmt.Sprintf("delete dnat eip %s to vmi ip %s: %s", mgr.ExternalIP.String(), mgr.InternalIP.String(), err.Error()))
		return err
	}

	if err := mgr.NatMgr.DeletePostroutingRule(mgr.ExternalIP, mgr.InternalIP, SET_INTERNAL); err != nil && !errhandle.IsIptablesRuleNotExist(err) {
		logger.Error(ctx, fmt.Sprintf("delete snat cmi ip %s to eip %s: %s", mgr.InternalIP.String(), mgr.ExternalIP.String(), err.Error()))
		return err
	}
	return nil
}

func (mgr *EipMgr) addPolicyRoute(ctx context.Context, vmiIP net.IP) error {
	if err := mgr.RouteMgr.AddEipRule(vmiIP); err != nil && !errhandle.IsNetlinkExistError(err) {
		logger.Error(ctx, fmt.Sprintf("add vmi %s eip rule: %s", vmiIP.String(), err.Error()))
		return err
	}

	return nil
}

func (mgr *EipMgr) deleteRoutesAndTable(ctx context.Context, vmiIP net.IP) error {
	if err := mgr.RouteMgr.DeleteEipRule(vmiIP); err != nil && !errhandle.IsRuleNotExist(err) {
		logger.Error(ctx, fmt.Sprintf("del vmi %s eip rule: %s", vmiIP.String(), err.Error()))
		return err
	}

	return nil
}

func (mgr *EipMgr) addEipToIface(ctx context.Context) error {
	if err := mgr.RouteMgr.AddEipToIface(mgr.ExternalIP); err != nil && !errhandle.IsNetlinkExistError(err) {
		logger.Error(ctx, fmt.Sprintf("add eip %s to hyper interface %s", mgr.ExternalIP.String(), err.Error()))
		return err
	}

	return nil
}

func (mgr *EipMgr) remoteEipFromIface(ctx context.Context) error {
	if err := mgr.RouteMgr.RemoveEipFromIface(mgr.ExternalIP); err != nil && !errhandle.IsIPNotExist(err) {
		logger.Error(ctx, fmt.Sprintf("remove eip %s from hyper interface %s", mgr.ExternalIP.String(), err.Error()))
		return err
	}

	return nil
}

func (mgr *EipMgr) addBgpRoute() error {
	return nil
}

func (mgr *EipMgr) deleteBgpRoute() error {
	return nil
}

func (mgr *EipMgr) BindEip() (int, error) {
	// Parse metadata
	ctx := ectx.NewTraceContext()
	md, err := parseMD(mgr.ExternalIP.String(), mgr.InternalIP.String())
	if err != nil {
		logger.Error(ctx, err.Error())
		return 0, err
	}

	if md == nil {
		// TODO(shawnlu) Check eip conflict

		// Init metadata
		md = new(EipMetadata)
		md.Phase = 0
		md.Status = MD_STATUS_FAILED
		md.ExternalIP = mgr.ExternalIP.String()
		md.InternalIP = mgr.InternalIP.String()
	}

	if md.Status == MD_STATUS_FINISHED {
		logger.Debug(ctx, fmt.Sprintf("eip %s alrady applyed, do nothing", mgr.ExternalIP.String()))
		return 0, nil
	}

	defer md.dumpMD()

	// Do binding phase
	for i := md.Phase; i < 5; i++ {
		switch i + 1 {
		case 1:
			// Add eip and vmi ip to ipset
			if err := mgr.addToSet(ctx); err != nil {
				md.Phase = 0
				return 1, err
			}
		case 2:
			// Add iptables rules
			if err := mgr.addNat(ctx); err != nil {
				md.Phase = 1
				return 2, err
			}
		case 3:
			// Add policy route
			if err := mgr.addPolicyRoute(ctx, mgr.InternalIP); err != nil {
				md.Phase = 2
				return 3, err
			}
		case 4:
			// Add eip to interface
			// If use arp snooping, no need do this
			if err := mgr.addEipToIface(ctx); err != nil {
				md.Phase = 3
				return 4, err
			}
		case 5:
			// Add bgp route
			if err := mgr.addBgpRoute(); err != nil {
				md.Phase = 4
				return 5, err
			}
		}
	}

	md.Phase = 5
	md.Status = MD_STATUS_FINISHED

	return 0, nil
}

func (mgr *EipMgr) UnbindEip() (int, error) {
	ctx := ectx.NewTraceContext()
	md, err := parseMD(mgr.ExternalIP.String(), mgr.InternalIP.String())
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("parse metadata file %s", err.Error()))
		return 0, err
	}

	if md == nil {
		logger.Debug(ctx, fmt.Sprintf("eip binding %s-%s alrady cleaned, do nothing", mgr.ExternalIP.String(), mgr.InternalIP.String()))
		return 0, nil
	}

	// Do unbind phase
	for i := md.Phase; i > 0; i-- {
		switch i {
		case 1:
			// Delete eip and vmi ip from ipset
			if err := mgr.deleteFromSet(ctx); err != nil {
				return 1, err
			}
		case 2:
			// Delete iptables rules
			if err := mgr.deleteNat(ctx); err != nil {
				return 2, err
			}
		case 3:
			// Delete policy route
			if err := mgr.deleteRoutesAndTable(ctx, mgr.InternalIP); err != nil {
				return 3, err
			}
		case 4:
			// Remove eip from interface
			if err := mgr.remoteEipFromIface(ctx); err != nil {
				return 4, err
			}
		case 5:
			// Delete bgp route
			if err := mgr.deleteBgpRoute(); err != nil {
				return 5, err
			}
		}
	}

	// Delete metadata file when clean up finished
	if err := md.deleteMD(); err != nil {
		logger.Error(ctx, fmt.Sprintf("delete eip %s metadata file failed", mgr.ExternalIP.String()))
		return 6, err
	}

	return 0, nil
}
