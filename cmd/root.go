package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"fmt"
)

var rootCmd = &cobra.Command{
	Use:   "nvmefc",
	Version:"1.0.0",
	//Aliases:[]string{},
	//DisableFlagParsing:true,

	Short: "nvmefc target / initiator config",
	//Long: `A Fast and Flexible Static Site Generator built with
     //           love by spf13 and friends in Go.
     //           Complete documentation is available at http://hugo.spf13.com`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	 Do Stuff Here
	//},
}


//var log = golog.New()
//var logger = utils.NEW()



func Execute() {
	//fmt.Printf("%p\n", golog.Default)

	//golog.SetTimeFormat("03/01/2006 15:04")

	//log := golog.New()

	// Default Output is `os.Stderr`,
	// but you can change it:
	// log.SetOutput(os.Stdout)

	// Level defaults to "info",
	// but you can change it:
	//log.SetLevel("debug")
	//
	//log.Println("This is a raw message, no levels, no colors.")
	//log.Info("This is an info message, with colors (if the output is terminal)")
	//log.Warn("This is a warning message")
	//log.Error("This is an error message")
	//log.Debug("This is a debug message")


	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}