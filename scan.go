package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

var targetHost string
var targetPort string

// TODO allow multiple ports to be scanner
func init() {
	flag.StringVar(&targetHost, "H", "", "specify target host")
	flag.StringVar(&targetPort, "p", "", "specify target port")
}

// TODO call portScan and not connectionScan
func main() {
	flag.Parse()
	if targetHost == "" || targetPort == "" {
		flag.Usage()
	}
	fmt.Println(connectionScan(targetHost, targetPort))
}

// TODO move hostname & address checks from connectionScan to portScan
func portScan(host string, ports []string) []string {
	results := make([]string, len(ports))
	for i, port := range ports {
		results[i] = connectionScan(host, port)
	}
	return results
}

func connectionScan(host string, port string) string {

	var targetHost string
	ipAddress, err := net.LookupHost(host)
	if err != nil {
		// TODO iterate through all returned ip Addresses
		targetHost = ipAddress[0]
	} else {
		fmt.Printf("Cannot resolve '%s': Unknown host\n", host)
		targetHost = host
	}
	targetPort := port
	target := fmt.Sprintf("%s:%s", targetHost, targetPort)

	_, err = net.Dial("tcp", target)
	if err != nil {
		return "TCP Closed"
	} else {
		if err != nil {
			log.Fatal(err)
		}
		return "TCP Open"
	}

}
