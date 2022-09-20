package util

import (
	"bufio"
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

// Read file utility function
func ReadFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bin := make([]byte, size)
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bin)

	return string(bin), err
}
