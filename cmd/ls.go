package cmd

import (
	"github.com/spf13/cobra"
	"nvmefc/target"
	"github.com/kataras/golog"
)

func init() {
	targetCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "list",
	Aliases:[]string{"ls", "l"},
	Short: "List target with all subsystems",
	//Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		//if verbose {
		//	fmt.Println("ls")
		//}

		//logger.Debug("ls cmd")
		golog.Debug("ls cmd")
		//fmt.Printf("%p\n", golog.Default)
		//fmt.Printf("short %p\n", golog.Default.Child("short"))
		//fmt.Printf("logger %p\n", logger)
		//golog.Warnf("Route %s regirested", "/mypath")

		target := target.New()
		target.Scan()
		target.PrintTable()
	},
}