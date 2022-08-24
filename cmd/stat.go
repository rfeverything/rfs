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

// statCmd represents the stat command
var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "Get the file information from the rfs.",
	Long: `Get the file information from the rfs.
	For example:
		rfs stat test.txt
		`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rmtPath := args[0]
		c, err := client.NewRfsClient()
		if err != nil {
			log.Fatalln(err)
		}
		stat, err := c.Stat(context.Background(), rmtPath)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(stat)
	},
}

func init() {
	rootCmd.AddCommand(statCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
