package helper

import "strings"

func CheckPrivateNetwork(ip string) bool {
	return strings.HasPrefix(ip, "127.0.0") || strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "172.16.")
}
