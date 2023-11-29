package manager

import "github.com/lucheng0127/kube-eip/pkg/arp"

type ArpManager interface {
	AddTarget(string)
	DeleteTarget(string)
	Poisoning()
}

var ArpMgr ArpManager

func RegisterAndRunArpManager(dev string) error {
	mgr, err := arp.NewArpCracker(dev)
	if err != nil {
		return err
	}

	ArpMgr = mgr

	// TODO(shawnlu): get current eip from ipset

	go ArpMgr.Poisoning()
	return nil
}
