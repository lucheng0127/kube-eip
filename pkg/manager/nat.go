package manager

import (
	"fmt"
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

func (mgr *IptablesMgr) SetupChains(eipset, vmiSet string) error {
	subcmd := strings.Split(fmt.Sprintf("-t nat -N %s", CHAIN_PREROUTING), " ")
	_, err := mgr.mgr.Execute(subcmd...)
	if err != nil && !errhandle.IsExistError(err) {
		return err
	}

	subcmd = strings.Split(fmt.Sprintf("-t nat -N %s", CHAIN_POSTROUTING), " ")
	_, err = mgr.mgr.Execute(subcmd...)
	if err != nil && !errhandle.IsExistError(err) {
		return err
	}

	subcmd = strings.Split(fmt.Sprintf("-t nat -I PREROUTING 1 -m set --match-set %s dst -j %s", eipset, CHAIN_PREROUTING), " ")
	_, err = mgr.mgr.Execute(subcmd...)
	if err != nil && !errhandle.IsExistError(err) {
		return err
	}

	subcmd = strings.Split(fmt.Sprintf("-t nat -I POSTROUTING 1 -m set --match-set %s src -j %s", vmiSet, CHAIN_POSTROUTING), " ")
	_, err = mgr.mgr.Execute(subcmd...)
	if err != nil && !errhandle.IsExistError(err) {
		return err
	}

	return nil
}
