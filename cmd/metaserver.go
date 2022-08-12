/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// metaserverCmd represents the metaserver command
var metaserverCmd = &cobra.Command{
	Use:   "metaserver",
	Short: "Start meta Server.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("metaserver called")
	},
}

func init() {
	rootCmd.AddCommand(metaserverCmd)

	metaserverCmd.Flags().StringP("addr", "a", "", "The address of the metaserver.")
	metaserverCmd.Flags().StringP("port", "p", "", "The port of the metaserver.")
	metaserverCmd.Flags().BoolP("verbose", "v", false, "If the verbose is true, the metaserver will print the debug information.")
	metaserverCmd.Flags().BoolP("daemon", "d", false, "If the daemon is true, the metaserver will run as a daemon.")
}
