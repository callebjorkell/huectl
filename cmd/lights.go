package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newLightsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "lights",
		Aliases: []string{"light", "l"},
		Short:   "Manage Hue light bulbs",
		Args:    cobra.NoArgs,
		// If called with no sub-command, list lights instead of printing help.
		Run: func(*cobra.Command, []string) {
			if err := runListLightsCmd(); err != nil {
				logrus.Fatal(err)
			}
		},
	}
}

func newListLightsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List available lights",
		Args:    cobra.NoArgs,
		Run: func(*cobra.Command, []string) {
			if err := runListLightsCmd(); err != nil {
				logrus.Fatal(err)
			}
		},
	}
}

func runListLightsCmd() error {
	b, err := getBridge()
	if err != nil {
		return fmt.Errorf("unable to find hue brigde: %w", err)
	}

	lights, err := b.GetLights()
	if err != nil {
		return fmt.Errorf("unable to list lights: %w", err)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
	defer tw.Flush()

	fmt.Fprintln(tw, "ID\tNAME\tON\tREACHABLE\tBRIGHTNESS\tHUE\tSAT")

	for _, light := range lights {
		fmt.Fprintf(tw, "%d\t%s\t%t\t%t\t%d\t%d\t%d\n", light.ID, light.Name, light.State.On, light.State.Reachable, light.State.Bri, light.State.Hue, light.State.Sat)
	}

	return nil
}
