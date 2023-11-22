package manager

import (
	ipset "github.com/gmccue/go-ipset"
)

const (
	SET_INTERNAL string = "k8s_internal_net"
	SET_EIP      string = "kube-eip-eip"
	SET_VMI      string = "kube-eip-vmi"
)

type IpsetManager interface {
	AddSetAndEntries(string, string, ...string) error
	SetupIpset(string, string, ...string) error
}

type CmdIpsetMgr struct {
	mgr *ipset.IPSet
}

var IpsetMgr IpsetManager

func newIpsetManager() (IpsetManager, error) {
	ipsetMgr := new(CmdIpsetMgr)
	mgr, err := ipset.New()
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
	if err := ipset.mgr.Create(name, setType); err != nil {
		return err
	}

	for _, entry := range entries {
		if err := ipset.mgr.Add(name, entry); err != nil {
			return err
		}
	}

	return nil
}

func (ipset *CmdIpsetMgr) SetupIpset(name, setType string, entries ...string) error {
	return ipset.AddSetAndEntries(name, setType, entries...)
}
