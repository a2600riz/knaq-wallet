package network

import (
	"log"
	"net"
)

func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = conn.Close(); err != nil {
			log.Println("error")
		}
	}()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, err
}

func GetOutboundIPString() (string, error) {
	ip, err := GetOutboundIP()
	return ip.String(), err
}
