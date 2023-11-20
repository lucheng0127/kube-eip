package manager

import "github.com/lucheng0127/kube-eip/pkg/utils/errcheck"

func RegisterManagers(internal_net ...string) error {
	if err := RegisterIPSetMgr(); err != nil {
		return err
	}

	err := SetupIpset("k8s_internal_net", internal_net...)
	if err != nil && !errcheck.IsExistError(err) {
		return err
	}

	return nil
}
