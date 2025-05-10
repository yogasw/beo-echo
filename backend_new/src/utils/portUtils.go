package utils

import (
	"fmt"
	"net"
	"time"
)

// IsPortSafe checks if a port number is in the safe range
func IsPortSafe(port int) bool {
	return port >= 9001 && port <= 9999
}

// IsPortInUse checks if a port is already in use
func IsPortInUse(port int) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), timeout)
	if err != nil {
		return false
	}

	if conn != nil {
		conn.Close()
		return true
	}

	return false
}
