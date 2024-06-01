package scan_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"pscan/scan"
)

func TestAdd(t *testing.T) {
	testcases := []struct {
		name        string
		host        string
		expectedLen int
		expectedErr error
	}{
		{name: "AddNew", host: "host2", expectedLen: 2, expectedErr: nil},
		{name: "AddExisting", host: "host1", expectedLen: 1, expectedErr: fmt.Errorf("%w: host1", scan.ErrExists)},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			hl := scan.HostsList{Hosts: []string{"host1"}}
			err := hl.Add(tc.host)
			if err != nil && errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}

			if len(hl.Hosts) != tc.expectedLen {
				t.Errorf("expected %v, got %v", tc.expectedLen, len(hl.Hosts))
			}
		})
	}
}

func TestRemove(t *testing.T) {
	testcases := []struct {
		name        string
		host        string
		expectedLen int
		expectedErr error
	}{
		{name: "RemoveExisting", host: "host1", expectedLen: 0, expectedErr: nil},
		{name: "RemoveNonExisting", host: "host2", expectedLen: 1, expectedErr: fmt.Errorf("%w: host2", scan.ErrNoHost)},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			hl := scan.HostsList{Hosts: []string{"host1"}}
			err := hl.Remove(tc.host)
			if err != nil && errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}

			if len(hl.Hosts) != tc.expectedLen {
				t.Errorf("expected %v, got %v", tc.expectedLen, len(hl.Hosts))
			}
		})
	}
}

func TestSaveLoad(t *testing.T) {
	hl1 := scan.HostsList{}
	hl2 := scan.HostsList{}

	hl1.Add("host1")
	// create temp file
	tmpfile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	if err := hl1.Save(tmpfile.Name()); err != nil {
		t.Fatal(err)
	}

	if err := hl2.Load(tmpfile.Name()); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(hl1.Hosts, hl2.Hosts) {
		t.Errorf("expected %v, got %v", hl1.Hosts, hl2.Hosts)
	}
}

func TestLoadNoFile(t *testing.T) {
	tf, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Remove(tf.Name()); err != nil {
		t.Fatal(err)
	}

	hl := scan.HostsList{}
	if err := hl.Load(tf.Name()); err != nil {
		t.Fatal(err)
	}

	if len(hl.Hosts) != 0 {
		t.Errorf("expected 0, got %v", len(hl.Hosts))
	}
}
