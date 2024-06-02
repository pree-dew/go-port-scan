package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"pscan/scan"
)

func setup(t *testing.T, hosts []string, initList bool) (string, func()) {
	tf, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal(err)
	}

	tf.Close()

	if initList {
		hl := &scan.HostsList{}

		for _, h := range hosts {
			hl.Add(h)
		}

		if err := hl.Save(tf.Name()); err != nil {
			t.Fatal(err)
		}

	}

	return tf.Name(), func() {
		os.Remove(tf.Name())
	}
}

func TestHostActions(t *testing.T) {
	hosts := []string{"host1", "host2", "host3"}
	testCases := []struct {
		name           string
		args           []string
		expectedOut    string
		initList       bool
		actionFunction func(io.Writer, string, []string) error
	}{
		{
			name:           "Add Hosts",
			args:           hosts,
			expectedOut:    "host1\nhost2\nhost3\n",
			initList:       false,
			actionFunction: addAction,
		},
		{
			name:           "List Hosts",
			args:           []string{},
			expectedOut:    "host1\nhost2\nhost3\n",
			initList:       true,
			actionFunction: listAction,
		},
		{
			name:           "Delete Hosts",
			args:           []string{"host1", "host3"},
			expectedOut:    "host1\nhost3\n",
			initList:       true,
			actionFunction: deleteAction,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hostsFile, cleanup := setup(t, hosts, tc.initList)
			defer cleanup()

			out := &bytes.Buffer{}
			if err := tc.actionFunction(out, hostsFile, tc.args); err != nil {
				t.Fatal(err)
			}

			if out.String() != tc.expectedOut {
				t.Fatalf("expected %s, got %s", tc.expectedOut, out.String())
			}
		})
	}
}

func TestIntegration(t *testing.T) {
	hosts := []string{"host1", "host2", "host3"}
	hostsFile, cleanup := setup(t, hosts, false)
	defer cleanup()

	delHost := "host2"
	hostsEnd := []string{"host1", "host3"}

	out := &bytes.Buffer{}
	expectedOut := ""
	for _, v := range hosts {
		expectedOut += fmt.Sprintf("%s\n", v)
	}

	expectedOut += strings.Join(hosts, "\n")
	expectedOut += fmt.Sprintln()
	expectedOut += fmt.Sprintf("%s\n", delHost)
	expectedOut += strings.Join(hostsEnd, "\n")
	expectedOut += fmt.Sprintln()

	if err := addAction(out, hostsFile, hosts); err != nil {
		t.Fatal(err)
	}

	if err := listAction(out, hostsFile, []string{}); err != nil {
		t.Fatal(err)
	}

	if err := deleteAction(out, hostsFile, []string{delHost}); err != nil {
		t.Fatal(err)
	}

	if err := listAction(out, hostsFile, []string{}); err != nil {
		t.Fatal(err)
	}

	if out.String() != expectedOut {
		t.Fatalf("expected %s, got %s", expectedOut, out.String())
	}
}
