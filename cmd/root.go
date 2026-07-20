package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/septrum101/zteOnu/app/factory"
	"github.com/septrum101/zteOnu/app/telnet"
	"github.com/septrum101/zteOnu/version"
)

var (
	// Used for flags.
	user       string
	passwd     string
	ip         string
	port       int
	permTelnet bool
	telnetPort int
	newMode    bool
	iface      string

	rootCmd = &cobra.Command{
		Use: "zteOnu",
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(); err != nil {
				fmt.Println(err)
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "telecomadmin", "factory mode auth username")
	rootCmd.PersistentFlags().StringVarP(&passwd, "pass", "p", "nE7jA%5m", "factory mode auth password")
	rootCmd.PersistentFlags().StringVarP(&ip, "ip", "i", "192.168.1.1", "ONU ip address")
	rootCmd.PersistentFlags().IntVar(&port, "port", 8080, "ONU http port")
	rootCmd.PersistentFlags().BoolVar(&permTelnet, "telnet", false, "permanent telnet (user: root, pass: Zte521)")
	rootCmd.PersistentFlags().IntVar(&telnetPort, "tp", 23, "ONU telnet port")
	rootCmd.PersistentFlags().BoolVar(&newMode, "new", false, "use new method to open telnet; the SendInfo payload is derived from the current interface MAC")
	rootCmd.PersistentFlags().StringVar(&iface, "iface", "", "network interface to read the MAC from (default: first non-loopback interface)")
}

func run() error {
	version.Show()

	if newMode {
		// The SendInfo payload is now derived from the current interface MAC
		// (see factory.MacToMagicBytes), so any MAC the device accepts works.
		// We only need to make sure a usable MAC is actually present.
		if _, err := factory.LocalMAC(iface); err != nil {
			return fmt.Errorf("new mode requires a usable network interface MAC: %w", err)
		}
	}

	tlUser, tlPass, err := factory.New(user, passwd, ip, port, iface).Handle()
	if err != nil {
		return err
	}

	if permTelnet {
		// create telnet conn
		t, err := telnet.New(tlUser, tlPass, ip, telnetPort)
		if err != nil {
			return err
		}
		defer t.Conn.Close()

		// handle permanent telnet
		if err := t.PermTelnet(); err != nil {
			return err
		}
		fmt.Println("Permanent Telnet succeed\r\nuser: root, pass: Zte521")

		// reboot device
		fmt.Println("wait reboot..")
		time.Sleep(time.Second)
		if err := t.Reboot(); err != nil {
			return err
		}
	} else {
		fmt.Printf("user: %s\npass: %s", tlUser, tlPass)
	}
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
