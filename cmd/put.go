/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"os"

	"github.com/rfeverything/rfs/internal/config"
	"github.com/rfeverything/rfs/pkg/client"
	"github.com/spf13/cobra"
)

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Upload a file to the rfs.",
	Long: `Upload a file to the rfs. For example:
		rfs put test.txt.
		`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		localPath := args[0]
		remotePath := args[1]
		c, err := client.NewRfsClient()
		if err != nil {
			log.Fatalln(err)
		}
		file, err := os.Open(localPath)
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()
		err = c.PutFile(context.Background(), remotePath, file)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("put file success")
	},
}

func init() {
	rootCmd.AddCommand(putCmd)

	putCmd.Flags().StringP("name", "n", "", "If the file name is not specified, the file name will be the file name of the path.")
	putCmd.Flags().StringP("server", "s", "", "The remote metaserver address.If the server address is not specified, the server address will be localhost.")
	putCmd.Flags().StringP("port", "p", "", "The remote metaserver port.If the port is not specified, the port will be 8080.")

	config.Global().BindPFlag("server", putCmd.Flags().Lookup("client.server"))
	config.Global().BindPFlag("port", putCmd.Flags().Lookup("client.port"))
}
