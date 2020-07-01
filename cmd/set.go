package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	//"nvmefc/target"
	"nvmefc/target"
)

func init() {
	targetCmd.AddCommand(setCmd)
	setCmd.AddCommand(setDiskCmd)
	setDiskCmd.AddCommand(setDiskEnableCmd, setDiskDisableCmd)

	setCmd.AddCommand(setOffloadCmd)
	setOffloadCmd.AddCommand(setDisableOffloadCmd, setEnableOffloadCmd)
}

var enableAliases = []string{"e", "E", "on", "Enable", "Enable", "true"}
var disableAliases = []string{"d", "D", "off", "Disable", "DISABLE", "false"}

var setCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"s", "Set", "modify"},
	Short:   "Set / Modify subsystems",
	Long:    `This command set subsystems context`,
	Example: "nvmefc target set offload disable",
	Args:    cobra.NoArgs,
	Version: "0.1.0",
}

var setDiskCmd = &cobra.Command{
	Use:     "disk",
	Short:   "Modify all xsubsystems disk state",
	Long:    `This command set subsystems disk context.`,
	Example: "nvmefc target set disk",

	Args: cobra.NoArgs,
}

var setDiskEnableCmd = &cobra.Command{
	Use:     "enable",
	Aliases: enableAliases,
	Short:   "Enable disk for all subsystems",

	Long:    `This command set subsystems disk context`,
	Example: "nvmefc target set disk enable",

	Args: cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {

		if verbose {
			fmt.Println("setEnableCmd", args)
		}
		target := target.New()
		target.Scan()
		target.SetSubsystemsDeviceSysFsEnable(true)
	},
}

var setDiskDisableCmd = &cobra.Command{
	Use:     "disable",
	Aliases: disableAliases,

	Short: "Disable disk for all subsystems",

	Long:    `This command set subsystems disk context`,
	Example: "nvmefc target set disk disable",

	Args: cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {

		if verbose {
			fmt.Println("setDisableCmd", args)
		}
		target := target.New()
		target.Scan()
		target.SetSubsystemsDeviceSysFsEnable(false)
	},
}

var setOffloadCmd = &cobra.Command{
	Use:     "offload",
	Aliases: []string{"o", "O", "Offload"},

	Short:   "Modify all subsystems offload state",
	Long:    `This command set subsystems context offload. This works only if the controller is not in use (initiators are not connected to)`,
	Example: "nvmefc target set offload",

	Args: cobra.NoArgs,
}

var setDisableOffloadCmd = &cobra.Command{
	Use:     "disable",
	Aliases: disableAliases,
	Short:   "Disable offload for all subsystems",
	Long:    `This command set subsystems offload context`,
	Example: "nvmefc target set offload disable",

	Args: cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {

		if verbose {
			fmt.Println("setDisableOffloadCmd", args)
		}
		target := target.New()
		target.Scan()

		target.SetSubsystemsDeviceSysFsEnable(false)
		target.UnlinkSubsystemsPorts()

		target.SetSubsystemOffload(false)

		target.SetSubsystemsDeviceSysFsEnable(true)
		target.LinkSubsystemsPorts()
	},
}

var setEnableOffloadCmd = &cobra.Command{
	Use:     "enable",
	Aliases: enableAliases,

	Short:   "Enable offload for all subsystems",
	Long:    `This command set subsystems offload context`,
	Example: "nvmefc target set enable disable",

	Args: cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {

		if verbose {
			fmt.Println("setEnableOffloadCmd", args)
		}
		target := target.New()
		target.Scan()

		target.SetSubsystemsDeviceSysFsEnable(false)
		target.UnlinkSubsystemsPorts()

		target.SetSubsystemOffload(true)

		target.SetSubsystemsDeviceSysFsEnable(true)
		target.LinkSubsystemsPorts()

	},
}
