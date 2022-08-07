/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// moveCmd represents the move command
var moveCmd = &cobra.Command{
	Use:     "move",
	Aliases: []string{"mv, rename"},
	Short:   "Move a file or directory in the rfs.",
	Long: `Move or rename a file or directory in the rfs.
	. For example:
		rfs move test.txt test2.txt
		rfs mv test.txt ./your/path/
		rfs mv test.txt ./your/path/test2.txt
		`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("move called")
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// moveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// moveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
