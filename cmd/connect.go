package cmd

import (
	"fmt"
	"github.com/amimof/huego"
	"github.com/callebjorkell/huectl/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newConnectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connect",
		Short: "Connects to a bridge, storing the new user in a config. The link button must be pressed before running this command.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := connectBridge(); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func connectBridge() error {
	bridge, err := huego.Discover()
	if err != nil {
		return fmt.Errorf("could not find a bridge: %w", err)
	}
	user, err := bridge.CreateUser("huectl")
	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}

	c := config.Config{
		BridgeAddress: bridge.Host,
		ClientID:      user,
	}

	err = c.Write()
	if err != nil {
		return err
	}

	fmt.Printf("Connected to bridge at %v\n", bridge.Host)
	return nil
}
