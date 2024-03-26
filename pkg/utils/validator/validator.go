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

func ValidatePort(port int) int {
	if port > 0 && port < 65535 {
		return port
	}

	return 0
}
