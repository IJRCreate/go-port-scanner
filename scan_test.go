package main

import (
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
	"time"
)

type testData struct {
	host     string
	ports    []string
	expected []string
}

var passingTests = []testData{
	{host: "127.0.0.1", ports: []string{"8000"}, expected: []string{"TCP Open"}},
}

var failingTests = []testData{
	{host: "127.0.0.1", ports: []string{"8001"}, expected: []string{"TCP Closed"}},
}

func openPort(port string) *exec.Cmd {
	// TODO Use commented commands to emulate banners for later grabbing
	// echo := exec.Command("echo", test.expected)
	netcat := exec.Command("nc", "-lvp", port)
	//netcat.Stdin, _ = echo.StdoutPipe()
	err := netcat.Start()
	if err != nil {
		log.Fatal(err)
	}
	// Required to allow netcat to start in time
	time.Sleep(5 * time.Millisecond)
	/* err = echo.Run()
	if err != nil {
		log.Fatal(err)
	} */
	return netcat
}

func closePort(netcat *exec.Cmd) {
	netcat.Process.Signal(os.Interrupt)
}

func TestConnectionScan(t *testing.T) {

	for _, test := range passingTests {

		listener := openPort(test.ports[0])

		actual := strings.TrimSpace(connectionScan(test.host, test.ports[0]))
		if actual != test.expected[0] {
			t.Errorf("Expected: %s; Actual: %s", test.expected, actual)
		}

		closePort(listener)
	}

	for _, test := range failingTests {
		actual := strings.TrimSpace(connectionScan(test.host, test.ports[0]))
		if actual != test.expected[0] {
			t.Errorf("Expected: %s; Actual: %s", test.expected, actual)
		}
	}

}

func TestPortScan(t *testing.T) {
	for _, test := range passingTests {
		var actual []string
		listeners := make([]*exec.Cmd, len(test.ports))
		for i, port := range test.ports {
			listeners[i] = openPort(port)
		}

		actual = portScan(test.host, test.ports)

		if !slices.Equal(actual, test.expected) {
			t.Errorf("Expected: %s; Actual: %s", test.expected, actual)
		}

		for _, listener := range listeners {
			closePort(listener)
		}

	}

	for _, test := range failingTests {
		var actual []string
		actual = portScan(test.host, test.ports)
		if !slices.Equal(actual, test.expected) {
			t.Errorf("Expected: %s; Actual: %s", test.expected, actual)
		}
	}
}
