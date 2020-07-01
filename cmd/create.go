package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"nvmefc/target"
	"strconv"
	"log"
	"github.com/kataras/golog"
)

var offload bool
var dryRun bool
var ips = target.GetIps()
var ports []int
var prefix = fmt.Sprintf("sub_%s", target.GetHostname())

func init() {
	targetCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVarP(&offload, "offload", "o", false, "Create with Offload (default false)")
	createCmd.Flags().BoolVarP(&dryRun, "dryRun", "d", false, "dryRun (default false)")
	createCmd.Flags().StringArrayVarP(&ips, "ips", "I", ips, "Ip-Address Array")
	createCmd.Flags().IntSliceVarP(&ports, "ports", "P", []int{3001, 4001}, "Ports Array")
	createCmd.Flags().StringVarP(&prefix, "prefix", "p", prefix, "Prefix for subsystem name. Default is sub_hostname_diskName")
	//createCmd.MarkFlagRequired("ips")
}


//func dd(disks []target.Disk) {
//	h, err := os.Hostname()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//
//	for i := range disks {
//		name := fmt.Sprintf("sub_%s_%s", h, disks[i].Name())
//		if err := target.CreateSubsystem(name, offload); err != nil {
//			log.Println(err)
//			return
//		}
//		//idx++
//		id := strconv.Itoa(i)
//
//		target.CreateNameSpace(name, id, disks[i].Path(), target.NewUUID(), target.ByteOne)
//	}
//
//}

var createCmd = &cobra.Command{
	Use:   "create",
	//Use: "print [OPTIONS] [COMMANDS]",
	Aliases:[]string{"cr", "c"},
	Short: "Create target with all subsystems",
	Args: cobra.ExactArgs(0),

	//Long:  `All software has versions. This is Hugo's`,
	//ValidArgs:
	Run: func(cmd *cobra.Command, args []string) {
		var id int

		if verbose {
			//fmt.Println("Create createCmd")
			golog.Debug("Create createCmd")
		}

		if len(ips) != len(ports) {
			//fmt.Printf("Number ip-address is: %d and number of ports is: %d \n", len(ips), len(ports))
			golog.Warnf("Ip-address are: %v and number of ports are: %v \n", ips, ports)
			golog.Errorf("Number ip-address is: %d and number of ports is: %d \n", len(ips), len(ports))
			if verbose {
				//fmt.Printf("Ip-address are: %v and number of ports are: %v \n", ips, ports)
			}
			return
		}

		if len(ips) < 1 {
			fmt.Printf("Number of ip-address is: %d \n", len(ips))
			if verbose {
				fmt.Printf("Ip-address are: %v \n", ports)
			}
			return
		}

		if len(ports) < 1 {
			fmt.Printf("Number of ports is: %d \n", len(ports))
			if verbose {
				fmt.Printf("Ports are: %v \n", ports)
			}
			return
		}

		t := target.New()
		t.Scan()
		disks := t.Disks()
		for i := range disks {
			id = i
			id++
			// name := fmt.Sprintf("sub_%s_%s", h, (*disks)[i].Name())
			name := fmt.Sprintf("%s_%s", prefix, disks[i].Name())

			if dryRun {
				fmt.Printf("dryRun CreateSubsystem mkdir %s %s\n", name, disks[i].Name())
				goto DRY0
			}

			// fmt.Printf("mkdir %s %s\n", name, disks[i].Name())

			if err := target.CreateSubsystem(name, offload); err != nil {
				log.Println(err)
				return
			}

			if err := target.CreateNameSpace(name, id, disks[i].Path(), target.NewUUID(), target.ByteOne); err != nil {
				log.Println(err)
				return
			}

		DRY0:
			for i := range ips {
				port := strconv.Itoa(ports[i])

				if dryRun {
					fmt.Printf("dryRun mkdir link %s %s %s\n", name, ips[i], port)
					fmt.Printf("dryRun link %s %s %s\n", name, ips[i], port)
					goto DRY1
				}

				if err := target.CreateLinkPort(name, ips[i], port); err != nil {
					log.Println(err)
					return
				}

			DRY1:
				ports[i]++
			}
		}
	},
}