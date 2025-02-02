package cli

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "yc-traefik-cli",
	Short: "CLI for Yandex.Cloud Traefik HTTP Provider",
	Run: func(cmd *cobra.Command, args []string) {
		if dsn == "" {
			panic("dsn is required")
		}
	},
}

var dsn string
var saFile string
var token string

func init() {
	RootCmd.PersistentFlags().StringVar(&dsn, "ydb", "", "ydb connection string")
	RootCmd.PersistentFlags().StringVar(&saFile, "sa-file", "", "service account key file")
}
