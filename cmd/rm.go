package cmd

import (
	"github.com/spf13/cobra"
	"nvmefc/target"
	"github.com/kataras/golog"
)

//var force bool

func init() {
	targetCmd.AddCommand(rmCmd)
	rmCmd.Flags().BoolVarP(&force, "force", "f", false, "Remove Force")
}

var rmCmd = &cobra.Command{
	Use:   "remove",
	Aliases:[]string{"rm", "r"},
	Short: "The rm command will remove all configuration",
	//Long:  `All software has versions. This is Hugo's`,
	//ValidArgs:
	Run: func(cmd *cobra.Command, args []string) {

		//if verbose {
		//	fmt.Printf("Remove configuration force is: %v \n", force)
		//
		//}

		golog.Debug("rm cmd")

		if force {
			golog.Warnf("Remove configuration force is: %v \n", force)
			target.DeleteForce()
			return
		}

		t := target.New()
		t.Scan()
		t.Destroy()
		//t.DestroyDisks()
		//t.DestroyPorts()
	},
}