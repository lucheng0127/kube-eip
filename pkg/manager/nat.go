package manager

import (
	"fmt"
	"net"
	"strings"

	"github.com/lucheng0127/kube-eip/pkg/utils/cmd"
	"github.com/lucheng0127/kube-eip/pkg/utils/errhandle"
)

const (
	CHAIN_PREROUTING  string = "KUBE-EIP-PREROUTING"
	CHAIN_POSTROUTING string = "KUBE-EIP-POSTROUTING"
)

type NatManager interface {
	SetupChains(string, string) error
	AddPreroutingRule(net.IP, net.IP) error
	AddPostroutingRule(net.IP, net.IP, string) error
	DeletePreroutingRule(net.IP, net.IP) error
	DeletePostroutingRule(net.IP, net.IP, string) error
}

type IptablesMgr struct {
	mgr *cmd.CmdMgr
}

var IptablesNatMgr NatManager

func newNatManager() (NatManager, error) {
	iptablesMgr := new(IptablesMgr)
	mgr, err := cmd.NewCmdMgr("iptables")
	if err != nil {
		return nil, err
	}

	iptablesMgr.mgr = mgr
	return iptablesMgr, nil
}

func RegisterIptablesMgr() error {
	mgr, err := newNatManager()
	if err != nil {
		return err
	}

	IptablesNatMgr = mgr
	return nil
}

func (iptables *IptablesMgr) SetupChains(eipset, vmiSet string) error {
	subcmd := strings.Split(fmt.Sprintf("-t nat -N %s", CHAIN_PREROUTING), " ")
	_, err := iptables.mgr.Execute(subcmd...)
	if err != nil && !errhandle.IsExistError(err) {
		return err
	}

	subcmd = strings.Split(fmt.Sprintf("-t nat -N %s", CHAIN_POSTROUTING), " ")
	_, err = iptables.mgr.Execute(subcmd...)
	if err != nil && !errhandle.IsExistError(err) {
		return err
	}

	subcmd = strings.Split("-t nat -S PREROUTING", " ")
	currentPreroutingRules, err := iptables.mgr.Execute(subcmd...)
	if err != nil {
		return err
	}

	if !strings.Contains(currentPreroutingRules, CHAIN_PREROUTING) {
		subcmd = strings.Split(fmt.Sprintf("-t nat -I PREROUTING 1 -m set --match-set %s dst -j %s", eipset, CHAIN_PREROUTING), " ")
		_, err = iptables.mgr.Execute(subcmd...)
		if err != nil && !errhandle.IsExistError(err) {
			return err
		}
	}

	subcmd = strings.Split("-t nat -S POSTROUTING", " ")
	currentPostroutingRules, err := iptables.mgr.Execute(subcmd...)
	if err != nil {
		return err
	}

	if !strings.Contains(currentPostroutingRules, CHAIN_POSTROUTING) {
		subcmd = strings.Split(fmt.Sprintf("-t nat -I POSTROUTING 1 -m set --match-set %s src -j %s", vmiSet, CHAIN_POSTROUTING), " ")
		_, err = iptables.mgr.Execute(subcmd...)
		if err != nil && !errhandle.IsExistError(err) {
			return err
		}
	}

	return nil
}

func (iptables *IptablesMgr) AddPreroutingRule(eip, vmiIp net.IP) error {
	subcmd := strings.Split(fmt.Sprintf("-t nat -A %s -d %s -j DNAT --to-destination %s", CHAIN_PREROUTING, eip.String(), vmiIp.String()), " ")
	_, err := iptables.mgr.Execute(subcmd...)
	if err != nil && !errhandle.IsExistError(err) {
		return err
	}

	return nil
}

func (iptables *IptablesMgr) AddPostroutingRule(eip, vmiIp net.IP, internalSet string) error {
	subcmd := strings.Split(fmt.Sprintf("-t nat -A %s -s %s -m set ! --match-set %s dst -j SNAT --to %s", CHAIN_POSTROUTING, vmiIp.String(), internalSet, eip.String()), " ")
	_, err := iptables.mgr.Execute(subcmd...)
	if err != nil && !errhandle.IsExistError(err) {
		return err
	}

	return nil
}

func (iptables *IptablesMgr) DeletePreroutingRule(eip, vmiIp net.IP) error {
	subcmd := strings.Split(fmt.Sprintf("-t nat -D %s -d %s -j DNAT --to-destination %s", CHAIN_PREROUTING, eip.String(), vmiIp.String()), " ")
	_, err := iptables.mgr.Execute(subcmd...)
	if err != nil && !errhandle.IsExistError(err) {
		return err
	}

	return nil
}

func (iptables *IptablesMgr) DeletePostroutingRule(eip, vmiIp net.IP, internalSet string) error {
	subcmd := strings.Split(fmt.Sprintf("-t nat -D %s -s %s -m set ! --match-set %s dst -j SNAT --to %s", CHAIN_POSTROUTING, vmiIp.String(), internalSet, eip.String()), " ")
	_, err := iptables.mgr.Execute(subcmd...)
	if err != nil && !errhandle.IsExistError(err) {
		return err
	}

	return nil
}

//TODO(shawnlu): Add chain inject filter table FORWARD, provent traffice match CNI-ISOLATION-STAGE-1
