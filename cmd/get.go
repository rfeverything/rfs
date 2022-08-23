/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/rfeverything/rfs/pkg/client"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Download the file from the rfs.",
	Long: `Download the file from the rfs. For example:
		rfs get test.txt.
	`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		remotePath := args[0]
		localPath := args[1]
		c, err := client.NewRfsClient()
		if err != nil {
			log.Fatalln(err)
		}
		filename, file, err := c.GetFile(context.Background(), remotePath)
		if err != nil {
			log.Fatalln(err)
		}

		if localPath == "" {
			localPath = filename
		}
		out, err := os.Create(localPath)
		if err != nil {
			log.Fatalln(err)
		}
		defer out.Close()

		fmt.Printf("download file %s success\n", filename)
		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("get file success")
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
