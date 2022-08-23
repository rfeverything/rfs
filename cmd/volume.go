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

	volumeCmd.Flags().StringP("addr", "a", "", "The address of the metaserver.")
	volumeCmd.Flags().StringP("port", "p", "", "The port of the metaserver.")
	volumeCmd.Flags().BoolP("verbose", "v", false, "If the verbose is true, the metaserver will print the debug information.")
	volumeCmd.Flags().BoolP("daemon", "d", false, "If the daemon is true, the metaserver will run as a daemon.")

	config.Global().BindPFlag("addr", volumeCmd.Flags().Lookup("addr"))
	config.Global().BindPFlag("volume.port", volumeCmd.Flags().Lookup("port"))
}
