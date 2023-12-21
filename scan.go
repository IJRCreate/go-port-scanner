package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

type scanResult struct {
	status string
	banner string
}

var targetHost string
var targetPort string

// TODO allow multiple ports to be scanner
func init() {
	flag.StringVar(&targetHost, "H", "", "specify target host")
	flag.StringVar(&targetPort, "p", "", "specify target ports")
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
		fmt.Println(strings.TrimSpace(result.status))
		if result.banner != "" {
			fmt.Println("[+]", strings.TrimSpace(result.banner))
		}
	}
}

func portScan(host string, ports []string) []scanResult {

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

	results := make([]scanResult, len(ports))
	for i, port := range ports {
		results[i] = connectionScan(host, port)
	}
	return results
}

func connectionScan(host string, port string) scanResult {

	targetHost := host
	targetPort := port
	target := fmt.Sprintf("%s:%s", targetHost, targetPort)

	connection, err := net.Dial("tcp", target)
	if err != nil {
		return scanResult{fmt.Sprintf("Scanning port %s\n[-] TCP Closed", targetPort), ""}
	} else {
		defer connection.Close()
		connection.Write([]byte("ViolentPython\r\n"))
		buffer := make([]byte, 100)
		messageLength, err := connection.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		results := strings.TrimSpace(string(buffer[:messageLength]))
		return scanResult{fmt.Sprintf("Scanning port %s\n[+] TCP Open", targetPort), results}
	}

}
