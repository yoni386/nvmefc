package cmd

import (

	"github.com/spf13/cobra"
	"github.com/kataras/golog"
)

var verbose bool
var force bool

func init() {
	rootCmd.AddCommand(targetCmd)
	//targetCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose info")
	targetCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose info (default false)")
	//targetCmd.AddCommand(lsCmd, createCmd, stopCmd, rmCmd)
}

var targetCmd = &cobra.Command{
	Use:   "target",
	Aliases:[]string{"t", "T", "tar"},
	Short: "target context",
	//Long:  `All software has versions. This is Hugo's`,
	//TraverseChildren:true,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	//},

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		//fmt.Printf("Inside rootCmd PersistentPreRun with args: %v\n", args)

		if verbose {
			//logger.Log1.Default.Level = golog.DebugLevel
			//utils.Log1.Default.Level = golog.DebugLevel
			//utils.Log1.SetLevel()
			//fmt.Printf("verbose: %v\n", verbose)
			//logger.Level = utils.DebugLevel
			//golog.Default.Level = utils.DebugLevel
			//logger.Debug("Set golog.Default.Level to golog.DebugLevel")
			//logger.Debug("Set golog.Default.Level to golog.DebugLevel")

			golog.Default.Level = golog.DebugLevel
			golog.Default.Child("simple").Level = golog.DebugLevel
			golog.Debug("Set golog.Default.Level to golog.DebugLevel")
			golog.Debug("Set golog.Default.Level to golog.DebugLevel")
			golog.Child("simple").Debug("Set golog.Default.Level to golog.DebugLevel")
		}
	},
}