/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"

	"github.com/rfeverything/rfs/pkg/client"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List the files in the current directory of rfs.",
	Long: `List the files in the current directory of rfs.
	For example:
		rfs list /
		rfs list /your/path/
		`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rmtPath := args[0]
		c, err := client.NewRfsClient()
		if err != nil {
			log.Fatalln(err)
		}
		list, err := c.List(context.Background(), rmtPath)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(list)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
