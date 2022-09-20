package util

import (
	"net"
	"os"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GetLocalIP returns the non loopback local IP of the host
func PrintLocalIPAddresses() string {
	PrintGreen("Local IP Address:")

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				PrintGreen("🔌 > " + ipnet.IP.String())
			}
		}
	}
	return ""
}
