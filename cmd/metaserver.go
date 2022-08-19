/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"net"

	"github.com/rfeverything/rfs/internal/logger"
	mpb "github.com/rfeverything/rfs/internal/proto/meta_server"
	server "github.com/rfeverything/rfs/internal/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
		lis, err := net.Listen("tcp", ":"+cmd.Flag("port").Value.String())
		if err != nil {
			logger.Global().Fatal("failed to listen", zap.Error(err))
		}
		s := grpc.NewServer()
		ms := server.NewMetaServer()
		mpb.RegisterMetaServerServer(s, ms)
		if err := s.Serve(lis); err != nil {
			logger.Global().Fatal("failed to serve", zap.Error(err))
		}
	},
}

func init() {
	rootCmd.AddCommand(metaserverCmd)

	metaserverCmd.Flags().StringP("addr", "a", "", "The address of the metaserver.")
	metaserverCmd.Flags().StringP("port", "p", "", "The port of the metaserver.")
	metaserverCmd.Flags().BoolP("verbose", "v", false, "If the verbose is true, the metaserver will print the debug information.")
	metaserverCmd.Flags().BoolP("daemon", "d", false, "If the daemon is true, the metaserver will run as a daemon.")
}
