package cmd

import (
	"fmt"
	"github.com/nurislam03/agent/api"
	"github.com/nurislam03/agent/backend"
	"github.com/nurislam03/agent/config"
	"github.com/nurislam03/agent/repo/memory"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start API server",
	Long:  `Start the API server`,
	Run:   serve,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		portStr := viper.GetString("SERVER_PORT")
		lsnr, err := net.Listen("tcp", ":"+portStr)
		if err != nil {
			return fmt.Errorf("Port %s is not available", portStr)
		}
		_ = lsnr.Close()
		return nil
	},
}

func init() {
	// serveCmd.PersistentFlags().IntP("port", "p", 8080, "port on which the server will listen")
	// serveCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yml", "config file")
	// viper.BindPFlag("port", serveCmd.PersistentFlags().Lookup("port"))
	RootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	cfg := config.NewConfig()
	tsk := memory.NewTaskStore()
	obj := memory.NewObjectStore()

	api := api.NewAPI(cfg, tsk, obj)

	backend.NewServer(cfg, api).Serve()
}
