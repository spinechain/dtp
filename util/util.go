package util

import (
	"bufio"
	"net"
	"os"
	"strings"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func IsIPLocalIP(ip string) bool {

	// Get addresses
	addresses := GetLocalIPAddresses()

	// loop and print addresses
	for _, address := range addresses {

		// split ip to port and ip
		ipParts := strings.Split(ip, ":")

		if len(ipParts) == 2 {
			ip = ipParts[0]
		} else {
			ip = ipParts[0]
		}

		if ip == address {
			return true
		}
	}

	return false
}

// Get all Local IP addresses as a list
func GetLocalIPAddresses() []string {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	var addresses []string
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				addresses = append(addresses, ipnet.IP.String())
			}
		}
	}

	return addresses
}

func PrintLocalIPAddresses() {
	PrintGreen("Local IP Address:")

	// Get addresses
	addresses := GetLocalIPAddresses()

	// loop and print addresses
	for _, address := range addresses {
		PrintGreen("🔌 > " + address)
	}

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

func FirstWords(value string, count int) string {
	// Loop over all indexes in the string.
	for i := range value {
		// If we encounter a space, reduce the count.
		if value[i] == ' ' {
			count -= 1
			// When no more words required, return a substring.
			if count == 0 {
				return value[0:i]
			}
		}
	}
	// Return the entire string.
	return value
}
