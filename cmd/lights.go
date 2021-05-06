package cmd

import (
	"fmt"
	"github.com/amimof/huego"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var lightIds []int

func newLightsCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:     "lights",
		Aliases: []string{"light", "l"},
		Short:   "Manage individual lights",
		Args:    cobra.NoArgs,
		// If called with no sub-command, list lights instead of printing help.
		Run: func(*cobra.Command, []string) {
			if err := listLightsCmd(); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.PersistentFlags().IntSliceVar(&lightIds, "id", nil, "List the targeted light IDs. If empty, then all lights will be targeted.")
	cmd.AddCommand(newListLightsCmd())
	cmd.AddCommand(newLightOnCmd())
	cmd.AddCommand(newLightOffCmd())
	cmd.AddCommand(newLightToggleCmd())
	cmd.AddCommand(newLightBrightnessCmd())

	return &cmd
}

func newListLightsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List available lights",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := listLightsCmd(); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func listLightsCmd() error {
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

func newLightOnCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "on",
		Short: "Turn on lights",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := simpleLightCommand(lightIds, lightOn()); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func newLightOffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "off",
		Short: "Turn off lights",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := simpleLightCommand(lightIds, lightOff()); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func newLightToggleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "toggle",
		Short: "Turn off lights",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := simpleLightCommand(lightIds, lightToggle()); err != nil {
				log.Fatal(err)
			}
		},
	}
}

var increaseValue, decreaseValue bool

func newLightBrightnessCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:     "brightness",
		Aliases: []string{"bri"},
		Short:   "Set the brightness",
		Long:    "Set the current brightness on a scale between 1 and 255. The inc and dec flags can be used to set the value relative to the current value.",
		Args:    uintArgs(),
		Run: func(cmd *cobra.Command, args []string) {
			if increaseValue && decreaseValue {
				log.Fatal("cannot both increase and decrease value")
			}

			val := toUint8(args[0])
			var transformer valueTransformer
			if increaseValue {
				transformer = addValue(val)
			} else if decreaseValue {
				transformer = subValue(val)
			} else {
				transformer = setValue(val)
			}

			if err := simpleLightCommand(lightIds, lightBrightness(transformer)); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().BoolVarP(&increaseValue, "inc", "i", false, "Increase the current value.")
	cmd.Flags().BoolVarP(&decreaseValue, "dec", "d", false, "Decrease the current value.")

	return &cmd
}

func uintArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("exactly one argument is required")
		}

		if a, err := strconv.ParseUint(args[0], 10, 8); err != nil {
			return fmt.Errorf("argument needs to be a uint8")
		} else if a < 1 {
			return fmt.Errorf("argument needs to be greater than 0")
		}

		return nil
	}
}

type lightCommand func(*huego.Light) error

func simpleLightCommand(ids []int, cmd lightCommand) error {
	b, err := getBridge()
	if err != nil {
		return fmt.Errorf("unable to find hue brigde: %w", err)
	}

	var lights []huego.Light
	if ids == nil {
		if lights, err = b.GetLights(); err != nil {
			return err
		}
	} else {
		for _, id := range ids {
			l, err := b.GetLight(id)
			if err != nil {
				return fmt.Errorf("could not fetch light %d: %w", id, err)
			}
			lights = append(lights, *l)
		}
	}

	for _, l := range lights {
		if err = cmd(&l); err != nil {
			log.Warningf(err.Error())
		} else {
			log.Debugf("Executed command on light %v", l.ID)
		}
	}

	return nil
}

func lightOn() lightCommand {
	return func(l *huego.Light) error {
		return l.On()
	}
}

type valueTransformer func(uint8) uint8

func toUint8(arg string) uint8 {
	a, err := strconv.ParseUint(arg, 10, 8)
	if err != nil {
		log.Fatal(err)
	}
	return uint8(a)
}

func addValue(a uint8) valueTransformer {
	return func(b uint8) uint8 {
		if a > 255 - b {
			return 255
		}
		return a + b
	}
}

func subValue(a uint8) valueTransformer {
	return func(b uint8) uint8 {
		if b <= a {
			return 1
		}
		return b - a
	}
}

func setValue(a uint8) valueTransformer {
	return func(uint8) uint8 {
		return a
	}
}

func lightBrightness(transformer valueTransformer) lightCommand {
	return func(l *huego.Light) error {
		value := l.State.Bri
		return l.Bri(transformer(value))
	}
}

func lightOff() lightCommand {
	return func(l *huego.Light) error {
		return l.Off()
	}
}

func lightToggle() lightCommand {
	return func(l *huego.Light) error {
		if l.IsOn() {
			return l.Off()
		}
		return l.On()
	}
}

