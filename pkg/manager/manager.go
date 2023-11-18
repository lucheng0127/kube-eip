package manager

import "sync"

func RegisterManagers(internal_net ...string) error {
	if err := RegisterIPSetMgr(); err != nil {
		return err
	}

	var setupIpsetErr error
	var once sync.Once

	once.Do(func() {
		setupIpsetErr = SetupIpset("k8s_internal_net", internal_net...)
	})
	if setupIpsetErr != nil {
		return setupIpsetErr
	}

	return nil
}
