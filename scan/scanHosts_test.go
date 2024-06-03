package scan_test

import (
	"net"
	"strconv"
	"testing"

	"pscan/scan"
)

func TestStateString(t *testing.T) {
	ps := scan.PortState{}

	if ps.Open.String() != "closed" {
		t.Errorf("Expected closed, got %s", ps.Open)
	}

	ps.Open = true
	if ps.Open.String() != "open" {
		t.Errorf("Expected open, got %s", ps.Open)
	}
}

func TestRunHostFound(t *testing.T) {
	testCases := []struct {
		name        string
		expectState string
	}{
		{name: "OpenPort", expectState: "open"},
		{name: "ClosedPort", expectState: "closed"},
	}

	host := "localhost"
	hl := &scan.HostsList{}
	hl.Add(host)

	ports := []int{}

	for _, tc := range testCases {
		ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatal(err)
		}

		defer ln.Close()

		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}

		ports = append(ports, port)

		if tc.name == "ClosedPort" {
			ln.Close()
		}
	}

	results := scan.Run(hl, ports)
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Host != host {
		t.Errorf("Expected host %s, got %s", host, results[0].Host)
	}

	if results[0].NotFound {
		t.Errorf("Expected host %s to be found", host)
	}

	if len(results[0].PortStates) != 2 {
		t.Fatalf("Expected 2 port state, got %d", len(results[0].PortStates))
	}

	for i, tc := range testCases {
		if results[0].PortStates[i].Port != ports[i] {
			t.Errorf("Expected port %d, got %d", ports[i], results[0].PortStates[i].Port)
		}

		if tc.name == "OpenPort" && !results[0].PortStates[i].Open {
			t.Errorf("Expected port %d to be open", i)
		}

		if tc.name == "ClosedPort" && results[0].PortStates[i].Open {
			t.Errorf("Expected port %d to be closed", i)
		}
	}
}

func TestRunHostNotFound(t *testing.T) {
	host := "389.389.389.389"
	hl := &scan.HostsList{}
	hl.Add(host)

	results := scan.Run(hl, []int{})

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Host != host {
		t.Errorf("Expected host %s, got %s", host, results[0].Host)
	}

	if !results[0].NotFound {
		t.Errorf("Expected host %s to be not found", host)
	}

	if len(results[0].PortStates) != 0 {
		t.Fatalf("Expected 0 port state, got %d", len(results[0].PortStates))
	}
}
