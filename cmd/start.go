package cmd

import (
	"fmt"
	"nvmefc/target"
	"runtime"

	"github.com/spf13/cobra"
)

var offloadMemSize = 8192
var numP2pQueues = 2
var offloadBufferSize int
var platform = runtime.GOARCH
var offloadMemStart = fmt.Sprintf("%#x", target.Kernel()) // no error handle if target.Kernel() return 0 or 01 then there is problem
var numberDisks = len(target.GetDisks())

// pool_size int // add ?
// buf_size int // add ?

func init() {
	targetCmd.AddCommand(startCmd)
	startCmd.AddCommand(staticCmd, dynamicCmd, nonOffloadCmd)

	if platform == "ppc64le" {
		offloadMemStart = fmt.Sprintf("%#x", 576460757672132608)
		offloadMemSize = 512
	}

	if numberDisks > 0 {
		//	offloadBufferSize = offloadMemSize / numberDisks
		offloadBufferSize = offloadMemSize / numberDisks / numP2pQueues
	}

	startCmd.PersistentFlags().IntVarP(&offloadMemSize, "mem", "M", offloadMemSize, "offloadMemSize")
	startCmd.PersistentFlags().IntVarP(&offloadBufferSize, "offloadBufferSize", "B", offloadBufferSize, "offloadBufferSize is (offloadMemSize / numberDisks / numP2pQueues)")
	startCmd.PersistentFlags().IntVarP(&numP2pQueues, "p2p", "P", numP2pQueues, "numP2pQueues")
	startCmd.PersistentFlags().IntVarP(&numberDisks, "numberDisks", "D", numberDisks, "numberDisks can help to calculate offloadBufferSize pool_size (offloadMemSize / numberDisks)")

	staticCmd.Flags().StringVarP(&offloadMemStart, "offloadMemStart", "O", offloadMemStart, "offloadMemStart")

	//startCmd.PersistentFlags().StringVarP(&platform, "platform", "C", platform, "platform")

	//fmt.Printf("%v\n", false)
}

var startCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"st", "s", "S", "Start", "START"},
	Short:   "Start driver",
	Long:    `Load nvme_rdma nvmet_rdma nvmet nvme nvme_core`,
}

var staticCmd = &cobra.Command{
	Use:     "static",
	Aliases: []string{"st", "s"},
	Short:   "The will start driver in offload static mode",
	Run: func(cmd *cobra.Command, args []string) {

		if verbose {
			fmt.Printf("Start offload static platform: %s\n", platform)
			fmt.Println(offloadMemSize, numP2pQueues, offloadBufferSize, platform, offloadMemStart, numberDisks)
			fmt.Printf("%p\n", &offloadMemStart)
		}

		//if numberDisks > 0 {
		//	offloadBufferSize = offloadMemSize / numberDisks
		//}
		//
		if platform == "ppc64le" {
			target.LoadModulesppc64le(offloadMemSize, offloadBufferSize, numP2pQueues, offloadMemStart)
			return
		}

		//offloadBufferSize = offloadMemSize / 6 // need to be changed

		// pass target.Kernel() uint64 instead of offloadMemStart string ?
		target.LoadModules(offloadMemSize, offloadBufferSize, numP2pQueues, target.Kernel())
	},
}

var dynamicCmd = &cobra.Command{
	Use:     "dynamic",
	Aliases: []string{"dy", "d", "D"},
	Short:   "The will start driver in offload dynamic mode",
	//Long:  `All software has versions. This is Hugo's`,
	//ValidArgs:
	Run: func(cmd *cobra.Command, args []string) {

		if offloadBufferSize < 1 {
			offloadBufferSize = 256
		}

		if verbose {
			fmt.Println("Start offload dynamic")
			fmt.Println(offloadMemSize, numP2pQueues, offloadBufferSize, platform, offloadMemStart, numberDisks)
		}

		target.LoadModulesDynamic(numP2pQueues, offloadBufferSize)

		//target.UnloadModules()
	},
}

var nonOffloadCmd = &cobra.Command{
	Use:     "nonOffload",
	Aliases: []string{"nonoffload", "non_offload", "non-offload", "nono", "r", "n"},
	Short:   "The will start driver in non-offload mode",
	//Long:  `All software has versions. This is Hugo's`,
	//ValidArgs:
	Run: func(cmd *cobra.Command, args []string) {

		if verbose {
			fmt.Println("Start non-offload")
		}

		target.LoadModulesDynamic(numP2pQueues, offloadBufferSize)
		//target.UnloadModules()
	},
}
