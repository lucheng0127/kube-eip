package manager

import (
	"fmt"
	"net"
	"strings"

	"github.com/lucheng0127/kube-eip/pkg/utils/cmd"
)

const (
	SET_INTERNAL string = "k8s_internal_net"
	SET_EIP      string = "kube-eip-eip"
	SET_VMI      string = "kube-eip-vmi"
)

type IpsetManager interface {
	AddSetAndEntries(string, string, ...string) error
	SetupIpset(string, string, ...string) error
	AddIPToSet(string, net.IP) error
	DeleteFromSet(string, net.IP) error
}

type CmdIpsetMgr struct {
	mgr *cmd.CmdMgr
}

var IpsetMgr IpsetManager

func newIpsetManager() (IpsetManager, error) {
	ipsetMgr := new(CmdIpsetMgr)
	mgr, err := cmd.NewCmdMgr("ipset")
	if err != nil {
		return nil, err
	}

	ipsetMgr.mgr = mgr
	return ipsetMgr, nil
}

func RegisterIPSetMgr() error {
	mgr, err := newIpsetManager()
	if err != nil {
		return err
	}

	IpsetMgr = mgr
	return nil
}

func (ipset *CmdIpsetMgr) AddSetAndEntries(name, setType string, entries ...string) error {
	subcmd := strings.Split(fmt.Sprintf("create %s %s", name, setType), " ")
	if _, err := ipset.mgr.Execute(subcmd...); err != nil {
		return err
	}

	for _, entry := range entries {
		subcmd = strings.Split(fmt.Sprintf("add %s %s", name, entry), " ")
		if _, err := ipset.mgr.Execute(subcmd...); err != nil {
			return err
		}
	}

	return nil
}

func (ipset *CmdIpsetMgr) SetupIpset(name, setType string, entries ...string) error {
	return ipset.AddSetAndEntries(name, setType, entries...)
}

func (ipset *CmdIpsetMgr) AddIPToSet(name string, ip net.IP) error {
	subcmd := strings.Split(fmt.Sprintf("add %s %s", name, ip.String()), " ")
	_, err := ipset.mgr.Execute(subcmd...)
	return err
}

func (ipset *CmdIpsetMgr) DeleteFromSet(name string, ip net.IP) error {
	subcmd := strings.Split(fmt.Sprintf("del %s %s", name, ip.String()), " ")
	_, err := ipset.mgr.Execute(subcmd...)
	return err
}
