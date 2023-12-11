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
	{host: "127.0.0.1", ports: []string{"8000", "8001", "8002"}, expected: []string{"TCP Open", "TCP Open", "TCP Open"}},
	{host: "127.0.0.1", ports: []string{"8003", "8004", "8005"}, expected: []string{"TCP Open", "TCP Open", "TCP Open"}},
}

var failingTests = []testData{
	{host: "127.0.0.1", ports: []string{"8006", "8007", "8008"}, expected: []string{"TCP Closed", "TCP Closed", "TCP Closed"}},
	{host: "127.0.0.1", ports: []string{"8009", "8010", "8011"}, expected: []string{"TCP Closed", "TCP Closed", "TCP Closed"}},
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
	time.Sleep(100 * time.Millisecond)
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

	test := passingTests[0]
	listener := openPort(test.ports[0])
	actual := strings.TrimSpace(connectionScan(test.host, test.ports[0]))
	if actual != test.expected[0] {
		t.Errorf("Expected: %s; Actual: %s", test.expected[0], actual)
	}
	closePort(listener)

	test = failingTests[0]
	actual = strings.TrimSpace(connectionScan(test.host, test.ports[0]))
	if actual != test.expected[0] {
		t.Errorf("Expected: %s; Actual: %s", test.expected[0], actual)
	}
}

func TestOneHostOnePortPortScan(t *testing.T) {

	test := passingTests[0]
	var actual []string
	port := test.ports[0]
	listener := openPort(port)
	actual = portScan(test.host, []string{port})
	closePort(listener)
	if !slices.Equal(actual, []string{test.expected[0]}) {
		t.Errorf("Expected: %s; Actual: %s", []string{test.expected[0]}, actual)
	}

	test = failingTests[0]
	port = test.ports[0]
	actual = portScan(test.host, []string{port})
	if !slices.Equal(actual, []string{test.expected[0]}) {
		t.Errorf("Expected: %s; Actual: %s", []string{test.expected[0]}, actual)
	}
}

func TestOneHostMultiplePortScans(t *testing.T) {

	test := passingTests[0]
	var actual []string
	listeners := make([]*exec.Cmd, len(test.ports))
	for i, port := range test.ports {
		listeners[i] = openPort(port)
	}
	actual = portScan(test.host, test.ports)
	for _, listener := range listeners {
		closePort(listener)
	}
	if !slices.Equal(actual, test.expected) {
		t.Errorf("Expected: %s; Actual: %s", test.expected, actual)
	}

	test = failingTests[0]
	actual = portScan(test.host, test.ports)
	if !slices.Equal(actual, test.expected) {
		t.Errorf("Expected: %s; Actual: %s", test.expected, actual)
	}
}

func TestMultipleHostsMultiplePortScans(t *testing.T) {

	for _, test := range passingTests {
		var actual []string
		listeners := make([]*exec.Cmd, len(test.ports))
		for i, port := range test.ports {
			listeners[i] = openPort(port)
		}
		actual = portScan(test.host, test.ports)
		for _, listener := range listeners {
			closePort(listener)
		}
		if !slices.Equal(actual, test.expected) {
			t.Errorf("Expected: %s; Actual: %s", test.expected, actual)
		}
	}

	for _, test := range failingTests {
		actual := portScan(test.host, test.ports)
		if !slices.Equal(actual, test.expected) {
			t.Errorf("Expected: %s; Actual: %s", test.expected, actual)
		}
	}
}
