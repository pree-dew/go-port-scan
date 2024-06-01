package scan

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

var (
	ErrExists = fmt.Errorf("host already exists")
	ErrNoHost = fmt.Errorf("host not found")
)

type HostsList struct {
	Hosts []string
}

func (hl *HostsList) Search(host string) (bool, int) {
	// sort the lost of hosts
	sort.Strings(hl.Hosts)

	// binary Search
	i := sort.SearchStrings(hl.Hosts, host)
	if i < len(hl.Hosts) && hl.Hosts[i] == host {
		return true, i
	}

	return false, -1
}

func (hl *HostsList) Add(host string) error {
	exists, _ := hl.Search(host)
	if exists {
		return fmt.Errorf("%w: %s", ErrExists, host)
	}

	hl.Hosts = append(hl.Hosts, host)
	return nil
}

func (hl *HostsList) Remove(host string) error {
	exists, i := hl.Search(host)
	if !exists {
		return fmt.Errorf("%w: %s", ErrNoHost, host)
	}

	hl.Hosts = append(hl.Hosts[:i], hl.Hosts[i+1:]...)
	return nil
}

func (hl *HostsList) Load(hostFile string) error {
	f, err := os.Open(hostFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		hl.Add(scanner.Text())
	}

	return nil
}

func (hl *HostsList) Save(hostFile string) error {
	output := ""
	for _, host := range hl.Hosts {
		output += fmt.Sprintln(host)
	}

	ioutil.WriteFile(hostFile, []byte(output), 0o644)
	return nil
}
