package manager

import (
	ipset "github.com/gmccue/go-ipset"
)

type IpsetManager interface {
	AddSetAndEntries(string, ...string) error
}

type CmdIpsetMgr struct {
	mgr *ipset.IPSet
}

var IpsetMgr CmdIpsetMgr

func RegisterIPSetMgr() error {
	mgr, err := ipset.New()
	if err != nil {
		return err
	}

	IpsetMgr.mgr = mgr
	return nil
}

func (ipset *CmdIpsetMgr) AddSetAndEntries(name string, entries ...string) error {
	if err := ipset.mgr.Create(name, "hash:net"); err != nil {
		return err
	}

	for _, entry := range entries {
		if err := ipset.mgr.Add(name, entry); err != nil {
			return err
		}
	}

	return nil
}

func SetupIpset(name string, entries ...string) error {
	return IpsetMgr.AddSetAndEntries(name, entries...)
}
