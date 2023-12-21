package main

import (
	"log"
	"os"
	"os/exec"
	"slices"
	"testing"
	"time"
)

/* type expected struct {
	status string
	banner string
} */

type testData struct {
	host     string
	ports    []string
	expected []scanResult
}

var passingTests = []testData{
	{
		host:  "127.0.0.1",
		ports: []string{"8000"},
		expected: []scanResult{
			{
				status: "Scanning port 8000\n[+] TCP Open",
				banner: "Banner 0",
			},
		},
	},
	{
		host:  "127.0.0.1",
		ports: []string{"8001", "8002", "8003"},
		expected: []scanResult{
			{
				status: "Scanning port 8001\n[+] TCP Open",
				banner: "Banner 1",
			},
			{
				status: "Scanning port 8002\n[+] TCP Open",
				banner: "Banner 2",
			},
			{
				status: "Scanning port 8003\n[+] TCP Open",
				banner: "Banner 3",
			},
		},
	},
	{
		host:  "127.0.0.1",
		ports: []string{"8004", "8005", "8006"},
		expected: []scanResult{
			{
				status: "Scanning port 8004\n[+] TCP Open",
				banner: "Banner 4",
			},
			{
				status: "Scanning port 8005\n[+] TCP Open",
				banner: "Banner 5",
			},
			{
				status: "Scanning port 8006\n[+] TCP Open",
				banner: "Banner 6",
			},
		},
	},
}

var failingTests = []testData{
	{
		host:  "127.0.0.1",
		ports: []string{"8007"},
		expected: []scanResult{
			{
				status: "Scanning port 8007\n[-] TCP Closed",
				banner: "",
			},
		},
	},
	{
		host:  "127.0.0.1",
		ports: []string{"8008", "8009", "8010"},
		expected: []scanResult{
			{
				status: "Scanning port 8008\n[-] TCP Closed",
				banner: "",
			},
			{
				status: "Scanning port 8009\n[-] TCP Closed",
				banner: "",
			},
			{
				status: "Scanning port 8010\n[-] TCP Closed",
				banner: "",
			},
		},
	},
	{
		host:  "127.0.0.1",
		ports: []string{"8011", "8012", "8013"},
		expected: []scanResult{
			{
				status: "Scanning port 8011\n[-] TCP Closed",
				banner: "",
			},
			{
				status: "Scanning port 8012\n[-] TCP Closed",
				banner: "",
			},
			{
				status: "Scanning port 8013\n[-] TCP Closed",
				banner: "",
			},
		},
	},
}

func openPort(port string, test testData, index int) *exec.Cmd {
	// TODO Use commented commands to emulate banners for later grabbing
	echo := exec.Command("echo", test.expected[index].banner)
	netcat := exec.Command("nc", "-lvp", port)
	netcat.Stdin, _ = echo.StdoutPipe()
	err := netcat.Start()
	if err != nil {
		log.Fatal(err)
	}
	// Required to allow netcat to start in time
	time.Sleep(100 * time.Millisecond)
	err = echo.Run()
	if err != nil {
		log.Fatal(err)
	}
	return netcat
}

func closePort(netcat *exec.Cmd) {
	netcat.Process.Signal(os.Interrupt)
}

func TestConnectionScan(t *testing.T) {

	test := passingTests[0]
	listener := openPort(test.ports[0], test, 0)
	actual := connectionScan(test.host, test.ports[0])
	closePort(listener)
	if actual.status != test.expected[0].status || actual.banner != test.expected[0].banner {
		t.Errorf("Expected: %#v; Actual: %#v", test.expected[0], actual)
	}

	test = failingTests[0]
	actual = connectionScan(test.host, test.ports[0])
	if actual.status != test.expected[0].status || actual.banner != test.expected[0].banner {
		t.Errorf("Expected: %#v; Actual: %#v", test.expected[0], actual)
	}
}

func TestOneHostOnePortPortScan(t *testing.T) {

	test := passingTests[0]
	var actual []scanResult
	port := test.ports[0]
	listener := openPort(port, test, 0)
	actual = portScan(test.host, []string{port})
	closePort(listener)
	if !slices.Equal(actual, test.expected) {
		t.Errorf("Expected: %s; Actual: %s", test.expected, actual)
	}

	test = failingTests[0]
	port = test.ports[0]
	actual = portScan(test.host, []string{port})
	if !slices.Equal(actual, test.expected) {
		t.Errorf("Expected: %s; Actual: %s", test.expected, actual)
	}
}

func TestOneHostMultiplePortScans(t *testing.T) {

	test := passingTests[1]
	var actual []scanResult
	listeners := make([]*exec.Cmd, len(test.ports))
	for i, port := range test.ports {
		listeners[i] = openPort(port, test, i)
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
		var actual []scanResult
		listeners := make([]*exec.Cmd, len(test.ports))
		for i, port := range test.ports {
			listeners[i] = openPort(port, test, i)
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
