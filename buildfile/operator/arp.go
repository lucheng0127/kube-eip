package manager

type ArpManager interface {
	AddTarget(string)
	DeleteTarget(string)
	Poisoning()
}

var ArpMgr ArpManager

func RegisterAndRunArpManager(dev string) error {
	return nil
}
