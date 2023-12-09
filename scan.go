package main

import (
	"flag"
	"fmt"
)

var targetHost string
var targetPort string

func init() {
	flag.StringVar(&targetHost, "H", "", "specify target host")
	flag.StringVar(&targetPort, "p", "", "specify target port")
}

func main() {
	flag.Parse()
	if targetHost == "" || targetPort == "" {
		flag.Usage()
	}
	scan(flag.Args())
}

func scan(args []string) string {
	return fmt.Sprintf("%v", args)
	//return fmt.Sprintf("Target Host: %s\nTargetPort\n: %d", targetHost, targetPort)
}
