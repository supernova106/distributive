// network.go provides filesystem related thunks.
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"regexp"
	"strconv"
)

// Port parses /proc/net/tcp to determine if a given port is in an open state
// and returns an error if it is not.
func Port(port int) Thunk {
	// strHexToDecimal converts from string containing hex number to int
	strHexToDecimal := func(hex string) int {
		portInt, err := strconv.ParseInt(hex, 16, 64)
		fatal(err)
		return int(portInt)
	}

	// getHexPorts gets all open ports as hex strings from /proc/net/tcp
	getHexPorts := func() (ports []string) {
		toReturn := []string{}
		tcp, err := ioutil.ReadFile("/proc/net/tcp")
		fatal(err)
		// matches only the beginnings of lines
		lines := bytes.Split(tcp, []byte("\n"))
		portRe, err := regexp.Compile(":([0-9A-F]{4})")
		for _, line := range lines {
			port := portRe.Find(line) // only get first port, which is local
			if port == nil {
				continue
			}
			portString := string(port[1:])
			fatal(err)
			toReturn = append(toReturn, portString)
		}
		return toReturn
	}

	// getOpenPorts gets a list of open/listening ports as integers
	getOpenPorts := func() (ports []int) {
		for _, port := range getHexPorts() {
			ports = append(ports, strHexToDecimal(port))
		}
		return ports

	}

	return func() (exitCode int, exitMessage string) {
		for _, p := range getOpenPorts() {
			if p == port {
				return 0, ""
			}
		}
		return 1, "Port " + fmt.Sprint(port) + " did not respond."
	}
}

// Interface detects if a network interface exists
func Interface(name string) Thunk {
	// getInterfaceNames returns the names of all network interfaces
	getInterfaceNames := func() (interfaces []string) {
		ifaces, err := net.Interfaces()
		fatal(err)
		for _, iface := range ifaces {
			interfaces = append(interfaces, iface.Name)
		}
		return
	}
	return func() (exitCode int, exitMessage string) {
		for _, iface := range getInterfaceNames() {
			if iface == name {
				return 0, ""
			}
		}
		return 1, "Interface does not exist: " + name
	}
}

// Up determines if a network interface is up and running or not
func Up(name string) Thunk {
	return func() (exitCode int, exitMessage string) {
		interfaces, err := net.Interfaces()
		fatal(err)
		for _, iface := range interfaces {
			if iface.Name == name && iface.Flags&net.FlagUp != 0 {
				return 0, ""
			}
		}
		return 1, "Interface is down: " + name
	}
}

// Ip4 checks to see if this network interface has this ipv4 address
func Ip4(name string, address string) Thunk {
	return func() (exitCode int, exitMessage string) {
		interfaces, err := net.Interfaces()
		fatal(err)
		for _, iface := range interfaces {
			addresses, err := iface.Addrs
			fatal(err)
			for _, addr := range addresses {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				if ip.To4().String() == address && iface.Name == name {
					return 0, ""
				}
			}
		}
		return 1, "Interface does not have IP: " + name + " " + address
	}
}