package main

import (
	"bufio"
	"fmt"
	"os"

	"libvirt.org/go/libvirt"
)

func main() {
	callback := func(creds []*libvirt.ConnectCredential) {
		for _, cred := range creds {
			if cred.Type == libvirt.CRED_AUTHNAME {
				cred.Result = ""
				cred.ResultLen = len(cred.Result)
			} else if cred.Type == libvirt.CRED_PASSPHRASE {
				cred.Result = ""
				cred.ResultLen = len(cred.Result)
			}
		}
	}
	auth := &libvirt.ConnectAuth{
		CredType: []libvirt.ConnectCredentialType{
			libvirt.CRED_AUTHNAME, libvirt.CRED_PASSPHRASE,
		},
		Callback: callback,
	}
	virConn, err := libvirt.NewConnectWithAuth("qemu:///system", auth, 0)
	if err != nil {
		panic(err)
	}
	input := bufio.NewScanner(os.Stdin)

	statsTypes := libvirt.DOMAIN_STATS_DIRTYRATE
	flags := libvirt.CONNECT_GET_ALL_DOMAINS_STATS_RUNNING
	doms := []*libvirt.Domain{}

	fmt.Println("Ready to find the leak")
	for input.Scan() {
		domss, err := virConn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Found %d domains", len(domss))

		for i := range domss {
			err := domss[i].Free()
			if err != nil {
				fmt.Printf("ERROR freeing domain: %s \n", err)
			}
		}

		fmt.Println("Getting stats")
		domStats, err := virConn.GetAllDomainStats(doms, statsTypes, flags)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Found %d stats for domains", len(domStats))

		for i := range domStats {
			err := domStats[i].Domain.Free()
			if err != nil {
				fmt.Printf("ERROR freeing domain: %s \n", err)
			}
		}
	}
}
