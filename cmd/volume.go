/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"net"

	"github.com/rfeverything/rfs/internal/config"
	"github.com/rfeverything/rfs/internal/logger"
	"github.com/rfeverything/rfs/internal/metrics"
	vpb "github.com/rfeverything/rfs/internal/proto/volume_server"
	server "github.com/rfeverything/rfs/internal/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// volumeserverCmd represents the volumeserver command
var volumeCmd = &cobra.Command{
	Use:     "volume",
	Aliases: []string{"vol", "volumeserver"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		lis, err := net.Listen("tcp", ":"+config.Global().GetString("volume.port"))
		if err != nil {
			logger.Global().Fatal("failed to listen", zap.Error(err))
		}
		logger.Global().Info("listening on " + config.Global().GetString("volume.port"))
		s := grpc.NewServer()
		vs, err := server.NewVolumeServer()
		if err != nil {
			logger.Global().Fatal("failed to create metaserver", zap.Error(err))
		}
		mhost := config.Global().GetString("metrics.host")
		mport := config.Global().GetString("metrics.port")
		go metrics.StartMetricsServer(mhost, mport)
		vpb.RegisterVolumeServerServer(s, vs)
		if err := s.Serve(lis); err != nil {
			logger.Global().Fatal("failed to serve", zap.Error(err))
		}
	},
}

func init() {
	rootCmd.AddCommand(volumeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// volumeserverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// volumeserverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
