package main

import (
	"fmt"
	"os"
	"path"

	"github.com/mkideal/cli"
)

type restoreT struct {
	cli.Helper
	RestoreIPtables bool `cli:"t,iptables" usage:"Restore iptables" dft:"false"`
	RestoreIPset    bool `cli:"s,ipset" usage:"Restore ipset" dft:"true"`
}

var restoreCMD = &cli.Command{
	Name:    "restore",
	Aliases: []string{"res", "restore"},
	Desc:    "restore ipset and iptables",
	Argv:    func() interface{} { return new(restoreT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*restoreT)
		_, configFile := createAndValidateConfigFile("")
		restoreIPs(configFile, argv.RestoreIPset, argv.RestoreIPtables)
		return nil
	},
}

func restoreIPs(configFile string, restoreIPset, restoreIPtables bool) {
	configFolder, _ := path.Split(configFile)
	iptablesFile := configFolder + "iptables.bak"
	ipsetFile := configFolder + "ipset.bak"

	if restoreIPset {

		_, err := os.Stat(ipsetFile)
		if err != nil {
			_, err = os.Create(ipsetFile)
			fmt.Println("Thereis no ipset backup!")
		} else {
			_, err = runCommand(nil, "ipset restore < "+ipsetFile)
			if err != nil {
				fmt.Println("Error restoring ipset:", err.Error())
			} else {
				fmt.Println("Successfully restored ipset")
			}
		}
	}

	if restoreIPtables {
		_, err := os.Stat(iptablesFile)
		if err != nil {
			fmt.Println("There is no iptables backup!")
		} else {
			_, err = runCommand(nil, "iptables-restore < "+iptablesFile)
			if err != nil {
				fmt.Println("Error restoring iptables:", err.Error())
			} else {
				fmt.Println("Successfully restored iptables")
			}
		}
	}
}
