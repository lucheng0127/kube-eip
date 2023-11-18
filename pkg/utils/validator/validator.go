package validator

import "net"

func ValidateAction(action string) bool {
	if action != "bind" && action != "unbind" {
		return false
	}

	return true
}

func ValidateIPv4(ip string) net.IP {
	return net.ParseIP(ip)
}
