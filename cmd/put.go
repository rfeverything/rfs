/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Upload a file to the rfs.",
	Long: `Upload a file to the rfs. For example:
		rfs put test.txt.
		`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("put called", args)
	},
}

func init() {
	rootCmd.AddCommand(putCmd)

	putCmd.Flags().StringP("name", "n", "", "If the file name is not specified, the file name will be the file name of the path.")
}
