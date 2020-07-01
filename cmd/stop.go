package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"nvmefc/target"
)

//var force bool

func init() {
	targetCmd.AddCommand(stopCmd)
	stopCmd.Flags().BoolVarP(&force, "force", "f", false, "Stop Force will delete all configuration and unload driver")
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Aliases:[]string{"st", "s", "Stop", "STOP"},
	Short: "The stop command will stop all",
	Long:  `Unload nvme_rdma nvmet_rdma nvmet nvme nvme_core`,
	//ValidArgs:
	Run: func(cmd *cobra.Command, args []string) {

		if verbose {
			fmt.Printf("Stop Driver force is: %v \n", force)
		}

		if force {
			target.DeleteForce()
		}

		target.UnloadModules()
	},
}