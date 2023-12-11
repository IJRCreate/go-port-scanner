package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
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
	targetPorts := strings.Split(targetPort, ",")
	results := portScan(targetHost, targetPorts)
	for _, result := range results {
		fmt.Println(result)
	}
}

func portScan(host string, ports []string) []string {

	targetIP, err := net.LookupHost(host)
	if err != nil {
		fmt.Printf("[-] Cannot resolve '%s': Unknown host\n", host)
		log.Fatal(err)
	}

	targetName, err := net.LookupAddr(targetIP[0])
	if err != nil {
		fmt.Printf("[+] Scan Results for: %s\n", targetIP[0])
	} else {
		fmt.Printf("[+] Scan Results for: %s\n", targetName[0])
	}

	results := make([]string, len(ports))
	for i, port := range ports {
		results[i] = connectionScan(host, port)
	}
	return results
}

func connectionScan(host string, port string) string {

	targetHost := host
	targetPort := port
	target := fmt.Sprintf("%s:%s", targetHost, targetPort)

	_, err := net.Dial("tcp", target)
	if err != nil {
		return "TCP Closed"
	} else {
		if err != nil {
			log.Fatal(err)
		}
		return "TCP Open"
	}

}
