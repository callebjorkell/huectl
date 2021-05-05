package cmd

import (
	"errors"
	"github.com/amimof/huego"
	"github.com/callebjorkell/huectl/config"
	"github.com/spf13/cobra"
	"os"
)

func Huectl() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "huectl",
		Short: "huectl controls a Philips Hue installation",
	}

	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newLightsCmd())

	return rootCmd
}

func getBridge() (*huego.Bridge, error) {
	cfg, err := config.Read()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("config file does not exist, please create one")
		}

		return nil, err
	}

	bridge := huego.New(cfg.BridgeAddress, cfg.ClientID)

	return bridge, nil
}